package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/NhatHaoDev3324/zizone-be/factory"
)

var MailSvc *MailService

type EmailJob struct {
	Subject string
	Body    string
	To      []string
}

type MailService struct {
	from     string
	password string
	host     string
	port     string
	queue    chan EmailJob
	wg       sync.WaitGroup
	once     sync.Once
}

func NewMailService(workerCount int) *MailService {
	from := os.Getenv("MAIL_USER")
	pass := os.Getenv("MAIL_PASS")

	m := &MailService{
		from:     from,
		password: pass,
		host:     "smtp.gmail.com",
		port:     "587",
		queue:    make(chan EmailJob, 1000),
	}

	for i := 0; i < workerCount; i++ {
		m.wg.Add(1)
		go m.worker()
	}

	MailSvc = m
	return m
}

func (m *MailService) worker() {
	defer m.wg.Done()
	var client *smtp.Client
	var err error

	for job := range m.queue {
		if client == nil {
			client, err = m.createClient()
			if err != nil {
				factory.LogError("SMTP Reconnection failed: " + err.Error())
				time.Sleep(2 * time.Second)
				continue
			}
		}

		if err := m.send(client, job); err != nil {
			factory.LogError("Send mail failed: " + err.Error())
			client.Close()
			client = nil

			client, err = m.createClient()
			if err == nil {
				if err := m.send(client, job); err != nil {
					factory.LogError("Persistent failure sending to: " + strings.Join(job.To, ","))
				} else {
					factory.LogSuccess("Sent mail to (retry): " + strings.Join(job.To, ","))
				}
			}
		} else {
			factory.LogSuccess("Sent mail to: " + strings.Join(job.To, ","))
		}
	}

	if client != nil {
		client.Quit()
	}
}

func (m *MailService) createClient() (*smtp.Client, error) {
	addr := fmt.Sprintf("%s:%s", m.host, m.port)
	c, err := smtp.Dial(addr)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		ServerName: m.host,
	}

	if err := c.StartTLS(config); err != nil {
		c.Close()
		return nil, err
	}

	auth := smtp.PlainAuth("", m.from, m.password, m.host)
	if err := c.Auth(auth); err != nil {
		c.Close()
		return nil, err
	}

	return c, nil
}

func (m *MailService) send(c *smtp.Client, job EmailJob) error {
	if err := c.Mail(m.from); err != nil {
		return err
	}

	for _, addr := range job.To {
		if err := c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	header := make(map[string]string)
	header["From"] = "NhatHao <" + m.from + ">"
	header["To"] = strings.Join(job.To, ",")
	header["Subject"] = job.Subject
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/html; charset=\"UTF-8\""

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + job.Body

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func (m *MailService) SendAsync(subject, body string, to []string) {
	m.queue <- EmailJob{
		Subject: subject,
		Body:    body,
		To:      to,
	}
}

func SendAsync(subject, body string, to []string) {
	if MailSvc != nil {
		MailSvc.SendAsync(subject, body, to)
	} else {
		factory.LogError("MailService not initialized")
	}
}

func (m *MailService) Close() {
	m.once.Do(func() {
		close(m.queue)
		m.wg.Wait()
	})
}
