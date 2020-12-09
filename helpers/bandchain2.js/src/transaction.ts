import { Client } from 'index'
import { Msg } from './message'
import { PublicKey } from './wallet'

export default class Transaction {
  msgs: Msg[] = []
  accountNum?: number
  sequence?: number
  chainID?: string
  fee: number = 0
  gas: number = 200000
  memo: string = ''

  withMessages(...msg: Msg[]): Transaction {
    this.msgs.push(...msg)
    return this
  }

  async withAuto(client: Client): Promise<Transaction> {
    if (this.msgs.length == 0)
      throw Error(
        'Message is empty, please use withMessages at least 1 message.',
      )

    let addr = this.msgs[0].getSender()
    let account = await client.getAccount(addr)
    if (!account) throw Error(`Account doesn't exist.`)
    this.accountNum = account.accountNumber
    this.sequence = account.sequence
    return this
  }

  withAccountNum(accountNum: number): Transaction {
    if (!Number.isInteger(accountNum)) {
      throw Error('accountNum is not an integer')
    }
    this.accountNum = accountNum
    return this
  }

  withSequence(sequence: number): Transaction {
    if (!Number.isInteger(sequence)) {
      throw Error('sequence is not an integer')
    }
    this.sequence = sequence
    return this
  }

  withChainID(chainID: string): Transaction {
    this.chainID = chainID
    return this
  }

  withFee(fee: number): Transaction {
    if (!Number.isInteger(fee)) {
      throw Error('fee is not an integer')
    }
    this.fee = fee
    return this
  }

  withGas(gas: number): Transaction {
    if (!Number.isInteger(gas)) {
      throw Error('gas is not an integer')
    }
    this.gas = gas
    return this
  }

  withMemo(memo: string): Transaction {
    this.memo = memo
    return this
  }

  getSignData(): Buffer {
    if (this.msgs.length == 0) {
      throw Error('message is empty')
    }

    if (this.accountNum == null) {
      throw Error('accountNum should be defined')
    }

    if (this.sequence == null) {
      throw Error('sequence should be defined')
    }

    if (this.chainID == null) {
      throw Error('chainID should be defined')
    }

    this.msgs.forEach((msg) => msg.validate())

    let messageJson: { [key: string]: any } = {
      chain_id: this.chainID,
      account_number: this.accountNum.toString(),
      fee: {
        amount: [
          {
            amount: this.fee.toString(),
            denom: 'uband',
          },
        ],
        gas: this.gas.toString(),
      },
      memo: this.memo,
      sequence: this.sequence.toString(),
      msgs: this.msgs.map((msg) => msg.asJson()),
    }

    const sortedKey = Object.keys(messageJson).sort()
    const result: { [key: string]: any } = {}
    sortedKey.forEach((key) => (result[key] = messageJson[key]))

    return Buffer.from(JSON.stringify(result))
  }

  getTxData(signature: Buffer, pubkey: PublicKey): Object {
    if (this.accountNum == null) {
      throw Error('accountNum should be defined')
    }

    if (this.sequence == null) {
      throw Error('sequence should be defined')
    }

    return {
      fee: {
        amount: [{ amount: this.fee.toString(), denom: 'uband' }],
        gas: this.gas.toString(),
      },
      memo: this.memo,
      msg: this.msgs.map((msg) => msg.asJson()),
      signatures: [
        {
          signature: signature.toString('base64'),
          pub_key: {
            type: 'tendermint/PubKeySecp256k1',
            value: Buffer.from(pubkey.toHex(), 'hex').toString('base64'),
          },
          account_number: this.accountNum.toString(),
          sequence: this.sequence.toString(),
        },
      ],
    }
  }
}
