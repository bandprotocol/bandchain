import * as gulp from "gulp";
import { Service, project } from "@wasm/studio-utils";

function getAllFiles() {
  return [
    project.getFile("src/lib.rs"),
    project.getFile("src/logic.rs"),
    project.getFile("Cargo.toml")
  ];
}

// Random session id
if (!window.localStorage.getItem("BUILD_SESSION"))
  window.localStorage.setItem(
    "BUILD_SESSION",
    Array(40)
      .fill(0)
      .map(() => Math.floor(Math.random() * 36).toString(36))
      .join("")
  );

gulp.task("build", async () => {
  console.log("Build session", window.localStorage.getItem("BUILD_SESSION"));

  const options = {
    debug: true,
    cargo: true,
    session: window.localStorage.getItem("BUILD_SESSION")
  };
  const data = await Service.compileFiles(
    getAllFiles(),
    "rust",
    "wasm",
    options
  );
  const outWasm = project.newFile("out.wasm", "wasm", true);
  outWasm.setData(data["a.wasm"]);
});

gulp.task("test", async () => {
  console.log("Test session", window.localStorage.getItem("BUILD_SESSION"));

  const options = {
    debug: true,
    cargo: true,
    session: window.localStorage.getItem("BUILD_SESSION")
  };
  const data = await Service.testCargo(getAllFiles(), options);
  return data;
});

gulp.task("run", async () => {
  console.log("Run out.wasm", project.getFile("out.wasm"));

  const options = {};
  const data = await Service.wasmRun(project.getFile("out.wasm"), options);
  return data;
});

gulp.task("default", ["build", "test", "run"], async () => {});
