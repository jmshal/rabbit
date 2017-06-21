const WebSocket = require('ws');

const server = new WebSocket.Server({
  port: 9000,
});

server.on('connection', function connection(ws) {
  ws.on('message', function incoming(message) {
    ws.send('ok, ' + message);
  });
  ws.send('welcome');
});
