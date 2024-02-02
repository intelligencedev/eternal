const uploadButton = document.getElementById('upload');
const fileInput = document.getElementById('file-input');

uploadButton.addEventListener('click', function () {
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
  console.log("Appending image to chat...")
  console.log(file.name)
  await insertImageIntoChat("./public/uploads/" + file.name);

  // Scroll to the bottom of the chat
  resetScroll(document.getElementById('chat-view'));

  // Run inference over image
  await fetchImageProcessingResult("./public/uploads/" + file.name);
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

