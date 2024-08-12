package services

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/troneras/ghost-listmonk-connector/models"
)

func TestSonStorage(t *testing.T) {
	// Create a temporary database for testing
	tempDBPath := "test_ghost_listmonk.db"
	storage, err := NewSonStorage(tempDBPath)
	assert.NoError(t, err)
	defer os.Remove(tempDBPath) // Clean up the temporary database file after tests

	// Test Create
	son := models.Son{
		Name:    "Test Son",
		Trigger: models.TriggerMemberCreated,
		Delay:   models.Duration(5 * time.Minute),
		Actions: []models.Action{
			{
				Type: models.ActionSendTransactionalEmail,
				Parameters: map[string]interface{}{
					"template_id": float64(1), // Use float64 to match JSON unmarshaling behavior
				},
			},
		},
	}

	err = storage.Create(&son)
	assert.NoError(t, err)
	assert.NotEmpty(t, son.ID) // Ensure an ID was generated

	// Test Get
	retrievedSon, err := storage.Get(son.ID)
	assert.NoError(t, err)
	assert.Equal(t, son.Name, retrievedSon.Name)
	assert.Equal(t, son.Trigger, retrievedSon.Trigger)
	assert.Equal(t, son.Delay, retrievedSon.Delay)
	assert.Len(t, retrievedSon.Actions, 1)
	assert.Equal(t, son.Actions[0].Type, retrievedSon.Actions[0].Type)
	assert.Equal(t, son.Actions[0].Parameters["template_id"], retrievedSon.Actions[0].Parameters["template_id"])

	// Test Update
	son.Name = "Updated Test Son"
	err = storage.Update(son)
	assert.NoError(t, err)

	updatedSon, err := storage.Get(son.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Test Son", updatedSon.Name)

	// Test List
	sons, err := storage.List()
	assert.NoError(t, err)
	assert.Len(t, sons, 1)
	assert.Equal(t, son.ID, sons[0].ID)

	// Test Delete
	err = storage.Delete(son.ID)
	assert.NoError(t, err)

	_, err = storage.Get(son.ID)
	assert.Error(t, err)
	assert.Equal(t, ErrSonNotFound, err)

	// Test error cases
	err = storage.Create(&son)
	assert.NoError(t, err)

	// Attempt to update a non-existent son (should fail)
	nonExistentSon := models.Son{ID: "non-existent", Name: "Non-existent Son"}
	err = storage.Update(nonExistentSon)
	assert.Equal(t, ErrSonNotFound, err)

	err = storage.Delete("non-existent")
	assert.Equal(t, ErrSonNotFound, err)
}
