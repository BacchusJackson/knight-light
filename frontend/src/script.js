let API_GET_STATUS = COUNTER_API_ENDPOINT ? `${COUNTER_API_ENDPOINT}/status` : undefined;
let API_POST_ADD = COUNTER_API_ENDPOINT ? `${COUNTER_API_ENDPOINT}/add` : undefined;
let API_COUNTER_SSE_STATUS = COUNTER_API_ENDPOINT ? `${COUNTER_API_ENDPOINT}/sse` : undefined;
let CHAT_API = CHAT_API_ENDPOINT ? `${CHAT_API_ENDPOINT}/ws` : undefined;

let socket


// Increment on the server and update
const increment = async (i, element) => {
  console.log(`Sending Request with Count: ${i}`);
  const addResData = await requestAdd(i);

};


// GET the status on the server

const requestStatus = async () => {
  if (!API_GET_STATUS) {
    console.error("API_GET_STATUS is not defined")
    return {};
  }

  try {
    console.log(`Fetching from ${API_GET_STATUS}`);
    const response = await fetch(API_GET_STATUS,{method: "GET"});
    return response.json();

  } catch (error) {
    console.error(error);
    return {};
  }

}

// POST the count to the server
const requestAdd = async (i) => {
  if (!API_POST_ADD) {
    console.error("API_POST_ADD is not defined")
    return {};
  }

  console.log(`POST to ${API_POST_ADD}`)
  try {
    const requestBody = JSON.stringify({count: i});

    const response = await fetch(API_POST_ADD, {method: "POST", body: requestBody});

    return response.json();

  } catch (error) {
    console.error(error);
    return {};
  }
}

const sendMessage = async () => {
  const inputBox = document.getElementById("chat-input")
  const msg = inputBox.value

  inputBox.value = ""

  socket.send(msg);

}

const setupCounterSSE = () => {

  const source = new EventSource(API_COUNTER_SSE_STATUS);

  source.onopen = (event) => {
    console.log("SSE Connected!");
  }

  source.onmessage = (event) => {
    const countStatus = JSON.parse(event.data)
    const countStr = `<b>Total:</b> ${countStatus.currentCount}<br><b>Last Updated:</b> ${countStatus.lastUpdated.substring(0,19)}`
    document.getElementById("count-1").innerHTML = countStr;
  }

  source.onerror = (event) => {

    console.log("SSE error received");
    source.close();
  }

}

const setupChatBox = () => {
  const chatBox = document.getElementById("chat-box");
  chatBox.innerHTML += `[client] Connecting to websocket ${CHAT_API}...` + "<br>";

  socket = new WebSocket(CHAT_API);

  socket.onerror = (event) => {
    chatBox.innerHTML += `[client] Connection ${CHAT_API} failed. See log for details.` + "<br>";
    console.log(event);
  }

  socket.addEventListener("open", (event) => {
    chatBox.innerHTML += "[client] Chat API Websocket connected!<br>";
  })

  socket.addEventListener("message", (event) => {
    chatBox.innerHTML = `${chatBox.innerHTML}<br> ${event.data}`;
  })

  document.getElementById("chat-input").addEventListener("keyup", event => {
    if(event.key !== "Enter") return; // Use `.key` instead.
    document.getElementById("chat-send-button").click();
    event.preventDefault(); // No need to `return false;`.
  });
}

const setupButtonClickAnimation = () => {
  const incrementButtons = document.querySelectorAll('.increment');

  incrementButtons.forEach((button) => {
    button.addEventListener('click', () => {
      button.style.animation = "none";
      button.style.animation = "clickAnimation 300ms";

      setTimeout(() => {button.style.animation = "none"}, 300);
    });
  });

}

window.addEventListener("load", (event) => {
  setupChatBox();
  setupCounterSSE();
  setupButtonClickAnimation();
})

// Add button click animation


  // incrementButtons.forEach((button) => {
  //   button.addEventListener('click', () => {
  //   });
  // })
  //
  //
