const http = require('http');

const server = http.createServer(function(req, res) {
  res.writeHead(200, {
    'Content-Type': 'text/html',
  });
  res.write('Hello World!');
  res.end();
});

server.listen(3000);
