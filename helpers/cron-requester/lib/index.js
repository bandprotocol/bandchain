const BandChain = require('@bandprotocol/bandchain.js')
const { CronJob } = require('cron')
const fs = require('fs')

// Due to how bandchain.js is written, we cannot batch all the requests in one tx yet
// TODO: Fix this
async function runJob(bandchain, mnemonic, validatorCounts, requests) {
  let count = 0
  for (let request of requests) {
    try {
      const requestId = await bandchain.submitRequestTx(
        request.oracleScript,
        request.params,
        validatorCounts,
        mnemonic,
      )

      count++
      console.log(
        '∟ ✅ requestId = %s | oracleScript #%d %s',
        requestId,
        request.oracleScriptId,
        JSON.stringify(request.params),
      )
    } catch {
      console.log(
        '∟ ⛔️ request failed | oracleScript #%d %s',
        request.oracleScriptId,
        JSON.stringify(request.params),
      )
    }
  }
  console.log(
    '%s [%d/%d] requests was submitted',
    count === requests.length ? '⛳️' : '👎',
    count,
    requests.length,
  )
  console.log('--------------------------------------------------------')
}

async function start(configFilePath) {
  let config = {}

  try {
    // Load config file
    const configFile = fs.readFileSync(configFilePath)
    config = JSON.parse(configFile)
  } catch {
    throw new Error('Incorrect config file path/format')
  }

  const {
    mnemonic,
    chainId,
    endpoint,
    cronPattern,
    validatorCounts,
    requests,
  } = config

  // Check config file content format
  if (typeof mnemonic !== 'string')
    throw new Error('config.mnemonic has to be string')
  if (typeof chainId !== 'string')
    throw new Error('config.chainId has to be string')
  if (typeof endpoint !== 'string')
    throw new Error('config.endpoint has to be string')
  if (typeof cronPattern !== 'string')
    throw new Error('config.cronPattern has to be string')
  if (typeof validatorCounts !== 'object')
    throw new Error('config.validatorCounts has to be an object')
  if (!Array.isArray(requests))
    throw new Error('config.requests has to be an array')

  // Instantiate BandChain object with the specified chain ID And REST Endpoint
  // TODO: Remove dependencies on chainId once #1951 goes live on devnet and bandchain.js supports it
  const bandchain = new BandChain(chainId, endpoint)

  // Format requests
  const formattedRequests = await Promise.all(
    requests.map(async (request) => {
      return {
        oracleScript: await bandchain.getOracleScript(request.oracleScriptId),
        ...request,
      }
    }),
  )

  // Start cronjob
  const cronJob = new CronJob(
    cronPattern,
    () => {
      console.log('⏰ Requests start at %s', new Date().toLocaleString())
      runJob(bandchain, mnemonic, validatorCounts, formattedRequests)
    },
    null,
    true,
  )

  // Log the start of the program
  console.log('--------------------------------------------------------')
  console.log(
    '⭐️ Cron is running! Your requests will be executed with the cron pattern %s',
    cronPattern,
  )
  console.log('📆 Your first requests will start at %s', cronJob.nextDates())
  console.log('--------------------------------------------------------')
}

module.exports = start
