import * as bip39 from 'bip39'
import * as bip32 from 'bip32'
import * as bech32 from 'bech32'
import secp256k1, { signatureImport } from 'secp256k1'
import crypto from 'crypto'
import { ECPair } from 'bitcoinjs-lib'
import TransportNodeHid from '@ledgerhq/hw-transport-node-hid'
// import TransportWebHid from '@ledgerhq/hw-transport-webhid'
import TransportWebUSB from '@ledgerhq/hw-transport-webusb'
import { CosmosApp, LedgerResponse, AppInfo } from 'ledger-cosmos-js'
import { Transaction } from 'index'
import {isBip44, bip44ToArray} from './helpers'
// import CosmosApp from 'ledger-cosmos-js'


const BECH32_PUBKEY_ACC_PREFIX = 'bandpub'
const BECH32_PUBKEY_VAL_PREFIX = 'bandvaloperpub'
const BECH32_PUBKEY_CONS_PREFIX = 'bandvalconspub'

const BECH32_ADDR_ACC_PREFIX = 'band'
const BECH32_ADDR_VAL_PREFIX = 'bandvaloper'
const BECH32_ADDR_CONS_PREFIX = 'bandvalcons'

const DEFAULT_DERIVATION_PATH = "m/44'/494'/0'/0'/0"
const DEFAULT_DERIVATION_PATH_LEDGER = "m/44'/118'/0'/0'/0"

export class Ledger {
  private cosmosApp?: CosmosApp
  public hdPath: string = DEFAULT_DERIVATION_PATH_LEDGER

  private constructor() {}

  static async connectLedgerWeb(hdPath: string): Promise<Ledger> {
    if (!isBip44(hdPath)) throw Error('Not BIP 44')

    let ledger = new Ledger()
    ledger.hdPath = hdPath
    await ledger.connectLedgerWeb()

    return ledger
  }

  static async connectLedgerNode(hdPath: string): Promise<Ledger> {
    if (!isBip44(hdPath)) throw Error('Not BIP 44')

    let ledger = new Ledger()
    ledger.hdPath = hdPath
    await ledger.connectLedgerNode()

    return ledger
  }

  checkLedgerError(response: LedgerResponse): void {
    switch(response.error_message) {
      case `No errors`:
        return
      default:
        throw new Error(response.error_message)
    }
  }

  private async connectLedgerNode(): Promise<void> {
    if (this.cosmosApp) return

    const transport = await TransportNodeHid.create(3, 3)
    const ledgerCosmosApp = new CosmosApp(transport)

    this.cosmosApp = ledgerCosmosApp

    await this.isCosmosAppOpen()
  }

  private async connectLedgerWeb(): Promise<void> {
    if (this.cosmosApp) return

    const transport = await TransportWebUSB.create(3, 3)
    const ledgerCosmosApp = new CosmosApp(transport)

    this.cosmosApp = ledgerCosmosApp

    await this.isCosmosAppOpen()
  }

  async isCosmosAppOpen(): Promise<boolean> {
    // await this.connect()

    const response = await this.cosmosApp!.appInfo()
    this.checkLedgerError(response)

    const { appName } = response
    if (appName.toLowerCase() !== 'cosmos') throw new Error(`Please close ${appName} and open the Cosmos app.`)

    return true
  }

  // async getCosmosAppVersion(): Promise<string> {
  //   await this.connect

  //   const response = await this.cosmosApp!.getVersion()
  //   this.checkLedgerError(response)

  //   const {major, minor, patch} = response
  //   const version = `${major}.${minor}.${patch}`

  //   return version
  // }

  async appInfo(): Promise<AppInfo> {
    const response = await this.cosmosApp!.appInfo()
    this.checkLedgerError(response)

    return response
  }

  // async isReady(): Promise<void> {
  //   const _ = await this.getCosmosAppVersion()

  // }

  // async get

  // async isReady(): Promise<boolean> {
  //   const _ = await this.

  //   return true
  // }

  async sign(transaction: Transaction): Promise<Buffer> {
    const signature = await this.cosmosApp!.sign(bip44ToArray(this.hdPath), transaction.getSignData())
    return Buffer.from(signatureImport(signature.signature))
  }

