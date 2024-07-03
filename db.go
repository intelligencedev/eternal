package main

import (
	"database/sql"
	"errors"
	"eternal/pkg/llm"
	"eternal/pkg/sd"
	"fmt"
	"time"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	_ "modernc.org/sqlite"
)

type SQLiteDB struct {
	db *sql.DB
}

type ChatSession struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ChatTurn struct {
	ID         int64
	SessionID  int64
	UserPrompt string
}

type ChatResponse struct {
	ID        int64
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

type ModelParams struct {
	ID         int              `json:"id"`
	Name       string           `json:"name"`
	Homepage   string           `json:"homepage"`
	GGUFInfo   string           `json:"gguf,omitempty"`
	Downloads  string           `json:"downloads,omitempty"`
	Downloaded bool             `json:"downloaded"`
	Options    *llm.GGUFOptions `json:"options"`
}

type ImageModel struct {
	ID         int          `json:"id"`
	Name       string       `json:"name"`
	Homepage   string       `json:"homepage"`
	Prompt     string       `json:"prompt"`
	Downloads  string       `json:"downloads,omitempty"`
	Downloaded bool         `json:"downloaded"`
	Options    *sd.SDParams `json:"options"`
}

type SelectedModels struct {
	ID        int    `json:"id"`
	ModelName string `json:"modelName"`
	Action    string `json:"action"`
}

type Chat struct {
	ID        int64  `json:"id"`
	Prompt    string `json:"prompt"`
	Response  string `json:"response"`
	ModelName string `json:"modelName"`
}

type ProjectTool struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Enable    bool   `json:"enable"`
	ProjectID uint   `json:"projectId"`
}

type URLTracking struct {
	ID  int64  `json:"id"`
	URL string `json:"url"`
}

func NewSQLiteDB(dataPath string) (*SQLiteDB, error) {
	dbPath := fmt.Sprintf("file:%s/eternaldata.db", dataPath)
	db, err := sql.Open("libsql", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}
	return &SQLiteDB{db: db}, nil
}

