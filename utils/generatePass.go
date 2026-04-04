package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"math/big"
	"os"
	"sync"
	"time"

	"github.com/NhatHaoDev3324/zizone-be/factory"
)

var (
	passTmpl *template.Template
	passOnce sync.Once
)

var randReader = rand.Reader

func GeneratePassword() string {
	const length = 12
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := rand.Int(randReader, big.NewInt(int64(len(charset))))
		if err != nil {
			log.Printf("⚠️ Crypto rand failed: %v", err)
			return fmt.Sprintf("%012d", time.Now().UnixNano()%1000000000000)
		}
		b[i] = charset[n.Int64()]
	}
	return string(b)
}

func SendPassword(email string, name string) (string, error) {
	frontendURL := os.Getenv("FRONTEND_DOMAIN")
	password := GeneratePassword()

	go func() {
		subject := "Zizone - Cung cấp mật khẩu tài khoản mới"

		passOnce.Do(func() {
			var err error
			passTmpl, err = template.ParseFiles("template/sendPassword.html")
			if err != nil {
				factory.LogError("Failed to parse password email template: " + err.Error())
			}
		})

		if passTmpl == nil {
			factory.LogError("Password email template is not initialized")
			return
		}

		var body bytes.Buffer
		data := struct {
			Password string
			Name     string
			Email    string
			LoginURL string
		}{Password: password, Name: name, Email: email, LoginURL: frontendURL}

		if err := passTmpl.Execute(&body, data); err != nil {
			factory.LogError("Failed to execute password template: " + err.Error())
			return
		}

		SendAsync(subject, body.String(), []string{email})
	}()

	return password, nil
}
