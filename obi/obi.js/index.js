class ObiSpec {
  static impls = []

  static fromSpec(schema) {
    for (let impl of ObiSpec.impls) {
      if (schema.match(impl.REGEX)) {
        return new impl(schema)
      }
    }

    throw new Error(`No schema matched: <${schema}>`)
  }
}

class ObiInteger {
  static REGEX = /^(u|i)(8|16|32|64|128|256)$/

  constructor(schema) {
    this.isSigned = schema[0] === 'i'
    this.sizeInBytes = parseInt(schema.slice(1)) / 8
  }

  encode(value) {
    value = BigInt(value)
    return Buffer.from(
      [...Array(this.sizeInBytes)].map(() => {
        const byte = value % BigInt(1 << 8)
        value /= BigInt(1 << 8)
        return parseInt(byte)
      }),
    )
  }

  decode(buff) {
    let value = BigInt(0)
    let multiplier = BigInt(1)
    for (let i = 0; i < this.sizeInBytes; i++) {
      value += BigInt(buff.readUInt8(i)) * multiplier
      multiplier = multiplier * BigInt(1 << 8)
    }
    return [value, buff.slice(this.sizeInBytes)]
  }
}

class ObiVector {
  static REGEX = /^\[.*\]$/
  constructor(schema) {
    this.internalObi = ObiSpec.fromSpec(schema.slice(1, -1))
  }
  encode(value) {
    return Buffer.concat([
      new ObiInteger('u32').encode(value.length),
      ...value.map((item) => this.internalObi.encode(item)),
    ])
  }
  decode(buff) {
    let [length, remaining] = new ObiInteger('u32').decode(buff)
    let value = []
    for (let i = 0; i < parseInt(length); i++) {
      const decodeInternalResult = this.internalObi.decode(remaining)
      value.push(decodeInternalResult[0])
      remaining = decodeInternalResult[1]
    }
    return [value, remaining]
  }
}

class ObiStruct {
  static REGEX = /^{.*}$/
  constructor(schema) {
    this.internalObiKvs = []

    let curlyCount = 0
    let kv = ['', ''],
      fill = 0 // 0 = k, 1 = v
    for (let c of schema.slice(1)) {
      if (c == '{') curlyCount++
      else if (curlyCount && c == '}') curlyCount--
      else if (!curlyCount && c === ':') {
        fill = 1
        continue
      } else if (!curlyCount && (c === ',' || c === '}')) {
        kv[1] = ObiSpec.fromSpec(kv[1])
        this.internalObiKvs.push(kv)
        kv = ['', '']
        fill = 0
        continue
      }

      kv[fill] += c
    }
  }
  encode(value) {
    return Buffer.concat(
      this.internalObiKvs.map(([k, obi]) => obi.encode(value[k])),
    )
  }
  decode(buff) {
    let value = {}
    let remaining = buff
    for (let i = 0; i < this.internalObiKvs.length; i++) {
      const [k, obi] = this.internalObiKvs[i]
      const decodeInternalResult = obi.decode(remaining)
      value[k] = decodeInternalResult[0]
      remaining = decodeInternalResult[1]
    }
    return [value, remaining]
  }
}

class ObiString {
  static REGEX = /^string$/
  constructor(schema) {}
  encode(value) {
    return Buffer.concat([
      new ObiInteger('u32').encode(value.length),
      Buffer.from(value),
    ])
  }
  decode(buff) {
    let [length, remaining] = new ObiInteger('u32').decode(buff)
    return [
      remaining.slice(0, parseInt(length)).toString(),
      remaining.slice(parseInt(length)),
    ]
  }
}

class ObiBytes {
  static REGEX = /^bytes$/
  constructor(schema) {
    this.internalObiKvs
  }
  encode(value) {
    return Buffer.concat([
      new ObiInteger('u32').encode(value.length),
      Buffer.from(value),
    ])
  }
  decode(buff) {
    let [length, remaining] = new ObiInteger('u32').decode(buff)
    return [
      remaining.slice(0, parseInt(length)),
      remaining.slice(parseInt(length)),
    ]
  }
}

class Obi {
  constructor(schema) {
    const normalizedSchema = schema.replace(/\s+/g, '')
    const tokens = normalizedSchema.split('/')
    this.inputObi = ObiSpec.fromSpec(tokens[0])
    this.outputObi = ObiSpec.fromSpec(tokens[1])
  }

  encodeInput(value) {
    return this.inputObi.encode(value)
  }

  decodeInput(buff) {
    const [value, remaining] = this.inputObi.decode(buff)
    if (remaining)
      throw new Error('Not all data is consumed after decoding output')
    return value
  }

  encodeOutput(value) {
    return this.outputObi.encode(value)
  }

  decodeOutput(buff) {
    const [value, remaining] = this.outputObi.decode(buff)
    if (remaining)
      throw new Error('Not all data is consumed after decoding output')
    return value
  }
}

ObiSpec.impls = [ObiInteger, ObiVector, ObiStruct, ObiString, ObiBytes]

module.exports = {
  Obi,
  ObiSpec,
  ObiInteger,
  ObiVector,
  ObiStruct,
  ObiString,
  ObiBytes,
}
