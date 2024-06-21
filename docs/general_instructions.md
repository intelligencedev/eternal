# Eternal

Eternal has been tested on MacOS using a modern M-Series Macbook, Ubuntu Linux using CUDA, and Windows WSL2 Ubuntu using CUDA. 

Eternal has the following general requirements:

- [Go 1.21.1](https://go.dev/doc/install)
- [Python 3](https://www.python.org/downloads/)
- [Google Chrome Browser](https://www.google.com/chrome/)
- `cmake`

There may be other requirements not documented here, please file an issue if this is the case or you have trouble getting started.

See the `cuda.md` doc in this repo for general instructions on how to setup CUDA drivers and other requirements for CUDA inference in a Linux or Windows WSL2 environment.

Note that image generation is significantly faster on a CUDA system. On a modern Apple M3 Series Macbook Pro, the average image generation time is about 30 seconds to 1 minute.

## Important!

Upon first execution of the application, Eternal will create a folder to stage all necessary files for the application to function. The location of all files is located as follows:

MacOS: `/Users/$USER.eternal`

Linux / Windows WSL2: `/home/$USER/.eternal`

Replace the `$USER` variable with your actual username. Please note this can be changed in the application configuration explained in the next steps.

## Getting Started

First, rename the provided `.config.yml` to `config.yml` and modify the contents for your environment. Multi-node workflows are not currently implemented so Eternal will only use the primary node for each configuration.

Be mindful of the preconfigured models. The provided config defaults to `Q8_0` quants which requires more compute power. If you want to download a lower quant model, we provide the link to the GGUF source that will list the available quants for the particular model.

It is recommended the image generation model not be changed since the preconfigured image generation settings are tuned to work with that particular model. We will expose these configurations via the UI in a future update to make it easier to experiment with different image models.

1. Download the application binary and drop it in your desired location.
2. Create a `config.yml` and put it in the same path as the `eternal` binary.
3. Open a terminal window, change into the binary path and run the application to start the service: `$ ./eternal`
4. Open your desired web browser and navigate to the configured host and port in the application configuration, by default: `http://localhost:8080` 
5. Click the models button on the bottom right of the interface and select one of the preconfigured models. Automatic download will occur for local models. Once the download completes, refresh the page, open the models view, and select the model. Monitor the terminal window in case there are issues with the download. If for any reason the download is interrupted, delete the model folder that was created in the application configuration path: `config_path/models/<model_name>` and retry the download.

In general, if a bug is encountered or there are issues, the best thing to do is quit the application in the terminal using `CTRL+C`, then delete the entire application configuration folder. In order to avoid having to download models again, you may opt to delete all the contents of the application configuration folder except the `models` subfolder.

If you encounter a bug, please open an issue.

## LLM Parameters

### Creative Writing (e.g., story generation, poetry)
- **TopK Range:** 40-100
- **TopP Range:** 0.8-0.95
- **Reason:** Higher TopP values encourage more diversity and creativity, allowing for unexpected and novel combinations of words.

### Conversational Agents (e.g., chatbots)
- **TopK Range:** 10-50
- **TopP Range:** 0.7-0.9
- **Reason:** A moderate TopP range ensures responses are relevant and contextually appropriate while maintaining some level of diversity.

### Technical Writing (e.g., code generation, documentation)
- **TopK Range:** 1-10
- **TopP Range:** 0.5-0.7
- **Reason:** Lower TopP values ensure precision and correctness, as the output needs to be highly accurate and specific.

### Summarization
- **TopK Range:** 5-20
- **TopP Range:** 0.6-0.8
- **Reason:** A balance is needed to ensure the summary is concise and relevant without losing important details.

### Translation
- **TopK Range:** 5-30
- **TopP Range:** 0.6-0.85
- **Reason:** Ensures the translation is accurate and contextually appropriate, with some flexibility for natural language variations.

### Question Answering
- **TopK Range:** 1-5
- **TopP Range:** 0.3-0.6
- **Reason:** Requires highly accurate and specific answers, so a very low TopP range is ideal to avoid irrelevant or incorrect responses.


## Tool Configuration

The tools are highly experimental and may not work as intended in some cases. It is recommended that the memory tool be disabled when enabling the web tools to prevent Eternal from running that embeddings workflow for memory which can be time consuming on large documents. This will be improved in future updates.

Deploying public web tools and scraping the public internet carries risk of your public IP and/or client getting banned from services. Always monitor the CLI logs as the tools execute their workflows to ensure your client does not spam. If this occurs, quit the application via the CLI using CTRL+C . Disable the web tools and open an Issue so we can implement a diligent fix.

The web retrieval tools require a Google Chrome installation. Search works without requiring any APIs or paid services and runs entirely local by making calls to a popular and private search engine. We ask that you give the [search platform your support](https://duckduckgo.com/donations) for providing a great service.


# Disclaimer

Eternal is provided as-is and its primary purpose is personal use to experiment with machine learning models and interesting workflows. Never attempt to serve it's API over the public internet or for any commercial use case. Never use this application with malicious intent or to spam public services.