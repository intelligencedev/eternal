# Roadmap Generate by Eternal for Eternal

Here’s a breakdown of functions that should be tested, along with strategies for testing the main package:

Functions to Test

loadConfig:

Why: This function is crucial for loading the application’s configuration. Errors in configuration loading can lead to major issues.
How:
Create test cases with valid and invalid YAML files.
Assert that the loaded configuration matches the expected values.
Verify error handling for invalid YAML or file not found scenarios.
createModelParam:

Why: This function converts a llm.Model into a ModelParams object, which is used for database interactions.
How:
Create test cases with different llm.Model objects.
Assert that the generated ModelParams object has the correct values.
EnsureDataPath:

Why: This function ensures the data directory exists.
How:
Use a temporary directory for testing.
Assert that the directory is created successfully.
Verify error handling if the directory cannot be created.
CopyFiles:

Why: This function copies files from an embedded filesystem.
How:
Create a temporary directory and copy files into it.
Assert that the files are copied correctly.
Verify error handling for file not found or permission issues.
IncrementTurn:

Why: This function increments a counter for chat elements.
How:
Call the function multiple times.
Assert that the returned value is incremented correctly.
NewSQLiteDB:

Why: This function initializes the SQLite database connection.
How:
Create a temporary database file.
Assert that the database connection is established successfully.
Verify error handling for invalid database paths.
AutoMigrate:

Why: This function migrates the database schema.
How:
Create a temporary database file.
Call the function with different models.
Verify that the database schema is updated correctly.
GetAppConfig:

Why: This function retrieves the application configuration from the database.
How:
Create a temporary database file with a known configuration.
Assert that the retrieved configuration matches the expected values.
GetProjects:

Why: This function retrieves all projects from the database.
How:
Create a temporary database file with known projects.
Assert that the retrieved projects match the expected values.
CreateProject:

Why: This function inserts a new project into the database.
How:
Create a temporary database file.
Call the function with a new project.
Verify that the project is inserted correctly.
DeleteProject:

Why: This function removes a project from the database.
How:
Create a temporary database file with a known project.
Call the function with the project name.
Verify that the project is deleted correctly.
ListProjects:

Why: This function retrieves all projects from the database.
How:
Create a temporary database file with known projects.
Assert that the retrieved projects match the expected values.
Create:

Why: This function inserts a new record into the database.
How:
Create a temporary database file.
Call the function with a new record.
Verify that the record is inserted correctly.
Find:

Why: This function retrieves records from the database.
How:
Create a temporary database file with known records.
Call the function with different criteria.
Assert that the retrieved records match the expected values.
First:

Why: This function retrieves the first record matching a specific criteria.
How:
Create a temporary database file with known records.
Call the function with different criteria.
Assert that the retrieved record matches the expected values.
FindByID:

Why: This function retrieves a record by its ID.
How:
Create a temporary database file with known records.
Call the function with different IDs.
Assert that the retrieved record matches the expected values.
UpdateByName:

Why: This function updates a record by its name.
How:
Create a temporary database file with a known record.
Call the function with the record name and updated values.
Verify that the record is updated correctly.
UpdateDownloadedByName:

Why: This function updates the downloaded state of a model by its name.
How:
Create a temporary database file with a known model.
Call the function with the model name and different downloaded states.
Verify that the downloaded state is updated correctly.
Delete:

Why: This function deletes a record from the database.
How:
Create a temporary database file with a known record.
Call the function with the record ID.
Verify that the record is deleted correctly.
LoadModelDataToDB:

Why: This function loads model data into the database.
How:
Create a temporary database file.
Call the function with a list of ModelParams.
Verify that the models are loaded correctly.
LoadImageModelDataToDB:

Why: This function loads image model data into the database.
How:
Create a temporary database file.
Call the function with a list of ImageModel.
Verify that the image models are loaded correctly.
AddSelectedModel:

Why: This function adds a selected model to the database.
How:
Create a temporary database file.
Call the function with a model name.
Verify that the model is added correctly.
RemoveSelectedModel:

