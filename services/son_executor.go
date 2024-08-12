// services/son_executor.go
package services

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"time"

	"github.com/hibiken/asynq"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

const (
	TypeSendTransactionalEmail = "send_transactional_email"
	TypeManageSubscriber       = "manage_subscriber"
	TypeCreateCampaign         = "create_campaign"
)

type SonExecutor struct {
	listmonkClient  *ListmonkClient
	asyncClient     *asynq.Client
	asyncServer     *asynq.Server
	executionLogger *SonExecutionLogger
}

func NewSonExecutor(listmonkClient *ListmonkClient, redisAddr string, executionLogger *SonExecutionLogger) (*SonExecutor, error) {
	asyncClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	asyncServer := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)

	return &SonExecutor{
		listmonkClient:  listmonkClient,
		asyncClient:     asyncClient,
		asyncServer:     asyncServer,
		executionLogger: executionLogger,
	}, nil
}

func (e *SonExecutor) Start() error {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeSendTransactionalEmail, e.handleSendTransactionalEmail)
	mux.HandleFunc(TypeManageSubscriber, e.handleManageSubscriber)
	mux.HandleFunc(TypeCreateCampaign, e.handleCreateCampaign)

	return e.asyncServer.Start(mux)
}

func (e *SonExecutor) Stop() {
	e.asyncServer.Shutdown()
	e.asyncClient.Close()
}

func (e *SonExecutor) ExecuteSon(son models.Son, data map[string]interface{}, webhookLogID string) {
	executionID, err := e.executionLogger.LogSonExecution(son.ID, webhookLogID, "success", "")
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to log son execution: %v", err)
		return
	}

	for _, action := range son.Actions {
		payload, err := json.Marshal(map[string]interface{}{
			"action":       action,
			"data":         data,
			"execution_id": executionID,
		})
		if err != nil {
			utils.ErrorLogger.Errorf("Failed to marshal action payload: %v", err)
			e.executionLogger.LogActionExecution(executionID, string(action.Type), "failure", err.Error())
			continue
		}

		var task *asynq.Task
		switch action.Type {
		case models.ActionSendTransactionalEmail:
			task = asynq.NewTask(TypeSendTransactionalEmail, payload)
		case models.ActionManageSubscriber:
			task = asynq.NewTask(TypeManageSubscriber, payload)
		case models.ActionCreateCampaign:
			task = asynq.NewTask(TypeCreateCampaign, payload)
		default:
			utils.ErrorLogger.Errorf("Unknown action type: %s", action.Type)
			e.executionLogger.LogActionExecution(executionID, string(action.Type), "failure", "Unknown action type")
			continue
		}

		delay, err := son.GetParsedDelay()
		if err != nil {
			utils.ErrorLogger.Errorf("Failed to parse delay: %v", err)
			delay = 0
		}

		info, err := e.asyncClient.Enqueue(task, asynq.ProcessIn(delay), asynq.MaxRetry(3), asynq.Queue("default"))
		if err != nil {
			utils.ErrorLogger.Errorf("Failed to enqueue task: %v", err)
			e.executionLogger.LogActionExecution(executionID, string(action.Type), "failure", err.Error())
		} else {
			utils.InfoLogger.Infof("Enqueued task: id=%s queue=%s", info.ID, info.Queue)
			e.executionLogger.LogActionExecution(executionID, string(action.Type), "queued", "")
		}
	}
}

func (e *SonExecutor) handleSendTransactionalEmail(ctx context.Context, t *asynq.Task) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	executionID, ok := payload["execution_id"].(string)
	if !ok {
		return fmt.Errorf("invalid execution_id in payload")
	}

	action, ok := payload["action"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid action in payload")
	}

	params, ok := action["parameters"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid parameters in action")
	}

	data, ok := payload["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid data in payload")
	}

	err := e.sendTransactionalEmail(params, data)
	if err != nil {
		e.executionLogger.LogActionExecution(executionID, "send_transactional_email", "failure", err.Error())
		return err
	}

	e.executionLogger.LogActionExecution(executionID, "send_transactional_email", "success", "")
	return nil
}

func (e *SonExecutor) handleManageSubscriber(ctx context.Context, t *asynq.Task) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	executionID, ok := payload["execution_id"].(string)
	if !ok {
		return fmt.Errorf("invalid execution_id in payload")
	}

	action, ok := payload["action"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid action in payload")
	}

	params, ok := action["parameters"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid parameters in action")
	}

	data, ok := payload["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid data in payload")
	}

	err := e.manageSubscriber(params, data)
	if err != nil {
		e.executionLogger.LogActionExecution(executionID, "manage_subscriber", "failure", err.Error())
		return err
	}

	e.executionLogger.LogActionExecution(executionID, "manage_subscriber", "success", "")
	return nil
}

