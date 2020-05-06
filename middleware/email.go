package middleware

import "fmt"

func SendResetPasswordEmail(token string) {
	return
}

func SendEmailConfirmation(token string, email string) {
	fmt.Println("Email Confirmation Sent to " + email + ". Token: " + token)
	return
}

func SendPhoneConfirmation(token string, phone string) {
	fmt.Println("Phone Confirmation Sent to " + phone + ". Token: " + token)
	return
}
