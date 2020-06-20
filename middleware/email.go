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
)

func SendEmailConfirmation(token string, email string) {

	from := "2gaijin@kitalabs.com"
	pass := "4Managing2GaijinEmail2020!"
	to := email
	body := "You can reset the password Token using the information below:\n" + token

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Reset Password Request\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}

func SendPhoneConfirmation(token string, phone string) {
	accountSid := "ACd93fe1eee224f1fcddd98f1149190302"
	authToken := "e79c7cd73ca3706771c74ff720db10ef"
	urlStr := "https://api.twilio.com/2010-04-01/Accounts/" + accountSid + "/Messages.json"
	phoneNumber := phone
	if strings.HasPrefix(phone, "0") {
		phoneNumber = strings.TrimPrefix(phoneNumber, "0")
		phoneNumber = "+81" + phoneNumber
	}

	// Create possible message bodies
	body := "https://go.2gaijin.com/confirm_phone?phone=" + phoneNumber + "&confirm_token=" + token

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

func SendResetPasswordEmail(token string, email string) {
	from := "2gaijin@kitalabs.com"
	pass := "4Managing2GaijinEmail2020!"
	to := email
	body := "You can reset the password Token using the information below:\n" + token

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Reset Password Request\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}
