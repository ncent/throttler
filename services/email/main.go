package email

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/apex/invoke"
	"github.com/aws/aws-lambda-go/events"
	Throttler "gitlab.com/ncent/throttler/services/throttler"
)

var (
	throttler = Throttler.New()
)

func IsValidIncomeEmail(record events.SimpleEmailRecord) bool {
	var from string
	var to string
	var bcc string
	log.Printf("Headers: %+v", record.SES.Mail.Headers)
	for _, header := range record.SES.Mail.Headers {
		if header.Name == "From" {
			log.Printf("From: %+v", header.Value)
			from = getStringInBetween(header.Value, "<", ">")
		} else if header.Name == "To" {
			log.Printf("To: %+v", header.Value)
			to = header.Value
		} else if header.Name == "Bcc" {
			log.Printf("Bcc: %+v", header.Value)
			bcc = header.Value
		}
	}

	isValid := isValidEmail(from, to, bcc)

	if !isValid {
		sendThrottleNotificationEmail(from)
	}

	return isValid
}

func sendThrottleNotificationEmail(from string) error {

	emailHtml := `<h2>Daily Quota Exceeded</h2>
		<p>You have exceeded the number of emails for the day.</p>
		<p>You will be able to send more emails in about 24 hours from now.</p>

		<p>Best,</p>
		<p>KK</p>
		<p>CEO <a href="http://arber.redb.ai">Arber, an nCent Labs Application</a></p>
	`

	body, err := json.Marshal(map[string]interface{}{
		"recipient": from,
		"sender":    "no-reply@redb.ai",
		"html":      emailHtml,
		"subject":   "Arber Daily Quota Exceeded",
	})

	lambdaName, _ := os.LookupEnv("SEND_EMAIL_LAMBDA")
	err = invoke.AsyncQualifier(lambdaName, "$LATEST", Payload{Body: string(body)})
	if err != nil {
		log.Printf("Failed to send email: %+v", err)
		return err
	}

	return nil
}

func isValidEmail(from string, to string, bcc string) bool {
	var key string
	hours := 24
	var err error
	if to == "start@redb.ai" {
		key = from + "-start"
		err = throttler.LimitToNTimesByNHours(key, 3, hours)
	} else if strings.Contains(bcc, "share") {
		key = from + "-share"
		err = throttler.LimitToNTimesByNHours(key, 20, hours)
	}
	log.Printf("Key: %+v", key)
	log.Printf("Err: %+v", err)
	return err == nil
}

func getStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str, end)
	if e == -1 {
		return
	}
	return str[s:e]
}
