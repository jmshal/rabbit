const WebSocket = require('ws');

// TLS
process.env.NODE_TLS_REJECT_UNAUTHORIZED = '0';
const ws = new WebSocket('wss://:10443/ws');

// or HTTP
// const ws = new WebSocket('ws://:1080/ws');

ws.on('open', function open() {
  ws.send('hello!');
});

ws.on('message', function incoming(message) {
  console.log('message from server: ' + message);
});