Why: This function removes a selected model from the database.
How:
Create a temporary database file with a selected model.
Call the function with the model name.
Verify that the model is removed correctly.
GetSelectedModels:

Why: This function retrieves the list of selected models from the database.
How:
Create a temporary database file with selected models.
Assert that the retrieved models match the expected values.
CreateChat:

Why: This function inserts a new chat record into the database.
How:
Create a temporary database file.
Call the function with a new chat record.
Verify that the chat record is inserted correctly.
GetChats:

Why: This function retrieves all chat records from the database.
How:
Create a temporary database file with known chat records.
Assert that the retrieved chat records match the expected values.
GetChatByID:

Why: This function retrieves a chat record by its ID.
How:
Create a temporary database file with a known chat record.
Call the function with the chat record ID.
Assert that the retrieved chat record matches the expected values.
UpdateChat:

Why: This function updates an existing chat record in the database.
How:
Create a temporary database file with a known chat record.
Call the function with the chat record ID and updated values.
Verify that the chat record is updated correctly.
DeleteChat:

Why: This function removes a chat record from the database.
How:
Create a temporary database file with a known chat record.
Call the function with the chat record ID.
Verify that the chat record is deleted correctly.
CreateURLTracking:

Why: This function inserts a new URL into the URLTracking table.
How:
Create a temporary database file.
Call the function with a new URL.
Verify that the URL is inserted correctly.
ListURLTrackings:

Why: This function retrieves all URLs from the URLTracking table.
How:
Create a temporary database file with known URLs.
Assert that the retrieved URLs match the expected values.
DeleteURLTracking:

Why: This function removes a URL from the URLTracking table.
How:
Create a temporary database file with a known URL.
Call the function with the URL.
Verify that the URL is deleted correctly.
CreateDevTeam:

Why: This function creates a new development team in the database.
How:
Create a temporary database file.
Call the function.
Verify that the team is created correctly.
logFatalError:

Why: This function logs an error and exits the program.
How:
Call the function with different error messages.
Verify that the error message is logged correctly and the program exits.
deleteFile:

Why: This function deletes a file or directory.
How:
Create a temporary file or directory.
Call the function with the file or directory path.
Verify that the file or directory is deleted correctly.
fileExists:

Why: This function checks if a file or directory exists.
How:
Create a temporary file or directory.
Call the function with the file or directory path.
Assert that the function returns the correct result.
displayBanner:

Why: This function displays a banner.
How:
Call the function.
Verify that the banner is displayed correctly.
displayProjects:

Why: This function displays a table of projects.
How:
Call the function with a list of projects.
Verify that the table is displayed correctly.
GetHostInfo:

Why: This function retrieves information about the host system.
How:
Call the function.
Assert that the returned host information is valid.
getMacOSGPUInfo:

Why: This function retrieves GPU information for macOS.
How:
Call the function on a macOS system.
Assert that the returned GPU information is valid.
parseGPUInfo:

Why: This function parses GPU information from system_profiler output.
How:
Call the function with different system_profiler outputs.
Assert that the parsed GPU information is valid.
CreateProjectFolder:

Why: This function creates a project folder and data in the database.
How:
Create a temporary directory.
Call the function with a new project.
Verify that the project folder is created and the data is inserted into the database.
GetProjectByName:

Why: This function retrieves a project by its name.
How:
Create a temporary database file with a known project.
Call the function with the project name.
Assert that the retrieved project matches the expected values.
performToolWorkflow:

Why: This function performs the tool workflow on a chat message.
How:
Create test cases with different chat messages and tool configurations.
Assert that the workflow is executed correctly and the chat message is modified as expected.
handleToolToggle:

Why: This function toggles the state of various tools.
How:
Create test cases with different tool names and enabled states.
Assert that the tool state is updated correctly.
handleToolList:

Why: This function retrieves and returns a list of tools.
How:
Call the function.
Assert that the returned list of tools matches the expected values.
queuePrompt:

Why: This function queues a prompt for processing by the ComfyUI service.
How:
Create test cases with different prompts and client IDs.
Assert that the prompt is queued correctly and the response is as expected.
getHistory:

Why: This function retrieves the history of prompts and outputs from the ComfyUI service.
How:
Create test cases with different history entries.
Assert that the retrieved history matches the expected values.
getImage:

Why: This function retrieves an image from the ComfyUI service.
How:
Create test cases with different image filenames and subfolders.
Assert that the retrieved image data is valid.
getImages:

Why: This function retrieves multiple images from the ComfyUI service.
How:
Create test cases with different prompts.
Assert that the retrieved image data is valid.
handleRoleSelection:

Why: This function handles the selection of assistant roles.
How:
Create test cases with different role names.
Assert that the current role is updated correctly.
fetchURL:

Why: This function fetches content from a URL.
How:
Create test cases with different URLs.
Assert that the retrieved content is valid.
pollURL:

Why: This function polls a URL until a successful response is received.
How:
Create test cases with different URLs and timeouts.
Assert that the retrieved content is valid and the function returns within the specified timeout.
handleImgSetWorkflow:

Why: This function sets the workflow for the image generation tool.
How:
Create test cases with different workflow names.
Assert that the workflow is set correctly.
handleImgSetResolution:

Why: This function sets the resolution for the image generation tool.
How:
Create test cases with different width and height values.
Assert that the resolution is set correctly.
performImageGen:

Why: This function performs image generation using the ComfyUI service.
How:
Create test cases with different prompts and configurations.
Assert that the image is generated correctly and the response is as expected.
loadPromptText:

Why: This function loads prompt text from a JSON file.
How:
Create test cases with different JSON files.
Assert that the loaded prompt text is valid.
SaveBytesAsImage:

Why: This function saves image data as a PNG file.
How:
Create test cases with different image data.
Assert that the image is saved correctly.
handleGetTools:

Why: This function retrieves and renders tools.
How:
Call the function.
Assert that the rendered HTML is valid.
handleGetChats:

Why: This function retrieves and returns all chat records.
How:
Create a temporary database file with known chat records.
Assert that the retrieved chat records match the expected values.
handleGetChatByID:

Why: This function retrieves and returns a chat record by its ID.
How:
Create a temporary database file with a known chat record.
Call the function with the chat record ID.
Assert that the retrieved chat record matches the expected values.
handleUpdateChat:

Why: This function updates a chat record by its ID.
How:
Create a temporary database file with a known chat record.
Call the function with the chat record ID and updated values.
Verify that the chat record is updated correctly.
handleDeleteChat:

Why: This function handles the deletion of a chat by its ID.
How:
Create a temporary database file with a known chat record.
Call the function with the chat record ID.
Verify that the chat record is deleted correctly.
handleListProjects:

Why: This function retrieves and returns a list of projects from the database.
How:
Create a temporary database file with known projects.
Assert that the retrieved projects match the expected values.
handleUpload:

Why: This function handles file uploads and saves them to the specified directory.
How:
Create a temporary directory.
Call the function with a file.
Verify that the file is uploaded correctly.
handleDPSearch:

Why: This function handles search requests using DuckDuckGo.
How:
Create test cases with different search queries.
Assert that the returned search results are valid.
handleSSEUpdates:

Why: This function handles Server-Sent Events (SSE) for updates.
How:
Create a test client that connects to the SSE endpoint.
Verify that the client receives updates as expected.
handleWebSocket:

Why: This function handles WebSocket connections for general use.
How:
Create a test client that connects to the WebSocket endpoint.
Send messages to the server and verify that the responses are as expected.
handleWebSocketConnection:

Why: This function handles WebSocket connections and processes messages.
How:
Create a test client that connects to the WebSocket endpoint.
Send messages to the server and verify that the responses are as expected.
readAndUnmarshalMessage:

Why: This function reads and unmarshals a WebSocket message.
How:
Create test cases with different WebSocket messages.
Assert that the unmarshaled message is valid.
handleError:

Why: This function handles errors that occur during message processing.
How:
Create test cases with different error messages.
Verify that the error is handled correctly.
handleChatTurnFinished:

