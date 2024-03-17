package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var db *gorm.DB

func TestChatSession(t *testing.T) {

	// Test creating a new chat session
	session := ChatSession{
		ChatTurns: []ChatTurn{},
	}

	err := db.Create(&session).Error
	assert.NoError(t, err)

	// Test fetching the created session
	var fetchedSession ChatSession
	err = db.First(&fetchedSession, session.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, session.ID, fetchedSession.ID)

	// Test updating the session
	newTurn := ChatTurn{
		SessionID:  session.ID,
		UserPrompt: "Hello",
	}
	session.ChatTurns = append(session.ChatTurns, newTurn)

	err = db.Save(&session).Error
	assert.NoError(t, err)

	// Test deleting the session
	err = db.Delete(&session).Error
	assert.NoError(t, err)

	// Test session was deleted
	err = db.First(&fetchedSession, session.ID).Error
	assert.Error(t, err)
}

func TestChatTurn(t *testing.T) {

	// Create test session
	session := ChatSession{}
	db.Create(&session)

	// Test creating a new chat turn
	turn := ChatTurn{
		SessionID:  session.ID,
		UserPrompt: "Hello",
	}

	err := db.Create(&turn).Error
	assert.NoError(t, err)

	// Test fetching the created turn
	var fetchedTurn ChatTurn
	err = db.First(&fetchedTurn, turn.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, turn.ID, fetchedTurn.ID)

	// Test updating the turn
	newResponse := ChatResponse{
		TurnID:  turn.ID,
		Content: "Hi there!",
	}
	turn.Responses = append(turn.Responses, newResponse)

	err = db.Save(&turn).Error
	assert.NoError(t, err)

	// Test deleting the turn
	err = db.Delete(&turn).Error
	assert.NoError(t, err)

	// Test turn was deleted
	err = db.First(&fetchedTurn, turn.ID).Error
	assert.Error(t, err)

	// Clean up
	db.Delete(&session)
}

func TestChatResponse(t *testing.T) {

	// Create test session and turn
	session := ChatSession{}
	db.Create(&session)

	turn := ChatTurn{
		SessionID:  session.ID,
		UserPrompt: "Hello",
	}
	db.Create(&turn)

	// Test creating a new response
	response := ChatResponse{
		TurnID:  turn.ID,
		Content: "Hi there!",
	}

	err := db.Create(&response).Error
	assert.NoError(t, err)

	// Test fetching the created response
	var fetchedResponse ChatResponse
	err = db.First(&fetchedResponse, response.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, response.ID, fetchedResponse.ID)

	// Test updating the response
	response.Content = "Hello!"
	err = db.Save(&response).Error
	assert.NoError(t, err)

	// Test deleting the response
	err = db.Delete(&response).Error
	assert.NoError(t, err)

	// Test response was deleted
	err = db.First(&fetchedResponse, response.ID).Error
	assert.Error(t, err)

	// Clean up
	db.Delete(&turn)
	db.Delete(&session)
}
