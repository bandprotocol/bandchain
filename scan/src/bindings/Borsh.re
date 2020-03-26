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
  let newValue = borsh.deserialize(schema, Test, Buffer.from(buf))
  console.log(buf)
  console.log(newValue)
}
|}
];

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
    let newValue = borsh.deserialize(schemaMap, model, data)
    return schemaMap.get(model).fields.map(([fieldName, _]) => {
      return [fieldName, newValue[fieldName].toString()]
    });
  } catch(err) {
    return undefined
  }
}
|}
];
