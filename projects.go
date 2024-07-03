package main

import (
	"fmt"
	"os"
)

// Domain represents a knowledge domain or area of expertise
type Domain struct {
	ID          uint
	Name        string
	Description string
}

type ToolMemory struct {
	Enabled bool `yaml:"enabled"`
	TopN    int  `yaml:"top_n"`
}

type ToolWebGet struct {
	Enabled bool `yaml:"enabled"`
}

type ToolWebSearch struct {
	Enabled  bool   `yaml:"enabled"`
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	TopN     int    `yaml:"top_n"`
}

type ToolImgGen struct {
	Enabled bool `yaml:"enabled"`
}

type Tools struct {
	Memory    ToolMemory    `yaml:"memory"`
	WebGet    ToolWebGet    `yaml:"webget"`
	WebSearch ToolWebSearch `yaml:"websearch"`
	ImgGen    ToolImgGen    `yaml:"img_gen"`
}

// Role defines a template for an assistant's behavior
type Role struct {
	ID           uint
	Name         string
	Instructions string
	DomainID     uint   // Foreign key
	Domain       Domain // `foreignKey:DomainID`
}

// LLMParams represents the configuration parameters for an LLM
type LLMParams struct {
	Model       string
	Temperature float64
	MaxTokens   int
}

// Assistant represents an LLM with its configuration and role
type Assistant struct {
	ID     uint
	Name   string
	RoleID uint
	Role   Role      // `foreignKey:RoleID`
	Params LLMParams `embedded`
	TeamID uint      // Foreign key
}

// Team is a collection of assistants
type Team struct {
	ID         uint
	Name       string
	Assistants []Assistant // `foreignKey:TeamID`
}

// Workflow represents the sequence and connections between assistants
type Workflow struct {
	ID          uint
	Name        string
	Description string
	Steps       []WorkflowStep // `foreignKey:WorkflowID`
}

// WorkflowStep represents a single step in a workflow
type WorkflowStep struct {
	ID            uint
	WorkflowID    uint
	AssistantID   uint
	Order         int
	InputSources  []WorkflowStepSource // `foreignKey:StepID`
	OutputTargets []WorkflowStepTarget // `foreignKey:StepID`
}

// WorkflowStepSource represents the source steps for a WorkflowStep
type WorkflowStepSource struct {
	ID       uint
	StepID   uint
	SourceID uint
}

// WorkflowStepTarget represents the target steps for a WorkflowStep
type WorkflowStepTarget struct {
	ID       uint
	StepID   uint
	TargetID uint
}

// FileType is an enum for supported file types
type FileType int

const (
	TextFile FileType = iota
	MarkdownFile
	JSONFile
	PythonFile
	GoFile
	HTMLFile
	CSSFile
	JSFile
)

// File represents a file in the project
type File struct {
	ID        uint
	Name      string
	Path      string
	Type      FileType
	Content   string
	ProjectID uint
}

// Project represents the overall configuration for a goal
type Project struct {
	ID          uint
	Name        string // `unique;not null`
	Description string
	TeamID      uint
	Team        Team
	WorkflowID  uint
	Workflow    Workflow
	Files       []File
}

// Processor is an interface for objects that can process files
type Processor interface {
	Process(file File) error
}

// Inferencer is an interface for objects that can run inference
type Inferencer interface {
	RunInference(prompt string) (string, error)
}

// AssistantManager handles CRUD operations for assistants
type AssistantManager interface {
	CreateAssistant(assistant Assistant) error
	GetAssistant(id uint) (Assistant, error)
	UpdateAssistant(assistant Assistant) error
	DeleteAssistant(id uint) error
}

// ProjectManager handles CRUD operations for projects
type ProjectManager interface {
	CreateProject(project Project) error
	GetProjectByName(name string) (Project, error)
	UpdateProject(project Project) error
	DeleteProject(id uint) error
}

// CreateProjectFolder creates the resources associated with the project.
func CreateProjectFolder(config AppConfig, p *Project) error {
	// Create the Project data in the database
	if err := sqliteDB.CreateProject(p); err != nil {
		return err
	}

	// Create the project folder
	projectPath := fmt.Sprintf("%s/projects/%s", config.DataPath, p.Name)
	return os.MkdirAll(projectPath, os.ModePerm)
}