  async toPubKey(): Promise<PublicKey> {
    const response = await this.cosmosApp!.getAddressAndPubKey(bip44ToArray(this.hdPath), 'band')
    return PublicKey.fromHex(response.compressed_pk.toString('hex'))
  }
}

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

  toPubkey(): PublicKey {
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

  private static fromBech32(bech: string, _prefix: string): PublicKey {
    const { prefix, words } = bech32.decode(bech)
    if (prefix != _prefix) throw new Error('Invalid bech32 prefix')
    if (words.length == 0) throw new Error('Cannot decode bech32')

    return new PublicKey(Buffer.from(bech32.fromWords(words).slice(5)))
  }

  static fromHex(pub: string): PublicKey {
    return new PublicKey(Buffer.from(pub, 'hex'))
  }

  static fromAccBech32(bech: string): PublicKey {
    return this.fromBech32(bech, BECH32_PUBKEY_ACC_PREFIX)
  }

  static fromValBech32(bech: string): PublicKey {
    return this.fromBech32(bech, BECH32_PUBKEY_VAL_PREFIX)
  }

  static fromConsBech32(bech: string): PublicKey {
    return this.fromBech32(bech, BECH32_PUBKEY_CONS_PREFIX)
  }

  private toBech32(prefix: string): string {
    const hex = Buffer.concat([
      Buffer.from('eb5ae98721', 'hex'),
      this.verifyKey,
    ])
    const words = bech32.toWords(Buffer.from(hex))
    if (words.length == 0) throw new Error('Unsuccessful bech32.toWords call')

    return bech32.encode(prefix, words)
  }

  toAccBech32(): string {
    return this.toBech32(BECH32_PUBKEY_ACC_PREFIX)
  }

  toValBech32(): string {
    return this.toBech32(BECH32_PUBKEY_VAL_PREFIX)
  }

  toConsBech32(): string {
    return this.toBech32(BECH32_PUBKEY_CONS_PREFIX)
  }

  toHex(): string {
    return this.verifyKey.toString('hex')
  }

  toAddress(): Address {
    const hash = crypto.createHash('sha256').update(this.verifyKey).digest()

    return Address.fromHex(
      crypto.createHash('ripemd160').update(hash).digest('hex'),
    )
  }

  verify(msg: Buffer, sig: Buffer): boolean {
    const hash = crypto.createHash('sha256').update(msg).digest('hex')
    const buf = Buffer.from(hash, 'hex')
    return secp256k1.ecdsaVerify(sig, buf, this.verifyKey)
  }
}

export class Address {
  private addr: Buffer

  private constructor(addr: Buffer) {
    this.addr = addr
  }

  private static fromBech32(bech: string, _prefix: string): Address {
    const { prefix, words } = bech32.decode(bech)
    if (prefix != _prefix) throw new Error('Invalid bech32 prefix')
    if (words.length == 0) throw new Error('Cannot decode bech32')

    return new Address(Buffer.from(bech32.fromWords(words)))
  }

  static fromHex(hex: string): Address {
    return new Address(Buffer.from(hex, 'hex'))
  }

  static fromAccBech32(bech: string): Address {
    return this.fromBech32(bech, BECH32_ADDR_ACC_PREFIX)
  }

  static fromValBech32(bech: string): Address {
    return this.fromBech32(bech, BECH32_ADDR_VAL_PREFIX)
  }

  static fromConsBech32(bech: string): Address {
    return this.fromBech32(bech, BECH32_ADDR_CONS_PREFIX)
  }

  private toBech32(prefix: string): string {
    const words = bech32.toWords(this.addr)
    if (words.length == 0) throw new Error('Unsuccessful bech32.toWords call')

    return bech32.encode(prefix, words)
  }

  toAccBech32(): string {
    return this.toBech32(BECH32_ADDR_ACC_PREFIX)
  }

  toValBech32(): string {
    return this.toBech32(BECH32_ADDR_VAL_PREFIX)
  }

  toConsBech32(): string {
    return this.toBech32(BECH32_ADDR_CONS_PREFIX)
  }

  toHex(): string {
    return this.addr.toString('hex')
  }
}
