const Koa = require('koa')
const bodyParse = require('koa-bodyparser')
const cors = require('@koa/cors')

const router = require('./src/routes')

const app = new Koa()
app
  .use(cors())
  .use(bodyParse())
  .use(router.routes())
  .use(router.allowedMethods())
app.listen(8082)

console.log('OWASM Studio server is listening at http://localhost:8082')
