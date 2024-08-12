package models

import (
	"encoding/json"
	"errors"
	"time"
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
	Name      string      `json:"name"`
	Trigger   TriggerType `json:"trigger"`
	Delay     Duration    `json:"delay"`
	Actions   []Action    `json:"actions"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}

type Action struct {
	Type       ActionType     `json:"type"`
	Parameters map[string]any `json:"parameters"`
}

// Duration is a custom type to handle time.Duration in JSON
type Duration time.Duration

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).Minutes())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value) * time.Minute)
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}
