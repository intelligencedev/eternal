# Detect the operating system
OS := $(shell uname -s)

# Path variables
LLAMA_DIR := pkg/llm/local/gguf
SD_DIR := pkg/sd/sdcpp
ARTIFACTS_DIR := pkg/llm/local/bin
LLAMA_BUILD_DIR := $(LLAMA_DIR)/build
SD_BUILD_DIR := $(SD_DIR)/build

# Define the location of binaries based on OS and build type
LLAMA_BINARIES_DIR := $(LLAMA_DIR) # Default for macOS
ifeq ($(OS),Linux)
    LLAMA_BINARIES_DIR := $(LLAMA_BUILD_DIR)/bin
endif

# Build commands
ifeq ($(OS),Darwin) # macOS
    LLAMA_BUILD_CMD = make -C $(LLAMA_DIR)
    SD_BUILD_CMD = cmake -S $(SD_DIR) -B $(SD_BUILD_DIR) -DSD_METAL=ON && cmake --build $(SD_BUILD_DIR) --config Release
else ifeq ($(OS),Linux)
    LLAMA_BUILD_CMD = cmake -S $(LLAMA_DIR) -B $(LLAMA_BUILD_DIR) -DLLAMA_CUBLAS=ON -DCMAKE_CUDA_COMPILER:PATH=/usr/local/cuda/bin/nvcc && cmake --build $(LLAMA_BUILD_DIR) --config Release
    SD_BUILD_CMD = cmake -S $(SD_DIR) -B $(SD_BUILD_DIR) -DSD_CUBLAS=ON -DCMAKE_CUDA_COMPILER:PATH=/usr/local/cuda/bin/nvcc && cmake --build $(SD_BUILD_DIR) --config Release
else
$(error Unsupported operating system)
endif

# Adjusted copy command for SD binary
COPY_SD_CMD = cp $(SD_BUILD_DIR)/bin/sd $(ARTIFACTS_DIR)

# Copy all llama artifacts to the artifacts directory
COPY_LLAMA_CMD = cp -R $(LLAMA_BINARIES_DIR) $(ARTIFACTS_DIR)

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
    $(SD_BUILD_CMD)

# Copy artifacts
.PHONY: copy-artifacts
copy-artifacts:
    mkdir -p $(ARTIFACTS_DIR)
    $(COPY_LLAMA_CMD)
    $(COPY_SD_CMD)

# Build the eternal application
.PHONY: eternal
eternal: #copy-artifacts
    go build -tags netgo -ldflags="-s -w" -trimpath -o eternal

# Default target
.PHONY: all
all: init-submodules deps llama sd copy-artifacts eternal

# Clean up build artifacts
.PHONY: clean
clean:
    rm -rf $(LLAMA_BUILD_DIR)
    rm -rf $(SD_BUILD_DIR)
    rm -rf $(ARTIFACTS_DIR)
    rm -f eternal