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