<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" height="200px" srcset="./public/img/eternal.png">
    <img alt="logo" height="200px" src="./public/img/eternal.png">
  </picture>
</div>

# Eternal

Eternal is an experimental platform for machine learning workflows.

<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./public/img/chat.png">
    <img alt="logo" src="./public/img/chat.png" style="width:50%; height:50%"
  </picture>
</div>

<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="./public/img/models.png">
    <img alt="logo" src="./public/img/models.png" style="width:50%; height:50%"
  </picture>
</div>

NOTE: This app is a work in progress and not stable. Please consider this repo for your reference. We
welcome contributors and constructive feedback. You are also welcome to use it as reference for your own projects.

Eternal integrates various projects such as `llama.cpp`, `stable diffusion.cpp` and `codapi` among many other projects whose
developers were kind enough to share with the world. All credit belongs to the respective contributors of all dependencies this
repo relies on. Thank you for sharing your projects with the world.

The Eternal frontend is rendered with the legendary `HTMX` framework.

IMPORTANT:

Configure the quant level of the models in your `config.yml` appropriately for your system specs. If a local fails to run, investigate the reason by viewing the generated `main.log` file. The most common reason is insufficient RAM or incorrect prompt template. We will implement more robust error handling and logging in a future commit.

## Features

- Language model catalog for easy download and configuration.
- Text generation using local language models, OpenAI GPT-4 an Google Gemini 1.5 Pro.
- Web retrieval that fetches URL content for LLM to reference.
- Web Search to automatically retrieve top results for a user's prompt for LLM to reference. _Requires Chrome browser installation._
- Image generation using Stable Diffusion backend.

## Showcase

### Web Retrieval - Search

Prompt for a URL and Eternal automatically fetches and sanitizes the page content for reference.
The search button can be enabled to fetch one result from a popular search engine. (Limited to 1 while context management features are matured.) _Requires Chrome browser installation._

<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" height="400px" srcset="./public/img/web.png">
    <img alt="logo" height="400px" src="./public/img/web.png" style="width:50%; height:50%">
  </picture>
</div>

### Code Execution

Execute and edit LLM generated code in the chat view in a secure sandbox. For now, JavaScript is implemented via WASM. More languages coming soon!

<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" height="400px" srcset="./public/img/code_fixed.png">
    <img alt="logo" height="400px" src="./public/img/code_fixed.png" style="width:50%; height:50%">
  </picture>
</div>

### Image Generation

Embedded Stable Diffusion for easy high quality image generation.
<div align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" height="400px" srcset="./public/img/sd.png">
    <img alt="logo" height="400px" src="./public/img/imggen1.png" style="width:50%; height:50%">
  </picture>
  <picture>
    <source media="(prefers-color-scheme: dark)" height="400px" srcset="./public/img/sd.png">
    <img alt="logo" height="400px" src="./public/img/imggen2.png" style="width:50%; height:50%">
  </picture>
</div>

## Configuration

Rename the provided `.config.yml` file to `config.yml` and place it in the same path as the application binary. Modify the contents for your environment and use case.

## Build

Eternal currently supports building on Linux or Windows WSL using CUDA (nVidia GPU required) or MacOS/Metal (M-series Max required).

To build the application:

```
$ git clone https://github.com/intelligencedev/eternal.git
$ cd eternal
$ git submodule update --init --recursive
$ make all
```

Please submit an issue if you encounter any issues with the build process.

## Disclaimer

This README is a high-level overview of the Eternal application. Detailed setup instructions and a complete list of features, dependencies, and configurations should be consulted in the actual application documentation.
