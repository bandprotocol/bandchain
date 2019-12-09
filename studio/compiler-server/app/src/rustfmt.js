const fs = require('fs');
const { exec } = require('child_process');

const { rustfmtCmd, tempDir } = require("../config.js");
const { joinCmd } = require("./common.js");

module.exports = function rustfmt(source, options, callback) {
  var baseName = tempDir + '/rustfmt_h_' + Math.random().toString(36).slice(2);
  var rustFile = baseName + '.rs';
  fs.writeFile(rustFile, source, (err) => {
    if (err) return callback(err);

    exec(
      joinCmd([rustfmtCmd, rustFile]),
      (err, stdout, stderr) => {
        var console = stdout.toString() + stderr.toString();
        fs.readFile(rustFile, (err, content) => {
          var success = !err;
          var formattedContent = err ? undefined : content.toString('base64');
          if (!err) fs.unlink(rustFile, () => {});
          fs.unlink(rustFile, () => {});
          callback(null, { success, output: formattedContent, console, });
        });
      }
    );
  });
};