func (sqldb *SQLiteDB) AutoMigrate() error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS projects (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE,
            description TEXT
        )`,
		`CREATE TABLE IF NOT EXISTS teams (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL
        )`,
		`CREATE TABLE IF NOT EXISTS roles (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            instructions TEXT
        )`,
		`CREATE TABLE IF NOT EXISTS assistants (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL,
            role_id INTEGER,
            team_id INTEGER,
            FOREIGN KEY (role_id) REFERENCES roles (id),
            FOREIGN KEY (team_id) REFERENCES teams (id)
        )`,
		`CREATE TABLE IF NOT EXISTS model_params (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE,
            homepage TEXT,
            gguf_info TEXT,
            downloads TEXT,
            downloaded BOOLEAN,
			prompt TEXT
        )`,
		`CREATE TABLE IF NOT EXISTS image_models (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE,
            homepage TEXT,
            downloads TEXT,
            downloaded BOOLEAN
        )`,
		`CREATE TABLE IF NOT EXISTS selected_models (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            model_name TEXT NOT NULL UNIQUE
        )`,
		`CREATE TABLE IF NOT EXISTS chats (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            prompt TEXT,
            response TEXT,
            model_name TEXT
        )`,
		`CREATE TABLE IF NOT EXISTS url_tracking (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            url TEXT NOT NULL UNIQUE
        )`,
	}

	for _, query := range queries {
		_, err := sqldb.db.Exec(query)
		if err != nil {
			return err
		}
	}

	return nil
}

func (sqldb *SQLiteDB) GetProjectByName(name string) (Project, error) {
	var project Project
	err := sqldb.db.QueryRow("SELECT id, name, description FROM projects WHERE name = ?", name).Scan(
		&project.ID, &project.Name, &project.Description)
	if err != nil {
		return Project{}, err
	}
	return project, nil
}

func (sqldb *SQLiteDB) GetProjects() ([]Project, error) {
	rows, err := sqldb.db.Query("SELECT * FROM projects")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		err := rows.Scan(&p.ID, &p.Name, &p.Description)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

func (sqldb *SQLiteDB) CreateProject(project *Project) error {
	_, err := sqldb.db.Exec("INSERT INTO projects (name, description) VALUES (?, ?)",
		project.Name, project.Description)
	return err
}

func (sqldb *SQLiteDB) DeleteProject(name string) error {
	_, err := sqldb.db.Exec("DELETE FROM projects WHERE name = ?", name)
	return err
}

func (sqldb *SQLiteDB) ListProjects() ([]Project, error) {
	return sqldb.GetProjects()
}

func (sqldb *SQLiteDB) Create(record interface{}) error {
	switch r := record.(type) {
	case Project:
		_, err := sqldb.db.Exec("INSERT INTO projects (name, description) VALUES (?, ?)", r.Name, r.Description)
		return err
	case ModelParams:
		_, err := sqldb.db.Exec("INSERT INTO model_params (name, homepage, gguf_info, downloads, downloaded) VALUES (?, ?, ?, ?, ?)",
			r.Name, r.Homepage, r.GGUFInfo, r.Downloads, r.Downloaded)
		return err
	case ImageModel:
		_, err := sqldb.db.Exec("INSERT INTO image_models (name, homepage, prompt, downloads, downloaded) VALUES (?, ?, ?, ?, ?)",
			r.Name, r.Homepage, r.Prompt, r.Downloads, r.Downloaded)
		return err
	default:
		return errors.New("not implemented for this type")
	}
}

func (sqldb *SQLiteDB) Find(out interface{}) error {
	switch out := out.(type) {
	case *[]Project:
		rows, err := sqldb.db.Query("SELECT id, name, description FROM projects")
		if err != nil {
			return err
		}
		defer rows.Close()
		var projects []Project
		for rows.Next() {
			var p Project
			err := rows.Scan(&p.ID, &p.Name, &p.Description)
			if err != nil {
				return err
			}
			projects = append(projects, p)
		}
		*out = projects
	case *[]ModelParams:
		rows, err := sqldb.db.Query("SELECT id, name, homepage, gguf_info, downloads, downloaded FROM model_params")
		if err != nil {
			return err
		}
		defer rows.Close()
		var models []ModelParams
		for rows.Next() {
			var m ModelParams
			err := rows.Scan(&m.ID, &m.Name, &m.Homepage, &m.GGUFInfo, &m.Downloads, &m.Downloaded)
			if err != nil {
				return err
			}
			models = append(models, m)
		}
		*out = models
	case *[]ImageModel:
		rows, err := sqldb.db.Query("SELECT id, name, homepage, prompt, downloads, downloaded FROM image_models")
		if err != nil {
			return err
		}
		defer rows.Close()
		var models []ImageModel
		for rows.Next() {
			var m ImageModel
			err := rows.Scan(&m.ID, &m.Name, &m.Homepage, &m.Prompt, &m.Downloads, &m.Downloaded)
			if err != nil {
				return err
			}
			models = append(models, m)
		}
		*out = models
	default:
		return errors.New("not implemented for this type")
	}
	return nil
}

func (sqldb *SQLiteDB) First(name string, out interface{}) error {
	switch out := out.(type) {
	case *ModelParams:
		err := sqldb.db.QueryRow("SELECT id, name, homepage, gguf_info, downloads, downloaded FROM model_params WHERE name = ?", name).Scan(
			&out.ID, &out.Name, &out.Homepage, &out.GGUFInfo, &out.Downloads, &out.Downloaded)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("record not found")
			}
			return err
		}
	case *ImageModel:
		err := sqldb.db.QueryRow("SELECT id, name, homepage, prompt, downloads, downloaded FROM image_models WHERE name = ?", name).Scan(
			&out.ID, &out.Name, &out.Homepage, &out.Prompt, &out.Downloads, &out.Downloaded)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("record not found")
			}
			return err
		}
	default:
		return errors.New("not implemented for this type")
	}
	return nil
}

func (sqldb *SQLiteDB) FindByID(id uint, out interface{}) error {
	switch out := out.(type) {
	case *Project:
		err := sqldb.db.QueryRow("SELECT id, name, description FROM projects WHERE id = ?", id).Scan(
			&out.ID, &out.Name, &out.Description)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("record not found")
			}
			return err
		}
	case *ModelParams:
		err := sqldb.db.QueryRow("SELECT id, name, homepage, gguf_info, downloads, downloaded FROM model_params WHERE id = ?", id).Scan(
			&out.ID, &out.Name, &out.Homepage, &out.GGUFInfo, &out.Downloads, &out.Downloaded)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("record not found")
			}
			return err
		}
	case *ImageModel:
		err := sqldb.db.QueryRow("SELECT id, name, homepage, prompt, downloads, downloaded FROM image_models WHERE id = ?", id).Scan(
			&out.ID, &out.Name, &out.Homepage, &out.Prompt, &out.Downloads, &out.Downloaded)
		if err != nil {
			if err == sql.ErrNoRows {
				return errors.New("record not found")
			}
			return err
		}
	default:
		return errors.New("not implemented for this type")
	}
	return nil
}

func (sqldb *SQLiteDB) UpdateByName(name string, updatedRecord interface{}) error {
	switch record := updatedRecord.(type) {
	case ModelParams:
		_, err := sqldb.db.Exec("UPDATE model_params SET homepage = ?, gguf_info = ?, downloads = ?, downloaded = ? WHERE name = ?",
			record.Homepage, record.GGUFInfo, record.Downloads, record.Downloaded, name)
		return err
	case ImageModel:
		_, err := sqldb.db.Exec("UPDATE image_models SET homepage = ?, prompt = ?, downloads = ?, downloaded = ? WHERE name = ?",
			record.Homepage, record.Prompt, record.Downloads, record.Downloaded, name)
		return err
	default:
		return errors.New("not implemented for this type")
	}
}

func (sqldb *SQLiteDB) Delete(id uint, model interface{}) error {
	switch model.(type) {
	case Project:
		_, err := sqldb.db.Exec("DELETE FROM projects WHERE id = ?", id)
		return err
	case ModelParams:
		_, err := sqldb.db.Exec("DELETE FROM model_params WHERE id = ?", id)
		return err
	case ImageModel:
		_, err := sqldb.db.Exec("DELETE FROM image_models WHERE id = ?", id)
		return err
	default:
		return errors.New("not implemented for this type")
	}
}

func LoadModelDataToDB(db *SQLiteDB, models []ModelParams) error {
	for _, model := range models {
		var existingModel ModelParams
		err := db.db.QueryRow("SELECT id, name, homepage, gguf_info, downloads, downloaded FROM model_params WHERE name = ?", model.Name).Scan(
			&existingModel.ID, &existingModel.Name, &existingModel.Homepage, &existingModel.GGUFInfo,
			&existingModel.Downloads, &existingModel.Downloaded)

		if err != nil {
			if err == sql.ErrNoRows {
				// If the model is not found, create a new one
				_, err := db.db.Exec("INSERT INTO model_params (name, homepage, gguf_info, downloads, downloaded) VALUES (?, ?, ?, ?, ?)",
					model.Name, model.Homepage, model.GGUFInfo, model.Downloads, model.Downloaded)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			// If the model exists, update it
			_, err := db.db.Exec("UPDATE model_params SET homepage = ?, gguf_info = ?, downloads = ?, downloaded = ? WHERE name = ?",
				model.Homepage, model.GGUFInfo, model.Downloads, model.Downloaded, model.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func LoadImageModelDataToDB(db *SQLiteDB, models []ImageModel) error {
	for _, model := range models {
		var existingModel ImageModel
		err := db.db.QueryRow("SELECT id, name, homepage, prompt, downloads, downloaded FROM image_models WHERE name = ?", model.Name).Scan(
			&existingModel.ID, &existingModel.Name, &existingModel.Homepage, &existingModel.Prompt,
			&existingModel.Downloads, &existingModel.Downloaded)

		if err != nil {
			if err == sql.ErrNoRows {
				// If the model is not found, create a new one
				_, err := db.db.Exec("INSERT INTO image_models (name, homepage, prompt, downloads, downloaded) VALUES (?, ?, ?, ?, ?)",
					model.Name, model.Homepage, model.Prompt, model.Downloads, model.Downloaded)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		} else {
			// If the model exists, update it
			_, err := db.db.Exec("UPDATE image_models SET homepage = ?, prompt = ?, downloads = ?, downloaded = ? WHERE name = ?",
				model.Homepage, model.Prompt, model.Downloads, model.Downloaded, model.Name)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func AddSelectedModel(db *sql.DB, modelName string) error {
	// Remove any existing selected model from the database
	_, err := db.Exec("DELETE FROM selected_models")
	if err != nil {
		return err
	}

	// Create a new selected model
	_, err = db.Exec("INSERT INTO selected_models (model_name) VALUES (?)", modelName)
	return err
}

func RemoveSelectedModel(db *sql.DB, modelName string) error {
	_, err := db.Exec("DELETE FROM selected_models WHERE model_name = ?", modelName)
	return err
}

func GetSelectedModels(db *sql.DB) ([]SelectedModels, error) {
	rows, err := db.Query("SELECT id, model_name FROM selected_models")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var selectedModels []SelectedModels
	for rows.Next() {
		var sm SelectedModels
		if err := rows.Scan(&sm.ID, &sm.ModelName); err != nil {
			return nil, err
		}
		selectedModels = append(selectedModels, sm)
	}
	return selectedModels, nil
}

func CreateChat(db *sql.DB, prompt, response, model string) (Chat, error) {
	chat := Chat{Prompt: prompt, Response: response, ModelName: model}
	result, err := db.Exec("INSERT INTO chats (prompt, response, model_name) VALUES (?, ?, ?)",
		chat.Prompt, chat.Response, chat.ModelName)
	if err != nil {
		return chat, err
	}
	id, err := result.LastInsertId()
	chat.ID = id
	return chat, err
}

func GetChats(db *sql.DB) ([]Chat, error) {
	rows, err := db.Query("SELECT id, prompt, response, model_name FROM chats")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []Chat
	for rows.Next() {
		var c Chat
		if err := rows.Scan(&c.ID, &c.Prompt, &c.Response, &c.ModelName); err != nil {
			return nil, err
		}
		chats = append(chats, c)
	}
	return chats, nil
}

func GetChatByID(db *sql.DB, id int64) (Chat, error) {
	var chat Chat
	err := db.QueryRow("SELECT id, prompt, response, model_name FROM chats WHERE id = ?", id).Scan(
		&chat.ID, &chat.Prompt, &chat.Response, &chat.ModelName)
	return chat, err
}

func UpdateChat(db *sql.DB, id int64, newPrompt, newResponse, newModel string) error {
	_, err := db.Exec("UPDATE chats SET prompt = ?, response = ?, model_name = ? WHERE id = ?",
		newPrompt, newResponse, newModel, id)
	return err
}

func DeleteChat(db *sql.DB, id int64) error {
	_, err := db.Exec("DELETE FROM chats WHERE id = ?", id)
	return err
}

func (sqldb *SQLiteDB) CreateURLTracking(url string) error {
	var count int
	err := sqldb.db.QueryRow("SELECT COUNT(*) FROM url_tracking WHERE url = ?", url).Scan(&count)
	if err != nil {
		return err
	}

	if count == 0 {
		_, err = sqldb.db.Exec("INSERT INTO url_tracking (url) VALUES (?)", url)
		return err
	}

	return nil
}

func (sqldb *SQLiteDB) ListURLTrackings() ([]URLTracking, error) {
	rows, err := sqldb.db.Query("SELECT * FROM url_tracking")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var urlTrackings []URLTracking
	for rows.Next() {
		var ut URLTracking
		err := rows.Scan(&ut.ID, &ut.URL)
		if err != nil {
			return nil, err
		}
		urlTrackings = append(urlTrackings, ut)
	}
	return urlTrackings, nil
}

func (sqldb *SQLiteDB) DeleteURLTracking(url string) error {
	_, err := sqldb.db.Exec("DELETE FROM url_tracking WHERE url = ?", url)
	return err
}

func CreateDevTeam(sqldb *SQLiteDB) error {
	tx, err := sqldb.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO teams (name) VALUES (?)", "Dev Team")
	if err != nil {
		return err
	}

	var teamID int64
	err = tx.QueryRow("SELECT last_insert_rowid()").Scan(&teamID)
	if err != nil {
		return err
	}

	assistants := []struct {
		name, roleName, instructions string
	}{
		{"Senior Developer", "Software Developer", "Help with software development tasks"},
		{"Code Reviewer", "Code Reviewer", "Review code changes and provide feedback"},
	}

	for _, a := range assistants {
		_, err = tx.Exec("INSERT INTO roles (name, instructions) VALUES (?, ?)", a.roleName, a.instructions)
		if err != nil {
			return err
		}

		var roleID int64
		err = tx.QueryRow("SELECT last_insert_rowid()").Scan(&roleID)
		if err != nil {
			return err
		}

		_, err = tx.Exec("INSERT INTO assistants (name, role_id, team_id) VALUES (?, ?, ?)", a.name, roleID, teamID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
