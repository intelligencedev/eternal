// db.go

package main

import (
	"errors"
	"eternal/pkg/llm"
	"eternal/pkg/sd"
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

type ImageModel struct {
	ID         int          `gorm:"primaryKey;autoIncrement"`
	Name       string       `yaml:"name"`
	Homepage   string       `yaml:"homepage"`
	Prompt     string       `yaml:"prompt"`
	Downloads  string       `yaml:"downloads,omitempty"`
	Downloaded bool         `yaml:"downloaded"`
	Options    *sd.SDParams `gorm:"embedded"`
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

type Project struct {
	gorm.Model
	Name  string
	Tools []ProjectTool `gorm:"foreignKey:ProjectID"`
	Files []File        `gorm:"foreignKey:ProjectID"`
}

type ProjectTool struct {
	gorm.Model
	Name      string
	Enable    bool
	ProjectID uint // Foreign key that refers to Project
}

type File struct {
	gorm.Model
	Path      string
	Content   string
	ProjectID uint // Foreign key that refers to Project
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

// CreateProject inserts a new project into the database.
func (sqldb *SQLiteDB) CreateProject(project *Project) error {
	return sqldb.db.Create(project).Error
}

// DeleteProject removes a project from the database.
func (sqldb *SQLiteDB) DeleteProject(name string) error {
	return sqldb.db.Where("name = ?", name).Delete(&Project{}).Error
}

// ListProjects retrieves all projects from the database.
func (sqldb *SQLiteDB) ListProjects() ([]Project, error) {
	var projects []Project
	err := sqldb.db.Find(&projects).Error
	return projects, err
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

func (sqldb *SQLiteDB) FindByID(id uint, out interface{}) error {
	return sqldb.db.First(out, id).Error
}

func (sqldb *SQLiteDB) UpdateByName(name string, updatedRecord interface{}) error {
	// Assuming 'Name' is the field in your model that holds the model's name.
	// The method first finds the record by name and then applies the updates.
	return sqldb.db.Model(updatedRecord).Where("name = ?", name).Updates(updatedRecord).Error
}

func (sqldb *SQLiteDB) UpdateDownloadedByName(name string, downloaded bool) error {
	return sqldb.db.Model(&ModelParams{}).Where("name = ?", name).Update("downloaded", downloaded).Error
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

func LoadImageModelDataToDB(db *SQLiteDB, models []ImageModel) error {
	for _, model := range models {
		var existingModel ImageModel
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
	// Remove any existing selected model from the database
	if err := db.Where("1 = 1").Delete(&SelectedModels{}).Error; err != nil {
		return err
	}

	// Create a new selected model
	selectedModel := SelectedModels{
		ModelName: modelName,
	}

	// Add the new selected model to the database
	return db.Create(&selectedModel).Error
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

// UpdateModelDownloadedState updates the downloaded state of a model in the database.
// func UpdateModelDownloadedState(db *gorm.DB, dataPath string, modelName string, downloaded bool) error {
// 	db, err := NewSQLiteDB(dataPath)
// 	if err != nil {
// 		return fmt.Errorf("failed to open database: %w", err)
// 	}
// 	defer db.Close()

// 	err = db.UpdateDownloadedByName(modelName, downloaded)
// 	if err != nil {
// 		return fmt.Errorf("failed to update model downloaded state: %w", err)
// 	}

// 	return nil
// }
