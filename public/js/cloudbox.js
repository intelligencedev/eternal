
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