func (e *SonExecutor) handleCreateCampaign(ctx context.Context, t *asynq.Task) error {
	var payload map[string]interface{}
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %v", err)
	}

	executionID, ok := payload["execution_id"].(string)
	if !ok {
		return fmt.Errorf("invalid execution_id in payload")
	}

	action, ok := payload["action"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid action in payload")
	}

	params, ok := action["parameters"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid parameters in action")
	}

	data, ok := payload["data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid data in payload")
	}

	// Parse the template
	body, ok := params["body"].(string)
	if !ok {
		return fmt.Errorf("invalid body in parameters")
	}

	postData := map[string]interface{}{
		"Title":         data["post"].(map[string]interface{})["current"].(map[string]interface{})["title"],
		"FeatureImage":  data["post"].(map[string]interface{})["current"].(map[string]interface{})["feature_image"],
		"Slug":          data["post"].(map[string]interface{})["current"].(map[string]interface{})["slug"],
		"CustomExcerpt": data["post"].(map[string]interface{})["current"].(map[string]interface{})["custom_excerpt"],
		"Html":          template.HTML(data["post"].(map[string]interface{})["current"].(map[string]interface{})["html"].(string)),
		"PlainText":     data["post"].(map[string]interface{})["current"].(map[string]interface{})["plaintext"],
		"PublishedAt":   data["post"].(map[string]interface{})["current"].(map[string]interface{})["published_at"],
	}

	parsedBody, err := utils.ParseTemplate(body, postData)
	if err != nil {
		e.executionLogger.LogActionExecution(executionID, "create_campaign", "failure", fmt.Sprintf("Failed to parse template: %v", err))
		return err
	}

	// Update the params with the parsed body
	params["body"] = parsedBody

	campaignID, err := e.createCampaign(params, data)
	if err != nil {
		e.executionLogger.LogActionExecution(executionID, "create_campaign", "failure", err.Error())
		return err
	}

	// Update the campaign status to 'scheduled'
	err = e.listmonkClient.UpdateCampaignStatus(campaignID, "scheduled")
	if err != nil {
		e.executionLogger.LogActionExecution(executionID, "update_campaign_status", "failure", err.Error())
		return err
	}

	e.executionLogger.LogActionExecution(executionID, "create_campaign", "success", "")
	return nil
}

func (e *SonExecutor) sendTransactionalEmail(params map[string]interface{}, data map[string]interface{}) error {
	templateID, ok := params["template_id"].(float64)
	if !ok {
		return fmt.Errorf("invalid or missing template_id")
	}

	subscriberEmail, err := getSubscriberEmail(data)
	if err != nil {
		return err
	}

	headers, err := getHeaders(params)
	if err != nil {
		return err
	}

	additionalData, err := getAdditionalData(params)
	if err != nil {
		return err
	}

	mergedData := mergeData(data, additionalData)

	utils.InfoLogger.Infof("Sending transactional email to %s using template %d", subscriberEmail, int(templateID))
	return e.listmonkClient.SendTransactionalEmail(int(templateID), subscriberEmail, mergedData, headers)
}

func (e *SonExecutor) manageSubscriber(params map[string]interface{}, data map[string]interface{}) error {
	member, ok := data["member"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid member data")
	}

	current, ok := member["current"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid current member data")
	}

	email, _ := current["email"].(string)
	name, _ := current["name"].(string)

	utils.InfoLogger.Infof("Managing subscriber %s", email)
	utils.InfoLogger.Infof("Params: %v", params)

	status := "enabled" // params["status"].(string)

	lists := []int{}
	if listsParam, ok := params["lists"]; ok && listsParam != nil {
		if listSlice, ok := listsParam.([]interface{}); ok {
			for _, listID := range listSlice {
				if id, ok := listID.(float64); ok {
					lists = append(lists, int(id))
				}
			}
		}
	}

	var geoLocation map[string]interface{}
	if geoStr, ok := current["geolocation"].(string); ok {
		err := json.Unmarshal([]byte(geoStr), &geoLocation)
		if err != nil {
			utils.ErrorLogger.Errorf("Error parsing geolocation data: %v", err)
		}
	}

	attributes := make(map[string]interface{})
	if geoLocation != nil {
		attributes["city"] = geoLocation["city"]
		attributes["country"] = geoLocation["country"]
		attributes["latitude"] = geoLocation["latitude"]
		attributes["longitude"] = geoLocation["longitude"]
		attributes["timezone"] = geoLocation["timezone"]
	}

	utils.InfoLogger.Infof("Managing subscriber %s with status %s, lists %v, and attributes %v", email, status, lists, attributes)
	return e.listmonkClient.ManageSubscriber(email, name, status, lists, attributes)
}

