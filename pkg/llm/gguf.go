package llm

import (
	"bufio"
	"bytes"
	_ "embed"

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
	cmdPath = filepath.Join(cmdPath, "gguf/main")

	ctxSize := fmt.Sprintf("%d", options.CtxSize)
	temp := fmt.Sprintf("%f", options.Temp)
	repeatPenalty := fmt.Sprintf("%f", options.RepeatPenalty)
	topP := fmt.Sprintf("%f", options.TopP)
	topK := fmt.Sprintf("%d", options.TopK)

	cmdArgs := []string{
		"--no-display-prompt",
		"-m", options.Model,
		"-p", options.Prompt,
		"-c", ctxSize,
		"--repeat-penalty", repeatPenalty,
		"--top-p", topP,
		"--top-k", topK,
		//"--n-gpu-layers", "-1",
		"--reverse-prompt", "<|eot_id|>",
		"--multiline-input",
		"--temp", temp,
		// "--mlock",
		"--seed", "-1",
		//"--ignore-eos",
		//"--no-mmap",
		"--simple-io",
		// "--rope-scaling", "linear",
		// "--rope-scale", "8",
		// "--rope-freq-base", "1",
		// "--rope-freq-scale", "1",
		// "--yarn-orig-ctx", "4096",
		//"--yarn-ext-factor", "1.0",
		//"--yarn-attn-factor", "1.0",
		//"--yarn-beta-slow", "0.0",
		//"--yarn-beta-fast", "0.0",
		//"-cml",
		//"--keep", "2048",
		//"--prompt-cache", "cache",
		//"--verbose-prompt",
		//"--in-prefix", "\n<|start_header_id|>user<|end_header_id|>\n\n",
		//"--in-prefix-bos",
		//"--in-suffix", "<|eot_id|><|start_header_id|>assistant<|end_header_id|>\n\n",
		//"--grammar-file", "./json.gbnf",
		//"--override-kv", "llama.expert_used_count=int:3", // mixtral only
		//"--override-kv", "tokenizer.ggml.pre=str:llama3",
	}

	return exec.Command(cmdPath, cmdArgs...)
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
