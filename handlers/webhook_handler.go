package handlers

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/services"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type WebhookHandler struct {
	sonStorage     *services.SonStorage
	executor       *services.SonExecutor
	webhookService *services.WebhookService
	webhookLogger  *services.WebhookLogger
}

func NewWebhookHandler(sonStorage *services.SonStorage, executor *services.SonExecutor, webhookService *services.WebhookService, webhookLogger *services.WebhookLogger) *WebhookHandler {
	return &WebhookHandler{
		sonStorage:     sonStorage,
		executor:       executor,
		webhookService: webhookService,
		webhookLogger:  webhookLogger,
	}
}

func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	startTime := time.Now()

	endpoint := c.Param("endpoint")
	webhook, err := h.webhookService.GetWebhookByEndpoint(endpoint)
	if err != nil {
		utils.ErrorLogger.Printf("Webhook not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Webhook not found"})
		return
	}

	// Read the request body
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to read request body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	// Create initial webhook log
	webhookLogID, err := h.webhookLogger.CreateWebhookLog(webhook.UserID, c.Request, body)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create webhook log: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to log webhook"})
		return
	}

	// Verify the webhook signature
	signature := c.GetHeader("x-ghost-signature")
	if !verifySignature(signature, body, webhook.Secret) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		h.webhookLogger.UpdateWebhookLog(webhookLogID, http.StatusUnauthorized, gin.H{"error": "Invalid signature"}, time.Since(startTime))
		return
	}

	var webhookData map[string]interface{}
	if err := json.Unmarshal(body, &webhookData); err != nil {
		utils.ErrorLogger.Errorf("Invalid webhook data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook data"})
		h.webhookLogger.UpdateWebhookLog(webhookLogID, http.StatusBadRequest, gin.H{"error": "Invalid webhook data"}, time.Since(startTime))
		return
	}

	utils.InfoLogger.Infof("Received webhook: %s", utils.PrettyPrint(webhookData))

	triggerType, err := determineTriggerType(webhookData)
	if err != nil {
		utils.ErrorLogger.Errorf("Unable to determine trigger type: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to determine trigger type"})
		h.webhookLogger.UpdateWebhookLog(webhookLogID, http.StatusBadRequest, gin.H{"error": "Unable to determine trigger type"}, time.Since(startTime))
		return
	}

	utils.InfoLogger.Infof("Determined trigger type: %s", triggerType)

	// Find and execute relevant Sons
	sons, err := h.sonStorage.List(webhook.UserID)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to list Sons: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list Sons"})
		h.webhookLogger.UpdateWebhookLog(webhookLogID, http.StatusInternalServerError, gin.H{"error": "Failed to list Sons"}, time.Since(startTime))
		return
	}

	executedCount := 0
	for _, son := range sons {
		if son.Trigger == triggerType && son.Enabled {
			go func(s models.Son) {
				utils.InfoLogger.Infof("Executing Son %s for trigger %s", s.ID, triggerType)
				h.executor.ExecuteSon(s, webhookData, webhookLogID)
			}(son)
			executedCount++
		}
	}

	utils.InfoLogger.Infof("Executed %d Sons for trigger %s", executedCount, triggerType)

	// Prepare response
	response := gin.H{"message": "Webhook processed successfully", "sons_executed": executedCount}
	c.JSON(http.StatusOK, response)

	// Update the webhook log
	duration := time.Since(startTime)
	err = h.webhookLogger.UpdateWebhookLog(webhookLogID, http.StatusOK, response, duration)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to update webhook log: %v", err)
	}
}

