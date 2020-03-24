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

let decode: (string, string, JsBuffer.t) => option(array((string, string))) = [%bs.raw
  {|
function(_schema, cls, data) {
  const borsh = require('borsh')


  window[cls] = function(x,a) {
    this.x = x
    this.a = a
  }

  var model = window[cls]

  var instance = new model(1, "YO");

  const schema = new Map([
    [
      model,
      {
        kind: 'struct',
        fields: [
          ['x', 'U8'],
          ['a', 'String']
        ]
      }
    ]
  ])
  console.log(instance, model)
  let buf = borsh.serialize(schema, instance)
  console.log(buf)
  try {
    console.log(cls,data)
    let new_value = borsh.deserialize(schema, model, data)
    console.log(data)
    console.log(new_value)
    return schema.get(model).fields.map(([fieldName, _]) => {
        console.log(fieldName,new_value[fieldName])
        return [fieldName,new_value[fieldName]]
    });
  } catch {
    return null
  }
}
|}
];
