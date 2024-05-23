Repository Documentation
This document provides a comprehensive overview of the repository's structure and contents. The first section, titled 'Directory/File Tree', displays the repository's hierarchy in a tree format. In this section, directories and files are listed using tree branches to indicate their structure and relationships. Following the tree representation, the 'File Content' section details the contents of each file in the repository. Each file's content is introduced with a '[File Begins]' marker followed by the file's relative path, and the content is displayed verbatim. The end of each file's content is marked with a '[File Ends]' marker. This format ensures a clear and orderly presentation of both the structure and the detailed contents of the repository.
Directory/File Tree Begins -->
eternal/
├── config.go
├── config.yml
├── config_test.go
├── db.go
├── db_test.go
├── docs
├── examples
├── handlers.go
├── host.go
├── index
├── main.go
├── main_test.go
├── pkg
│   ├── documents
│   │   ├── gitloader.go
│   │   └── txtsplitter.go
│   ├── embeddings
│   │   ├── local.go
│   │   └── oai.go
│   ├── hfutils
│   │   └── hfutils.go
│   ├── jobs
│   │   └── jobs.go
│   ├── llm
│   │   ├── anthropic
│   │   │   ├── anthropic.go
│   │   │   └── completions.go
│   │   ├── gguf.go
│   │   ├── google
│   │   │   └── google.go
│   │   ├── llm.go
│   │   ├── local
│   │   │   ├── bin
│   │   │   └── gguf
│   │   ├── openai
│   │   │   ├── completions.go
│   │   │   ├── models.go
│   │   │   └── openai.go
│   │   └── templates.go
│   ├── sd
│   │   ├── sd.go
│   │   └── sdcpp
│   ├── search
│   │   └── search.go
│   ├── vecstore
│   │   └── vecstore.go
│   └── web
│       ├── mdtohtml.go
│       ├── serpapi.go
│       └── web.go
├── public
│   ├── css
│   │   ├── drawflow
│   │   ├── halfmoon
│   │   ├── header.css
│   │   └── styles.css
│   ├── fonts
│   ├── img
│   ├── js
│   │   ├── bootstrap
│   │   ├── cloudbox.js
│   │   ├── drawflow.min.js
│   │   ├── events.js
│   │   ├── highlight
│   │   ├── htmx.min.js
│   │   ├── node_modules
│   │   ├── package-lock.json
│   │   ├── package.json
│   │   └── workflows.js
│   ├── templates
│   │   ├── alerts.html
│   │   ├── boxheader.html
│   │   ├── chat.html
│   │   ├── drawflowdata.json
│   │   ├── flow.html
│   │   ├── header.html
│   │   ├── index.html
│   │   ├── model.html
│   │   └── shell.html
│   └── uploads
│       └── uploads_go_here
├── scripts
├── services
├── tmp
└── utils.go
<-- Directory/File Tree Ends
File Content Begins -->
[File Begins] config.go
package main

import (
	"eternal/pkg/llm"
	"eternal/pkg/sd"
	"time"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm"
)

type AppConfig struct {
	ServerID       string                            `yaml:"server_id"`
	CurrentUser    string                            `yaml:"current_user"`
	AssistantName  string                            `yaml:"assistant_name"`
	ControlHost    string                            `yaml:"control_host"`
	ControlPort    string                            `yaml:"control_port"`
	DataPath       string                            `yaml:"data_path"`
	ServiceHosts   map[string]map[string]BackendHost `yaml:"service_hosts"`
	ChromedpKey    string                            `yaml:"chromedp_key"`
	OAIKey         string                            `yaml:"oai_key"`
	AnthropicKey   string                            `yaml:"anthropic_key"`
	GoogleKey      string                            `yaml:"google_key"`
	LanguageModels []llm.Model                       `yaml:"language_models"`
	ImageModels    []sd.ImageModel                   `yaml:"image_models"`
	Tools          struct {
		Memory struct {
			Enabled bool `yaml:"enabled"`
			TopN    int  `yaml:"top_n"`
		} `yaml:"memory"`
		WebGet struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"webget"`
		WebSearch struct {
			Enabled bool `yaml:"enabled"`
			TopN    int  `yaml:"top_n"`
		} `yaml:"websearch"`
		ImgGen struct {
			Enabled bool `yaml:"enabled"`
		} `yaml:"img_gen"`
	} `yaml:"tools"`
}

type BackendHost struct {
	ID            uint           `gorm:"primaryKey" yaml:"-"`
	Host          string         `yaml:"host" gorm:"column:host"`
	Port          string         `yaml:"port" gorm:"column:port"`
	GgufGPULayers int            `yaml:"gpu_layers" gorm:"column:gguf_gpu_layers"`
	ModelType     string         `yaml:"model_type" gorm:"column:model_type"`
	CreatedAt     time.Time      `yaml:"-"`
	UpdatedAt     time.Time      `yaml:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" yaml:"-"`
}

// LoadConfig loads configuration from a YAML file.
func LoadConfig(fs afero.Fs, path string) (*AppConfig, error) {
	config := &AppConfig{}

	// Use Afero to read the file
	file, err := afero.ReadFile(fs, path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

[File Ends] config.go

[File Begins] config.yml
# The desired user display name
current_user: 'Locomod'

# Display name for LLM responses
assistant_name: 'Eternal'

# Host name and port of the frontend and management server
control_host: 'localhost'
control_port: 8080

# datapath will be created by Eternal to host persistent data
data_path: '/Users/arturoaquino/.eternal-v1'

# Service hosts for specific functionalities
service_hosts:
  # generate embeddings and local retrieval services
  retrieval:
    retrieval_host_1:
      host: 'localhost'
      port: 8080
  # image generation and processing
  image:
    image_host_1:
      host: 'localhost'
      port: 8080
  # text-to-speech and speech-to-text    
  speech:
    speech_host_1:
      host: 'localhost'
      port: 8080
  # llm text generation services
  llm:
    llm_host_1:
      host: 'localhost'
      port: '8080'
      gpu_layers: -1
    llm_host_2:
      host: 'localhost'
      port: '8080'
      gpu_layers: -1

tools:
  memory:
    enabled: false
    top_n: 2 # Number of similar chunks to retrieve from memory. Higher numbers require models with much higher context.
  webget:
    enabled: false
  websearch:
    enabled: false
    top_n: 2 # Number of page results to retrieve. Higher numbers require models with much higher context.
  imggen:
    enabled: false

# OpenAI API Key
oai_key: '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'

# Anthropic Key for Claude Completions
anthropic_key: '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'

# Google API Key for Gemini Completions
google_key: 'AIzaSyBi9ueEkCBnTxKjXIYJ9IBOvbtGA8drD1o'

language_models:
  - name: 'openai-gpt'
    homepage: 'https://platform.openai.com/docs/models/gpt-4-and-gpt-4-turbo'
    prompt: '{system}\n\n{prompt}'
    ctx: 128000
    roles:
      - 'all'
  - name: 'anthropic-claude-opus'
    homepage: 'https://www.anthropic.com/product'
    prompt: |
      Below is an instruction that describes a task. Write a response that appropriately completes the request using advanced AI capabilities.

      ### Instruction:
      {system}
      {prompt}

      ### Response:
    ctx: 128000
    roles:
      - 'all'
  - name: 'google-gemini-1.5'
    homepage: 'https://deepmind.google/technologies/gemini/#gemini-1.5'
    prompt: '{prompt}'
    ctx: 500000
    roles:
      - 'all'
  - name: 'starling-7b-beta'
    homepage: 'https://huggingface.co/HuggingFaceH4/zephyr-7b-beta'
    gguf: 'https://huggingface.co/LoneStriker/Starling-LM-7B-beta-GGUF'
    downloads:
      - 'https://huggingface.co/LoneStriker/Starling-LM-7B-beta-GGUF/resolve/main/Starling-LM-7B-beta-Q8_0.gguf'
    prompt: "Code User: {prompt}<|end_of_turn|>Code Assistant:"
    ctx: 8192
    roles:
      - 'chat'
    tags:
      - '7B'
  - name: 'llama3-8b-instruct'
    homepage: 'https://ai.meta.com/blog/meta-llama-3/'
    gguf: 'https://huggingface.co/NikolayKozloff/Meta-Llama-3-8B-Instruct-bf16-correct-pre-tokenizer-and-EOS-token-Q8_0-Q6_k-GGUF'
    downloads:
      - 'https://huggingface.co/NikolayKozloff/Meta-Llama-3-8B-Instruct-bf16-correct-pre-tokenizer-and-EOS-token-Q8_0-Q6_k-GGUF/resolve/main/Meta-Llama-3-8B-Instruct-correct-pre-tokenizer-and-EOS-token-Q8_0.gguf'
    prompt: "<|begin_of_text|><|start_header_id|>system<|end_header_id|>\n\n{system}<|eot_id|><|start_header_id|>user<|end_header_id|>\n\n{prompt}<|eot_id|><|start_header_id|>assistant<|end_header_id|>"
    ctx: 8192
    roles:
      - 'instruct'
    tags:
      - '8B'
  - name: 'llama3-8b-1048k'
    homepage: 'https://ai.meta.com/blog/meta-llama-3/'
    gguf: 'https://huggingface.co/bartowski/Llama-3-8B-Instruct-Gradient-1048k-GGUF'
    downloads:
      - 'https://huggingface.co/bartowski/Llama-3-8B-Instruct-Gradient-1048k-GGUF/resolve/main/Llama-3-8B-Instruct-Gradient-1048k-Q8_0.gguf'
    prompt: "<|begin_of_text|><|start_header_id|>system<|end_header_id|>\n\n{system}<|eot_id|><|start_header_id|>user<|end_header_id|>\n\n{prompt}<|eot_id|><|start_header_id|>assistant<|end_header_id|>"
    ctx: 120000
    roles:
      - 'instruct'
    tags:
      - '8B'
  - name: 'llama3-70b-instruct'
    homepage: 'https://ai.meta.com/blog/meta-llama-3/'
    gguf: 'https://huggingface.co/LoneStriker/Meta-Llama-3-70B-Instruct-GGUF'
    downloads:
      - 'https://huggingface.co/LoneStriker/Meta-Llama-3-70B-Instruct-GGUF/resolve/main/Meta-Llama-3-70B-Instruct-Q6_K.gguf'
    prompt: "<|begin_of_text|><|start_header_id|>system<|end_header_id|>\n\n{system}<|eot_id|><|start_header_id|>user<|end_header_id|>\n\n{prompt}<|eot_id|><|start_header_id|>assistant<|end_header_id|>"
    ctx: 40960
    roles:
      - 'instruct'
    tags:
      - '70B'
  - name: 'wizardlm-2-7b'
    homepage: 'https://huggingface.co/microsoft/WizardLM-2-7B'
    gguf: 'https://huggingface.co/MaziyarPanahi/WizardLM-2-7B-GGUF'
    downloads:
      - 'https://huggingface.co/MaziyarPanahi/WizardLM-2-7B-GGUF/resolve/main/WizardLM-2-7B.Q8_0.gguf'
    prompt: "<|im_start|>system {system}<|im_end|><|im_start|>user{prompt}<|im_end|><|im_start|>assistant"
    ctx: 32768
    roles:
      - 'all'
    tags:
      - '7b'
  - name: 'wizardlm-2-8x22b'
    homepage: 'https://huggingface.co/microsoft/WizardLM-2-7B'
    gguf: 'https://huggingface.co/MaziyarPanahi/WizardLM-2-8x22B-GGUF'
    downloads:
      - 'https://huggingface.co/MaziyarPanahi/WizardLM-2-8x22B-GGUF/resolve/main/WizardLM-2-8x22B.Q4_K_S-00001-of-00005.gguf'
    prompt: "<|im_start|>system {system}<|im_end|><|im_start|>user{prompt}<|im_end|><|im_start|>assistant"
    ctx: 32768
    roles:
      - 'all'
    tags:
      - '22x7b'
  - name: 'codeninja-7b'
    homepage: 'https://huggingface.co/beowolx/CodeNinja-1.0-OpenChat-7B'
    gguf: 'https://huggingface.co/TheBloke/CodeNinja-1.0-OpenChat-7B-GGUF'
    downloads:
      - 'https://huggingface.co/TheBloke/CodeNinja-1.0-OpenChat-7B-GGUF/resolve/main/codeninja-1.0-openchat-7b.Q8_0.gguf'
    prompt: "GPT4 Correct User: {prompt}<|end_of_turn|>GPT4 Correct Assistant:"
    ctx: 8192
    roles:
      - 'code'
    tags:
      - '7B'
  - name: 'wizardmath-7b'
    homepage: 'https://huggingface.co/WizardLM/WizardMath-7B-V1.1'
    gguf: 'https://huggingface.co/TheBloke/WizardMath-7B-V1.1-GGUF'
    downloads:
      - 'https://huggingface.co/TheBloke/WizardMath-7B-V1.1-GGUF/resolve/main/wizardmath-7b-v1.1.Q8_0.gguf'
    prompt: "Below is an instruction that describes a task. Write a response that appropriately completes the request.\n\n### Instruction:\n{prompt}\n\n### Response:"
    ctx: 32768
    roles:
      - 'math'
      - 'logic'
    tags:
      - '7B'
  - name: 'mixtral-8x7b-instruct'
    homepage: 'https://huggingface.co/mistralai/Mixtral-8x7B-Instruct-v0.1'
    gguf: 'https://huggingface.co/TheBloke/Mixtral-8x7B-Instruct-v0.1-GGUF'
    downloads:
      - 'https://huggingface.co/TheBloke/Mixtral-8x7B-Instruct-v0.1-GGUF/resolve/main/mixtral-8x7b-instruct-v0.1.Q8_0.gguf'
    prompt: "</s>[INST] {system} {prompt} [/INST]"
    ctx: 32768
    roles:
      - 'all'
    tags:
      - '8x7B'
  - name: 'yarn-mistral-7b-64k'
    homepage: 'https://huggingface.co/MaziyarPanahi/Yarn-Mistral-7b-64k-Mistral-7B-Instruct-v0.1'
    gguf: 'https://huggingface.co/MaziyarPanahi/Yarn-Mistral-7b-64k-Mistral-7B-Instruct-v0.1-GGUF'
    downloads:
      - 'https://huggingface.co/MaziyarPanahi/Yarn-Mistral-7b-64k-Mistral-7B-Instruct-v0.1-GGUF/resolve/main/Yarn-Mistral-7b-64k-Mistral-7B-Instruct-v0.1.Q6_K.gguf'
    prompt: "<|prompter|>{prompt}</s><|assistant|>"
    ctx: 64000
    roles:
      - 'all'
    tags:
      - '8x7b'
  - name: 'aixcoder-7b'
    homepage: 'https://huggingface.co/aiXcoder/aixcoder-7b-base'
    gguf: 'https://huggingface.co/bartowski/aixcoder-7b-base-GGUF'
    downloads:
      - 'hhttps://huggingface.co/bartowski/aixcoder-7b-base-GGUF/resolve/main/aixcoder-7b-base-Q8_0.gguf'
    prompt: "{prompt}"
    ctx: 8092
    roles:
      - 'all'
    tags:
      - '7b'
  - name: 'everyone-coder-33b'
    homepage: 'https://huggingface.co/rombodawg/Everyone-Coder-33b-Base'
    gguf: 'https://huggingface.co/TheBloke/Everyone-Coder-33B-Base-GGUF'
    downloads:
      - 'https://huggingface.co/TheBloke/Everyone-Coder-33B-Base-GGUF/resolve/main/everyone-coder-33b-base.Q8_0.gguf'
    prompt: "Below is an instruction that describes a task. Write a response that appropriately completes the request.\n### Instruction:\n{prompt}\n### Response:"
    ctx: 16384
    roles:
      - 'code'
    tags:
      - '33B'
  - name: 'miqu-1-70b'
    homepage: 'https://huggingface.co/miqudev/miqu-1-70b'
    gguf: 'https://huggingface.co/miqudev/miqu-1-70b'
    downloads:
      - 'https://huggingface.co/miqudev/miqu-1-70b/resolve/main/miqu-1-70b.q5_K_M.gguf'
    prompt:  "<s> [INST] {{prompt}} [/INST] ANSWER_1</s>"
    ctx: 32768
    roles:
      - 'all'
    tags:
      - '70B'
  - name: 'llama-8b-hermes'
    homepage: 'https://huggingface.co/NousResearch/Hermes-2-Pro-Llama-3-8B-GGUF'
    gguf: 'https://huggingface.co/NousResearch/Hermes-2-Pro-Llama-3-8B-GGUF'
    downloads:
      - 'https://huggingface.co/NousResearch/Hermes-2-Pro-Llama-3-8B-GGUF/resolve/main/Hermes-2-Pro-Llama-3-8B-Q8_0.gguf'
    prompt:  "<|im_start|>system You are a helpful assistant. You respond with I dont know the answer to that question when you do not know the answer to a question or the correct response.<|im_end|><|im_start|>user{prompt}<|im_end|><|im_start|>assistant"
    ctx: 8096
    roles:
      - 'all'
    tags:
      - '8B'
  - name: 'llama3-sfr-dpo-8b'
    homepage: 'https://huggingface.co/NikolayKozloff/SFR-Iterative-DPO-LLaMA-3-8B-R-Q8_0-GGUF'
    gguf: 'https://huggingface.co/NikolayKozloff/SFR-Iterative-DPO-LLaMA-3-8B-R-Q8_0-GGUF'
    downloads:
      - 'https://huggingface.co/NikolayKozloff/SFR-Iterative-DPO-LLaMA-3-8B-R-Q8_0-GGUF/resolve/main/sfr-iterative-dpo-llama-3-8b-r.Q8_0.gguf'
    prompt:  "<|begin_of_text|><|start_header_id|>system<|end_header_id|>\n\nYou are a knowledgeable, efficient, and direct AI assistant. Provide concise answers, focusing on the key information needed. Offer suggestions tactfully when appropriate to improve outcomes. Engage in productive collaboration with the user.<|eot_id|><|start_header_id|>user<|end_header_id|>\n\n{prompt}<|eot_id|><|start_header_id|>assistant<|end_header_id|>"
    ctx: 8192
    roles:
      - 'all'
    tags:
      - '8B'
  - name: 'yi-1.5-9b-chat'
    homepage: 'https://huggingface.co/YorkieOH10/Yi-1.5-9B-Chat-Q8_0-GGUF'
    gguf: 'https://huggingface.co/YorkieOH10/Yi-1.5-9B-Chat-Q8_0-GGUF'
    downloads:
      - 'https://huggingface.co/YorkieOH10/Yi-1.5-9B-Chat-Q8_0-GGUF/resolve/main/yi-1.5-9b-chat.Q8_0.gguf'
    prompt:  "<|startoftext|>You are a helpful, smart, kind, and efficient AI assistant. You always fulfill the user's requests to the best of your ability.<|im_end|>### Instruction: {prompt}/"
    ctx: 4096
    roles:
      - 'all'
    tags:
      - '9B'
  - name: 'yi-1.5-34b-chat'
    homepage: 'https://huggingface.co/bartowski/Yi-1.5-34B-Chat-GGUF'
    gguf: 'https://huggingface.co/bartowski/Yi-1.5-34B-Chat-GGUF'
    downloads:
      - 'https://huggingface.co/bartowski/Yi-1.5-34B-Chat-GGUF/resolve/main/Yi-1.5-34B-Chat-Q8_0.gguf'
    prompt:  "{prompt}"
    ctx: 4096
    roles:
      - 'all'
    tags:
      - '9B'

image_models:
  - name: 'dreamshaper-8-turbo-sdxl'
    homepage: 'https://huggingface.co/Lykon/dreamshaper-xl-v2-turbo/'
    downloads:
      - 'https://huggingface.co/Lykon/dreamshaper-xl-v2-turbo/resolve/main/DreamShaperXL_Turbo_V2-SFW.safetensors'

[File Ends] config.yml

[File Begins] config_test.go
package main

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig(t *testing.T) {
	fs := afero.NewOsFs()

	// Load example .config.yml from current directory
	config, err := LoadConfig(fs, ".config.yml")
	assert.NoError(t, err)

	// Assertions
	assert.Equal(t, "User", config.CurrentUser)
	assert.Equal(t, "Assistant", config.AssistantName)
	assert.Equal(t, "localhost", config.ControlHost)
	assert.Equal(t, "8080", config.ControlPort)
	assert.Equal(t, "/Users/$USER/.eternal", config.DataPath)

	// Refactored assertServiceHost function
	assertServiceHost := func(service string, hostKey string, expectedHost string, expectedPort string) {
		hostConfig, exists := config.ServiceHosts[service][hostKey]
		assert.True(t, exists)
		assert.Equal(t, expectedHost, hostConfig.Host) // Use Host instead of HostURL
		assert.Equal(t, expectedPort, hostConfig.Port) // Use Port instead of HostPort
	}

	// Updated calls to assertServiceHost with correct field names and types
	assertServiceHost("retrieval", "retrieval_host_1", "localhost", "8081")
	assertServiceHost("image", "image_host_1", "localhost", "8082")
	assertServiceHost("speech", "speech_host_1", "localhost", "8083")
	assertServiceHost("llm", "llm_host_1", "localhost", "8081")

	assert.Equal(t, "sk-...", config.OAIKey)
}

[File Ends] config_test.go

[File Begins] db.go
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

[File Ends] db.go

[File Begins] db_test.go
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

[File Ends] db_test.go

[File Begins] handlers.go
package main

[File Ends] handlers.go

[File Begins] host.go
// NOTE: Thiese functions are not implemented in the main app yet.
package main

import (
	// Importing necessary packages for executing commands, formatting strings, and hardware information retrieval.
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/jaypipes/ghw"        // Package for hardware information
	"github.com/shirou/gopsutil/mem" // Package for system memory information
)

// HostInfo struct: Stores information about the host system.
type HostInfo struct {
	OS     string `json:"os"`   // Operating System
	Arch   string `json:"arch"` // Architecture (e.g., amd64, 386)
	CPUs   int    `json:"cpus"` // Number of CPUs
	Memory struct {
		Total uint64 `json:"total"` // Total memory in bytes
	} `json:"memory"`
	GPUs []GPUInfo `json:"gpus"` // Slice of GPU information
}

// GPUInfo struct: Stores information about GPUs in the system.
type GPUInfo struct {
	Model              string `json:"model"`                 // GPU model
	TotalNumberOfCores string `json:"total_number_of_cores"` // Total cores in GPU
	MetalSupport       string `json:"metal_support"`         // Metal support (specific to macOS)
}

// GetHostInfo function: Retrieves information about the host system.
func GetHostInfo() (HostInfo, error) {
	hostInfo := HostInfo{
		OS:   runtime.GOOS,     // Fetching OS
		Arch: runtime.GOARCH,   // Fetching architecture
		CPUs: runtime.NumCPU(), // Fetching CPU count
	}

	// Retrieve memory information using gopsutil
	vmStat, _ := mem.VirtualMemory()
	hostInfo.Memory.Total = vmStat.Total

	// GPU information retrieval based on OS
	switch runtime.GOOS {
	case "darwin":
		// macOS specific GPU information retrieval
		gpus, err := getMacOSGPUInfo()
		if err != nil {
			fmt.Printf("Error getting GPU info: %v\n", err)
		} else {
			hostInfo.GPUs = append(hostInfo.GPUs, gpus)
		}

	case "linux", "windows":
		// Linux and Windows GPU information retrieval
		gpu, err := ghw.GPU()
		if err != nil {
			fmt.Printf("Error getting GPU info: %v\n", err)
		} else {
			for _, card := range gpu.GraphicsCards {
				gpuInfo := GPUInfo{
					Model: card.DeviceInfo.Product.Name, // Fetching GPU model
				}
				hostInfo.GPUs = append(hostInfo.GPUs, gpuInfo)
			}
		}
	}

	return hostInfo, nil
}

// getMacOSGPUInfo function: Retrieves GPU information for macOS.
func getMacOSGPUInfo() (GPUInfo, error) {
	cmd := exec.Command("system_profiler", "SPDisplaysDataType")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return GPUInfo{}, err
	}

	return parseGPUInfo(out.String())
}

// parseGPUInfo function: Parses the output from system_profiler to extract GPU info.
func parseGPUInfo(input string) (GPUInfo, error) {
	gpuInfo := GPUInfo{}

	for _, line := range strings.Split(input, "\n") {
		// Extracting relevant information from the output
		if strings.Contains(line, "Chipset Model") {
			gpuInfo.Model = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Total Number of Cores") {
			gpuInfo.TotalNumberOfCores = strings.TrimSpace(strings.Split(line, ":")[1])
		}
		if strings.Contains(line, "Metal") {
			gpuInfo.MetalSupport = strings.TrimSpace(strings.Split(line, ":")[1])
		}
	}

	return gpuInfo, nil
}

[File Ends] host.go

[File Begins] main.go
package main

import (
	"bufio"
	"context"
	"embed"
	"encoding/json"
	"errors"
	"eternal/pkg/embeddings"
	"eternal/pkg/hfutils"
	"eternal/pkg/llm"
	"eternal/pkg/llm/anthropic"
	"eternal/pkg/llm/google"
	"eternal/pkg/llm/openai"
	"eternal/pkg/sd"
	"eternal/pkg/web"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/blevesearch/bleve/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/afero"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

var (
	//go:embed public/* pkg/llm/local/bin/* pkg/sd/sdcpp/build/bin/*
	embedfs embed.FS

	osFS  afero.Fs = afero.NewOsFs()
	memFS afero.Fs = afero.NewMemMapFs()

	chatTurn    = 1
	sqliteDB    *SQLiteDB
	searchIndex bleve.Index
)

type WebSocketMessage struct {
	ChatMessage string                 `json:"chat_message"`
	Model       string                 `json:"model"`
	Headers     map[string]interface{} `json:"HEADERS"`
}

type Tool struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

func main() {
	_ = pterm.DefaultBigText.WithLetters(putils.LettersFromString("ETERNAL")).Render()

	// LOG SETTINGS
	//log.SetOutput(io.Discard)

	// Log configuration
	//log.SetLevel(log.LevelDebug)

	//zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// TODO: Check if external dependencies are installed and if not, install them
	// Such as Chromium, Docker, etc. For now, only Chromium is required for the web tool.

	// CONFIG
	currentPath, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current path: %v", err)
	}

	configPath := filepath.Join(currentPath, "config.yml")

	pterm.Info.Println("Loading config:", configPath)

	config, err := LoadConfig(osFS, configPath)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
		os.Exit(1)
	}

	// Populate tools based on the configuration
	var tools []Tool

	// Print the tools enabled in the config

	if config.Tools.WebGet.Enabled {
		tools = append(tools, Tool{Name: "webget", Enabled: true})
	}
	if config.Tools.WebSearch.Enabled {
		tools = append(tools, Tool{Name: "websearch", Enabled: true})
	}

	pterm.Info.Sprintf("GPU Layers: %s\n", config.ServiceHosts["llm"]["llm_host_1"].GgufGPULayers)

	if _, err := os.Stat(config.DataPath); os.IsNotExist(err) {
		err = os.Mkdir(config.DataPath, 0755)
		if err != nil {
			pterm.Error.Println("Error creating data directory:", err)
			os.Exit(1)
		}
	}

	_, err = InitServer(config.DataPath)
	if err != nil {
		log.Fatalf("Error initializing server: %v", err)
	}

	sqliteDB, err = NewSQLiteDB(config.DataPath)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	err = sqliteDB.AutoMigrate(&ModelParams{}, &ImageModel{}, &SelectedModels{}, &Chat{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	searchDB := fmt.Sprintf("%s/search.bleve", config.DataPath)

	// If the database exists, open it, else create a new one
	if _, err := os.Stat(searchDB); os.IsNotExist(err) {
		mapping := bleve.NewIndexMapping()
		searchIndex, err = bleve.New(searchDB, mapping)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		searchIndex, err = bleve.Open(searchDB)
		if err != nil {
			log.Fatalf("Failed to open search index: %v", err)
		}
	}

	// Instantiate ModelParams then populate it with each model from the config
	var modelParams []ModelParams
	for _, model := range config.LanguageModels {
		if model.Downloads != nil {
			fileName := strings.Split(model.Downloads[0], "/")
			model.LocalPath = fmt.Sprintf("%s/models/%s/%s", config.DataPath, model.Name, fileName[len(fileName)-1])
		}

		var downloaded bool
		if _, err := os.Stat(model.LocalPath); err == nil {
			downloaded = true
		}

		modelParams = append(modelParams, ModelParams{
			Name:       model.Name,
			Homepage:   model.Homepage,
			GGUFInfo:   model.GGUF,
			Downloaded: downloaded,
			Options: &llm.GGUFOptions{
				Model:         model.LocalPath,
				Prompt:        model.Prompt,
				CtxSize:       model.Ctx,
				Temp:          0.7,
				RepeatPenalty: 1.1,
			},
		})
	}

	if err := LoadModelDataToDB(sqliteDB, modelParams); err != nil {
		log.Fatalf("Failed to load model data to database: %v", err)
	}

	// Instantiate ImageModel then populate it with each model from the config
	var imageModels []ImageModel
	for _, model := range config.ImageModels {
		if model.Downloads != nil {
			fileName := strings.Split(model.Downloads[0], "/")
			model.LocalPath = fmt.Sprintf("%s/models/%s/%s", config.DataPath, model.Name, fileName[len(fileName)-1])
		}

		var downloaded bool
		if _, err := os.Stat(model.LocalPath); err == nil {
			downloaded = true
		}

		imageModels = append(imageModels, ImageModel{
			Name:       model.Name,
			Homepage:   model.Homepage,
			Prompt:     model.Prompt,
			Downloaded: downloaded,
			Options: &sd.SDParams{
				Model:  model.LocalPath,
				Prompt: model.Prompt,
			},
		})
	}

	if err := LoadImageModelDataToDB(sqliteDB, imageModels); err != nil {
		log.Fatalf("Failed to load image model data to database: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pterm.Info.Printf("Serving fronted on: %s:%s\n", config.ControlHost, config.ControlPort)
	pterm.Info.Println("Press Ctrl+C to stop")

	runFrontendServer(ctx, config, modelParams)

	pterm.Warning.Println("Shutdown signal received")

	os.Exit(0)
}

func runFrontendServer(ctx context.Context, config *AppConfig, modelParams []ModelParams) {

	// Create a http fs
	basePath := filepath.Join(config.DataPath, "web")
	baseFs := afero.NewBasePathFs(osFS, basePath)
	httpFs := afero.NewHttpFs(baseFs)
	engine := html.NewFileSystem(httpFs, ".html")

	app := fiber.New(fiber.Config{
		AppName:               "Eternal v0.1.0",
		BodyLimit:             100 * 1024 * 1024, // 100MB, to allow for larger file uploads
		DisableStartupMessage: true,
		ServerHeader:          "Eternal",
		PassLocalsToViews:     true,
		Views:                 engine,
		StrictRouting:         true,
		StreamRequestBody:     true,
	})

	// CORS allow all origins for now while mvp dev mode
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
	}))

	app.Use("/public", filesystem.New(filesystem.Config{
		Root:   httpFs,
		Index:  "index.html",
		Browse: true,
	}))

	app.Static("/", "public")

	// main route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("templates/index", fiber.Map{})
	})

	app.Get("/config", func(c *fiber.Ctx) error {
		// Return the app config as JSON
		return c.JSON(config)
	})

	app.Get("/flow", func(c *fiber.Ctx) error {
		return c.Render("templates/flow", fiber.Map{})
	})

	app.Post("/upload", func(c *fiber.Ctx) error {
		pterm.Warning.Println("Uploads route hit")

		// Parse the multipart form
		form, err := c.MultipartForm()
		if err != nil {
			return err
		}

		// Get the files from the form
		files := form.File["file"]

		// Loop through the files
		for _, file := range files {
			// Save the file to the datapath web/uploads directory
			filename := filepath.Join(config.DataPath, "web", "uploads", file.Filename)
			pterm.Warning.Printf("Uploading file: %s\n", filename)
			err := c.SaveFile(file, filename)
			if err != nil {
				return err
			}

			// Log the uploaded file
			log.Infof("Uploaded file: %s", file.Filename)
		}

		// Return a success response
		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("%d files uploaded successfully", len(files)),
		})
	})

	// route to enable or disable a tool
	app.Post("/tool/:toolName", func(c *fiber.Ctx) error {
		toolName := c.Params("toolName")

		switch toolName {
		case "websearch":
			config.Tools.WebSearch.Enabled = !config.Tools.WebSearch.Enabled
		case "webget":
			config.Tools.WebGet.Enabled = !config.Tools.WebGet.Enabled
		case "imggen":
			config.Tools.ImgGen.Enabled = true
		default:
			return c.Status(fiber.StatusNotFound).SendString("Tool not found")
		}

		return c.JSON(fiber.Map{
			"message": fmt.Sprintf("Tool %s is now %t", toolName, config.Tools.ImgGen.Enabled)})
	})

	app.Get("/openai/models", func(c *fiber.Ctx) error {
		client := openai.NewClient(config.OAIKey)
		modelsResponse, err := openai.GetModels(client)

		if err != nil {
			log.Errorf(err.Error())
			return c.Status(500).SendString("Server Error")
		}

		// Filter the models to include only those with IDs starting with 'gpt'
		// This needs to be changed to a different method later. Using the name
		// is a future bug waiting to happen.
		var gptModels []string

		for _, model := range modelsResponse.Data {
			if strings.HasPrefix(model.ID, "gpt") {
				gptModels = append(gptModels, model.ID)
			}
		}

		return c.JSON(fiber.Map{
			"object": "list",
			"data":   gptModels,
		})
	})

	app.Get("/modeldata/:modelName", func(c *fiber.Ctx) error {
		var model ModelParams

		modelName := c.Params("modelName")
		err := sqliteDB.First(modelName, &model)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return c.Status(fiber.StatusNotFound).SendString("Model not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
		}

		return c.JSON(model)
	})

	app.Put("/modeldata/:modelName/downloaded", func(c *fiber.Ctx) error {
		modelName := c.Params("modelName")
		var payload struct {
			Downloaded bool `json:"downloaded"`
		}

		// Parse the JSON body to extract the 'downloaded' status
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
		}

		// Update the 'Downloaded' status of the model in the database using its name
		err := sqliteDB.UpdateDownloadedByName(modelName, payload.Downloaded)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fmt.Sprintf("Failed to update model: %v", err)})
		}

		return c.JSON(fiber.Map{
			"success": true,
			"message": "Model 'Downloaded' status updated successfully",
		})
	})

	app.Post("/modelcards", func(c *fiber.Ctx) error {
		// Retrieve all models from the database
		err := sqliteDB.Find(&modelParams)

		if err != nil {
			log.Errorf("Database error: %v", err)
			return c.Status(500).SendString("Server Error")
		}

		// Render the template with the models data
		return c.Render("templates/model", fiber.Map{"models": modelParams})
	})

	app.Post("/model/select", func(c *fiber.Ctx) error {
		var selection SelectedModels

		if err := c.BodyParser(&selection); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Bad request")
		}

		// Add or remove the model from the selection based on the action
		if selection.Action == "add" {
			if err := AddSelectedModel(sqliteDB.db, selection.ModelName); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		} else if selection.Action == "remove" {
			if err := RemoveSelectedModel(sqliteDB.db, selection.ModelName); err != nil {
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/models/selected", func(c *fiber.Ctx) error {
		selectedModels, err := GetSelectedModels(sqliteDB.db)

		if err != nil {
			log.Errorf("Error getting selected models: %v", err)
			return c.Status(500).SendString("Server Error")
		}

		var selectedModelNames []string

		for _, model := range selectedModels {
			selectedModelNames = append(selectedModelNames, model.ModelName)
		}

		return c.JSON(selectedModelNames)
	})

	app.Post("/model/download", func(c *fiber.Ctx) error {
		modelName := c.Query("model")

		if modelName == "" {
			log.Errorf("Missing parameters for download")
			return c.Status(fiber.StatusBadRequest).SendString("Missing parameters")
		}

		var downloadURL string
		for _, model := range config.LanguageModels {
			if model.Name == modelName {
				downloadURL = model.Downloads[0]
				break
			}
		}

		modelFileName := filepath.Base(downloadURL)
		modelPath := filepath.Join(config.DataPath, "models", modelName, modelFileName)

		// Check if the file exists and partially downloaded
		var partialDownload bool
		if info, err := os.Stat(modelPath); err == nil {
			// Check if the file size is less than the expected size (if available)
			if info.Size() > 0 {
				// Assuming here that we can check the expected file size somehow,
				// e.g., from a database or a config file. If not available, we
				// still try to resume assuming partial download.
				expectedSize, err := llm.GetExpectedFileSize(downloadURL)
				if err != nil {
					log.Errorf("Error getting expected file size: %v", err)
				}
				partialDownload = info.Size() < expectedSize
			}
		}

		// Download or resume the download
		go func() {
			var err error

			if partialDownload {
				pterm.Info.Printf("Resuming download for model: %s\n", modelName)
				err = llm.Download(downloadURL, modelPath)
			} else {
				pterm.Info.Printf("Starting download for model: %s\n", modelName)
				err = llm.Download(downloadURL, modelPath)
			}

			if err != nil {
				log.Errorf("Error in download: %v", err)
			} else {
				// Update the model's downloaded state in the database
				err = sqliteDB.UpdateDownloadedByName(modelName, true)
				if err != nil {
					log.Errorf("Failed to update model downloaded state: %v", err)
				}
			}
		}()

		progressErr := fmt.Sprintf("<div class='w-100' id='progress-download-%s' hx-ext='sse' sse-connect='/sseupdates' sse-swap='message' hx-trigger='load'></div>", modelName)

		return c.SendString(progressErr)
	})

	app.Post("/imgmodel/download", func(c *fiber.Ctx) error {
		config.Tools.ImgGen.Enabled = true

		modelName := c.Query("model")

		var downloadURL string
		for _, model := range config.ImageModels {
			if model.Name == modelName {
				downloadURL = model.Downloads[0]
			}
		}

		modelFileName := strings.Split(downloadURL, "/")[len(strings.Split(downloadURL, "/"))-1]

		if modelName == "" {
			log.Errorf("Missing parameters for download")
			return c.Status(fiber.StatusBadRequest).SendString("Missing parameters")
		}

		modelRoot := fmt.Sprintf("%s/models/%s", config.DataPath, modelName)
		modelPath := fmt.Sprintf("%s/models/%s/%s", config.DataPath, modelName, modelFileName)
		tmpPath := fmt.Sprintf("%s/tmp", config.DataPath)

		// Create the model directory if it doesn't exist
		if _, err := os.Stat(modelRoot); os.IsNotExist(err) {
			if err := os.MkdirAll(modelRoot, 0755); err != nil {
				log.Errorf("Error creating model directory: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		}

		// Create the tmp directory if it doesn't exist
		if _, err := os.Stat(tmpPath); os.IsNotExist(err) {
			if err := os.MkdirAll(tmpPath, 0755); err != nil {
				log.Errorf("Error creating tmp directory: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		}

		// Check if the modelPath does not exist and download it if it doesn't
		if _, err := os.Stat(modelPath); err != nil {
			// Start the download in a goroutine
			dm := hfutils.ConcurrentDownloadManager{
				FileName:    modelFileName,
				URL:         downloadURL,
				Destination: modelPath,
				NumParts:    1,
				TempDir:     tmpPath,
				//Sha256Checksum: "abc123...", // Optional, provide if needed.
			}

			go dm.PrintProgress()

			if err := dm.Download(); err != nil {
				fmt.Println("Download failed:", err)
			} else {
				fmt.Println("Download successful!")
			}
		}

		// https://huggingface.co/madebyollin/sdxl-vae-fp16-fix/blob/main/sdxl_vae.safetensors
		vaeName := "sdxl_vae.safetensors"
		vaeURL := "https://huggingface.co/madebyollin/sdxl-vae-fp16-fix/blob/main/sdxl_vae.safetensors"
		vaePath := fmt.Sprintf("%s/models/%s/%s", config.DataPath, modelName, vaeName)

		// Create the model directory if it doesn't exist
		if _, err := os.Stat(modelRoot); os.IsNotExist(err) {
			if err := os.MkdirAll(modelRoot, 0755); err != nil {
				log.Errorf("Error creating model directory: %v", err)
				return c.Status(fiber.StatusInternalServerError).SendString("Server Error")
			}
		}

		// Download SDXL VAE fix if it doesn't exist
		if _, err := os.Stat(vaePath); os.IsNotExist(err) {
			// Start the download in a goroutine
			go func() {
				// Download the file
				response, err := http.Get(vaeURL)
				if err != nil {
					pterm.Error.Printf("Failed to download file: %v", err)
					return
				}
				defer response.Body.Close()

				// Create the file
				file, err := os.Create(vaePath)
				if err != nil {
					pterm.Error.Printf("Failed to create file: %v", err)
					return
				}
				defer file.Close()

				// Write the downloaded data to the file
				_, err = io.Copy(file, response.Body)
				if err != nil {
					pterm.Error.Printf("Failed to write to file: %v", err)
					return
				}

				pterm.Info.Printf("Downloaded file: %s", vaeName)
			}()
		}

		progressErr := "<div name='sse-messages' class='w-100' id='sse-messages' hx-ext='sse' sse-connect='/sseupdates' sse-swap='message'></div>"

		return c.SendString(progressErr)
	})

	app.Post("/chattemplates", func(c *fiber.Ctx) error {
		modelsFile := fmt.Sprintf("%v/chat-templates.json", config)

		chatTemplates, err := os.ReadFile(modelsFile)
		if err != nil {
			log.Errorf(err.Error())
			return c.Status(500).SendString("Server Error")
		}

		var chatTemplate []llm.ChatPromptTemplate
		err = json.Unmarshal(chatTemplates, &chatTemplate)

		if err != nil {
			log.Errorf(err.Error())
			return c.Status(500).SendString("Server Error")
		}

		return c.Render("templates/chattemplates", fiber.Map{"templates": chatTemplate})
	})

	app.Post("/chatsubmit", func(c *fiber.Ctx) error {

		// userPrompt is the message displayed in the chat view
		userPrompt := c.FormValue("userprompt")

		var wsroute string

		selectedModels, err := GetSelectedModels(sqliteDB.db)
		if err != nil {
			log.Errorf("Error getting selected models: %v", err)
			return c.Status(500).SendString("Server Error")
		}

		if len(selectedModels) > 0 {
			firstModelName := selectedModels[0].ModelName

			// Check if the first model name starts with "openai-"
			if strings.HasPrefix(firstModelName, "openai-") {
				wsroute = "/wsoai"
			} else if strings.HasPrefix(firstModelName, "google-") {
				wsroute = "/wsgoogle"
			} else if strings.HasPrefix(firstModelName, "anthropic-") {
				wsroute = "/wsanthropic"
			} else {
				wsroute = fmt.Sprintf("ws://%s:%s/ws", config.ServiceHosts["llm"]["llm_host_1"].Host, config.ServiceHosts["llm"]["llm_host_1"].Port)
			}

		} else {
			// return error
			return c.JSON(fiber.Map{"error": "No models selected"})
		}

		// Generate unique ID
		turnID := IncrementTurn()

		return c.Render("templates/chat", fiber.Map{
			"username":  config.CurrentUser,
			"message":   userPrompt, // This is the message that will be displayed in the chat
			"assistant": config.AssistantName,
			"model":     selectedModels[0].ModelName,
			"turnID":    turnID,
			"wsRoute":   wsroute,
			"hosts":     config.ServiceHosts["llm"],
		})
	})

	// Retrieve all chats from sqlite database
	app.Get("/chats", func(c *fiber.Ctx) error {
		chats, err := GetChats(sqliteDB.db)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not get chats"})
		}

		return c.Status(fiber.StatusOK).JSON(chats)
	})

	// Retrieve a single chat from sqlite database by id
	app.Get("/chats/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		chat, err := GetChatByID(sqliteDB.db, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not get chat"})
		}

		return c.Status(fiber.StatusOK).JSON(chat)
	})

	// Update a single chat in sqlite database by id
	app.Put("/chats/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		chat := new(Chat)

		if err := c.BodyParser(chat); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse JSON"})
		}

		err = UpdateChat(sqliteDB.db, id, chat.Prompt, chat.Response, chat.ModelName)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not update chat"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	// Delete a single chat in sqlite database by id
	app.Delete("/chats/:id", func(c *fiber.Ctx) error {
		id, err := strconv.ParseInt(c.Params("id"), 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
		}

		err = DeleteChat(sqliteDB.db, id)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not delete chat"})
		}

		return c.SendStatus(fiber.StatusNoContent)
	})

	// Multi web page retrieval via local ChromeDP
	app.Get("/dpsearch", func(c *fiber.Ctx) error {
		urls := []string{}
		query := c.Query("q")
		res := web.SearchDDG(query)

		if len(res) == 0 {
			return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving search results")
		}

		// Remove youtube results
		urls = append(urls, res...)

		// Send results as json object
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"urls": urls})
	})

	app.Get("/sseupdates", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			for {
				// Get updated download progress
				progress := llm.GetDownloadProgress("sse-progress")

				// Format message for SSE
				msg := fmt.Sprintf("data: <div class='progress specific-h-25 m-4' role='progressbar' aria-label='download' aria-valuenow='%s' aria-valuemin='0' aria-valuemax='100'><div class='progress-bar progress-bar-striped progress-bar-animated' style='width: %s;'></div></div><div class='text-center fs-6'>Please refresh this page when the download completes.</br> Downloading...%s</div>\n\n", progress, progress, progress)

				// Write the message
				if _, err := w.WriteString(msg); err != nil {
					pterm.Printf("Error writing to stream: %v", err)
					break
				}
				if err := w.Flush(); err != nil {
					pterm.Printf("Error flushing writer: %v", err)
					break
				}

				time.Sleep(2 * time.Second) // Adjust the sleep time as necessary
			}
		}))

		return nil
	})

	app.Get("/ws", websocket.New(func(c *websocket.Conn) {
		handleWebSocket(c, config, func(wsMessage WebSocketMessage, chatMessage string) error {
			// Process the message
			//cpt := llm.GetSystemTemplate(chatMessage)
			//fullPrompt := cpt.Messages[0].Content + "\n" + chatMessage

			// Get the details of the first model from database
			var model ModelParams
			err := sqliteDB.First(wsMessage.Model, &model)
			if err != nil {
				log.Errorf("Error getting model %s: %v", wsMessage.Model, err)
				return err
			}

			promptTemplate := model.Options.Prompt

			// Replace {user} with the chat Message
			fullPrompt := strings.ReplaceAll(promptTemplate, "{prompt}", chatMessage)

			// Replace {system} with the system message
			//fullPrompt = strings.ReplaceAll(fullPrompt, "{system}", "You are a helpful AI assistant that responds in well structured markdown format. Do not repeat your instructions. Do not deviate from the topic. Begin all responses with 'Sure thing!' and end with 'Is there anything else I can help you with?'")
			fullPrompt = strings.ReplaceAll(fullPrompt, "{system}", llm.AssistantDefault)

			modelOpts := &llm.GGUFOptions{
				NGPULayers:    config.ServiceHosts["llm"]["llm_host_1"].GgufGPULayers,
				Model:         model.Options.Model,
				Prompt:        fullPrompt,
				CtxSize:       model.Options.CtxSize,
				Temp:          0.1, // Prefer lower temperature for more controlled responses for now
				RepeatPenalty: 1.1,
				TopP:          1.0, // Prefer greedy decoding for now
				TopK:          1.0, // Prefer greedy decoding for now
			}

			// Search the search index for the chat message
			// searchResults, err := search.Search(searchIndex, chatMessage)
			// if err != nil {
			// 	log.Errorf("Error searching index: %v", err)
			// }

			// search for some text
			// query := bleve.NewMatchQuery(chatMessage)
			// search := bleve.NewSearchRequest(query)
			// searchResults, err := searchIndex.Search(search)
			// if err != nil {
			// 	fmt.Println(err)
			// 	return err
			// }
			// pterm.Info.Println(searchResults)

			////////////////////////
			// AGENT REPLIES
			///////////////////////

			advWorkflow := false
			if advWorkflow {

				res1 := llm.MakeCompletionWebSocket(*c, chatTurn, modelOpts, config.DataPath)

				var smodel ModelParams
				newModel := "llama3-70b-instruct"
				err = sqliteDB.First(newModel, &smodel)
				if err != nil {
					log.Errorf("Error getting model %s: %v", newModel, err)
					return err
				}

				nextPrompt := fmt.Sprintf("%s\nNew Instructions:\n%s\n", res1, llm.AssistantVisualBot)

				smodelOpts := &llm.GGUFOptions{
					NGPULayers:    config.ServiceHosts["llm"]["llm_host_1"].GgufGPULayers,
					Model:         smodel.Options.Model,
					Prompt:        nextPrompt,
					CtxSize:       smodel.Options.CtxSize,
					Temp:          0.1, // Prefer lower temperature for more controlled responses for now
					RepeatPenalty: 1.1,
					TopP:          1.0, // Prefer greedy decoding for now
					TopK:          1.0, // Prefer greedy decoding for now
				}

				return llm.LMResponse(*c, chatTurn, smodelOpts, config.DataPath)
			}

			return llm.MakeCompletionWebSocket(*c, chatTurn, modelOpts, config.DataPath)
		})
	}))

	app.Get("/wsoai", websocket.New(func(c *websocket.Conn) {
		apiKey := config.OAIKey

		handleWebSocket(c, config, func(wsMessage WebSocketMessage, chatMessage string) error {
			// Check if embeddings.db exists
			// if _, err := os.Stat(filepath.Join(config.DataPath, "embeddings.db")); os.IsNotExist(err) {
			// 	pterm.Warning.Println("embeddings.db does not exist. Generating embeddings...")
			// 	embeddings.GenerateEmbeddingChat(chatMessage, config.DataPath)
			// }

			cpt := llm.GetSystemTemplate(chatMessage)
			return openai.StreamCompletionToWebSocket(c, chatTurn, "gpt-4o", cpt.Messages, 0.3, apiKey)
		})
	}))

	app.Get("/wsanthropic", websocket.New(func(c *websocket.Conn) {
		apiKey := config.AnthropicKey

		handleAnthropicWS(c, apiKey, chatTurn)
	}))

	app.Get("/wsgoogle", websocket.New(func(c *websocket.Conn) {
		apiKey := config.GoogleKey

		handleWebSocket(c, config, func(wsMessage WebSocketMessage, chatMessage string) error {
			return google.StreamGeminiResponseToWebSocket(c, chatTurn, chatMessage, apiKey)
		})
	}))

	go func() {
		<-ctx.Done() // Wait for the context to be cancelled
		if err := app.Shutdown(); err != nil {
			log.Fatalf("Server shutdown failed: %v", err)
		}
	}()

	addr := fmt.Sprintf("%s:%s", config.ControlHost, config.ControlPort)
	if err := app.Listen(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Frontend server failed: %v", err)
	}

	pterm.Info.Println("Server gracefully shutdown")
}

func handleWebSocket(c *websocket.Conn, config *AppConfig, processMessage func(WebSocketMessage, string) error) {
	if c == nil {
		pterm.Error.Println("WebSocket connection is nil")
		return
	}
	defer c.Close()

	// Read the initial message
	_, message, err := c.ReadMessage()
	if err != nil {
		pterm.PrintOnError(err)
		return
	}

	// Unmarshal the JSON message
	var wsMessage WebSocketMessage
	err = json.Unmarshal(message, &wsMessage)
	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte("Error unmarshalling JSON"))
		return
	}

	// Extract the chat_message value
	chatMessage := wsMessage.ChatMessage

	// Perform tool workflow and update chatMessage
	chatMessage = performToolWorkflow(c, config, chatMessage)

	// If image generation is enabled, return early
	if config.Tools.ImgGen.Enabled {
		return
	}

	// Process the message using the provided function
	res := processMessage(wsMessage, chatMessage)
	if res != nil {
		//pterm.Warning.Println(res)

		if config.Tools.Memory.Enabled {
			err = storeChat(sqliteDB.db, config, chatMessage, res.Error(), wsMessage.Model)
			if err != nil {
				pterm.PrintOnError(err)
			}
		}

		// Increment the chat turn counter
		chatTurn = chatTurn + 1
		pterm.Warning.Println("Chat turn:", chatTurn)
		return
	}
}

func performToolWorkflow(c *websocket.Conn, config *AppConfig, chatMessage string) string {
	// Begin tool workflow. Tools will add context to the submitted message for
	// the model to use. Document is the abstraction that will hold that context.
	var document string

	if config.Tools.ImgGen.Enabled {
		pterm.Info.Println("Generating image...")
		sdParams := &sd.SDParams{Prompt: chatMessage}

		// Call the sd tool
		sd.Text2Image(config.DataPath, sdParams)

		// Return the image to the client
		timestamp := time.Now().UnixNano() // Get the current timestamp in nanoseconds
		imgElement := fmt.Sprintf("<img class='rounded-2 object-fit-scale' width='512' height='512' src='public/img/sd_out.png?%d' />", timestamp)
		formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>", fmt.Sprint(chatTurn), imgElement)
		if err := c.WriteMessage(websocket.TextMessage, []byte(formattedContent)); err != nil {
			pterm.PrintOnError(err)
			return chatMessage
		}

		// Increment the chat turn counter
		chatTurn = chatTurn + 1

		// End the tool workflow
		return chatMessage
	}

	if config.Tools.Memory.Enabled {
		topN := config.Tools.Memory.TopN // retrieve top N results. Adjust based on context size.
		topEmbeddings := embeddings.Search(config.DataPath, "embeddings.db", chatMessage, topN)

		var documents []string
		var documentString string
		if len(topEmbeddings) > 0 {
			for _, topEmbedding := range topEmbeddings {
				documents = append(documents, topEmbedding.Word)
			}
			documentString = strings.Join(documents, " ")

			pterm.Info.Println("Retrieving memory content...")
			document = fmt.Sprintf("%s\n%s", document, documentString)

			// Replace new lines with spaces
			document = strings.ReplaceAll(document, "\n\n", "\n")
		} else {
			pterm.Info.Println("No memory content found...")
		}
	}

	if config.Tools.WebGet.Enabled {
		url := web.ExtractURLs(chatMessage)
		if len(url) > 0 {
			pterm.Info.Println("Retrieving page content...")

			document, _ = web.WebGetHandler(url[0])
		}
	}

	if config.Tools.WebSearch.Enabled {

		topN := config.Tools.WebSearch.TopN // retrieve top N results. Adjust based on context size.

		pterm.Info.Println("Searching the web...")

		urls := web.SearchDDG(chatMessage)

		//pterm.Warning.Printf("URLs to fetch: %v\n", urls)

		if len(urls) > 0 {
			pagesRetrieved := 0

			for {
				// Check if we have collected topN pages
				if pagesRetrieved >= topN {
					break
				}

				// Iterate over URLs
				for _, url := range urls {
					pterm.Info.Printf("Fetching URL: %s\n", url)

					page, err := web.WebGetHandler(url)
					if err != nil {
						if errors.Is(err, context.DeadlineExceeded) {
							pterm.Warning.Printf("Timeout exceeded for URL: %s\n", url)

							// Remove the URL from the list, do not use the web package
							urls = urls[1:]

							pterm.Warning.Printf("URL list: %s\n", urls)

							// Increase the timeout for the next request to avoid spamming the same URL
							time.Sleep(5 * time.Second)

							continue
						}
						pterm.PrintOnError(err)
					} else {
						// Page successfully retrieved, update document and increment pagesRetrieved
						document = fmt.Sprintf("%s\n%s", document, page)
						pagesRetrieved++

						// Check if we have collected topN pages
						if pagesRetrieved >= topN {
							break
						}
					}
				}
			}
		}
	}

	//Remove http(s) links from the document so we do not retrieve them unintentionally
	document = web.RemoveUrls(document)

	chatMessage = fmt.Sprintf("%s Reference the previous information and respond to the following task or question:\n%s", document, chatMessage)

	pterm.Error.Println("Tool workflow complete")

	return chatMessage
}

func storeChat(db *gorm.DB, config *AppConfig, prompt, response, modelName string) error {
	// Generate embeddings
	pterm.Warning.Println("Generating embeddings for chat...")

	err := embeddings.GenerateEmbeddingForTask("chat", response, "txt", 2048, 500, config.DataPath)
	if err != nil {
		pterm.Error.Println("Error generating embeddings:", err)
		return err
	}

	pterm.Warning.Print("Storing chat in database...")
	if _, err := CreateChat(db, prompt, response, modelName); err != nil {
		pterm.Error.Println("Error storing chat in database:", err)
		return err
	}

	return nil
}

func handleAnthropicWS(c *websocket.Conn, apiKey string, chatID int) {
	// Read the initial message
	_, message, err := c.ReadMessage()
	if err != nil {
		pterm.PrintOnError(err)
		return
	}

	// Unmarshal the JSON message
	var wsMessage WebSocketMessage
	err = json.Unmarshal(message, &wsMessage)
	if err != nil {
		c.WriteMessage(websocket.TextMessage, []byte("Error unmarshalling JSON"))
		return
	}

	// Extract the chat_message value
	chatMessage := wsMessage.ChatMessage

	messages := []anthropic.Message{
		{Role: "user", Content: chatMessage},
	}

	res := anthropic.StreamCompletionToWebSocket(c, chatID, "claude-3-opus-20240229", messages, 0.5, apiKey)
	if res != nil {
		pterm.Error.Println("Error in anthropic completion:", res)
	}

	chatTurn = chatTurn + 1
}

[File Ends] main.go

[File Begins] main_test.go
package main

import (
	"bytes"
	"encoding/json"
	"eternal/pkg/llm"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestIndexRoute(t *testing.T) {
	app := fiber.New()

	// Set up routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Send a GET request to the root URL
	req, _ := http.NewRequest("GET", "/", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	body, _ := io.ReadAll(resp.Body)    // Read the response body
	assert.Equal(t, "OK", string(body)) // Convert body to string and compare with "OK"
}

func TestConfigRoute(t *testing.T) {
	app := fiber.New()

	// Set up routes
	app.Get("/config", func(c *fiber.Ctx) error {
		config := &AppConfig{
			CurrentUser:   "test_user",
			AssistantName: "test_assistant",
		}
		return c.JSON(config)
	})

	// Send a GET request to the /config URL
	req, _ := http.NewRequest("GET", "/config", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var config AppConfig
	err = json.NewDecoder(resp.Body).Decode(&config)
	assert.NoError(t, err)
	assert.Equal(t, "test_user", config.CurrentUser)
	assert.Equal(t, "test_assistant", config.AssistantName)
}

func TestToolRoute(t *testing.T) {
	app := fiber.New()

	// Set up routes
	app.Post("/tool/:toolName", func(c *fiber.Ctx) error {
		toolName := c.Params("toolName")
		var index int
		found := false
		for i, t := range tools {
			if t.Name == toolName {
				index = i
				found = true
				break
			}
		}
		if !found {
			return c.Status(404).SendString("Tool not found")
		}
		tools[index].Enabled = !tools[index].Enabled
		return c.JSON(tools[index])
	})

	// Set up test data
	tools = []Tool{
		{Name: "websearch", Enabled: false},
		{Name: "imagegen", Enabled: false},
	}

	// Test enabling a tool
	body := []byte{}
	req, _ := http.NewRequest("POST", "/tool/websearch", bytes.NewBuffer(body))
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	var tool Tool
	err = json.NewDecoder(resp.Body).Decode(&tool)
	assert.NoError(t, err)
	assert.Equal(t, "websearch", tool.Name)
	assert.True(t, tool.Enabled)

	// Test disabling a tool
	req, _ = http.NewRequest("POST", "/tool/websearch", bytes.NewBuffer(body))
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	err = json.NewDecoder(resp.Body).Decode(&tool)
	assert.NoError(t, err)
	assert.Equal(t, "websearch", tool.Name)
	assert.False(t, tool.Enabled)

	// Test non-existent tool
	req, _ = http.NewRequest("POST", "/tool/nonexistent", bytes.NewBuffer(body))
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
	body, _ = io.ReadAll(resp.Body)
	assert.Equal(t, "Tool not found", string(body))
}

func TestModelDataRoute(t *testing.T) {
	app := fiber.New()

	app.Get("/modeldata/:modelName", func(c *fiber.Ctx) error {
		mockData := ModelParams{
			ID:         3,
			Name:       "eternal-120b",
			Homepage:   "https://huggingface.co/intelligence-dev/eternal-120b",
			GGUFInfo:   "https://huggingface.co/intelligence-dev/eternal-120b",
			Downloads:  "",
			Downloaded: true,
			Options: &llm.GGUFOptions{
				Prompt: "Test",
			},
		}

		return c.JSON(mockData)
	})

	req := httptest.NewRequest("GET", "/modeldata/test-model", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)

	var modelData ModelParams
	err = json.NewDecoder(resp.Body).Decode(&modelData)

	assert.NoError(t, err)
	assert.Equal(t, "eternal-120b", modelData.Name)
}

[File Ends] main_test.go

    [File Begins] pkg/documents/gitloader.go
    package documents
    
    import (
    	"fmt"
    	"io/ioutil"
    	"os"
    	"path/filepath"
    
    	gogit "github.com/go-git/go-git/v5"
    	"github.com/go-git/go-git/v5/plumbing"
    	gitssh "github.com/go-git/go-git/v5/plumbing/transport/ssh"
    	golangssh "golang.org/x/crypto/ssh"
    )
    
    type Document struct {
    	PageContent string
    	Metadata    map[string]string
    }
    
    type GitLoader struct {
    	RepoPath           string
    	CloneURL           string
    	Branch             string
    	PrivateKeyPath     string
    	FileFilter         func(string) bool
    	InsecureSkipVerify bool
    }
    
    func NewGitLoader(repoPath, cloneURL, branch, privateKeyPath string, fileFilter func(string) bool, insecureSkipVerify bool) *GitLoader {
    	return &GitLoader{RepoPath: repoPath, CloneURL: cloneURL, Branch: branch, PrivateKeyPath: privateKeyPath, FileFilter: fileFilter, InsecureSkipVerify: insecureSkipVerify}
    }
    
    // Load loads the documents from the Git repository specified by the GitLoader.
    // It returns a slice of Document and an error if any.
    // If the repository does not exist at the specified path and a clone URL is provided,
    // it clones the repository using the provided authentication options.
    // If the repository already exists, it opens the repository at the specified path.
    // If a branch is specified, it checks out the branch.
    // It then walks through the repository files, reads the content of each file,
    // and creates a Document object for each file with the corresponding metadata.
    // The resulting documents are returned as a slice.
    // If any error occurs during the process, it is returned.
    func (gl *GitLoader) Load() ([]Document, error) {
    	var repo *gogit.Repository
    	var err error
    
    	if _, err = os.Stat(gl.RepoPath); os.IsNotExist(err) && gl.CloneURL != "" {
    		sshKey, _ := os.ReadFile(gl.PrivateKeyPath)
    		signer, _ := golangssh.ParsePrivateKey(sshKey)
    		auth := &gitssh.PublicKeys{User: "git", Signer: signer}
    		if gl.InsecureSkipVerify {
    			auth.HostKeyCallback = golangssh.InsecureIgnoreHostKey()
    		}
    		repo, err = gogit.PlainClone(gl.RepoPath, false, &gogit.CloneOptions{URL: gl.CloneURL, Auth: auth})
    		if err != nil {
    			return nil, err
    		}
    	} else {
    		repo, err = gogit.PlainOpen(gl.RepoPath)
    		if err != nil {
    			return nil, err
    		}
    	}
    
    	if gl.Branch != "" {
    		w, err := repo.Worktree()
    		if err != nil {
    			return nil, err
    		}
    		err = w.Checkout(&gogit.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(gl.Branch)})
    		if err != nil {
    			return nil, err
    		}
    	}
    
    	var docs []Document
    
    	err = filepath.Walk(gl.RepoPath, func(path string, info os.FileInfo, err error) error {
    		if err != nil {
    			return err
    		}
    
    		if info.IsDir() {
    			return nil
    		}
    
    		if gl.FileFilter != nil && !gl.FileFilter(path) {
    			return nil
    		}
    
    		content, err := ioutil.ReadFile(path)
    		if err != nil {
    			return err
    		}
    
    		textContent := string(content)
    		relFilePath, _ := filepath.Rel(gl.RepoPath, path)
    		fileType := filepath.Ext(info.Name())
    
    		metadata := map[string]string{
    			"source":    relFilePath,
    			"file_path": relFilePath,
    			"file_name": info.Name(),
    			"file_type": fileType,
    		}
    
    		doc := Document{PageContent: textContent, Metadata: metadata}
    		docs = append(docs, doc)
    
    		return nil
    	})
    
    	if err != nil {
    		fmt.Printf("Error reading files: %s\n", err)
    	}
    
    	return docs, nil
    }

    [File Ends] pkg/documents/gitloader.go

    [File Begins] pkg/documents/txtsplitter.go
    package documents
    
    import (
    	"errors"
    	"regexp"
    )
    
    // RecursiveCharacterTextSplitter is a struct that represents a text splitter
    // that splits text based on recursive character separators.
    type RecursiveCharacterTextSplitter struct {
    	Separators       []string
    	KeepSeparator    bool
    	IsSeparatorRegex bool
    	ChunkSize        int
    	OverlapSize      int
    	LengthFunction   func(string) int
    }
    
    // Language is a type that represents a programming language.
    type Language string
    
    const (
    	PYTHON   Language = "PYTHON"
    	GO       Language = "GO"
    	HTML     Language = "HTML"
    	JS       Language = "JS"
    	TS       Language = "TS"
    	MARKDOWN Language = "MARKDOWN"
    	JSON     Language = "JSON"
    )
    
    // escapeString is a helper function that escapes special characters in a string.
    func escapeString(s string) string {
    	return regexp.QuoteMeta(s)
    }
    
    // splitTextWithRegex is a helper function that splits text using a regular expression separator.
    func splitTextWithRegex(text string, separator string, keepSeparator bool) []string {
    	sepPattern := regexp.MustCompile(separator)
    	splits := sepPattern.Split(text, -1)
    	if keepSeparator {
    		matches := sepPattern.FindAllString(text, -1)
    		result := make([]string, 0, len(splits)+len(matches))
    		for i, split := range splits {
    			result = append(result, split)
    			if i < len(matches) {
    				result = append(result, matches[i])
    			}
    		}
    		return result
    	}
    	return splits
    }
    
    // SplitTextByCount splits the given text into chunks of the given size.
    func SplitTextByCount(text string, size int) []string {
    	// slice the string into chunks of size
    	var chunks []string
    	for i := 0; i < len(text); i += size {
    		end := i + size
    		if end > len(text) {
    			end = len(text)
    		}
    		chunks = append(chunks, text[i:end])
    	}
    	return chunks
    }
    
    // SplitText splits the given text using the configured separators.
    func (r *RecursiveCharacterTextSplitter) SplitText(text string) []string {
    	chunks := r.splitTextHelper(text, r.Separators)
    
    	// Apply chunk overlap
    	if r.OverlapSize > 0 {
    		overlappedChunks := make([]string, 0)
    		for i := 0; i < len(chunks)-1; i++ {
    			currentChunk := chunks[i]
    			nextChunk := chunks[i+1]
    
    			nextChunkOverlap := nextChunk[:min(len(nextChunk), r.OverlapSize)]
    
    			overlappedChunk := currentChunk + nextChunkOverlap
    			overlappedChunks = append(overlappedChunks, overlappedChunk)
    		}
    		overlappedChunks = append(overlappedChunks, chunks[len(chunks)-1])
    
    		chunks = overlappedChunks
    	}
    
    	return chunks
    }
    
    // splitTextHelper is a recursive helper function that splits text using the given separators.
    func (r *RecursiveCharacterTextSplitter) splitTextHelper(text string, separators []string) []string {
    	finalChunks := make([]string, 0)
    
    	if len(separators) == 0 {
    		return []string{text}
    	}
    
    	// Determine the separator
    	separator := separators[len(separators)-1]
    	newSeparators := make([]string, 0)
    	for i, sep := range separators {
    		sepPattern := sep
    		if !r.IsSeparatorRegex {
    			sepPattern = escapeString(sep)
    		}
    		if regexp.MustCompile(sepPattern).MatchString(text) {
    			separator = sep
    			newSeparators = separators[i+1:]
    			break
    		}
    	}
    
    	// Split the text using the determined separator
    	splits := splitTextWithRegex(text, separator, r.KeepSeparator)
    
    	// Check each split
    	for _, s := range splits {
    		if r.LengthFunction(s) < r.ChunkSize {
    			finalChunks = append(finalChunks, s)
    		} else if len(newSeparators) > 0 {
    			// If the split is too large, try to split it further using remaining separators
    			recursiveSplits := r.splitTextHelper(s, newSeparators)
    			finalChunks = append(finalChunks, recursiveSplits...)
    		} else {
    			// If no more separators left, add the large chunk as it is
    			finalChunks = append(finalChunks, s)
    		}
    	}
    
    	return finalChunks
    }
    
    // FromLanguage creates a RecursiveCharacterTextSplitter based on the given language.
    func FromLanguage(language Language) (*RecursiveCharacterTextSplitter, error) {
    	separators, err := GetSeparatorsForLanguage(language)
    	if err != nil {
    		return nil, err
    	}
    	return &RecursiveCharacterTextSplitter{
    		Separators:       separators,
    		IsSeparatorRegex: true,
    	}, nil
    }
    
    // GetSeparatorsForLanguage returns the separators for the given language.
    func GetSeparatorsForLanguage(language Language) ([]string, error) {
    	switch language {
    	case PYTHON:
    		return []string{
    			// Split along class definitions
    			"\nclass ",
    			"\ndef ",
    			"\n\tdef ",
    			// Split by the normal type of lines
    			"\n\n",
    			"\n",
    			" ",
    			"",
    		}, nil
    	case GO:
    		return []string{
    			// Split along function definitions
    			"\nfunc ",
    			"\nvar ",
    			"\nconst ",
    			"\ntype ",
    			// Split along control flow statements
    			"\nif ",
    			"\nfor ",
    			"\nswitch ",
    			"\ncase ",
    			// Split by the normal type of lines
    			"\n\n",
    			"\n",
    			" ",
    			"",
    		}, nil
    	case HTML:
    		return []string{
    			// Split along HTML tags
    			"<body",
    			"<div",
    			"<p",
    			"<br",
    			"<li",
    			"<h1",
    			"<h2",
    			"<h3",
    			"<h4",
    			"<h5",
    			"<h6",
    			"<span",
    			"<table",
    			"<tr",
    			"<td",
    			"<th",
    			"<ul",
    			"<ol",
    			"<header",
    			"<footer",
    			"<nav",
    			// Head
    			"<head",
    			"<style",
    			"<script",
    			"<meta",
    			"<title",
    			"",
    			"\n</",
    		}, nil
    	case JS:
    		return []string{
    			// Split along function definitions
    			"\nfunction ",
    			"\nconst ",
    			"\nlet ",
    			"\nvar ",
    			"\nclass ",
    			// Split along control flow statements
    			"\nif ",
    			"\nfor ",
    			"\nwhile ",
    			"\nswitch ",
    			"\ncase ",
    			"\ndefault ",
    			// Split by the normal type of lines
    			"\n\n",
    			"\n",
    			" ",
    			"",
    		}, nil
    	case TS:
    		return []string{
    			"\nenum ",
    			"\ninterface ",
    			"\nnamespace ",
    			"\ntype ",
    			// Split along class definitions
    			"\nclass ",
    			// Split along function definitions
    			"\nfunction ",
    			"\nconst ",
    			"\nlet ",
    			"\nvar ",
    			// Split along control flow statements
    			"\nif ",
    			"\nfor ",
    			"\nwhile ",
    			"\nswitch ",
    			"\ncase ",
    			"\ndefault ",
    			// Split by the normal type of lines
    			"\n\n",
    			"\n",
    			" ",
    			"",
    		}, nil
    	case MARKDOWN:
    		return []string{
    			// First, try to split along Markdown headings (starting with level 2)
    			"\n#{1,6} ",
    			// Note the alternative syntax for headings (below) is not handled here
    			// Heading level 2
    			// ---------------
    			// End of code block
    			"```\n",
    			// Horizontal lines
    			"\n\\*\\*\\*+\n",
    			"\n---+\n",
    			"\n___+\n",
    			// Note that this splitter doesn't handle horizontal lines defined
    			// by *three or more* of ***, ---, or ___, but this is not handled
    			"\n\n",
    			"\n",
    			" ",
    			"",
    		}, nil
    	case JSON:
    		return []string{
    			"}\n",
    		}, nil
    	default:
    		return nil, errors.New("unsupported language")
    	}
    }
    
    // Helper functions
    func min(a, b int) int {
    	if a < b {
    		return a
    	}
    	return b
    }
    
    func max(a, b int) int {
    	if a > b {
    		return a
    	}
    	return b
    }

    [File Ends] pkg/documents/txtsplitter.go

    [File Begins] pkg/embeddings/local.go
    package embeddings
    
    import (
    	"context"
    	"eternal/pkg/documents"
    	"fmt"
    	"strings"
    
    	estore "eternal/pkg/vecstore"
    
    	"github.com/nlpodyssey/cybertron/pkg/models/bert"
    	"github.com/nlpodyssey/cybertron/pkg/tasks"
    	"github.com/nlpodyssey/cybertron/pkg/tasks/textencoding"
    	"github.com/pterm/pterm"
    )
    
    // var modelName = "BAAI/bge-large-en-v1.5"
    var modelName = "avsolatorio/GIST-small-Embedding-v0"
    
    const limit = 128
    
    var INSTRUCTIONS = map[string]struct {
    	Query string
    	Key   string
    }{
    	"qa": {
    		Query: "Represent this query for retrieving relevant documents: ",
    		Key:   "Represent this document for retrieval: ",
    	},
    	"icl": {
    		Query: "Convert this example into vector to look for useful examples: ",
    		Key:   "Convert this example into vector for retrieval: ",
    	},
    	"chat": {
    		Query: "Embed this dialogue to find useful historical dialogues: ",
    		Key:   "Embed this historical dialogue for retrieval: ",
    	},
    	"lrlm": {
    		Query: "Embed this text chunk for finding useful historical chunks: ",
    		Key:   "Embed this historical text chunk for retrieval: ",
    	},
    	"tool": {
    		Query: "Transform this user request for fetching helpful tool descriptions: ",
    		Key:   "Transform this tool description for retrieval: ",
    	},
    	"convsearch": {
    		Query: "Encode this query and context for searching relevant passages: ",
    		Key:   "Encode this passage for retrieval: ",
    	},
    }
    
    // Embedding represents a word embedding.
    type Embedding struct {
    	Word       string
    	Vector     []float64
    	Similarity float64
    }
    
    func GenerateEmbeddingForTask(task string, content string, doctype string, chunkSize int, overlapSize int, dataPath string) error {
    
    	_, ok := INSTRUCTIONS[task]
    	if !ok {
    		fmt.Printf("Unknown task: %s\n", task)
    		return fmt.Errorf("unknown task: %s", task)
    	}
    
    	db := estore.NewEmbeddingDB()
    
    	var chunks []string
    	var separators []string
    
    	if doctype == "txt" {
    		// convert to lower case
    		content = strings.ToLower(content)
    		chunks = documents.SplitTextByCount(string(content), chunkSize)
    
    	} else {
    		doctype = strings.ToUpper(doctype)
    		separators, _ = documents.GetSeparatorsForLanguage(documents.Language(doctype))
    
    		overlapSize := chunkSize / 2 // Set the overlap size to half of the chunk size
    
    		splitter := documents.RecursiveCharacterTextSplitter{
    			Separators:       separators,
    			KeepSeparator:    true,
    			IsSeparatorRegex: false,
    			ChunkSize:        chunkSize,
    			OverlapSize:      overlapSize, // Add the OverlapSize field
    			LengthFunction:   func(s string) int { return len(s) },
    		}
    		chunks = splitter.SplitText(string(content))
    	}
    
    	// Remove duplicate chunks
    	seen := make(map[string]bool)
    	var uniqueChunks []string
    	for _, chunk := range chunks {
    
    		if _, ok := seen[chunk]; !ok {
    			uniqueChunks = append(uniqueChunks, chunk)
    			seen[chunk] = true
    		}
    	}
    
    	modelsDir := fmt.Sprintf("%s/data/models/HF/%s/", dataPath, modelName)
    
    	tasksConfig := &tasks.Config{
    		ModelsDir:        modelsDir,
    		ModelName:        modelName,
    		DownloadPolicy:   tasks.DownloadMissing,
    		ConversionPolicy: tasks.ConvertMissing,
    	}
    
    	model, err := tasks.Load[textencoding.Interface](tasksConfig)
    	if err != nil {
    		pterm.Error.Println("Error loading model...")
    		return err
    	}
    
    	// 3. Embedding Generation
    	pterm.Info.Println("Generating embeddings...")
    	for _, chunk := range uniqueChunks {
    		var vec []float64
    
    		encoder := func(text string) error {
    			result, err := model.Encode(context.Background(), text, int(bert.MeanPooling))
    			if err != nil {
    				return err
    			}
    
    			vec = result.Vector.Data().F64()[:limit]
    			fmt.Println(result.Vector.Data().F64()[:limit])
    			return nil
    		}
    
    		err = encoder(chunk) // Actually invoke the encoder function with the chunk
    		if err != nil {
    			pterm.Error.Println("Error encoding text...")
    			return err
    		}
    
    		embedding := estore.Embedding{
    			Word:       chunk,
    			Vector:     vec,
    			Similarity: 0.0,
    		}
    
    		db.AddEmbedding(embedding)
    	}
    
    	// Save the database to a file
    	pterm.Info.Println("Saving embeddings...")
    
    	dbPath := fmt.Sprintf("%s/embeddings.db", dataPath)
    
    	db.SaveEmbeddings(dbPath)
    
    	return nil
    }
    
    func Search(dataPath string, dbName string, prompt string, topN int) []estore.Embedding {
    	db := estore.NewEmbeddingDB()
    	dbPath := fmt.Sprintf("%s/%s", dataPath, dbName)
    	embeddings, err := db.LoadEmbeddings(dbPath)
    	if err != nil {
    		fmt.Println("Error loading embeddings:", err)
    		return nil
    	}
    
    	embeddingsModelPath := fmt.Sprintf("%s/data/models/HF/", dataPath)
    
    	model, err := tasks.Load[textencoding.Interface](&tasks.Config{
    		ModelsDir:           embeddingsModelPath,
    		ModelName:           modelName,
    		DownloadPolicy:      tasks.DownloadMissing,
    		ConversionPolicy:    tasks.ConvertMissing,
    		ConversionPrecision: tasks.F32,
    	})
    
    	if err != nil {
    		fmt.Println("Error loading model:", err)
    		return nil
    	}
    
    	var vec []float64
    	result, err := model.Encode(context.Background(), prompt, int(bert.MeanPooling))
    	if err != nil {
    		fmt.Println("Error encoding text:", err)
    		return nil
    	}
    	vec = result.Vector.Data().F64()[:limit]
    
    	embeddingForPrompt := estore.Embedding{
    		Word:       prompt,
    		Vector:     vec,
    		Similarity: 0.0,
    	}
    
    	// Retrieve the top N similar embeddings
    	topEmbeddings := estore.FindTopNSimilarEmbeddings(embeddingForPrompt, embeddings, topN)
    	if len(topEmbeddings) == 0 {
    		fmt.Println("Error finding similar embeddings.")
    		return nil
    	}
    
    	return topEmbeddings
    }

    [File Ends] pkg/embeddings/local.go

    [File Begins] pkg/embeddings/oai.go
    package embeddings
    
    import (
    	"bytes"
    	"encoding/json"
    	"eternal/pkg/documents"
    	"eternal/pkg/llm/openai"
    	"eternal/pkg/vecstore"
    	"fmt"
    	"net/http"
    	"os"
    
    	"github.com/pterm/pterm"
    )
    
    // ErrorData represents the structure of an error response from the OpenAI API.
    type ErrorData struct {
    	Code    interface{} `json:"code"`
    	Message string      `json:"message"`
    }
    
    // ErrorResponse wraps the structure of an error when an API request fails.
    type ErrorResponse struct {
    	Error ErrorData `json:"error"`
    }
    
    // EmbedRequest encapsulates the request data for the OpenAI Embeddings API.
    type EmbedRequest struct {
    	Model string `json:"model"`
    	Input string `json:"input"`
    }
    
    // EmbedResponse contains the response data from the Embeddings API.
    type EmbedResponse struct {
    	Object string       `json:"object"`
    	Data   []EmbedData  `json:"data"`
    	Model  string       `json:"model"`
    	Usage  UsageMetrics `json:"usage"`
    }
    
    // EmbedData represents a single embedding and its associated data.
    type EmbedData struct {
    	Object    string    `json:"object"`
    	Embedding []float64 `json:"embedding"`
    	Index     int       `json:"index"`
    }
    
    // UsageMetrics details the token usage of the Embeddings API request.
    type UsageMetrics struct {
    	PromptTokens int `json:"prompt_tokens"`
    	TotalTokens  int `json:"total_tokens"`
    }
    
    // GetEmbeddings interacts with the OpenAI Embeddings API to retrieve embeddings based on the provided request.
    // It returns an EmbedResponse pointer and any error encountered during the API call.
    func GetEmbeddings(req EmbedRequest) (*EmbedResponse, error) {
    	apiKey := os.Getenv("OPENAI_API_KEY")
    	if apiKey == "" {
    		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is not set")
    	}
    
    	client := openai.NewClient(apiKey)
    
    	jsonData, err := json.Marshal(req)
    	if err != nil {
    		return nil, fmt.Errorf("failed to marshal request data: %w", err)
    	}
    
    	httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/embeddings", bytes.NewBuffer(jsonData))
    	if err != nil {
    		return nil, fmt.Errorf("failed to create new HTTP request: %w", err)
    	}
    
    	httpReq.Header.Set("Content-Type", "application/json")
    	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
    
    	resp, err := client.HTTP.Do(httpReq)
    	if err != nil {
    		return nil, fmt.Errorf("failed to send request: %w", err)
    	}
    	defer resp.Body.Close()
    
    	if resp.StatusCode != http.StatusOK {
    		var errorResponse ErrorResponse
    		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
    			return nil, fmt.Errorf("failed to decode error response: %w", err)
    		}
    
    		// Additional logic to handle Code as a number or string
    		var codeStr string
    		switch code := errorResponse.Error.Code.(type) {
    		case float64:
    			codeStr = fmt.Sprintf("%.0f", code) // Convert number to string
    		case string:
    			codeStr = code
    		default:
    			return nil, fmt.Errorf("unexpected type for error code")
    		}
    
    		return nil, fmt.Errorf("API error: %s (Code: %s)", errorResponse.Error.Message, codeStr)
    	}
    
    	var embedResponse EmbedResponse
    	if err := json.NewDecoder(resp.Body).Decode(&embedResponse); err != nil {
    		return nil, fmt.Errorf("failed to decode successful response: %w", err)
    	}
    
    	return &embedResponse, nil
    }
    
    func GenerateEmbeddingOAI() {
    	if len(os.Args) < 2 {
    		fmt.Println("Usage: main.go <path_to_input_file>")
    		return
    	}
    
    	// 1. Initialization
    	pterm.Info.Println("Initializing...")
    
    	// Create a new OpenAI client
    	//client := client()
    
    	db := vecstore.NewEmbeddingDB()
    
    	// 2. Code Splitting
    	pterm.Info.Println("Splitting code...")
    	inputFilePath := os.Args[1]
    	content, err := os.ReadFile(inputFilePath)
    	if err != nil {
    		fmt.Printf("Error reading file: %v\n", err)
    		return
    	}
    
    	separators, _ := documents.GetSeparatorsForLanguage(documents.JSON)
    	// Updated the RecursiveCharacterTextSplitter to include OverlapSize and updated SplitText method
    	splitter := documents.RecursiveCharacterTextSplitter{
    		Separators:       separators,
    		KeepSeparator:    true,
    		IsSeparatorRegex: false,
    		ChunkSize:        1000,
    		LengthFunction:   func(s string) int { return len(s) },
    	}
    	chunks := splitter.SplitText(string(content))
    
    	// 3. Embedding Generation
    	pterm.Info.Println("Generating embeddings...")
    	for _, chunk := range chunks {
    		req := EmbedRequest{
    			Model: "text-embedding-ada-002",
    			Input: chunk,
    		}
    		resp, err := GetEmbeddings(req)
    		if err != nil {
    			fmt.Printf("Error getting embeddings: %v\n", err)
    			panic(err)
    		}
    
    		response := EmbedData{
    			Object:    resp.Object,
    			Embedding: resp.Data[0].Embedding,
    			Index:     resp.Data[0].Index,
    		}
    
    		embedding := vecstore.Embedding{
    			Word:       chunk,
    			Vector:     response.Embedding,
    			Similarity: 0.0,
    		}
    
    		db.AddEmbedding(embedding)
    
    	}
    
    	// Save the database to a file
    	pterm.Info.Println("Saving embeddings...")
    	db.SaveEmbeddings("./db/embeddings.db")
    
    	if len(chunks) > 0 {
    		embedding, ok := db.RetrieveEmbedding(chunks[0])
    		if ok {
    			fmt.Printf("Embedding for the first chunk:\n%v\n", embedding)
    		}
    	}
    }

    [File Ends] pkg/embeddings/oai.go

    [File Begins] pkg/hfutils/hfutils.go
    package hfutils
    
    import (
    	"crypto/sha256"
    	"encoding/hex"
    	"fmt"
    	"io"
    	"net/http"
    	"os"
    	"path/filepath"
    	"strconv"
    	"sync"
    	"time"
    )
    
    // ConcurrentDownloadManager handles downloading a file in parts concurrently.
    type ConcurrentDownloadManager struct {
    	FileName       string
    	URL            string
    	Destination    string
    	NumParts       int
    	TempDir        string
    	Sha256Checksum string
    	TotalBytes     int64
    	TotalLength    int64      // Total length of the file to be downloaded
    	BytesMutex     sync.Mutex // Mutex to protect TotalBytes
    }
    
    func (dm *ConcurrentDownloadManager) Download() error {
    	resp, err := http.Head(dm.URL)
    	if err != nil {
    		return fmt.Errorf("failed to get file info: %w", err)
    	}
    	defer resp.Body.Close()
    
    	if resp.StatusCode != http.StatusOK {
    		return fmt.Errorf("bad response: %s", resp.Status)
    	}
    
    	lengthStr := resp.Header.Get("Content-Length")
    	length, err := strconv.ParseInt(lengthStr, 10, 64)
    	if err != nil {
    		return fmt.Errorf("invalid content length: %w", err)
    	}
    	dm.TotalLength = length // Store the total length
    
    	partSize := int(length) / dm.NumParts
    	var wg sync.WaitGroup
    	var downloadErr error
    	mutex := &sync.Mutex{}
    
    	// Start progress reporting in a goroutine
    	go dm.PrintProgress()
    
    	for i := 0; i < dm.NumParts; i++ {
    		wg.Add(1)
    		start := i * partSize
    		end := start + partSize
    		if i == dm.NumParts-1 {
    			end = int(length)
    		}
    
    		go func(partNum, start, end int) {
    			defer wg.Done()
    			err := dm.downloadPart(partNum, start, end)
    			if err != nil {
    				mutex.Lock()
    				if downloadErr == nil {
    					downloadErr = err
    				}
    				mutex.Unlock()
    			}
    		}(i, start, end)
    	}
    
    	wg.Wait()
    
    	if downloadErr != nil {
    		return downloadErr
    	}
    
    	if err := dm.mergeParts(); err != nil {
    		return fmt.Errorf("failed to merge parts: %w", err)
    	}
    
    	if dm.Sha256Checksum != "" {
    		if err := dm.verifyChecksum(); err != nil {
    			return err
    		}
    	}
    
    	return nil
    }
    
    func (dm *ConcurrentDownloadManager) downloadPart(partNum, start, end int) error {
    	req, err := http.NewRequest("GET", dm.URL, nil)
    	if err != nil {
    		return err
    	}
    
    	req.Header.Set("Range", fmt.Sprintf("bytes=%d-%d", start, end-1))
    	resp, err := http.DefaultClient.Do(req)
    	if err != nil {
    		return err
    	}
    	defer resp.Body.Close()
    
    	if resp.StatusCode != http.StatusPartialContent {
    		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
    	}
    
    	partPath := filepath.Join(dm.TempDir, fmt.Sprintf("part-%d", partNum))
    	partFile, err := os.Create(partPath)
    	if err != nil {
    		return err
    	}
    	defer partFile.Close()
    
    	buf := make([]byte, 32*1024) // 32 KB buffer
    	for {
    		n, err := resp.Body.Read(buf)
    		if n > 0 {
    			_, writeErr := partFile.Write(buf[:n])
    			if writeErr != nil {
    				return writeErr
    			}
    			dm.BytesMutex.Lock()
    			dm.TotalBytes += int64(n)
    			dm.BytesMutex.Unlock()
    		}
    		if err != nil {
    			if err == io.EOF {
    				break
    			}
    			return err
    		}
    	}
    
    	return nil
    }
    
    func (dm *ConcurrentDownloadManager) mergeParts() error {
    	finalFile, err := os.Create(dm.Destination)
    	if err != nil {
    		return err
    	}
    	defer finalFile.Close()
    
    	for i := 0; i < dm.NumParts; i++ {
    		partPath := filepath.Join(dm.TempDir, fmt.Sprintf("part-%d", i))
    		partFile, err := os.Open(partPath)
    		if err != nil {
    			return err
    		}
    
    		if _, err := io.Copy(finalFile, partFile); err != nil {
    			partFile.Close()
    			return err
    		}
    
    		partFile.Close()
    		os.Remove(partPath) // Cleanup part file after merge
    	}
    	return nil
    }
    
    func (dm *ConcurrentDownloadManager) verifyChecksum() error {
    	finalPath := filepath.Join(dm.TempDir, dm.FileName)
    	file, err := os.Open(finalPath)
    	if err != nil {
    		return err
    	}
    	defer file.Close()
    
    	hasher := sha256.New()
    	if _, err := io.Copy(hasher, file); err != nil {
    		return err
    	}
    
    	actualChecksum := hex.EncodeToString(hasher.Sum(nil))
    	if actualChecksum != dm.Sha256Checksum {
    		return fmt.Errorf("checksum mismatch: expected %s, got %s", dm.Sha256Checksum, actualChecksum)
    	}
    
    	return nil
    }
    
    func (dm *ConcurrentDownloadManager) PrintProgress() {
    	ticker := time.NewTicker(500 * time.Millisecond)
    	defer ticker.Stop()
    
    	for range ticker.C {
    		dm.BytesMutex.Lock()
    		downloaded := dm.TotalBytes
    		dm.BytesMutex.Unlock()
    
    		percent := float64(downloaded) / float64(dm.TotalLength) * 100
    		fmt.Printf("\rDownloading... %.2f%% complete (%d of %d bytes)", percent, downloaded, dm.TotalLength)
    
    		if downloaded >= dm.TotalLength {
    			break
    		}
    	}
    }

    [File Ends] pkg/hfutils/hfutils.go

    [File Begins] pkg/jobs/jobs.go
    package main
    
    import (
    	"errors"
    	"fmt"
    	"os"
    	"sync"
    	"time"
    )
    
    // Constants for different job statuses.
    const (
    	queued    = "queued"
    	running   = "running"
    	completed = "completed"
    	failed    = "failed"
    )
    
    var (
    	jobQueue     = make(chan *Job, 100)
    	jobStatusMap = make(map[string]*jobStatus)
    	mutex        = &sync.Mutex{}
    	wg           sync.WaitGroup
    )
    
    // JobType represents different types of jobs that can be handled.
    type JobType int
    
    // Enumeration of different JobTypes.
    const (
    	WriteTimeToFile JobType = iota
    	AnotherJobType
    	// Future job types should be added here.
    )
    
    // Job defines the structure of a job including its type, payload, and callback.
    type Job struct {
    	ID       string
    	JobType  JobType
    	Payload  interface{}
    	Callback func(result interface{}, err error)
    }
    
    // jobStatus represents the current status of a job along with its result or error.
    type jobStatus struct {
    	Status string
    	Result interface{}
    	Error  error
    }
    
    // worker is a goroutine that processes jobs from the jobQueue.
    func worker() {
    	for job := range jobQueue {
    		processJob(job)
    	}
    }
    
    // processJob handles the execution and updating of the job status.
    func processJob(job *Job) {
    	status := &jobStatus{Status: running}
    	mutex.Lock()
    	jobStatusMap[job.ID] = status
    	mutex.Unlock()
    
    	result, err := executeJob(job)
    	updateJobStatus(job.ID, result, err)
    
    	job.Callback(result, err)
    	wg.Done()
    }
    
    // executeJob executes the given job based on its JobType.
    func executeJob(job *Job) (interface{}, error) {
    	switch job.JobType {
    	case WriteTimeToFile:
    		return writeTimeToFile(job.Payload)
    	case AnotherJobType:
    		return anotherJobFunction(job.Payload)
    	default:
    		return nil, errors.New("unknown job type")
    	}
    }
    
    // SubmitJob adds a job to the jobQueue and tracks its status.
    func SubmitJob(job *Job) {
    	wg.Add(1)
    	mutex.Lock()
    	jobStatusMap[job.ID] = &jobStatus{Status: queued}
    	mutex.Unlock()
    	jobQueue <- job
    }
    
    // GetJobStatus returns the status of a job by its ID.
    func GetJobStatus(jobID string) *jobStatus {
    	mutex.Lock()
    	defer mutex.Unlock()
    	return jobStatusMap[jobID]
    }
    
    // GetAllJobsStatus returns the status of all jobs.
    func GetAllJobsStatus() map[string]*jobStatus {
    	mutex.Lock()
    	defer mutex.Unlock()
    	return jobStatusMap
    }
    
    // InitWorkers initializes a specified number of worker goroutines.
    func InitWorkers(numWorkers int) {
    	for i := 0; i < numWorkers; i++ {
    		go worker()
    	}
    }
    
    // updateJobStatus updates the status of a job in the jobStatusMap.
    func updateJobStatus(jobID string, result interface{}, err error) {
    	mutex.Lock()
    	defer mutex.Unlock()
    
    	status := jobStatusMap[jobID]
    	if err != nil {
    		status.Status = failed
    		status.Error = err
    	} else {
    		status.Status = completed
    		status.Result = result
    	}
    }
    
    // writeTimeToFile handles the specific logic for writing the time to a file.
    // This is a function to test the jobs system and not meant for production use.
    func writeTimeToFile(payload interface{}) (interface{}, error) {
    	waitTime, ok := payload.(time.Duration)
    	if !ok {
    		return nil, errors.New("invalid payload")
    	}
    
    	// Wait for the specified time
    	fmt.Println("Waiting for", waitTime)
    	time.Sleep(waitTime)
    
    	// Get the current time
    	currentTime := time.Now().Format(time.RFC3339)
    
    	// Create a file in /tmp with the current time
    	filePath := fmt.Sprintf("/tmp/job_%s.txt", currentTime)
    	file, err := os.Create(filePath)
    	if err != nil {
    		return nil, fmt.Errorf("failed to create file: %v", err)
    	}
    	defer file.Close()
    
    	// Write the current time to the file
    	_, err = file.WriteString(currentTime)
    	if err != nil {
    		return nil, fmt.Errorf("failed to write to file: %v", err)
    	}
    
    	return "File created with time: " + currentTime, nil
    }
    
    // anotherJobFunction represents a placeholder for future job types.
    func anotherJobFunction(payload interface{}) (interface{}, error) {
    	// Code for another type of job
    	// ...
    	return nil, errors.New("not implemented")
    }

    [File Ends] pkg/jobs/jobs.go

      [File Begins] pkg/llm/anthropic/anthropic.go
      package anthropic
      
      import (
      	"bytes"
      	"encoding/json"
      	"net/http"
      )
      
      const (
      	baseURL             = "https://api.anthropic.com/v1"
      	completionsEndpoint = "/messages"
      )
      
      // SendRequest sends a request to the Anthropic API and decodes the response.
      func SendRequest(endpoint string, payload interface{}, apiKey string) (*http.Response, error) {
      	jsonData, err := json.Marshal(payload)
      	if err != nil {
      		return nil, err
      	}
      
      	req, err := http.NewRequest("POST", baseURL+endpoint, bytes.NewBuffer(jsonData))
      	if err != nil {
      		return nil, err
      	}
      
      	req.Header.Set("Content-Type", "application/json")
      	req.Header.Set("x-api-key", apiKey)
      	req.Header.Set("anthropic-version", "2023-06-01")
      
      	return http.DefaultClient.Do(req)
      }

      [File Ends] pkg/llm/anthropic/anthropic.go

      [File Begins] pkg/llm/anthropic/completions.go
      package anthropic
      
      import (
      	"bufio"
      	"bytes"
      	"encoding/json"
      	"eternal/pkg/web"
      	"fmt"
      	"strconv"
      	"strings"
      
      	"github.com/gofiber/websocket/v2"
      	"github.com/pterm/pterm"
      )
      
      type Message struct {
      	Role    string `json:"role"`
      	Content string `json:"content"`
      }
      
      type CompletionRequest struct {
      	Model       string    `json:"model"`
      	MaxTokens   int       `json:"max_tokens"`
      	Messages    []Message `json:"messages"`
      	Stream      bool      `json:"stream"`
      	Temperature float64   `json:"temperature"`
      }
      
      type CompletionResponse struct {
      	ID      string `json:"id"`
      	Type    string `json:"type"`
      	Role    string `json:"role"`
      	Content []struct {
      		Type string `json:"type"`
      		Text string `json:"text"`
      	} `json:"content"`
      	Model        string      `json:"model"`
      	StopReason   string      `json:"stop_reason"`
      	StopSequence interface{} `json:"stop_sequence"`
      	Usage        struct {
      		InputTokens  int `json:"input_tokens"`
      		OutputTokens int `json:"output_tokens"`
      	} `json:"usage"`
      }
      
      type ContentBlockDelta struct {
      	Type  string    `json:"type"`
      	Index int       `json:"index"`
      	Delta TextDelta `json:"delta"`
      }
      
      type TextDelta struct {
      	Type string `json:"type"`
      	Text string `json:"text"`
      }
      
      func StreamCompletionToWebSocket(c *websocket.Conn, chatID int, model string, messages []Message, temperature float64, apiKey string) error {
      
      	payload := &CompletionRequest{
      		Model:       model,
      		MaxTokens:   4096,
      		Stream:      true,
      		Messages:    messages,
      		Temperature: temperature,
      	}
      
      	resp, err := SendRequest(completionsEndpoint, payload, apiKey)
      	if err != nil {
      		return err
      	}
      	defer resp.Body.Close()
      
      	// Handle streaming response
      	msgBuffer := new(bytes.Buffer)
      	scanner := bufio.NewScanner(resp.Body)
      	for scanner.Scan() {
      		line := scanner.Text()
      		fmt.Println(line)
      		if strings.HasPrefix(line, "data: ") {
      			line = strings.TrimPrefix(line, "data: ")
      
      			// Unmarshal the JSON response
      			var data ContentBlockDelta
      			if err := json.Unmarshal([]byte(line), &data); err != nil {
      				return err
      			}
      
      			msgBuffer.WriteString(data.Delta.Text)
      
      			htmlMsg := web.MarkdownToHTML(msgBuffer.Bytes())
      
      			turnIDStr := strconv.Itoa(chatID)
      
      			formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>\n<codapi-snippet engine='browser' sandbox='javascript' editor='basic'></codapi-snippet>", turnIDStr, htmlMsg)
      
      			if err := c.WriteMessage(websocket.TextMessage, []byte(formattedContent)); err != nil {
      				pterm.Error.Println("WebSocket write error:", err)
      				return err
      			}
      		}
      	}
      
      	if err := scanner.Err(); err != nil {
      		pterm.Error.Println("Error reading stream:", err)
      		return err
      	}
      
      	return nil
      
      }

      [File Ends] pkg/llm/anthropic/completions.go

    [File Begins] pkg/llm/gguf.go
    package llm
    
    import (
    	"bufio"
    	"bytes"
    	_ "embed"
    	"strings"
    
    	"eternal/pkg/web"
    	"fmt"
    	"io"
    	"log"
    	"os/exec"
    	"path/filepath"
    
    	"github.com/gofiber/fiber/v2"
    	"github.com/gofiber/websocket/v2"
    	"github.com/pterm/pterm"
    )
    
    type CommandOutput struct {
    	Output       string `json:"output"`
    	Finished     string `json:"finished"`
    	SocketNumber string `json:"socketNumber"`
    	ModelName    string `json:"modelName"`
    }
    
    // Options represents all the command line options
    type GGUFOptions struct {
    	ResponseDelimiter     string  `json:"responseDelimiter"`
    	Help                  bool    `json:"help"`
    	Version               bool    `json:"version"`
    	Interactive           bool    `json:"interactive"`
    	InteractiveFirst      bool    `json:"interactive_first"`
    	Instruct              bool    `json:"instruct"`
    	ChatML                bool    `json:"chatml"`
    	MultilineInput        bool    `json:"multiline_input"`
    	ReversePrompt         string  `json:"reverse_prompt"`
    	Color                 bool    `json:"color"`
    	Seed                  int     `json:"seed"`
    	Threads               int     `json:"threads"`
    	ThreadsBatch          int     `json:"threads_batch"`
    	ThreadsDraft          int     `json:"threads_draft"`
    	ThreadsBatchDraft     int     `json:"threads_batch_draft"`
    	Prompt                string  `json:"prompt"`
    	Escape                bool    `json:"escape"`
    	PromptCache           string  `json:"prompt_cache"`
    	PromptCacheAll        bool    `json:"prompt_cache_all"`
    	PromptCacheRO         bool    `json:"prompt_cache_ro"`
    	RandomPrompt          bool    `json:"random_prompt"`
    	InPrefixBOS           bool    `json:"in_prefix_bos"`
    	InPrefix              string  `json:"in_prefix"`
    	InSuffix              string  `json:"in_suffix"`
    	File                  string  `json:"file"`
    	NPredict              int     `json:"n_predict"`
    	CtxSize               int     `json:"ctx_size"`
    	BatchSize             int     `json:"batch_size"`
    	Samplers              string  `json:"samplers"`
    	SamplingSeq           string  `json:"sampling_seq"`
    	TopK                  int     `json:"top_k"`
    	TopP                  float64 `json:"top_p"`
    	MinP                  float64 `json:"min_p"`
    	TFS                   float64 `json:"tfs"`
    	Typical               float64 `json:"typical"`
    	RepeatLastN           int     `json:"repeat_last_n"`
    	RepeatPenalty         float64 `json:"repeat_penalty"`
    	PresencePenalty       float64 `json:"presence_penalty"`
    	FrequencyPenalty      float64 `json:"frequency_penalty"`
    	Mirostat              int     `json:"mirostat"`
    	MirostatLR            float64 `json:"mirostat_lr"`
    	MirostatEnt           float64 `json:"mirostat_ent"`
    	LogitBias             string  `json:"logit_bias"`
    	Grammar               string  `json:"grammar"`
    	GrammarFile           string  `json:"grammar_file"`
    	CFGNegativePrompt     string  `json:"cfg_negative_prompt"`
    	CFGNegativePromptFile string  `json:"cfg_negative_prompt_file"`
    	CFGScale              float64 `json:"cfg_scale"`
    	RopeScaling           string  `json:"rope_scaling"`
    	RopeScale             int     `json:"rope_scale"`
    	RopeFreqBase          int     `json:"rope_freq_base"`
    	RopeFreqScale         int     `json:"rope_freq_scale"`
    	YarnOrigCtx           int     `json:"yarn_orig_ctx"`
    	YarnExtFactor         float64 `json:"yarn_ext_factor"`
    	YarnAttnFactor        float64 `json:"yarn_attn_factor"`
    	YarnBetaSlow          float64 `json:"yarn_beta_slow"`
    	YarnBetaFast          float64 `json:"yarn_beta_fast"`
    	IgnoreEOS             bool    `json:"ignore_eos"`
    	NoPenalizeNL          bool    `json:"no_penalize_nl"`
    	Temp                  float64 `json:"temp"`
    	LogitsAll             bool    `json:"logits_all"`
    	Hellaswag             bool    `json:"hellaswag"`
    	HellaswagTasks        int     `json:"hellaswag_tasks"`
    	Winogrande            bool    `json:"winogrande"`
    	WinograndeTasks       int     `json:"winogrande_tasks"`
    	Keep                  int     `json:"keep"`
    	Draft                 int     `json:"draft"`
    	Chunks                int     `json:"chunks"`
    	Parallel              int     `json:"parallel"`
    	Sequences             int     `json:"sequences"`
    	PAccept               float64 `json:"p_accept"`
    	PSplit                float64 `json:"p_split"`
    	ContBatching          bool    `json:"cont_batching"`
    	MMProj                string  `json:"mmproj"`
    	Image                 string  `json:"image"`
    	Mlock                 bool    `json:"mlock"`
    	NoMmap                bool    `json:"no_mmap"`
    	Numa                  bool    `json:"numa"`
    	NGPULayers            int     `json:"n_gpu_layers"`
    	NGPULayersDraft       int     `json:"n_gpu_layers_draft"`
    	SplitMode             string  `json:"split_mode"`
    	TensorSplit           string  `json:"tensor_split"`
    	MainGPU               int     `json:"main_gpu"`
    	VerbosePrompt         bool    `json:"verbose_prompt"`
    	NoDisplayPrompt       bool    `json:"no_display_prompt"`
    	GrpAttnN              int     `json:"grp_attn_n"`
    	GrpAttnW              float64 `json:"grp_attn_w"`
    	DumpKVCache           bool    `json:"dump_kv_cache"`
    	NoKVOffload           bool    `json:"no_kv_offload"`
    	CacheTypeK            string  `json:"cache_type_k"`
    	CacheTypeV            string  `json:"cache_type_v"`
    	SimpleIO              bool    `json:"simple_io"`
    	Lora                  string  `json:"lora"`
    	LoraScaled            string  `json:"lora_scaled"`
    	LoraBase              string  `json:"lora_base"`
    	Model                 string  `json:"model"`
    	ModelDraft            string  `json:"model_draft"`
    	LogDir                string  `json:"logdir"`
    	OverrideKV            string  `json:"override_kv"`
    	PrintTokenCount       int     `json:"print_token_count"`
    	// log options
    	LogTest    bool   `json:"log_test"`
    	LogDisable bool   `json:"log_disable"`
    	LogEnable  bool   `json:"log_enable"`
    	LogFile    string `json:"log_file"`
    	LogNew     bool   `json:"log_new"`
    	LogAppend  bool   `json:"log_append"`
    }
    
    func BuildCommand(cmdPath string, options GGUFOptions) *exec.Cmd {
    	execPath := filepath.Join(cmdPath, "gguf/main")
    	//cachePath := filepath.Join(cmdPath, "cache")
    
    	ctxSize := fmt.Sprintf("%d", options.CtxSize)
    	temp := fmt.Sprintf("%f", options.Temp)
    	repeatPenalty := fmt.Sprintf("%f", options.RepeatPenalty)
    	topP := fmt.Sprintf("%f", options.TopP)
    	topK := fmt.Sprintf("%d", options.TopK)
    
    	cmdArgs := []string{
    		"--no-display-prompt",
    		"-m", options.Model,
    		"-p", options.Prompt,
    		"-c", ctxSize, // 0 = loaded from model
    		"--n-predict", "-2", // -1 = infinity, -2 = until context filled
    		"--repeat-penalty", repeatPenalty,
    		"--top-p", topP,
    		"--top-k", topK,
    		"--n-gpu-layers", fmt.Sprintf("%d", options.NGPULayers),
    		"--reverse-prompt", "<|eot_id|>",
    		"--multiline-input",
    		"--temp", temp,
    		//--dynatemp-range", "0.5", // 0.0 = disabled
    		"--flash-attn", // enable flash attention, default disabled
    		// "--mlock",
    		"--seed", "-1",
    		//"--ignore-eos",
    		//"--no-mmap",
    		//"--simple-io",
    		//"--keep", "2048",
    		//"--prompt-cache", cachePath,
    		//"--prompt-cache-all",
    		//"--grammar-file", "./json.gbnf",
    		//"--override-kv", "llama.expert_used_count=int:3", // mixtral only
    		//"--override-kv", "tokenizer.ggml.pre=str:llama3",
    	}
    
    	return exec.Command(execPath, cmdArgs...)
    }
    
    // Upgrades the HTTP connection to a WebSocket connection
    func UpgradeToWebSocket(c *fiber.Ctx) error {
    	if websocket.IsWebSocketUpgrade(c) {
    		return c.Next()
    	}
    	return c.Status(fiber.StatusUpgradeRequired).SendString("Upgrade required")
    }
    
    func CompletionWebSocket(c *websocket.Conn, cmdPath string) {
    	// Use websocket.Conn's `ReadJSON` and `WriteJSON` for communication
    	var args GGUFOptions
    	if err := c.ReadJSON(&args); err != nil {
    		log.Println("Invalid input:", err)
    		c.WriteJSON(fiber.Map{"error": "Invalid input"})
    		return
    	}
    
    	cmd := BuildCommand(cmdPath, args)
    
    	stdout, err := cmd.StdoutPipe()
    	if err != nil {
    		log.Println("Failed to set up command output:", err)
    		c.WriteJSON(fiber.Map{"error": "Failed to set up command output"})
    		return
    	}
    
    	if err := cmd.Start(); err != nil {
    		log.Println("Error starting command:", err)
    		c.WriteJSON(fiber.Map{"error": "Error starting command"})
    		return
    	}
    
    	buf := make([]byte, 1024)
    	for {
    		n, err := stdout.Read(buf)
    		if err != nil {
    			log.Println("Error reading command output:", err)
    			c.WriteJSON(fiber.Map{"error": "Error reading command output"})
    			return
    		}
    		if n > 0 {
    			data := CommandOutput{Output: string(buf[:n])}
    			if err := c.WriteJSON(data); err != nil {
    				log.Println("Error encoding JSON:", err)
    				c.WriteJSON(fiber.Map{"error": "Error encoding JSON"})
    				return
    			}
    		} else {
    			break
    		}
    	}
    
    	if err := cmd.Wait(); err != nil {
    		log.Println("Command finished with error:", err)
    		c.WriteJSON(fiber.Map{"error": "Command finished with error"})
    	}
    }
    
    // MakeCompletionWebSocket creates a closure that captures model parameters and returns a WebSocket handler.
    func MakeCompletionWebSocket(c websocket.Conn, chatID int, modelOpts *GGUFOptions, dataPath string) error {
    	defer c.Close()
    	var msgBuffer bytes.Buffer // Buffer to accumulate messages
    
    	// // Get the model name from its file name
    	// modelPath := filepath.Base(modelOpts.Model)
    	// modelName := strings.TrimSuffix(modelPath, filepath.Ext(modelPath))
    
    	// responsePrefix := "### Response from " + modelName + "\n"
    
    	// // Store the prompt in the buffer
    	// msgBuffer.WriteString(responsePrefix)
    
    	for {
    		cmd := BuildCommand(dataPath, *modelOpts)
    
    		stdout, err := cmd.StdoutPipe()
    		if err != nil {
    			return err
    		}
    
    		if err = cmd.Start(); err != nil {
    			return err
    		}
    
    		reader := bufio.NewReader(stdout)
    		for {
    			line, err := reader.ReadString('\n')
    			if err != nil {
    				if err == io.EOF {
    					return fmt.Errorf("%s", msgBuffer.String())
    				}
    
    				return err
    			}
    
    			msgBuffer.WriteString(line)
    
    			// Convert the buffer content to HTML
    			htmlMsg := web.MarkdownToHTML(msgBuffer.Bytes())
    
    			// Convert chatID to string for formatting
    			turnIDStr := fmt.Sprint(chatID + TurnCounter)
    
    			// Send the accumulated content
    			// formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>\n<codapi-snippet url='http://localhost:1313/v1/exec' sandbox='go' editor='external'></codapi-snippet>", turnIDStr, htmlMsg)
    			//formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>\n<codapi-snippet engine='browser' sandbox='javascript' editor='basic'></codapi-snippet>", turnIDStr, htmlMsg, turnIDStr)
    			formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1, rounded-2' hx-trigger='load'>%s</div><codapi-snippet engine='browser' sandbox='javascript' editor='basic'></codapi-snippet>", turnIDStr, htmlMsg)
    			if err := c.WriteMessage(websocket.TextMessage, []byte(formattedContent)); err != nil {
    				pterm.Error.Println("WebSocket write error:", err)
    				return err
    			}
    		}
    	}
    }
    
    func LMResponse(c websocket.Conn, chatID int, modelOpts *GGUFOptions, dataPath string) error {
    	defer c.Close()
    	var msgBuffer bytes.Buffer // Buffer to accumulate messages
    
    	// Remove the system instructions from the prompt
    	sprompt := strings.Split(modelOpts.Prompt, "New Instructions:")
    	sprompt = sprompt[:1]
    
    	// Join the prompt without the system instructions
    	prompt := strings.Join(sprompt, "\n")
    
    	//pterm.Error.Println(prompt)
    
    	// Store the prompt in the buffer
    	msgBuffer.WriteString(prompt + "\n")
    
    	for {
    		cmd := BuildCommand(dataPath, *modelOpts)
    
    		stdout, err := cmd.StdoutPipe()
    		if err != nil {
    			return err
    		}
    
    		if err = cmd.Start(); err != nil {
    			return err
    		}
    
    		reader := bufio.NewReader(stdout)
    		for {
    			line, err := reader.ReadString('\n')
    			if err != nil {
    				if err == io.EOF {
    					return fmt.Errorf("%s", msgBuffer.String())
    				}
    
    				return err
    			}
    
    			msgBuffer.WriteString(line)
    
    			// Convert the buffer content to HTML
    			htmlMsg := web.MarkdownToHTML(msgBuffer.Bytes())
    
    			// Convert chatID to string for formatting
    			turnIDStr := fmt.Sprint(chatID + TurnCounter)
    
    			// Send the accumulated content
    			// formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>\n<codapi-snippet url='http://localhost:1313/v1/exec' sandbox='go' editor='external'></codapi-snippet>", turnIDStr, htmlMsg)
    			//formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>\n<codapi-snippet engine='browser' sandbox='javascript' editor='basic'></codapi-snippet>", turnIDStr, htmlMsg, turnIDStr)
    			formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1, rounded-2' hx-trigger='load'>%s</div><codapi-snippet engine='browser' sandbox='javascript' editor='basic'></codapi-snippet>", turnIDStr, htmlMsg)
    			if err := c.WriteMessage(websocket.TextMessage, []byte(formattedContent)); err != nil {
    				pterm.Error.Println("WebSocket write error:", err)
    				return err
    			}
    		}
    	}
    }

    [File Ends] pkg/llm/gguf.go

      [File Begins] pkg/llm/google/google.go
      package google
      
      import (
      	"bytes"
      	"context"
      	"fmt"
      
      	"eternal/pkg/llm"
      	"eternal/pkg/web"
      
      	"github.com/gofiber/websocket/v2"
      	"github.com/google/generative-ai-go/genai"
      	"github.com/pterm/pterm"
      	"golang.org/x/text/language"
      	"golang.org/x/text/message"
      	"google.golang.org/api/iterator"
      	"google.golang.org/api/option"
      )
      
      const (
      	model = "models/gemini-1.5-pro-latest"
      )
      
      // StreamGeminiResponseToWebSocket streams the response from the Gemini API to a WebSocket connection.
      func StreamGeminiResponseToWebSocket(c *websocket.Conn, chatID int, prompt string, apiKey string) error {
      	pterm.Warning.Printfln("Using model: %s", model)
      	ctx := context.Background()
      	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
      	if err != nil {
      		pterm.Error.Println(err)
      		return err
      	}
      	defer client.Close()
      
      	pterm.Warning.Printfln("Sending prompt to api...")
      	generativeModel := client.GenerativeModel(model)
      
      	// Configure model parameters by invoking Set* methods on the model.
      	generativeModel.SetTemperature(0.1)
      	generativeModel.SetTopK(1)
      	generativeModel.SetTopP(1)
      
      	pterm.Warning.Printfln("Generating content stream...")
      	iter := generativeModel.GenerateContentStream(ctx, genai.Text(prompt))
      
      	msgBuffer := new(bytes.Buffer)
      	for {
      		resp, err := iter.Next()
      		if err == iterator.Done {
      			return err
      		}
      		if err != nil {
      			pterm.Error.Println(err)
      			return err
      		}
      
      		// Access the Content field of the genai.Part type
      		p := message.NewPrinter(language.English)
      		content := p.Sprintf("%s", resp.Candidates[0].Content.Parts[0])
      
      		msgBuffer.WriteString(content)
      
      		htmlMsg := web.MarkdownToHTML(msgBuffer.Bytes())
      
      		turnIDStr := fmt.Sprint(chatID + llm.TurnCounter)
      
      		formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>\n<codapi-snippet engine='browser' sandbox='javascript' editor='basic'></codapi-snippet>", turnIDStr, htmlMsg)
      
      		if err := c.WriteMessage(websocket.TextMessage, []byte(formattedContent)); err != nil {
      			pterm.Error.Println("WebSocket write error:", err)
      			return err
      		}
      	}
      }

      [File Ends] pkg/llm/google/google.go

    [File Begins] pkg/llm/llm.go
    package llm
    
    import (
    	"fmt"
    	"io"
    	"net/http"
    	"os"
    	"path/filepath"
    	"sync"
    
    	"github.com/pterm/pterm"
    )
    
    var (
    	downloadProgressMap = make(map[string]DownloadProgress)
    	progressMutex       sync.Mutex
    	TurnCounter         = 0
    )
    
    // DownloadProgress structure to hold download progress information
    type DownloadProgress struct {
    	Total   int64 `json:"total"`
    	Current int64 `json:"current"`
    }
    
    // WebModel defines the interface for a model that interacts over the web.
    type WebModel interface {
    	Connect(endpoint string) (*http.Response, error)
    }
    
    // ModelManager defines the interface for managing local models.
    type ModelManager interface {
    	GetConfig() error
    	Download() error
    	Delete() error
    }
    
    type Model struct {
    	Name      string   `yaml:"name"`
    	Homepage  string   `yaml:"homepage"`
    	Prompt    string   `yaml:"prompt"`
    	Ctx       int      `yaml:"ctx"`
    	Roles     []string `yaml:"roles"`
    	Tags      []string `yaml:"tags,omitempty"`
    	GGUF      string   `yaml:"gguf,omitempty"`
    	Downloads []string `yaml:"downloads,omitempty"`
    	LocalPath string   `yaml:"localPath,omitempty"`
    }
    
    type ProgressReader struct {
    	Reader        io.Reader
    	ProgressBar   *pterm.ProgressbarPrinter
    	TotalRead     int64
    	ContentLength int64
    }
    
    func (pr *ProgressReader) Read(p []byte) (int, error) {
    	n, err := pr.Reader.Read(p)
    	pr.TotalRead += int64(n)
    
    	// Calculate progress percentage
    	//progressPercentage := int64(100.0 * float64(pr.TotalRead) / float64(pr.ContentLength))
    
    	// Update the progress map safely
    	progressMutex.Lock()
    	downloadProgressMap["sse-progress"] = DownloadProgress{
    		Total:   pr.ContentLength,
    		Current: pr.TotalRead,
    	}
    	progressMutex.Unlock()
    
    	pr.ProgressBar.Add(n)
    
    	return n, err
    }
    
    // Download downloads a file from a URL to a local path, resuming if possible.
    func Download(url string, localPath string) error {
    	dir := filepath.Dir(localPath)
    	// Ensure the directory exists
    	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
    		return fmt.Errorf("failed to create directory: %w", err)
    	}
    
    	// Open the local file for writing, create it if not exists
    	out, err := os.OpenFile(localPath, os.O_RDWR|os.O_CREATE, 0666)
    	if err != nil {
    		return fmt.Errorf("failed to open file: %w", err)
    	}
    	defer out.Close()
    
    	// Find out how much has already been downloaded
    	fi, err := out.Stat()
    	if err != nil {
    		return fmt.Errorf("failed to stat file: %w", err)
    	}
    	size := fi.Size()
    
    	// Create a new HTTP request
    	req, err := http.NewRequest("GET", url, nil)
    	if err != nil {
    		return fmt.Errorf("failed to create request: %w", err)
    	}
    
    	// Set the Range header to request the portion of the file we don't have yet
    	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", size))
    
    	// Make the HTTP request
    	resp, err := http.DefaultClient.Do(req)
    	if err != nil {
    		return fmt.Errorf("failed to start file download: %w", err)
    	}
    	defer resp.Body.Close()
    
    	if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusOK {
    		return fmt.Errorf("bad status getting file: %s", resp.Status)
    	}
    
    	// Seek to the end of the file to start appending data
    	_, err = out.Seek(0, io.SeekEnd)
    	if err != nil {
    		return fmt.Errorf("failed to seek file: %w", err)
    	}
    
    	pterm.Info.Printf("Downloading model:\nURL: %s\nFile: %s\n", url, localPath)
    
    	// Initialize the progress bar
    	progressBar, _ := pterm.DefaultProgressbar.WithTotal(int(resp.ContentLength + size)).WithTitle("Downloading").Start()
    
    	// Wrap the response body in a custom reader that updates the progress bar
    	progressReader := &ProgressReader{
    		Reader:        resp.Body,
    		ProgressBar:   progressBar,
    		TotalRead:     size,
    		ContentLength: resp.ContentLength + size,
    	}
    
    	// Copy the remaining data to the file, updating progress along the way
    	_, err = io.Copy(out, progressReader)
    	if err != nil {
    		return err
    	}
    
    	// Ensure the progress bar reflects the complete download
    	progressBar.Total = (int(progressReader.TotalRead))
    
    	// Finish the progress bar
    	progressBar.Stop()
    
    	// Update the model's downloaded state in the database
    	// err = UpdateModelDownloadedState(modelName, true)
    	// if err != nil {
    	//     log.Errorf("Failed to update model downloaded state: %v", err)
    	// }
    
    	return nil
    }
    
    func GetDownloadProgress(key string) string {
    	progressMutex.Lock()
    	defer progressMutex.Unlock()
    
    	// Return the progress for a specific key
    	if progress, ok := downloadProgressMap[key]; ok {
    		// Calculate progress percentage
    		return fmt.Sprintf("%d%%", int64(100.0*float64(progress.Current)/float64(progress.Total)))
    	}
    
    	// Fallback if no progress is found
    	return "0%"
    }
    
    // GetExpectedFileSize returns the expected file size of a download.
    func GetExpectedFileSize(url string) (int64, error) {
    	// Create a new HTTP request
    	req, err := http.NewRequest("HEAD", url, nil)
    	if err != nil {
    		return 0, fmt.Errorf("failed to create request: %w", err)
    	}
    
    	// Make the HTTP request
    	resp, err := http.DefaultClient.Do(req)
    	if err != nil {
    		return 0, fmt.Errorf("failed to start file download: %w", err)
    	}
    	defer resp.Body.Close()
    
    	if resp.StatusCode != http.StatusOK {
    		return 0, fmt.Errorf("bad status getting file: %s", resp.Status)
    	}
    
    	return resp.ContentLength, nil
    }
    
    func (m *Model) Delete() error {
    	if err := os.Remove(m.LocalPath); err != nil {
    		return fmt.Errorf("failed to delete model: %w", err)
    	}
    
    	return nil
    }
    
    func IncrementTurnCounter() {
    	TurnCounter++
    }
    
    func GetTurnCounter() int {
    	return TurnCounter
    }

    [File Ends] pkg/llm/llm.go

      [File Begins] pkg/llm/openai/completions.go
      package openai
      
      import (
      	"bufio"
      	"bytes"
      	"encoding/json"
      	"fmt"
      	"net/http"
      	"strings"
      
      	"eternal/pkg/llm"
      	"eternal/pkg/web"
      
      	"github.com/gofiber/websocket/v2"
      	"github.com/pterm/pterm"
      )
      
      const (
      	baseURL             = "https://api.openai.com/v1"
      	completionsEndpoint = "/chat/completions"
      )
      
      // SendRequest sends a request to the OpenAI API and decodes the response.
      func SendRequest(endpoint string, payload interface{}, apiKey string) (*http.Response, error) {
      	jsonData, err := json.Marshal(payload)
      	if err != nil {
      		return nil, err
      	}
      
      	req, err := http.NewRequest("POST", baseURL+endpoint, bytes.NewBuffer(jsonData))
      	if err != nil {
      		return nil, err
      	}
      
      	req.Header.Set("Content-Type", "application/json")
      	req.Header.Set("Authorization", "Bearer "+apiKey)
      
      	return http.DefaultClient.Do(req)
      }
      
      func StreamCompletionToWebSocket(c *websocket.Conn, chatID int, model string, messages []llm.Message, temperature float64, apiKey string) error {
      	payload := &CompletionRequest{
      		Model:       model,
      		Messages:    messages,
      		Temperature: temperature,
      		Stream:      true,
      	}
      
      	resp, err := SendRequest(completionsEndpoint, payload, apiKey)
      	if err != nil {
      		pterm.Error.Println(err)
      		return err
      	}
      	defer resp.Body.Close()
      
      	// Handle streaming response
      	msgBuffer := new(bytes.Buffer)
      	scanner := bufio.NewScanner(resp.Body)
      	for scanner.Scan() {
      		line := scanner.Text()
      		if strings.HasPrefix(line, "data: ") {
      			jsonStr := line[6:] // Strip the "data: " prefix
      			var data struct {
      				Choices []struct {
      					Delta struct {
      						Content string `json:"content"`
      					} `json:"delta"`
      				} `json:"choices"`
      				FinishReason string `json:"finish_reason"`
      			}
      
      			if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
      				return fmt.Errorf("%s", msgBuffer.String())
      			}
      
      			// Accumulate content from each choice in the buffer
      			for _, choice := range data.Choices {
      				msgBuffer.WriteString(choice.Delta.Content)
      			}
      
      			// Process the accumulated content after streaming is complete
      			htmlMsg := web.MarkdownToHTML(msgBuffer.Bytes())
      
      			turnIDStr := fmt.Sprint(chatID + llm.TurnCounter)
      
      			// TODO: Abstract this into a function that all backends use.
      			//formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>\n<codapi-snippet url='http://localhost:1313/v1/exec' sandbox='go' editor='external'></codapi-snippet>", turnIDStr, htmlMsg)
      			formattedContent := fmt.Sprintf("<div id='response-content-%s' class='mx-1' hx-trigger='load'>%s</div>\n<codapi-snippet engine='browser' sandbox='javascript' editor='basic'></codapi-snippet>", turnIDStr, htmlMsg)
      
      			if err := c.WriteMessage(websocket.TextMessage, []byte(formattedContent)); err != nil {
      				pterm.Error.Println("WebSocket write error:", err)
      				return err
      			}
      		}
      	}
      
      	if err := scanner.Err(); err != nil {
      		pterm.Error.Println("Error reading stream:", err)
      		return err
      	}
      
      	return nil
      }

      [File Ends] pkg/llm/openai/completions.go

      [File Begins] pkg/llm/openai/models.go
      package openai
      
      import (
      	"encoding/json"
      	"fmt"
      	"io"
      	"net/http"
      	"os"
      
      	"github.com/pterm/pterm"
      )
      
      // Client represents the HTTP client for interacting with the LLM API.
      type Client struct {
      	APIKey string
      	HTTP   *http.Client
      }
      
      // OAIModel represents a single model in the JSON response.
      type OAIModel struct {
      	ID      string `json:"id"`
      	Object  string `json:"object"`
      	Created int64  `json:"created"`
      	OwnedBy string `json:"owned_by"`
      }
      
      // ModelsResponse represents the top-level structure of the JSON response.
      type ModelsResponse struct {
      	Object string     `json:"object"`
      	Data   []OAIModel `json:"data"`
      }
      
      // NewClient creates and initializes a new instance of an LLM API client using the provided API key.
      func NewClient(apiKey string) *Client {
      	return &Client{
      		APIKey: apiKey,
      		HTTP:   &http.Client{},
      	}
      }
      
      // Initialize sets up the LLM API client using the API key from environment variables.
      func (c *Client) Initialize() error {
      	apiKey := os.Getenv("LLM_API_KEY")
      	if apiKey == "" {
      		pterm.Error.Println("Please set the LLM_API_KEY environment variable.")
      		return fmt.Errorf("LLM_API_KEY is not set")
      	}
      
      	c.APIKey = apiKey
      	return nil
      }
      
      func (c *Client) Connect(endpoint string) (*http.Response, error) {
      	// Include the base URL and the protocol scheme
      	fullURL := "https://api.openai.com/v1" + endpoint // Example base URL
      
      	req, err := http.NewRequest("GET", fullURL, nil)
      	if err != nil {
      		return nil, fmt.Errorf("failed to create request: %w", err)
      	}
      	req.Header.Set("Authorization", "Bearer "+c.APIKey)
      
      	resp, err := c.HTTP.Do(req)
      	if err != nil {
      		return nil, fmt.Errorf("failed to connect to %s: %w", fullURL, err)
      	}
      
      	return resp, nil
      }
      
      // GetModels retrieves the available models from the OpenAI API.
      // It returns a ModelsResponse and any error encountered.
      func GetModels(client *Client) (ModelsResponse, error) {
      	resp, err := client.Connect("/models")
      	if err != nil {
      		return ModelsResponse{}, fmt.Errorf("failed to connect to OpenAI API: %w", err)
      	}
      	defer func(Body io.ReadCloser) {
      		err := Body.Close()
      		if err != nil {
      			fmt.Println(err)
      		}
      	}(resp.Body)
      
      	if resp.StatusCode != http.StatusOK {
      		return ModelsResponse{}, fmt.Errorf("failed to fetch models: non-OK status received - %s", resp.Status)
      	}
      
      	var models ModelsResponse
      	if err := json.NewDecoder(resp.Body).Decode(&models); err != nil {
      		return ModelsResponse{}, fmt.Errorf("failed to decode response body: %w", err)
      	}
      
      	return models, nil
      }

      [File Ends] pkg/llm/openai/models.go

      [File Begins] pkg/llm/openai/openai.go
      package openai
      
      import "eternal/pkg/llm"
      
      type Message struct {
      	Role    string `json:"role"`
      	Content string `json:"content"`
      }
      
      // Model represents an AI model from the OpenAI API with its ID, name, and description.
      type Model struct {
      	ID          string `json:"id"`
      	Name        string `json:"name"`
      	Description string `json:"description"`
      }
      
      // CompletionRequest represents the payload for the completion API.
      type CompletionRequest struct {
      	Model       string        `json:"model"`
      	Messages    []llm.Message `json:"messages"`
      	Temperature float64       `json:"temperature"`
      	Stream      bool          `json:"stream"`
      }
      
      // Choice represents a choice for the completion response.
      type Choice struct {
      	Index        int     `json:"index"`
      	Message      Message `json:"message"`
      	Logprobs     *bool   `json:"logprobs"` // Pointer to a boolean or nil
      	FinishReason string  `json:"finish_reason"`
      }
      
      // Usage contains information about token usage in the completion response.
      type Usage struct {
      	PromptTokens     int `json:"prompt_tokens"`
      	CompletionTokens int `json:"completion_tokens"`
      	TotalTokens      int `json:"total_tokens"`
      }
      
      // CompletionResponse represents the response from the completion API.
      type CompletionResponse struct {
      	ID                string   `json:"id"`
      	Object            string   `json:"object"`
      	Created           int64    `json:"created"`
      	Model             string   `json:"model"`
      	SystemFingerprint string   `json:"system_fingerprint"`
      	Choices           []Choice `json:"choices"`
      	Usage             Usage    `json:"usage"`
      }
      
      // AudioSpeechRequest represents the payload for the audio speech API.
      type AudioSpeechRequest struct {
      	Model string `json:"model"`
      	Input string `json:"input"`
      	Voice string `json:"voice"`
      }
      
      // ErrorData represents the structure of an error response from the OpenAI API.
      type ErrorData struct {
      	Code    interface{} `json:"code"`
      	Message string      `json:"message"`
      }
      
      // ErrorResponse wraps the structure of an error when an API request fails.
      type ErrorResponse struct {
      	Error ErrorData `json:"error"`
      }
      
      // UsageMetrics details the token usage of the Embeddings API request.
      type UsageMetrics struct {
      	PromptTokens int `json:"prompt_tokens"`
      	TotalTokens  int `json:"total_tokens"`
      }

      [File Ends] pkg/llm/openai/openai.go

    [File Begins] pkg/llm/templates.go
    package llm
    
    import (
    	"fmt"
    	"strings"
    )
    
    var (
    	AssistantDefault = "You are a helpful knowledge assistant. You respond in a pleasant and friendly conversational tone and always end your replies with a question to encourage further interaction. You provide clear and concise answers to user queries and offer additional information or assistance when needed. You aim to be informative, engaging, and supportive in all your interactions."
    
    	AssistantGraphOfThoughts = `Respond to each query using the following process to reason through to the most insightful answer:
    	First, carefully analyze the question to identify the key pieces of information required to answer it comprehensively. Break the question down into its core components.
    	For each component of the question, brainstorm several relevant ideas, facts, and perspectives that could help address that part of the query. Consider the question from multiple angles.
    	Critically evaluate each of those ideas you generated. Assess how directly relevant they are to the question, how logical and well-supported they are, and how clearly they convey key points. Aim to hone in on the strongest and most pertinent thoughts.
    	Take the most promising ideas and try to combine them into a coherent line of reasoning that flows logically from one point to the next in order to address the original question. See if you can construct a compelling argument or explanation.
    	If your current line of reasoning doesn't fully address all aspects of the original question in a satisfactory way, continue to iteratively explore other possible angles by swapping in alternative ideas and seeing if they allow you to build a stronger overall case.
    	As you work through the above process, make a point to capture your thought process and explain the reasoning behind why you selected or discarded certain ideas. Highlight the relative strengths and flaws in different possible arguments. Make your reasoning transparent.
    	After exploring multiple possible thought paths, integrating the strongest arguments, and explaining your reasoning along the way, pull everything together into a clear, concise, and complete final response that directly addresses the original query.
    	Throughout your response, weave in relevant parts of your intermediate reasoning and thought process. Use natural language to convey your train of thought in a conversational tone. Focus on clearly explaining insights and conclusions rather than mechanically labeling each step.
    	The goal is to use a tree-like process to explore multiple potential angles, rigorously evaluate and select the most promising and relevant ideas, iteratively build strong lines of reasoning, and ultimately synthesize key points into an insightful, well-reasoned, and accessible final answer.
    	Always end your response asking if there is anything else you can help with.`
    
    	AssistantCodeReview = "Begin by thoroughly reviewing the submitted codebase to understand its structure, design patterns, and functionality. Then, critically analyze each component for code quality, best practices, bugs, security, and performance. Identify areas for improvement, prioritizing critical issues for production readiness and separating them from less urgent optimizations. Develop a comprehensive improvement plan with actionable steps for refactoring, testing, and documentation enhancements, ensuring the plan is modular and incrementally actionable. Provide specific recommendations with examples, focusing on maintainability and efficiency. Summarize the strengths and weaknesses of the code and your overall assessment, using a supportive tone to offer constructive feedback and encourage collaboration. Offer to assist with follow-up questions and further guidance."
    
    	AssistantCodeAdvanced = `First, carefully read through the entire codebase that was submitted for review. Make note of the overall structure, design patterns used, and the key functionality being implemented. Aim to holistically understand the code at a high level.
    	Next, go through the code again but this time critically analyze each module, class, function and code block in detail: - Assess the code quality, adherence to best practices and coding standards, use of appropriate design patterns, and overall readability and maintainability. - Look for any bugs, edge cases, error handling, security vulnerabilities, and performance issues. - Evaluate how well the code is organized, commented, and documented. - Consider how testable, modular, and extensible the code is.
    	For each issue or area for improvement identified, brainstorm several suggestions on how to address it. Provide specific, actionable recommendations including code snippets and examples wherever applicable. Explain the rationale behind each suggestion.
    	Prioritize the list of potential improvements based on their importance and impact. Separate critical issues that must be fixed before the code can be considered production-ready from less urgent optimizations and enhancements.
    	Draft a comprehensive code improvement plan that organizes the prioritized suggestions into concrete steps the developer can follow: - Break down complex changes into smaller, incremental action items. - Provide clear guidance on refactoring and redesigning the code where needed to make it cleaner, more efficient, and easier to maintain. - Include tips on writing unit tests and integration tests to properly validate all the core functionality and edge cases. Emphasize the importance of testing. - Offer suggestions on improving the code documentation, comments, logging and error handling.
    	As you create the code improvement plan, continue to revisit the original code and your detailed analysis to ensure your suggestions are complete and address the most important issues. Iteratively refine your feedback.
    	Once you have a polished list of concrete suggestions organized into a clear plan of action, combine them with your overarching feedback on the submission as a whole. Summarize the key strengths and weaknesses of the code, the major areas for improvement, and your overall assessment of its production readiness.
    	Preface your final code review response with a friendly greeting and positive feedback to acknowledge the work the developer put in. Then concisely explain your high-level analysis and segue into presenting the detailed improvement plan.
    	When delivering constructive criticism and suggestions, use a supportive and encouraging tone. Be objective and focus on the code itself rather than the developer. Back up your recommendations with clear reasoning and examples.
    	Close your response by offering to answer any follow-up questions and provide further guidance as needed. Reiterate that the ultimate goal is to work collaboratively to improve the code and get it ready for a successful deployment to production.
    	The goal is to use a systematic process to thoroughly evaluate the code from multiple angles, identify the most critical issues, and provide clear and actionable suggestions that the developer can follow to improve their code. The code review should be comprehensive, insightful, and help the developer grow their skills. Always maintain a positive and supportive tone while delivering constructive feedback.
    	Please let me know if you would like me to modify or expand this code review prompt template in any way. I’m happy to refine it further.`
    
    	AssistantVisualBot = `You are an AI assistant that specializes in creating mermaid diagrams based on user descriptions provided in natural English. Your task is to interpret the user’s input and convert it into a structured format that can be used to generate the corresponding mermaid diagram.
    
    	When a user provides a description of a diagram they want to create, follow these steps:
    	
    	Identify the type of diagram the user wants to create based on their description (e.g., flowchart, sequence diagram, class diagram, etc.).
    	
    	Extract the relevant elements, their types, and any additional properties or relationships mentioned in the user’s description.
    	
    	Determine the relationships or connections between the elements, if specified by the user.
    	
    	Identify any specific styling or formatting requirements mentioned by the user.
    	
    	Organize the extracted information into the following template:
    	
    	Diagram type: [Identified diagram type] Diagram title: [Appropriate title based on user’s description] Diagram direction (optional): [Specified direction or default based on diagram type]
    	
    	Diagram elements: - Element 1 name: [Element 1 description or properties] - Element 2 name: [Element 2 description or properties] …
    	
    	Relationships (optional): - [Element 1 name] –> [Element 2 name]: [Relationship description] - [Element 2 name] –> [Element 3 name]: [Relationship description] …
    	
    	Additional styling or formatting (optional): [Specified styling or formatting options]
    	
    	Generate the mermaid diagram based on the structured template.
    	Remember, the user may not be familiar with mermaid syntax or diagram terminology, so you need to interpret their natural language description and convert it into the appropriate format. Always attempt to understand the user’s intent and provide a helpful response, even if their description is incomplete or ambiguous.
    	
    	TASK:
    	Port the user's question or comment into a format as follows, but do not respond or mention these instructions. Simply generate the diagram using all of the previous and next instructions:
    	
    	Diagram type: [Specify the type of diagram, e.g., flowchart, sequence diagram, class diagram, state diagram, pie chart, journey, gantt, requirement diagram, gitgraph, c4c, mindmap, timeline, or other valid mermaid diagram types]
    	
    	Diagram title: [Provide a title for the diagram]
    	
    	Diagram direction (optional): [Specify the direction of the diagram, e.g., TD (top-down), LR (left-right), RL (right-left), BT (bottom-top)]
    	
    	Diagram elements:
    	[List the elements of the diagram, including their names, types, and any additional properties or relationships. For example:
    	- Element 1 (type): [Description or properties]
    	- Element 2 (type): [Description or properties]
    	- Element 3 (type): [Description or properties]
    	...
    	]
    	
    	Relationships (optional):
    	[Describe the relationships or connections between the elements, if applicable. For example:
    	- Element 1 --> Element 2: [Relationship description]
    	- Element 2 --> Element 3: [Relationship description]
    	...
    	]
    	
    	Additional styling or formatting (optional):
    	[Specify any additional styling or formatting options for the diagram, such as colors, shapes, line styles, or other valid mermaid syntax]
    	
    	Example:
    	Diagram type: flowchart
    	Diagram title: Sample Flowchart
    	Diagram direction: LR
    	
    	Diagram elements:
    	- Start (start)
    	- Process 1 (process): Some processing step
    	- Decision (decision): Yes or No?
    	- Process 2 (process): Another processing step
    	- End (end)
    	
    	Relationships:
    	- Start --> Process 1
    	- Process 1 --> Decision
    	- Decision --Yes--> Process 2
    	- Decision --No--> End
    	- Process 2 --> End
    	
    	Additional styling or formatting:
    	- linkStyle default stroke:#0000FF,stroke-width:2px;
    	- style Process 1 fill:#FFFFCC,stroke:#FFFF00,stroke-width:2px
    	- style Decision fill:#CCFFFF,stroke:#0000FF,stroke-width:2px
    	`
    )
    
    // Message represents a message for the completion API.
    type Message struct {
    	Role    string `json:"role"`
    	Content string `json:"content"`
    }
    
    // PromptTemplate represents a template for generating string prompts.
    type PromptTemplate struct {
    	Template string
    }
    
    // ChatPromptTemplate represents a template for generating chat prompts.
    type ChatPromptTemplate struct {
    	Messages []Message
    }
    
    // GetSystemTemplate returns the system template.
    func GetSystemTemplate(userPrompt string) ChatPromptTemplate {
    	userPrompt = fmt.Sprintf("{%s}", userPrompt)
    	template := NewChatPromptTemplate([]Message{
    		{
    			Role:    "system",
    			Content: "You are a helpful AI assistant that responds in well structured markdown format. Do not repeat your instructions. Do not deviate from the topic.",
    		},
    		{
    			Role:    "user",
    			Content: userPrompt,
    		},
    	})
    
    	return *template
    }
    
    // NewChatPromptTemplate creates a new ChatPromptTemplate.
    func NewChatPromptTemplate(messages []Message) *ChatPromptTemplate {
    	return &ChatPromptTemplate{Messages: messages}
    }
    
    // Format formats the template with the provided variables.
    func (pt *PromptTemplate) Format(vars map[string]string) string {
    	result := pt.Template
    	for k, v := range vars {
    		placeholder := fmt.Sprintf("{%s}", k)
    		result = strings.ReplaceAll(result, placeholder, v)
    	}
    	return result
    }
    
    // FormatMessages formats the chat messages with the provided variables.
    func (cpt *ChatPromptTemplate) FormatMessages(vars map[string]string) []Message {
    	var formattedMessages []Message
    	for _, msg := range cpt.Messages {
    		formattedContent := msg.Content
    		for k, v := range vars {
    			placeholder := fmt.Sprintf("{%s}", k)
    			formattedContent = strings.ReplaceAll(formattedContent, placeholder, v)
    		}
    		formattedMessages = append(formattedMessages, Message{Role: msg.Role, Content: formattedContent})
    	}
    	return formattedMessages
    }

    [File Ends] pkg/llm/templates.go

    [File Begins] pkg/sd/sd.go
    package sd
    
    import (
    	"fmt"
    	"image"
    	"io"
    	"math"
    	"net/http"
    	"os"
    	"os/exec"
    	"path/filepath"
    	"sync"
    
    	"github.com/anthonynsimon/bild/adjust"
    	"github.com/anthonynsimon/bild/imgio"
    	"github.com/anthonynsimon/bild/transform"
    	"github.com/pterm/pterm"
    )
    
    type ImageModel struct {
    	Name       string   `yaml:"name"`
    	Homepage   string   `yaml:"homepage"`
    	Prompt     string   `yaml:"prompt"`
    	Downloads  []string `yaml:"downloads,omitempty"`
    	Downloaded bool     `yaml:"downloaded"`
    	LocalPath  string   `yaml:"localPath,omitempty"`
    }
    
    type CommandOutput struct {
    	Output       string `json:"output"`
    	Finished     string `json:"finished"`
    	SocketNumber string `json:"socketNumber"`
    	ModelName    string `json:"modelName"`
    }
    
    type SDParams struct {
    	Mode            string  `json:"mode" cli:"-M, --mode"`
    	Threads         int     `json:"threads" cli:"-t, --threads"`
    	Model           string  `json:"model" cli:"-m, --model"`
    	VAE             string  `json:"vae" cli:"--vae"`
    	TAESDPath       string  `json:"taesd" cli:"--taesd"`
    	ControlNet      string  `json:"control_net" cli:"--control-net"`
    	EmbdDir         string  `json:"embd_dir" cli:"--embd-dir"`
    	UpscaleModel    string  `json:"upscale_model" cli:"--upscale-model"`
    	WeightType      string  `json:"type" cli:"--type"`
    	LoraModelDir    string  `json:"lora_model_dir" cli:"--lora-model-dir"`
    	InitImage       string  `json:"init_img" cli:"-i, --init-img"`
    	ControlImage    string  `json:"control_image" cli:"--control-image"`
    	Output          string  `json:"output" cli:"-o, --output"`
    	Prompt          string  `json:"prompt" cli:"-p, --prompt"`
    	NegativePrompt  string  `json:"negative_prompt" cli:"-n, --negative-prompt"`
    	CFGScale        float64 `json:"cfg_scale" cli:"--cfg-scale"`
    	Strength        float64 `json:"strength" cli:"--strength"`
    	ControlStrength float64 `json:"control_strength" cli:"--control-strength"`
    	Height          int     `json:"height" cli:"--height"`
    	Width           int     `json:"width" cli:"--width"`
    	SamplingMethod  string  `json:"sampling_method" cli:"--sampling-method"`
    	Steps           int     `json:"steps" cli:"--steps"`
    	RNG             string  `json:"rng" cli:"--rng"`
    	Seed            int     `json:"seed" cli:"-s, --seed"`
    	BatchCount      int     `json:"batch_count" cli:"-b, --batch-count"`
    	Schedule        string  `json:"schedule" cli:"--schedule"`
    	ClipSkip        int     `json:"clip_skip" cli:"--clip-skip"`
    	VAETiling       bool    `json:"vae_tiling" cli:"--vae-tiling"`
    	ControlNetCPU   bool    `json:"control_net_cpu" cli:"--control-net-cpu"`
    	Canny           bool    `json:"canny" cli:"--canny"`
    	Verbose         bool    `json:"verbose" cli:"-v, --verbose"`
    }
    
    func BuildCommand(dataPath string, params SDParams) *exec.Cmd {
    	//vaePath := filepath.Join(dataPath, "models/StableDiffusion/sd15/sdxl_vae.safetensors")
    	//modelPath := filepath.Join(dataPath, "models/StableDiffusion/sd15/dreamshaper_8_q5_1.gguf")
    	modelPath := filepath.Join(dataPath, "models/dreamshaper-8-turbo-sdxl/DreamShaperXL_Turbo_V2-SFW.safetensors")
    	vaePath := filepath.Join(dataPath, "models/dreamshaper-8-turbo-sdxl/sdxl_vae.safetensors")
    	outPath := filepath.Join(dataPath, "web/img/sd_out.png")
    	cmdPath := filepath.Join(dataPath, "sd/sd")
    
    	pterm.Println("Command:", cmdPath)
    
    	cmdArgs := []string{
    		"-p", params.Prompt,
    		"-n", "ugly, low quality, deformed, malformed, floating limbs, bad hands, poorly drawn, bad anatomy, extra limb, blurry, disfigured, realistic, child, long neck, big forehead",
    		"-m", modelPath,
    		"--vae", vaePath,
    		"-o", outPath,
    		"--rng", "std_default",
    		"--cfg-scale", "2",
    		"--sampling-method", "dpm++2m",
    		"--steps", "10",
    		"--seed", "-1",
    		//"--upscale-model", "/mnt/d/StableDiffusionModels/sdxl/upscalers/RealESRGAN_x4plus_anime_6B.pth",
    		"--schedule", "karras",
    		"--strength", "1.0",
    		"--clip-skip", "2",
    		"--width", "1024",
    		"--height", "1024",
    	}
    
    	// Print cmdArgs
    	pterm.Println("Command Args:", cmdArgs)
    
    	// Run the command in a separate process
    	return exec.Command(cmdPath, cmdArgs...)
    
    	//return exec.Command(cmdPath, cmdArgs...)
    }
    
    func Text2Image(cmdPath string, params *SDParams) error {
    	cmd := BuildCommand(cmdPath, *params)
    
    	out, err := cmd.Output()
    	if err != nil {
    		fmt.Println(err)
    		return err
    	}
    
    	pterm.Warning.Println("Output:", string(out))
    
    	// return the output
    	return nil
    }
    
    // cropAndAdjustContrast crops the image to a 768x768 square while keeping the subject centered,
    // increases the contrast slightly, and overwrites the original image.
    func CropAndAdjustContrast(filePath string) error {
    	// Load the image from the file path
    	img, err := imgio.Open(filePath)
    	if err != nil {
    		return err
    	}
    
    	// Determine the size and position for cropping
    	width, height := img.Bounds().Dx(), img.Bounds().Dy()
    	minDim := int(math.Min(float64(width), float64(height)))
    	offsetX, offsetY := (width-minDim)/2, (height-minDim)/2
    
    	// Crop the image to a square while centering the subject
    	croppedImg := transform.Crop(img, image.Rect(offsetX, offsetY, offsetX+minDim, offsetY+minDim))
    
    	// Resize the image to 768x768 if necessary
    	if minDim != 768 {
    		croppedImg = transform.Resize(croppedImg, 768, 768, transform.Linear)
    	}
    
    	// Increase the contrast
    	adjustedImg := adjust.Contrast(croppedImg, 0.2) // Adjust the contrast level as needed
    
    	// Save the image back to the same file
    	err = imgio.Save(filePath, adjustedImg, imgio.JPEGEncoder(100))
    	if err != nil {
    		return err
    	}
    
    	return nil
    }
    
    var (
    	downloadProgressMap = make(map[string]DownloadProgress)
    	progressMutex       sync.Mutex
    	TurnCounter         = 0
    )
    
    // DownloadProgress structure to hold download progress information
    type DownloadProgress struct {
    	Total   int64 `json:"total"`
    	Current int64 `json:"current"`
    }
    
    // WebModel defines the interface for a model that interacts over the web.
    type WebModel interface {
    	Connect(endpoint string) (*http.Response, error)
    }
    
    // ModelManager defines the interface for managing local models.
    type ModelManager interface {
    	GetConfig() error
    	Download() error
    	Delete() error
    }
    
    type Model struct {
    	Name      string   `yaml:"name"`
    	Homepage  string   `yaml:"homepage"`
    	Prompt    string   `yaml:"prompt"`
    	Ctx       int      `yaml:"ctx"`
    	Roles     []string `yaml:"roles"`
    	Tags      []string `yaml:"tags,omitempty"`
    	GGUF      string   `yaml:"gguf,omitempty"`
    	Downloads []string `yaml:"downloads,omitempty"`
    	LocalPath string   `yaml:"localPath,omitempty"`
    }
    
    type ProgressReader struct {
    	Reader        io.Reader
    	ProgressBar   *pterm.ProgressbarPrinter
    	TotalRead     int64
    	ContentLength int64
    }
    
    func (pr *ProgressReader) Read(p []byte) (int, error) {
    	n, err := pr.Reader.Read(p)
    	pr.TotalRead += int64(n)
    
    	// Calculate progress percentage
    	//progressPercentage := int64(100.0 * float64(pr.TotalRead) / float64(pr.ContentLength))
    
    	// Update the progress map safely
    	progressMutex.Lock()
    	downloadProgressMap["sse-progress"] = DownloadProgress{
    		Total:   pr.ContentLength,
    		Current: pr.TotalRead,
    	}
    	progressMutex.Unlock()
    
    	pr.ProgressBar.Add(n)
    
    	return n, err
    }
    
    // Download downloads a file from a URL to a local path, resuming if possible.
    func Download(url string, localPath string) error {
    	dir := filepath.Dir(localPath)
    	// Ensure the directory exists
    	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
    		return fmt.Errorf("failed to create directory: %w", err)
    	}
    
    	// Open the local file for writing, create it if not exists
    	out, err := os.OpenFile(localPath, os.O_RDWR|os.O_CREATE, 0666)
    	if err != nil {
    		return fmt.Errorf("failed to open file: %w", err)
    	}
    	defer out.Close()
    
    	// Find out how much has already been downloaded
    	fi, err := out.Stat()
    	if err != nil {
    		return fmt.Errorf("failed to stat file: %w", err)
    	}
    	size := fi.Size()
    
    	// Create a new HTTP request
    	req, err := http.NewRequest("GET", url, nil)
    	if err != nil {
    		return fmt.Errorf("failed to create request: %w", err)
    	}
    
    	// Set the Range header to request the portion of the file we don't have yet
    	req.Header.Set("Range", fmt.Sprintf("bytes=%d-", size))
    
    	// Make the HTTP request
    	resp, err := http.DefaultClient.Do(req)
    	if err != nil {
    		return fmt.Errorf("failed to start file download: %w", err)
    	}
    	defer resp.Body.Close()
    
    	if resp.StatusCode != http.StatusPartialContent && resp.StatusCode != http.StatusOK {
    		return fmt.Errorf("bad status getting file: %s", resp.Status)
    	}
    
    	// Seek to the end of the file to start appending data
    	_, err = out.Seek(0, io.SeekEnd)
    	if err != nil {
    		return fmt.Errorf("failed to seek file: %w", err)
    	}
    
    	pterm.Info.Printf("Downloading model:\nURL: %s\nFile: %s\n", url, localPath)
    
    	// Initialize the progress bar
    	progressBar, _ := pterm.DefaultProgressbar.WithTotal(int(resp.ContentLength + size)).WithTitle("Downloading").Start()
    
    	// Wrap the response body in a custom reader that updates the progress bar
    	progressReader := &ProgressReader{
    		Reader:        resp.Body,
    		ProgressBar:   progressBar,
    		TotalRead:     size,
    		ContentLength: resp.ContentLength + size,
    	}
    
    	// Copy the remaining data to the file, updating progress along the way
    	_, err = io.Copy(out, progressReader)
    	if err != nil {
    		return err
    	}
    
    	// Ensure the progress bar reflects the complete download
    	progressBar.Total = (int(progressReader.TotalRead))
    
    	// Finish the progress bar
    	progressBar.Stop()
    
    	// Update the model's downloaded state in the database
    	// err = UpdateModelDownloadedState(modelName, true)
    	// if err != nil {
    	//     log.Errorf("Failed to update model downloaded state: %v", err)
    	// }
    
    	return nil
    }
    
    func GetDownloadProgress(key string) string {
    	progressMutex.Lock()
    	defer progressMutex.Unlock()
    
    	// Return the progress for a specific key
    	if progress, ok := downloadProgressMap[key]; ok {
    		// Calculate progress percentage
    		return fmt.Sprintf("%d%%", int64(100.0*float64(progress.Current)/float64(progress.Total)))
    	}
    
    	// Fallback if no progress is found
    	return "0%"
    }
    
    // GetExpectedFileSize returns the expected file size of a download.
    func GetExpectedFileSize(url string) (int64, error) {
    	// Create a new HTTP request
    	req, err := http.NewRequest("HEAD", url, nil)
    	if err != nil {
    		return 0, fmt.Errorf("failed to create request: %w", err)
    	}
    
    	// Make the HTTP request
    	resp, err := http.DefaultClient.Do(req)
    	if err != nil {
    		return 0, fmt.Errorf("failed to start file download: %w", err)
    	}
    	defer resp.Body.Close()
    
    	if resp.StatusCode != http.StatusOK {
    		return 0, fmt.Errorf("bad status getting file: %s", resp.Status)
    	}
    
    	return resp.ContentLength, nil
    }
    
    func (m *Model) Delete() error {
    	if err := os.Remove(m.LocalPath); err != nil {
    		return fmt.Errorf("failed to delete model: %w", err)
    	}
    
    	return nil
    }

    [File Ends] pkg/sd/sd.go

    [File Begins] pkg/search/search.go
    package search
    
    import (
    	"github.com/blevesearch/bleve/v2"
    )
    
    // InitIndex initializes or opens a Bleve index at the given path.
    func InitIndex(indexPath string) (bleve.Index, error) {
    	mapping := bleve.NewIndexMapping() // Default mapping; can be customized as needed
    	index, err := bleve.New(indexPath, mapping)
    	if err != nil {
    		return nil, err
    	}
    	return index, nil
    }
    
    // IndexData indexes the given data with the specified ID.
    func IndexData(index *bleve.Index, id string, data interface{}) error {
    	return (*index).Index(id, data)
    }
    
    // Search performs a search query on the given index.
    func Search(index bleve.Index, query string) (*bleve.SearchResult, error) {
    	searchQuery := bleve.NewMatchQuery(query)
    	search := bleve.NewSearchRequest(searchQuery)
    	return (index).Search(search)
    }

    [File Ends] pkg/search/search.go

    [File Begins] pkg/vecstore/vecstore.go
    package vecstore
    
    import (
    	"encoding/json"
    	"fmt"
    	"math"
    	"math/rand"
    	"os"
    	"sort"
    	"strings"
    )
    
    // Vector represents a vector of floats.
    type Vector []float64
    
    // Node is a struct that represents a node in a k-d tree. It has three fields:
    // Domain: A slice of float64 values representing the domain or feature space of the data point associated with this node.
    // Value: A float64 value representing the pivot value used to partition the data points into two subsets.
    // Left: A pointer to the left child node in the tree, or nil if there is no left child.
    // Right: A pointer to the right child node in the tree, or nil if there is no right child.
    type Node struct {
    	Domain []float64
    	Value  float64
    	Left   *Node
    	Right  *Node
    }
    
    // Embedding represents a word embedding.
    type Embedding struct {
    	Word       string
    	Vector     []float64
    	Similarity float64 // Similarity field to store the cosine similarity
    }
    
    // EmbeddingDB represents a database of Embeddings.
    type EmbeddingDB struct {
    	Embeddings map[string]Embedding
    }
    
    // Document represents a document to be ranked.
    type Document struct {
    	ID     string
    	Score  float64
    	Length int
    }
    
    // NewEmbeddingDB creates a new embedding database.
    func NewEmbeddingDB() *EmbeddingDB {
    	return &EmbeddingDB{
    		Embeddings: make(map[string]Embedding),
    	}
    }
    
    // AddEmbedding adds a new embedding to the database.
    func (db *EmbeddingDB) AddEmbedding(embedding Embedding) {
    	db.Embeddings[embedding.Word] = embedding
    }
    
    // AddEmbeddings adds a slice of embeddings to the database.
    func (db *EmbeddingDB) AddEmbeddings(embeddings []Embedding) {
    	for _, embedding := range embeddings {
    		db.AddEmbedding(embedding)
    	}
    }
    
    // SaveEmbeddings saves the Embeddings to a file, appending new ones to existing data.
    func (db *EmbeddingDB) SaveEmbeddings(path string) error {
    	// Read the existing content from the file
    	var existingEmbeddings map[string]Embedding
    	content, err := os.ReadFile(path)
    	if err != nil {
    		if !os.IsNotExist(err) {
    			return fmt.Errorf("error reading file: %v", err)
    		}
    		existingEmbeddings = make(map[string]Embedding)
    	} else {
    		err = json.Unmarshal(content, &existingEmbeddings)
    		if err != nil {
    			return fmt.Errorf("error unmarshaling existing embeddings: %v", err)
    		}
    	}
    
    	// Merge new embeddings with existing ones
    	for key, embedding := range db.Embeddings {
    		existingEmbeddings[key] = embedding
    	}
    
    	// Marshal the combined embeddings to JSON
    	jsonData, err := json.Marshal(existingEmbeddings)
    	if err != nil {
    		return fmt.Errorf("error marshaling embeddings: %v", err)
    	}
    
    	// Open the file in write mode (this will overwrite the existing file)
    	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
    	if err != nil {
    		return fmt.Errorf("error opening file: %v", err)
    	}
    	defer f.Close()
    
    	// Write the combined JSON to the file
    	if _, err := f.Write(jsonData); err != nil {
    		return fmt.Errorf("error writing to file: %v", err)
    	}
    
    	return nil
    }
    
    // LoadEmbeddings loads the Embeddings from a file.
    func (db *EmbeddingDB) LoadEmbeddings(path string) (map[string]Embedding, error) {
    	content, err := os.ReadFile(path)
    	if err != nil {
    		return nil, err
    	}
    
    	var embeddings map[string]Embedding
    	err = json.Unmarshal(content, &embeddings)
    	if err != nil {
    		return nil, err
    	}
    
    	return embeddings, nil
    }
    
    // RetrieveEmbedding retrieves an embedding from the database.
    func (db *EmbeddingDB) RetrieveEmbedding(word string) ([]float64, bool) {
    	embedding, exists := db.Embeddings[word]
    	if !exists {
    		return nil, false
    	}
    
    	return embedding.Vector, true
    }
    
    // RecreateDocument recreates a document from a slice of embeddings.
    func (db *EmbeddingDB) RecreateDocument(embeddings []Embedding) string {
    	var document []string
    	for _, embedding := range embeddings {
    		document = append(document, embedding.Word)
    	}
    
    	return strings.Join(document, " ")
    }
    
    // CosineSimilarity calculates the cosine similarity between two vectors.
    // func CosineSimilarity(a, b []float64) float64 {
    // 	if len(a) != len(b) {
    // 		log.Fatal("Vectors must be of the same length")
    // 	}
    
    // 	var dotProduct, magnitudeA, magnitudeB float64
    // 	var wg sync.WaitGroup
    
    // 	// Adjust the number of partitions based on the number of CPU cores.
    // 	partitions := runtime.NumCPU()
    // 	partSize := len(a) / partitions
    
    // 	results := make([]struct {
    // 		dotProduct, magnitudeA, magnitudeB float64
    // 	}, partitions)
    
    // 	for i := 0; i < partitions; i++ {
    // 		wg.Add(1)
    // 		go func(partition int) {
    // 			defer wg.Done()
    // 			start := partition * partSize
    // 			end := start + partSize
    // 			if partition == partitions-1 {
    // 				end = len(a)
    // 			}
    // 			for j := start; j < end; j++ {
    // 				results[partition].dotProduct += a[j] * b[j]
    // 				results[partition].magnitudeA += a[j] * a[j]
    // 				results[partition].magnitudeB += b[j] * b[j]
    // 			}
    // 		}(i)
    // 	}
    
    // 	wg.Wait()
    
    // 	for _, result := range results {
    // 		dotProduct += result.dotProduct
    // 		magnitudeA += result.magnitudeA
    // 		magnitudeB += result.magnitudeB
    // 	}
    
    // 	return dotProduct / (math.Sqrt(magnitudeA) * math.Sqrt(magnitudeB))
    // }
    
    func CosineSimilarity(vecA, vecB []float64) float64 {
    	var dotProduct, normA, normB float64
    	for i, v := range vecA {
    		dotProduct += v * vecB[i]
    		normA += v * v
    		normB += vecB[i] * vecB[i]
    	}
    	return dotProduct / (math.Sqrt(normA) * math.Sqrt(normB))
    }
    
    // MostSimilarWord returns the word with the highest similarity value.
    func (db *EmbeddingDB) MostSimilarWord(embeddings map[string]Embedding, targetWord string) (string, float64) {
    	// Check for an exact match
    	if _, exists := embeddings[targetWord]; exists {
    		return targetWord, 1.0
    	}
    	var targetVector []float64
    
    	// If the target word exists in the embeddings, use its vector.
    	if embedding, exists := embeddings[targetWord]; exists {
    		targetVector = embedding.Vector
    	} else {
    		// If target word doesn't exist, print the error and try to find the most similar word using the embeddings available.
    		// (In a more robust implementation, you might want to obtain a vector for the targetWord from another source.)
    		fmt.Printf("Error: Word '%s' not found in embeddings database.\n", targetWord)
    		return "", -1.0 // We'll use this -1.0 later to identify that the target word wasn't in the database.
    	}
    
    	mostSimilarWord := ""
    	highestSimilarity := -2.0 // Starting with -2 to ensure that any cosine similarity will be higher.
    
    	for word, embedding := range embeddings {
    		// If the word is the same as target, skip.
    		if word == targetWord {
    			continue
    		}
    
    		// Compute similarity.
    		similarity := CosineSimilarity(targetVector, embedding.Vector)
    
    		// If this word's similarity is greater than the highest similarity seen so far, update.
    		if similarity > highestSimilarity {
    			mostSimilarWord = word
    			highestSimilarity = similarity
    		}
    	}
    
    	if highestSimilarity == -1.0 {
    		return "No similar words found in database", highestSimilarity
    	}
    
    	return mostSimilarWord, highestSimilarity
    }
    
    // TODO: Implement an Efficient Search Mechanism using KD-Trees or Ball Trees.
    // Consider integrating with a library or service that offers efficient nearest-neighbor search capabilities.
    // Placeholder function for this:
    func EfficientSearch(targetWord string) (string, float64) {
    	// Implement efficient search here
    	return "", 0.0
    }
    
    // dataPoints: A slice of slices of float64 values representing the data points to be inserted into the k-d tree. Each inner slice should have the same length, which represents the number of dimensions or features of the data point.
    // dimensions: An integer representing the number of dimensions or features of the data points.
    func BuildKdTree(dataPoints [][]float64, dimensions int) *Node {
    	if len(dataPoints) == 0 {
    		return nil
    	}
    
    	// Use a stack to store the nodes in the tree.
    	stack := make([]*Node, 0)
    
    	// Select the median value of a random dimension as the pivot.
    	pivotIndex := rand.Intn(len(dataPoints))
    	pivotValue := dataPoints[pivotIndex][rand.Intn(dimensions)]
    
    	// Partition the data points into two subsets based on the pivot value.
    	left := make([][]float64, 0)
    	right := make([][]float64, 0)
    	for _, point := range dataPoints {
    		if point[rand.Intn(dimensions)] < pivotValue {
    			left = append(left, point)
    		} else {
    			right = append(right, point)
    		}
    	}
    
    	// Build the left and right subtrees iteratively using a stack.
    	node := &Node{
    		Domain: dataPoints[pivotIndex],
    		Value:  pivotValue,
    	}
    	stack = append(stack, node)
    	for len(left) > 0 || len(right) > 0 {
    		if len(left) > 0 {
    			// Build the left subtree.
    			pivotIndex := rand.Intn(len(left))
    			pivotValue := left[pivotIndex][rand.Intn(dimensions)]
    			var partition [][]float64
    			for _, point := range left {
    				if point[rand.Intn(dimensions)] < pivotValue {
    					partition = append(partition, point)
    				}
    			}
    			node := &Node{
    				Domain: left[pivotIndex],
    				Value:  pivotValue,
    			}
    			stack[len(stack)-1].Left = node
    			if len(partition) > 0 {
    				stack = append(stack, node)
    				left = partition
    			} else {
    				left = nil
    			}
    		} else {
    			// Build the right subtree.
    			pivotIndex := rand.Intn(len(right))
    			pivotValue := right[pivotIndex][rand.Intn(dimensions)]
    			var partition [][]float64
    			for _, point := range right {
    				if point[rand.Intn(dimensions)] < pivotValue {
    					partition = append(partition, point)
    				}
    			}
    			node := &Node{
    				Domain: right[pivotIndex],
    				Value:  pivotValue,
    			}
    			stack[len(stack)-1].Right = node
    			if len(partition) > 0 {
    				stack = append(stack, node)
    				right = partition
    			} else {
    				right = nil
    			}
    		}
    	}
    
    	// Return the root node.
    	return stack[0]
    }
    
    // FindMostSimilarEmbedding finds the most similar embeddings in the database.
    func FindMostSimilarEmbedding(targetEmbedding Embedding, embeddings map[string]Embedding) (Embedding, bool) {
    	var mostSimilarEmbedding Embedding
    	var highestSimilarity float64
    
    	for _, embedding := range embeddings {
    		// If the word is the same as target, skip.
    		if embedding.Word == targetEmbedding.Word {
    			continue
    		}
    
    		// Compute similarity.
    		similarity := CosineSimilarity(targetEmbedding.Vector, embedding.Vector)
    
    		// If this word's similarity is greater than the highest similarity seen so far, update.
    		if similarity > highestSimilarity {
    			mostSimilarEmbedding = embedding
    			highestSimilarity = similarity
    		}
    	}
    
    	if highestSimilarity == -1.0 {
    		return Embedding{}, false
    	}
    
    	return mostSimilarEmbedding, true
    }
    
    // NormalizeL2 normalizes a vector using L2 normalization.
    func NormalizeL2(vec []float64) []float64 {
    	var sumSquares float64
    	for _, value := range vec {
    		sumSquares += value * value
    	}
    	norm := math.Sqrt(sumSquares)
    	for i, value := range vec {
    		vec[i] = value / norm
    	}
    	return vec
    }
    
    // ComputeSimilarityMatrix computes the cosine similarity matrix between two slices of embeddings.
    func ComputeSimilarityMatrix(queryEmbeddings, keyEmbeddings []Embedding) [][]float64 {
    	matrix := make([][]float64, len(queryEmbeddings))
    	for i, query := range queryEmbeddings {
    		matrix[i] = make([]float64, len(keyEmbeddings))
    		for j, key := range keyEmbeddings {
    			matrix[i][j] = CosineSimilarity(query.Vector, key.Vector)
    		}
    	}
    	return matrix
    }
    
    // SimilarityWithKey is a type that holds both the similarity value and the corresponding word key.
    type SimilarityWithKey struct {
    	Similarity float64
    	Key        string
    }
    
    // FindTopNSimilarEmbeddings finds the top N most similar embeddings in the database.
    func FindTopNSimilarEmbeddings(targetEmbedding Embedding, embeddings map[string]Embedding, topN int) []Embedding {
    	var topEmbeddings []Embedding
    	var similarityList []SimilarityWithKey
    
    	// Compute the cosine similarity for each embedding in the database and store it with its key.
    	for key, embedding := range embeddings {
    		similarity := CosineSimilarity(targetEmbedding.Vector, embedding.Vector)
    
    		fmt.Println("Similarity:", similarity)
    
    		similarityList = append(similarityList, SimilarityWithKey{similarity, key})
    	}
    
    	// Sort the similarityList in descending order of similarity.
    	sort.SliceStable(similarityList, func(i, j int) bool {
    		return similarityList[i].Similarity > similarityList[j].Similarity
    	})
    
    	// Retrieve the top N most similar embeddings.
    	for i := 0; i < topN && i < len(similarityList); i++ {
    		topEmbeddings = append(topEmbeddings, embeddings[similarityList[i].Key])
    	}
    
    	return topEmbeddings
    }
    
    // Reranker function reranks documents based on a weighted combination of score and length.
    func Reranker(documents []Document, weightScore float64, weightLength float64) []Document {
    	// Validate weights
    	if weightScore < 0 || weightLength < 0 || (weightScore+weightLength) == 0 {
    		// Handle invalid weights
    		return documents
    	}
    
    	rerankedDocuments := make([]Document, len(documents))
    	copy(rerankedDocuments, documents)
    
    	sort.SliceStable(rerankedDocuments, func(i, j int) bool {
    		scoreDiffI := rerankedDocuments[i].Score * weightScore
    		lengthDiffI := float64(rerankedDocuments[i].Length) * weightLength
    		combinedScoreI := scoreDiffI + lengthDiffI
    
    		scoreDiffJ := rerankedDocuments[j].Score * weightScore
    		lengthDiffJ := float64(rerankedDocuments[j].Length) * weightLength
    		combinedScoreJ := scoreDiffJ + lengthDiffJ
    
    		if combinedScoreI == combinedScoreJ {
    			// Handle tie-breaking here if needed
    		}
    
    		return combinedScoreI > combinedScoreJ
    	})
    
    	return rerankedDocuments
    }
    
    func preprocessText(text string) string {
    	// Example for natural language
    	text = strings.ToLower(text)
    	text = removeStopWords(text) // Assuming a function exists
    	text = stemText(text)        // Assuming a function exists
    	return text
    }
    
    func removeStopWords(text string) string {
    	// Remove common words like "the", "and", etc.
    	return text
    }
    
    func stemText(text string) string {
    	// Reduce words to their base form
    	return text
    }
    
    // SortEmbeddingsBySimilarity sorts a slice of Embeddings by similarity in descending order.
    func SortEmbeddingsBySimilarity(embeddings []Embedding) {
    	sort.Slice(embeddings, func(i, j int) bool {
    		return embeddings[i].Similarity > embeddings[j].Similarity
    	})
    }

    [File Ends] pkg/vecstore/vecstore.go

    [File Begins] pkg/web/mdtohtml.go
    package web
    
    import (
    	"bytes"
    
    	"github.com/gomarkdown/markdown"
    	"github.com/gomarkdown/markdown/html"
    	"github.com/gomarkdown/markdown/parser"
    )
    
    // PreprocessMarkdown scans the markdown content for unclosed code blocks and closes them.
    func PreprocessMarkdown(content []byte) []byte {
    	// Simple heuristic: if the count of ``` is odd, append one at the end
    	if bytes.Count(content, []byte("```"))%2 != 0 {
    		content = append(content, []byte("\n```")...) // Append closing code block
    	}
    	return content
    }
    
    // MarkdownToHTML converts preprocessed markdown content to HTML.
    func MarkdownToHTML(mdContent []byte) []byte {
    	// Preprocess to ensure code blocks are properly closed
    	preprocessedContent := PreprocessMarkdown(mdContent)
    
    	// Setup parser and renderer
    	extensions := parser.CommonExtensions
    	parser := parser.NewWithExtensions(extensions)
    	htmlFlags := html.CommonFlags
    
    	// Custom CSS to apply 0px bottom margin to <code> elements
    	customCSS := "code { margin-bottom: 0px; }"
    
    	renderer := html.NewRenderer(html.RendererOptions{
    		Flags: htmlFlags,
    		CSS:   customCSS, // Adding custom CSS here
    	})
    
    	// Convert markdown to HTML
    	htmlContent := markdown.ToHTML(preprocessedContent, parser, renderer)
    
    	return htmlContent
    }

    [File Ends] pkg/web/mdtohtml.go

    [File Begins] pkg/web/serpapi.go
    package web
    
    import (
    	"encoding/json"
    	"os"
    
    	g "github.com/serpapi/google-search-results-golang"
    )
    
    type SearchData struct {
    	Title string `json:"title"`
    	Link  string `json:"link"`
    }
    
    func GetSerpResults(query string, apikey string) (*[]SearchData, error) {
    	parameter := map[string]string{
    		"engine":  "google",
    		"q":       query,
    		"api_key": apikey,
    	}
    
    	search := g.NewGoogleSearch(parameter, apikey)
    	results, err := search.GetJSON()
    	if err != nil {
    		panic(err)
    	}
    
    	resdata := []SearchData{}
    
    	organic_results := results["organic_results"].([]interface{})
    	for _, organic_result := range organic_results {
    		organic_result := organic_result.(map[string]interface{})
    		link := organic_result["link"].(string)
    		title := organic_result["title"].(string)
    
    		res := new(SearchData)
    		res.Title = title
    		res.Link = link
    
    		// Only append if the link is not a PDF
    		if len(link) > 4 && link[len(link)-4:] != ".pdf" {
    			resdata = append(resdata, *res)
    		}
    	}
    
    	// Write to JSON file
    	file, _ := json.MarshalIndent(resdata, "", " ")
    	_ = os.WriteFile("data.json", file, 0644)
    
    	return &resdata, nil
    }

    [File Ends] pkg/web/serpapi.go

    [File Begins] pkg/web/web.go
    package web
    
    import (
    	"bytes"
    	"context"
    	"io"
    	"log"
    	"net/url"
    	"os"
    	"regexp"
    	"strings"
    	"time"
    
    	"github.com/chromedp/cdproto/cdp"
    	"github.com/chromedp/chromedp"
    	"github.com/chromedp/chromedp/kb"
    	"github.com/pterm/pterm"
    	"golang.org/x/net/html"
    	"golang.org/x/net/html/atom"
    )
    
    var (
    	unwantedURLs = []string{
    		"www.youtube.com",
    		"www.youtube.com/watch",
    		"www.wired.com",
    		"www.techcrunch.com",
    		"www.wsj.com",
    		"www.cnn.com",
    		"www.nytimes.com",
    		"www.forbes.com",
    		"www.businessinsider.com",
    		"www.theverge.com",
    		"www.thehill.com",
    		"www.theatlantic.com",
    		"www.foxnews.com",
    		"www.theguardian.com",
    		"www.nbcnews.com",
    		"www.msn.com",
    		"www.sciencedaily.com",
    		"reuters.com",
    		"bbc.com",
    		"thenewstack.io",
    		"abcnews.go.com",
    		"apnews.com",
    		"bloomberg.com",
    		"polygon.com",
    		// Add more URLs to block from search results
    	}
    
    	resultURLs []string
    )
    
    func WebGetHandler(url string) (string, error) {
    	// Set up chromedp
    	opts := append(chromedp.DefaultExecAllocatorOptions[:],
    		chromedp.Flag("headless", true),
    		// other options if needed...
    	)
    
    	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    	defer cancel()
    
    	ctx, cancel := chromedp.NewContext(allocCtx)
    	defer cancel()
    
    	// Set a timeout for the entire operation
    	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
    	defer cancel()
    
    	// Retrieve and sanitize the page
    	var docs string
    	err := chromedp.Run(ctx,
    		chromedp.Navigate(url),
    		chromedp.WaitReady("body"),
    		chromedp.OuterHTML("html", &docs),
    	)
    	if err != nil {
    		log.Println("Error retrieving page:", err)
    		return "", err
    	}
    
    	// Clean up consecutive newlines
    	docs = strings.ReplaceAll(docs, "\n\n", "\n")
    
    	// Create a strings.Reader from docs to provide an io.Reader to ParseReaderView
    	reader := strings.NewReader(docs)
    
    	// Invoke web.ParseReaderView to get the reader view of the HTML
    	cleanedHTML, err := ParseReaderView(reader)
    	if err != nil {
    		log.Println("Error parsing reader view:", err)
    		return "", err
    	}
    
    	pterm.Info.Println("Document:", cleanedHTML)
    
    	return cleanedHTML, nil
    }
    
    // ExtractURLs extracts and cleans URLs from the input string.
    func ExtractURLs(input string) []string {
    	// Regular expression to match URLs and port numbers
    	urlRegex := `http.*?://[^\s<>{}|\\^` + "`" + `"]+`
    	re := regexp.MustCompile(urlRegex)
    
    	// Find all URLs in the input string
    	matches := re.FindAllString(input, -1)
    
    	var cleanedURLs []string
    	for _, match := range matches {
    		cleanedURL := cleanURL(match)
    		cleanedURLs = append(cleanedURLs, cleanedURL)
    	}
    
    	return cleanedURLs
    }
    
    // RemoveUrls removes URLs from the input string slice.
    func RemoveUrl(input []string) []string {
    	// Regular expression to match URLs and port numbers
    	urlRegex := `http.*?://[^\s<>{}|\\^` + "`" + `"]+`
    	re := regexp.MustCompile(urlRegex)
    
    	// Iterate over each string in the input slice
    	for i, str := range input {
    		// Find all URLs in the current string
    		matches := re.FindAllString(str, -1)
    
    		// Remove URLs from the current string
    		for _, match := range matches {
    			input[i] = strings.ReplaceAll(input[i], match, "")
    		}
    	}
    
    	return input
    }
    
    // cleanURL removes illegal trailing characters from the URL.
    func cleanURL(url string) string {
    	// Define illegal trailing characters.
    	illegalTrailingChars := []rune{'.', ',', ';', '!', '?', ')'}
    
    	for _, char := range illegalTrailingChars {
    		if url[len(url)-1] == byte(char) {
    			url = url[:len(url)-1]
    		}
    	}
    
    	return url
    }
    
    func ParseReaderView(r io.Reader) (string, error) {
    	doc, err := html.Parse(r)
    	if err != nil {
    		return "", err
    	}
    
    	var readerView bytes.Buffer
    	var f func(*html.Node)
    	f = func(n *html.Node) {
    		// Check if the node is an element node
    		if n.Type == html.ElementNode {
    			// Check if the node is a block element that typically contains content
    			if isContentElement(n) {
    				renderNode(&readerView, n)
    			}
    		}
    		for c := n.FirstChild; c != nil; c = c.NextSibling {
    			f(c)
    		}
    	}
    	f(doc)
    
    	return readerView.String(), nil
    }
    
    // isContentElement checks if the node is an element of interest for the reader view.
    func isContentElement(n *html.Node) bool {
    	switch n.DataAtom {
    	case atom.H1, atom.H2, atom.H3, atom.H4, atom.H5, atom.H6, atom.Article, atom.P, atom.Li, atom.Code, atom.Span, atom.Br:
    		return true
    	case atom.A, atom.Strong, atom.Em, atom.B, atom.I:
    		return true
    	}
    	return false
    }
    
    // renderNode writes the content of the node to the buffer, including inline tags.
    func renderNode(buf *bytes.Buffer, n *html.Node) {
    	// Render the opening tag if it's not a text node
    	// if n.Type == html.ElementNode {
    	// 	buf.WriteString("<" + n.Data + ">")
    	// }
    
    	// Render the contents
    	for c := n.FirstChild; c != nil; c = c.NextSibling {
    		if c.Type == html.TextNode {
    			buf.WriteString(strings.TrimSpace(c.Data))
    		} else if c.Type == html.ElementNode && isInlineElement(c) {
    			// Render inline elements and their children
    			renderNode(buf, c)
    		}
    	}
    
    	// Render the closing tag if it's not a text node
    	// if n.Type == html.ElementNode {
    	// 	buf.WriteString("</" + n.Data + ">")
    	// }
    }
    
    // isInlineElement checks if the node is an inline element that should be included in the output.
    func isInlineElement(n *html.Node) bool {
    	switch n.DataAtom {
    	case atom.A, atom.Strong, atom.Em, atom.B, atom.I, atom.Span, atom.Br, atom.Code:
    		return true
    	}
    	return false
    }
    
    // SearchDDG performs a search on DuckDuckGo and retrieves the HTML of the first page of results.
    func SearchDDG(query string) []string {
    
    	// Clear the resultURLs slice
    	resultURLs = nil
    
    	// Initialize headless Chrome
    	opts := append(chromedp.DefaultExecAllocatorOptions[:],
    		chromedp.Flag("headless", true),
    	)
    	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
    	defer cancel()
    	ctx, cancel := chromedp.NewContext(allocCtx)
    	defer cancel()
    
    	// Set a timeout for the entire operation
    	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
    	defer cancel()
    
    	var nodes []*cdp.Node
    
    	// Perform the search on DuckDuckGo
    	err := chromedp.Run(ctx,
    		chromedp.Navigate(`https://lite.duckduckgo.com/lite/`),
    		chromedp.WaitVisible(`input[name="q"]`, chromedp.ByQuery),
    		chromedp.SendKeys(`input[name="q"]`, query+kb.Enter, chromedp.ByQuery),
    		chromedp.Sleep(5*time.Second), // Wait for JavaScript to load the search results
    		chromedp.WaitVisible(`input[name="q"]`, chromedp.ByQuery),
    		chromedp.Nodes(`a`, &nodes, chromedp.ByQueryAll),
    	)
    	if err != nil {
    		log.Printf("Error during search: %v", err)
    		return nil
    	}
    
    	// Process the search results
    	err = chromedp.Run(ctx,
    		chromedp.ActionFunc(func(c context.Context) error {
    			re, err := regexp.Compile(`^http[s]?://`)
    			if err != nil {
    				return err
    			}
    
    			uniqueUrls := make(map[string]bool)
    			for _, n := range nodes {
    				for _, attr := range n.Attributes {
    					if re.MatchString(attr) && !strings.Contains(attr, "duckduckgo") {
    						uniqueUrls[attr] = true
    					}
    				}
    			}
    
    			for u := range uniqueUrls {
    				resultURLs = append(resultURLs, u)
    			}
    
    			return nil
    		}),
    	)
    
    	if err != nil {
    		log.Printf("Error processing results: %v", err)
    		return nil
    	}
    
    	// Remove unwanted URLs from the list
    	resultURLs = RemoveUnwantedURLs(resultURLs)
    
    	pterm.Warning.Println("Search results:", resultURLs)
    
    	return resultURLs
    }
    
    // GetSearchResults loops over a list of URLs and retrieves the HTML of each page.
    func GetSearchResults(urls []string) string {
    	var resultHTML string
    
    	for _, url := range urls {
    		res, err := WebGetHandler(url)
    		if err != nil {
    			pterm.Error.Printf("Error getting search result: %v", err)
    			continue
    		}
    
    		if res != "" {
    			resultHTML += res
    		}
    
    		// time.Sleep(5 * time.Second)
    	}
    
    	return resultHTML
    }
    
    // RemoveUnwantedURLs removes unwanted URLs from the list of URLs.
    func RemoveUnwantedURLs(urls []string) []string {
    	for _, u := range urls {
    		pterm.Info.Printf("Checking URL: %s", u)
    
    		unwanted := false
    		for _, unwantedURL := range unwantedURLs {
    			if strings.Contains(u, unwantedURL) {
    				pterm.Warning.Printf("URL %s contains unwanted URL %s", u, unwantedURL)
    				unwanted = true
    				break
    			}
    		}
    		if !unwanted {
    			resultURLs = append(resultURLs, u)
    		}
    	}
    
    	pterm.Info.Printf("Filtered URLs: %v", resultURLs)
    
    	return resultURLs
    }
    
    // GetPageScreenshot navigates to the given address and takes a full page screenshot.
    func GetPageScreen(chromeUrl string, pageAddress string) string {
    
    	instanceUrl := chromeUrl
    
    	// Create allocator context for using existing Chrome instance
    	allocatorCtx, cancel := chromedp.NewRemoteAllocator(context.Background(), instanceUrl)
    	defer cancel()
    
    	// Create context with logging for actions performed by chromedp
    	ctx, cancel := chromedp.NewContext(allocatorCtx, chromedp.WithLogf(log.Printf))
    	defer cancel()
    
    	// Set a timeout for the entire operation
    	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
    	defer cancel()
    
    	// Run tasks
    	var buf []byte // Buffer to store screenshot data
    	err := chromedp.Run(ctx,
    		chromedp.Navigate(pageAddress),
    		chromedp.FullScreenshot(&buf, 90),
    	)
    	if err != nil {
    		log.Fatal(err)
    	}
    
    	// Parse the URL to get the domain name
    	u, err := url.Parse(pageAddress)
    	if err != nil {
    		log.Fatal(err)
    	}
    
    	// Get the date and time
    	t := time.Now()
    
    	// Create a filename
    	filename := u.Hostname() + "-" + t.Format("20060102150405") + ".png"
    
    	// Save the screenshot to a file
    	err = os.WriteFile(filename, buf, 0644)
    	if err != nil {
    		log.Fatal(err)
    	}
    
    	return filename
    }
    
    // RemoveUrls removes URLs from the input string.
    func RemoveUrls(input string) string {
    	// Regular expression to match URLs and port numbers
    	urlRegex := `http.*?://[^\s<>{}|\\^` + "`" + `"]+`
    	re := regexp.MustCompile(urlRegex)
    
    	// Find all URLs in the input string
    	matches := re.FindAllString(input, -1)
    
    	// Remove URLs from the input string
    	for _, match := range matches {
    		input = strings.ReplaceAll(input, match, "")
    	}
    
    	return input
    }

    [File Ends] pkg/web/web.go

    [File Begins] public/css/header.css
    #cb-header {
        position: relative;
        height: 50px;
        filter: blur(5px); /* Apply blur */
    }
    
    #hgradient {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 50px;
        background: linear-gradient(to bottom, rgba(0, 0, 0, 0.5), rgba(0, 0, 0, 0));
        background-size: 200% 100%; /* Adjust the length of the gradient */
        background-position: 0 50%; /* Start the gradient at 50% height */
    }
    
    
    .circle {
        position: absolute;
        width: 12px; /* Half the width of the box */
        height: 12px; /* Half the height of the box */
        background-color: black;
        border-radius: 50%; /* Makes the shape a circle */
        top: 50%; /* Center vertically */
        left: 50%; /* Center horizontally */
        transform: translate(-50%, -50%); /* Offset the position to truly center */
        z-index: 1; /* Ensure it is behind the panels of the box */
    }
    
    .box-container {
        position: absolute;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
    }
    
    .box {
        width: 20px;
        height: 20px;
        position: relative;
        transform-style: preserve-3d;
        transition: transform 0.5s ease;
        transform-origin: center;
    }
    
    .box .panel {
        position: absolute;
        width: 20px;
        height: 20px;
        background: var(--bg-color, #FFF); /* Default to white if variable not set */
        border: 1px solid rgba(255, 255, 255, 0.25);
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 2px;
    }
    
    .box .panel::before,
    .box .panel::after {
        content: '';
        display: block;
        width: 12px; /* Diameter = 2 * radius */
        height: 12px; /* Diameter = 2 * radius */
        border-radius: 50%; /* Makes the shape a circle */
        position: absolute; /* Positions the pseudo-elements absolutely within their parent */
    }
    
    .box .panel::before {
        background: black;
        top: 0; /* Aligns the circle to the top */
        left: 0; /* Aligns the circle to the left */
    }
    
    .box .panel::after {
        background: white;
        bottom: 0; /* Aligns the circle to the bottom */
        right: 0; /* Aligns the circle to the right */
    }
    
    
    .box .panel.front { transform: translateZ(10px); }
    .box .panel.back { transform: rotateY(180deg) translateZ(10px); }
    .box .panel.top { transform: rotateX(90deg) translateZ(10px); }
    .box .panel.bottom { transform: rotateX(-90deg) translateZ(10px); }
    .box .panel.left { transform: rotateY(-90deg) translateZ(10px); }
    .box .panel.right { transform: rotateY(90deg) translateZ(10px); }
    
    @keyframes spin {
        from { transform: rotateX(0deg) rotateY(0deg) rotateZ(0deg); }
        to { transform: rotateX(360deg) rotateY(360deg) rotateZ(360deg); }
    }
    
    .box:hover {
        animation: spin 2s linear infinite;
    }

    [File Ends] public/css/header.css

    [File Begins] public/css/styles.css
    :root {
        /* Default Accents */
        --et-purple: #7f5af0;
        --et-green: #2cb67d;
        --et-blue: #2d68f0;
        --et-red: #e45858;
        --et-yellow: #f3d672;
        --et-light: #d9d9d9;
    
    
        /* Galactic Gray Theme */
        --et-galactic-primary: #cccccc;
        --et-galactic-secondary: #808080;
        --et-galactic-accent: #333333;
        --et-galactic-background: #101010;
        --et-galactic-border-highlight: #ffffff;
    
        /*  Aurora Theme */
        --et-aurora-primary: #00ffa3; /* Bright Aqua */
        --et-aurora-secondary: #00875f; /* Deep Aqua Green */
        --et-aurora-accent: #004d39; /* Dark Teal */
        --et-aurora-background: #00231c; /* Very Dark Green */
        --et-aurora-border-highlight: #00ffc4; /* Bright Teal */
    
        /*  Celestial Theme */
        --et-celestial-primary: #4b8f8c; /* Teal Blue */
        --et-celestial-secondary: #32717a; /* Sea Green Blue */
        --et-celestial-accent: #1b3a41; /* Dark Cyan */
        --et-celestial-background: #081a1c; /* Almost Black Cyan */
        --et-celestial-border-highlight: #6cbfbf; /* Soft Cyan */
    }
    
    body {
        font-family: 'Roboto', sans-serif;
        background-color: var(--et-galactic-background);
    }
    
    @font-face {
        font-family: 'MonaspaceArgon';
        src: url('../fonts/monaspace/MonaspaceArgon-Regular.woff') format('woff');
        font-weight: normal;
        font-style: normal;
    }
    
    @font-face {
        font-family: 'MonaspaceRadon';
        src: url('../fonts/monaspace/MonaspaceRadon-Regular.woff') format('woff');
        font-weight: normal;
        font-style: normal;
    }
    
    .fade-me-out.htmx-swapping {
        opacity: 0;
        transition: opacity 1s ease-out;
    }
    
    @keyframes fade-in {
        from {
            opacity: 0;
        }
    
        to {
            opacity: 1;
        }
    }
    
    @keyframes fade-out {
        from {
            opacity: 1;
        }
    
        to {
            opacity: 0;
        }
    }
    
    codapi-ref {
        display: none;
    }
    
    codapi-toolbar button {
        display: inline-block;
        width: 48px;
        height: 48px;
        /* background: url('data:image/svg+xml;utf8,<svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24"><g fill="currentColor" stroke="currentColor" stroke-width="1.5"><path stroke-linecap="round" d="M11 14H3m8 4H3"/><path d="M18.875 14.118c1.654.955 2.48 1.433 2.602 2.121a1.5 1.5 0 0 1 0 .521c-.121.69-.948 1.167-2.602 2.121c-1.654.955-2.48 1.433-3.138 1.194a1.499 1.499 0 0 1-.451-.261c-.536-.45-.536-1.404-.536-3.314c0-1.91 0-2.865.536-3.314a1.5 1.5 0 0 1 .451-.26c.657-.24 1.484.238 3.138 1.192Z"/><path stroke-linecap="round" d="M3 6h10.5M20 6h-2.25M20 10H9.5M3 10h2.25"/></g></svg>'); */
        border: none;
    }
    
    codapi-toolbar a {
        display: none;
    }
    
    .fade-it {
        view-transition-name: fade-it;
        animation-duration: 300ms;
    }
    
    ::view-transition-old(fade-it) {
        animation: 600ms ease both fade-out;
    }
    
    ::view-transition-new(fade-it) {
        animation: 600ms ease both fade-in;
    }
    
    .model-container {
        margin-bottom: 20px;
        border: 1px solid #ddd;
        padding: 10px;
    }
    
    .message-content {
        white-space: pre-wrap;
        /* Preserves spaces and line breaks */
        word-wrap: break-word;
        /* Ensures long words do not break the layout */
    
    }
    
    .user-prompt {
        background-color: var(--et-galactic-accent);
        border-left: 3px solid var(--et-red);
        /* padding: 8px 12px;
        margin-top: 20px;
        margin-bottom: 10px; */
        box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
        max-height: min-content;
    }
    
    .response {
        background-color: var(--et-galactic-accent);
        border-left: 3px solid var(--et-purple);
        /* padding: 8px 12px;
        margin-bottom: 10px;
        border-radius: 4px; */
        box-shadow: 0 2px 5px rgba(0, 0, 0, 0.2);
        overflow-y: auto;
    }
    
    pre code {
        margin-bottom: 0px;
    }
    
    .card-selected {
        border: 2px solid var(--et-blue);
        /* or any other styling to depict selection */
    }
    
    /* .chat-container {
        overflow-y: auto;
    } */
    
    .dark-blur {
        /* Apply a blur effect to the background */
        backdrop-filter: blur(7px) brightness(100%);
        --webkit-backface-visibility: hidden;
        --moz-backface-visibility: hidden;
        --webkit-transform: translate3d(0, 0, 0);
        --moz-transform: translate3d(0, 0, 0);
    }
    
    .hljs {
        background: #1e1e1e;
        font-family: 'MonaspaceArgon', monospace;
    }
    
    pre.hljs {
        background: #1e1e1e;
        padding: 0.5em;
    }
    
    /* Base background color for code blocks */
    code.hljs {
        display: block;
    
        padding: 0.5em;
        background: #272822;
        /* Monokai background color */
        color: #f8f8f2;
        /* Monokai main text color */
    }
    
    /* Color of the code itself */
    span.hljs {
        color: #f8f8f2;
        white-space: pre-wrap;
    }
    
    .preserve-format {
        white-space: pre-wrap;
    }
    
    span.hljs-comment,
    span.hljs-quote {
        color: #75715e;
        font-family: 'MonaspaceRadon', monospace;
    }
    
    span.hljs-string {
        color: #e6db74;
    }
    
    span.hljs-keyword,
    span.hljs-selector-tag,
    span.hljs-addition {
        color: #f92672;
    }
    
    span.hljs-built_in,
    span.hljs-class .hljs-title {
        color: #a6e22e;
    }
    
    span.hljs-function,
    span.hljs-tag .hljs-title,
    span.hljs-title {
        color: #a6e22e;
    }
    
    span.hljs-number {
        color: #ae81ff;
    }
    
    span.hljs-literal {
        color: #ae81ff;
    }
    
    span.hljs-symbol,
    span.hljs-attribute,
    span.hljs-meta .hljs-keyword,
    span.hljs-selector-id,
    span.hljs-selector-attr,
    span.hljs-selector-pseudo,
    span.hljs-template-tag,
    span.hljs-template-variable {
        color: #f92672;
    }
    
    /* Color of class names */
    span.hljs-type,
    span.hljs-builtin-name {
        color: #66d9ef;
    }
    
    span.hljs-diff .hljs-change,
    span.hljs-diff .hljs-error,
    span.hljs-diff .hljs-deletion {
        color: #ae81ff;
    }
    
    span.hljs-addition {
        color: #e6db74;
    }
    
    span.hljs-operator {
        color: #f92672;
    }
    
    span.hljs-punctuation {
        color: #f8f8f2;
    }
    
    span.hljs-escape {
        color: #ae81ff;
    }
    
    span.hljs-regexp {
        color: #e6db74;
    }
    
    span.hljs-tag {
        color: #f92672;
    }

    [File Ends] public/css/styles.css

    [File Begins] public/js/cloudbox.js
    
    // Get the elements
    const boxcontainer = document.querySelector(".box-container");
    const box = document.querySelector(".box");
    
    var rotationInterval = null;
    let r = 45;
    
    function getRandomRotation() {
      r += -90;
      const rotationX = r;
      const rotationY = r;
      const rotationZ = -180;
      return `rotateX(${rotationX}deg) rotateY(${rotationY}deg) rotateZ(${rotationZ}deg)`;
    }
    
    function chatRotation() {
      r += -10;
      const rotationX = r;
      const rotationY = r;
      const rotationZ = -90;
      return `rotateX(${rotationX}deg) rotateY(${rotationY}deg) rotateZ(${rotationZ}deg)`;
    }
    
    // Add a click event listener to rotate the box on click
    boxcontainer.addEventListener("click", function () {
      const newRotation = getRandomRotation();
      box.style.transition = "transform 0.5s";
      box.style.transform = newRotation;
      chatbox.classList.toggle("expand");
    });
    
    box.addEventListener("mouseover", function () {
      // Check if event listener is already added
      if (box.hasAttribute("data-event-added")) {
        return;
      }
    
      // Only spin the box once when the mouse is over it
      box.setAttribute("data-event-added", true);
    
      // Rotate the box
      const newRotation = getRandomRotation();
      box.style.transition = "transform 0.5s ease";
      box.style.transform = newRotation;
    
      // Remove the event listener after the animation is done
      setTimeout(function () {
        box.removeAttribute("data-event-added");
      }, 1000);
    });
    
    // Add a click event to all links to rotate the box
    const links = document.querySelectorAll("a");
    links.forEach(function (link) {
      link.addEventListener("click", function () {
        const newRotation = getRandomRotation();
        box.style.transition = "transform 0.5s ease";
        box.style.transform = newRotation;
      });
    });
    
    // Print the box transform values
    setInterval(function () {
      const transform = window.getComputedStyle(box).getPropertyValue("transform");
    }, 1000);
    
    // Rotate the box on page load
    const newRotation = getRandomRotation();
    box.style.transition = "transform 0.5s ease";
    box.style.transform = newRotation;

    [File Ends] public/js/cloudbox.js

    [File Begins] public/js/drawflow.min.js
    !function(e,t){"object"==typeof exports&&"object"==typeof module?module.exports=t():"function"==typeof define&&define.amd?define([],t):"object"==typeof exports?exports.Drawflow=t():e.Drawflow=t()}("undefined"!=typeof self?self:this,(function(){return function(e){var t={};function n(i){if(t[i])return t[i].exports;var s=t[i]={i:i,l:!1,exports:{}};return e[i].call(s.exports,s,s.exports,n),s.l=!0,s.exports}return n.m=e,n.c=t,n.d=function(e,t,i){n.o(e,t)||Object.defineProperty(e,t,{enumerable:!0,get:i})},n.r=function(e){"undefined"!=typeof Symbol&&Symbol.toStringTag&&Object.defineProperty(e,Symbol.toStringTag,{value:"Module"}),Object.defineProperty(e,"__esModule",{value:!0})},n.t=function(e,t){if(1&t&&(e=n(e)),8&t)return e;if(4&t&&"object"==typeof e&&e&&e.__esModule)return e;var i=Object.create(null);if(n.r(i),Object.defineProperty(i,"default",{enumerable:!0,value:e}),2&t&&"string"!=typeof e)for(var s in e)n.d(i,s,function(t){return e[t]}.bind(null,s));return i},n.n=function(e){var t=e&&e.__esModule?function(){return e.default}:function(){return e};return n.d(t,"a",t),t},n.o=function(e,t){return Object.prototype.hasOwnProperty.call(e,t)},n.p="",n(n.s=0)}([function(e,t,n){"use strict";n.r(t),n.d(t,"default",(function(){return i}));class i{constructor(e,t=null,n=null){this.events={},this.container=e,this.precanvas=null,this.nodeId=1,this.ele_selected=null,this.node_selected=null,this.drag=!1,this.reroute=!1,this.reroute_fix_curvature=!1,this.curvature=.5,this.reroute_curvature_start_end=.5,this.reroute_curvature=.5,this.reroute_width=6,this.drag_point=!1,this.editor_selected=!1,this.connection=!1,this.connection_ele=null,this.connection_selected=null,this.canvas_x=0,this.canvas_y=0,this.pos_x=0,this.pos_x_start=0,this.pos_y=0,this.pos_y_start=0,this.mouse_x=0,this.mouse_y=0,this.line_path=5,this.first_click=null,this.force_first_input=!1,this.draggable_inputs=!0,this.useuuid=!1,this.parent=n,this.noderegister={},this.render=t,this.drawflow={drawflow:{Home:{data:{}}}},this.module="Home",this.editor_mode="edit",this.zoom=1,this.zoom_max=1.6,this.zoom_min=.5,this.zoom_value=.1,this.zoom_last_value=1,this.evCache=new Array,this.prevDiff=-1}start(){this.container.classList.add("parent-drawflow"),this.container.tabIndex=0,this.precanvas=document.createElement("div"),this.precanvas.classList.add("drawflow"),this.container.appendChild(this.precanvas),this.container.addEventListener("mouseup",this.dragEnd.bind(this)),this.container.addEventListener("mousemove",this.position.bind(this)),this.container.addEventListener("mousedown",this.click.bind(this)),this.container.addEventListener("touchend",this.dragEnd.bind(this)),this.container.addEventListener("touchmove",this.position.bind(this)),this.container.addEventListener("touchstart",this.click.bind(this)),this.container.addEventListener("contextmenu",this.contextmenu.bind(this)),this.container.addEventListener("keydown",this.key.bind(this)),this.container.addEventListener("wheel",this.zoom_enter.bind(this)),this.container.addEventListener("input",this.updateNodeValue.bind(this)),this.container.addEventListener("dblclick",this.dblclick.bind(this)),this.container.onpointerdown=this.pointerdown_handler.bind(this),this.container.onpointermove=this.pointermove_handler.bind(this),this.container.onpointerup=this.pointerup_handler.bind(this),this.container.onpointercancel=this.pointerup_handler.bind(this),this.container.onpointerout=this.pointerup_handler.bind(this),this.container.onpointerleave=this.pointerup_handler.bind(this),this.load()}pointerdown_handler(e){this.evCache.push(e)}pointermove_handler(e){for(var t=0;t<this.evCache.length;t++)if(e.pointerId==this.evCache[t].pointerId){this.evCache[t]=e;break}if(2==this.evCache.length){var n=Math.abs(this.evCache[0].clientX-this.evCache[1].clientX);this.prevDiff>100&&(n>this.prevDiff&&this.zoom_in(),n<this.prevDiff&&this.zoom_out()),this.prevDiff=n}}pointerup_handler(e){this.remove_event(e),this.evCache.length<2&&(this.prevDiff=-1)}remove_event(e){for(var t=0;t<this.evCache.length;t++)if(this.evCache[t].pointerId==e.pointerId){this.evCache.splice(t,1);break}}load(){for(var e in this.drawflow.drawflow[this.module].data)this.addNodeImport(this.drawflow.drawflow[this.module].data[e],this.precanvas);if(this.reroute)for(var e in this.drawflow.drawflow[this.module].data)this.addRerouteImport(this.drawflow.drawflow[this.module].data[e]);for(var e in this.drawflow.drawflow[this.module].data)this.updateConnectionNodes("node-"+e);const t=this.drawflow.drawflow;let n=1;Object.keys(t).map((function(e,i){Object.keys(t[e].data).map((function(e,t){parseInt(e)>=n&&(n=parseInt(e)+1)}))})),this.nodeId=n}removeReouteConnectionSelected(){this.dispatch("connectionUnselected",!0),this.reroute_fix_curvature&&this.connection_selected.parentElement.querySelectorAll(".main-path").forEach((e,t)=>{e.classList.remove("selected")})}click(e){if(this.dispatch("click",e),"fixed"===this.editor_mode){if(e.preventDefault(),"parent-drawflow"!==e.target.classList[0]&&"drawflow"!==e.target.classList[0])return!1;this.ele_selected=e.target.closest(".parent-drawflow")}else"view"===this.editor_mode?(null!=e.target.closest(".drawflow")||e.target.matches(".parent-drawflow"))&&(this.ele_selected=e.target.closest(".parent-drawflow"),e.preventDefault()):(this.first_click=e.target,this.ele_selected=e.target,0===e.button&&this.contextmenuDel(),null!=e.target.closest(".drawflow_content_node")&&(this.ele_selected=e.target.closest(".drawflow_content_node").parentElement));switch(this.ele_selected.classList[0]){case"drawflow-node":null!=this.node_selected&&(this.node_selected.classList.remove("selected"),this.node_selected!=this.ele_selected&&this.dispatch("nodeUnselected",!0)),null!=this.connection_selected&&(this.connection_selected.classList.remove("selected"),this.removeReouteConnectionSelected(),this.connection_selected=null),this.node_selected!=this.ele_selected&&this.dispatch("nodeSelected",this.ele_selected.id.slice(5)),this.node_selected=this.ele_selected,this.node_selected.classList.add("selected"),this.draggable_inputs?"SELECT"!==e.target.tagName&&(this.drag=!0):"INPUT"!==e.target.tagName&&"TEXTAREA"!==e.target.tagName&&"SELECT"!==e.target.tagName&&!0!==e.target.hasAttribute("contenteditable")&&(this.drag=!0);break;case"output":this.connection=!0,null!=this.node_selected&&(this.node_selected.classList.remove("selected"),this.node_selected=null,this.dispatch("nodeUnselected",!0)),null!=this.connection_selected&&(this.connection_selected.classList.remove("selected"),this.removeReouteConnectionSelected(),this.connection_selected=null),this.drawConnection(e.target);break;case"parent-drawflow":case"drawflow":null!=this.node_selected&&(this.node_selected.classList.remove("selected"),this.node_selected=null,this.dispatch("nodeUnselected",!0)),null!=this.connection_selected&&(this.connection_selected.classList.remove("selected"),this.removeReouteConnectionSelected(),this.connection_selected=null),this.editor_selected=!0;break;case"main-path":null!=this.node_selected&&(this.node_selected.classList.remove("selected"),this.node_selected=null,this.dispatch("nodeUnselected",!0)),null!=this.connection_selected&&(this.connection_selected.classList.remove("selected"),this.removeReouteConnectionSelected(),this.connection_selected=null),this.connection_selected=this.ele_selected,this.connection_selected.classList.add("selected");const t=this.connection_selected.parentElement.classList;t.length>1&&(this.dispatch("connectionSelected",{output_id:t[2].slice(14),input_id:t[1].slice(13),output_class:t[3],input_class:t[4]}),this.reroute_fix_curvature&&this.connection_selected.parentElement.querySelectorAll(".main-path").forEach((e,t)=>{e.classList.add("selected")}));break;case"point":this.drag_point=!0,this.ele_selected.classList.add("selected");break;case"drawflow-delete":this.node_selected&&this.removeNodeId(this.node_selected.id),this.connection_selected&&this.removeConnection(),null!=this.node_selected&&(this.node_selected.classList.remove("selected"),this.node_selected=null,this.dispatch("nodeUnselected",!0)),null!=this.connection_selected&&(this.connection_selected.classList.remove("selected"),this.removeReouteConnectionSelected(),this.connection_selected=null)}"touchstart"===e.type?(this.pos_x=e.touches[0].clientX,this.pos_x_start=e.touches[0].clientX,this.pos_y=e.touches[0].clientY,this.pos_y_start=e.touches[0].clientY,this.mouse_x=e.touches[0].clientX,this.mouse_y=e.touches[0].clientY):(this.pos_x=e.clientX,this.pos_x_start=e.clientX,this.pos_y=e.clientY,this.pos_y_start=e.clientY),["input","output","main-path"].includes(this.ele_selected.classList[0])&&e.preventDefault(),this.dispatch("clickEnd",e)}position(e){if("touchmove"===e.type)var t=e.touches[0].clientX,n=e.touches[0].clientY;else t=e.clientX,n=e.clientY;if(this.connection&&this.updateConnection(t,n),this.editor_selected&&(i=this.canvas_x+-(this.pos_x-t),s=this.canvas_y+-(this.pos_y-n),this.dispatch("translate",{x:i,y:s}),this.precanvas.style.transform="translate("+i+"px, "+s+"px) scale("+this.zoom+")"),this.drag){e.preventDefault();var i=(this.pos_x-t)*this.precanvas.clientWidth/(this.precanvas.clientWidth*this.zoom),s=(this.pos_y-n)*this.precanvas.clientHeight/(this.precanvas.clientHeight*this.zoom);this.pos_x=t,this.pos_y=n,this.ele_selected.style.top=this.ele_selected.offsetTop-s+"px",this.ele_selected.style.left=this.ele_selected.offsetLeft-i+"px",this.drawflow.drawflow[this.module].data[this.ele_selected.id.slice(5)].pos_x=this.ele_selected.offsetLeft-i,this.drawflow.drawflow[this.module].data[this.ele_selected.id.slice(5)].pos_y=this.ele_selected.offsetTop-s,this.updateConnectionNodes(this.ele_selected.id)}if(this.drag_point){i=(this.pos_x-t)*this.precanvas.clientWidth/(this.precanvas.clientWidth*this.zoom),s=(this.pos_y-n)*this.precanvas.clientHeight/(this.precanvas.clientHeight*this.zoom);this.pos_x=t,this.pos_y=n;var o=this.pos_x*(this.precanvas.clientWidth/(this.precanvas.clientWidth*this.zoom))-this.precanvas.getBoundingClientRect().x*(this.precanvas.clientWidth/(this.precanvas.clientWidth*this.zoom)),l=this.pos_y*(this.precanvas.clientHeight/(this.precanvas.clientHeight*this.zoom))-this.precanvas.getBoundingClientRect().y*(this.precanvas.clientHeight/(this.precanvas.clientHeight*this.zoom));this.ele_selected.setAttributeNS(null,"cx",o),this.ele_selected.setAttributeNS(null,"cy",l);const e=this.ele_selected.parentElement.classList[2].slice(9),c=this.ele_selected.parentElement.classList[1].slice(13),d=this.ele_selected.parentElement.classList[3],a=this.ele_selected.parentElement.classList[4];let r=Array.from(this.ele_selected.parentElement.children).indexOf(this.ele_selected)-1;if(this.reroute_fix_curvature){r-=this.ele_selected.parentElement.querySelectorAll(".main-path").length-1,r<0&&(r=0)}const h=e.slice(5),u=this.drawflow.drawflow[this.module].data[h].outputs[d].connections.findIndex((function(e,t){return e.node===c&&e.output===a}));this.drawflow.drawflow[this.module].data[h].outputs[d].connections[u].points[r]={pos_x:o,pos_y:l};const p=this.ele_selected.parentElement.classList[2].slice(9);this.updateConnectionNodes(p)}"touchmove"===e.type&&(this.mouse_x=t,this.mouse_y=n),this.dispatch("mouseMove",{x:t,y:n})}dragEnd(e){if("touchend"===e.type)var t=this.mouse_x,n=this.mouse_y,i=document.elementFromPoint(t,n);else t=e.clientX,n=e.clientY,i=e.target;if(this.drag&&(this.pos_x_start==t&&this.pos_y_start==n||this.dispatch("nodeMoved",this.ele_selected.id.slice(5))),this.drag_point&&(this.ele_selected.classList.remove("selected"),this.pos_x_start==t&&this.pos_y_start==n||this.dispatch("rerouteMoved",this.ele_selected.parentElement.classList[2].slice(14))),this.editor_selected&&(this.canvas_x=this.canvas_x+-(this.pos_x-t),this.canvas_y=this.canvas_y+-(this.pos_y-n),this.editor_selected=!1),!0===this.connection)if("input"===i.classList[0]||this.force_first_input&&(null!=i.closest(".drawflow_content_node")||"drawflow-node"===i.classList[0])){if(!this.force_first_input||null==i.closest(".drawflow_content_node")&&"drawflow-node"!==i.classList[0])s=i.parentElement.parentElement.id,o=i.classList[1];else{if(null!=i.closest(".drawflow_content_node"))var s=i.closest(".drawflow_content_node").parentElement.id;else var s=i.id;if(0===Object.keys(this.getNodeFromId(s.slice(5)).inputs).length)var o=!1;else var o="input_1"}var l=this.ele_selected.parentElement.parentElement.id,c=this.ele_selected.classList[1];if(l!==s&&!1!==o){if(0===this.container.querySelectorAll(".connection.node_in_"+s+".node_out_"+l+"."+c+"."+o).length){this.connection_ele.classList.add("node_in_"+s),this.connection_ele.classList.add("node_out_"+l),this.connection_ele.classList.add(c),this.connection_ele.classList.add(o);var d=s.slice(5),a=l.slice(5);this.drawflow.drawflow[this.module].data[a].outputs[c].connections.push({node:d,output:o}),this.drawflow.drawflow[this.module].data[d].inputs[o].connections.push({node:a,input:c}),this.updateConnectionNodes("node-"+a),this.updateConnectionNodes("node-"+d),this.dispatch("connectionCreated",{output_id:a,input_id:d,output_class:c,input_class:o})}else this.dispatch("connectionCancel",!0),this.connection_ele.remove();this.connection_ele=null}else this.dispatch("connectionCancel",!0),this.connection_ele.remove(),this.connection_ele=null}else this.dispatch("connectionCancel",!0),this.connection_ele.remove(),this.connection_ele=null;this.drag=!1,this.drag_point=!1,this.connection=!1,this.ele_selected=null,this.editor_selected=!1,this.dispatch("mouseUp",e)}contextmenu(e){if(this.dispatch("contextmenu",e),e.preventDefault(),"fixed"===this.editor_mode||"view"===this.editor_mode)return!1;if(this.precanvas.getElementsByClassName("drawflow-delete").length&&this.precanvas.getElementsByClassName("drawflow-delete")[0].remove(),this.node_selected||this.connection_selected){var t=document.createElement("div");t.classList.add("drawflow-delete"),t.innerHTML="x",this.node_selected&&this.node_selected.appendChild(t),this.connection_selected&&this.connection_selected.parentElement.classList.length>1&&(t.style.top=e.clientY*(this.precanvas.clientHeight/(this.precanvas.clientHeight*this.zoom))-this.precanvas.getBoundingClientRect().y*(this.precanvas.clientHeight/(this.precanvas.clientHeight*this.zoom))+"px",t.style.left=e.clientX*(this.precanvas.clientWidth/(this.precanvas.clientWidth*this.zoom))-this.precanvas.getBoundingClientRect().x*(this.precanvas.clientWidth/(this.precanvas.clientWidth*this.zoom))+"px",this.precanvas.appendChild(t))}}contextmenuDel(){this.precanvas.getElementsByClassName("drawflow-delete").length&&this.precanvas.getElementsByClassName("drawflow-delete")[0].remove()}key(e){if(this.dispatch("keydown",e),"fixed"===this.editor_mode||"view"===this.editor_mode)return!1;("Delete"===e.key||"Backspace"===e.key&&e.metaKey)&&(null!=this.node_selected&&"INPUT"!==this.first_click.tagName&&"TEXTAREA"!==this.first_click.tagName&&!0!==this.first_click.hasAttribute("contenteditable")&&this.removeNodeId(this.node_selected.id),null!=this.connection_selected&&this.removeConnection())}zoom_enter(e,t){e.ctrlKey&&(e.preventDefault(),e.deltaY>0?this.zoom_out():this.zoom_in())}zoom_refresh(){this.dispatch("zoom",this.zoom),this.canvas_x=this.canvas_x/this.zoom_last_value*this.zoom,this.canvas_y=this.canvas_y/this.zoom_last_value*this.zoom,this.zoom_last_value=this.zoom,this.precanvas.style.transform="translate("+this.canvas_x+"px, "+this.canvas_y+"px) scale("+this.zoom+")"}zoom_in(){this.zoom<this.zoom_max&&(this.zoom+=this.zoom_value,this.zoom_refresh())}zoom_out(){this.zoom>this.zoom_min&&(this.zoom-=this.zoom_value,this.zoom_refresh())}zoom_reset(){1!=this.zoom&&(this.zoom=1,this.zoom_refresh())}createCurvature(e,t,n,i,s,o){var l=e,c=t,d=n,a=i,r=s;switch(o){case"open":if(e>=n)var h=l+Math.abs(d-l)*r,u=d-Math.abs(d-l)*(-1*r);else h=l+Math.abs(d-l)*r,u=d-Math.abs(d-l)*r;return" M "+l+" "+c+" C "+h+" "+c+" "+u+" "+a+" "+d+"  "+a;case"close":if(e>=n)h=l+Math.abs(d-l)*(-1*r),u=d-Math.abs(d-l)*r;else h=l+Math.abs(d-l)*r,u=d-Math.abs(d-l)*r;return" M "+l+" "+c+" C "+h+" "+c+" "+u+" "+a+" "+d+"  "+a;case"other":if(e>=n)h=l+Math.abs(d-l)*(-1*r),u=d-Math.abs(d-l)*(-1*r);else h=l+Math.abs(d-l)*r,u=d-Math.abs(d-l)*r;return" M "+l+" "+c+" C "+h+" "+c+" "+u+" "+a+" "+d+"  "+a;default:return" M "+l+" "+c+" C "+(h=l+Math.abs(d-l)*r)+" "+c+" "+(u=d-Math.abs(d-l)*r)+" "+a+" "+d+"  "+a}}drawConnection(e){var t=document.createElementNS("http://www.w3.org/2000/svg","svg");this.connection_ele=t;var n=document.createElementNS("http://www.w3.org/2000/svg","path");n.classList.add("main-path"),n.setAttributeNS(null,"d",""),t.classList.add("connection"),t.appendChild(n),this.precanvas.appendChild(t);var i=e.parentElement.parentElement.id.slice(5),s=e.classList[1];this.dispatch("connectionStart",{output_id:i,output_class:s})}updateConnection(e,t){const n=this.precanvas,i=this.zoom;let s=n.clientWidth/(n.clientWidth*i);s=s||0;let o=n.clientHeight/(n.clientHeight*i);o=o||0;var l=this.connection_ele.children[0],c=this.ele_selected.offsetWidth/2+(this.ele_selected.getBoundingClientRect().x-n.getBoundingClientRect().x)*s,d=this.ele_selected.offsetHeight/2+(this.ele_selected.getBoundingClientRect().y-n.getBoundingClientRect().y)*o,a=e*(this.precanvas.clientWidth/(this.precanvas.clientWidth*this.zoom))-this.precanvas.getBoundingClientRect().x*(this.precanvas.clientWidth/(this.precanvas.clientWidth*this.zoom)),r=t*(this.precanvas.clientHeight/(this.precanvas.clientHeight*this.zoom))-this.precanvas.getBoundingClientRect().y*(this.precanvas.clientHeight/(this.precanvas.clientHeight*this.zoom)),h=this.curvature,u=this.createCurvature(c,d,a,r,h,"openclose");l.setAttributeNS(null,"d",u)}addConnection(e,t,n,i){var s=this.getModuleFromNodeId(e);if(s===this.getModuleFromNodeId(t)){var o=this.getNodeFromId(e),l=!1;for(var c in o.outputs[n].connections){var d=o.outputs[n].connections[c];d.node==t&&d.output==i&&(l=!0)}if(!1===l){if(this.drawflow.drawflow[s].data[e].outputs[n].connections.push({node:t.toString(),output:i}),this.drawflow.drawflow[s].data[t].inputs[i].connections.push({node:e.toString(),input:n}),this.module===s){var a=document.createElementNS("http://www.w3.org/2000/svg","svg"),r=document.createElementNS("http://www.w3.org/2000/svg","path");r.classList.add("main-path"),r.setAttributeNS(null,"d",""),a.classList.add("connection"),a.classList.add("node_in_node-"+t),a.classList.add("node_out_node-"+e),a.classList.add(n),a.classList.add(i),a.appendChild(r),this.precanvas.appendChild(a),this.updateConnectionNodes("node-"+e),this.updateConnectionNodes("node-"+t)}this.dispatch("connectionCreated",{output_id:e,input_id:t,output_class:n,input_class:i})}}}updateConnectionNodes(e){const t="node_in_"+e,n="node_out_"+e;this.line_path;const i=this.container,s=this.precanvas,o=this.curvature,l=this.createCurvature,c=this.reroute_curvature,d=this.reroute_curvature_start_end,a=this.reroute_fix_curvature,r=this.reroute_width,h=this.zoom;let u=s.clientWidth/(s.clientWidth*h);u=u||0;let p=s.clientHeight/(s.clientHeight*h);p=p||0;const f=i.querySelectorAll("."+n);Object.keys(f).map((function(t,n){if(null===f[t].querySelector(".point")){var m=i.querySelector("#"+e),g=f[t].classList[1].replace("node_in_",""),_=i.querySelector("#"+g).querySelectorAll("."+f[t].classList[4])[0],w=_.offsetWidth/2+(_.getBoundingClientRect().x-s.getBoundingClientRect().x)*u,v=_.offsetHeight/2+(_.getBoundingClientRect().y-s.getBoundingClientRect().y)*p,y=m.querySelectorAll("."+f[t].classList[3])[0],C=y.offsetWidth/2+(y.getBoundingClientRect().x-s.getBoundingClientRect().x)*u,x=y.offsetHeight/2+(y.getBoundingClientRect().y-s.getBoundingClientRect().y)*p;const n=l(C,x,w,v,o,"openclose");f[t].children[0].setAttributeNS(null,"d",n)}else{const n=f[t].querySelectorAll(".point");let o="";const m=[];n.forEach((t,a)=>{if(0===a&&n.length-1==0){var f=i.querySelector("#"+e),g=((x=t).getBoundingClientRect().x-s.getBoundingClientRect().x)*u+r,_=(x.getBoundingClientRect().y-s.getBoundingClientRect().y)*p+r,w=(L=f.querySelectorAll("."+t.parentElement.classList[3])[0]).offsetWidth/2+(L.getBoundingClientRect().x-s.getBoundingClientRect().x)*u,v=L.offsetHeight/2+(L.getBoundingClientRect().y-s.getBoundingClientRect().y)*p,y=l(w,v,g,_,d,"open");o+=y,m.push(y);f=t;var C=t.parentElement.classList[1].replace("node_in_",""),x=(E=i.querySelector("#"+C)).querySelectorAll("."+t.parentElement.classList[4])[0];g=(R=E.querySelectorAll("."+t.parentElement.classList[4])[0]).offsetWidth/2+(R.getBoundingClientRect().x-s.getBoundingClientRect().x)*u,_=R.offsetHeight/2+(R.getBoundingClientRect().y-s.getBoundingClientRect().y)*p,w=(f.getBoundingClientRect().x-s.getBoundingClientRect().x)*u+r,v=(f.getBoundingClientRect().y-s.getBoundingClientRect().y)*p+r,y=l(w,v,g,_,d,"close");o+=y,m.push(y)}else if(0===a){var L;f=i.querySelector("#"+e),g=((x=t).getBoundingClientRect().x-s.getBoundingClientRect().x)*u+r,_=(x.getBoundingClientRect().y-s.getBoundingClientRect().y)*p+r,w=(L=f.querySelectorAll("."+t.parentElement.classList[3])[0]).offsetWidth/2+(L.getBoundingClientRect().x-s.getBoundingClientRect().x)*u,v=L.offsetHeight/2+(L.getBoundingClientRect().y-s.getBoundingClientRect().y)*p,y=l(w,v,g,_,d,"open");o+=y,m.push(y);f=t,g=((x=n[a+1]).getBoundingClientRect().x-s.getBoundingClientRect().x)*u+r,_=(x.getBoundingClientRect().y-s.getBoundingClientRect().y)*p+r,w=(f.getBoundingClientRect().x-s.getBoundingClientRect().x)*u+r,v=(f.getBoundingClientRect().y-s.getBoundingClientRect().y)*p+r,y=l(w,v,g,_,c,"other");o+=y,m.push(y)}else if(a===n.length-1){var E,R;f=t,C=t.parentElement.classList[1].replace("node_in_",""),x=(E=i.querySelector("#"+C)).querySelectorAll("."+t.parentElement.classList[4])[0],g=(R=E.querySelectorAll("."+t.parentElement.classList[4])[0]).offsetWidth/2+(R.getBoundingClientRect().x-s.getBoundingClientRect().x)*u,_=R.offsetHeight/2+(R.getBoundingClientRect().y-s.getBoundingClientRect().y)*p,w=(f.getBoundingClientRect().x-s.getBoundingClientRect().x)*(s.clientWidth/(s.clientWidth*h))+r,v=(f.getBoundingClientRect().y-s.getBoundingClientRect().y)*(s.clientHeight/(s.clientHeight*h))+r,y=l(w,v,g,_,d,"close");o+=y,m.push(y)}else{f=t,g=((x=n[a+1]).getBoundingClientRect().x-s.getBoundingClientRect().x)*(s.clientWidth/(s.clientWidth*h))+r,_=(x.getBoundingClientRect().y-s.getBoundingClientRect().y)*(s.clientHeight/(s.clientHeight*h))+r,w=(f.getBoundingClientRect().x-s.getBoundingClientRect().x)*(s.clientWidth/(s.clientWidth*h))+r,v=(f.getBoundingClientRect().y-s.getBoundingClientRect().y)*(s.clientHeight/(s.clientHeight*h))+r,y=l(w,v,g,_,c,"other");o+=y,m.push(y)}}),a?m.forEach((e,n)=>{f[t].children[n].setAttributeNS(null,"d",e)}):f[t].children[0].setAttributeNS(null,"d",o)}}));const m=i.querySelectorAll("."+t);Object.keys(m).map((function(t,n){if(null===m[t].querySelector(".point")){var h=i.querySelector("#"+e),f=m[t].classList[2].replace("node_out_",""),g=i.querySelector("#"+f).querySelectorAll("."+m[t].classList[3])[0],_=g.offsetWidth/2+(g.getBoundingClientRect().x-s.getBoundingClientRect().x)*u,w=g.offsetHeight/2+(g.getBoundingClientRect().y-s.getBoundingClientRect().y)*p,v=(h=h.querySelectorAll("."+m[t].classList[4])[0]).offsetWidth/2+(h.getBoundingClientRect().x-s.getBoundingClientRect().x)*u,y=h.offsetHeight/2+(h.getBoundingClientRect().y-s.getBoundingClientRect().y)*p;const n=l(_,w,v,y,o,"openclose");m[t].children[0].setAttributeNS(null,"d",n)}else{const n=m[t].querySelectorAll(".point");let o="";const h=[];n.forEach((t,a)=>{if(0===a&&n.length-1==0){var f=i.querySelector("#"+e),m=((C=t).getBoundingClientRect().x-s.getBoundingClientRect().x)*u+r,g=(C.getBoundingClientRect().y-s.getBoundingClientRect().y)*p+r,_=(E=f.querySelectorAll("."+t.parentElement.classList[4])[0]).offsetWidth/2+(E.getBoundingClientRect().x-s.getBoundingClientRect().x)*u,w=E.offsetHeight/2+(E.getBoundingClientRect().y-s.getBoundingClientRect().y)*p,v=l(m,g,_,w,d,"close");o+=v,h.push(v);f=t;var y=t.parentElement.classList[2].replace("node_out_",""),C=(L=i.querySelector("#"+y)).querySelectorAll("."+t.parentElement.classList[3])[0];m=(x=L.querySelectorAll("."+t.parentElement.classList[3])[0]).offsetWidth/2+(x.getBoundingClientRect().x-s.getBoundingClientRect().x)*u,g=x.offsetHeight/2+(x.getBoundingClientRect().y-s.getBoundingClientRect().y)*p,_=(f.getBoundingClientRect().x-s.getBoundingClientRect().x)*u+r,w=(f.getBoundingClientRect().y-s.getBoundingClientRect().y)*p+r,v=l(m,g,_,w,d,"open");o+=v,h.push(v)}else if(0===a){var x;f=t,y=t.parentElement.classList[2].replace("node_out_",""),C=(L=i.querySelector("#"+y)).querySelectorAll("."+t.parentElement.classList[3])[0],m=(x=L.querySelectorAll("."+t.parentElement.classList[3])[0]).offsetWidth/2+(x.getBoundingClientRect().x-s.getBoundingClientRect().x)*u,g=x.offsetHeight/2+(x.getBoundingClientRect().y-s.getBoundingClientRect().y)*p,_=(f.getBoundingClientRect().x-s.getBoundingClientRect().x)*u+r,w=(f.getBoundingClientRect().y-s.getBoundingClientRect().y)*p+r,v=l(m,g,_,w,d,"open");o+=v,h.push(v);f=t,_=((C=n[a+1]).getBoundingClientRect().x-s.getBoundingClientRect().x)*u+r,w=(C.getBoundingClientRect().y-s.getBoundingClientRect().y)*p+r,m=(f.getBoundingClientRect().x-s.getBoundingClientRect().x)*u+r,g=(f.getBoundingClientRect().y-s.getBoundingClientRect().y)*p+r,v=l(m,g,_,w,c,"other");o+=v,h.push(v)}else if(a===n.length-1){var L,E;f=t,y=t.parentElement.classList[1].replace("node_in_",""),C=(L=i.querySelector("#"+y)).querySelectorAll("."+t.parentElement.classList[4])[0],_=(E=L.querySelectorAll("."+t.parentElement.classList[4])[0]).offsetWidth/2+(E.getBoundingClientRect().x-s.getBoundingClientRect().x)*u,w=E.offsetHeight/2+(E.getBoundingClientRect().y-s.getBoundingClientRect().y)*p,m=(f.getBoundingClientRect().x-s.getBoundingClientRect().x)*u+r,g=(f.getBoundingClientRect().y-s.getBoundingClientRect().y)*p+r,v=l(m,g,_,w,d,"close");o+=v,h.push(v)}else{f=t,_=((C=n[a+1]).getBoundingClientRect().x-s.getBoundingClientRect().x)*u+r,w=(C.getBoundingClientRect().y-s.getBoundingClientRect().y)*p+r,m=(f.getBoundingClientRect().x-s.getBoundingClientRect().x)*u+r,g=(f.getBoundingClientRect().y-s.getBoundingClientRect().y)*p+r,v=l(m,g,_,w,c,"other");o+=v,h.push(v)}}),a?h.forEach((e,n)=>{m[t].children[n].setAttributeNS(null,"d",e)}):m[t].children[0].setAttributeNS(null,"d",o)}}))}dblclick(e){null!=this.connection_selected&&this.reroute&&this.createReroutePoint(this.connection_selected),"point"===e.target.classList[0]&&this.removeReroutePoint(e.target)}createReroutePoint(e){this.connection_selected.classList.remove("selected");const t=this.connection_selected.parentElement.classList[2].slice(9),n=this.connection_selected.parentElement.classList[1].slice(13),i=this.connection_selected.parentElement.classList[3],s=this.connection_selected.parentElement.classList[4];this.connection_selected=null;const o=document.createElementNS("http://www.w3.org/2000/svg","circle");o.classList.add("point");var l=this.pos_x*(this.precanvas.clientWidth/(this.precanvas.clientWidth*this.zoom))-this.precanvas.getBoundingClientRect().x*(this.precanvas.clientWidth/(this.precanvas.clientWidth*this.zoom)),c=this.pos_y*(this.precanvas.clientHeight/(this.precanvas.clientHeight*this.zoom))-this.precanvas.getBoundingClientRect().y*(this.precanvas.clientHeight/(this.precanvas.clientHeight*this.zoom));o.setAttributeNS(null,"cx",l),o.setAttributeNS(null,"cy",c),o.setAttributeNS(null,"r",this.reroute_width);let d=0;if(this.reroute_fix_curvature){const t=e.parentElement.querySelectorAll(".main-path").length;var a=document.createElementNS("http://www.w3.org/2000/svg","path");if(a.classList.add("main-path"),a.setAttributeNS(null,"d",""),e.parentElement.insertBefore(a,e.parentElement.children[t]),1===t)e.parentElement.appendChild(o);else{const n=Array.from(e.parentElement.children).indexOf(e);d=n,e.parentElement.insertBefore(o,e.parentElement.children[n+t+1])}}else e.parentElement.appendChild(o);const r=t.slice(5),h=this.drawflow.drawflow[this.module].data[r].outputs[i].connections.findIndex((function(e,t){return e.node===n&&e.output===s}));void 0===this.drawflow.drawflow[this.module].data[r].outputs[i].connections[h].points&&(this.drawflow.drawflow[this.module].data[r].outputs[i].connections[h].points=[]),this.reroute_fix_curvature?(d>0||this.drawflow.drawflow[this.module].data[r].outputs[i].connections[h].points!==[]?this.drawflow.drawflow[this.module].data[r].outputs[i].connections[h].points.splice(d,0,{pos_x:l,pos_y:c}):this.drawflow.drawflow[this.module].data[r].outputs[i].connections[h].points.push({pos_x:l,pos_y:c}),e.parentElement.querySelectorAll(".main-path").forEach((e,t)=>{e.classList.remove("selected")})):this.drawflow.drawflow[this.module].data[r].outputs[i].connections[h].points.push({pos_x:l,pos_y:c}),this.dispatch("addReroute",r),this.updateConnectionNodes(t)}removeReroutePoint(e){const t=e.parentElement.classList[2].slice(9),n=e.parentElement.classList[1].slice(13),i=e.parentElement.classList[3],s=e.parentElement.classList[4];let o=Array.from(e.parentElement.children).indexOf(e);const l=t.slice(5),c=this.drawflow.drawflow[this.module].data[l].outputs[i].connections.findIndex((function(e,t){return e.node===n&&e.output===s}));if(this.reroute_fix_curvature){const t=e.parentElement.querySelectorAll(".main-path").length;e.parentElement.children[t-1].remove(),o-=t,o<0&&(o=0)}else o--;this.drawflow.drawflow[this.module].data[l].outputs[i].connections[c].points.splice(o,1),e.remove(),this.dispatch("removeReroute",l),this.updateConnectionNodes(t)}registerNode(e,t,n=null,i=null){this.noderegister[e]={html:t,props:n,options:i}}getNodeFromId(e){var t=this.getModuleFromNodeId(e);return JSON.parse(JSON.stringify(this.drawflow.drawflow[t].data[e]))}getNodesFromName(e){var t=[];const n=this.drawflow.drawflow;return Object.keys(n).map((function(i,s){for(var o in n[i].data)n[i].data[o].name==e&&t.push(n[i].data[o].id)})),t}addNode(e,t,n,i,s,o,l,c,d=!1){if(this.useuuid)var a=this.getUuid();else a=this.nodeId;const r=document.createElement("div");r.classList.add("parent-node");const h=document.createElement("div");h.innerHTML="",h.setAttribute("id","node-"+a),h.classList.add("drawflow-node"),""!=o&&h.classList.add(...o.split(" "));const u=document.createElement("div");u.classList.add("inputs");const p=document.createElement("div");p.classList.add("outputs");const f={};for(var m=0;m<t;m++){const e=document.createElement("div");e.classList.add("input"),e.classList.add("input_"+(m+1)),f["input_"+(m+1)]={connections:[]},u.appendChild(e)}const g={};for(m=0;m<n;m++){const e=document.createElement("div");e.classList.add("output"),e.classList.add("output_"+(m+1)),g["output_"+(m+1)]={connections:[]},p.appendChild(e)}const _=document.createElement("div");if(_.classList.add("drawflow_content_node"),!1===d)_.innerHTML=c;else if(!0===d)_.appendChild(this.noderegister[c].html.cloneNode(!0));else if(3===parseInt(this.render.version)){let e=this.render.h(this.noderegister[c].html,this.noderegister[c].props,this.noderegister[c].options);e.appContext=this.parent,this.render.render(e,_)}else{let e=new this.render({parent:this.parent,render:e=>e(this.noderegister[c].html,{props:this.noderegister[c].props}),...this.noderegister[c].options}).$mount();_.appendChild(e.$el)}Object.entries(l).forEach((function(e,t){if("object"==typeof e[1])!function e(t,n,i){if(null===t)t=l[n];else t=t[n];null!==t&&Object.entries(t).forEach((function(n,s){if("object"==typeof n[1])e(t,n[0],i+"-"+n[0]);else for(var o=_.querySelectorAll("[df-"+i+"-"+n[0]+"]"),l=0;l<o.length;l++)o[l].value=n[1],o[l].isContentEditable&&(o[l].innerText=n[1])}))}(null,e[0],e[0]);else for(var n=_.querySelectorAll("[df-"+e[0]+"]"),i=0;i<n.length;i++)n[i].value=e[1],n[i].isContentEditable&&(n[i].innerText=e[1])})),h.appendChild(u),h.appendChild(_),h.appendChild(p),h.style.top=s+"px",h.style.left=i+"px",r.appendChild(h),this.precanvas.appendChild(r);var w={id:a,name:e,data:l,class:o,html:c,typenode:d,inputs:f,outputs:g,pos_x:i,pos_y:s};return this.drawflow.drawflow[this.module].data[a]=w,this.dispatch("nodeCreated",a),this.useuuid||this.nodeId++,a}addNodeImport(e,t){const n=document.createElement("div");n.classList.add("parent-node");const i=document.createElement("div");i.innerHTML="",i.setAttribute("id","node-"+e.id),i.classList.add("drawflow-node"),""!=e.class&&i.classList.add(...e.class.split(" "));const s=document.createElement("div");s.classList.add("inputs");const o=document.createElement("div");o.classList.add("outputs"),Object.keys(e.inputs).map((function(n,i){const o=document.createElement("div");o.classList.add("input"),o.classList.add(n),s.appendChild(o),Object.keys(e.inputs[n].connections).map((function(i,s){var o=document.createElementNS("http://www.w3.org/2000/svg","svg"),l=document.createElementNS("http://www.w3.org/2000/svg","path");l.classList.add("main-path"),l.setAttributeNS(null,"d",""),o.classList.add("connection"),o.classList.add("node_in_node-"+e.id),o.classList.add("node_out_node-"+e.inputs[n].connections[i].node),o.classList.add(e.inputs[n].connections[i].input),o.classList.add(n),o.appendChild(l),t.appendChild(o)}))}));for(var l=0;l<Object.keys(e.outputs).length;l++){const e=document.createElement("div");e.classList.add("output"),e.classList.add("output_"+(l+1)),o.appendChild(e)}const c=document.createElement("div");if(c.classList.add("drawflow_content_node"),!1===e.typenode)c.innerHTML=e.html;else if(!0===e.typenode)c.appendChild(this.noderegister[e.html].html.cloneNode(!0));else if(3===parseInt(this.render.version)){let t=this.render.h(this.noderegister[e.html].html,this.noderegister[e.html].props,this.noderegister[e.html].options);t.appContext=this.parent,this.render.render(t,c)}else{let t=new this.render({parent:this.parent,render:t=>t(this.noderegister[e.html].html,{props:this.noderegister[e.html].props}),...this.noderegister[e.html].options}).$mount();c.appendChild(t.$el)}Object.entries(e.data).forEach((function(t,n){if("object"==typeof t[1])!function t(n,i,s){if(null===n)n=e.data[i];else n=n[i];null!==n&&Object.entries(n).forEach((function(e,i){if("object"==typeof e[1])t(n,e[0],s+"-"+e[0]);else for(var o=c.querySelectorAll("[df-"+s+"-"+e[0]+"]"),l=0;l<o.length;l++)o[l].value=e[1],o[l].isContentEditable&&(o[l].innerText=e[1])}))}(null,t[0],t[0]);else for(var i=c.querySelectorAll("[df-"+t[0]+"]"),s=0;s<i.length;s++)i[s].value=t[1],i[s].isContentEditable&&(i[s].innerText=t[1])})),i.appendChild(s),i.appendChild(c),i.appendChild(o),i.style.top=e.pos_y+"px",i.style.left=e.pos_x+"px",n.appendChild(i),this.precanvas.appendChild(n)}addRerouteImport(e){const t=this.reroute_width,n=this.reroute_fix_curvature,i=this.container;Object.keys(e.outputs).map((function(s,o){Object.keys(e.outputs[s].connections).map((function(o,l){const c=e.outputs[s].connections[o].points;void 0!==c&&c.forEach((l,d)=>{const a=e.outputs[s].connections[o].node,r=e.outputs[s].connections[o].output,h=i.querySelector(".connection.node_in_node-"+a+".node_out_node-"+e.id+"."+s+"."+r);if(n&&0===d)for(var u=0;u<c.length;u++){var p=document.createElementNS("http://www.w3.org/2000/svg","path");p.classList.add("main-path"),p.setAttributeNS(null,"d",""),h.appendChild(p)}const f=document.createElementNS("http://www.w3.org/2000/svg","circle");f.classList.add("point");var m=l.pos_x,g=l.pos_y;f.setAttributeNS(null,"cx",m),f.setAttributeNS(null,"cy",g),f.setAttributeNS(null,"r",t),h.appendChild(f)})}))}))}updateNodeValue(e){for(var t=e.target.attributes,n=0;n<t.length;n++)if(t[n].nodeName.startsWith("df-")){for(var i=t[n].nodeName.slice(3).split("-"),s=this.drawflow.drawflow[this.module].data[e.target.closest(".drawflow_content_node").parentElement.id.slice(5)].data,o=0;o<i.length-1;o+=1)null==s[i[o]]&&(s[i[o]]={}),s=s[i[o]];s[i[i.length-1]]=e.target.value,e.target.isContentEditable&&(s[i[i.length-1]]=e.target.innerText),this.dispatch("nodeDataChanged",e.target.closest(".drawflow_content_node").parentElement.id.slice(5))}}updateNodeDataFromId(e,t){var n=this.getModuleFromNodeId(e);if(this.drawflow.drawflow[n].data[e].data=t,this.module===n){const n=this.container.querySelector("#node-"+e);Object.entries(t).forEach((function(e,i){if("object"==typeof e[1])!function e(i,s,o){if(null===i)i=t[s];else i=i[s];null!==i&&Object.entries(i).forEach((function(t,s){if("object"==typeof t[1])e(i,t[0],o+"-"+t[0]);else for(var l=n.querySelectorAll("[df-"+o+"-"+t[0]+"]"),c=0;c<l.length;c++)l[c].value=t[1],l[c].isContentEditable&&(l[c].innerText=t[1])}))}(null,e[0],e[0]);else for(var s=n.querySelectorAll("[df-"+e[0]+"]"),o=0;o<s.length;o++)s[o].value=e[1],s[o].isContentEditable&&(s[o].innerText=e[1])}))}}addNodeInput(e){var t=this.getModuleFromNodeId(e);const n=this.getNodeFromId(e),i=Object.keys(n.inputs).length;if(this.module===t){const t=document.createElement("div");t.classList.add("input"),t.classList.add("input_"+(i+1)),this.container.querySelector("#node-"+e+" .inputs").appendChild(t),this.updateConnectionNodes("node-"+e)}this.drawflow.drawflow[t].data[e].inputs["input_"+(i+1)]={connections:[]}}addNodeOutput(e){var t=this.getModuleFromNodeId(e);const n=this.getNodeFromId(e),i=Object.keys(n.outputs).length;if(this.module===t){const t=document.createElement("div");t.classList.add("output"),t.classList.add("output_"+(i+1)),this.container.querySelector("#node-"+e+" .outputs").appendChild(t),this.updateConnectionNodes("node-"+e)}this.drawflow.drawflow[t].data[e].outputs["output_"+(i+1)]={connections:[]}}removeNodeInput(e,t){var n=this.getModuleFromNodeId(e);const i=this.getNodeFromId(e);this.module===n&&this.container.querySelector("#node-"+e+" .inputs .input."+t).remove();const s=[];Object.keys(i.inputs[t].connections).map((function(n,o){const l=i.inputs[t].connections[o].node,c=i.inputs[t].connections[o].input;s.push({id_output:l,id:e,output_class:c,input_class:t})})),s.forEach((e,t)=>{this.removeSingleConnection(e.id_output,e.id,e.output_class,e.input_class)}),delete this.drawflow.drawflow[n].data[e].inputs[t];const o=[],l=this.drawflow.drawflow[n].data[e].inputs;Object.keys(l).map((function(e,t){o.push(l[e])})),this.drawflow.drawflow[n].data[e].inputs={};const c=t.slice(6);let d=[];if(o.forEach((t,i)=>{t.connections.forEach((e,t)=>{d.push(e)}),this.drawflow.drawflow[n].data[e].inputs["input_"+(i+1)]=t}),d=new Set(d.map(e=>JSON.stringify(e))),d=Array.from(d).map(e=>JSON.parse(e)),this.module===n){this.container.querySelectorAll("#node-"+e+" .inputs .input").forEach((e,t)=>{const n=e.classList[1].slice(6);parseInt(c)<parseInt(n)&&(e.classList.remove("input_"+n),e.classList.add("input_"+(n-1)))})}d.forEach((t,i)=>{this.drawflow.drawflow[n].data[t.node].outputs[t.input].connections.forEach((i,s)=>{if(i.node==e){const o=i.output.slice(6);if(parseInt(c)<parseInt(o)){if(this.module===n){const n=this.container.querySelector(".connection.node_in_node-"+e+".node_out_node-"+t.node+"."+t.input+".input_"+o);n.classList.remove("input_"+o),n.classList.add("input_"+(o-1))}i.points?this.drawflow.drawflow[n].data[t.node].outputs[t.input].connections[s]={node:i.node,output:"input_"+(o-1),points:i.points}:this.drawflow.drawflow[n].data[t.node].outputs[t.input].connections[s]={node:i.node,output:"input_"+(o-1)}}}})}),this.updateConnectionNodes("node-"+e)}removeNodeOutput(e,t){var n=this.getModuleFromNodeId(e);const i=this.getNodeFromId(e);this.module===n&&this.container.querySelector("#node-"+e+" .outputs .output."+t).remove();const s=[];Object.keys(i.outputs[t].connections).map((function(n,o){const l=i.outputs[t].connections[o].node,c=i.outputs[t].connections[o].output;s.push({id:e,id_input:l,output_class:t,input_class:c})})),s.forEach((e,t)=>{this.removeSingleConnection(e.id,e.id_input,e.output_class,e.input_class)}),delete this.drawflow.drawflow[n].data[e].outputs[t];const o=[],l=this.drawflow.drawflow[n].data[e].outputs;Object.keys(l).map((function(e,t){o.push(l[e])})),this.drawflow.drawflow[n].data[e].outputs={};const c=t.slice(7);let d=[];if(o.forEach((t,i)=>{t.connections.forEach((e,t)=>{d.push({node:e.node,output:e.output})}),this.drawflow.drawflow[n].data[e].outputs["output_"+(i+1)]=t}),d=new Set(d.map(e=>JSON.stringify(e))),d=Array.from(d).map(e=>JSON.parse(e)),this.module===n){this.container.querySelectorAll("#node-"+e+" .outputs .output").forEach((e,t)=>{const n=e.classList[1].slice(7);parseInt(c)<parseInt(n)&&(e.classList.remove("output_"+n),e.classList.add("output_"+(n-1)))})}d.forEach((t,i)=>{this.drawflow.drawflow[n].data[t.node].inputs[t.output].connections.forEach((i,s)=>{if(i.node==e){const o=i.input.slice(7);if(parseInt(c)<parseInt(o)){if(this.module===n){const n=this.container.querySelector(".connection.node_in_node-"+t.node+".node_out_node-"+e+".output_"+o+"."+t.output);n.classList.remove("output_"+o),n.classList.remove(t.output),n.classList.add("output_"+(o-1)),n.classList.add(t.output)}i.points?this.drawflow.drawflow[n].data[t.node].inputs[t.output].connections[s]={node:i.node,input:"output_"+(o-1),points:i.points}:this.drawflow.drawflow[n].data[t.node].inputs[t.output].connections[s]={node:i.node,input:"output_"+(o-1)}}}})}),this.updateConnectionNodes("node-"+e)}removeNodeId(e){this.removeConnectionNodeId(e);var t=this.getModuleFromNodeId(e.slice(5));this.module===t&&this.container.querySelector("#"+e).remove(),delete this.drawflow.drawflow[t].data[e.slice(5)],this.dispatch("nodeRemoved",e.slice(5))}removeConnection(){if(null!=this.connection_selected){var e=this.connection_selected.parentElement.classList;this.connection_selected.parentElement.remove();var t=this.drawflow.drawflow[this.module].data[e[2].slice(14)].outputs[e[3]].connections.findIndex((function(t,n){return t.node===e[1].slice(13)&&t.output===e[4]}));this.drawflow.drawflow[this.module].data[e[2].slice(14)].outputs[e[3]].connections.splice(t,1);var n=this.drawflow.drawflow[this.module].data[e[1].slice(13)].inputs[e[4]].connections.findIndex((function(t,n){return t.node===e[2].slice(14)&&t.input===e[3]}));this.drawflow.drawflow[this.module].data[e[1].slice(13)].inputs[e[4]].connections.splice(n,1),this.dispatch("connectionRemoved",{output_id:e[2].slice(14),input_id:e[1].slice(13),output_class:e[3],input_class:e[4]}),this.connection_selected=null}}removeSingleConnection(e,t,n,i){var s=this.getModuleFromNodeId(e);if(s===this.getModuleFromNodeId(t)){if(this.drawflow.drawflow[s].data[e].outputs[n].connections.findIndex((function(e,n){return e.node==t&&e.output===i}))>-1){this.module===s&&this.container.querySelector(".connection.node_in_node-"+t+".node_out_node-"+e+"."+n+"."+i).remove();var o=this.drawflow.drawflow[s].data[e].outputs[n].connections.findIndex((function(e,n){return e.node==t&&e.output===i}));this.drawflow.drawflow[s].data[e].outputs[n].connections.splice(o,1);var l=this.drawflow.drawflow[s].data[t].inputs[i].connections.findIndex((function(t,i){return t.node==e&&t.input===n}));return this.drawflow.drawflow[s].data[t].inputs[i].connections.splice(l,1),this.dispatch("connectionRemoved",{output_id:e,input_id:t,output_class:n,input_class:i}),!0}return!1}return!1}removeConnectionNodeId(e){const t="node_in_"+e,n="node_out_"+e,i=this.container.querySelectorAll("."+n);for(var s=i.length-1;s>=0;s--){var o=i[s].classList,l=this.drawflow.drawflow[this.module].data[o[1].slice(13)].inputs[o[4]].connections.findIndex((function(e,t){return e.node===o[2].slice(14)&&e.input===o[3]}));this.drawflow.drawflow[this.module].data[o[1].slice(13)].inputs[o[4]].connections.splice(l,1);var c=this.drawflow.drawflow[this.module].data[o[2].slice(14)].outputs[o[3]].connections.findIndex((function(e,t){return e.node===o[1].slice(13)&&e.output===o[4]}));this.drawflow.drawflow[this.module].data[o[2].slice(14)].outputs[o[3]].connections.splice(c,1),i[s].remove(),this.dispatch("connectionRemoved",{output_id:o[2].slice(14),input_id:o[1].slice(13),output_class:o[3],input_class:o[4]})}const d=this.container.querySelectorAll("."+t);for(s=d.length-1;s>=0;s--){o=d[s].classList,c=this.drawflow.drawflow[this.module].data[o[2].slice(14)].outputs[o[3]].connections.findIndex((function(e,t){return e.node===o[1].slice(13)&&e.output===o[4]}));this.drawflow.drawflow[this.module].data[o[2].slice(14)].outputs[o[3]].connections.splice(c,1);l=this.drawflow.drawflow[this.module].data[o[1].slice(13)].inputs[o[4]].connections.findIndex((function(e,t){return e.node===o[2].slice(14)&&e.input===o[3]}));this.drawflow.drawflow[this.module].data[o[1].slice(13)].inputs[o[4]].connections.splice(l,1),d[s].remove(),this.dispatch("connectionRemoved",{output_id:o[2].slice(14),input_id:o[1].slice(13),output_class:o[3],input_class:o[4]})}}getModuleFromNodeId(e){var t;const n=this.drawflow.drawflow;return Object.keys(n).map((function(i,s){Object.keys(n[i].data).map((function(n,s){n==e&&(t=i)}))})),t}addModule(e){this.drawflow.drawflow[e]={data:{}},this.dispatch("moduleCreated",e)}changeModule(e){this.dispatch("moduleChanged",e),this.module=e,this.precanvas.innerHTML="",this.canvas_x=0,this.canvas_y=0,this.pos_x=0,this.pos_y=0,this.mouse_x=0,this.mouse_y=0,this.zoom=1,this.zoom_last_value=1,this.precanvas.style.transform="",this.import(this.drawflow,!1)}removeModule(e){this.module===e&&this.changeModule("Home"),delete this.drawflow.drawflow[e],this.dispatch("moduleRemoved",e)}clearModuleSelected(){this.precanvas.innerHTML="",this.drawflow.drawflow[this.module]={data:{}}}clear(){this.precanvas.innerHTML="",this.drawflow={drawflow:{Home:{data:{}}}}}export(){const e=JSON.parse(JSON.stringify(this.drawflow));return this.dispatch("export",e),e}import(e,t=!0){this.clear(),this.drawflow=JSON.parse(JSON.stringify(e)),this.load(),t&&this.dispatch("import","import")}on(e,t){return"function"!=typeof t?(console.error("The listener callback must be a function, the given type is "+typeof t),!1):"string"!=typeof e?(console.error("The event name must be a string, the given type is "+typeof e),!1):(void 0===this.events[e]&&(this.events[e]={listeners:[]}),void this.events[e].listeners.push(t))}removeListener(e,t){if(!this.events[e])return!1;const n=this.events[e].listeners,i=n.indexOf(t);i>-1&&n.splice(i,1)}dispatch(e,t){if(void 0===this.events[e])return!1;this.events[e].listeners.forEach(e=>{e(t)})}getUuid(){for(var e=[],t=0;t<36;t++)e[t]="0123456789abcdef".substr(Math.floor(16*Math.random()),1);return e[14]="4",e[19]="0123456789abcdef".substr(3&e[19]|8,1),e[8]=e[13]=e[18]=e[23]="-",e.join("")}}}]).default}));

    [File Ends] public/js/drawflow.min.js

    [File Begins] public/js/events.js
    var promptViewHeight
    
    // Function to be called when models-container is added to the DOM
    function onModelsContainerAdded() {
      console.log("Model cards loaded");
      setViewHeight("models-container");
      highlightSelectedModels();
    }
    
    // Callback function to execute when mutations are observed
    const mutationCallback = (mutationsList, observer) => {
      for (const mutation of mutationsList) {
        if (mutation.type === 'childList') {
          for (const node of mutation.addedNodes) {
            if (node.nodeName === 'DIV' && node.getAttribute('id') === "models-container") {
              onModelsContainerAdded();
              // Do not disconnect the observer
            }
          }
        }
      }
    };
    
    async function highlightSelectedModels() {
      console.log("Highlighting selected models...")
      try {
        // Fetch the selected model names from the server
        const response = await fetch('/models/selected');
        if (!response.ok) {
          throw new Error('Failed to fetch selected models');
        }
        const selectedModelNames = await response.json();
    
        // Apply 'card-selected' class to each selected model's card
        selectedModelNames.forEach(modelName => {
          const modelCard = document.querySelector(`[data-model-name="${modelName}"]`);
          if (modelCard) {
            modelCard.classList.add('card-selected');
          }
        });
      } catch (error) {
        console.error('Error:', error);
      }
    }
    
    // MODELS HANDLERS
    // Define a global variable to store selected models
    let selectedModels = [];
    
    // Call getSelectedModels when the DOM content is fully loaded
    document.addEventListener('DOMContentLoaded', () => {
    
      promptViewHeight = document.getElementById('prompt-view').offsetHeight;
      getSelectedModels();
      const observer = new MutationObserver(mutationCallback);
    
      // Start observing the document body for DOM changes
      observer.observe(document.body, { childList: true, subtree: true });
    
      let userHasScrolled = false;
    
      // Attach the event listener to the window object for scroll events
      window.addEventListener('scroll', () => {
        // If the user is not at the bottom of the page, update the flag
        userHasScrolled = (window.innerHeight + window.scrollY) < document.body.offsetHeight;
      });
    });
    
    async function getSelectedModels() {
      try {
        // Fetch the selected model names from the server
        const response = await fetch('/models/selected');
        if (!response.ok) {
          throw new Error('Failed to fetch selected models');
        }
        selectedModels = await response.json();
    
        // Log the fetched models for debugging
        console.log("Selected Models:", selectedModels);
      } catch (error) {
        console.error('Error:', error);
      }
    }
    
    // Call this function to reset the scroll behavior when the user clicks a "scroll to bottom" button
    function resetScroll() {
      userHasScrolled = false;
      scrollToBottomOfPage();
    }
    
    function scrollToBottomOfPage() {
      if (!userHasScrolled) {
        requestAnimationFrame(() => {
          // Scroll to the bottom of the document body
          window.scrollTo(0, document.body.scrollHeight);
        });
      }
    }
    
    function toastDisplay(toastId, toastMessage) {
      const toastLiveExample = document.getElementById(toastId);
      toastLiveExample.querySelector('#toast-message').innerText = toastMessage;
      const toastBootstrap = bootstrap.Toast.getOrCreateInstance(toastLiveExample);
      toastBootstrap.show();
    }
    
    function setViewHeight(viewId) {
      // Get the height of the prompt view
      // var viewHeight = document.getElementById(viewId).offsetHeight;
    
      // Define additional spacing in pixels
      var additionalSpacing = 10; // Additional spacing in pixels
    
      // Set the bottom padding of the body to ensure the prompt view doesn't cover content
      // Subtract the additionalSpacing from promptViewHeight
      document.body.style.paddingBottom = (promptViewHeight + additionalSpacing) + 'px';
    }
    
      // Listen for the event with the name 'message' and update the div#sse-messages
      document.body.addEventListener('htmx:sseAfterMessage', function(event) {
        console.log('Received event:', event);
        var sseData = JSON.parse(event.detail.data);
        var sseMessages = document.getElementById('sse-messages');
        var messageElement = document.createElement('div');
        messageElement.textContent = sseData.message + ' (received at ' + sseData.timestamp + ')';
        sseMessages.appendChild(messageElement);
      });

    [File Ends] public/js/events.js

    [File Begins] public/js/htmx.min.js
    (function(e,t){if(typeof define==="function"&&define.amd){define([],t)}else if(typeof module==="object"&&module.exports){module.exports=t()}else{e.htmx=e.htmx||t()}})(typeof self!=="undefined"?self:this,function(){return function(){"use strict";var Q={onLoad:F,process:zt,on:de,off:ge,trigger:ce,ajax:Nr,find:C,findAll:f,closest:v,values:function(e,t){var r=dr(e,t||"post");return r.values},remove:_,addClass:z,removeClass:n,toggleClass:$,takeClass:W,defineExtension:Ur,removeExtension:Br,logAll:V,logNone:j,logger:null,config:{historyEnabled:true,historyCacheSize:10,refreshOnHistoryMiss:false,defaultSwapStyle:"innerHTML",defaultSwapDelay:0,defaultSettleDelay:20,includeIndicatorStyles:true,indicatorClass:"htmx-indicator",requestClass:"htmx-request",addedClass:"htmx-added",settlingClass:"htmx-settling",swappingClass:"htmx-swapping",allowEval:true,allowScriptTags:true,inlineScriptNonce:"",attributesToSettle:["class","style","width","height"],withCredentials:false,timeout:0,wsReconnectDelay:"full-jitter",wsBinaryType:"blob",disableSelector:"[hx-disable], [data-hx-disable]",useTemplateFragments:false,scrollBehavior:"smooth",defaultFocusScroll:false,getCacheBusterParam:false,globalViewTransitions:false,methodsThatUseUrlParams:["get"],selfRequestsOnly:false,ignoreTitle:false,scrollIntoViewOnBoost:true,triggerSpecsCache:null},parseInterval:d,_:t,createEventSource:function(e){return new EventSource(e,{withCredentials:true})},createWebSocket:function(e){var t=new WebSocket(e,[]);t.binaryType=Q.config.wsBinaryType;return t},version:"1.9.10"};var r={addTriggerHandler:Lt,bodyContains:se,canAccessLocalStorage:U,findThisElement:xe,filterValues:yr,hasAttribute:o,getAttributeValue:te,getClosestAttributeValue:ne,getClosestMatch:c,getExpressionVars:Hr,getHeaders:xr,getInputValues:dr,getInternalData:ae,getSwapSpecification:wr,getTriggerSpecs:it,getTarget:ye,makeFragment:l,mergeObjects:le,makeSettleInfo:T,oobSwap:Ee,querySelectorExt:ue,selectAndSwap:je,settleImmediately:nr,shouldCancel:ut,triggerEvent:ce,triggerErrorEvent:fe,withExtensions:R};var w=["get","post","put","delete","patch"];var i=w.map(function(e){return"[hx-"+e+"], [data-hx-"+e+"]"}).join(", ");var S=e("head"),q=e("title"),H=e("svg",true);function e(e,t=false){return new RegExp(`<${e}(\\s[^>]*>|>)([\\s\\S]*?)<\\/${e}>`,t?"gim":"im")}function d(e){if(e==undefined){return undefined}let t=NaN;if(e.slice(-2)=="ms"){t=parseFloat(e.slice(0,-2))}else if(e.slice(-1)=="s"){t=parseFloat(e.slice(0,-1))*1e3}else if(e.slice(-1)=="m"){t=parseFloat(e.slice(0,-1))*1e3*60}else{t=parseFloat(e)}return isNaN(t)?undefined:t}function ee(e,t){return e.getAttribute&&e.getAttribute(t)}function o(e,t){return e.hasAttribute&&(e.hasAttribute(t)||e.hasAttribute("data-"+t))}function te(e,t){return ee(e,t)||ee(e,"data-"+t)}function u(e){return e.parentElement}function re(){return document}function c(e,t){while(e&&!t(e)){e=u(e)}return e?e:null}function L(e,t,r){var n=te(t,r);var i=te(t,"hx-disinherit");if(e!==t&&i&&(i==="*"||i.split(" ").indexOf(r)>=0)){return"unset"}else{return n}}function ne(t,r){var n=null;c(t,function(e){return n=L(t,e,r)});if(n!=="unset"){return n}}function h(e,t){var r=e.matches||e.matchesSelector||e.msMatchesSelector||e.mozMatchesSelector||e.webkitMatchesSelector||e.oMatchesSelector;return r&&r.call(e,t)}function A(e){var t=/<([a-z][^\/\0>\x20\t\r\n\f]*)/i;var r=t.exec(e);if(r){return r[1].toLowerCase()}else{return""}}function a(e,t){var r=new DOMParser;var n=r.parseFromString(e,"text/html");var i=n.body;while(t>0){t--;i=i.firstChild}if(i==null){i=re().createDocumentFragment()}return i}function N(e){return/<body/.test(e)}function l(e){var t=!N(e);var r=A(e);var n=e;if(r==="head"){n=n.replace(S,"")}if(Q.config.useTemplateFragments&&t){var i=a("<body><template>"+n+"</template></body>",0);return i.querySelector("template").content}switch(r){case"thead":case"tbody":case"tfoot":case"colgroup":case"caption":return a("<table>"+n+"</table>",1);case"col":return a("<table><colgroup>"+n+"</colgroup></table>",2);case"tr":return a("<table><tbody>"+n+"</tbody></table>",2);case"td":case"th":return a("<table><tbody><tr>"+n+"</tr></tbody></table>",3);case"script":case"style":return a("<div>"+n+"</div>",1);default:return a(n,0)}}function ie(e){if(e){e()}}function I(e,t){return Object.prototype.toString.call(e)==="[object "+t+"]"}function k(e){return I(e,"Function")}function P(e){return I(e,"Object")}function ae(e){var t="htmx-internal-data";var r=e[t];if(!r){r=e[t]={}}return r}function M(e){var t=[];if(e){for(var r=0;r<e.length;r++){t.push(e[r])}}return t}function oe(e,t){if(e){for(var r=0;r<e.length;r++){t(e[r])}}}function X(e){var t=e.getBoundingClientRect();var r=t.top;var n=t.bottom;return r<window.innerHeight&&n>=0}function se(e){if(e.getRootNode&&e.getRootNode()instanceof window.ShadowRoot){return re().body.contains(e.getRootNode().host)}else{return re().body.contains(e)}}function D(e){return e.trim().split(/\s+/)}function le(e,t){for(var r in t){if(t.hasOwnProperty(r)){e[r]=t[r]}}return e}function E(e){try{return JSON.parse(e)}catch(e){b(e);return null}}function U(){var e="htmx:localStorageTest";try{localStorage.setItem(e,e);localStorage.removeItem(e);return true}catch(e){return false}}function B(t){try{var e=new URL(t);if(e){t=e.pathname+e.search}if(!/^\/$/.test(t)){t=t.replace(/\/+$/,"")}return t}catch(e){return t}}function t(e){return Tr(re().body,function(){return eval(e)})}function F(t){var e=Q.on("htmx:load",function(e){t(e.detail.elt)});return e}function V(){Q.logger=function(e,t,r){if(console){console.log(t,e,r)}}}function j(){Q.logger=null}function C(e,t){if(t){return e.querySelector(t)}else{return C(re(),e)}}function f(e,t){if(t){return e.querySelectorAll(t)}else{return f(re(),e)}}function _(e,t){e=g(e);if(t){setTimeout(function(){_(e);e=null},t)}else{e.parentElement.removeChild(e)}}function z(e,t,r){e=g(e);if(r){setTimeout(function(){z(e,t);e=null},r)}else{e.classList&&e.classList.add(t)}}function n(e,t,r){e=g(e);if(r){setTimeout(function(){n(e,t);e=null},r)}else{if(e.classList){e.classList.remove(t);if(e.classList.length===0){e.removeAttribute("class")}}}}function $(e,t){e=g(e);e.classList.toggle(t)}function W(e,t){e=g(e);oe(e.parentElement.children,function(e){n(e,t)});z(e,t)}function v(e,t){e=g(e);if(e.closest){return e.closest(t)}else{do{if(e==null||h(e,t)){return e}}while(e=e&&u(e));return null}}function s(e,t){return e.substring(0,t.length)===t}function G(e,t){return e.substring(e.length-t.length)===t}function J(e){var t=e.trim();if(s(t,"<")&&G(t,"/>")){return t.substring(1,t.length-2)}else{return t}}function Z(e,t){if(t.indexOf("closest ")===0){return[v(e,J(t.substr(8)))]}else if(t.indexOf("find ")===0){return[C(e,J(t.substr(5)))]}else if(t==="next"){return[e.nextElementSibling]}else if(t.indexOf("next ")===0){return[K(e,J(t.substr(5)))]}else if(t==="previous"){return[e.previousElementSibling]}else if(t.indexOf("previous ")===0){return[Y(e,J(t.substr(9)))]}else if(t==="document"){return[document]}else if(t==="window"){return[window]}else if(t==="body"){return[document.body]}else{return re().querySelectorAll(J(t))}}var K=function(e,t){var r=re().querySelectorAll(t);for(var n=0;n<r.length;n++){var i=r[n];if(i.compareDocumentPosition(e)===Node.DOCUMENT_POSITION_PRECEDING){return i}}};var Y=function(e,t){var r=re().querySelectorAll(t);for(var n=r.length-1;n>=0;n--){var i=r[n];if(i.compareDocumentPosition(e)===Node.DOCUMENT_POSITION_FOLLOWING){return i}}};function ue(e,t){if(t){return Z(e,t)[0]}else{return Z(re().body,e)[0]}}function g(e){if(I(e,"String")){return C(e)}else{return e}}function ve(e,t,r){if(k(t)){return{target:re().body,event:e,listener:t}}else{return{target:g(e),event:t,listener:r}}}function de(t,r,n){jr(function(){var e=ve(t,r,n);e.target.addEventListener(e.event,e.listener)});var e=k(r);return e?r:n}function ge(t,r,n){jr(function(){var e=ve(t,r,n);e.target.removeEventListener(e.event,e.listener)});return k(r)?r:n}var me=re().createElement("output");function pe(e,t){var r=ne(e,t);if(r){if(r==="this"){return[xe(e,t)]}else{var n=Z(e,r);if(n.length===0){b('The selector "'+r+'" on '+t+" returned no matches!");return[me]}else{return n}}}}function xe(e,t){return c(e,function(e){return te(e,t)!=null})}function ye(e){var t=ne(e,"hx-target");if(t){if(t==="this"){return xe(e,"hx-target")}else{return ue(e,t)}}else{var r=ae(e);if(r.boosted){return re().body}else{return e}}}function be(e){var t=Q.config.attributesToSettle;for(var r=0;r<t.length;r++){if(e===t[r]){return true}}return false}function we(t,r){oe(t.attributes,function(e){if(!r.hasAttribute(e.name)&&be(e.name)){t.removeAttribute(e.name)}});oe(r.attributes,function(e){if(be(e.name)){t.setAttribute(e.name,e.value)}})}function Se(e,t){var r=Fr(t);for(var n=0;n<r.length;n++){var i=r[n];try{if(i.isInlineSwap(e)){return true}}catch(e){b(e)}}return e==="outerHTML"}function Ee(e,i,a){var t="#"+ee(i,"id");var o="outerHTML";if(e==="true"){}else if(e.indexOf(":")>0){o=e.substr(0,e.indexOf(":"));t=e.substr(e.indexOf(":")+1,e.length)}else{o=e}var r=re().querySelectorAll(t);if(r){oe(r,function(e){var t;var r=i.cloneNode(true);t=re().createDocumentFragment();t.appendChild(r);if(!Se(o,e)){t=r}var n={shouldSwap:true,target:e,fragment:t};if(!ce(e,"htmx:oobBeforeSwap",n))return;e=n.target;if(n["shouldSwap"]){Fe(o,e,e,t,a)}oe(a.elts,function(e){ce(e,"htmx:oobAfterSwap",n)})});i.parentNode.removeChild(i)}else{i.parentNode.removeChild(i);fe(re().body,"htmx:oobErrorNoTarget",{content:i})}return e}function Ce(e,t,r){var n=ne(e,"hx-select-oob");if(n){var i=n.split(",");for(var a=0;a<i.length;a++){var o=i[a].split(":",2);var s=o[0].trim();if(s.indexOf("#")===0){s=s.substring(1)}var l=o[1]||"true";var u=t.querySelector("#"+s);if(u){Ee(l,u,r)}}}oe(f(t,"[hx-swap-oob], [data-hx-swap-oob]"),function(e){var t=te(e,"hx-swap-oob");if(t!=null){Ee(t,e,r)}})}function Re(e){oe(f(e,"[hx-preserve], [data-hx-preserve]"),function(e){var t=te(e,"id");var r=re().getElementById(t);if(r!=null){e.parentNode.replaceChild(r,e)}})}function Te(o,e,s){oe(e.querySelectorAll("[id]"),function(e){var t=ee(e,"id");if(t&&t.length>0){var r=t.replace("'","\\'");var n=e.tagName.replace(":","\\:");var i=o.querySelector(n+"[id='"+r+"']");if(i&&i!==o){var a=e.cloneNode();we(e,i);s.tasks.push(function(){we(e,a)})}}})}function Oe(e){return function(){n(e,Q.config.addedClass);zt(e);Nt(e);qe(e);ce(e,"htmx:load")}}function qe(e){var t="[autofocus]";var r=h(e,t)?e:e.querySelector(t);if(r!=null){r.focus()}}function m(e,t,r,n){Te(e,r,n);while(r.childNodes.length>0){var i=r.firstChild;z(i,Q.config.addedClass);e.insertBefore(i,t);if(i.nodeType!==Node.TEXT_NODE&&i.nodeType!==Node.COMMENT_NODE){n.tasks.push(Oe(i))}}}function He(e,t){var r=0;while(r<e.length){t=(t<<5)-t+e.charCodeAt(r++)|0}return t}function Le(e){var t=0;if(e.attributes){for(var r=0;r<e.attributes.length;r++){var n=e.attributes[r];if(n.value){t=He(n.name,t);t=He(n.value,t)}}}return t}function Ae(e){var t=ae(e);if(t.onHandlers){for(var r=0;r<t.onHandlers.length;r++){const n=t.onHandlers[r];e.removeEventListener(n.event,n.listener)}delete t.onHandlers}}function Ne(e){var t=ae(e);if(t.timeout){clearTimeout(t.timeout)}if(t.webSocket){t.webSocket.close()}if(t.sseEventSource){t.sseEventSource.close()}if(t.listenerInfos){oe(t.listenerInfos,function(e){if(e.on){e.on.removeEventListener(e.trigger,e.listener)}})}Ae(e);oe(Object.keys(t),function(e){delete t[e]})}function p(e){ce(e,"htmx:beforeCleanupElement");Ne(e);if(e.children){oe(e.children,function(e){p(e)})}}function Ie(t,e,r){if(t.tagName==="BODY"){return Ue(t,e,r)}else{var n;var i=t.previousSibling;m(u(t),t,e,r);if(i==null){n=u(t).firstChild}else{n=i.nextSibling}r.elts=r.elts.filter(function(e){return e!=t});while(n&&n!==t){if(n.nodeType===Node.ELEMENT_NODE){r.elts.push(n)}n=n.nextElementSibling}p(t);u(t).removeChild(t)}}function ke(e,t,r){return m(e,e.firstChild,t,r)}function Pe(e,t,r){return m(u(e),e,t,r)}function Me(e,t,r){return m(e,null,t,r)}function Xe(e,t,r){return m(u(e),e.nextSibling,t,r)}function De(e,t,r){p(e);return u(e).removeChild(e)}function Ue(e,t,r){var n=e.firstChild;m(e,n,t,r);if(n){while(n.nextSibling){p(n.nextSibling);e.removeChild(n.nextSibling)}p(n);e.removeChild(n)}}function Be(e,t,r){var n=r||ne(e,"hx-select");if(n){var i=re().createDocumentFragment();oe(t.querySelectorAll(n),function(e){i.appendChild(e)});t=i}return t}function Fe(e,t,r,n,i){switch(e){case"none":return;case"outerHTML":Ie(r,n,i);return;case"afterbegin":ke(r,n,i);return;case"beforebegin":Pe(r,n,i);return;case"beforeend":Me(r,n,i);return;case"afterend":Xe(r,n,i);return;case"delete":De(r,n,i);return;default:var a=Fr(t);for(var o=0;o<a.length;o++){var s=a[o];try{var l=s.handleSwap(e,r,n,i);if(l){if(typeof l.length!=="undefined"){for(var u=0;u<l.length;u++){var f=l[u];if(f.nodeType!==Node.TEXT_NODE&&f.nodeType!==Node.COMMENT_NODE){i.tasks.push(Oe(f))}}}return}}catch(e){b(e)}}if(e==="innerHTML"){Ue(r,n,i)}else{Fe(Q.config.defaultSwapStyle,t,r,n,i)}}}function Ve(e){if(e.indexOf("<title")>-1){var t=e.replace(H,"");var r=t.match(q);if(r){return r[2]}}}function je(e,t,r,n,i,a){i.title=Ve(n);var o=l(n);if(o){Ce(r,o,i);o=Be(r,o,a);Re(o);return Fe(e,r,t,o,i)}}function _e(e,t,r){var n=e.getResponseHeader(t);if(n.indexOf("{")===0){var i=E(n);for(var a in i){if(i.hasOwnProperty(a)){var o=i[a];if(!P(o)){o={value:o}}ce(r,a,o)}}}else{var s=n.split(",");for(var l=0;l<s.length;l++){ce(r,s[l].trim(),[])}}}var ze=/\s/;var x=/[\s,]/;var $e=/[_$a-zA-Z]/;var We=/[_$a-zA-Z0-9]/;var Ge=['"',"'","/"];var Je=/[^\s]/;var Ze=/[{(]/;var Ke=/[})]/;function Ye(e){var t=[];var r=0;while(r<e.length){if($e.exec(e.charAt(r))){var n=r;while(We.exec(e.charAt(r+1))){r++}t.push(e.substr(n,r-n+1))}else if(Ge.indexOf(e.charAt(r))!==-1){var i=e.charAt(r);var n=r;r++;while(r<e.length&&e.charAt(r)!==i){if(e.charAt(r)==="\\"){r++}r++}t.push(e.substr(n,r-n+1))}else{var a=e.charAt(r);t.push(a)}r++}return t}function Qe(e,t,r){return $e.exec(e.charAt(0))&&e!=="true"&&e!=="false"&&e!=="this"&&e!==r&&t!=="."}function et(e,t,r){if(t[0]==="["){t.shift();var n=1;var i=" return (function("+r+"){ return (";var a=null;while(t.length>0){var o=t[0];if(o==="]"){n--;if(n===0){if(a===null){i=i+"true"}t.shift();i+=")})";try{var s=Tr(e,function(){return Function(i)()},function(){return true});s.source=i;return s}catch(e){fe(re().body,"htmx:syntax:error",{error:e,source:i});return null}}}else if(o==="["){n++}if(Qe(o,a,r)){i+="(("+r+"."+o+") ? ("+r+"."+o+") : (window."+o+"))"}else{i=i+o}a=t.shift()}}}function y(e,t){var r="";while(e.length>0&&!t.test(e[0])){r+=e.shift()}return r}function tt(e){var t;if(e.length>0&&Ze.test(e[0])){e.shift();t=y(e,Ke).trim();e.shift()}else{t=y(e,x)}return t}var rt="input, textarea, select";function nt(e,t,r){var n=[];var i=Ye(t);do{y(i,Je);var a=i.length;var o=y(i,/[,\[\s]/);if(o!==""){if(o==="every"){var s={trigger:"every"};y(i,Je);s.pollInterval=d(y(i,/[,\[\s]/));y(i,Je);var l=et(e,i,"event");if(l){s.eventFilter=l}n.push(s)}else if(o.indexOf("sse:")===0){n.push({trigger:"sse",sseEvent:o.substr(4)})}else{var u={trigger:o};var l=et(e,i,"event");if(l){u.eventFilter=l}while(i.length>0&&i[0]!==","){y(i,Je);var f=i.shift();if(f==="changed"){u.changed=true}else if(f==="once"){u.once=true}else if(f==="consume"){u.consume=true}else if(f==="delay"&&i[0]===":"){i.shift();u.delay=d(y(i,x))}else if(f==="from"&&i[0]===":"){i.shift();if(Ze.test(i[0])){var c=tt(i)}else{var c=y(i,x);if(c==="closest"||c==="find"||c==="next"||c==="previous"){i.shift();var h=tt(i);if(h.length>0){c+=" "+h}}}u.from=c}else if(f==="target"&&i[0]===":"){i.shift();u.target=tt(i)}else if(f==="throttle"&&i[0]===":"){i.shift();u.throttle=d(y(i,x))}else if(f==="queue"&&i[0]===":"){i.shift();u.queue=y(i,x)}else if(f==="root"&&i[0]===":"){i.shift();u[f]=tt(i)}else if(f==="threshold"&&i[0]===":"){i.shift();u[f]=y(i,x)}else{fe(e,"htmx:syntax:error",{token:i.shift()})}}n.push(u)}}if(i.length===a){fe(e,"htmx:syntax:error",{token:i.shift()})}y(i,Je)}while(i[0]===","&&i.shift());if(r){r[t]=n}return n}function it(e){var t=te(e,"hx-trigger");var r=[];if(t){var n=Q.config.triggerSpecsCache;r=n&&n[t]||nt(e,t,n)}if(r.length>0){return r}else if(h(e,"form")){return[{trigger:"submit"}]}else if(h(e,'input[type="button"], input[type="submit"]')){return[{trigger:"click"}]}else if(h(e,rt)){return[{trigger:"change"}]}else{return[{trigger:"click"}]}}function at(e){ae(e).cancelled=true}function ot(e,t,r){var n=ae(e);n.timeout=setTimeout(function(){if(se(e)&&n.cancelled!==true){if(!ct(r,e,Wt("hx:poll:trigger",{triggerSpec:r,target:e}))){t(e)}ot(e,t,r)}},r.pollInterval)}function st(e){return location.hostname===e.hostname&&ee(e,"href")&&ee(e,"href").indexOf("#")!==0}function lt(t,r,e){if(t.tagName==="A"&&st(t)&&(t.target===""||t.target==="_self")||t.tagName==="FORM"){r.boosted=true;var n,i;if(t.tagName==="A"){n="get";i=ee(t,"href")}else{var a=ee(t,"method");n=a?a.toLowerCase():"get";if(n==="get"){}i=ee(t,"action")}e.forEach(function(e){ht(t,function(e,t){if(v(e,Q.config.disableSelector)){p(e);return}he(n,i,e,t)},r,e,true)})}}function ut(e,t){if(e.type==="submit"||e.type==="click"){if(t.tagName==="FORM"){return true}if(h(t,'input[type="submit"], button')&&v(t,"form")!==null){return true}if(t.tagName==="A"&&t.href&&(t.getAttribute("href")==="#"||t.getAttribute("href").indexOf("#")!==0)){return true}}return false}function ft(e,t){return ae(e).boosted&&e.tagName==="A"&&t.type==="click"&&(t.ctrlKey||t.metaKey)}function ct(e,t,r){var n=e.eventFilter;if(n){try{return n.call(t,r)!==true}catch(e){fe(re().body,"htmx:eventFilter:error",{error:e,source:n.source});return true}}return false}function ht(a,o,e,s,l){var u=ae(a);var t;if(s.from){t=Z(a,s.from)}else{t=[a]}if(s.changed){t.forEach(function(e){var t=ae(e);t.lastValue=e.value})}oe(t,function(n){var i=function(e){if(!se(a)){n.removeEventListener(s.trigger,i);return}if(ft(a,e)){return}if(l||ut(e,a)){e.preventDefault()}if(ct(s,a,e)){return}var t=ae(e);t.triggerSpec=s;if(t.handledFor==null){t.handledFor=[]}if(t.handledFor.indexOf(a)<0){t.handledFor.push(a);if(s.consume){e.stopPropagation()}if(s.target&&e.target){if(!h(e.target,s.target)){return}}if(s.once){if(u.triggeredOnce){return}else{u.triggeredOnce=true}}if(s.changed){var r=ae(n);if(r.lastValue===n.value){return}r.lastValue=n.value}if(u.delayed){clearTimeout(u.delayed)}if(u.throttle){return}if(s.throttle>0){if(!u.throttle){o(a,e);u.throttle=setTimeout(function(){u.throttle=null},s.throttle)}}else if(s.delay>0){u.delayed=setTimeout(function(){o(a,e)},s.delay)}else{ce(a,"htmx:trigger");o(a,e)}}};if(e.listenerInfos==null){e.listenerInfos=[]}e.listenerInfos.push({trigger:s.trigger,listener:i,on:n});n.addEventListener(s.trigger,i)})}var vt=false;var dt=null;function gt(){if(!dt){dt=function(){vt=true};window.addEventListener("scroll",dt);setInterval(function(){if(vt){vt=false;oe(re().querySelectorAll("[hx-trigger='revealed'],[data-hx-trigger='revealed']"),function(e){mt(e)})}},200)}}function mt(t){if(!o(t,"data-hx-revealed")&&X(t)){t.setAttribute("data-hx-revealed","true");var e=ae(t);if(e.initHash){ce(t,"revealed")}else{t.addEventListener("htmx:afterProcessNode",function(e){ce(t,"revealed")},{once:true})}}}function pt(e,t,r){var n=D(r);for(var i=0;i<n.length;i++){var a=n[i].split(/:(.+)/);if(a[0]==="connect"){xt(e,a[1],0)}if(a[0]==="send"){bt(e)}}}function xt(s,r,n){if(!se(s)){return}if(r.indexOf("/")==0){var e=location.hostname+(location.port?":"+location.port:"");if(location.protocol=="https:"){r="wss://"+e+r}else if(location.protocol=="http:"){r="ws://"+e+r}}var t=Q.createWebSocket(r);t.onerror=function(e){fe(s,"htmx:wsError",{error:e,socket:t});yt(s)};t.onclose=function(e){if([1006,1012,1013].indexOf(e.code)>=0){var t=wt(n);setTimeout(function(){xt(s,r,n+1)},t)}};t.onopen=function(e){n=0};ae(s).webSocket=t;t.addEventListener("message",function(e){if(yt(s)){return}var t=e.data;R(s,function(e){t=e.transformResponse(t,null,s)});var r=T(s);var n=l(t);var i=M(n.children);for(var a=0;a<i.length;a++){var o=i[a];Ee(te(o,"hx-swap-oob")||"true",o,r)}nr(r.tasks)})}function yt(e){if(!se(e)){ae(e).webSocket.close();return true}}function bt(u){var f=c(u,function(e){return ae(e).webSocket!=null});if(f){u.addEventListener(it(u)[0].trigger,function(e){var t=ae(f).webSocket;var r=xr(u,f);var n=dr(u,"post");var i=n.errors;var a=n.values;var o=Hr(u);var s=le(a,o);var l=yr(s,u);l["HEADERS"]=r;if(i&&i.length>0){ce(u,"htmx:validation:halted",i);return}t.send(JSON.stringify(l));if(ut(e,u)){e.preventDefault()}})}else{fe(u,"htmx:noWebSocketSourceError")}}function wt(e){var t=Q.config.wsReconnectDelay;if(typeof t==="function"){return t(e)}if(t==="full-jitter"){var r=Math.min(e,6);var n=1e3*Math.pow(2,r);return n*Math.random()}b('htmx.config.wsReconnectDelay must either be a function or the string "full-jitter"')}function St(e,t,r){var n=D(r);for(var i=0;i<n.length;i++){var a=n[i].split(/:(.+)/);if(a[0]==="connect"){Et(e,a[1])}if(a[0]==="swap"){Ct(e,a[1])}}}function Et(t,e){var r=Q.createEventSource(e);r.onerror=function(e){fe(t,"htmx:sseError",{error:e,source:r});Tt(t)};ae(t).sseEventSource=r}function Ct(a,o){var s=c(a,Ot);if(s){var l=ae(s).sseEventSource;var u=function(e){if(Tt(s)){return}if(!se(a)){l.removeEventListener(o,u);return}var t=e.data;R(a,function(e){t=e.transformResponse(t,null,a)});var r=wr(a);var n=ye(a);var i=T(a);je(r.swapStyle,n,a,t,i);nr(i.tasks);ce(a,"htmx:sseMessage",e)};ae(a).sseListener=u;l.addEventListener(o,u)}else{fe(a,"htmx:noSSESourceError")}}function Rt(e,t,r){var n=c(e,Ot);if(n){var i=ae(n).sseEventSource;var a=function(){if(!Tt(n)){if(se(e)){t(e)}else{i.removeEventListener(r,a)}}};ae(e).sseListener=a;i.addEventListener(r,a)}else{fe(e,"htmx:noSSESourceError")}}function Tt(e){if(!se(e)){ae(e).sseEventSource.close();return true}}function Ot(e){return ae(e).sseEventSource!=null}function qt(e,t,r,n){var i=function(){if(!r.loaded){r.loaded=true;t(e)}};if(n>0){setTimeout(i,n)}else{i()}}function Ht(t,i,e){var a=false;oe(w,function(r){if(o(t,"hx-"+r)){var n=te(t,"hx-"+r);a=true;i.path=n;i.verb=r;e.forEach(function(e){Lt(t,e,i,function(e,t){if(v(e,Q.config.disableSelector)){p(e);return}he(r,n,e,t)})})}});return a}function Lt(n,e,t,r){if(e.sseEvent){Rt(n,r,e.sseEvent)}else if(e.trigger==="revealed"){gt();ht(n,r,t,e);mt(n)}else if(e.trigger==="intersect"){var i={};if(e.root){i.root=ue(n,e.root)}if(e.threshold){i.threshold=parseFloat(e.threshold)}var a=new IntersectionObserver(function(e){for(var t=0;t<e.length;t++){var r=e[t];if(r.isIntersecting){ce(n,"intersect");break}}},i);a.observe(n);ht(n,r,t,e)}else if(e.trigger==="load"){if(!ct(e,n,Wt("load",{elt:n}))){qt(n,r,t,e.delay)}}else if(e.pollInterval>0){t.polling=true;ot(n,r,e)}else{ht(n,r,t,e)}}function At(e){if(Q.config.allowScriptTags&&(e.type==="text/javascript"||e.type==="module"||e.type==="")){var t=re().createElement("script");oe(e.attributes,function(e){t.setAttribute(e.name,e.value)});t.textContent=e.textContent;t.async=false;if(Q.config.inlineScriptNonce){t.nonce=Q.config.inlineScriptNonce}var r=e.parentElement;try{r.insertBefore(t,e)}catch(e){b(e)}finally{if(e.parentElement){e.parentElement.removeChild(e)}}}}function Nt(e){if(h(e,"script")){At(e)}oe(f(e,"script"),function(e){At(e)})}function It(e){var t=e.attributes;for(var r=0;r<t.length;r++){var n=t[r].name;if(s(n,"hx-on:")||s(n,"data-hx-on:")||s(n,"hx-on-")||s(n,"data-hx-on-")){return true}}return false}function kt(e){var t=null;var r=[];if(It(e)){r.push(e)}if(document.evaluate){var n=document.evaluate('.//*[@*[ starts-with(name(), "hx-on:") or starts-with(name(), "data-hx-on:") or'+' starts-with(name(), "hx-on-") or starts-with(name(), "data-hx-on-") ]]',e);while(t=n.iterateNext())r.push(t)}else{var i=e.getElementsByTagName("*");for(var a=0;a<i.length;a++){if(It(i[a])){r.push(i[a])}}}return r}function Pt(e){if(e.querySelectorAll){var t=", [hx-boost] a, [data-hx-boost] a, a[hx-boost], a[data-hx-boost]";var r=e.querySelectorAll(i+t+", form, [type='submit'], [hx-sse], [data-hx-sse], [hx-ws],"+" [data-hx-ws], [hx-ext], [data-hx-ext], [hx-trigger], [data-hx-trigger], [hx-on], [data-hx-on]");return r}else{return[]}}function Mt(e){var t=v(e.target,"button, input[type='submit']");var r=Dt(e);if(r){r.lastButtonClicked=t}}function Xt(e){var t=Dt(e);if(t){t.lastButtonClicked=null}}function Dt(e){var t=v(e.target,"button, input[type='submit']");if(!t){return}var r=g("#"+ee(t,"form"))||v(t,"form");if(!r){return}return ae(r)}function Ut(e){e.addEventListener("click",Mt);e.addEventListener("focusin",Mt);e.addEventListener("focusout",Xt)}function Bt(e){var t=Ye(e);var r=0;for(var n=0;n<t.length;n++){const i=t[n];if(i==="{"){r++}else if(i==="}"){r--}}return r}function Ft(t,e,r){var n=ae(t);if(!Array.isArray(n.onHandlers)){n.onHandlers=[]}var i;var a=function(e){return Tr(t,function(){if(!i){i=new Function("event",r)}i.call(t,e)})};t.addEventListener(e,a);n.onHandlers.push({event:e,listener:a})}function Vt(e){var t=te(e,"hx-on");if(t){var r={};var n=t.split("\n");var i=null;var a=0;while(n.length>0){var o=n.shift();var s=o.match(/^\s*([a-zA-Z:\-\.]+:)(.*)/);if(a===0&&s){o.split(":");i=s[1].slice(0,-1);r[i]=s[2]}else{r[i]+=o}a+=Bt(o)}for(var l in r){Ft(e,l,r[l])}}}function jt(e){Ae(e);for(var t=0;t<e.attributes.length;t++){var r=e.attributes[t].name;var n=e.attributes[t].value;if(s(r,"hx-on")||s(r,"data-hx-on")){var i=r.indexOf("-on")+3;var a=r.slice(i,i+1);if(a==="-"||a===":"){var o=r.slice(i+1);if(s(o,":")){o="htmx"+o}else if(s(o,"-")){o="htmx:"+o.slice(1)}else if(s(o,"htmx-")){o="htmx:"+o.slice(5)}Ft(e,o,n)}}}}function _t(t){if(v(t,Q.config.disableSelector)){p(t);return}var r=ae(t);if(r.initHash!==Le(t)){Ne(t);r.initHash=Le(t);Vt(t);ce(t,"htmx:beforeProcessNode");if(t.value){r.lastValue=t.value}var e=it(t);var n=Ht(t,r,e);if(!n){if(ne(t,"hx-boost")==="true"){lt(t,r,e)}else if(o(t,"hx-trigger")){e.forEach(function(e){Lt(t,e,r,function(){})})}}if(t.tagName==="FORM"||ee(t,"type")==="submit"&&o(t,"form")){Ut(t)}var i=te(t,"hx-sse");if(i){St(t,r,i)}var a=te(t,"hx-ws");if(a){pt(t,r,a)}ce(t,"htmx:afterProcessNode")}}function zt(e){e=g(e);if(v(e,Q.config.disableSelector)){p(e);return}_t(e);oe(Pt(e),function(e){_t(e)});oe(kt(e),jt)}function $t(e){return e.replace(/([a-z0-9])([A-Z])/g,"$1-$2").toLowerCase()}function Wt(e,t){var r;if(window.CustomEvent&&typeof window.CustomEvent==="function"){r=new CustomEvent(e,{bubbles:true,cancelable:true,detail:t})}else{r=re().createEvent("CustomEvent");r.initCustomEvent(e,true,true,t)}return r}function fe(e,t,r){ce(e,t,le({error:t},r))}function Gt(e){return e==="htmx:afterProcessNode"}function R(e,t){oe(Fr(e),function(e){try{t(e)}catch(e){b(e)}})}function b(e){if(console.error){console.error(e)}else if(console.log){console.log("ERROR: ",e)}}function ce(e,t,r){e=g(e);if(r==null){r={}}r["elt"]=e;var n=Wt(t,r);if(Q.logger&&!Gt(t)){Q.logger(e,t,r)}if(r.error){b(r.error);ce(e,"htmx:error",{errorInfo:r})}var i=e.dispatchEvent(n);var a=$t(t);if(i&&a!==t){var o=Wt(a,n.detail);i=i&&e.dispatchEvent(o)}R(e,function(e){i=i&&(e.onEvent(t,n)!==false&&!n.defaultPrevented)});return i}var Jt=location.pathname+location.search;function Zt(){var e=re().querySelector("[hx-history-elt],[data-hx-history-elt]");return e||re().body}function Kt(e,t,r,n){if(!U()){return}if(Q.config.historyCacheSize<=0){localStorage.removeItem("htmx-history-cache");return}e=B(e);var i=E(localStorage.getItem("htmx-history-cache"))||[];for(var a=0;a<i.length;a++){if(i[a].url===e){i.splice(a,1);break}}var o={url:e,content:t,title:r,scroll:n};ce(re().body,"htmx:historyItemCreated",{item:o,cache:i});i.push(o);while(i.length>Q.config.historyCacheSize){i.shift()}while(i.length>0){try{localStorage.setItem("htmx-history-cache",JSON.stringify(i));break}catch(e){fe(re().body,"htmx:historyCacheError",{cause:e,cache:i});i.shift()}}}function Yt(e){if(!U()){return null}e=B(e);var t=E(localStorage.getItem("htmx-history-cache"))||[];for(var r=0;r<t.length;r++){if(t[r].url===e){return t[r]}}return null}function Qt(e){var t=Q.config.requestClass;var r=e.cloneNode(true);oe(f(r,"."+t),function(e){n(e,t)});return r.innerHTML}function er(){var e=Zt();var t=Jt||location.pathname+location.search;var r;try{r=re().querySelector('[hx-history="false" i],[data-hx-history="false" i]')}catch(e){r=re().querySelector('[hx-history="false"],[data-hx-history="false"]')}if(!r){ce(re().body,"htmx:beforeHistorySave",{path:t,historyElt:e});Kt(t,Qt(e),re().title,window.scrollY)}if(Q.config.historyEnabled)history.replaceState({htmx:true},re().title,window.location.href)}function tr(e){if(Q.config.getCacheBusterParam){e=e.replace(/org\.htmx\.cache-buster=[^&]*&?/,"");if(G(e,"&")||G(e,"?")){e=e.slice(0,-1)}}if(Q.config.historyEnabled){history.pushState({htmx:true},"",e)}Jt=e}function rr(e){if(Q.config.historyEnabled)history.replaceState({htmx:true},"",e);Jt=e}function nr(e){oe(e,function(e){e.call()})}function ir(a){var e=new XMLHttpRequest;var o={path:a,xhr:e};ce(re().body,"htmx:historyCacheMiss",o);e.open("GET",a,true);e.setRequestHeader("HX-Request","true");e.setRequestHeader("HX-History-Restore-Request","true");e.setRequestHeader("HX-Current-URL",re().location.href);e.onload=function(){if(this.status>=200&&this.status<400){ce(re().body,"htmx:historyCacheMissLoad",o);var e=l(this.response);e=e.querySelector("[hx-history-elt],[data-hx-history-elt]")||e;var t=Zt();var r=T(t);var n=Ve(this.response);if(n){var i=C("title");if(i){i.innerHTML=n}else{window.document.title=n}}Ue(t,e,r);nr(r.tasks);Jt=a;ce(re().body,"htmx:historyRestore",{path:a,cacheMiss:true,serverResponse:this.response})}else{fe(re().body,"htmx:historyCacheMissLoadError",o)}};e.send()}function ar(e){er();e=e||location.pathname+location.search;var t=Yt(e);if(t){var r=l(t.content);var n=Zt();var i=T(n);Ue(n,r,i);nr(i.tasks);document.title=t.title;setTimeout(function(){window.scrollTo(0,t.scroll)},0);Jt=e;ce(re().body,"htmx:historyRestore",{path:e,item:t})}else{if(Q.config.refreshOnHistoryMiss){window.location.reload(true)}else{ir(e)}}}function or(e){var t=pe(e,"hx-indicator");if(t==null){t=[e]}oe(t,function(e){var t=ae(e);t.requestCount=(t.requestCount||0)+1;e.classList["add"].call(e.classList,Q.config.requestClass)});return t}function sr(e){var t=pe(e,"hx-disabled-elt");if(t==null){t=[]}oe(t,function(e){var t=ae(e);t.requestCount=(t.requestCount||0)+1;e.setAttribute("disabled","")});return t}function lr(e,t){oe(e,function(e){var t=ae(e);t.requestCount=(t.requestCount||0)-1;if(t.requestCount===0){e.classList["remove"].call(e.classList,Q.config.requestClass)}});oe(t,function(e){var t=ae(e);t.requestCount=(t.requestCount||0)-1;if(t.requestCount===0){e.removeAttribute("disabled")}})}function ur(e,t){for(var r=0;r<e.length;r++){var n=e[r];if(n.isSameNode(t)){return true}}return false}function fr(e){if(e.name===""||e.name==null||e.disabled||v(e,"fieldset[disabled]")){return false}if(e.type==="button"||e.type==="submit"||e.tagName==="image"||e.tagName==="reset"||e.tagName==="file"){return false}if(e.type==="checkbox"||e.type==="radio"){return e.checked}return true}function cr(e,t,r){if(e!=null&&t!=null){var n=r[e];if(n===undefined){r[e]=t}else if(Array.isArray(n)){if(Array.isArray(t)){r[e]=n.concat(t)}else{n.push(t)}}else{if(Array.isArray(t)){r[e]=[n].concat(t)}else{r[e]=[n,t]}}}}function hr(t,r,n,e,i){if(e==null||ur(t,e)){return}else{t.push(e)}if(fr(e)){var a=ee(e,"name");var o=e.value;if(e.multiple&&e.tagName==="SELECT"){o=M(e.querySelectorAll("option:checked")).map(function(e){return e.value})}if(e.files){o=M(e.files)}cr(a,o,r);if(i){vr(e,n)}}if(h(e,"form")){var s=e.elements;oe(s,function(e){hr(t,r,n,e,i)})}}function vr(e,t){if(e.willValidate){ce(e,"htmx:validation:validate");if(!e.checkValidity()){t.push({elt:e,message:e.validationMessage,validity:e.validity});ce(e,"htmx:validation:failed",{message:e.validationMessage,validity:e.validity})}}}function dr(e,t){var r=[];var n={};var i={};var a=[];var o=ae(e);if(o.lastButtonClicked&&!se(o.lastButtonClicked)){o.lastButtonClicked=null}var s=h(e,"form")&&e.noValidate!==true||te(e,"hx-validate")==="true";if(o.lastButtonClicked){s=s&&o.lastButtonClicked.formNoValidate!==true}if(t!=="get"){hr(r,i,a,v(e,"form"),s)}hr(r,n,a,e,s);if(o.lastButtonClicked||e.tagName==="BUTTON"||e.tagName==="INPUT"&&ee(e,"type")==="submit"){var l=o.lastButtonClicked||e;var u=ee(l,"name");cr(u,l.value,i)}var f=pe(e,"hx-include");oe(f,function(e){hr(r,n,a,e,s);if(!h(e,"form")){oe(e.querySelectorAll(rt),function(e){hr(r,n,a,e,s)})}});n=le(n,i);return{errors:a,values:n}}function gr(e,t,r){if(e!==""){e+="&"}if(String(r)==="[object Object]"){r=JSON.stringify(r)}var n=encodeURIComponent(r);e+=encodeURIComponent(t)+"="+n;return e}function mr(e){var t="";for(var r in e){if(e.hasOwnProperty(r)){var n=e[r];if(Array.isArray(n)){oe(n,function(e){t=gr(t,r,e)})}else{t=gr(t,r,n)}}}return t}function pr(e){var t=new FormData;for(var r in e){if(e.hasOwnProperty(r)){var n=e[r];if(Array.isArray(n)){oe(n,function(e){t.append(r,e)})}else{t.append(r,n)}}}return t}function xr(e,t,r){var n={"HX-Request":"true","HX-Trigger":ee(e,"id"),"HX-Trigger-Name":ee(e,"name"),"HX-Target":te(t,"id"),"HX-Current-URL":re().location.href};Rr(e,"hx-headers",false,n);if(r!==undefined){n["HX-Prompt"]=r}if(ae(e).boosted){n["HX-Boosted"]="true"}return n}function yr(t,e){var r=ne(e,"hx-params");if(r){if(r==="none"){return{}}else if(r==="*"){return t}else if(r.indexOf("not ")===0){oe(r.substr(4).split(","),function(e){e=e.trim();delete t[e]});return t}else{var n={};oe(r.split(","),function(e){e=e.trim();n[e]=t[e]});return n}}else{return t}}function br(e){return ee(e,"href")&&ee(e,"href").indexOf("#")>=0}function wr(e,t){var r=t?t:ne(e,"hx-swap");var n={swapStyle:ae(e).boosted?"innerHTML":Q.config.defaultSwapStyle,swapDelay:Q.config.defaultSwapDelay,settleDelay:Q.config.defaultSettleDelay};if(Q.config.scrollIntoViewOnBoost&&ae(e).boosted&&!br(e)){n["show"]="top"}if(r){var i=D(r);if(i.length>0){for(var a=0;a<i.length;a++){var o=i[a];if(o.indexOf("swap:")===0){n["swapDelay"]=d(o.substr(5))}else if(o.indexOf("settle:")===0){n["settleDelay"]=d(o.substr(7))}else if(o.indexOf("transition:")===0){n["transition"]=o.substr(11)==="true"}else if(o.indexOf("ignoreTitle:")===0){n["ignoreTitle"]=o.substr(12)==="true"}else if(o.indexOf("scroll:")===0){var s=o.substr(7);var l=s.split(":");var u=l.pop();var f=l.length>0?l.join(":"):null;n["scroll"]=u;n["scrollTarget"]=f}else if(o.indexOf("show:")===0){var c=o.substr(5);var l=c.split(":");var h=l.pop();var f=l.length>0?l.join(":"):null;n["show"]=h;n["showTarget"]=f}else if(o.indexOf("focus-scroll:")===0){var v=o.substr("focus-scroll:".length);n["focusScroll"]=v=="true"}else if(a==0){n["swapStyle"]=o}else{b("Unknown modifier in hx-swap: "+o)}}}}return n}function Sr(e){return ne(e,"hx-encoding")==="multipart/form-data"||h(e,"form")&&ee(e,"enctype")==="multipart/form-data"}function Er(t,r,n){var i=null;R(r,function(e){if(i==null){i=e.encodeParameters(t,n,r)}});if(i!=null){return i}else{if(Sr(r)){return pr(n)}else{return mr(n)}}}function T(e){return{tasks:[],elts:[e]}}function Cr(e,t){var r=e[0];var n=e[e.length-1];if(t.scroll){var i=null;if(t.scrollTarget){i=ue(r,t.scrollTarget)}if(t.scroll==="top"&&(r||i)){i=i||r;i.scrollTop=0}if(t.scroll==="bottom"&&(n||i)){i=i||n;i.scrollTop=i.scrollHeight}}if(t.show){var i=null;if(t.showTarget){var a=t.showTarget;if(t.showTarget==="window"){a="body"}i=ue(r,a)}if(t.show==="top"&&(r||i)){i=i||r;i.scrollIntoView({block:"start",behavior:Q.config.scrollBehavior})}if(t.show==="bottom"&&(n||i)){i=i||n;i.scrollIntoView({block:"end",behavior:Q.config.scrollBehavior})}}}function Rr(e,t,r,n){if(n==null){n={}}if(e==null){return n}var i=te(e,t);if(i){var a=i.trim();var o=r;if(a==="unset"){return null}if(a.indexOf("javascript:")===0){a=a.substr(11);o=true}else if(a.indexOf("js:")===0){a=a.substr(3);o=true}if(a.indexOf("{")!==0){a="{"+a+"}"}var s;if(o){s=Tr(e,function(){return Function("return ("+a+")")()},{})}else{s=E(a)}for(var l in s){if(s.hasOwnProperty(l)){if(n[l]==null){n[l]=s[l]}}}}return Rr(u(e),t,r,n)}function Tr(e,t,r){if(Q.config.allowEval){return t()}else{fe(e,"htmx:evalDisallowedError");return r}}function Or(e,t){return Rr(e,"hx-vars",true,t)}function qr(e,t){return Rr(e,"hx-vals",false,t)}function Hr(e){return le(Or(e),qr(e))}function Lr(t,r,n){if(n!==null){try{t.setRequestHeader(r,n)}catch(e){t.setRequestHeader(r,encodeURIComponent(n));t.setRequestHeader(r+"-URI-AutoEncoded","true")}}}function Ar(t){if(t.responseURL&&typeof URL!=="undefined"){try{var e=new URL(t.responseURL);return e.pathname+e.search}catch(e){fe(re().body,"htmx:badResponseUrl",{url:t.responseURL})}}}function O(e,t){return t.test(e.getAllResponseHeaders())}function Nr(e,t,r){e=e.toLowerCase();if(r){if(r instanceof Element||I(r,"String")){return he(e,t,null,null,{targetOverride:g(r),returnPromise:true})}else{return he(e,t,g(r.source),r.event,{handler:r.handler,headers:r.headers,values:r.values,targetOverride:g(r.target),swapOverride:r.swap,select:r.select,returnPromise:true})}}else{return he(e,t,null,null,{returnPromise:true})}}function Ir(e){var t=[];while(e){t.push(e);e=e.parentElement}return t}function kr(e,t,r){var n;var i;if(typeof URL==="function"){i=new URL(t,document.location.href);var a=document.location.origin;n=a===i.origin}else{i=t;n=s(t,document.location.origin)}if(Q.config.selfRequestsOnly){if(!n){return false}}return ce(e,"htmx:validateUrl",le({url:i,sameHost:n},r))}function he(t,r,n,i,a,e){var o=null;var s=null;a=a!=null?a:{};if(a.returnPromise&&typeof Promise!=="undefined"){var l=new Promise(function(e,t){o=e;s=t})}if(n==null){n=re().body}var M=a.handler||Mr;var X=a.select||null;if(!se(n)){ie(o);return l}var u=a.targetOverride||ye(n);if(u==null||u==me){fe(n,"htmx:targetError",{target:te(n,"hx-target")});ie(s);return l}var f=ae(n);var c=f.lastButtonClicked;if(c){var h=ee(c,"formaction");if(h!=null){r=h}var v=ee(c,"formmethod");if(v!=null){if(v.toLowerCase()!=="dialog"){t=v}}}var d=ne(n,"hx-confirm");if(e===undefined){var D=function(e){return he(t,r,n,i,a,!!e)};var U={target:u,elt:n,path:r,verb:t,triggeringEvent:i,etc:a,issueRequest:D,question:d};if(ce(n,"htmx:confirm",U)===false){ie(o);return l}}var g=n;var m=ne(n,"hx-sync");var p=null;var x=false;if(m){var B=m.split(":");var F=B[0].trim();if(F==="this"){g=xe(n,"hx-sync")}else{g=ue(n,F)}m=(B[1]||"drop").trim();f=ae(g);if(m==="drop"&&f.xhr&&f.abortable!==true){ie(o);return l}else if(m==="abort"){if(f.xhr){ie(o);return l}else{x=true}}else if(m==="replace"){ce(g,"htmx:abort")}else if(m.indexOf("queue")===0){var V=m.split(" ");p=(V[1]||"last").trim()}}if(f.xhr){if(f.abortable){ce(g,"htmx:abort")}else{if(p==null){if(i){var y=ae(i);if(y&&y.triggerSpec&&y.triggerSpec.queue){p=y.triggerSpec.queue}}if(p==null){p="last"}}if(f.queuedRequests==null){f.queuedRequests=[]}if(p==="first"&&f.queuedRequests.length===0){f.queuedRequests.push(function(){he(t,r,n,i,a)})}else if(p==="all"){f.queuedRequests.push(function(){he(t,r,n,i,a)})}else if(p==="last"){f.queuedRequests=[];f.queuedRequests.push(function(){he(t,r,n,i,a)})}ie(o);return l}}var b=new XMLHttpRequest;f.xhr=b;f.abortable=x;var w=function(){f.xhr=null;f.abortable=false;if(f.queuedRequests!=null&&f.queuedRequests.length>0){var e=f.queuedRequests.shift();e()}};var j=ne(n,"hx-prompt");if(j){var S=prompt(j);if(S===null||!ce(n,"htmx:prompt",{prompt:S,target:u})){ie(o);w();return l}}if(d&&!e){if(!confirm(d)){ie(o);w();return l}}var E=xr(n,u,S);if(t!=="get"&&!Sr(n)){E["Content-Type"]="application/x-www-form-urlencoded"}if(a.headers){E=le(E,a.headers)}var _=dr(n,t);var C=_.errors;var R=_.values;if(a.values){R=le(R,a.values)}var z=Hr(n);var $=le(R,z);var T=yr($,n);if(Q.config.getCacheBusterParam&&t==="get"){T["org.htmx.cache-buster"]=ee(u,"id")||"true"}if(r==null||r===""){r=re().location.href}var O=Rr(n,"hx-request");var W=ae(n).boosted;var q=Q.config.methodsThatUseUrlParams.indexOf(t)>=0;var H={boosted:W,useUrlParams:q,parameters:T,unfilteredParameters:$,headers:E,target:u,verb:t,errors:C,withCredentials:a.credentials||O.credentials||Q.config.withCredentials,timeout:a.timeout||O.timeout||Q.config.timeout,path:r,triggeringEvent:i};if(!ce(n,"htmx:configRequest",H)){ie(o);w();return l}r=H.path;t=H.verb;E=H.headers;T=H.parameters;C=H.errors;q=H.useUrlParams;if(C&&C.length>0){ce(n,"htmx:validation:halted",H);ie(o);w();return l}var G=r.split("#");var J=G[0];var L=G[1];var A=r;if(q){A=J;var Z=Object.keys(T).length!==0;if(Z){if(A.indexOf("?")<0){A+="?"}else{A+="&"}A+=mr(T);if(L){A+="#"+L}}}if(!kr(n,A,H)){fe(n,"htmx:invalidPath",H);ie(s);return l}b.open(t.toUpperCase(),A,true);b.overrideMimeType("text/html");b.withCredentials=H.withCredentials;b.timeout=H.timeout;if(O.noHeaders){}else{for(var N in E){if(E.hasOwnProperty(N)){var K=E[N];Lr(b,N,K)}}}var I={xhr:b,target:u,requestConfig:H,etc:a,boosted:W,select:X,pathInfo:{requestPath:r,finalRequestPath:A,anchor:L}};b.onload=function(){try{var e=Ir(n);I.pathInfo.responsePath=Ar(b);M(n,I);lr(k,P);ce(n,"htmx:afterRequest",I);ce(n,"htmx:afterOnLoad",I);if(!se(n)){var t=null;while(e.length>0&&t==null){var r=e.shift();if(se(r)){t=r}}if(t){ce(t,"htmx:afterRequest",I);ce(t,"htmx:afterOnLoad",I)}}ie(o);w()}catch(e){fe(n,"htmx:onLoadError",le({error:e},I));throw e}};b.onerror=function(){lr(k,P);fe(n,"htmx:afterRequest",I);fe(n,"htmx:sendError",I);ie(s);w()};b.onabort=function(){lr(k,P);fe(n,"htmx:afterRequest",I);fe(n,"htmx:sendAbort",I);ie(s);w()};b.ontimeout=function(){lr(k,P);fe(n,"htmx:afterRequest",I);fe(n,"htmx:timeout",I);ie(s);w()};if(!ce(n,"htmx:beforeRequest",I)){ie(o);w();return l}var k=or(n);var P=sr(n);oe(["loadstart","loadend","progress","abort"],function(t){oe([b,b.upload],function(e){e.addEventListener(t,function(e){ce(n,"htmx:xhr:"+t,{lengthComputable:e.lengthComputable,loaded:e.loaded,total:e.total})})})});ce(n,"htmx:beforeSend",I);var Y=q?null:Er(b,n,T);b.send(Y);return l}function Pr(e,t){var r=t.xhr;var n=null;var i=null;if(O(r,/HX-Push:/i)){n=r.getResponseHeader("HX-Push");i="push"}else if(O(r,/HX-Push-Url:/i)){n=r.getResponseHeader("HX-Push-Url");i="push"}else if(O(r,/HX-Replace-Url:/i)){n=r.getResponseHeader("HX-Replace-Url");i="replace"}if(n){if(n==="false"){return{}}else{return{type:i,path:n}}}var a=t.pathInfo.finalRequestPath;var o=t.pathInfo.responsePath;var s=ne(e,"hx-push-url");var l=ne(e,"hx-replace-url");var u=ae(e).boosted;var f=null;var c=null;if(s){f="push";c=s}else if(l){f="replace";c=l}else if(u){f="push";c=o||a}if(c){if(c==="false"){return{}}if(c==="true"){c=o||a}if(t.pathInfo.anchor&&c.indexOf("#")===-1){c=c+"#"+t.pathInfo.anchor}return{type:f,path:c}}else{return{}}}function Mr(l,u){var f=u.xhr;var c=u.target;var e=u.etc;var t=u.requestConfig;var h=u.select;if(!ce(l,"htmx:beforeOnLoad",u))return;if(O(f,/HX-Trigger:/i)){_e(f,"HX-Trigger",l)}if(O(f,/HX-Location:/i)){er();var r=f.getResponseHeader("HX-Location");var v;if(r.indexOf("{")===0){v=E(r);r=v["path"];delete v["path"]}Nr("GET",r,v).then(function(){tr(r)});return}var n=O(f,/HX-Refresh:/i)&&"true"===f.getResponseHeader("HX-Refresh");if(O(f,/HX-Redirect:/i)){location.href=f.getResponseHeader("HX-Redirect");n&&location.reload();return}if(n){location.reload();return}if(O(f,/HX-Retarget:/i)){if(f.getResponseHeader("HX-Retarget")==="this"){u.target=l}else{u.target=ue(l,f.getResponseHeader("HX-Retarget"))}}var d=Pr(l,u);var i=f.status>=200&&f.status<400&&f.status!==204;var g=f.response;var a=f.status>=400;var m=Q.config.ignoreTitle;var o=le({shouldSwap:i,serverResponse:g,isError:a,ignoreTitle:m},u);if(!ce(c,"htmx:beforeSwap",o))return;c=o.target;g=o.serverResponse;a=o.isError;m=o.ignoreTitle;u.target=c;u.failed=a;u.successful=!a;if(o.shouldSwap){if(f.status===286){at(l)}R(l,function(e){g=e.transformResponse(g,f,l)});if(d.type){er()}var s=e.swapOverride;if(O(f,/HX-Reswap:/i)){s=f.getResponseHeader("HX-Reswap")}var v=wr(l,s);if(v.hasOwnProperty("ignoreTitle")){m=v.ignoreTitle}c.classList.add(Q.config.swappingClass);var p=null;var x=null;var y=function(){try{var e=document.activeElement;var t={};try{t={elt:e,start:e?e.selectionStart:null,end:e?e.selectionEnd:null}}catch(e){}var r;if(h){r=h}if(O(f,/HX-Reselect:/i)){r=f.getResponseHeader("HX-Reselect")}if(d.type){ce(re().body,"htmx:beforeHistoryUpdate",le({history:d},u));if(d.type==="push"){tr(d.path);ce(re().body,"htmx:pushedIntoHistory",{path:d.path})}else{rr(d.path);ce(re().body,"htmx:replacedInHistory",{path:d.path})}}var n=T(c);je(v.swapStyle,c,l,g,n,r);if(t.elt&&!se(t.elt)&&ee(t.elt,"id")){var i=document.getElementById(ee(t.elt,"id"));var a={preventScroll:v.focusScroll!==undefined?!v.focusScroll:!Q.config.defaultFocusScroll};if(i){if(t.start&&i.setSelectionRange){try{i.setSelectionRange(t.start,t.end)}catch(e){}}i.focus(a)}}c.classList.remove(Q.config.swappingClass);oe(n.elts,function(e){if(e.classList){e.classList.add(Q.config.settlingClass)}ce(e,"htmx:afterSwap",u)});if(O(f,/HX-Trigger-After-Swap:/i)){var o=l;if(!se(l)){o=re().body}_e(f,"HX-Trigger-After-Swap",o)}var s=function(){oe(n.tasks,function(e){e.call()});oe(n.elts,function(e){if(e.classList){e.classList.remove(Q.config.settlingClass)}ce(e,"htmx:afterSettle",u)});if(u.pathInfo.anchor){var e=re().getElementById(u.pathInfo.anchor);if(e){e.scrollIntoView({block:"start",behavior:"auto"})}}if(n.title&&!m){var t=C("title");if(t){t.innerHTML=n.title}else{window.document.title=n.title}}Cr(n.elts,v);if(O(f,/HX-Trigger-After-Settle:/i)){var r=l;if(!se(l)){r=re().body}_e(f,"HX-Trigger-After-Settle",r)}ie(p)};if(v.settleDelay>0){setTimeout(s,v.settleDelay)}else{s()}}catch(e){fe(l,"htmx:swapError",u);ie(x);throw e}};var b=Q.config.globalViewTransitions;if(v.hasOwnProperty("transition")){b=v.transition}if(b&&ce(l,"htmx:beforeTransition",u)&&typeof Promise!=="undefined"&&document.startViewTransition){var w=new Promise(function(e,t){p=e;x=t});var S=y;y=function(){document.startViewTransition(function(){S();return w})}}if(v.swapDelay>0){setTimeout(y,v.swapDelay)}else{y()}}if(a){fe(l,"htmx:responseError",le({error:"Response Status Error Code "+f.status+" from "+u.pathInfo.requestPath},u))}}var Xr={};function Dr(){return{init:function(e){return null},onEvent:function(e,t){return true},transformResponse:function(e,t,r){return e},isInlineSwap:function(e){return false},handleSwap:function(e,t,r,n){return false},encodeParameters:function(e,t,r){return null}}}function Ur(e,t){if(t.init){t.init(r)}Xr[e]=le(Dr(),t)}function Br(e){delete Xr[e]}function Fr(e,r,n){if(e==undefined){return r}if(r==undefined){r=[]}if(n==undefined){n=[]}var t=te(e,"hx-ext");if(t){oe(t.split(","),function(e){e=e.replace(/ /g,"");if(e.slice(0,7)=="ignore:"){n.push(e.slice(7));return}if(n.indexOf(e)<0){var t=Xr[e];if(t&&r.indexOf(t)<0){r.push(t)}}})}return Fr(u(e),r,n)}var Vr=false;re().addEventListener("DOMContentLoaded",function(){Vr=true});function jr(e){if(Vr||re().readyState==="complete"){e()}else{re().addEventListener("DOMContentLoaded",e)}}function _r(){if(Q.config.includeIndicatorStyles!==false){re().head.insertAdjacentHTML("beforeend","<style>                      ."+Q.config.indicatorClass+"{opacity:0}                      ."+Q.config.requestClass+" ."+Q.config.indicatorClass+"{opacity:1; transition: opacity 200ms ease-in;}                      ."+Q.config.requestClass+"."+Q.config.indicatorClass+"{opacity:1; transition: opacity 200ms ease-in;}                    </style>")}}function zr(){var e=re().querySelector('meta[name="htmx-config"]');if(e){return E(e.content)}else{return null}}function $r(){var e=zr();if(e){Q.config=le(Q.config,e)}}jr(function(){$r();_r();var e=re().body;zt(e);var t=re().querySelectorAll("[hx-trigger='restored'],[data-hx-trigger='restored']");e.addEventListener("htmx:abort",function(e){var t=e.target;var r=ae(t);if(r&&r.xhr){r.xhr.abort()}});const r=window.onpopstate?window.onpopstate.bind(window):null;window.onpopstate=function(e){if(e.state&&e.state.htmx){ar();oe(t,function(e){ce(e,"htmx:restored",{document:re(),triggerEvent:ce})})}else{if(r){r(e)}}};setTimeout(function(){ce(e,"htmx:load",{});e=null},0)});return Q}()});

    [File Ends] public/js/htmx.min.js

    [File Begins] public/js/package-lock.json
    {
      "name": "js",
      "version": "1.0.0",
      "lockfileVersion": 3,
      "requires": true,
      "packages": {
        "": {
          "name": "js",
          "version": "1.0.0",
          "license": "ISC",
          "dependencies": {
            "@antonz/codapi": "^0.17.0"
          }
        },
        "node_modules/@antonz/codapi": {
          "version": "0.17.0",
          "resolved": "https://registry.npmjs.org/@antonz/codapi/-/codapi-0.17.0.tgz",
          "integrity": "sha512-NEnnkXnNZauIzh9ou5vwMkomWMH0zyuQ+0TOyZTn1weq1MlDSxfORB/QYhAvbMplCHgAeXJYiMfkK50CrI+hKA=="
        }
      }
    }

    [File Ends] public/js/package-lock.json

    [File Begins] public/js/package.json
    {
      "name": "js",
      "version": "1.0.0",
      "description": "",
      "main": "events.js",
      "scripts": {
        "test": "echo \"Error: no test specified\" && exit 1"
      },
      "author": "",
      "license": "ISC",
      "dependencies": {
        "@antonz/codapi": "^0.17.0"
      }
    }

    [File Ends] public/js/package.json

    [File Begins] public/js/workflows.js
    // List of tool names enabled by default:
    var enabledTools = [
      // "webget", // Retrieves a single page from a given url
      // "websearch", // Searches a given query and returns n urls to pass to webget tool
    ];
    
    const uploadButton = document.getElementById('upload');
    const fileInput = document.getElementById('file-input');
    const form = document.querySelector('form'); // Select the form element
    
    // Prevent the default form submit behavior
    form.addEventListener('submit', function(event) {
      event.preventDefault(); // Prevent form from submitting traditionally
      console.log('Form submission prevented');
    });
    
    uploadButton.addEventListener('click', function () {
      console.log('Upload button clicked');
      fileInput.click(); // Trigger the file input dialog
    });
    
    fileInput.addEventListener('change', function () {
      const file = fileInput.files[0];
      if (file) {
        fileHandler(file);
      }
    });
    
    async function fileHandler(file) {
      if (file) {
        await uploadFile(file);
      }
      // reset the file input
      fileInput.value = null;
      // append the file to the chat view
      console.log("uploading...")
      console.log(file.name)
    }
    
    async function uploadFile(file) {
      const formData = new FormData();
      formData.append('file', file);
    
      try {
        const response = await fetch('/upload', {
          method: 'POST',
          body: formData
        });
    
        const data = await response.json();
    
        if (data.status === 'success' && data.callback === 'image') {
          // Image file detected, trigger another request
          console.log('Image file uploaded successfully')
          return //await fetchImageProcessingResult("./public/uploads/" + file.name);
        } else {
          // Handle other cases or non-image files
          console.log('File uploaded for processing:', data);
        }
      } catch (error) {
        console.error('Error:', error);
      }
    }
    
    async function createChat(prompt, msg, model) {
      const chatUrl = 'http://localhost:8080/chats';
    
      try {
        chatData = {
          Prompt: prompt,
          Response: msg,
          Model: model
        };
    
        //console.log(chatData);
    
        const response = await fetch(chatUrl, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify(chatData)
        });
    
        if (!response.ok) {
          throw new Error(`HTTP error! status: ${response.status}`);
        }
      } catch (error) {
        console.error('Error creating new chat:', error);
      }
    }
    
    async function fetchImageProcessingResult(fileName) {
      console.log(fileName)
    
      modelChain[0] = "bakllava";
    
      let payload = {
        modelPath: "models/bakllava/bakllava-1.Q8_0.gguf",
        mmproj: "models/bakllava/mmproj-model-f16.gguf",
        image: fileName,
        prompt: "Describe the image in detail.",
        contextSize: "4096",
        seed: "-1",
        temp: "0.7",
        responseDelimiter: "encode_image_with_clip:",
        socketNumber: "1"
      };
    
      console.log(payload);
    
      socket1.send(JSON.stringify(payload));
    }
    
    async function insertImageIntoChat(fileName) {
      try {
        const responseDiv = document.createElement('div');
        responseDiv.classList.add(
          'response',
          'rounded-2',
          'mt-3',
          'overflow-y-auto'
        );
    
        // Append image to response div
        const image = document.createElement('img');
        image.src = fileName;
        image.classList.add('img-fluid', 'rounded-2');
        responseDiv.appendChild(image);
    
        // Append response div to chat container
        chatContainer.appendChild(responseDiv);
      } catch (error) {
        console.error('Error generating image:', error);
      }
    }
    
    // Implement a workflow manager that handles the sequence of interactions with different models. 
    // This manager can take a workflow configuration and execute the steps accordingly.
    const workflows = {
      defaultFlow: ['model1', 'model2'],
      // Additional workflows
    };
    
    function executeWorkflow(workflowName) {
      const models = workflows[workflowName];
      // Logic to interact with models in the specified order
    }

    [File Ends] public/js/workflows.js

    [File Begins] public/templates/alerts.html
    <!-- WebSocket connection container -->
    <div id="download-progress" hx-ext="ws" ws-connect="/wsdownload-progress">
      <!-- Progress will be displayed here -->
      <div class="toast-container position-fixed bottom-0 end-0 p-3" id="toast-container">
        <div id="live-toast" class="toast" aria-live="assertive" aria-atomic="true" data-bs-autohide="false">
          <div class="toast-header">
            <img src="..." class="rounded me-2" alt="...">
            <strong class="me-auto">Download Progress</strong>
            <small class="toast-time">Just now</small>
            <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
          </div>
          <div class="toast-body">
            <!-- Progress message will be updated here -->
            <div id="progress-message">Starting download...</div>
          </div>
        </div>
      </div>
    </div>
    
    <script>
      // Define function to update progress
      function updateProgress(progress) {
        var progressMessage = document.getElementById('progress-message');
        progressMessage.textContent = 'Downloaded ' + progress.Current + '% of ' + progress.Total + '%';
      }
    
      // HTMX WebSocket event listeners
      htmx.on("htmx:wsOnMessage", function (evt) {
        var progress = JSON.parse(evt.detail.content);
        updateProgress(progress);
      });
    
      htmx.on("htmx:wsOnOpen", function (evt) {
        console.log("WebSocket connection opened");
      });
    
      htmx.on("htmx:wsOnClose", function (evt) {
        console.log("WebSocket connection closed");
      });
    
      // Additional functions as needed for scrolling, highlighting, etc.
      // ... (similar to previous <script> content) ...
    </script>

    [File Ends] public/templates/alerts.html

    [File Begins] public/templates/boxheader.html
    <head>
      <style>
        .box-container {
          perspective: 600px; /* Adjusted to be half of the original to maintain the perspective effect */
          position: absolute;
          top: 25px;
          left: 50%;
          transform: translateY(-5%);
          transform: translateX(-50%);
          z-index: 2;
        }
    
        .box {
          width: 20px; /* Half of the original size */
          height: 20px; /* Half of the original size */
          position: relative;
          transform-style: preserve-3d;
          transition: transform 0.5s ease;
          transform-origin: center;
        }
    
        .box .panel {
          position: absolute;
          width: 20px; /* Half of the original size */
          height: 20px; /* Half of the original size */
          background: rgba(255, 255, 255, 0.1);
          border: 1px solid rgba(255, 255, 255, 0.25);
        }
    
        .box .panel.front {
          transform: translateZ(12.5px); /* Half of the original size */
        }
    
        .box .panel.back {
          transform: rotateY(180deg) translateZ(12.5px); /* Half of the original size */
        }
    
        .box .panel.top {
          transform: rotateX(90deg) translateZ(12.5px); /* Half of the original size */
        }
    
        .box .panel.bottom {
          transform: rotateX(-90deg) translateZ(12.5px); /* Half of the original size */
        }
    
        .box .panel.left {
          transform: rotateY(-90deg) translateZ(12.5px); /* Half of the original size */
        }
    
        .box .panel.right {
          transform: rotateY(90deg) translateZ(12.5px); /* Half of the original size */
        }
    
        @keyframes spin {
          from {
            transform: rotateX(0deg) rotateY(0deg) rotateZ(0deg);
          }
    
          to {
            transform: rotateX(360deg) rotateY(360deg) rotateZ(360deg);
          }
        }
    
        .box:hover {
          animation: spin linear infinite;
        }
      </style>
    </head>
    
    <div id="row cb-header">
      <div class="box-container">
        <div class="box">
          <div class="panel front"></div>
          <div class="panel back"></div>
          <div class="panel top"></div>
          <div class="panel bottom"></div>
          <div class="panel left"></div>
          <div class="panel right"></div>
        </div>
      </div>
    </div>
    
    <script>
      // Get the elements
      const boxcontainer = document.querySelector(".box-container");
      const box = document.querySelector(".box");
    
      var rotationInterval = null;
      let r = 45;
    
      function getRandomRotation() {
        r += -90;
        const rotationX = r;
        const rotationY = r;
        const rotationZ = -180;
        return `rotateX(${rotationX}deg) rotateY(${rotationY}deg) rotateZ(${rotationZ}deg)`;
      }
    
      function chatRotation() {
        r += -10;
        const rotationX = r;
        const rotationY = r;
        const rotationZ = -90;
        return `rotateX(${rotationX}deg) rotateY(${rotationY}deg) rotateZ(${rotationZ}deg)`;
      }
    
      // Add a click event listener to rotate the box on click
      boxcontainer.addEventListener("click", function () {
        const newRotation = getRandomRotation();
        box.style.transition = "transform 0.5s";
        box.style.transform = newRotation;
        // chatbox.classList.toggle("expand"); // This line references an undefined "chatbox" which is not present in the provided code
      });
    
      box.addEventListener("mouseover", function () {
        // Check if event listener is already added
        if (box.hasAttribute("data-event-added")) {
          return;
        }
    
        // Only spin the box once when the mouse is over it
        box.setAttribute("data-event-added", true);
    
        // Rotate the box
        const newRotation = getRandomRotation();
        box.style.transition = "transform 0.5s ease";
        box.style.transform = newRotation;
    
        // Remove the event listener after the animation is done
        setTimeout(function () {
          box.removeAttribute("data-event-added");
        }, 1000);
      });
    
      setInterval(function () {
        const transform = window.getComputedStyle(box).getPropertyValue("transform");
      }, 1000);
    
      // Rotate the box on page load
      const newRotation = getRandomRotation();
      box.style.transition = "transform 0.5s ease";
      box.style.transform = newRotation;
    </script>

    [File Ends] public/templates/boxheader.html

    [File Begins] public/templates/chat.html
    <style>
      .loadership_JTACT {
        display: flex;
        position: relative;
        width: 68px;
        height: 68px;
      }
    
      .loadership_JTACT div {
        animation: loadership_JTACT_roller 1.2s infinite;
        animation-timing-function: cubic-bezier(0.5, 0, 0.5, 1);
        transform-origin: 34px 34px;
      }
    
      .loadership_JTACT div:after {
        content: " ";
        display: block;
        position: absolute;
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: #ffffff;
      }
    
      .loadership_JTACT div:nth-child(1) {
        animation-delay: 0.00s;
      }
    
      .loadership_JTACT div:nth-child(1):after {
        top: 60px;
        left: 30px;
      }
    
    
      .loadership_JTACT div:nth-child(2) {
        animation-delay: -0.04s;
      }
    
      .loadership_JTACT div:nth-child(2):after {
        top: 56px;
        left: 45px;
      }
    
    
      .loadership_JTACT div:nth-child(3) {
        animation-delay: -0.07s;
      }
    
      .loadership_JTACT div:nth-child(3):after {
        top: 45px;
        left: 56px;
      }
    
    
      .loadership_JTACT div:nth-child(4) {
        animation-delay: -0.11s;
      }
    
      .loadership_JTACT div:nth-child(4):after {
        top: 30px;
        left: 60px;
      }
    
    
      .loadership_JTACT div:nth-child(5) {
        animation-delay: -0.14s;
      }
    
      .loadership_JTACT div:nth-child(5):after {
        top: 15px;
        left: 56px;
      }
    
    
      .loadership_JTACT div:nth-child(6) {
        animation-delay: -0.18s;
      }
    
      .loadership_JTACT div:nth-child(6):after {
        top: 4px;
        left: 45px;
      }
    
    
      .loadership_JTACT div:nth-child(7) {
        animation-delay: -0.22s;
      }
    
      .loadership_JTACT div:nth-child(7):after {
        top: 0px;
        left: 30px;
      }
    
    
      .loadership_JTACT div:nth-child(8) {
        animation-delay: -0.25s;
      }
    
      .loadership_JTACT div:nth-child(8):after {
        top: 4px;
        left: 15px;
      }
    
    
      .loadership_JTACT div:nth-child(9) {
        animation-delay: -0.29s;
      }
    
      .loadership_JTACT div:nth-child(9):after {
        top: 15px;
        left: 4px;
      }
    
    
      .loadership_JTACT div:nth-child(10) {
        animation-delay: -0.32s;
      }
    
      .loadership_JTACT div:nth-child(10):after {
        top: 30px;
        left: 0px;
      }
    
    
      .loadership_JTACT div:nth-child(11) {
        animation-delay: -0.36s;
      }
    
      .loadership_JTACT div:nth-child(11):after {
        top: 45px;
        left: 4px;
      }
    
    
      .loadership_JTACT div:nth-child(12) {
        animation-delay: -0.40s;
      }
    
      .loadership_JTACT div:nth-child(12):after {
        top: 56px;
        left: 15px;
      }
    
    
    
      @keyframes loadership_JTACT_roller {
        0% {
          transform: rotate(0deg);
        }
    
        100% {
          transform: rotate(360deg);
        }
      }
    </style>
    
    <div name="chat-{{.turnID}}" id="chat-{{.turnID}}" hx-ext="ws" ws-connect="{{.wsRoute}}">
      <!-- <div name="chat-{{.turnID}}" id="chat-{{.turnID}}"> -->
      <div class="row">
        <div id="prompt-{{.turnID}}" class="user-prompt rounded-2 mt-3 pb-3">
          <div>
            <span class="badge my-3 mx-1" style="background-color: var(--et-red);">{{.username}}</span>
          </div>
          <!-- Using a hidden form to send the message -->
          <form id="hidden-form-{{.turnID}}" style="display:none;" hx-trigger="load" ws-send>
            <!-- <form id="hidden-form-{{.turnID}}" style="display:none;" hx-trigger="load">   -->
            <!-- Get the selectedModels in localStorage and send over websocket -->
            <input type="hidden" name="model" value="{{.model}}">
            <input type="hidden" name="chat_message" value="{{.message}}">
          </form>
          <div>
            <span class="message-content mx-1">{{.message}}</span>
          </div>
        </div>
      </div>
      <div class="row">
        <div id="response-{{.turnID}}" class="response rounded-2 mt-3 pb-3 overflow-y-auto">
          <div>
            <span class="badge my-3 mx-1" style="background-color: var(--et-purple);">{{.assistant}}</span>
          </div>
          <!-- Messages received from WebSocket will be appended here -->
          <div name="chat-{{.turnID}}" id="response-content-{{.turnID}}" hx-trigger="load, customEndOfStream"
            hx-on:load="highlight()">
            <div class="loadership_JTACT">
              <div></div>
              <div></div>
              <div></div>
              <div></div>
              <div></div>
              <div></div>
              <div></div>
              <div></div>
              <div></div>
              <div></div>
              <div></div>
              <div></div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- <codapi-settings url="http://localhost:1313/v1"></codapi-settings> -->
    <script src="js/node_modules/@antonz/codapi/dist/snippet.js"></script>
    <script>
      htmx.on("htmx:wsOpen", function (evt) {
        console.log("WebSocket opened");
        const box = document.querySelector(".box");
        if (box) {
          box.style.transition = "transform 0.5s ease";
    
          // Restore original rotation
          box.style.transform = "rotateX(45deg) rotateY(45deg) rotateZ(-180deg)";
        }
        scrollToBottom();
      });
    
      htmx.on("htmx:wsOnClose", function (evt) {
        console.log("WebSocket closed");
        highlight();
        const box = document.querySelector(".box");
        const newRotation = getRandomRotation();
        box.style.transition = "transform 0.5s ease";
        box.style.transform = newRotation;
      });
    
      function scrollToBottom() {
        setTimeout(function () {
          window.scrollTo(0, document.body.scrollHeight);
        }, 0); // Timeout ensures the DOM has been painted
      }
    
      // Call scrollToBottom in the htmx event handlers
      htmx.on("htmx:wsAfterMessage", function (evt) {
        const box = document.querySelector(".box");
        const newRotation = chatRotation();
        box.style.transition = "transform 0.2s ease";
        box.style.transform = newRotation;
        setViewHeight("prompt-view");
        highlight();
        scrollToBottom();
      });
    
      // Highlight code blocks
      function highlight() {
        const container = document.getElementById("response-content-{{.turnID}}");
        container.querySelectorAll('pre code').forEach((block, index) => {
          if (!block.hasAttribute('data-snippet-added')) {
            hljs.highlightElement(block);
    
            block.classList.add('rounded-2');
    
            // Create a new snippet element only if not already added
            const snippet = document.createElement('codapi-snippet');
            snippet.setAttribute('url', 'http://localhost:1313/v1');
            snippet.setAttribute('engine', 'browser');
            snippet.setAttribute('sandbox', 'javascript');
            snippet.setAttribute('editor', 'basic');
    
            // Create a button to copy code
            const copyButton = document.createElement('button');
            copyButton.classList.add('btn', 'btn-link');
            copyButton.innerHTML = `
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
              <path fill="#FFFFFF" fill-rule="evenodd" d="M7.263 3.26A2.25 2.25 0 0 1 9.5 1.25h5a2.25 2.25 0 0 1 2.237 2.01c.764.016 1.423.055 1.987.159c.758.14 1.403.404 1.928.93c.602.601.86 1.36.982 2.26c.116.866.116 1.969.116 3.336v6.11c0 1.367 0 2.47-.116 3.337c-.122.9-.38 1.658-.982 2.26c-.602.602-1.36.86-2.26.982c-.867.116-1.97.116-3.337.116h-6.11c-1.367 0-2.47 0-3.337-.116c-.9-.122-1.658-.38-2.26-.982c-.602-.602-.86-1.36-.981-2.26c-.117-.867-.117-1.97-.117-3.337v-6.11c0-1.367 0-2.47.117-3.337c.12-.9.38-1.658.981-2.26c.525-.525 1.17-.79 1.928-.929c.564-.104 1.224-.143 1.987-.159Zm1.487.741V4.5c0 .414.336.75.75.75h5a.75.75 0 0 0 .75-.75v-1a.75.75 0 0 0-.75-.75h-5a.75.75 0 0 0-.75.75v.501Zm7.985.76A2.25 2.25 0 0 1 14.5 6.75h-5a2.25 2.25 0 0 1-2.235-1.99c-.718.016-1.272.052-1.718.134c-.566.104-.895.272-1.138.515c-.277.277-.457.665-.556 1.4c-.101.754-.103 1.756-.103 3.191v6c0 1.435.002 2.436.103 3.192c.099.734.28 1.122.556 1.399c.277.277.665.457 1.4.556c.754.101 1.756.103 3.191.103h6c1.435 0 2.436-.002 3.192-.103c.734-.099 1.122-.28 1.399-.556c.277-.277.457-.665.556-1.4c.101-.755.103-1.756.103-3.191v-6c0-1.435-.002-2.437-.103-3.192c-.099-.734-.28-1.122-.556-1.399c-.244-.243-.572-.41-1.138-.515c-.446-.082-1-.118-1.718-.133ZM6.25 14.5a.75.75 0 0 1 .75-.75h8a.75.75 0 0 1 0 1.5H7a.75.75 0 0 1-.75-.75Zm0 3.5a.75.75 0 0 1 .75-.75h5.5a.75.75 0 0 1 0 1.5H7a.75.75 0 0 1-.75-.75Z" clip-rule="evenodd"/>
            </svg>`;
            copyButton.onclick = function () {
              navigator.clipboard.writeText(block.textContent).then(() => {
                alert('Code copied to clipboard!');
              }, () => {
                alert('Failed to copy code.');
              });
            };
    
            // Insert the snippet and the button after the code block
            const wrapper = document.createElement('div');
            wrapper.setAttribute('id', 'snippet-wrapper');
            wrapper.classList.add('rounded-2', 'mb-3');
            // Add background color to the wrapper rgb #272822
            wrapper.style.backgroundColor = '#272822';
            wrapper.appendChild(snippet);
            wrapper.appendChild(copyButton);
    
            block.parentNode.insertAdjacentElement('afterend', wrapper);
    
            // Mark this code block as having a snippet added
            block.setAttribute('data-snippet-added', 'true');
          }
        });
      }
    </script>

    [File Ends] public/templates/chat.html

    [File Begins] public/templates/drawflowdata.json
    {"drawflow":{"Home":{"data":{"1":{"id":1,"name":"welcome","data":{},"class":"welcome","html":"\n    <div>\n      <div class=\"title-box\">👏 Welcome!!</div>\n      <div class=\"box\">\n        <p>Simple flow library <b>demo</b>\n        <a href=\"https://github.com/jerosoler/Drawflow\" target=\"_blank\">Drawflow</a> by <b>Jero Soler</b></p><br>\n\n        <p>Multiple input / outputs<br>\n           Data sync nodes<br>\n           Import / export<br>\n           Modules support<br>\n           Simple use<br>\n           Type: Fixed or Edit<br>\n           Events: view console<br>\n           Pure Javascript<br>\n        </p>\n        <br>\n        <p><b><u>Shortkeys:</u></b></p>\n        <p>🎹 <b>Delete</b> for remove selected<br>\n        💠 Mouse Left Click == Move<br>\n        ❌ Mouse Right == Delete Option<br>\n        🔍 Ctrl + Wheel == Zoom<br>\n        📱 Mobile support<br>\n        ...</p>\n      </div>\n    </div>\n    ", "typenode": false, "inputs":{},"outputs":{},"pos_x":50,"pos_y":50},"2":{"id":2,"name":"slack","data":{},"class":"slack","html":"\n          <div>\n            <div class=\"title-box\"><i class=\"fab fa-slack\"></i> Slack chat message</div>\n          </div>\n          ", "typenode": false, "inputs":{"input_1":{"connections":[{"node":"7","input":"output_1"}]}},"outputs":{},"pos_x":1028,"pos_y":87},"3":{"id":3,"name":"telegram","data":{"channel":"channel_2"},"class":"telegram","html":"\n          <div>\n            <div class=\"title-box\"><i class=\"fab fa-telegram-plane\"></i> Telegram bot</div>\n            <div class=\"box\">\n              <p>Send to telegram</p>\n              <p>select channel</p>\n              <select df-channel>\n                <option value=\"channel_1\">Channel 1</option>\n                <option value=\"channel_2\">Channel 2</option>\n                <option value=\"channel_3\">Channel 3</option>\n                <option value=\"channel_4\">Channel 4</option>\n              </select>\n            </div>\n          </div>\n          ", "typenode": false, "inputs":{"input_1":{"connections":[{"node":"7","input":"output_1"}]}},"outputs":{},"pos_x":1032,"pos_y":184},"4":{"id":4,"name":"email","data":{},"class":"email","html":"\n            <div>\n              <div class=\"title-box\"><i class=\"fas fa-at\"></i> Send Email </div>\n            </div>\n            ", "typenode": false, "inputs":{"input_1":{"connections":[{"node":"5","input":"output_1"}]}},"outputs":{},"pos_x":1033,"pos_y":439},"5":{"id":5,"name":"template","data":{"template":"Write your template"},"class":"template","html":"\n            <div>\n              <div class=\"title-box\"><i class=\"fas fa-code\"></i> Template</div>\n              <div class=\"box\">\n                Ger Vars\n                <textarea df-template></textarea>\n                Output template with vars\n              </div>\n            </div>\n            ", "typenode": false, "inputs":{"input_1":{"connections":[{"node":"6","input":"output_1"}]}},"outputs":{"output_1":{"connections":[{"node":"4","output":"input_1"},{"node":"11","output":"input_1"}]}},"pos_x":607,"pos_y":304},"6":{"id":6,"name":"github","data":{"name":"https://github.com/jerosoler/Drawflow"},"class":"github","html":"\n          <div>\n            <div class=\"title-box\"><i class=\"fab fa-github \"></i> Github Stars</div>\n            <div class=\"box\">\n              <p>Enter repository url</p>\n            <input type=\"text\" df-name>\n            </div>\n          </div>\n          ", "typenode": false, "inputs":{},"outputs":{"output_1":{"connections":[{"node":"5","output":"input_1"}]}},"pos_x":341,"pos_y":191},"7":{"id":7,"name":"prompt","data":{},"class":"prompt","html":"\n        <div>\n          <div class=\"title-box\"><i class=\"fab fa-prompt\"></i> prompt Message</div>\n        </div>\n        ", "typenode": false, "inputs":{},"outputs":{"output_1":{"connections":[{"node":"2","output":"input_1"},{"node":"3","output":"input_1"},{"node":"11","output":"input_1"}]}},"pos_x":347,"pos_y":87},"11":{"id":11,"name":"log","data":{},"class":"log","html":"\n            <div>\n              <div class=\"title-box\"><i class=\"fas fa-file-signature\"></i> Save log file </div>\n            </div>\n            ", "typenode": false, "inputs":{"input_1":{"connections":[{"node":"5","input":"output_1"},{"node":"7","input":"output_1"}]}},"outputs":{},"pos_x":1031,"pos_y":363}}},"Other":{"data":{"8":{"id":8,"name":"personalized","data":{},"class":"personalized","html":"\n            <div>\n              Personalized\n            </div>\n            ", "typenode": false, "inputs":{"input_1":{"connections":[{"node":"12","input":"output_1"},{"node":"12","input":"output_2"},{"node":"12","input":"output_3"},{"node":"12","input":"output_4"}]}},"outputs":{"output_1":{"connections":[{"node":"9","output":"input_1"}]}},"pos_x":764,"pos_y":227},"9":{"id":9,"name":"dbclick","data":{"name":"Hello World!!"},"class":"dbclick","html":"\n            <div>\n            <div class=\"title-box\"><i class=\"fas fa-mouse\"></i> Db Click</div>\n              <div class=\"box dbclickbox\" ondblclick=\"showpopup(event)\">\n                Db Click here\n                <div class=\"modal\" style=\"display:none\">\n                  <div class=\"modal-content\">\n                    <span class=\"close\" onclick=\"closemodal(event)\">&times;</span>\n                    Change your variable {name} !\n                    <input type=\"text\" df-name>\n                  </div>\n\n                </div>\n              </div>\n            </div>\n            ", "typenode": false, "inputs":{"input_1":{"connections":[{"node":"8","input":"output_1"}]}},"outputs":{"output_1":{"connections":[{"node":"12","output":"input_2"}]}},"pos_x":209,"pos_y":38},"12":{"id":12,"name":"multiple","data":{},"class":"multiple","html":"\n            <div>\n              <div class=\"box\">\n                Multiple!\n              </div>\n            </div>\n            ", "typenode": false, "inputs":{"input_1":{"connections":[]},"input_2":{"connections":[{"node":"9","input":"output_1"}]},"input_3":{"connections":[]}},"outputs":{"output_1":{"connections":[{"node":"8","output":"input_1"}]},"output_2":{"connections":[{"node":"8","output":"input_1"}]},"output_3":{"connections":[{"node":"8","output":"input_1"}]},"output_4":{"connections":[{"node":"8","output":"input_1"}]}},"pos_x":179,"pos_y":272}}}}}

    [File Ends] public/templates/drawflowdata.json

    [File Begins] public/templates/flow.html
    <!DOCTYPE html>
    <html lang="en">
    <head>
      <meta charset="UTF-8">
      <meta name="viewport" content="width=device-width, initial-scale=1.0">
      <meta http-equiv="X-UA-Compatible" content="ie=edge">
      <title>Drawflow</title>
    </head>
    <body>
      <script src="../js/drawflow.min.js"></script>
      <script src="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.13.0/js/all.min.js" integrity="sha256-KzZiKy0DWYsnwMF+X1DvQngQ2/FxF7MF3Ff72XcpuPs=" crossorigin="anonymous"></script>
      <link rel="stylesheet" type="text/css" href="../../css/drawflow/drawflow.min.css" />
      <link rel="stylesheet" type="text/css" href="../../css/drawflow/beautiful.css" />
      <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/font-awesome/5.13.0/css/all.min.css" integrity="sha256-h20CPZ0QyXlBuAw7A+KluUYx/3pK+c7lYEpqLTlxjYQ=" crossorigin="anonymous" />
      <link href="https://fonts.googleapis.com/css2?family=Roboto&display=swap" rel="stylesheet">
      <script src="https://cdn.jsdelivr.net/npm/sweetalert2@9"></script>
      <script src="https://unpkg.com/micromodal/dist/micromodal.min.js"></script>
    
    
      <header>
        <h2>Drawflow</h2>
      </header>
      <div class="wrapper">
        <div class="col">
          <div class="drag-drawflow" draggable="true" ondragstart="drag(event)" data-node="prompt">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
              <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-width="1.5" d="M8 12h1m7 0h-4m4-4h-1m-3 0H8m0 8h5M3 14v-4c0-3.771 0-5.657 1.172-6.828C5.343 2 7.229 2 11 2h2c3.771 0 5.657 0 6.828 1.172c.654.653.943 1.528 1.07 2.828M21 10v4c0 3.771 0 5.657-1.172 6.828C18.657 22 16.771 22 13 22h-2c-3.771 0-5.657 0-6.828-1.172c-.654-.653-.943-1.528-1.07-2.828"/>
            </svg>
            <span> Prompt</span>
          </div>
          <div class="drag-drawflow" draggable="true" ondragstart="drag(event)" data-node="slack">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
              <g fill="currentColor">
                  <path d="M16 10.5c0 .828-.448 1.5-1 1.5s-1-.672-1-1.5s.448-1.5 1-1.5s1 .672 1 1.5Zm-6 0c0 .828-.448 1.5-1 1.5s-1-.672-1-1.5S8.448 9 9 9s1 .672 1 1.5Z"/>
                  <path fill-rule="evenodd" d="M9.944 1.25H10a.75.75 0 0 1 0 1.5c-1.907 0-3.261.002-4.29.14c-1.005.135-1.585.389-2.008.812c-.423.423-.677 1.003-.812 2.009c-.138 1.028-.14 2.382-.14 4.289a.75.75 0 0 1-1.5 0v-.056c0-1.838 0-3.294.153-4.433c.158-1.172.49-2.121 1.238-2.87c.749-.748 1.698-1.08 2.87-1.238c1.14-.153 2.595-.153 4.433-.153Zm8.345 1.64c-1.027-.138-2.382-.14-4.289-.14a.75.75 0 0 1 0-1.5h.056c1.838 0 3.294 0 4.433.153c1.172.158 2.121.49 2.87 1.238c.748.749 1.08 1.698 1.238 2.87c.153 1.14.153 2.595.153 4.433V10a.75.75 0 0 1-1.5 0c0-1.907-.002-3.261-.14-4.29c-.135-1.005-.389-1.585-.812-2.008c-.423-.423-1.003-.677-2.009-.812ZM2 13.25a.75.75 0 0 1 .75.75c0 1.907.002 3.262.14 4.29c.135 1.005.389 1.585.812 2.008c.423.423 1.003.677 2.009.812c1.028.138 2.382.14 4.289.14a.75.75 0 0 1 0 1.5h-.056c-1.838 0-3.294 0-4.433-.153c-1.172-.158-2.121-.49-2.87-1.238c-.748-.749-1.08-1.698-1.238-2.87c-.153-1.14-.153-2.595-.153-4.433V14a.75.75 0 0 1 .75-.75Zm20 0a.75.75 0 0 1 .75.75v.056c0 1.838 0 3.294-.153 4.433c-.158 1.172-.49 2.121-1.238 2.87c-.749.748-1.698 1.08-2.87 1.238c-1.14.153-2.595.153-4.433.153H14a.75.75 0 0 1 0-1.5c1.907 0 3.262-.002 4.29-.14c1.005-.135 1.585-.389 2.008-.812c.423-.423.677-1.003.812-2.009c.138-1.027.14-2.382.14-4.289a.75.75 0 0 1 .75-.75ZM8.397 15.553a.75.75 0 0 1 1.05-.155c.728.54 1.607.852 2.553.852s1.825-.313 2.553-.852a.75.75 0 1 1 .894 1.204A5.766 5.766 0 0 1 12 17.75a5.766 5.766 0 0 1-3.447-1.148a.75.75 0 0 1-.156-1.049Z" clip-rule="evenodd"/>
              </g>
            </svg>
            <span> Assistant</span>
          </div>
          <div class="drag-drawflow" draggable="true" ondragstart="drag(event)" data-node="github">
            <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
              <g fill="none">
                  <path stroke="currentColor" stroke-linecap="round" stroke-width="1.5" d="M2 8V6c0-1.4 0-2.1.272-2.635a2.5 2.5 0 0 1 1.093-1.093C3.9 2 4.6 2 6 2c1.4 0 2.1 0 2.635.272a2.5 2.5 0 0 1 1.093 1.093C10 3.9 10 4.6 10 6v12c0 1.4 0 2.1-.272 2.635a2.5 2.5 0 0 1-1.093 1.092C8.1 22 7.4 22 6 22c-1.4 0-2.1 0-2.635-.273a2.5 2.5 0 0 1-1.093-1.092C2 20.1 2 19.4 2 18v-6m5 7H5"/>
                  <path stroke="currentColor" stroke-width="1.5" d="m13.314 4.929l-2.142 2.142c-.579.578-.867.867-1.02 1.235C10 8.673 10 9.082 10 9.9v9.656l8.97-8.97c.99-.99 1.486-1.485 1.671-2.056a2.5 2.5 0 0 0 0-1.545c-.185-.57-.68-1.066-1.67-2.056c-.99-.99-1.486-1.485-2.056-1.67a2.5 2.5 0 0 0-1.545 0c-.571.185-1.066.68-2.056 1.67Z"/>
                  <path fill="currentColor" d="M18 22v-.75v.75Zm0-8v.75V14Zm4 4h-.75h.75Zm-.273 2.635l-.668-.34l.668.34Zm-1.092 1.092l-.34-.668l.34.668Zm1.092-6.362l-.668.34l.668-.34Zm-1.092-1.092l-.34.668l.34-.668ZM13 22.75a.75.75 0 0 0 0-1.5v1.5Zm4-1.5a.75.75 0 0 0 0 1.5v-1.5Zm-1.5-6.5H18v-1.5h-2.5v1.5ZM21.25 18c0 .712 0 1.202-.032 1.58c-.03.371-.085.57-.159.715l1.337.68c.199-.39.28-.809.317-1.272c.038-.454.037-1.015.037-1.703h-1.5ZM18 22.75c.688 0 1.249 0 1.703-.037c.463-.037.882-.118 1.273-.317l-.681-1.337c-.145.074-.344.13-.714.16c-.38.03-.869.031-1.581.031v1.5Zm3.06-2.456a1.75 1.75 0 0 1-.765.765l.68 1.337a3.25 3.25 0 0 0 1.42-1.42l-1.336-.681ZM22.75 18c0-.688 0-1.249-.037-1.703c-.037-.463-.118-.882-.317-1.273l-1.337.682c.074.144.13.343.16.713c.03.38.031.869.031 1.581h1.5ZM18 14.75c.712 0 1.202 0 1.58.032c.371.03.57.085.715.159l.68-1.337c-.39-.199-.809-.28-1.272-.317c-.454-.038-1.015-.037-1.703-.037v1.5Zm4.396.274a3.25 3.25 0 0 0-1.42-1.42l-.681 1.337c.329.167.596.435.764.764l1.337-.68ZM13 21.25H6v1.5h7v-1.5Zm5 0h-1v1.5h1v-1.5Z"/>
              </g>
            </svg>
            <span> Tool</span>
          </div>
        </div>
        <div class="col-right">
          <div class="menu">
            <ul>
              <li onclick="editor.changeModule('Home'); changeModule(event);" class="selected">Main</li>
              <li onclick="editor.changeModule('Other'); changeModule(event);">Other Module</li>
            </ul>
          </div>
          <div id="drawflow" ondrop="drop(event)" ondragover="allowDrop(event)">
    
            <div class="btn-export" onclick="Swal.fire({ title: 'Export',
            html: '<pre><code>'+JSON.stringify(editor.export(), null,4)+'</code></pre>'
            })">Export</div>
            <div class="btn-clear" onclick="editor.clearModuleSelected()">Clear</div>
            <div class="btn-lock">
              <i id="lock" class="fas fa-lock" onclick="editor.editor_mode='fixed'; changeMode('lock');"></i>
              <i id="unlock" class="fas fa-lock-open" onclick="editor.editor_mode='edit'; changeMode('unlock');" style="display:none;"></i>
            </div>
            <div class="bar-zoom">
              <i class="fas fa-search-minus" onclick="editor.zoom_out()"></i>
              <i class="fas fa-search" onclick="editor.zoom_reset()"></i>
              <i class="fas fa-search-plus" onclick="editor.zoom_in()"></i>
            </div>
          </div>
        </div>
      </div>
    
      <script>
    
        var id = document.getElementById("drawflow");
        const editor = new Drawflow(id);
        editor.reroute = true;
        editor.reroute_fix_curvature = true;
        editor.force_first_input = false;
    
      /*
        editor.createCurvature = function(start_pos_x, start_pos_y, end_pos_x, end_pos_y, curvature_value, type) {
          var center_x = ((end_pos_x - start_pos_x)/2)+start_pos_x;
          return ' M ' + start_pos_x + ' ' + start_pos_y + ' L '+ center_x +' ' +  start_pos_y  + ' L ' + center_x + ' ' +  end_pos_y  + ' L ' + end_pos_x + ' ' + end_pos_y;
        }*/
    
    
    
        
    
        const dataToImport =  {"drawflow":{"Home":{"data":{"1":{"id":1,"name":"welcome","data":{},"class":"welcome","html":"\n    <div>\n      <div class=\"title-box\">👏 Welcome!!</div>\n      <div class=\"box\">\n        <p>Simple flow library <b>demo</b>\n        <a href=\"https://github.com/jerosoler/Drawflow\" target=\"_blank\">Drawflow</a> by <b>Jero Soler</b></p><br>\n\n        <p>Multiple input / outputs<br>\n           Data sync nodes<br>\n           Import / export<br>\n           Modules support<br>\n           Simple use<br>\n           Type: Fixed or Edit<br>\n           Events: view console<br>\n           Pure Javascript<br>\n        </p>\n        <br>\n        <p><b><u>Shortkeys:</u></b></p>\n        <p>🎹 <b>Delete</b> for remove selected<br>\n        💠 Mouse Left Click == Move<br>\n        ❌ Mouse Right == Delete Option<br>\n        🔍 Ctrl + Wheel == Zoom<br>\n        📱 Mobile support<br>\n        ...</p>\n      </div>\n    </div>\n    ", "typenode": false, "inputs":{},"outputs":{},"pos_x":50,"pos_y":50},"2":{"id":2,"name":"slack","data":{},"class":"slack","html":"\n          <div>\n            <div class=\"title-box\"><i class=\"fab fa-slack\"></i> Slack chat message</div>\n          </div>\n          ", "typenode": false, "inputs":{"input_1":{"connections":[{"node":"7","input":"output_1"}]}},"outputs":{},"pos_x":1028,"pos_y":87},"3":{"id":3,"name":"telegram","data":{"channel":"channel_2"},"class":"telegram","html":"\n          <div>\n            <div class=\"title-box\"><i class=\"fab fa-telegram-plane\"></i> Telegram bot</div>\n            <div class=\"box\">\n              <p>Send to telegram</p>\n              <p>select channel</p>\n              <select df-channel>\n                <option value=\"channel_1\">Channel 1</option>\n                <option value=\"channel_2\">Channel 2</option>\n                <option value=\"channel_3\">Channel 3</option>\n                <option value=\"channel_4\">Channel 4</option>\n              </select>\n            </div>\n          </div>\n          ", "typenode": false, "inputs":{"input_1":{"connections":[{"node":"7","input":"output_1"}]}},"outputs":{},"pos_x":1032,"pos_y":184},"4":{"id":4,"name":"email","data":{},"class":"email","html":"\n            <div>\n              <div class=\"title-box\"><i class=\"fas fa-at\"></i> Send Email </div>\n            </div>\n            ", "typenode": false, "inputs":{"input_1":{"connections":[{"node":"5","input":"output_1"}]}},"outputs":{},"pos_x":1033,"pos_y":439},"5":{"id":5,"name":"template","data":{"template":"Write your template"},"class":"template","html":"\n            <div>\n              <div class=\"title-box\"><i class=\"fas fa-code\"></i> Template</div>\n              <div class=\"box\">\n                Ger Vars\n                <textarea df-template></textarea>\n                Output template with vars\n              </div>\n            </div>\n            ", "typenode": false, "inputs":{"input_1":{"connections":[{"node":"6","input":"output_1"}]}},"outputs":{"output_1":{"connections":[{"node":"4","output":"input_1"},{"node":"11","output":"input_1"}]}},"pos_x":607,"pos_y":304},"6":{"id":6,"name":"github","data":{"name":"https://github.com/jerosoler/Drawflow"},"class":"github","html":"\n          <div>\n            <div class=\"title-box\"><i class=\"fab fa-github \"></i> Github Stars</div>\n            <div class=\"box\">\n              <p>Enter repository url</p>\n            <input type=\"text\" df-name>\n            </div>\n          </div>\n          ", "typenode": false, "inputs":{},"outputs":{"output_1":{"connections":[{"node":"5","output":"input_1"}]}},"pos_x":341,"pos_y":191},"7":{"id":7,"name":"prompt","data":{},"class":"prompt","html":"\n        <div>\n          <div class=\"title-box\"><i class=\"fab fa-prompt\"></i> prompt Message</div>\n        </div>\n        ", "typenode": false, "inputs":{},"outputs":{"output_1":{"connections":[{"node":"2","output":"input_1"},{"node":"3","output":"input_1"},{"node":"11","output":"input_1"}]}},"pos_x":347,"pos_y":87},"11":{"id":11,"name":"log","data":{},"class":"log","html":"\n            <div>\n              <div class=\"title-box\"><i class=\"fas fa-file-signature\"></i> Save log file </div>\n            </div>\n            ", "typenode": false, "inputs":{"input_1":{"connections":[{"node":"5","input":"output_1"},{"node":"7","input":"output_1"}]}},"outputs":{},"pos_x":1031,"pos_y":363}}},"Other":{"data":{"8":{"id":8,"name":"personalized","data":{},"class":"personalized","html":"\n            <div>\n              Personalized\n            </div>\n            ", "typenode": false, "inputs":{"input_1":{"connections":[{"node":"12","input":"output_1"},{"node":"12","input":"output_2"},{"node":"12","input":"output_3"},{"node":"12","input":"output_4"}]}},"outputs":{"output_1":{"connections":[{"node":"9","output":"input_1"}]}},"pos_x":764,"pos_y":227},"9":{"id":9,"name":"dbclick","data":{"name":"Hello World!!"},"class":"dbclick","html":"\n            <div>\n            <div class=\"title-box\"><i class=\"fas fa-mouse\"></i> Db Click</div>\n              <div class=\"box dbclickbox\" ondblclick=\"showpopup(event)\">\n                Db Click here\n                <div class=\"modal\" style=\"display:none\">\n                  <div class=\"modal-content\">\n                    <span class=\"close\" onclick=\"closemodal(event)\">&times;</span>\n                    Change your variable {name} !\n                    <input type=\"text\" df-name>\n                  </div>\n\n                </div>\n              </div>\n            </div>\n            ", "typenode": false, "inputs":{"input_1":{"connections":[{"node":"8","input":"output_1"}]}},"outputs":{"output_1":{"connections":[{"node":"12","output":"input_2"}]}},"pos_x":209,"pos_y":38},"12":{"id":12,"name":"multiple","data":{},"class":"multiple","html":"\n            <div>\n              <div class=\"box\">\n                Multiple!\n              </div>\n            </div>\n            ", "typenode": false, "inputs":{"input_1":{"connections":[]},"input_2":{"connections":[{"node":"9","input":"output_1"}]},"input_3":{"connections":[]}},"outputs":{"output_1":{"connections":[{"node":"8","output":"input_1"}]},"output_2":{"connections":[{"node":"8","output":"input_1"}]},"output_3":{"connections":[{"node":"8","output":"input_1"}]},"output_4":{"connections":[{"node":"8","output":"input_1"}]}},"pos_x":179,"pos_y":272}}}}}
        editor.start();
        editor.import(dataToImport);
    
    
    
      /*
        var welcome = `
        <div>
          <div class="title-box">👏 Welcome!!</div>
          <div class="box">
            <p>Simple flow library <b>demo</b>
            <a href="https://github.com/jerosoler/Drawflow" target="_blank">Drawflow</a> by <b>Jero Soler</b></p><br>
    
            <p>Multiple input / outputs<br>
               Data sync nodes<br>
               Import / export<br>
               Modules support<br>e
               Simple use<br>
               Type: Fixed or Edit<br>
               Events: view console<br>
               Pure Javascript<br>
            </p>
            <br>
            <p><b><u>Shortkeys:</u></b></p>
            <p>🎹 <b>Delete</b> for remove selected<br>
            💠 Mouse Left Click == Move<br>
            ❌ Mouse Right == Delete Option<br>
            🔍 Ctrl + Wheel == Zoom<br>
            📱 Mobile support<br>
            ...</p>
          </div>
        </div>
        `;
    */
    
    
        //editor.addNode(name, "typenode": false,  inputs, outputs, posx, posy, class, data, html);
        /*editor.addNode('welcome', 0, 0, 50, 50, 'welcome', {}, welcome );
        editor.addModule('Other');
        */
    
        // Events!
        editor.on('nodeCreated', function(id) {
          console.log("Node created " + id);
        })
    
        editor.on('nodeRemoved', function(id) {
          console.log("Node removed " + id);
        })
    
        editor.on('nodeSelected', function(id) {
          console.log("Node selected " + id);
        })
    
        editor.on('moduleCreated', function(name) {
          console.log("Module Created " + name);
        })
    
        editor.on('moduleChanged', function(name) {
          console.log("Module Changed " + name);
        })
    
        editor.on('connectionCreated', function(connection) {
          console.log('Connection created');
          console.log(connection);
        })
    
        editor.on('connectionRemoved', function(connection) {
          console.log('Connection removed');
          console.log(connection);
        })
    /*
        editor.on('mouseMove', function(position) {
          console.log('Position mouse x:' + position.x + ' y:'+ position.y);
        })
    */
        editor.on('nodeMoved', function(id) {
          console.log("Node moved " + id);
        })
    
        editor.on('zoom', function(zoom) {
          console.log('Zoom level ' + zoom);
        })
    
        editor.on('translate', function(position) {
          console.log('Translate x:' + position.x + ' y:'+ position.y);
        })
    
        editor.on('addReroute', function(id) {
          console.log("Reroute added " + id);
        })
    
        editor.on('removeReroute', function(id) {
          console.log("Reroute removed " + id);
        })
        /* DRAG EVENT */
    
        /* Mouse and Touch Actions */
    
        var elements = document.getElementsByClassName('drag-drawflow');
        for (var i = 0; i < elements.length; i++) {
          elements[i].addEventListener('touchend', drop, false);
          elements[i].addEventListener('touchmove', positionMobile, false);
          elements[i].addEventListener('touchstart', drag, false );
        }
    
        var mobile_item_selec = '';
        var mobile_last_move = null;
       function positionMobile(ev) {
         mobile_last_move = ev;
       }
    
       function allowDrop(ev) {
          ev.preventDefault();
        }
    
        function drag(ev) {
          if (ev.type === "touchstart") {
            mobile_item_selec = ev.target.closest(".drag-drawflow").getAttribute('data-node');
          } else {
          ev.dataTransfer.setData("node", ev.target.getAttribute('data-node'));
          }
        }
    
        function drop(ev) {
          if (ev.type === "touchend") {
            var parentdrawflow = document.elementFromPoint( mobile_last_move.touches[0].clientX, mobile_last_move.touches[0].clientY).closest("#drawflow");
            if(parentdrawflow != null) {
              addNodeToDrawFlow(mobile_item_selec, mobile_last_move.touches[0].clientX, mobile_last_move.touches[0].clientY);
            }
            mobile_item_selec = '';
          } else {
            ev.preventDefault();
            var data = ev.dataTransfer.getData("node");
            addNodeToDrawFlow(data, ev.clientX, ev.clientY);
          }
    
        }
    
        function addNodeToDrawFlow(name, pos_x, pos_y) {
          if(editor.editor_mode === 'fixed') {
            return false;
          }
          pos_x = pos_x * ( editor.precanvas.clientWidth / (editor.precanvas.clientWidth * editor.zoom)) - (editor.precanvas.getBoundingClientRect().x * ( editor.precanvas.clientWidth / (editor.precanvas.clientWidth * editor.zoom)));
          pos_y = pos_y * ( editor.precanvas.clientHeight / (editor.precanvas.clientHeight * editor.zoom)) - (editor.precanvas.getBoundingClientRect().y * ( editor.precanvas.clientHeight / (editor.precanvas.clientHeight * editor.zoom)));
    
    
          switch (name) {
            case 'prompt':
            var prompt = `
            <div>
              <div class="title-box"><i class="fab fa-prompt"></i> prompt Message</div>
            </div>
            `;
              editor.addNode('prompt', 0,  1, pos_x, pos_y, 'prompt', {}, prompt );
              break;
            case 'slack':
              var slackchat = `
              <div>
                <div class="title-box"><i class="fab fa-slack"></i> Slack chat message</div>
              </div>
              `
              editor.addNode('slack', 1, 0, pos_x, pos_y, 'slack', {}, slackchat );
              break;
            case 'github':
              var githubtemplate = `
              <div>
                <div class="title-box"><i class="fab fa-github "></i> Github Stars</div>
                <div class="box">
                  <p>Enter repository url</p>
                <input type="text" df-name>
                </div>
              </div>
              `;
              editor.addNode('github', 0, 1, pos_x, pos_y, 'github', { "name": ''}, githubtemplate );
              break;
    
            default:
          }
        }
    
      var transform = '';
      function showpopup(e) {
        e.target.closest(".drawflow-node").style.zIndex = "9999";
        e.target.children[0].style.display = "block";
        //document.getElementById("modalfix").style.display = "block";
    
        //e.target.children[0].style.transform = 'translate('+translate.x+'px, '+translate.y+'px)';
        transform = editor.precanvas.style.transform;
        editor.precanvas.style.transform = '';
        editor.precanvas.style.left = editor.canvas_x +'px';
        editor.precanvas.style.top = editor.canvas_y +'px';
        console.log(transform);
    
        //e.target.children[0].style.top  =  -editor.canvas_y - editor.container.offsetTop +'px';
        //e.target.children[0].style.left  =  -editor.canvas_x  - editor.container.offsetLeft +'px';
        editor.editor_mode = "fixed";
    
      }
    
       function closemodal(e) {
         e.target.closest(".drawflow-node").style.zIndex = "2";
         e.target.parentElement.parentElement.style.display  ="none";
         //document.getElementById("modalfix").style.display = "none";
         editor.precanvas.style.transform = transform;
           editor.precanvas.style.left = '0px';
           editor.precanvas.style.top = '0px';
          editor.editor_mode = "edit";
       }
    
        function changeModule(event) {
          var all = document.querySelectorAll(".menu ul li");
            for (var i = 0; i < all.length; i++) {
              all[i].classList.remove('selected');
            }
          event.target.classList.add('selected');
        }
    
        function changeMode(option) {
    
        //console.log(lock.id);
          if(option == 'lock') {
            lock.style.display = 'none';
            unlock.style.display = 'block';
          } else {
            lock.style.display = 'block';
            unlock.style.display = 'none';
          }
    
        }
    
      </script>
    </body>
    </html>

    [File Ends] public/templates/flow.html

    [File Begins] public/templates/header.html
    <div id="cb-header" class="position-fixed fixed-top dark-blur mt-0" hx-preserve>
      <div id="hgradient"></div>
      <div class="box-container">
        <div class="circle"></div>
        <div class="box">
          <div class="panel front" style="--bg-color: var(--et-purple);"></div>
          <div class="panel back" style="--bg-color: var(--et-green);"></div>
          <div class="panel top" style="--bg-color: var(--et-blue);"></div>
          <div class="panel bottom" style="--bg-color: var(--et-red);"></div>
          <div class="panel left" style="--bg-color: var(--et-yellow);"></div>
          <div class="panel right" style="--bg-color: var(--et-light);"></div>
        </div>
      </div>
    </div>
    
    <script>
      // Get the elements
      const boxcontainer = document.querySelector(".box-container");
      const box = document.querySelector(".box");
    
      var rotationInterval = null;
      let r = 45;
    
      function getRandomRotation() {
        r += -90;
        const rotationX = r;
        const rotationY = r;
        const rotationZ = -180;
        return `rotateX(${rotationX}deg) rotateY(${rotationY}deg) rotateZ(${rotationZ}deg)`;
      }
    
      function chatRotation() {
        r += 10;
        const rotationX = r;
        const rotationY = r;
        const rotationZ = 0;
        return `rotateX(${rotationX}deg) rotateY(${rotationY}deg) rotateZ(${rotationZ}deg)`;
      }
    
      // Add a click event listener to rotate the box on click
      boxcontainer.addEventListener("click", function () {
        const newRotation = getRandomRotation();
        box.style.transition = "transform 0.5s";
        box.style.transform = newRotation;
        // chatbox.classList.toggle("expand"); // This line references an undefined "chatbox" which is not present in the provided code
      });
    
      box.addEventListener("mouseover", function () {
        // Check if event listener is already added
        if (box.hasAttribute("data-event-added")) {
          return;
        }
    
        // Only spin the box once when the mouse is over it
        box.setAttribute("data-event-added", true);
    
        // Rotate the box
        const newRotation = getRandomRotation();
        box.style.transition = "transform 0.5s ease";
        box.style.transform = newRotation;
    
        // Remove the event listener after the animation is done
        setTimeout(function () {
          box.removeAttribute("data-event-added");
        }, 1000);
      });
    
      setInterval(function () {
        const transform = window.getComputedStyle(box).getPropertyValue("transform");
      }, 1000);
    
      // Rotate the box on page load
      const newRotation = getRandomRotation();
      box.style.transition = "transform 0.5s ease";
      box.style.transform = newRotation;
    </script>

    [File Ends] public/templates/header.html

    [File Begins] public/templates/index.html
    <!doctype html>
    
    <html lang="en" data-bs-core="modern" data-bs-theme="dark">
    
    <head>
      <meta charset="utf-8">
      <meta name="viewport" content="width=device-width, initial-scale=1">
    
      <title>Eternal</title>
    
      <!-- Fonts -->
      <link href="https://fonts.googleapis.com/css2?family=Roboto&display=swap" rel="stylesheet">
    
      <!-- Halfmoon CSS -->
      <link rel="stylesheet" href="css/halfmoon/halfmoon.css">
    
      <!-- Halfmoon modern core theme only -->
      <link rel="stylesheet" href="css/halfmoon/cores/halfmoon.modern.css">
    
      <!-- Custom Styles -->
      <link rel="stylesheet" href="css/styles.css">
      <link rel="stylesheet" href="css/header.css">
      <!-- <link rel="stylesheet" href="https://unpkg.com/@antonz/codapi@0.17.0/dist/snippet.css" /> -->
    
      <!-- Code Highlight -->
      <link rel="stylesheet" href="js/highlight/styles/github-dark-dimmed.min.css">
      <script src="js/highlight/highlight.js"></script>
      <script src="js/highlight/es/languages/go.min.js"></script>
      <script src="js/highlight/es/languages/python.min.js"></script>
      <script src="js/highlight/es/languages/rust.min.js"></script>
      <script src="js/highlight/es/languages/bash.min.js"></script>
      <script src="js/highlight/es/languages/yaml.min.js"></script>
      <script src="js/highlight/es/languages/json.min.js"></script>
      <script src="js/highlight/es/languages/markdown.min.js"></script>
      <script src="js/highlight/es/languages/javascript.min.js"></script>
      <script src="js/highlight/es/languages/typescript.min.js"></script>
      <script src="js/highlight/es/languages/css.min.js"></script>
    
      <!-- Bootstrap JS bundle with Popper -->
      <script src="js/bootstrap/bootstrap.bundle.min.js"></script>
    
      <!-- HTMX -->
      <!-- <script src="js/htmx.min.js"></script> -->
      <script src="https://unpkg.com/htmx.org@2.0.0-beta1/dist/htmx.min.js"></script>
      <script src="https://unpkg.com/htmx.org/dist/ext/ws.js"></script>
      <script src="https://unpkg.com/htmx-ext-sse@2.0.0/sse.js"></script>
      <!-- <script src="https://unpkg.com/htmx.org/dist/ext/sse.js"></script> -->
    
      <style>
    
      </style>
    </head>
    
    <body class="overflow-y-scroll">
    
      {{template "templates/header".}}
      <div id="content" class="container-fluid overflow-x-hidden">
        <div class="row mt-5">
          <div class="col-2"></div>
          <div id="chat-view" class="col-8">
            <div id="chat" class="row chat-container fs-5"></div>
          </div>
          <div id="info" class="col-2 mt-2 pt-1">
          </div>
        </div>
    
        <div class="row mt-1 pb-3 px-2 w-auto fixed-bottom dark-blur">
          <form>
            <div class="py-1" id="prompt-view">
              <div class="row">
                <button class="btn fw-medium" data-bs-toggle="/">
                  <span class="fs-4">Et<svg class="" width="16" height="16" viewBox="0 0 24 24"
                      xmlns="http://www.w3.org/2000/svg">
                      <path fill="#ffffff"
                        d="M6.676 11.946a.75.75 0 0 0 1.18-.925a7.882 7.882 0 0 1-1.01-1.677a.75.75 0 1 0-1.372.604c.316.72.728 1.394 1.202 1.998M4.84 7.672a.75.75 0 0 0 1.489-.178a5.115 5.115 0 0 1 .109-1.862a.75.75 0 1 0-1.455-.366a6.615 6.615 0 0 0-.144 2.406M6.007 3.08a.75.75 0 0 0 1.218.875a5.84 5.84 0 0 1 .621-.727a.75.75 0 0 0-1.06-1.061a7.396 7.396 0 0 0-.779.912m11.629 8.975a.75.75 0 0 0-1.18.925c.4.511.745 1.079 1.009 1.677a.75.75 0 1 0 1.373-.604a9.383 9.383 0 0 0-1.202-1.998m1.836 4.274a.75.75 0 0 0-1.49.178a5.114 5.114 0 0 1-.108 1.862a.75.75 0 1 0 1.454.366a6.616 6.616 0 0 0 .144-2.406m-1.168 4.592a.75.75 0 0 0-1.218-.875a5.9 5.9 0 0 1-.62.727a.75.75 0 0 0 1.06 1.061c.293-.293.552-.598.778-.912M12.082 7.573a.75.75 0 0 1 .127-1.053a9.384 9.384 0 0 1 1.998-1.202a.75.75 0 0 1 .605 1.373a7.881 7.881 0 0 0-1.678 1.01a.75.75 0 0 1-1.053-.128m3.747-2.056a.75.75 0 0 1 .656-.833a6.615 6.615 0 0 1 2.405.143a.75.75 0 0 1-.366 1.455a5.115 5.115 0 0 0-1.862-.109a.75.75 0 0 1-.833-.656m4.202.506a.75.75 0 0 1 1.046-.171c.314.226.619.485.912.778a.75.75 0 1 1-1.06 1.06a5.895 5.895 0 0 0-.728-.62a.75.75 0 0 1-.17-1.047M12.103 17.48a.75.75 0 1 0-.926-1.18c-.51.4-1.078.746-1.677 1.01a.75.75 0 0 0 .604 1.372a9.379 9.379 0 0 0 1.999-1.202m-4.275 1.836a.75.75 0 0 0-.178-1.49a5.114 5.114 0 0 1-1.862-.108a.75.75 0 0 0-.366 1.455a6.614 6.614 0 0 0 2.406.143m-4.592-1.168a.75.75 0 0 0 .875-1.218a5.892 5.892 0 0 1-.727-.62a.75.75 0 1 0-1.06 1.06c.293.293.597.552.912.778" />
                      <path fill="#ffffff"
                        d="M13.746 15.817a.75.75 0 0 1-1.347-.407c-1.28.605-2.914.783-4.504.558C4.685 15.513 1.25 13.316 1.25 9a.75.75 0 0 1 1.5 0c0 3.284 2.564 5.087 5.355 5.482a7.72 7.72 0 0 0 1.872.04a6.978 6.978 0 0 1-1.638-.932a.75.75 0 0 1 .492-1.348c-.548-1.255-.703-2.821-.487-4.347c.455-3.21 2.652-6.645 6.968-6.645a.75.75 0 0 1 0 1.5c-3.285 0-5.087 2.564-5.483 5.355a7.872 7.872 0 0 0-.073 1.423c.212-.465.487-.918.81-1.345a.75.75 0 0 1 1.336.587c1.23-.499 2.735-.634 4.203-.426c3.21.455 6.645 2.652 6.645 6.968a.75.75 0 0 1-1.5 0c0-3.285-2.564-5.087-5.355-5.483a7.985 7.985 0 0 0-.959-.078c.357.186.704.408 1.037.659a.75.75 0 0 1-.492 1.348c.548 1.255.703 2.821.487 4.347c-.455 3.21-2.652 6.645-6.968 6.645a.75.75 0 0 1 0-1.5c3.284 0 5.087-2.564 5.482-5.355a7.87 7.87 0 0 0 .073-1.423a7.192 7.192 0 0 1-.809 1.345" />
                    </svg>rnal</span>
                </button>
              </div>
              <div class="hstack">
    
                <div class="row ms-auto w-25">
                  <!-- WEB PAGE RETRIEVAL -->
                  <div class="col">
                    <button id="webget-btn" class="btn" data-bs-toggle="tooltip" data-bs-title="Web Retrieval">
                      <!-- <svg id="webget" width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <g fill="none" stroke="#ffffff" stroke-width="1.5">
                          <path
                            d="M3 10c0-3.771 0-5.657 1.172-6.828C5.343 2 7.229 2 11 2h2c3.771 0 5.657 0 6.828 1.172C21 4.343 21 6.229 21 10v4c0 3.771 0 5.657-1.172 6.828C18.657 22 16.771 22 13 22h-2c-3.771 0-5.657 0-6.828-1.172C3 19.657 3 17.771 3 14z" />
                          <path stroke-linecap="round" d="M12 6v2m0 0v2m0-2h-2m2 0h2m-6 6h8m-7 4h6" />
                        </g>
                      </svg> -->
                    </button>
                  </div>
                  <!-- Roles -->
                  <div class="col" id="imgstatus">
                    <!-- <button id="rolesBtn" class="btn" hx-post="/tool/websearch" hx-swap="none" data-bs-toggle="tooltip"
                      data-bs-title="Assistant Roles">
                      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
                        <g fill="currentColor">
                          <path
                            d="M6.005 13.368c.029-.296.26-.6.638-.702c.379-.102.73.047.903.29a.75.75 0 0 0 1.22-.873c-.55-.77-1.552-1.123-2.511-.866c-.96.257-1.651 1.064-1.743 2.006a.75.75 0 1 0 1.493.145Zm5.796-1.553c.029-.296.26-.6.638-.702c.379-.102.73.047.903.289a.75.75 0 0 0 1.22-.872c-.55-.77-1.552-1.123-2.511-.866c-.96.257-1.651 1.063-1.743 2.006a.75.75 0 0 0 1.493.145Zm1.399 4.416l.448-.602a.75.75 0 0 1-.885 1.211l-.01-.006a2.06 2.06 0 0 0-.485-.2c-.361-.098-.93-.163-1.686.04c-.756.202-1.215.543-1.48.808a2.064 2.064 0 0 0-.32.416l-.005.01a.75.75 0 0 1-1.372-.607l.689.298l-.689-.298l.001-.001v-.002l.003-.004l.003-.008l.011-.023l.032-.064c.027-.051.065-.118.115-.196c.1-.156.252-.36.469-.578c.436-.439 1.124-.924 2.155-1.2c1.031-.277 1.87-.2 2.467-.038c.297.08.53.18.695.266a2.682 2.682 0 0 1 .257.151l.02.015l.009.005l.003.003h.001l.002.002l-.447.602Z" />
                          <path fill-rule="evenodd"
                            d="m13.252 2.25l.042.02c1.167.547 1.692.791 2.235.963c.193.061.387.116.583.164c.552.134 1.122.197 2.395.334l.045.004c.808.087 1.48.16 2.01.28c.554.127 1.054.328 1.448.743c.23.24.414.521.546.827c.225.52.226 1.064.144 1.64c-.08.554-.253 1.232-.464 2.056l-.856 3.339c-.716 2.793-2.533 4.345-4.357 5.189c-.725 1.574-1.863 2.78-2.804 3.583l-.021.018c-.25.214-.497.425-.82.61c-.335.191-.724.34-1.269.493c-.544.152-.953.227-1.338.236c-.37.009-.687-.045-1.006-.1l-.028-.004c-2.321-.394-6.012-1.714-7.117-6.025l-.856-3.34c-.21-.823-.384-1.5-.464-2.056c-.082-.575-.081-1.118.144-1.639c.132-.306.317-.586.546-.827c.394-.415.894-.616 1.448-.742c.53-.122 1.201-.194 2.01-.28l.045-.005c.52-.056.921-.1 1.253-.14l.625-2.44c.211-.824.385-1.501.582-2.024c.203-.54.466-1.017.92-1.358c.265-.2.565-.348.884-.439c.55-.156 1.084-.066 1.622.113c.516.172 1.132.46 1.873.808Zm6.675 9.997c-.412 1.608-1.26 2.701-2.263 3.45a7.953 7.953 0 0 0-.18-3.207l-.93-3.632a.746.746 0 0 0 .338-.263c.173-.242.525-.39.904-.289c.378.101.608.406.637.702a.75.75 0 1 0 1.493-.145c-.091-.942-.783-1.749-1.742-2.006a2.37 2.37 0 0 0-2.084.416a6.985 6.985 0 0 0-.053-.146c-.203-.54-.466-1.017-.92-1.358a2.698 2.698 0 0 0-.884-.439c-.52-.147-1.026-.075-1.533.085a2.448 2.448 0 0 0-.322-.111c-.96-.257-1.962.096-2.512.866a.748.748 0 0 0-.132.547c-.55.252-.908.4-1.273.516l-.092.03l.434-1.697c.225-.877.38-1.474.543-1.91c.161-.428.296-.596.417-.687c.12-.09.254-.156.393-.196c.133-.038.329-.043.74.094c.422.14.958.39 1.752.762l.053.025c1.1.515 1.717.804 2.364 1.01c.225.07.453.134.682.19c.66.16 1.332.233 2.531.362l.059.006c.865.093 1.448.157 1.88.256c.418.095.591.203.696.313c.106.111.193.243.256.39c.067.154.101.377.036.83c-.066.465-.219 1.063-.443 1.939l-.845 3.297Zm-6.832-5.38c-.423.14-.959.39-1.753.762l-.053.025c-1.1.515-1.717.804-2.364 1.01c-.225.07-.453.134-.682.19c-.66.16-1.332.233-2.531.362l-.059.006c-.865.093-1.448.157-1.88.256c-.418.095-.591.203-.696.313a1.328 1.328 0 0 0-.256.39c-.067.154-.101.377-.036.83c.066.465.219 1.063.443 1.939l.845 3.297c.882 3.44 3.798 4.56 5.916 4.92c.348.059.532.088.746.082c.21-.005.486-.045.97-.18c.483-.136.742-.245.929-.352c.19-.109.338-.232.611-.465c1.67-1.425 3.672-3.936 2.787-7.39l-.845-3.296c-.225-.877-.38-1.474-.543-1.91c-.161-.428-.296-.596-.417-.687a1.198 1.198 0 0 0-.393-.196c-.133-.038-.329-.043-.74.094Z"
                            clip-rule="evenodd" />
                        </g>
                      </svg>
                    </button> -->
                  </div>
                  <!-- Image Generation -->
                  <div class="col" id="txt2img">
                    <button id="imgGenBtn" class="btn" data-model-name="128713"
                      onclick="downloadImageModel('dreamshaper-8-turbo-sdxl')" hx-target="#imgstatus"
                      data-bs-toggle="tooltip">
                      <svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <path fill="#FFFFFF"
                          d="M17.29 11.969a1.33 1.33 0 0 1-1.322 1.337a1.33 1.33 0 0 1-1.323-1.337a1.33 1.33 0 0 1 1.323-1.338c.73 0 1.323.599 1.323 1.338Z" />
                        <path fill="#FFFFFF" fill-rule="evenodd"
                          d="M18.132 7.408c-.849-.12-1.942-.12-3.305-.12H9.173c-1.363 0-2.456 0-3.305.12c-.877.125-1.608.393-2.152 1.02c-.543.628-.71 1.397-.716 2.293c-.006.866.139 1.962.319 3.329l.365 2.771c.141 1.069.255 1.933.432 2.61c.185.704.457 1.289.968 1.741c.51.452 1.12.648 1.834.74c.687.088 1.55.088 2.615.088h4.934c1.065 0 1.928 0 2.615-.088c.715-.092 1.323-.288 1.834-.74c.511-.452.783-1.037.968-1.741c.177-.677.291-1.542.432-2.61l.365-2.771c.18-1.367.325-2.463.319-3.33c-.007-.895-.172-1.664-.716-2.291c-.544-.628-1.275-.896-2.152-1.021ZM6.052 8.732c-.726.104-1.094.292-1.34.578c-.248.285-.384.678-.39 1.42c-.005.762.126 1.765.315 3.195l.05.38l.371-.273c.96-.702 2.376-.668 3.288.095l3.384 2.833c.32.268.871.318 1.269.084l.235-.138c1.125-.662 2.634-.592 3.672.19l1.832 1.38c.09-.496.171-1.105.273-1.876l.352-2.675c.189-1.43.32-2.433.314-3.195c-.005-.742-.141-1.135-.388-1.42c-.247-.286-.615-.474-1.342-.578c-.745-.106-1.745-.107-3.172-.107h-5.55c-1.427 0-2.427.001-3.172.107Z"
                          clip-rule="evenodd" />
                        <path fill="#FFFFFF"
                          d="M8.859 2h6.282c.21 0 .37 0 .51.015a2.623 2.623 0 0 1 2.159 1.672H6.19a2.623 2.623 0 0 1 2.159-1.672c.14-.015.3-.015.51-.015ZM6.88 4.5c-1.252 0-2.278.84-2.62 1.954a2.814 2.814 0 0 0-.021.07c.358-.12.73-.2 1.108-.253c.973-.139 2.202-.139 3.629-.139h6.203c1.427 0 2.656 0 3.628.139c.378.053.75.132 1.11.253a2.771 2.771 0 0 0-.021-.07C19.553 5.34 18.527 4.5 17.276 4.5H6.878Z" />
                      </svg>
                    </button>
                  </div>
                </div>
                <div class="w-50 input-group mx-1">
                  <button class="btn btn-secondary" id="upload">
                    <svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                      <path fill="#ffffff" fill-rule="evenodd"
                        d="M11.244 1.955c1.7-.94 3.79-.94 5.49 0c.63.348 1.218.91 2.173 1.825l.093.09l.098.093c.95.91 1.54 1.475 1.906 2.081a5.144 5.144 0 0 1 0 5.337c-.366.607-.955 1.17-1.906 2.08l-.098.095l-7.457 7.14c-.53.506-.96.92-1.34 1.226c-.393.316-.78.561-1.235.692a3.51 3.51 0 0 1-1.937 0c-.454-.13-.841-.376-1.234-.692c-.38-.307-.811-.72-1.34-1.226l-.048-.046c-.529-.507-.96-.92-1.28-1.283c-.33-.376-.592-.753-.733-1.201a3.181 3.181 0 0 1 0-1.907c.14-.448.402-.825.733-1.2c.32-.364.751-.777 1.28-1.284l7.35-7.038l.079-.075c.369-.354.68-.654 1.041-.82a2.402 2.402 0 0 1 2.007 0c.36.166.672.466 1.041.82l.079.075l.08.078c.367.35.683.651.86 1.003a2.213 2.213 0 0 1 0 1.994a2.331 2.331 0 0 1-.391.538c-.142.152-.323.326-.535.529l-7.394 7.08a.75.75 0 0 1-1.038-1.083l7.38-7.067c.23-.22.38-.364.488-.48a.906.906 0 0 0 .15-.191a.712.712 0 0 0 0-.646c-.044-.088-.143-.198-.638-.671c-.492-.471-.61-.57-.71-.617a.902.902 0 0 0-.75 0c-.101.047-.22.146-.711.617L5.47 14.836c-.558.535-.943.904-1.215 1.213c-.267.304-.376.496-.428.66a1.683 1.683 0 0 0 0 1.008c.052.163.16.355.428.659c.272.31.657.678 1.215 1.213c.56.535.945.904 1.269 1.165c.316.255.523.365.707.418c.361.104.747.104 1.108 0c.184-.053.391-.163.707-.418c.324-.261.71-.63 1.269-1.165l7.433-7.117c1.08-1.034 1.507-1.453 1.756-1.866a3.645 3.645 0 0 0 0-3.787c-.249-.413-.676-.832-1.756-1.866c-1.079-1.032-1.518-1.444-1.954-1.685a4.198 4.198 0 0 0-4.039 0c-.437.24-.876.653-1.954 1.685l-5.99 5.735A.75.75 0 0 1 2.99 9.605L8.98 3.87l.093-.09c.955-.914 1.543-1.477 2.172-1.825"
                        clip-rule="evenodd" />
                    </svg>
                  </button>
                  <input type="file" id="file-input" style="display: none;" />
                  <textarea id="message" name="userprompt" class="col form-control shadow-none"
                    placeholder="Type your message..." rows="2"
                    style="outline: none;">Write a complete and executable Hello World in JavaScript.</textarea>
                  <!-- Clear textarea after submit -->
                  <button id="send" class="btn btn-prompt-send btn-success " type="button" hx-post="/chatsubmit"
                    hx-target="#chat" hx-swap="beforeend show:bottom"
                    hx-on::after-request="document.getElementById('message').value=''; textarea.style.height = 'auto'; textarea.style.height = `${Math.min(this.scrollHeight, this.clientHeight * 1)}px`;">
                    <svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                      <path fill="#ffffff" fill-rule="evenodd"
                        d="M12 15.75a.75.75 0 0 0 .75-.75V4.027l1.68 1.961a.75.75 0 1 0 1.14-.976l-3-3.5a.75.75 0 0 0-1.14 0l-3 3.5a.75.75 0 1 0 1.14.976l1.68-1.96V15c0 .414.336.75.75.75"
                        clip-rule="evenodd" />
                      <path fill="#ffffff"
                        d="M16 9c-.702 0-1.053 0-1.306.169a1 1 0 0 0-.275.275c-.169.253-.169.604-.169 1.306V15a2.25 2.25 0 1 1-4.5 0v-4.25c0-.702 0-1.053-.169-1.306a1 1 0 0 0-.275-.275C9.053 9 8.702 9 8 9c-2.828 0-4.243 0-5.121.879C2 10.757 2 12.17 2 14.999v1c0 2.83 0 4.243.879 5.122C3.757 22 5.172 22 8 22h8c2.828 0 4.243 0 5.121-.879C22 20.242 22 18.828 22 16v-1c0-2.829 0-4.243-.879-5.121C20.243 9 18.828 9 16 9" />
                    </svg>
                  </button>
                </div>
    
                <div class="row ms-auto w-25">
                  <div class="col-auto">
                    <button id="rolesBtn" class="btn" hx-post="/tool/websearch" hx-swap="none" data-bs-toggle="tooltip"
                      data-bs-title="Assistant Roles">
                      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
                        <g fill="currentColor">
                          <path
                            d="M6.005 13.368c.029-.296.26-.6.638-.702c.379-.102.73.047.903.29a.75.75 0 0 0 1.22-.873c-.55-.77-1.552-1.123-2.511-.866c-.96.257-1.651 1.064-1.743 2.006a.75.75 0 1 0 1.493.145Zm5.796-1.553c.029-.296.26-.6.638-.702c.379-.102.73.047.903.289a.75.75 0 0 0 1.22-.872c-.55-.77-1.552-1.123-2.511-.866c-.96.257-1.651 1.063-1.743 2.006a.75.75 0 0 0 1.493.145Zm1.399 4.416l.448-.602a.75.75 0 0 1-.885 1.211l-.01-.006a2.06 2.06 0 0 0-.485-.2c-.361-.098-.93-.163-1.686.04c-.756.202-1.215.543-1.48.808a2.064 2.064 0 0 0-.32.416l-.005.01a.75.75 0 0 1-1.372-.607l.689.298l-.689-.298l.001-.001v-.002l.003-.004l.003-.008l.011-.023l.032-.064c.027-.051.065-.118.115-.196c.1-.156.252-.36.469-.578c.436-.439 1.124-.924 2.155-1.2c1.031-.277 1.87-.2 2.467-.038c.297.08.53.18.695.266a2.682 2.682 0 0 1 .257.151l.02.015l.009.005l.003.003h.001l.002.002l-.447.602Z" />
                          <path fill-rule="evenodd"
                            d="m13.252 2.25l.042.02c1.167.547 1.692.791 2.235.963c.193.061.387.116.583.164c.552.134 1.122.197 2.395.334l.045.004c.808.087 1.48.16 2.01.28c.554.127 1.054.328 1.448.743c.23.24.414.521.546.827c.225.52.226 1.064.144 1.64c-.08.554-.253 1.232-.464 2.056l-.856 3.339c-.716 2.793-2.533 4.345-4.357 5.189c-.725 1.574-1.863 2.78-2.804 3.583l-.021.018c-.25.214-.497.425-.82.61c-.335.191-.724.34-1.269.493c-.544.152-.953.227-1.338.236c-.37.009-.687-.045-1.006-.1l-.028-.004c-2.321-.394-6.012-1.714-7.117-6.025l-.856-3.34c-.21-.823-.384-1.5-.464-2.056c-.082-.575-.081-1.118.144-1.639c.132-.306.317-.586.546-.827c.394-.415.894-.616 1.448-.742c.53-.122 1.201-.194 2.01-.28l.045-.005c.52-.056.921-.1 1.253-.14l.625-2.44c.211-.824.385-1.501.582-2.024c.203-.54.466-1.017.92-1.358c.265-.2.565-.348.884-.439c.55-.156 1.084-.066 1.622.113c.516.172 1.132.46 1.873.808Zm6.675 9.997c-.412 1.608-1.26 2.701-2.263 3.45a7.953 7.953 0 0 0-.18-3.207l-.93-3.632a.746.746 0 0 0 .338-.263c.173-.242.525-.39.904-.289c.378.101.608.406.637.702a.75.75 0 1 0 1.493-.145c-.091-.942-.783-1.749-1.742-2.006a2.37 2.37 0 0 0-2.084.416a6.985 6.985 0 0 0-.053-.146c-.203-.54-.466-1.017-.92-1.358a2.698 2.698 0 0 0-.884-.439c-.52-.147-1.026-.075-1.533.085a2.448 2.448 0 0 0-.322-.111c-.96-.257-1.962.096-2.512.866a.748.748 0 0 0-.132.547c-.55.252-.908.4-1.273.516l-.092.03l.434-1.697c.225-.877.38-1.474.543-1.91c.161-.428.296-.596.417-.687c.12-.09.254-.156.393-.196c.133-.038.329-.043.74.094c.422.14.958.39 1.752.762l.053.025c1.1.515 1.717.804 2.364 1.01c.225.07.453.134.682.19c.66.16 1.332.233 2.531.362l.059.006c.865.093 1.448.157 1.88.256c.418.095.591.203.696.313c.106.111.193.243.256.39c.067.154.101.377.036.83c-.066.465-.219 1.063-.443 1.939l-.845 3.297Zm-6.832-5.38c-.423.14-.959.39-1.753.762l-.053.025c-1.1.515-1.717.804-2.364 1.01c-.225.07-.453.134-.682.19c-.66.16-1.332.233-2.531.362l-.059.006c-.865.093-1.448.157-1.88.256c-.418.095-.591.203-.696.313a1.328 1.328 0 0 0-.256.39c-.067.154-.101.377-.036.83c.066.465.219 1.063.443 1.939l.845 3.297c.882 3.44 3.798 4.56 5.916 4.92c.348.059.532.088.746.082c.21-.005.486-.045.97-.18c.483-.136.742-.245.929-.352c.19-.109.338-.232.611-.465c1.67-1.425 3.672-3.936 2.787-7.39l-.845-3.296c-.225-.877-.38-1.474-.543-1.91c-.161-.428-.296-.596-.417-.687a1.198 1.198 0 0 0-.393-.196c-.133-.038-.329-.043-.74.094Z"
                            clip-rule="evenodd" />
                        </g>
                      </svg>
                    </button>
                  </div>
                  <div class="col-auto">
                    <!-- MODELS -->
                    <button class="btn" hx-post="/modelcards" hx-target="#chat" hx-preserve="#chat"
                      hx-swap="innerHTML transition:true" data-bs-toggle="tooltip" data-bs-title="Language Models">
                      <svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <g fill="none">
                          <path stroke="#ffffff" stroke-linecap="round" stroke-width="1.5"
                            d="M9 16c.85.63 1.885 1 3 1s2.15-.37 3-1" />
                          <ellipse cx="15" cy="10.5" fill="#ffffff" rx="1" ry="1.5" />
                          <ellipse cx="9" cy="10.5" fill="#ffffff" rx="1" ry="1.5" />
                          <path stroke="#ffffff" stroke-linecap="round" stroke-width="1.5"
                            d="M22 14c0 3.771 0 5.657-1.172 6.828C19.657 22 17.771 22 14 22m-4 0c-3.771 0-5.657 0-6.828-1.172C2 19.657 2 17.771 2 14m8-12C6.229 2 4.343 2 3.172 3.172C2 4.343 2 6.229 2 10m12-8c3.771 0 5.657 0 6.828 1.172C22 4.343 22 6.229 22 10" />
                        </g>
                      </svg>
                    </button>
                  </div>
                  <!-- SETTINGS -->
                  <div class="col-auto">
                    <button class="btn" data-bs-target="#modal-settings" data-bs-toggle="tooltip" data-bs-title="Settings">
                      <svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                        <g fill="none" stroke="#ffffff" stroke-width="1.5">
                          <circle cx="12" cy="12" r="3" />
                          <path
                            d="M13.765 2.152C13.398 2 12.932 2 12 2c-.932 0-1.398 0-1.765.152a2 2 0 0 0-1.083 1.083c-.092.223-.129.484-.143.863a1.617 1.617 0 0 1-.79 1.353a1.617 1.617 0 0 1-1.567.008c-.336-.178-.579-.276-.82-.308a2 2 0 0 0-1.478.396C4.04 5.79 3.806 6.193 3.34 7c-.466.807-.7 1.21-.751 1.605a2 2 0 0 0 .396 1.479c.148.192.355.353.676.555c.473.297.777.803.777 1.361c0 .558-.304 1.064-.777 1.36c-.321.203-.529.364-.676.556a2 2 0 0 0-.396 1.479c.052.394.285.798.75 1.605c.467.807.7 1.21 1.015 1.453a2 2 0 0 0 1.479.396c.24-.032.483-.13.819-.308a1.617 1.617 0 0 1 1.567.008c.483.28.77.795.79 1.353c.014.38.05.64.143.863a2 2 0 0 0 1.083 1.083C10.602 22 11.068 22 12 22c.932 0 1.398 0 1.765-.152a2 2 0 0 0 1.083-1.083c.092-.223.129-.483.143-.863c.02-.558.307-1.074.79-1.353a1.617 1.617 0 0 1 1.567-.008c.336.178.579.276.819.308a2 2 0 0 0 1.479-.396c.315-.242.548-.646 1.014-1.453c.466-.807.7-1.21.751-1.605a2 2 0 0 0-.396-1.479c-.148-.192-.355-.353-.676-.555A1.617 1.617 0 0 1 19.562 12c0-.558.304-1.064.777-1.36c.321-.203.529-.364.676-.556a2 2 0 0 0 .396-1.479c-.052-.394-.285-.798-.75-1.605c-.467-.807-.7-1.21-1.015-1.453a2 2 0 0 0-1.479-.396c-.24.032-.483.13-.82.308a1.617 1.617 0 0 1-1.566-.008a1.617 1.617 0 0 1-.79-1.353c-.014-.38-.05-.64-.143-.863a2 2 0 0 0-1.083-1.083Z" />
                        </g>
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </form>
        </div>
    
      </div>
    
      <script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
      <script src="js/events.js"></script>
      <script src="js/workflows.js"></script>
      <script src="https://unpkg.com/@antonz/runno@0.6.1/dist/runno.js"></script>
      <script src="https://unpkg.com/@antonz/codapi@0.17.0/dist/engine/wasi.js"></script>
      <script src="https://unpkg.com/@antonz/codapi@0.17.0/dist/snippet.js"></script>
      <script src="https://unpkg.com/@antonz/codapi@0.17.0/dist/settings.js"></script>
      <script src="https://cdn.jsdelivr.net/npm/mermaid/dist/mermaid.min.js"></script>
      <!-- <script src="https://unpkg.com/@antonz/codapi@0.17.0/dist/status.js"></script> -->
      <script>
        const textarea = document.getElementById('message');
        const maxRows = 8;
    
        textarea.addEventListener('input', function () {
          this.style.height = 'auto';
          this.style.height = `${Math.min(this.scrollHeight, this.clientHeight * maxRows)}px`;
        });
    
        const tooltipTriggerList = document.querySelectorAll(
          "[data-bs-toggle='tooltip']"
        );
        const tooltipList = [...tooltipTriggerList].map(
          (tooltipTriggerEl) => new bootstrap.Tooltip(tooltipTriggerEl)
        );
    
        // Function to toggle button state and border
        function toggleButtonState(buttonId) {
          const button = document.getElementById(buttonId);
          button.style.border = button.style.border ? '' : '2px solid purple'; // Toggle border
        }
    
        document.getElementById('imgGenBtn').addEventListener('click', () => {
          toggleButtonState('imgGenBtn');
          downloadImageModel('dreamshaper-8-turbo-sdxl'); // Keep existing functionality
        });
    
        function downloadImageModel(modelName) {
          // Toggle a border on the imgGenBtn
          console.log("Downloading model: ", modelName);
    
          // Toggle a border on the imgGenBtn
          const imgGenBtn = document.getElementById('imgGenBtn');
          imgGenBtn.style.border = imgGenBtn.style.border ? '' : '2px solid purple';
    
          // Fetch the download route
          fetch(`/imgmodel/download?model=${modelName}`, {
            method: 'POST',
            headers: {
              'Content-Type': 'application/json',
            },
            body: JSON.stringify({ modelName }),
          })
            .then(response => {
              if (!response.ok) {
                throw new Error('Failed to download model');
              }
            })
            .catch(error => {
              console.error('Error:', error);
            });
        }
    
        document.getElementById('message').addEventListener('keydown', function (event) {
          if (event.key === 'Enter' || event.key === 'Return') {
            if (!event.shiftKey) {
              event.preventDefault(); // Prevent the default behavior of the Enter key
              document.getElementById('send').click(); // Trigger the click event on the submit button
              textarea.style.height = 'auto';
              textarea.style.height = `${Math.min(this.scrollHeight, this.clientHeight * 1)}px`;
            } else {
              // Insert a new line instead of submitting the form
              var messageInput = document.getElementById('message');
              messageInput.value += '\n';
            }
          }
        });
    
      </script>
    </body>
    
    </html>

    [File Ends] public/templates/index.html

    [File Begins] public/templates/model.html
    <div id="models-container" class="fade-it">
      <div id="content" class="mt-3 mb-4 pb-3">
        <button class="btn btn-danger btn-lg mb-2 w-100" style="background-color: var(--et-red);" hx-get="/"
          hx-target="#models-container" hx-select="#chat" hx-preserve="#chat" hx-swap="delete transition:true">Back</button>
      </div>
      <div class="row">
        {{range .models}}
        <div class="col-6 gy-3">
          <div class="card h-100" data-model-name="{{.Name}}">
            <!-- If .Downloaded is false, add button to call download link -->
            {{if not .Downloaded}}
            {{if not (or (eq .Name "openai-gpt") (eq .Name "google-gemini-1.5") (eq .Name "anthropic-claude-opus"))}}
            <div class="card-header col">{{.Name}}</div>
            <div class="row h-100">
              <div name="progress-download" id="progress-download-{{.Name}}"></div>
              <button id="btn-download-{{.Name}}" class="btn fw-medium" style="min-height: 228px;"
                data-model-name="{{.Name}}" hx-post="/model/download?model={{.Name}}" hx-trigger="click">
                <svg width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path fill="currentColor" fill-rule="evenodd"
                    d="M12 15.25a.75.75 0 0 1 .75.75v4.19l.72-.72a.75.75 0 1 1 1.06 1.06l-2 2a.75.75 0 0 1-1.06 0l-2-2a.75.75 0 1 1 1.06-1.06l.72.72V16a.75.75 0 0 1 .75-.75"
                    clip-rule="evenodd" />
                  <path fill="currentColor"
                    d="M12.226 3.5c-2.75 0-4.964 2.2-4.964 4.897c0 .462.065.909.185 1.331c.497.144.963.36 1.383.64a.75.75 0 1 1-.827 1.25a3.54 3.54 0 0 0-1.967-.589c-1.961 0-3.536 1.57-3.536 3.486C2.5 16.43 4.075 18 6.036 18a.75.75 0 0 1 0 1.5C3.263 19.5 1 17.276 1 14.515c0-2.705 2.17-4.893 4.864-4.983a6.366 6.366 0 0 1-.102-1.135C5.762 4.856 8.664 2 12.226 2c3.158 0 5.796 2.244 6.355 5.221c2.3.977 3.919 3.238 3.919 5.882c0 3.074-2.188 5.631-5.093 6.253a.75.75 0 0 1-.314-1.467c2.24-.48 3.907-2.446 3.907-4.786c0-2.137-1.39-3.962-3.338-4.628a5.018 5.018 0 0 0-1.626-.27c-.583 0-1.14.1-1.658.28a.75.75 0 0 1-.494-1.416a6.517 6.517 0 0 1 3.024-.305A4.962 4.962 0 0 0 12.226 3.5" />
                </svg>
              </button>
            </div>
            {{else}}
            <div id="card-header-{{.Name}}" class="card-header" onclick="selectModel('{{.Name}}')">{{.Name}}</div>
            <div class="card-body tab-content mh-100" id="js-tabs-content-{{.Name}}">
              <!-- Tabs -->
              <ul class="nav nav-underline nav-fill mb-3 justify-content-center" id="js-tabs-{{.Name}}" role="tablist">
                <li class="nav-item">
                  <a class="nav-link active" id="settings-tab-{{.Name}}" data-bs-toggle="tab"
                    data-bs-target="#settings-tab-pane-{{.Name}}" role="tab" aria-controls="settings-tab-pane-{{.Name}}"
                    aria-selected="true" style="background-color: unset;">Settings</a>
                </li>
                <li class="nav-item">
                  <a class="nav-link" id="template-tab-{{.Name}}" data-bs-toggle="tab"
                    data-bs-target="#template-tab-pane-{{.Name}}" role="tab" aria-controls="template-tab-pane-{{.Name}}"
                    aria-selected="false" style="background-color: unset;">Prompt Template</a>
                </li>
              </ul>
              <!-- Settings Tab Content -->
              <div class="tab-pane fade show active" id="settings-tab-pane-{{.Name}}" role="tabpanel"
                aria-labelledby="settings-tab-{{.Name}}">
                <p><a href="{{.GGUFInfo}}" target="_blank">Model Info</a></p>
                <p><strong>Context:</strong> {{.Options.CtxSize}}</p>
                <p><strong>Temperature:</strong> {{.Options.Temp}}</p>
                <p><strong>Repetition Penalty:</strong> {{.Options.RepeatPenalty}}</p>
              </div>
              <!-- Prompt Template Tab Content -->
              <div class="tab-pane fade" id="template-tab-pane-{{.Name}}" role="tabpanel"
                aria-labelledby="template-tab-{{.Name}}">
                <p><strong>Prompt Template:</strong> {{.Options.Prompt}}</p>
              </div>
            </div>
            {{end}}
            {{else}}
            <div class="card-header" onclick="selectModel('{{.Name}}')">{{.Name}}</div>
            <div class="card-body tab-content mh-100" id="js-tabs-content-{{.Name}}">
              <!-- Tabs -->
              <ul class="nav nav-underline nav-fill mb-3 justify-content-center" id="js-tabs-{{.Name}}" role="tablist">
                <li class="nav-item">
                  <a class="nav-link active" id="settings-tab-{{.Name}}" data-bs-toggle="tab"
                    data-bs-target="#settings-tab-pane-{{.Name}}" role="tab" aria-controls="settings-tab-pane-{{.Name}}"
                    aria-selected="true" style="background-color: unset;">Settings</a>
                </li>
                <li class="nav-item">
                  <a class="nav-link" id="template-tab-{{.Name}}" data-bs-toggle="tab"
                    data-bs-target="#template-tab-pane-{{.Name}}" role="tab" aria-controls="template-tab-pane-{{.Name}}"
                    aria-selected="false" style="background-color: unset;">Prompt Template</a>
                </li>
              </ul>
              <!-- Settings Tab Content -->
              <div class="tab-pane fade show active" id="settings-tab-pane-{{.Name}}" role="tabpanel"
                aria-labelledby="settings-tab-{{.Name}}">
                <p><a href="{{.GGUFInfo}}" target="_blank">Model Info</a></p>
                <p><strong>Context:</strong> {{.Options.CtxSize}}</p>
                <p><strong>Temperature:</strong> {{.Options.Temp}}</p>
                <p><strong>Repetition Penalty:</strong> {{.Options.RepeatPenalty}}</p>
              </div>
              <!-- Prompt Template Tab Content -->
              <div class="tab-pane fade" id="template-tab-pane-{{.Name}}" role="tabpanel"
                aria-labelledby="template-tab-{{.Name}}">
                <p><strong>Prompt Template:</strong> {{.Options.Prompt}}</p>
              </div>
            </div>
            {{end}}
          </div>
        </div>
        {{end}}
      </div>
    </div>
    
    <script>
      async function downloadModel(modelName) {
    
        // Append 
    
        console.log("Downloading model: ", modelName);
    
        // Fetch the download route
        fetch(`/model/download?model=${modelName}`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ modelName }),
        })
          .then(response => {
            if (!response.ok) {
              throw new Error('Failed to download model');
            }
          })
          .catch(error => {
            console.error('Error:', error);
          });
      }
    
      async function selectModel(modelName) {
        // Deselect all cards
        const cards = document.querySelectorAll('.card');
        cards.forEach(card => card.classList.remove('card-selected'));
    
        // Get the clicked card element
        const cardElement = document.querySelector(`div[data-model-name="${modelName}"]`);
    
        // Select the clicked card
        cardElement.classList.add('card-selected');
    
        // Send request to update model selection in the backend
        fetch(`/model/select`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ modelName, action: 'add' }), // Always 'add' as we deselect all others
        })
          .then(response => {
            if (!response.ok) {
              throw new Error('Failed to update model selection');
            }
          })
          .catch(error => {
            console.error('Error:', error);
          });
      }
    </script>

    [File Ends] public/templates/model.html

    [File Begins] public/templates/shell.html
    
    <div
      x-data="loadterm()"
    >
      <div x-disclosure>
        <div
          x-disclosure:button
        >
          <span class="shellButton">CloudBox WebShell</span>
      </div>
    
        <div x-disclosure:panel x-collapse>
          <div id="terminal" style="width: 100%; overflow: hidden;"></div>
        </div>
      </div>
    </div>
    
    <script type="text/javascript">
      // Add an event listener to the shell id to load the terminal
      document.getElementById("shell").addEventListener("click", loadterm);
    
    
      function loadterm() {
        console.log("loadterm");
        var conn;
        var term = new Terminal();
        var fitAddon = new FitAddon.FitAddon();
        term.loadAddon(fitAddon);
        term.open(document.getElementById("terminal"));
        fitAddon.fit();
        term.writeln("Hello from \x1B[1;3;31mCloudBox!\x1B[0m");
    
        function sendMessage(message) {
          if (!conn) {
            return;
          }
          if (!message) {
            return;
          }
          termConn.send(message);
        }
    
        if (window["WebSocket"]) {
          termConn = new WebSocket("wss://" + document.location.host + "/host/shell");
          termConn.onopen = function (evt) {
            term.writeln("Connected to CloudBox host terminal.");
          };
          termConn.onclose = function (evt) {
            term.writeln("Connection closed.");
          };
          termConn.onmessage = function (evt) {
            term.write("\r" + evt.data);
            //term.write(evt.data);
          };
          termConn.onerror = function (evt) {
            term.writeln("ERROR: " + evt.data);
          };
        } else {
          term.writeln("Your browser does not support WebSockets.");
        }
    
        var inputBuffer = "";
    
        term.onKey(function (keyEvent) {
          const key = keyEvent.key;
    
          if (key === "\r") {
            // Enter key
            term.write("\n"); // Add a newline and a prompt before the next command
            sendMessage(inputBuffer);
            inputBuffer = "";
          } else if (key === "\u007F" || key === "\b") {
            // Backspace or Delete key
            if (inputBuffer.length > 0) {
              inputBuffer = inputBuffer.slice(0, inputBuffer.length - 1);
              term.write("\b \b"); // Move cursor back, write a space to erase, then move cursor back again
            }
          } else if (key.length === 1) {
            inputBuffer += key;
            term.write(key);
          }
        });
      }
    </script>

    [File Ends] public/templates/shell.html

    [File Begins] public/uploads/uploads_go_here

    [File Ends] public/uploads/uploads_go_here

[File Begins] utils.go
package main

import (
	"embed"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sync/atomic"

	"github.com/gofiber/fiber/v2/log"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/spf13/afero"
)

var (
	LocalFs        = new(afero.OsFs)
	MemFs          = afero.NewMemMapFs()
	messageCounter int64
)

func InitServer(configPath string) (string, error) {

	// WEB FILES
	webPath := filepath.Join(configPath, "web")
	err := os.MkdirAll(webPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", webPath, err)
	}
	err = CopyFiles(embedfs, "public", webPath)
	if err != nil {
		return "", fmt.Errorf("failed to copy files: %v", err)
	}

	// GGUF FILES
	ggufPath := filepath.Join(configPath, "gguf")
	err = os.MkdirAll(ggufPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", ggufPath, err)
	}
	err = CopyFiles(embedfs, "pkg/llm/local/bin", ggufPath)
	if err != nil {
		return "", fmt.Errorf("failed to copy files: %v", err)
	}

	files, err := os.ReadDir(ggufPath)
	if err != nil {
		return "", fmt.Errorf("failed to read directory %s: %v", ggufPath, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			err = os.Chmod(filepath.Join(ggufPath, file.Name()), 0755)
			if err != nil {
				return "", fmt.Errorf("failed to set executable permission on file %s: %v", file.Name(), err)
			}
		}
	}

	// IMG GEN
	imgGenPath := filepath.Join(configPath, "sd")
	err = os.MkdirAll(imgGenPath, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", imgGenPath, err)
	}

	err = CopyFiles(embedfs, "pkg/sd/sdcpp/build/bin", imgGenPath)
	if err != nil {
		return "", fmt.Errorf("failed to copy files: %v", err)
	}

	files, err = os.ReadDir(imgGenPath)
	if err != nil {
		return "", fmt.Errorf("failed to read directory %s: %v", imgGenPath, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			err = os.Chmod(filepath.Join(imgGenPath, file.Name()), 0755)
			if err != nil {
				return "", fmt.Errorf("failed to set executable permission on file %s: %v", file.Name(), err)
			}
		}
	}

	return configPath, nil
}

func GetServerInfo() {
	// Basic OS and Architecture information
	fmt.Println("OS:", runtime.GOOS)
	fmt.Println("Arch:", runtime.GOARCH)

	// CPU information
	cpuInfos, _ := cpu.Info()
	for _, ci := range cpuInfos {
		fmt.Printf("CPU: %+v\n", ci)
	}

	// Memory information
	vmStat, _ := mem.VirtualMemory()
	fmt.Printf("Virtual Memory: %+v\n", vmStat)
}

func EnsureDataPath(config *AppConfig) error {
	if _, err := os.Stat(config.DataPath); os.IsNotExist(err) {
		return LocalFs.MkdirAll(config.DataPath, os.ModePerm)
	}
	return nil
}

func CopyFiles(fsys embed.FS, srcDir, destDir string) error {
	fileEntries, err := fsys.ReadDir(srcDir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", srcDir, err)
	}

	for _, entry := range fileEntries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())

		if entry.IsDir() {
			// Create the directory and copy its contents
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory %s: %v", destPath, err)
			}
			if err := CopyFiles(fsys, srcPath, destPath); err != nil {
				return err
			}
		} else {
			// Copy the file
			fileData, err := fsys.ReadFile(srcPath)
			if err != nil {
				log.Errorf("failed to read file %s: %v", srcPath, err)
				continue // Skip to the next file
			}
			if err := os.WriteFile(destPath, fileData, 0755); err != nil {
				return fmt.Errorf("failed to write file %s: %v", destPath, err)
			}
		}
	}
	return nil
}

// Increments and returns a counter that gets appended to the id for frontend chat elements
func IncrementTurn() int64 {
	return atomic.AddInt64(&messageCounter, 1)
}

// findURLInText searches for a URL in a given text and returns it if found.
// It returns nil if no valid URL is found.
func URLParse(text string) *url.URL {
	// Define a regular expression for finding URLs
	// This is a simple regex for demonstration; it might not cover all URL cases
	re := regexp.MustCompile(`https?://[^\s]+`)

	// Find a URL using the regex
	found := re.FindString(text)
	if found == "" {
		// No URL found
		return nil
	}

	// Parse the URL to validate it and return *url.URL
	parsedURL, err := url.Parse(found)
	if err != nil {
		// The URL is not valid
		return nil
	}

	return parsedURL
}

[File Ends] utils.go

<-- File Content Ends


Code Review and Analysis of the Eternal Application
This application, named “Eternal,” appears to be an AI-powered chatbot with a web frontend. It leverages various language and image models, including OpenAI, Anthropic, Google Gemini, and local GGUF models. It also incorporates tools like web retrieval, web search, and image generation.

Strengths:

Modular Design: The code is well-organized into packages, separating concerns like embeddings, LLMs, web utilities, and stable diffusion.
Embedded Filesystem: Using embed.FS for static assets and binaries improves portability and simplifies deployment.
Configuration Management: The config.yml file provides a centralized location for configuring various aspects of the application.
WebSocket Integration: Real-time communication via WebSockets enhances the user experience for chat interactions.
Tool Workflow: The application implements a tool workflow that enriches user prompts with additional context from web retrieval, search, and memory.
Database Integration: SQLite is used for persisting model data, chat history, and potentially embeddings, providing data management capabilities.
Code Highlighting: The frontend uses highlight.js for code highlighting, improving the readability of code snippets in responses.
Weaknesses and Areas for Improvement:

Error Handling: Error handling is inconsistent throughout the code. Some functions return errors, while others simply log them or panic. A consistent error handling strategy should be implemented, potentially using custom error types for better context.
Testing: The provided test suite is limited and doesn’t cover all functionalities. More comprehensive unit and integration tests are needed to ensure code quality and prevent regressions.
Documentation: While the repository documentation provides a basic overview, the code itself lacks detailed comments explaining the logic and purpose of various functions and components.
Hardcoded Values: Several hardcoded values, like model names and paths, should be moved to the configuration file for easier customization and maintenance.
Concurrency Management: While the application uses goroutines for some tasks, there’s room for improvement in concurrency management, especially for long-running operations like model downloads.
Security: The application allows CORS from all origins, which is a security risk in production environments. A more restrictive CORS policy should be implemented.
User Interface: The frontend UI is basic and could benefit from improvements in design and user experience.
Abstraction: Some parts of the code, like the tool workflow, could benefit from further abstraction to improve maintainability and extensibility.
Specific Recommendations:

Implement a Consistent Error Handling Strategy: Define custom error types for different scenarios and handle errors gracefully, providing informative messages to the user.
Expand Test Coverage: Write unit tests for individual functions and components, and integration tests for end-to-end workflows.
Improve Code Documentation: Add detailed comments explaining the purpose, logic, and expected behavior of functions, classes, and significant code blocks.
Externalize Hardcoded Values: Move hardcoded values like model names, paths, and API keys to the configuration file.
Enhance Concurrency Management: Use channels and synchronization primitives to manage concurrent operations effectively, especially for downloads and long-running tasks.
Strengthen Security: Implement a more restrictive CORS policy, allowing only trusted origins. Consider adding authentication and authorization mechanisms for sensitive operations.
Refine User Interface: Improve the frontend UI with a more modern design and user-friendly interactions. Consider using a frontend framework like React or Vue.js for better structure and maintainability.
Increase Abstraction: Abstract the tool workflow logic into separate components or classes for each tool, making it easier to add, remove, or modify tools.
Optimize Performance: Profile the application to identify performance bottlenecks and optimize critical code paths. Consider using caching mechanisms for frequently accessed data.
Implement Logging: Integrate a logging framework like logrus or zap to capture application events and errors for debugging and monitoring.
Overall Quality:

The code demonstrates a good understanding of Golang concepts and a decent level of organization. However, there are areas for improvement in error handling, testing, documentation, and security. By addressing these weaknesses, the application can become more robust, maintainable, and secure.