func (e *SonExecutor) createCampaign(params map[string]interface{}, data map[string]interface{}) (int, error) {
	// Add error checking for each parameter
	name, ok := params["name"].(string)
	if !ok {
		return 0, fmt.Errorf("invalid or missing name parameter")
	}

	// Append a timestamp and random string to ensure uniqueness
	uniqueSuffix := fmt.Sprintf("_%s_%s", time.Now().Format("20060102_150405"), utils.GenerateRandomString(5))
	uniqueName := name + uniqueSuffix

	subject, ok := params["subject"].(string)
	if !ok {
		return 0, fmt.Errorf("invalid or missing subject parameter")
	}
	lists, ok := params["lists"].([]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid or missing lists parameter")
	}
	templateID, ok := params["template_id"].(float64)
	if !ok {
		return 0, fmt.Errorf("invalid or missing template_id parameter")
	}
	sendAt, ok := params["send_at"].(string)
	if !ok || sendAt == "" {
		// If send_at is not provided or is empty, set it to 5 minutes from now
		sendAt = time.Now().UTC().Add(5 * time.Minute).Format(time.RFC3339)
	} else {
		// If send_at is provided, ensure it's in the future
		scheduledTime, err := time.Parse(time.RFC3339, sendAt)
		if err != nil {
			return 0, fmt.Errorf("invalid send_at time format: %v", err)
		}
		if scheduledTime.Before(time.Now().UTC()) {
			// If the provided time is in the past, set it to 5 minutes from now
			sendAt = time.Now().UTC().Add(5 * time.Minute).Format(time.RFC3339)
		}
	}

	utils.InfoLogger.Infof("Scheduling campaign %s for %s", uniqueName, sendAt)
	body, ok := params["body"].(string)
	if !ok {
		return 0, fmt.Errorf("invalid or missing body parameter")
	}

	contentType, ok := params["content_type"].(string)
	if !ok {
		contentType = "html" // Default to HTML if not provided
	}

	// Convert lists to []int
	listIDs := make([]int, len(lists))
	for i, v := range lists {
		listID, ok := v.(float64)
		if !ok {
			return 0, fmt.Errorf("invalid list ID at index %d", i)
		}
		listIDs[i] = int(listID)
	}

	utils.InfoLogger.Infof("Creating campaign %s with subject %s", uniqueName, subject)
	return e.listmonkClient.CreateCampaign(uniqueName, subject, listIDs, int(templateID), sendAt, body, contentType)
}

func getSubscriberEmail(data map[string]interface{}) (string, error) {
	member, ok := data["member"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid member data")
	}

	current, ok := member["current"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid current member data")
	}

	email, ok := current["email"].(string)
	if !ok {
		return "", fmt.Errorf("invalid or missing email")
	}

	return email, nil
}

func getHeaders(params map[string]interface{}) ([]map[string]string, error) {
	headersRaw, ok := params["headers"]
	if !ok {
		return nil, nil // Headers are optional
	}

	headersJSON, err := json.Marshal(headersRaw)
	if err != nil {
		return nil, fmt.Errorf("error marshalling headers: %v", err)
	}

	var headers []map[string]string
	if err := json.Unmarshal(headersJSON, &headers); err != nil {
		return nil, fmt.Errorf("error unmarshalling headers: %v", err)
	}

	return headers, nil
}

func getAdditionalData(params map[string]interface{}) (map[string]interface{}, error) {
	additionalDataRaw, ok := params["data"]
	if !ok {
		return nil, nil // Additional data is optional
	}

	additionalDataJSON, err := json.Marshal(additionalDataRaw)
	if err != nil {
		return nil, fmt.Errorf("error marshalling additional data: %v", err)
	}

	var additionalData map[string]interface{}
	if err := json.Unmarshal(additionalDataJSON, &additionalData); err != nil {
		return nil, fmt.Errorf("error unmarshalling additional data: %v", err)
	}

	return additionalData, nil
}

func mergeData(data1, data2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range data1 {
		result[k] = v
	}

	for k, v := range data2 {
		result[k] = v
	}

	return result
}
