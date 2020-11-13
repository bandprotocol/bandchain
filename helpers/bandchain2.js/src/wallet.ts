import * as bip39 from 'bip39'
import * as bip32 from 'bip32'
import secp256k1 from 'secp256k1'
import crypto from 'crypto'
import { ECPair } from 'bitcoinjs-lib'

const DEFAULT_DERIVATION_PATH = "m/44'/494'/0'/0/0"

export class PrivateKey {
  private signingKey: Buffer

  private constructor(signingKey: Buffer) {
    this.signingKey = signingKey
  }

  static generate(path = DEFAULT_DERIVATION_PATH): [string, PrivateKey] {
    const phrase = bip39.generateMnemonic(256)
    return [phrase, this.fromMnemonic(phrase, path)]
  }

  static fromMnemonic(
    words: string,
    path = DEFAULT_DERIVATION_PATH,
  ): PrivateKey {
    const seed = bip39.mnemonicToSeedSync(words)
    const node = bip32.fromSeed(seed)
    const child = node.derivePath(path)

    if (!child.privateKey) throw new Error('Cannot create private key')
    const ecpair = ECPair.fromPrivateKey(child.privateKey, {
      compressed: false,
    })

    if (!ecpair.privateKey) throw new Error('Cannot create private key')
    return new PrivateKey(ecpair.privateKey)
  }

  static fromHex(priv: string): PrivateKey {
    return new PrivateKey(Buffer.from(priv, 'hex'))
  }

  toHex(): string {
    return this.signingKey.toString('hex')
  }

  toPubkey() {
    const pubKeyByte = secp256k1.publicKeyCreate(this.signingKey)
    return PublicKey.fromHex(Buffer.from(pubKeyByte).toString('hex'))
  }

  sign(msg: Buffer): Buffer {
    const hash = crypto.createHash('sha256').update(msg).digest('hex')
    const buf = Buffer.from(hash, 'hex')
    const { signature } = secp256k1.ecdsaSign(buf, this.signingKey)
    return Buffer.from(signature)
  }
}

export class PublicKey {
  private verifyKey: Buffer

  private constructor(verifyKey: Buffer) {
    this.verifyKey = verifyKey
  }

  static fromHex(pub: string): PublicKey {
    return new PublicKey(Buffer.from(pub, 'hex'))
  }

  toHex(): string {
    return this.verifyKey.toString('hex')
  }

  verify(msg: Buffer, sig: Buffer): boolean {
    const hash = crypto.createHash('sha256').update(msg).digest('hex')
    const buf = Buffer.from(hash, 'hex')
    return secp256k1.ecdsaVerify(sig, buf, this.verifyKey)
  }
}
