package models

import (
	"encoding/json"
	"time"

	"github.com/troneras/ghost-listmonk-connector/utils"
)

type TriggerType string
type ActionType string

const (
	TriggerMemberCreated TriggerType = "member_created"
	TriggerMemberDeleted TriggerType = "member_deleted"
	TriggerMemberUpdated TriggerType = "member_updated"
	TriggerPagePublished TriggerType = "page_published"
	TriggerPostPublished TriggerType = "post_published"
	TriggerPostScheduled TriggerType = "post_scheduled"
)

const (
	ActionSendTransactionalEmail ActionType = "send_transactional_email"
	ActionManageSubscriber       ActionType = "manage_subscriber"
	ActionCreateCampaign         ActionType = "create_campaign"
)

type Son struct {
	ID        string      `json:"id"`
	UserID    string      `json:"user_id"`
	Name      string      `json:"name"`
	Trigger   TriggerType `json:"trigger"`
	Delay     string      `json:"delay"`
	Actions   []Action    `json:"actions"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Enabled   bool        `json:"enabled"`
}

type Action struct {
	Type       ActionType     `json:"type"`
	Parameters map[string]any `json:"parameters"`
}

//

func (s *Son) UnmarshalJSON(data []byte) error {
	type Alias Son
	aux := &struct {
		*Alias
		Delay string `json:"delay"`
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.Delay = aux.Delay
	return nil
}

func (s Son) MarshalJSON() ([]byte, error) {
	type Alias Son
	return json.Marshal(&struct {
		*Alias
		Delay string `json:"delay"`
	}{
		Alias: (*Alias)(&s),
		Delay: s.Delay,
	})
}

func (s *Son) GetParsedDelay() (time.Duration, error) {
	return utils.ParseDuration(s.Delay)
}
