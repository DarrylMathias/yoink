package resend

import (
	"yoink/utils/env"

	"github.com/resend/resend-go/v3"
)

var ResendClient *resend.Client

func GetResendClient() {
    apiKey := env.EnvValue.ResendAPIKey
    client := resend.NewClient(apiKey)
	ResendClient = client
}

func SendEmail(text string, subject string) (error, string){
	params := &resend.SendEmailRequest{
        To:      []string{"darrylnevmat@gmail.com"},
        From:    "help@darrylmathias.tech",
        Text:    text,
        Subject: subject,
    }

    sent, err := ResendClient.Emails.Send(params)
    if err != nil {
        return err, ""
    }
    return nil, sent.Id
}