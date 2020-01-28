const { wasmPackCmd, cargoCmd, tempDir } = require('../config')
const {
  exec,
  execFile,
  joinCmd,
  writeFile,
  readFile,
  mkdir,
} = require('../system')

exports.compile = async (tar, options = {}) => {
  // Each session has its own folder which we'll keep overwrite with new version
  const crateName =
    'rustc_h_' +
    (options.session.replace(/[\.\/]/g, '') ||
      Math.random()
        .toString(36)
        .slice(2))

  const crateDir = tempDir + '/' + crateName

  try {
    await mkdir(crateDir)
  } catch (err) {
    // Folder alread exist, which is fine
  }

  const rustTar = crateDir + '/' + 'lib.tar'
  const wasmFile = crateDir + '/' + 'pkg/main_bg.wasm'

  await writeFile(rustTar, Buffer.from(tar, 'base64').toString('ascii'))
  await exec(joinCmd(['tar', 'xvf', rustTar, '-C', crateDir]))

  try {
    const args = ['build', '--out-name=main', crateDir]
    const opts = {}

    return {
      success: true,
      message: await execFile(wasmPackCmd, args, opts),
      output: (await readFile(wasmFile)).toString('base64'),
    }
  } catch (e) {
    return { success: false, message: 'Error: ' + e, output: '' }
  }
}

exports.test = async (tar, options = {}) => {
  // Each session has its own folder which we'll keep overwrite with new version
  const crateName =
    'rustc_h_' +
    (options.session.replace(/[\.\/]/g, '') ||
      Math.random()
        .toString(36)
        .slice(2))

  const crateDir = tempDir + '/' + crateName

  try {
    await mkdir(crateDir)
  } catch (err) {
    // Folder alread exist, which is fine
  }

  const rustTar = crateDir + '/' + 'lib.tar'

  await writeFile(rustTar, Buffer.from(tar, 'base64').toString('ascii'))
  await exec(joinCmd(['tar', 'xvf', rustTar, '-C', crateDir]))

  try {
    const args = [
      'test',
      '--manifest-path',
      crateDir + '/Cargo.toml',
      '--',
      '--nocapture',
    ]

    const opts = {}

    return {
      success: true,
      message: await execFile(cargoCmd, args, opts),
    }
  } catch (e) {
    return { success: false, message: 'Error: ' + e }
  }
}
