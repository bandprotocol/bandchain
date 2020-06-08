type field_key_type_t = {
  fieldName: string,
  fieldType: string,
};

type field_key_value_t = {
  fieldName: string,
  fieldValue: string,
};

let extractFields: (string, string) => option(array(field_key_type_t)) = [%bs.raw
  {|
  function(schema, t) {
    console.log(schema, t)
    try {
      const normalizedSchema = schema.replace(/\s+/g, '')
      const tokens = normalizedSchema.split('/')
      let val
      if (t === 'input') {
        val = tokens[0]
      } else if (t === 'output') {
        val = tokens[1]
      } else {
        return undefined
      }
      let specs = val.slice(1, val.length - 1).split(',')
      return specs.map((spec) => {
        let x = spec.split(':')
        return {fieldName: x[0], fieldType: x[1]}
      })
    } catch {
      return undefined
    }
  }
|}
];

let encode: (string, string, array(field_key_value_t)) => option(JsBuffer.t) = [%bs.raw
  {|
function(schema, t, data) {
  const { Obi } = require('obi.js')

  try {
    const obi = new Obi(schema)

    let payload = {}
    for (let x of data) {
      let value = x.fieldValue
      if (x.fieldValue.startsWith("[") || x.fieldValue.startsWith("{")) {
        value = JSON.parse(x.fieldValue)
      } else if (x.fieldValue.startsWith("0x")) {
        value = Buffer.from(x.fieldValue.slice(2), "hex")
      }
      payload = {...payload, [x.fieldName]: value}
    }

    let buf
    if (t === 'input') {
      buf = obi.encodeInput(payload)
    } else if (t === 'output') {
      buf = obi.encodeOutput(payload)
    } else {
      return undefined
    }

    return Buffer.from(buf)
  } catch(err) {
    console.error(`Error encode ${t}`, err)
    return undefined
  }
}
  |}
];

let decode: (string, string, JsBuffer.t) => option(array(field_key_value_t)) = [%bs.raw
  {|
function(schema, t, data) {
  const { Obi } = require('obi.js')

  function stringify(data) {
    if (Array.isArray(data)) {
      return "[" + [...data].map(stringify).join(",") + "]"
    } else if (typeof(data) === "bigint") {
      return data.toString()
    } else if (Buffer.isBuffer(data)) {
      return "0x" + data.toString('hex')
    } else if (typeof(data) === "object") {
      return "{" + Object.entries(data).map(([k,v]) => JSON.stringify(k)+ ":" + stringify(v)).join(",") + "}"
    } else {
      return JSON.stringify(data)
    }
  }

  try {
    const obi = new Obi(schema)
    let rawResult
    if (t === 'input') {
      rawResult = obi.decodeInput(data)
    } else if (t === 'output') {
      rawResult = obi.decodeOutput(data)
    } else {
      return undefined
    }

    let result = []
    for (let x of Object.entries(rawResult)) {
      let value = stringify(x[1])
      result = [...result, {fieldName: x[0], fieldValue: value}]
    }

    return result
  } catch(err) {
    console.error(`Error decode ${t}`, err)
    return undefined
  }
}
  |}
];
