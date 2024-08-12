package services

import (
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type Services struct {
	User               *UserService
	MagicLink          *MagicLinkService
	Email              *EmailService
	SonStorage         *SonStorage
	SonExecutor        *SonExecutor
	Webhook            *WebhookService
	ListmonkClient     *ListmonkClient
	WebhookLogger      *WebhookLogger
	SonExecutionLogger *SonExecutionLogger
	RecentActivity     *RecentActivityService
}

func NewServices(config *utils.Config) (*Services, error) {
	emailService, err := NewEmailService()
	if err != nil {
		return nil, err
	}

	webhookService := NewWebhookService()
	userService := NewUserService(webhookService)

	listmonkClient := NewListmonkClient(config)
	sonExecutionLogger := NewSonExecutionLogger(config.RedisAddr)

	recentActivity := NewRecentActivityService()

	sonExecutor, err := NewSonExecutor(listmonkClient, config.RedisAddr, sonExecutionLogger)
	if err != nil {
		return nil, err
	}

	return &Services{
		User:               userService,
		MagicLink:          NewMagicLinkService(),
		Email:              emailService,
		SonStorage:         NewSonStorage(recentActivity),
		SonExecutor:        sonExecutor,
		Webhook:            webhookService,
		ListmonkClient:     listmonkClient,
		WebhookLogger:      NewWebhookLogger(),
		SonExecutionLogger: sonExecutionLogger,
		RecentActivity:     recentActivity,
	}, nil
}
