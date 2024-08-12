package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/troneras/ghost-listmonk-connector/utils"
)

type ListmonkList struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ListmonkTemplate struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ListmonkClient struct {
	baseURL string
	client  *http.Client
}

func NewListmonkClient(config *utils.Config) *ListmonkClient {
	return &ListmonkClient{
		baseURL: config.ListmonkURL,
		client:  &http.Client{},
	}
}

func (c *ListmonkClient) GetLists() ([]ListmonkList, error) {
	resp, err := c.client.Get(c.baseURL + "/api/lists?page=1&per_page=100")
	if err != nil {
		return nil, fmt.Errorf("error fetching lists: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Data struct {
			Results []ListmonkList `json:"results"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return result.Data.Results, nil
}

func (c *ListmonkClient) GetTemplates() ([]ListmonkTemplate, error) {
	resp, err := c.client.Get(c.baseURL + "/api/templates?page=1&per_page=100")
	if err != nil {
		return nil, fmt.Errorf("error fetching templates: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Data []ListmonkTemplate `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	return result.Data, nil
}

func (c *ListmonkClient) SendTransactionalEmail(templateID int, subscriberEmail string, data map[string]interface{}) error {
	payload := map[string]interface{}{
		"subscriber_email": subscriberEmail,
		"template_id":      templateID,
		"data":             data,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to marshal payload: %v", err)
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := c.client.Post(c.baseURL+"/api/tx", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to send transactional email: %v", err)
		return fmt.Errorf("failed to send transactional email: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		utils.ErrorLogger.Errorf("Unexpected status code: %d, body: %s", resp.StatusCode, string(body))
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	utils.InfoLogger.Infof("Sent transactional email to %s using template %d", subscriberEmail, templateID)
	return nil
}

func (c *ListmonkClient) ManageSubscriber(email string, name string, status string, lists []int) error {
	payload := map[string]interface{}{
		"email":                    email,
		"name":                     name,
		"status":                   status,
		"lists":                    lists,
		"preconfirm_subscriptions": true,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to marshal payload: %v", err)
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// log the payload
	utils.InfoLogger.Infof("Payload: %s", string(jsonPayload))

	resp, err := c.client.Post(c.baseURL+"/api/subscribers", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to manage subscriber: %v", err)
		return fmt.Errorf("failed to manage subscriber: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		utils.ErrorLogger.Errorf("Unexpected status code: %d, body: %s", resp.StatusCode, string(body))
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	utils.InfoLogger.Infof("Managed subscriber %s with status %s", email, status)
	return nil
}

func (c *ListmonkClient) CreateCampaign(name string, subject string, lists []int, templateID int, sendAt string) error {
	payload := map[string]interface{}{
		"name":        name,
		"subject":     subject,
		"lists":       lists,
		"template_id": templateID,
		"send_at":     sendAt,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to marshal payload: %v", err)
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := c.client.Post(c.baseURL+"/api/campaigns", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to create campaign: %v", err)
		return fmt.Errorf("failed to create campaign: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		utils.ErrorLogger.Errorf("Unexpected status code: %d, body: %s", resp.StatusCode, string(body))
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	utils.InfoLogger.Infof("Created campaign %s with subject %s", name, subject)
	return nil
}
