package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	Email "gitlab.com/ncent/throttler/services/email"
)

func handler(ctx context.Context, event events.SimpleEmailEvent) (events.SimpleEmailDisposition, error) {
	for _, record := range event.Records {
		if !Email.IsValidIncomeEmail(record) {
			log.Printf("Stop Rule Set")
			return events.SimpleEmailDisposition{events.SimpleEmailStopRuleSet}, nil
		}
	}
	log.Printf("Continue")
	return events.SimpleEmailDisposition{events.SimpleEmailContinue}, nil
}

func main() {
	lambda.Start(handler)
}
