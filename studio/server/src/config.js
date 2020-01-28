const homeDir = process.env['HOME']

exports.cargoCmd = homeDir + '/.cargo/bin/cargo'
exports.wasmPackCmd = homeDir + '/.cargo/bin/wasm-pack'
exports.wasmRunCmd = __dirname + '/wasm-run'
exports.tempDir = '/tmp'
