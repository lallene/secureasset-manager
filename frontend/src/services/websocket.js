let socket = null;

export const connectWebSocket = (callback) => {
socket = new WebSocket("ws://localhost:8080/ws");

socket.onmessage = (event) => {
const data = JSON.parse(event.data);

callback(data);
};

socket.onclose = () => {
setTimeout(() => {
connectWebSocket(callback);
}, 3000);
};
};