let serialize: unit => unit = [%bs.raw
  {|
function() {
  const borsh = require('borsh')
  class Test {
    constructor(x, y, z, q) {
      this.x = x
      this.y = y
      this.z = z
      this.q = q
    }
  }

  const value = new Test(255, 20, '123', [1, 2, 3])

  const schema = new Map([
    [
      Test,
      {
        kind: 'struct',
        fields: [
          ['x', 'U8'],
          ['y', 'U64'],
          ['z', 'String'],
          ['q', [3]]
        ]
      }
    ]
  ])
  let buf = borsh.serialize(schema, value)
  let new_value = borsh.deserialize(schema, Test, Buffer.from(buf))
  console.log(buf)
  console.log(new_value)
}
|}
];
