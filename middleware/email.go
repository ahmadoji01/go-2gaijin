package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
	"time"

	"gitlab.com/kitalabs/go-2gaijin/config"
	"gitlab.com/kitalabs/go-2gaijin/templates"
)

var htmlMIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

func SendEmailConfirmation(token string, email string, source string) {

	var confirmLink string
	if source == "mobile_web_app" {
		confirmLink = config.MobileWebAppLink + "confirm_email/" + email + "/" + token
	} else if source == "android_app" {
		confirmLink = config.AndroidAppLink + "confirm_email/token=" + token
	} else if source == "ios_app" {
		confirmLink = config.IOSAppLink + "confirm_email/token=" + token
	} else if source == "desktop_web_app" {
		confirmLink = config.DesktopWebAppLink + "confirm_email?email=" + email + "&token=" + token
	} else {
		confirmLink = config.MobileWebAppLink + "confirm_email/" + email + "/" + token
	}

	from := "2gaijin@kitalabs.com"
	pass := "4Managing2GaijinEmail2020!"
	to := email
	body := templates.EmailConfirmation(confirmLink)

	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n" +
		"From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: 2Gaijin.com - Email Confirmation Request\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

func SendPhoneConfirmation(token string, phone string, source string) {
	accountSid := "ACd93fe1eee224f1fcddd98f1149190302"
	authToken := "e79c7cd73ca3706771c74ff720db10ef"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	phoneNumber := phone
	if strings.HasPrefix(phone, "0") {
		phoneNumber = strings.TrimPrefix(phoneNumber, "0")
		phoneNumber = "+81" + phoneNumber
	}

	var confirmLink string
	if source == "mobile_web_app" {
		confirmLink = config.MobileWebAppLink + "confirm_phone/" + phone + "/" + token
	} else if source == "android_app" {
		confirmLink = config.AndroidAppLink + "confirm_phone/token=" + token
	} else if source == "ios_app" {
		confirmLink = config.IOSAppLink + "confirm_phone/token=" + token
	} else if source == "desktop_web_app" {
		confirmLink = config.DesktopWebAppLink + "confirm_phone?token=" + token
	} else {
		confirmLink = config.MobileWebAppLink + "confirm_phone/token=" + token
	}

	// Create possible message bodies
	body := "To confirm your phone, click the link below:\n" + confirmLink

	// Set up rand
	rand.Seed(time.Now().Unix())

	msgData := url.Values{}
	msgData.Set("To", phoneNumber)
	msgData.Set("From", "+12513579601")
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request and return message SID
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data)
		}
	}
}

func SendPhoneConfirmationCode(code string, phone string) {
	accountSid := "ACd93fe1eee224f1fcddd98f1149190302"
	authToken := "e79c7cd73ca3706771c74ff720db10ef"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	phoneNumber := phone
	if strings.HasPrefix(phone, "0") {
		phoneNumber = strings.TrimPrefix(phoneNumber, "0")
		phoneNumber = "+81" + phoneNumber
	}

	// Create possible message bodies
	body := "2Gaijin.com - To confirm your phone, enter the following code:\n" + code

	// Set up rand
	rand.Seed(time.Now().Unix())

	msgData := url.Values{}
	msgData.Set("To", phoneNumber)
	msgData.Set("From", "+12513579601")
	msgData.Set("Body", body)
	msgDataReader := *strings.NewReader(msgData.Encode())

	client := &http.Client{}
	req, _ := http.NewRequest("POST", urlStr, &msgDataReader)
	req.SetBasicAuth(accountSid, authToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Make HTTP POST request and return message SID
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data["sid"])
		}
	} else {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			fmt.Println(data)
		}
	}
}

func SendResetPasswordEmail(token string, email string, source string) {
	var confirmLink string
	if source == "mobile_web_app" {
		confirmLink = config.MobileWebAppLink + "update-password/" + email + "/" + token
	} else if source == "android_app" {
		confirmLink = config.AndroidAppLink + "confirm_email/token=" + token
	} else if source == "ios_app" {
		confirmLink = config.IOSAppLink + "confirm_email/token=" + token
	} else if source == "desktop_web_app" {
		confirmLink = config.DesktopWebAppLink + "?email=" + email + "&reset_password_token=" + token
	} else {
		confirmLink = config.MobileWebAppLink + "update-password/" + email + "/" + token
	}

	from := "2gaijin@kitalabs.com"
	pass := "4Managing2GaijinEmail2020!"
	to := email
	body := templates.ResetPassword(confirmLink)

	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n" +
		"From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: 2Gaijin.com - Reset Password Request\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

func SendBuyingRequestEmail(email string, source string, itemName string) {
	var redirectLink string
	if source == "mobile_web_app" {
		redirectLink = config.MobileWebAppLink + "notification/"
	} else if source == "android_app" {
		redirectLink = config.AndroidAppLink + "notification/"
	} else if source == "ios_app" {
		redirectLink = config.IOSAppLink + "notification/"
	} else if source == "desktop_web_app" {
		redirectLink = config.DesktopWebAppLink + "notification/"
	} else {
		redirectLink = config.MobileWebAppLink + "notification/"
	}

	from := "2gaijin@kitalabs.com"
	pass := "4Managing2GaijinEmail2020!"
	to := email
	body := templates.OrderNotificationEmail(redirectLink, itemName)

	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n" +
		"From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Someone requested to buy your items!\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

func SendDeliveryRequestEmail(itemName string, name string, email string, phone string, wechat string, facebook string, destination string, deliveryTime string, notes string) {
	from := "2gaijin@kitalabs.com"
	pass := "4Managing2GaijinEmail2020!"
	to := "2gaijin@kitalabs.com"
	body := templates.DeliveryEmail(itemName, name, email, phone, wechat, facebook, destination, deliveryTime, notes)

	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n" +
		"From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Someone requested for delivery!\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

func SendTicketEmail(name string, email string, message string) {
	from := "2gaijin@kitalabs.com"
	pass := "4Managing2GaijinEmail2020!"
	to := "2gaijin@kitalabs.com"
	body := templates.TicketEmail(name, email, message)

	msg := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n" +
		"From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + name + " requested a ticket!\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}
