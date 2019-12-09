const {
  rustcCmd,
  wasmRunCmd,
  wasmGCCmd,
  tempDir,
  wasmBindgenCmd,
  wasmBindgenDeps
} = require("../config.js");
const {
  exec,
  execFile,
  joinCmd,
  exists,
  writeFile,
  readFile,
  unlink
} = require("./common.js");

async function wasmGC(wasmFile, callback) {
  if (!(await exists(wasmFile))) {
    throw new Error("wasm is not found");
  }
  await exec(joinCmd([wasmGCCmd, wasmFile]));
}

async function wasmRun(source, options = {}) {
  let wasmFile =
    tempDir +
    "/" +
    "owasm_" +
    Math.random()
      .toString(36)
      .slice(2) +
    ".wasm";

  const sourceBytes = Buffer.from(source, "base64");

  await writeFile(wasmFile, sourceBytes);

  try {
    let args = [wasmFile];
    let output;
    let success = false;
    let opts = {};
    try {
      output = await execFile(wasmRunCmd, args, opts);
      success = true;
    } catch (e) {
      output = "error: " + e;
    }
    try {
      if (!success) return { success, output: "", message: output };
      let ret = { success, message: output };
      ret.output = (await readFile(wasmFile)).toString("base64");
      return ret;
    } finally {
      // if (success)
      // await unlink(wasmFile);
    }
  } finally {
    // await unlink(rustFile);
  }
}

module.exports = function(source, options, callback) {
  wasmRun(source, options)
    .then(result => callback(null, result))
    .catch(err => callback(err, null));
};
