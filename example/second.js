const http = require('http');

const server = http.createServer(function(req, res) {
  res.writeHead(200, {
    'Content-Type': 'application/json',
  });
  res.write(JSON.stringify({
    url: req.url,
    headers: req.headers,
  }, null, '  '));
  res.end();
});

server.listen(5000);
