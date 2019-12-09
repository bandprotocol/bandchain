D3N OWASM Studio
====
[![Build Status](https://travis-ci.org/wasdk/WebAssemblyStudio.svg?branch=master)](https://travis-ci.org/wasdk/WebAssemblyStudio) [![Coverage Status](https://coveralls.io/repos/github/wasdk/WebAssemblyStudio/badge.svg)](https://coveralls.io/github/wasdk/WebAssemblyStudio) [![Maintainance Status](https://img.shields.io/badge/maintained-seldom-yellowgreen.svg)](https://github.com/wasdk/WebAssemblyStudio/issues/381)

This repository contains the [D3N OWASM Studio](https://owasm.bandprotocol.com) website source code, which is a modified version of [WebAssembly Studio](https://webassembly.studio).

Running your own local copy of the website
===

To run a local copy, you will need to install node.js and webpack on your computer, then run the following commands:

```
yarn install
```

Easiest way to start dev server on [https://localhost:28443/](https://localhost:28443/) is to run

```
yarn start
```

To build WebAssembly Studio whenever a file changes run:

```
yarn build-watch
```

To start a dev web serve only, run:

```
yarn dev-server
```

Before submitting a pull request run:

```
npm test
```

### Contributing

Please get familiar with the [contributing guide](https://github.com/wasdk/WebAssemblyStudio/wiki/Contributing).

Any doubts or questions? You can always find us on slack at http://wasm-studio.slack.com

Need a slack invite? https://wasm-studio-invite.herokuapp.com/

### Credits

This project depends on several excellent libraries and tools:

* [Monaco Editor](https://github.com/Microsoft/monaco-editor) is used for rich text editing, tree views and context menus.

* [WebAssembly Binary Toolkit](https://github.com/WebAssembly/wabt) is used to assemble and disassemble `.wasm` files.

* [Binaryen](https://github.com/WebAssembly/binaryen/) is used to validate and optimize `.wasm` files.

* [Clang Format](https://github.com/tbfleming/cib) is used to format C/C++ files.

* [Cassowary.js](https://github.com/slightlyoff/cassowary.js/) is used to make split panes work.

* [Showdown](https://github.com/showdownjs/showdown) is used to automatically preview `.md` files.

* [Capstone.js](https://alexaltea.github.io/capstone.js/) is used to disassemble `x86` code.

* LLVM, Rust, Emscripten running server side.

* And of course: React, WebPack, TypeScript and TSLint.
