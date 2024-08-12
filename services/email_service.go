package services

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type EmailService struct {
	sesClient *ses.SES
}

func NewEmailService() (*EmailService, error) {
	config := utils.GetConfig()

	awsConfig := &aws.Config{
		Region: aws.String(config.AWSRegion),
	}

	// Only set static credentials if both access key and secret key are provided
	if config.AWSAccessKey != "" && config.AWSSecretKey != "" {
		utils.InfoLogger.Printf("Using AWS credentials: %s, %s", config.AWSAccessKey, config.AWSSecretKey)

		awsConfig.Credentials = credentials.NewStaticCredentials(
			config.AWSAccessKey,
			config.AWSSecretKey,
			"", // token can be left empty for non-temporary credentials
		)
	} else {
		utils.InfoLogger.Println("Using default AWS credentials")
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		return nil, err
	}

	return &EmailService{
		sesClient: ses.New(sess),
	}, nil
}

func (s *EmailService) SendMagicLinkEmail(to, magicLink string) error {
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(to)},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(generateHTMLEmail(magicLink)),
				},
				Text: &ses.Content{
					Data: aws.String("Your magic link: " + magicLink),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("Your Magic Link"),
			},
		},
		Source: aws.String(utils.GetConfig().SESFromEmail),
	}

	_, err := s.sesClient.SendEmail(input)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to send email: %v", err)
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				utils.ErrorLogger.Printf("Message rejected: %v", aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				utils.ErrorLogger.Printf("Mail from domain not verified: %v", aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				utils.ErrorLogger.Printf("Configuration set does not exist: %v", aerr.Error())
			default:
				utils.ErrorLogger.Printf("Unknown error: %v", aerr.Error())
			}
		} else {
			utils.ErrorLogger.Printf("Unknown error: %v", err.Error())
		}
	}
	return err
}

func generateHTMLEmail(magicLink string) string {
	return `
		<html>
			<body>
				<h1>Your Magic Link</h1>
				<p>Click the button below to log in:</p>
				<a href="` + magicLink + `" style="background-color: #4CAF50; border: none; color: white; padding: 15px 32px; text-align: center; text-decoration: none; display: inline-block; font-size: 16px; margin: 4px 2px; cursor: pointer;">Log In</a>
			</body>
		</html>
	`
}
