const {
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
  unlink
} = require("./common.js");

async function wasmGC(wasmFile, callback) {
  if (!(await exists(wasmFile))) {
    throw new Error("wasm is not found");
  }
  await exec(joinCmd([wasmGCCmd, wasmFile]));
}

async function rustc(source, options = {}) {
  let baseName = __dirname + "/../../owasm-boilerplate/";
  let libFileName = "lib";
  let rustFile = baseName + "src/" + libFileName + ".rs";
  let wasmFile = baseName + "pkg/awesome_oracle_bg.wasm";

  await writeFile(rustFile, source);

  try {
    let args = ["build", baseName];
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
    // await unlink(rustFile);
  }
}

async function rustc_original(source, options = {}) {
  let crateName =
    "rustc_h_" +
    Math.random()
      .toString(36)
      .slice(2);
  let baseName = tempDir + "/" + crateName;
  let rustFile = baseName + ".rs";
  let wasmFile = baseName + ".wasm";
  await writeFile(rustFile, source);

  try {
    let args = [rustcCmd, rustFile];
    args.push("--target=wasm32-unknown-unknown");
    args.push("--crate-type=cdylib");
    if (options.lto) args.push("-Clto");
    if (options.debug) args.push("-g");
    switch (options.opt_level) {
      case "s":
      case "z":
      case "0":
      case "1":
      case "2":
      case "2":
        args.push("-Copt-level=" + options.opt_level);
        break;
    }
    args.push("-o");
    args.push(wasmFile);
    for (let i = 0; i < wasmBindgenDeps.length; i++) {
      args.push("-L");
      args.push(wasmBindgenDeps[i]);
    }
    let output;
    let success = false;
    let opts = {
      // env vars needed for #[wasm_bindgen]
      env: {
        CARGO_PKG_NAME: "main",
        CARGO_PKG_VERSION: "1.0.0"
      }
    };
    try {
      output = await exec(joinCmd(args), opts);
      success = true;
    } catch (e) {
      output = "error: " + e;
    }
    try {
      if (!success) return { success, output: "", message: output };
      let wasmBindgenJs = "";
      let wasm = await readFile(wasmFile);
      let m = await WebAssembly.compile(wasm);
      let ret = { success, message: output };
      if (
        WebAssembly.Module.customSections(m, "__wasm_bindgen_unstable")
          .length !== 0
      ) {
        await exec(
          joinCmd([
            wasmBindgenCmd,
            wasmFile,
            "--no-modules",
            "--out-dir",
            tempDir
          ])
        );
        wasm = await readFile(baseName + "_bg.wasm");
        ret.wasmBindgenJs = (await readFile(baseName + ".js")).toString();
      } else {
        await exec(joinCmd([wasmGCCmd, wasmFile]));
        wasm = await readFile(wasmFile);
      }
      ret.output = wasm.toString("base64");
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
  rustc(source, options)
    .then(result => callback(null, result))
    .catch(err => callback(err, null));
};
