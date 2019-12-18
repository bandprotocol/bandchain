const {
  cargoCmd,
  rustcCmd,
  wasmPackCmd,
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
  mkdir,
  unlink
} = require("./common.js");

function checkBuildPlan(plan) {
  let success = true;
  let invocations = plan["invocations"];

  var custom_build = invocations.find(function(element) {
    return element["target_kind"].includes("custom-build");
  });

  if (custom_build) {
    success = false;
    return { success, output: "", message: "the build includes custom builds" };
  }

  if (invocations.length > 1) {
    success = false;
    return {
      success,
      output: "",
      message: "dependencies are currently deactivated"
    };
  }

  return { success: true };
}

async function wasmGC(wasmFile, callback) {
  if (!(await exists(wasmFile))) {
    throw new Error("wasm is not found");
  }
  await exec(joinCmd([wasmGCCmd, wasmFile]));
}

async function cargo(tar, options = {}) {
  let crateName =
    "rustc_h_" +
    (options.session.replace(/[\.\/]/g, "") ||
      Math.random()
        .toString(36)
        .slice(2));
  let crateDir = tempDir + "/" + crateName;

  try {
    await mkdir(crateDir);
  } catch (err) {}

  let rustTar = crateDir + "/" + "lib.tar";
  let wasmFile = crateDir + "/" + "pkg/main_bg.wasm";
  await writeFile(rustTar, new Buffer(tar, "base64").toString("ascii"));

  let args = ["tar", "xvf", rustTar, "-C", crateDir];
  await exec(joinCmd(args));

  try {
    let args = ["build", "--out-name=main", crateDir];
    let output;
    let success = false;
    let opts = {};
    try {
      output = await execFile(wasmPackCmd, args, opts);
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
    //await unlink(crateDir);
  }
}

module.exports = function(source, options, callback) {
  cargo(source, options)
    .then(result => callback(null, result))
    .catch(err => callback(err, null));
};
