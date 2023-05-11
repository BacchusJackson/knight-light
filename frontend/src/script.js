let API_GET_STATUS = COUNTER_API_ENDPOINT ? `${COUNTER_API_ENDPOINT}/status` : undefined;
let API_POST_ADD = COUNTER_API_ENDPOINT ? `${COUNTER_API_ENDPOINT}/add` : undefined;
let CHAT_API = CHAT_API_ENDPOINT ? `${CHAT_API_ENDPOINT}/ws` : undefined;

let socket

// Run Update every 1 s
setInterval(() => {
  updateAll();
}, 5000);

// Increment on the server and update
const increment = async (i, element) => {
  console.log(`Sending Request with Count: ${i}`);
  const addResData = await requestAdd(i);
  console.log(addResData);

  // console.log(`Sending WS Message with Count: ${i}`);
  // socket.send(JSON.stringify({count: i}));

  update(element);
};


// GET the status on the server
const update = async (element) => {

  const data = await requestStatus();
  console.log(data);
  document.getElementById(element).innerHTML = counterString(data);
}

const updateAll = async() => {
  const data = await requestStatus();
  console.log(data);
  document.getElementById("count-1").innerHTML =counterString(data);
  document.getElementById("count-2").innerHTML =counterString(data);
}

const counterString = (data) => {
  console.log(data.lastUpdated)
  const date = new Date(data.lastUpdated);
  const month = date.toLocaleString('en-US', { month: 'short', timeZone: "America/New_York" });
  const dateStr = `${month} ${date.getDate()}, ${date.getFullYear()} ${date.getHours()}:${date.getMinutes()}`
  return `Total: ${data.currentCount.toString() || '-'} @ ${dateStr}`
}


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
    const response = await fetch(API_POST_ADD, {method: "POST", body: JSON.stringify({count: i})})
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

document.getElementById("chat-input").addEventListener("keyup", event => {
  if(event.key !== "Enter") return; // Use `.key` instead.
  document.getElementById("chat-send-button").click();
  event.preventDefault(); // No need to `return false;`.
});

window.addEventListener("load", (event) => {
  document.getElementById("chat-box").innerHTML = `[client] Connecting to websocket ${CHAT_API}...<br>`

  socket = new WebSocket(CHAT_API)

  socket.addEventListener("open", (event) => {
    document.getElementById("chat-box").innerHTML += "[client] Chat API Websocket connected!<br>"
  })

  socket.addEventListener("message", (event) => {
    console.log(`Message from Server: ${event.data}`)
    const chatBox = document.getElementById("chat-box")
    chatBox.innerHTML = `${chatBox.innerHTML}<br> ${event.data}`
  })

})


