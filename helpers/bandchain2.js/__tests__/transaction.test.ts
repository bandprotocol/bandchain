import { Message, Data, Transaction, Wallet } from '../src/index'
import { Address } from '../src/wallet'

const { MsgSend } = Message
const { Coin } = Data
const { PrivateKey } = Wallet

describe('Transaction', () => {
  const coin = new Coin(100000, 'uband')
  const msgSend = new MsgSend(
    Address.fromAccBech32('band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte'),
    Address.fromAccBech32('band1ksnd0f3xjclvg0d4z9w0v9ydyzhzfhuy47yx79'),
    [coin],
  )

  it('create successfully', () => {
    const tsc = new Transaction()

    expect(tsc.accountNum).toBeUndefined()
    expect(tsc.chainID).toBeUndefined()
    expect(tsc.sequence).toBeUndefined()
    expect(tsc.fee).toEqual(0)
    expect(tsc.gas).toEqual(200000)
    expect(tsc.withMessages.length).toEqual(0)
  })

  it('able to work "with" function', () => {
    const tsc = new Transaction()
      .withMessages(msgSend)
      .withAccountNum(100)
      .withSequence(30)
      .withChainID('bandchain')
      .withGas(500000)
      .withFee(10)
      .withMemo('bandchain2.js test')

    expect(tsc.msgs.length).toEqual(1)
    expect(tsc.accountNum).toEqual(100)
    expect(tsc.sequence).toEqual(30)
    expect(tsc.chainID).toEqual('bandchain')
    expect(tsc.fee).toEqual(10)
    expect(tsc.gas).toEqual(500000)
    expect(tsc.memo).toEqual('bandchain2.js test')
  })

  it('get error from checking integer', () => {
    const tsc = new Transaction().withMessages(msgSend)

    expect(() => {
      tsc.withAccountNum(100.5)
    }).toThrowError('accountNum is not an integer')

    expect(() => {
      tsc.withSequence(100.5)
    }).toThrowError('sequence is not an integer')

    expect(() => {
      tsc.withFee(100.5)
    }).toThrowError('fee is not an integer')

    expect(() => {
      tsc.withGas(100.5)
    }).toThrowError('gas is not an integer')
  })

  it('getSignData successfully', () => {
    const tsc = new Transaction()
      .withMessages(msgSend)
      .withAccountNum(100)
      .withSequence(30)
      .withChainID('bandchain')
      .withGas(500000)
      .withFee(10)
      .withMemo('bandchain2.js test')

    let expectedResult =
      '7b226163636f756e745f6e756d626572223a22313030222c22636861696e5f6964223a2262616e64636861696e222c22666565223a7b22616d6f756e74223a5b7b22616d6f756e74223a223130222c2264656e6f6d223a227562616e64227d5d2c22676173223a22353030303030227d2c226d656d6f223a2262616e64636861696e322e6a732074657374222c226d736773223a5b7b2274797065223a22636f736d6f732d73646b2f4d736753656e64222c2276616c7565223a7b22616d6f756e74223a5b7b22616d6f756e74223a22313030303030222c2264656e6f6d223a227562616e64227d5d2c2266726f6d5f61646472657373223a2262616e64316a7268757172796d7a74346d6e766777386376793373397a6878336a6a306471333071707465222c22746f5f61646472657373223a2262616e64316b736e64306633786a636c76673064347a39773076397964797a687a66687579343779783739227d7d5d2c2273657175656e6365223a223330227d'

    expect(tsc.getSignData().toString('hex')).toEqual(expectedResult)
  })

  it('getSignData fail', () => {
    let tsc = new Transaction()

    expect(() => tsc.getSignData()).toThrowError('message is empty')

    tsc = tsc.withMessages(msgSend)
    expect(() => tsc.getSignData()).toThrowError('accountNum should be defined')

    tsc = tsc.withAccountNum(100)
    expect(() => tsc.getSignData()).toThrowError('sequence should be defined')

    tsc = tsc.withSequence(30)
    expect(() => tsc.getSignData()).toThrowError('chainID should be defined')

    tsc = tsc.withChainID('bandchain')
    let expectedResult =
      '7b226163636f756e745f6e756d626572223a22313030222c22636861696e5f6964223a2262616e64636861696e222c22666565223a7b22616d6f756e74223a5b7b22616d6f756e74223a2230222c2264656e6f6d223a227562616e64227d5d2c22676173223a22323030303030227d2c226d656d6f223a22222c226d736773223a5b7b2274797065223a22636f736d6f732d73646b2f4d736753656e64222c2276616c7565223a7b22616d6f756e74223a5b7b22616d6f756e74223a22313030303030222c2264656e6f6d223a227562616e64227d5d2c2266726f6d5f61646472657373223a2262616e64316a7268757172796d7a74346d6e766777386376793373397a6878336a6a306471333071707465222c22746f5f61646472657373223a2262616e64316b736e64306633786a636c76673064347a39773076397964797a687a66687579343779783739227d7d5d2c2273657175656e6365223a223330227d'
    expect(tsc.getSignData().toString('hex')).toEqual(expectedResult)
  })

  it('getTxData successfully', () => {
    const tsc = new Transaction()
      .withMessages(msgSend)
      .withAccountNum(100)
      .withSequence(30)
      .withChainID('bandchain')
      .withGas(200000)

    const rawData = tsc.getSignData()
    const privKey = PrivateKey.fromMnemonic('s')
    const pubkey = privKey.toPubkey()
    const signature = privKey.sign(rawData)
    const rawTx = tsc.getTxData(signature, pubkey)

    expect(rawTx).toEqual({
      msg: [
        {
          type: 'cosmos-sdk/MsgSend',
          value: {
            amount: [
              {
                amount: '100000',
                denom: 'uband',
              },
            ],
            from_address: 'band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte',
            to_address: 'band1ksnd0f3xjclvg0d4z9w0v9ydyzhzfhuy47yx79',
          },
        },
      ],
      fee: { gas: '200000', amount: [{ denom: 'uband', amount: '0' }] },
      memo: '',
      signatures: [
        {
          signature:
            'LLnpUzoB7gmQd+9XxQiRyiEv3O04FWWgyX2Agm2xOoMF0iwab3BzG1L5Szl0OGfJezmMm016gF/Gjy0l9niB5w==',
          pub_key: {
            type: 'tendermint/PubKeySecp256k1',
            value: 'A/5wi9pmUk/SxrzpBoLjhVWoUeA9Ku5PYpsF3pD1Htm8',
          },
          account_number: '100',
          sequence: '30',
        },
      ],
    })
  })
})
