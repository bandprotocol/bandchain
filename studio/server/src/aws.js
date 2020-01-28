const AWS = require('aws-sdk')

const s3 = new AWS.S3({
  accessKeyId: process.env.AWS_ACCESS_KEY || 'AKIAJKUU2UFOS2DH3VEQ',
  secretAccessKey:
    process.env.AWS_SECRET_ACCESS_KEY ||
    '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',
})

exports.uploadFile = (name, content) => {
  const params = {
    Bucket: process.env.AWS_BUCKET || 'code.d3n.bandprotocol.com',
    Key: name,
    Body: content,
    ContentType: 'application/json',
  }

  return new Promise((resolve, reject) => {
    s3.upload(params, (err, data) => {
      if (err) reject(err)
      resolve(data.Location)
    })
  })
}