func (h *WebhookHandler) ReplayWebhook(c *gin.Context) {
	logID := c.Param("id")

	// Get the original webhook log
	log, err := h.webhookLogger.GetWebhookLogForReplay(logID)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to get webhook log for replay: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Webhook log not found"})
		return
	}

	// Construct the webhook URL
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	// Check for X-Forwarded-Proto header
	if forwardedProto := c.GetHeader("X-Forwarded-Proto"); forwardedProto != "" {
		scheme = forwardedProto
	}
	webhookURL := fmt.Sprintf("%s://%s%s", scheme, c.Request.Host, log.Path)
	req, err := http.NewRequest(log.Method, webhookURL, bytes.NewBuffer([]byte(log.Body)))
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to create replay request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create replay request"})
		return
	}

	// Set the original headers
	var headers map[string]string
	if err := json.Unmarshal([]byte(log.Headers), &headers); err != nil {
		utils.ErrorLogger.Errorf("Failed to unmarshal headers: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process headers"})
		return
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Add a custom header to indicate this is a replay
	req.Header.Set("X-Webhook-Replay", "true")

	// Send the request to our own webhook endpoint
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to send replay request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to replay webhook"})
		return
	}
	defer resp.Body.Close()

	// Read the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to read replay response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process replay response"})
		return
	}

	// Return the response from the webhook endpoint
	c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
}

func (h *WebhookHandler) GetWebhookInfo(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		utils.ErrorLogger.Println("User not found in context")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	currentUser := user.(*models.User)

	webhooks, err := h.webhookService.GetWebhooksByUserID(currentUser.ID)
	if err != nil {
		utils.ErrorLogger.Errorf("Failed to get webhooks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get webhooks"})
		return
	}

	if len(webhooks) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No webhook found for user"})
		return
	}

	webhooks[0].Endpoint = utils.GetConfig().FrontendURL + "/webhook/" + webhooks[0].Endpoint

	// For now, we'll just return the first webhook
	c.JSON(http.StatusOK, gin.H{"data": webhooks[0]})
}

func verifySignature(signature string, payload []byte, secret string) bool {
	// Split the signature into its components
	parts := strings.Split(signature, ", ")
	if len(parts) != 2 {
		utils.ErrorLogger.Println("Invalid signature format")
		return false
	}

	receivedSignature := strings.TrimPrefix(parts[0], "sha256=")
	timestamp := strings.TrimPrefix(parts[1], "t=")

	// Append the timestamp to the payload
	fullPayload := append(payload, []byte(timestamp)...)

	// Compute the expected signature
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(fullPayload)
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	// Compare the signatures
	if receivedSignature != expectedSignature {
		utils.ErrorLogger.Printf("Signature mismatch. Received: %s, Expected: %s", receivedSignature, expectedSignature)
		return false
	}

	return true
}

func determineTriggerType(webhookData map[string]interface{}) (models.TriggerType, error) {
	if member, ok := webhookData["member"].(map[string]interface{}); ok {
		if _, ok := member["current"].(map[string]interface{}); ok {
			if previous, ok := member["previous"].(map[string]interface{}); ok && len(previous) > 0 {
				utils.InfoLogger.Infof("Member updated. Previous data: %v", previous)
				return models.TriggerMemberUpdated, nil
			}
			utils.InfoLogger.Infof("New member created")
			return models.TriggerMemberCreated, nil
		}
		utils.InfoLogger.Infof("Member deleted. Previous data: %v", member["previous"])
		return models.TriggerMemberDeleted, nil
	}

	if post, ok := webhookData["post"].(map[string]interface{}); ok {
		if current, ok := post["current"].(map[string]interface{}); ok {
			if status, ok := current["status"].(string); ok {
				switch status {
				case "published":
					utils.InfoLogger.Infof("Post published")
					return models.TriggerPostPublished, nil
				case "scheduled":
					utils.InfoLogger.Infof("Post scheduled")
					return models.TriggerPostScheduled, nil
				default:
					utils.InfoLogger.Infof("Unhandled post status: %s", status)
				}
			}
		}
	}

	if page, ok := webhookData["page"].(map[string]interface{}); ok {
		if current, ok := page["current"].(map[string]interface{}); ok {
			if status, ok := current["status"].(string); ok && status == "published" {
				utils.InfoLogger.Infof("Page published")
				return models.TriggerPagePublished, nil
			}
		}
	}

	// Log the entire webhook data for debugging
	utils.ErrorLogger.Printf("Unable to determine trigger type. Webhook data: %+v", webhookData)

	return "", utils.NewError("UnknownTriggerType", "Unable to determine trigger type from webhook data")
}
