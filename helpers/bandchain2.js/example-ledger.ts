import { Wallet } from './src/index'

const { Ledger } = Wallet

async function ledger() {
  const ledger = await Ledger.connectLedgerNode()
  console.log(ledger)
}

ledger()