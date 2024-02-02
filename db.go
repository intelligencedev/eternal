package main

import (
	"errors"
	"eternal/pkg/llm"
	"fmt"
	"reflect"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type SQLiteDB struct {
	db *gorm.DB
}

// TEST
type ChatSession struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
	ChatTurns []ChatTurn `gorm:"foreignKey:SessionID"`
}

type ChatTurn struct {
	ID         int64 `gorm:"primaryKey;autoIncrement"`
	SessionID  int64
	UserPrompt string
	Responses  []ChatResponse `gorm:"foreignKey:TurnID"`
}

type ChatResponse struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	TurnID    int64
	Content   string
	Model     string // Identifier for the LLM model used
	Host      SystemInfo
	CreatedAt time.Time
}

type SystemInfo struct {
	OS     string `json:"os"`
	Arch   string `json:"arch"`
	CPUs   int    `json:"cpus"`
	Memory Memory `json:"memory"`
	GPUs   []GPU  `json:"gpus"`
}

type Memory struct {
	Total int64 `json:"total"`
}

type GPU struct {
	Model              string `json:"model"`
	TotalNumberOfCores string `json:"total_number_of_cores"`
	MetalSupport       string `json:"metal_support"`
}

// END TEST

type ModelParams struct {
	ID         int              `gorm:"primaryKey;autoIncrement"`
	Name       string           `yaml:"name"`
	Homepage   string           `yaml:"homepage"`
	GGUFInfo   string           `yaml:"gguf,omitempty"`
	Downloads  string           `yaml:"downloads,omitempty"`
	Downloaded bool             `yaml:"downloaded"`
	Options    *llm.GGUFOptions `gorm:"embedded"`
}

type SelectedModels struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	ModelName string `json:"modelName"`
	Action    string `json:"action"`
}

type Chat struct {
	ID        int64 `gorm:"primaryKey;autoIncrement"`
	Prompt    string
	Response  string
	ModelName string
}

func NewSQLiteDB(dataPath string) (*SQLiteDB, error) {

	// Silence gorm logs during this step
	newLogger := logger.Default.LogMode(logger.Silent)

	dbPath := fmt.Sprintf("%s/eternaldata.db", dataPath)
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	return &SQLiteDB{db: db}, nil
}

func (sqldb *SQLiteDB) AutoMigrate(models ...interface{}) error {
	for _, model := range models {
		if err := sqldb.db.AutoMigrate(model); err != nil {
			return fmt.Errorf("error migrating schema for %v: %v", reflect.TypeOf(model), err)
		}
	}
	return nil
}

func (sqldb *SQLiteDB) Create(record interface{}) error {
	return sqldb.db.Create(record).Error
}

func (sqldb *SQLiteDB) Find(out interface{}) error {
	return sqldb.db.Find(out).Error
}

func (sqldb *SQLiteDB) First(name string, out interface{}) error {
	return sqldb.db.Where("name = ?", name).First(out).Error
}

func (sqldb *SQLiteDB) Update(id uint, updatedRecord interface{}) error {
	return sqldb.db.Model(reflect.New(reflect.TypeOf(updatedRecord).Elem()).Interface()).Where("id = ?", id).Updates(updatedRecord).Error
}

func (sqldb *SQLiteDB) Delete(id uint, model interface{}) error {
	return sqldb.db.Delete(model, id).Error
}

func LoadModelDataToDB(db *SQLiteDB, models []ModelParams) error {
	for _, model := range models {
		var existingModel ModelParams
		result := db.db.Where("name = ?", model.Name).First(&existingModel)

		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				// If the model is not found, create a new one
				if err := db.Create(&model); err != nil {
					return err
				}
			} else {
				// Other errors
				return result.Error
			}
		} else {
			// If the model exists, update it
			if err := db.db.Model(&existingModel).Updates(&model).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func AddSelectedModel(db *gorm.DB, modelName string) error {
	selectedModel := SelectedModels{
		ModelName: modelName,
	}

	// Check if the model is already selected by the user
	var count int64
	db.Model(&SelectedModels{}).Where("model_name = ?", modelName).Count(&count)
	if count == 0 {
		// If not, add it to the database
		return db.Create(&selectedModel).Error
	}

	return nil // Model is already selected, no action needed
}

func RemoveSelectedModel(db *gorm.DB, modelName string) error {
	return db.Where("model_name = ?", modelName).Delete(&SelectedModels{}).Error
}

func GetSelectedModels(db *gorm.DB) ([]SelectedModels, error) {
	var selectedModels []SelectedModels
	err := db.Find(&selectedModels).Error
	return selectedModels, err
}

// CreateChat inserts a new chat into the database.
func CreateChat(db *gorm.DB, prompt, response, model string) (Chat, error) {
	chat := Chat{Prompt: prompt, Response: response, ModelName: model}
	result := db.Create(&chat)
	return chat, result.Error
}

// GetChats retrieves all chat entries from the database.
func GetChats(db *gorm.DB) ([]Chat, error) {
	var chats []Chat
	result := db.Find(&chats)
	return chats, result.Error
}

// GetChatByID retrieves a chat by its ID.
func GetChatByID(db *gorm.DB, id int64) (Chat, error) {
	var chat Chat
	result := db.First(&chat, id)
	return chat, result.Error
}

// UpdateChat updates an existing chat entry in the database without changing its ID.
func UpdateChat(db *gorm.DB, id int64, newPrompt, newResponse, newModel string) error {
	result := db.Model(&Chat{}).Where("id = ?", id).Updates(Chat{Prompt: newPrompt, Response: newResponse, ModelName: newModel})
	return result.Error
}

// DeleteChat removes a chat entry from the database.
func DeleteChat(db *gorm.DB, id int64) error {
	result := db.Delete(&Chat{}, id)
	return result.Error
}
