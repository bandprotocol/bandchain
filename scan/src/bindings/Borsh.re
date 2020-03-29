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

let encode: (string, string, array((string, string))) => option(JsBuffer.t) = [%bs.raw
  {|
function(_schema, cls, data) {
  const borsh = require('borsh')

  function gen(members) {
    return function(args) {
      for (let i in members) {
        this[members[i]] = args[members[i]]
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

    let rawValue = {}
    let specs = schemaMap.get(window[cls])

    if (!specs.fields) return undefined

    for (let i in data) {
      let isFound = specs.fields.some(x => x[0] == data[i][0])
      if (!isFound) return undefined

      rawValue[data[i][0]] = data[i][1]
    }
    let value = new window[cls](rawValue)

    if (specs.fields.length !== data.length) return undefined

    let buf = borsh.serialize(schemaMap, value)
    return Buffer.from(buf)
  } catch(err) {
    return undefined
  }
}
|}
];
