package services

import (
	"fmt"
	"time"

	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type SonExecutor struct {
	listmonkClient *ListmonkClient
}

func NewSonExecutor(listmonkClient *ListmonkClient) *SonExecutor {
	return &SonExecutor{
		listmonkClient: listmonkClient,
	}
}

func (e *SonExecutor) ExecuteSon(son models.Son, data map[string]interface{}) {
	utils.InfoLogger.Infof("Executing Son: %s", son.ID)

	if son.Delay > 0 {
		utils.InfoLogger.Infof("Delaying execution of Son %s for %v", son.ID, time.Duration(son.Delay))
		time.Sleep(time.Duration(son.Delay))
	}

	for _, action := range son.Actions {
		err := e.executeAction(action, data)
		if err != nil {
			utils.ErrorLogger.Errorf("Error executing action for Son %s: %v", son.ID, err)
		}
	}

	utils.InfoLogger.Infof("Finished executing Son: %s", son.ID)
}

func (e *SonExecutor) executeAction(action models.Action, data map[string]interface{}) error {
	utils.InfoLogger.Infof("Executing action: %s", action.Type)

	switch action.Type {
	case models.ActionSendTransactionalEmail:
		return e.sendTransactionalEmail(action.Parameters, data)
	case models.ActionManageSubscriber:
		return e.manageSubscriber(action.Parameters, data)
	case models.ActionCreateCampaign:
		return e.createCampaign(action.Parameters, data)
	default:
		err := fmt.Errorf("unknown action type: %s", action.Type)
		utils.ErrorLogger.Errorf("%v", err)
		return err
	}
}

func (e *SonExecutor) sendTransactionalEmail(params map[string]interface{}, data map[string]interface{}) error {
	templateID := int(params["template_id"].(float64))
	subscriberEmail := data["member"].(map[string]interface{})["current"].(map[string]interface{})["email"].(string)

	utils.InfoLogger.Infof("Sending transactional email to %s using template %d", subscriberEmail, templateID)
	return e.listmonkClient.SendTransactionalEmail(templateID, subscriberEmail, data)
}

func (e *SonExecutor) manageSubscriber(params map[string]interface{}, data map[string]interface{}) error {
	email := data["member"].(map[string]interface{})["current"].(map[string]interface{})["email"].(string)

	utils.InfoLogger.Infof("Managing subscriber %s", email)
	utils.InfoLogger.Infof("Params: %v", params)
	var name string

	if member, ok := data["member"].(map[string]interface{}); ok {
		if current, ok := member["current"].(map[string]interface{}); ok {
			if n, ok := current["name"].(string); ok {
				name = n
			}
		}
	}
	status := "enabled" //params["status"].(string)
	lists := []int{}
	for _, listID := range params["lists"].([]interface{}) {
		lists = append(lists, int(listID.(float64)))
	}

	utils.InfoLogger.Infof("Managing subscriber %s with status %s", email, status)
	return e.listmonkClient.ManageSubscriber(email, name, status, lists)
}

func (e *SonExecutor) createCampaign(params map[string]interface{}, data map[string]interface{}) error {
	name := params["name"].(string)
	subject := params["subject"].(string)
	lists := []int{}
	for _, listID := range params["lists"].([]interface{}) {
		lists = append(lists, int(listID.(float64)))
	}
	templateID := int(params["template_id"].(float64))
	sendAt := params["send_at"].(string)

	utils.InfoLogger.Infof("Creating campaign %s with subject %s", name, subject)
	return e.listmonkClient.CreateCampaign(name, subject, lists, templateID, sendAt)
}
