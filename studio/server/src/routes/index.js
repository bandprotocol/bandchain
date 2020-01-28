const compiler = require('./compiler')
const Router = require('@koa/router')
const AWS = require('../aws')
const axios = require('axios')
const { getD3NScriptHash } = require('../hash')

const router = new Router()

router.post('/compile', async (ctx, next) => {
  const { tar, options } = ctx.request.body
  ctx.body = await compiler.compile(tar, options)
})

router.post('/test', async (ctx, next) => {
  const { tar, options } = ctx.request.body
  ctx.body = await compiler.test(tar, options)
})

router.post('/upload', async (ctx, next) => {
  const { wasm, name, code } = ctx.request.body
  const wasmContent = Buffer.from(wasm, 'hex')
  const scriptHash = getD3NScriptHash(name, wasmContent)
  ctx.body = await AWS.uploadFile(scriptHash, code)
})

/** TODO: Remove once CORS on bandsv works */
router.post('/execute', async (ctx, next) => {
  const { code, params } = ctx.request.body

  ctx.body = (
    await axios.post('http://d3n-debug.bandprotocol.com:5000/execute', {
      code,
      params,
    })
  ).data
})

/** TODO: Remove once CORS on bandsv works */
router.post('/params-info', async (ctx, next) => {
  const { code } = ctx.request.body

  ctx.body = (
    await axios.post('http://d3n-debug.bandprotocol.com:5000/params-info', {
      code,
    })
  ).data
})

/** TODO: Remove once CORS on bandsv works */
router.post('/store', async (ctx, next) => {
  const { code, name } = ctx.request.body

  ctx.body = (
    await axios.post('http://d3n-debug.bandprotocol.com:5000/store', {
      code,
      name,
    })
  ).data
})

module.exports = router
