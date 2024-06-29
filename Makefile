# Detect the operating system
OS := $(shell uname -s)

# Path variables
LLAMA_BUILD_TARGETS := llama-baby-llama llama-batched llama-batched-bench llama-bench llama-benchmark-matmult llama-cli llama-convert-llama2c-to-ggml llama-cvector-generator llama-embedding llama-eval-callback llama-export-lora llama-finetune llama-gbnf-validator llama-gguf llama-gguf-split llama-gritlm llama-imatrix llama-infill llama-llava-cli llama-lookahead llama-lookup llama-lookup-create llama-lookup-merge llama-lookup-stats llama-parallel llama-passkey llama-perplexity llama-q8dot llama-quantize llama-quantize-stats llama-retrieval llama-save-load-state llama-server llama-simple llama-speculative llama-tokenize llama-train-text-from-scratch llama-vdot

# Directory variables
LLAMA_DIR := pkg/llm/local/gguf
ARTIFACTS_DIR := pkg/llm/local/bin
LLAMA_BUILD_DIR := $(LLAMA_DIR)/build

SD_DIR := pkg/sd/sdcpp
SD_BUILD_DIR := $(SD_DIR)/build

# Define the location of binaries based on OS and build type
LLAMA_BINARIES_DIR := $(LLAMA_DIR) # Default for macOS
ifeq ($(OS),Linux)
	LLAMA_BINARIES_DIR := $(LLAMA_BUILD_DIR)/bin
endif

# Build commands
ifeq ($(OS),Darwin) # macOS
	LLAMA_BUILD_CMD = make -C $(LLAMA_DIR)
	COPY_LLAMA_CMD = cp $(addprefix $(LLAMA_DIR)/, $(LLAMA_BUILD_TARGETS)) $(ARTIFACTS_DIR)
	SD_BUILD_CMD = cmake -S $(SD_DIR) -B $(SD_BUILD_DIR) -DSD_METAL=ON -DSD_FLASH_ATTN=ON && cmake --build $(SD_BUILD_DIR) --config Release
else ifeq ($(OS),Linux)
	LLAMA_BUILD_CMD = cmake -S $(LLAMA_DIR) -B $(LLAMA_BUILD_DIR) -DLLAMA_CUBLAS=ON -DCMAKE_CUDA_COMPILER:PATH=/usr/local/cuda/bin/nvcc && cmake --build $(LLAMA_BUILD_DIR) --config Release
    FILTERED_LLAMA_BUILD_TARGETS := $(filter-out libllava.a ggml-metal.metal ggml-common.h benchmark-matmult,$(LLAMA_BUILD_TARGETS))
    COPY_LLAMA_CMD = cp $(addprefix $(LLAMA_BUILD_DIR)/bin/,$(FILTERED_LLAMA_BUILD_TARGETS)) $(ARTIFACTS_DIR)
	SD_BUILD_CMD = cmake -S $(SD_DIR) -B $(SD_BUILD_DIR) -DSD_FLASH_ATTN=ON -DSD_CUBLAS=ON -DCMAKE_CUDA_COMPILER:PATH=/usr/local/cuda/bin/nvcc && cmake --build $(SD_BUILD_DIR) --config Release
else
	$(error Unsupported operating system)
endif

$(shell mkdir -p $(ARTIFACTS_DIR))

# Adjusted copy command for SD binary
COPY_SD_CMD = cp $(SD_BUILD_DIR)/bin/sd $(ARTIFACTS_DIR)

# Initialize and update git submodules
.PHONY: init-submodules
init-submodules:
	git submodule update --init --recursive

# Build dependencies
.PHONY: deps
deps: init-submodules llama sd

# Build llama
.PHONY: llama
llama:
	$(LLAMA_BUILD_CMD)

# Build stable-diffusion
.PHONY: sd
sd:
	rm -rf $(SD_BUILD_DIR)
	$(SD_BUILD_CMD)

# Define a macro for copying a build target
define COPY_BUILD_TARGET
$(ARTIFACTS_DIR)/$(1): | $(ARTIFACTS_DIR)
	cp $(LLAMA_DIR)/$(1) $(ARTIFACTS_DIR)/$(1)

.PHONY: $(ARTIFACTS_DIR)/$(1)
endef

# Generate the copy targets
$(foreach target,$(LLAMA_BUILD_TARGETS),$(eval $(call COPY_BUILD_TARGET,$(target))))

# Copy llama artifacts
.PHONY: copy-llama-artifacts
copy-llama-artifacts:
	$(COPY_LLAMA_CMD)

# Copy sd artifacts
.PHONY: copy-sd-artifacts
copy-sd-artifacts:
	$(COPY_SD_CMD)

# Build the eternal application
.PHONY: eternal
eternal: #copy-artifacts
	go build -tags netgo -ldflags="-s -w" -trimpath -o eternal

# Default target
.PHONY: all
all: init-submodules deps llama sd copy-llama-artifacts copy-sd-artifacts eternal

# Clean up build artifacts
.PHONY: clean
clean:
	rm -rf $(LLAMA_BUILD_DIR)
	rm -rf $(SD_BUILD_DIR)
	rm -rf $(ARTIFACTS_DIR)
	rm -f eternal