package mailService

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/gomail.v2"
)

func MailService(toMailer string, htmlContent string, subject string) bool {
	m := gomail.NewMessage()
	m.SetHeader("From", os.Getenv("EMAILID"))
	m.SetHeader("To", toMailer)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlContent)

	fmt.Println(os.Getenv("PASSWORD"))

	d := gomail.NewDialer("smtp.gmail.com", 465, os.Getenv("EMAILID"), os.Getenv("PASSWORD"))

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Could not send email: %v", err)
		return false
	}

	log.Println("Email sent successfully!")
	return true
}
