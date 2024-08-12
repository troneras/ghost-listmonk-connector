package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/troneras/ghost-listmonk-connector/models"
	"github.com/troneras/ghost-listmonk-connector/services"
	"github.com/troneras/ghost-listmonk-connector/utils"
)

type WebhookHandler struct {
	sonStorage *services.SonStorage
	executor   *services.SonExecutor
}

func NewWebhookHandler(sonStorage *services.SonStorage, executor *services.SonExecutor) *WebhookHandler {
	return &WebhookHandler{
		sonStorage: sonStorage,
		executor:   executor,
	}
}

func (h *WebhookHandler) HandleWebhook(c *gin.Context) {
	var webhookData map[string]interface{}
	if err := c.ShouldBindJSON(&webhookData); err != nil {
		utils.ErrorLogger.Errorf("Invalid webhook data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid webhook data"})
		return
	}

	utils.InfoLogger.Infof("Received webhook: %s", utils.PrettyPrint(webhookData))

	triggerType, err := determineTriggerType(webhookData)
	if err != nil {
		utils.ErrorLogger.Errorf("Unable to determine trigger type: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to determine trigger type"})
		return
	}

	utils.InfoLogger.Infof("Determined trigger type: %s", triggerType)

	// Find and execute relevant Sons
	sons, err := h.sonStorage.List()

	if err != nil {
		utils.ErrorLogger.Errorf("Failed to list Sons: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list Sons"})
		return
	}

	executedCount := 0
	for _, son := range sons {
		if son.Trigger == triggerType {
			go func(s models.Son) {
				utils.InfoLogger.Infof("Executing Son %s for trigger %s", s.ID, triggerType)
				h.executor.ExecuteSon(s, webhookData)
			}(son)
			executedCount++
		}
	}

	utils.InfoLogger.Infof("Executed %d Sons for trigger %s", executedCount, triggerType)
	c.JSON(http.StatusOK, gin.H{"message": "Webhook processed successfully", "sons_executed": executedCount})
}

func determineTriggerType(webhookData map[string]interface{}) (models.TriggerType, error) {
	if member, ok := webhookData["member"].(map[string]interface{}); ok {
		if _, ok := member["current"]; ok {
			if previous, ok := member["previous"].(map[string]interface{}); ok && len(previous) > 0 {
				// log member["previous"]
				utils.InfoLogger.Infof("Previous member data: %v", previous)
				return models.TriggerMemberUpdated, nil
			}
			return models.TriggerMemberCreated, nil
		}
		// Assuming member deletion doesn't have 'current' field
		return models.TriggerMemberDeleted, nil
	}

	if post, ok := webhookData["post"].(map[string]interface{}); ok {
		if status, ok := post["status"].(string); ok && status == "published" {
			return models.TriggerPostPublished, nil
		}
		if status, ok := post["status"].(string); ok && status == "scheduled" {
			return models.TriggerPostScheduled, nil
		}
	}

	if page, ok := webhookData["page"].(map[string]interface{}); ok {
		if status, ok := page["status"].(string); ok && status == "published" {
			return models.TriggerPagePublished, nil
		}
	}

	return "", utils.NewError("UnknownTriggerType", "Unable to determine trigger type from webhook data")
}
