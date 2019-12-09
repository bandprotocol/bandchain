var http = require('http');
var handleRequest = require('./src/handleRequest');

var port = process.env.PORT || 8082;
var ip = '0.0.0.0';

var server = http.createServer(handleRequest);

console.log("Listening on http://" + ip + ":" + port);
server.listen(port, ip);
