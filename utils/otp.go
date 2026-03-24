package utils

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/NhatHaoDev3324/goAuth/config"
	"github.com/NhatHaoDev3324/goAuth/factory"
	"github.com/redis/go-redis/v9"
)

var (
	otpTmpl *template.Template
	otpOnce sync.Once
)

const (
	OTPExpiration = 5 * time.Minute
	OTPRateLimit  = 1 * time.Minute
	OTPMaxFails   = 5
)

func GenerateOTP() string {
	const length = 6
	const charset = "0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Printf("⚠️ Crypto rand failed: %v", err)
			return fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
		}
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

func SendOTP(email string) (string, error) {
	lockKey := "otp_lock:" + email
	exists, _ := config.Redis.Exists(config.Ctx, lockKey).Result()
	if exists > 0 {
		return "", errors.New("please wait a minute before requesting a new OTP")
	}

	otp := GenerateOTP()

	key := "otp:" + email
	failKey := "otp_fail:" + email

	pipe := config.Redis.Pipeline()
	pipe.Set(config.Ctx, key, otp, OTPExpiration)
	pipe.Set(config.Ctx, lockKey, "1", OTPRateLimit)
	pipe.Set(config.Ctx, failKey, "0", OTPExpiration)
	_, err := pipe.Exec(config.Ctx)
	if err != nil {
		return "", fmt.Errorf("failed to save OTP: %v", err)
	}

	go func() {
		subject := "NhatHaoDev3324 - Xác thực tài khoản"

		otpOnce.Do(func() {
			var err error
			otpTmpl, err = template.ParseFiles("template/email.html")
			if err != nil {
				factory.LogError("Failed to parse email template: " + err.Error())
			}
		})

		if otpTmpl == nil {
			factory.LogError("Email template is not initialized")
			return
		}

		var body bytes.Buffer
		data := struct{ OTP string }{OTP: otp}
		if err := otpTmpl.Execute(&body, data); err != nil {
			factory.LogError("Failed to execute template: " + err.Error())
			return
		}

		SendAsync(subject, body.String(), []string{email})
	}()

	return otp, nil
}

func VerifyOTP(email, otp string) (bool, string, error) {
	key := "otp:" + email
	failKey := "otp_fail:" + email
	fails, _ := config.Redis.Get(config.Ctx, failKey).Int()
	if fails >= OTPMaxFails {
		return false, "too_many_attempts", errors.New("this OTP has been blocked due to too many incorrect attempts")
	}

	val, err := config.Redis.Get(config.Ctx, key).Result()
	if err == redis.Nil {
		return false, "expired", errors.New("OTP has expired or does not exist")
	} else if err != nil {
		return false, "internal_error", err
	}

	if val == otp {
		config.Redis.Del(config.Ctx, key, failKey, "otp_lock:"+email)
		return true, "success", nil
	}
	config.Redis.Incr(config.Ctx, failKey)

	remaining := OTPMaxFails - (fails + 1)
	if remaining <= 0 {
		config.Redis.Del(config.Ctx, key)
		return false, "blocked", errors.New("too many failed attempts. This OTP is now invalid")
	}

	return false, "invalid", fmt.Errorf("wrong OTP. You have %d attempts remaining", remaining)
}
