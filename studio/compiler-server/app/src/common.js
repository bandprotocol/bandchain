const fs = require("fs");
const cp = require("child_process");

function joinCmd(arr) {
  return arr.join(" ");
}

function exec(cmd, opts) {
  return new Promise((resolve, reject) => {
    cp.exec(cmd, opts, (err, stdout, stderr) => {
      if (err) {
        reject(err);
      } else {
        resolve(stdout + stderr);
      }
    });
  });
}
function execFile(cmd, args, opts) {
  return new Promise((resolve, reject) => {
    cp.execFile(cmd, args, opts, (err, stdout, stderr) => {
      if (err) {
        reject(err);
      } else {
        resolve(stdout + stderr);
      }
    });
  });
}

function exists(path) {
  return new Promise(function(resolve, reject) {
    fs.exists(path, (err, stats) => {
      resolve(!err);
    });
  });
}

function writeFile(path, contents) {
  return new Promise(function(resolve, reject) {
    fs.writeFile(path, contents, err => {
      if (err) {
        reject(err);
      } else {
        resolve(true);
      }
    });
  });
}

function readFile(path) {
  return new Promise(function(resolve, reject) {
    fs.readFile(path, (err, data) => {
      if (err) {
        reject(err);
      } else {
        resolve(data);
      }
    });
  });
}

function mkdir(path) {
  return new Promise(function(resolve, reject) {
    fs.mkdir(path, err => {
      if (err) {
        reject(err);
      } else {
        resolve(true);
      }
    });
  });
}

function unlink(path) {
  return new Promise(function(resolve, reject) {
    fs.unlink(path, err => {
      if (err) {
        reject(err);
      } else {
        resolve(true);
      }
    });
  });
}

module.exports = {
  joinCmd,
  exec,
  execFile,
  exists,
  writeFile,
  readFile,
  mkdir,
  unlink
};
