const crypto = require('crypto')

exports.getD3NScriptHash = (name, content) =>
  crypto
    .createHash('sha256')
    .update(Buffer.concat([Buffer.from(name), Buffer.from(content)]))
    .digest('hex')
