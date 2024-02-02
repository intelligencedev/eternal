package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T, models ...interface{}) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	for _, model := range models {
		err := db.AutoMigrate(model)
		assert.NoError(t, err)
	}

	return db
}

func TestCreateAndGetRecord(t *testing.T) {
	db := setupTestDB(t, &BackendHost{}, &Chat{})

	// Test CreateRecord
	backendHost := BackendHost{Host: "http://localhost:8080", ModelType: "TestModel"}
	err := CreateRecord(db, &backendHost)
	assert.NoError(t, err)

	// Test GetRecordByID
	var fetchedBackendHost BackendHost
	err = GetRecordByID(db, backendHost.ID, &fetchedBackendHost)
	assert.NoError(t, err)
	assert.Equal(t, backendHost.Host, fetchedBackendHost.Host)
	assert.Equal(t, backendHost.ModelType, fetchedBackendHost.ModelType)
}

func TestUpdateRecord(t *testing.T) {
	db := setupTestDB(t, &BackendHost{})

	// Setup test record
	backendHost := BackendHost{Host: "http://localhost:8080", ModelType: "TestModel"}
	db.Create(&backendHost)

	// Test UpdateRecord
	backendHost.Host = "http://localhost:9090"
	err := UpdateRecord(db, backendHost.ID, &backendHost)
	assert.NoError(t, err)

	// Verify update
	var updatedBackendHost BackendHost
	db.First(&updatedBackendHost, backendHost.ID)
	assert.Equal(t, "http://localhost:9090", updatedBackendHost.Host)
}

func TestDeleteRecord(t *testing.T) {
	db := setupTestDB(t, &BackendHost{})

	// Setup test record
	backendHost := BackendHost{Host: "http://localhost:8080", ModelType: "TestModel"}
	db.Create(&backendHost)

	// Test DeleteRecord
	err := DeleteRecord(db, backendHost.ID, &BackendHost{})
	assert.NoError(t, err)

	// Verify deletion
	var deletedBackendHost BackendHost
	result := db.First(&deletedBackendHost, backendHost.ID)
	assert.Error(t, result.Error)
}
