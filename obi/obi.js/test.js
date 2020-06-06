const {
  Obi,
  ObiSpec,
  ObiInteger,
  ObiVector,
  ObiStruct,
  ObiString,
  ObiBytes,
} = require('./index.js')

const check = (desc, fn) => console.log(`${fn() ? '✅' : '⛔️'} ${desc}`)

;(() => {
  const data = BigInt('12345678')
  const encode = new ObiInteger('u256').encode(data)
  const decode = new ObiInteger('u256').decode(encode)[0]
  check(`Encode/decode u256 with ${data}`, () => data == decode)
})()
;(() => {
  const data = 'everything is awesome'
  const encode = new ObiString().encode(data)
  const decode = new ObiString().decode(encode)[0]
  check(`Encode/decode string with "${data}"`, () => data == decode)
})()
;(() => {
  const data = [BigInt(1), BigInt(2), BigInt(3)]
  const encode = new ObiVector('[u8]').encode(data)
  const decode = new ObiVector('[u8]').decode(encode)[0]
  check(
    `Encode/decode [u8] with ${data}`,
    () => data[0] == decode[0] && data[1] == decode[1] && data[2] == decode[2],
  )
})()
;(() => {
  const data = { a: 1, b: 'hi' }
  const encode = new ObiStruct('{a:u8,b:string}').encode(data)
  const decode = new ObiStruct('{a:u8,b:string}').decode(encode)[0]
  check(
    `Encode/decode {a:u8,b:string} with ${JSON.stringify(data)}`,
    () => data.a == decode.a && data.b == decode.b,
  )
})()
;(() => {
  const data = Buffer.from([1, 2, 3, 4, 5])
  const encode = new ObiBytes().encode(data)
  const decode = new ObiBytes().decode(encode)[0]
  check(`Encode/decode bytes with ${data}`, () => data.equals(decode))
})()
;(() => {
  const obi = new Obi(`
  {
    symbol: string,
    multiplier: u64
  } / {
    price: u64,
    sources: [{ name: string, time: u64 }]
  }
`)

  check('Encode example output', () =>
    Buffer.from(
      '0000086df1baab000000000200000009436f696e4765636b6f000000005eca223d0000000d43727970746f436f6d70617265000000005eca2252',
      'hex',
    ).equals(
      obi.encodeOutput({
        price: 9268300000000,
        sources: [
          { name: 'CoinGecko', time: 1590305341 },
          { name: 'CryptoCompare', time: 1590305362 },
        ],
      }),
    ),
  )

  check('Encode example input', () =>
    Buffer.from('00000003425443000000003b9aca00', 'hex').equals(
      obi.encodeInput({ symbol: 'BTC', multiplier: BigInt('1000000000') }),
    ),
  )
})()
