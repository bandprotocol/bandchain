var qs = require("querystring");
var cargo = require("./cargo");
var cargoTest = require("./cargoTest");
var wasmRun = require("./wasmRun");
var rustc = require("./rustc");
var rustfmt = require("./rustfmt");

function notAllowed(res) {
  res.writeHead(503);
  res.end();
}

function showError(res, err) {
  res.writeHead(501);
  res.end(`<pre>${err.toString()}</pre>`);
}

function readFormData(request, callback) {
  var body = "";

  request.on("data", function(data) {
    body += data;

    // Too much POST data, kill the connection!
    // 1e6 === 1 * Math.pow(10, 6) === 1 * 1000000 ~~~ 1MB
    if (body.length > 1e6) request.connection.destroy();
  });

  request.on("end", function() {
    try {
      callback(null, JSON.parse(body));
    } catch (ex) {
      callback(ex);
    }
  });
}

module.exports = function handleRequest(req, res) {
  res.setHeader("Access-Control-Allow-Origin", "*");
  res.setHeader(
    "Access-Control-Allow-Headers",
    "Origin, X-Requested-With, Content-Type, Accept"
  );
  res.setHeader("Access-Control-Allow-Methods", "OPTIONS, GET, POST");

  if (req.method == "OPTIONS") {
    res.writeHead(200);
    res.end();
    return;
  }

  if (req.url == "/cargo-test") {
    if (req.method != "POST") return notAllowed(res);
    readFormData(req, (err, post) => {
      if (err) return showError(res, err);
      cargoTest(post.tar, post.options, (err, result) => {
        if (err) return showError(res, err);
        res.setHeader("Content-type", "application/json");
        res.writeHead(200);
        res.end(JSON.stringify(result));
      });
    });
    return;
  }

  if (req.url == "/cargo") {
    if (req.method != "POST") return notAllowed(res);
    readFormData(req, (err, post) => {
      if (err) return showError(res, err);
      cargo(post.tar, post.options, (err, result) => {
        if (err) return showError(res, err);
        res.setHeader("Content-type", "application/json");
        res.writeHead(200);
        res.end(JSON.stringify(result));
      });
    });
    return;
  }

  if (req.url == "/rustc") {
    if (req.method != "POST") return notAllowed(res);
    readFormData(req, (err, post) => {
      if (err) return showError(res, err);
      rustc(post.code, post.options, (err, result) => {
        if (err) return showError(res, err);
        res.setHeader("Content-type", "application/json");
        res.writeHead(200);
        res.end(JSON.stringify(result));
      });
    });
    return;
  }

  if (req.url == "/wasm-run") {
    if (req.method != "POST") return notAllowed(res);
    readFormData(req, (err, post) => {
      if (err) return showError(res, err);
      wasmRun(post.code, post.options, (err, result) => {
        if (err) return showError(res, err);
        res.setHeader("Content-type", "application/json");
        res.writeHead(200);
        res.end(JSON.stringify(result));
      });
    });
    return;
  }

  if (req.url == "/rustfmt") {
    if (req.method != "POST") return notAllowed(res);
    readFormData(req, (err, post) => {
      if (err) return showError(res, err);
      rustfmt(post.code, post.options, (err, result) => {
        if (err) return showError(res, err);
        res.setHeader("Content-type", "application/json");
        res.writeHead(200);
        res.end(JSON.stringify(result));
      });
    });
    return;
  }

  res.writeHead(404);
  res.end();
};
