const http = require('http');

const hostname = '127.0.0.1';
const port = 3000;

const server = http.createServer((req, res) => {
  let responseBody;
  let contentType;

  if (req.url === '/json') {
    responseBody = JSON.stringify('Hello from Node JSON endpoint');
    contentType = 'application/json';
  } else {
    responseBody = 'Hello from node ';
    contentType = 'text/plain';
  }

  res.statusCode = 200;
  res.setHeader('Content-Type', contentType);
  res.end(responseBody);
});

server.listen(port, hostname, () => {
  console.log(`Server running at http://${hostname}:${port}/`);
});
