const fs = require('fs')
const cp = require('child_process')

exports.joinCmd = arr => arr.join(' ')

exports.exec = (cmd, opts) =>
  new Promise((resolve, reject) => {
    cp.exec(cmd, opts, (err, stdout, stderr) => {
      if (err) {
        reject(err)
      } else {
        resolve(stdout + stderr)
      }
    })
  })

exports.execFile = (cmd, args, opts) =>
  new Promise((resolve, reject) => {
    cp.execFile(cmd, args, opts, (err, stdout, stderr) => {
      if (err) {
        reject(err)
      } else {
        resolve(stdout + stderr)
      }
    })
  })

exports.exists = path =>
  new Promise(function(resolve, reject) {
    fs.exists(path, err => {
      resolve(!err)
    })
  })

exports.writeFile = (path, contents) =>
  new Promise(function(resolve, reject) {
    fs.writeFile(path, contents, err => {
      if (err) {
        reject(err)
      } else {
        resolve(true)
      }
    })
  })

exports.readFile = path =>
  new Promise(function(resolve, reject) {
    fs.readFile(path, (err, data) => {
      console.log('path.length', data.length)
      if (err) {
        reject(err)
      } else {
        resolve(data)
      }
    })
  })

exports.mkdir = path => {
  return new Promise(function(resolve, reject) {
    fs.mkdir(path, err => {
      if (err) {
        reject(err)
      } else {
        resolve(true)
      }
    })
  })
}

exports.unlink = path =>
  new Promise(function(resolve, reject) {
    fs.unlink(path, err => {
      if (err) {
        reject(err)
      } else {
        resolve(true)
      }
    })
  })
