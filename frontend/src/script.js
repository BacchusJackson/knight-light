let API_GET_STATUS = COUNTER_API_ENDPOINT ? `${COUNTER_API_ENDPOINT}/status` : undefined;
let API_POST_ADD = COUNTER_API_ENDPOINT ? `${COUNTER_API_ENDPOINT}/add` : undefined;
let API_WS = COUNTER_WS_API_ENDPOINT ? `${COUNTER_WS_API_ENDPOINT}/ws` : undefined;
let socket

// Run Update every 1 s
setInterval(() => {
  update();
}, 5000);

// Increment on the server and update
const increment = async (i) => {
  console.log(`Sending Request with Count: ${i}`);
  const addResData = await requestAdd(i);
  console.log(addResData);

  console.log(`Sending WS Message with Count: ${i}`);
  socket.send(JSON.stringify({count: i}));

  update();
};


// GET the status on the server
const update = async () => {

  const data = await requestStatus();
  console.log(data);
  document.getElementById('count').innerHTML =`API Server Total: ${data.currentCount.toString() || "-"}`;
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

window.addEventListener("load", (event) => {
  console.log("Window Loaded")
  console.log(`Connecting to websocket ${API_WS}...`) 

  socket = new WebSocket(API_WS)


  socket.addEventListener("open", (event) => {
    socket.send("Knight Light Mono Connected")
  })

  socket.addEventListener("message", (event) => {
    console.log(`Message from Server: ${event.Data}`)
  })
  
})


