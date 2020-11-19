abstract class ObiBase {
  abstract encode(value: any): Buffer
  abstract decode(buff: Buffer): any[]
}

export class ObiInteger extends ObiBase {
  static REGEX = /^(u|i)(8|16|32|64|128|256)$/
  isSigned: boolean
  sizeInBytes: number

  constructor(schema: string) {
    super()
    this.isSigned = schema[0] === 'i'
    this.sizeInBytes = parseInt(schema.slice(1)) / 8
  }

  encode(value: BigInt) {
    let newValue = BigInt(value)
    return Buffer.from(
      [...Array(this.sizeInBytes)]
        .map(() => {
          const dataByte = newValue % BigInt(1 << 8)
          newValue /= BigInt(1 << 8)
          return Number(dataByte)
        })
        .reverse(),
    )
  }

  decode(buff: Buffer): any {
    let value = BigInt(0)
    let multiplier = BigInt(1)
    for (let i = 0; i < this.sizeInBytes; i++) {
      value += BigInt(buff.readUInt8(this.sizeInBytes - i - 1)) * multiplier
      multiplier *= BigInt(1 << 8)
    }
    return [value, buff.slice(this.sizeInBytes)]
  }
}

export class ObiVector {
  static REGEX = /^\[.*\]$/
  internalObi: any

  constructor(schema: string) {
    this.internalObi = ObiSpec.fromSpec(schema.slice(1, -1))
  }

  encode(value: any) {
    return Buffer.concat([
      new ObiInteger('u32').encode(value.length),
      ...value.map((item: any) => this.internalObi.encode(item)),
    ])
  }

  decode(buff: Buffer) {
    let [length, remaining] = new ObiInteger('u32').decode(buff)
    let value = []
    for (let i = 0; i < Number(length); i++) {
      const decodeInternalResult = this.internalObi.decode(remaining)
      value.push(decodeInternalResult[0])
      remaining = decodeInternalResult[1]
    }
    return [value, remaining]
  }
}

export class ObiStruct {
  static REGEX = /^{.*}$/
  internalObiKvs: any

  constructor(schema: string) {
    this.internalObiKvs = []

    let curlyCount = 0
    let kv: any[] = ['', ''],
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

  encode(value: any) {
    return Buffer.concat(
      this.internalObiKvs.map(([k, obi]: any) => obi.encode(value[k])),
    )
  }

  decode(buff: Buffer) {
    let value: any = {}
    let remaining = buff
    for (let [k, obi] of this.internalObiKvs) {
      const decodeInternalResult = obi.decode(remaining)
      value[k] = decodeInternalResult[0]
      remaining = decodeInternalResult[1]
    }
    return [value, remaining]
  }
}

export class ObiString {
  static REGEX = /^string$/

  encode(value: string) {
    return Buffer.concat([
      new ObiInteger('u32').encode(BigInt(value.length)),
      Buffer.from(value),
    ])
  }

  decode(buff: Buffer) {
    let [length, remaining] = new ObiInteger('u32').decode(buff)
    return [
      remaining.slice(0, parseInt(length)).toString(),
      remaining.slice(parseInt(length)),
    ]
  }
}

export class ObiBytes {
  static REGEX = /^bytes$/

  encode(value: any) {
    return Buffer.concat([
      new ObiInteger('u32').encode(value.length),
      Buffer.from(value),
    ])
  }

  decode(buff: Buffer) {
    let [length, remaining] = new ObiInteger('u32').decode(buff)
    return [
      remaining.slice(0, parseInt(length)),
      remaining.slice(parseInt(length)),
    ]
  }
}

export class Obi {
  inputObi: ObiBase
  outputObi: ObiBase

  constructor(schema: string) {
    const normalizedSchema = schema.replace(/\s+/g, '')
    const tokens = normalizedSchema.split('/')
    this.inputObi = ObiSpec.fromSpec(tokens[0])
    this.outputObi = ObiSpec.fromSpec(tokens[1])
  }

  encodeInput(value: any) {
    return this.inputObi.encode(value)
  }

  decodeInput(buff: Buffer) {
    const [value, remaining] = this.inputObi.decode(buff)
    if (remaining.length != 0)
      throw new Error('Not all data is consumed after decoding output')
    return value
  }

  encodeOutput(value: any) {
    return this.outputObi.encode(value)
  }

  decodeOutput(buff: Buffer) {
    const [value, remaining] = this.outputObi.decode(buff)
    if (remaining.length != 0)
      throw new Error('Not all data is consumed after decoding output')
    return value
  }
}

export class ObiSpec {
  static impls = [ObiInteger, ObiVector, ObiStruct, ObiString, ObiBytes]

  static fromSpec(schema: string) {
    for (let impl of ObiSpec.impls) {
      if (schema.match(impl.REGEX)) {
        return new impl(schema)
      }
    }

    throw new Error(`No schema matched: <${schema}>`)
  }
}
