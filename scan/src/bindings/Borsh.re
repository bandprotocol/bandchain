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

// Test by this test
// Borsh.decode(
//   {j|{"Input": "{ \\"kind\\": \\"struct\\", \\"fields\\": [ [\\"symbol\\", \\"string\\"], [\\"multiplier\\", \\"u64\\"], [\\"what\\", \\"u8\\"] ] }"}|j},
//   "Input",
//   "0x03000000425443320000000000000064" |> JsBuffer.fromHex,
// );
let decode: (string, string, JsBuffer.t) => option(array((string, string))) = [%bs.raw
  {|
function(_schema, cls, data) {
  const borsh = require('borsh')

  function gen(members) {
    return function(...args) {
      for (let i in members) {
        this[members[i]] = args[i]
      }
    }
  }

  try {
    let schema = JSON.parse(_schema)
    let schemaMap = new Map()
    for (let className in schema) {
      let t = JSON.parse(schema[className])
      if (t.kind == "struct") {
        window[className] = gen(t.fields.map(x => x[0]))
        t.fields = t.fields.map(x => {
          x[1] = x[1].replace(/^\w/, c => c.toUpperCase());
          return x
        })
        schemaMap.set(window[className], t)
      }
    }
    let model = window[cls]
    let new_value = borsh.deserialize(schemaMap, model, data)
    return schemaMap.get(model).fields.map(([fieldName, _]) => {
        return [fieldName, JSON.stringify(new_value[fieldName])]
    });
  } catch(err) {
    return undefined
  }
}
|}
];
