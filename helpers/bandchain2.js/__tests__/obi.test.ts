import { assert } from 'console'
import {
  Obi,
  ObiBytes,
  ObiInteger,
  ObiString,
  ObiStruct,
  ObiVector,
} from '../src/obi'

describe('ObiInteger', () => {
  it('encode / decode successfully', () => {
    const data = BigInt(12345678)
    const encode = new ObiInteger('u256').encode(data)
    const decode = new ObiInteger('u256').decode(encode)[0]

    expect(decode).toEqual(data)
  })
})

describe('ObiString', () => {
  it('encode / decode successfully', () => {
    const data = 'everything is awesome'
    const encode = new ObiString().encode(data)
    const decode = new ObiString().decode(encode)[0]
    expect(decode).toEqual(data)
  })
})

describe('ObiVector', () => {
  it('encode / decode successfully', () => {
    const data = [1, 2, 3]
    const encode = new ObiVector('[u8]').encode(data)
    const decode = new ObiVector('[u8]').decode(encode)[0]
    expect(decode.map((each) => Number(each))).toEqual(data)
  })
})

describe('ObiStruct', () => {
  it('encode / decode successfully', () => {
    const data = { a: BigInt(1), b: 'hi' }
    const encode = new ObiStruct('{a:u8,b:string}').encode(data)
    const decode = new ObiStruct('{a:u8,b:string}').decode(encode)[0]
    expect(decode).toEqual(data)
  })
})

describe('ObiBytes', () => {
  it('encode / decode successfully', () => {
    const data = Buffer.from([1, 2, 3, 4, 5])
    const encode = new ObiBytes().encode(data)
    const decode = new ObiBytes().decode(encode)[0]
    expect(decode).toEqual(data)
  })
})

describe('Obi', () => {
  const obi = new Obi(`
  {
    symbol: string,
    multiplier: u64
  } / {
    price: u64,
    sources: [{ name: string, time: u64 }]
  }
`)

  it('encode input successfully', () => {
    expect(
      obi
        .encodeInput({ symbol: 'BTC', multiplier: BigInt('1000000000') })
        .toString('hex'),
    ).toEqual('00000003425443000000003b9aca00')
  })

  it('decode input successfully', () => {
    expect(
      obi.decodeInput(Buffer.from('00000003425443000000003b9aca00', 'hex')),
    ).toEqual({
      symbol: 'BTC',
      multiplier: BigInt('1000000000'),
    })
  })

  it('encode output successfully', () => {
    expect(
      obi
        .encodeOutput({
          price: 9268300000000,
          sources: [
            { name: 'CoinGecko', time: 1590305341 },
            { name: 'CryptoCompare', time: 1590305362 },
          ],
        })
        .toString('hex'),
    ).toEqual(
      '0000086df1baab000000000200000009436f696e4765636b6f000000005eca223d0000000d43727970746f436f6d70617265000000005eca2252',
    )
  })

  it('decode successfully', () => {
    expect(
      obi.decodeOutput(
        Buffer.from(
          '0000086df1baab000000000200000009436f696e4765636b6f000000005eca223d0000000d43727970746f436f6d70617265000000005eca2252',
          'hex',
        ),
      ),
    ).toEqual({
      price: BigInt(9268300000000),
      sources: [
        { name: 'CoinGecko', time: BigInt(1590305341) },
        { name: 'CryptoCompare', time: BigInt(1590305362) },
      ],
    })
  })
})
