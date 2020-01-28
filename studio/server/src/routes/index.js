const compiler = require('./compiler')
const Router = require('@koa/router')
const AWS = require('../aws')
const axios = require('axios')

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
  const { code, hash } = ctx.request.body
  ctx.body = await AWS.uploadFile(hash, code)
})

module.exports = router