Why: This function handles the completion of a chat turn.
How:
Create test cases with different chat messages and errors.
Verify that the chat turn is stored correctly in the database.
handleChatMemory:

Why: This function retrieves and returns chat memory.
How:
Create test cases with different chat messages and tool configurations.
Assert that the retrieved chat memory is valid.
storePromptInBleve:

Why: This function stores the user’s prompt in the Bleve index.
How:
Create test cases with different prompts.
Verify that the prompt is stored correctly in the index.
handleTextSplitAndIndex:

Why: This function handles the splitting and indexing of text.
How:
Create test cases with different text and chunk sizes.
Verify that the text is split and indexed correctly.
searchSimilarEmbeddings:

Why: This function searches for similar embeddings in the database.
How:
Create test cases with different prompts and embedding configurations.
Assert that the returned similar embeddings are valid.
handleModelData:

Why: This function retrieves and returns data for a specific model.
How:
Create a temporary database file with a known model.
Call the function with the model name.
Assert that the retrieved model data matches the expected values.
handleModelDownloadUpdate:

Why: This function updates the download status of a model.
How:
Create a temporary database file with a known model.
Call the function with the model name and different downloaded states.
Verify that the downloaded state is updated correctly.
handleModelUpdate:

Why: This function updates the model data in the database.
How:
Create a temporary database file with a known model.
Call the function with the model name and updated values.
Verify that the model data is updated correctly.
handleModelCards:

Why: This function retrieves and renders model cards.
How:
Call the function.
Assert that the rendered HTML is valid.
handleModelSelect:

Why: This function handles the selection of models for use.
How:
Create test cases with different model names and actions.
Verify that the model selection is updated correctly.
handleSelectedModels:

Why: This function retrieves and returns the list of selected models.
How:
Create a temporary database file with selected models.
Assert that the retrieved models match the expected values.
handleModelDownload:

Why: This function handles the download of a specified model.
How:
Create test cases with different model names.
Verify that the model is downloaded correctly.
handleImgModelDownload:

Why: This function handles the download of image generation models.
How:
Create test cases with different model names.
Verify that the image generation model is downloaded correctly.
handleOpenAIModels:

Why: This function retrieves and returns a list of OpenAI models.
How:
Call the function.
Assert that the returned list of OpenAI models is valid.
handleGetRoles:

Why: This function retrieves and returns a list of roles.
How:
Call the function.
Assert that the returned list of roles is valid.
DownloadDefaultImageModel:

Why: This function downloads the default image model.
How:
Create test cases with different model names.
Verify that the image generation model is downloaded correctly.
initServer:

Why: This function initializes the server by setting up necessary directories and files.
How:
Create a temporary directory.
Call the function with the directory path.
Verify that the directory is created and the files are copied correctly.
setupDirectory:

Why: This function creates a directory and copies files into it.
How:
Create a temporary directory.
Call the function with the directory path and source directory.
Verify that the directory is created and the files are copied correctly.
installPythonRequirements:

Why: This function installs Python requirements from a requirements.txt file.
How:
Create a temporary requirements.txt file.
Call the function with the file path.
Verify that the requirements are installed correctly.
setExecutablePermissions:

Why: This function sets executable permissions on all files in a directory.
How:
Create a temporary directory with files.
Call the function with the directory path.
Verify that the files have executable permissions.
initializeApplication:

Why: This function initializes the application with the given configuration.
How:
Create a temporary directory.
Call the function with a configuration object.
Verify that the application is initialized correctly.
createDataDirectory:

Why: This function creates the data directory.
How:
Create a temporary directory.
Call the function with the directory path.
Verify that the directory is created correctly.
initializeServer:

Why: This function initializes the server by setting up necessary directories and files.
How:
Create a temporary directory.
Call the function with the directory path.
Verify that the server is initialized correctly.
initializeDatabase:

Why: This function initializes the SQLite database.
How:
Create a temporary directory.
Call the function with the directory path.
Verify that the database is initialized correctly.
initializeDefaultProject:

Why: This function initializes the default project.
How:
Create a temporary directory.
Call the function with a configuration object.
Verify that the default project is initialized correctly.
initializeSearchIndex:

Why: This function initializes the search index.
How:
Create a temporary directory.
Call the function with the directory path.
Verify that the search index is initialized correctly.
loadModelParams:

Why: This function loads model parameters from the configuration.
How:
Create a configuration object with known model parameters.
Call the function with the configuration object.
Verify that the model parameters are loaded correctly.
displayModelParams:

Why: This function displays model parameters in a table format.
How:
Call the function with a list of model parameters.
Verify that the table is displayed correctly.
startComfyUI:

Why: This function starts the ComfyUI service.
How:
Create a temporary directory.
Call the function with a configuration object.
Verify that the ComfyUI service is started correctly.
installComfyUIRequirements:

Why: This function installs ComfyUI requirements.
How:
Create a temporary directory with a requirements.txt file.
Call the function with the directory path.
Verify that the requirements are installed correctly.
runComfyUI:

Why: This function runs the ComfyUI service.
How:
Create a temporary directory.
Call the function with the directory path and port.
Verify that the ComfyUI service is running correctly.
monitorComfyUIOutput:

Why: This function monitors and logs the output of the ComfyUI service.
How:
Create a temporary directory.
Call the function with the stdout and stderr pipes of the ComfyUI process.
Verify that the output is logged correctly.
scanAndLog:

Why: This function reads and logs output from a reader.
How:
Create a test reader with known output.
Call the function with the reader and a logging function.
Verify that the output is logged correctly.
waitForComfyUI:

Why: This function waits for the ComfyUI service to become available.
How:
Create a temporary directory.
Start the ComfyUI service in the directory.
Call the function with a configuration object.
Verify that the function returns when the ComfyUI service is available.
stopComfyUI:

Why: This function stops the ComfyUI service gracefully.
How:
Create a temporary directory.
Start the ComfyUI service in the directory.
Call the function with a context object.
Verify that the ComfyUI service is stopped correctly.
setupComfyUI:

Why: This function sets up the ComfyUI directory.
How:
Create a temporary directory.
Call the function with the directory path.
Verify that the ComfyUI directory is created correctly.
setupImpactPack:

Why: This function sets up the Impact Pack for ComfyUI.
How:
Create a temporary directory.
Call the function with the directory path.
Verify that the Impact Pack directory is created correctly.
runFrontendServer:

Why: This function starts the Fiber frontend server.
How:
Create a temporary directory.
Call the function with a configuration object and a list of model parameters.
Verify that the frontend server is started correctly.
createFiberApp:

Why: This function initializes and returns a new Fiber application.
How:
Call the function with a configuration object.
Verify that the Fiber application is initialized correctly.
handleGracefulShutdown:

Why: This function handles the graceful shutdown of the application.
How:
Create a temporary directory.
Start the frontend server in the directory.
Call the function with a context object, a Fiber application, a configuration object, and a list of model parameters.
Verify that the application is shut down gracefully.
cleanupDevMode:

Why: This function performs cleanup tasks specific to development mode.
How:
Create a temporary directory.
Call the function with a configuration object and a list of model parameters.
Verify that the cleanup tasks are performed correctly.
Testing the main Package

Integration Tests:

Why: Integration tests are essential for verifying that different parts of your application work together as expected.
How:
Use a test HTTP client to interact with your API endpoints.
Mock external dependencies like databases or external services.
Test the entire workflow of your application, from API requests to database interactions and response generation.
End-to-End Tests:

Why: End-to-end tests are the most comprehensive type of test, simulating real user interactions with your application.
How:
Use a tool like Selenium or Cypress to automate browser interactions.
Test the entire user journey, from navigating your application to submitting forms and viewing results.
Key Considerations

Mocking: Use mocking libraries to isolate your tests and avoid dependencies on external services or databases.
Temporary Resources: Create temporary files, directories, and databases for your tests to avoid polluting your development environment.
Test Data: Use realistic test data to ensure your tests are comprehensive and representative of real-world usage.
Test Coverage: Aim for high test coverage to ensure that all critical parts of your application are tested.