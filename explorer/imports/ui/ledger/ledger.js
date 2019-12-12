// https://github.com/zondax/cosmos-delegation-js/
// https://github.com/cosmos/ledger-cosmos-js/blob/master/src/index.js
import 'babel-polyfill';
import Cosmos from "@lunie/cosmos-js"
import TransportWebUSB from "@ledgerhq/hw-transport-webusb";
import CosmosApp from "ledger-cosmos-js"
import { signatureImport } from "secp256k1"
import semver from "semver"
import bech32 from "bech32";
import secp256k1 from "secp256k1";
import sha256 from "crypto-js/sha256"
import ripemd160 from "crypto-js/ripemd160"
import CryptoJS from "crypto-js"

// TODO: discuss TIMEOUT value
const INTERACTION_TIMEOUT = 10000
const REQUIRED_COSMOS_APP_VERSION = "1.5.0"
const DEFAULT_DENOM = 'uatom';
const DEFAULT_GAS = 200000;
export const DEFAULT_GAS_PRICE = 0.025;
export const DEFAULT_MEMO = 'Sent via Big Dipper'

/*
HD wallet derivation path (BIP44)
DerivationPath{44, 118, account, 0, index}
*/

const HDPATH = [44, 118, 0, 0, 0]
const BECH32PREFIX = Meteor.settings.public.bech32PrefixAccAddr

function bech32ify(address, prefix) {
    const words = bech32.toWords(address)
    return bech32.encode(prefix, words)
}

export const toPubKey = (address) => {
    return bech32.decode(Meteor.settings.public.bech32PrefixAccAddr, address);
}

function createCosmosAddress(publicKey) {
    const message = CryptoJS.enc.Hex.parse(publicKey.toString(`hex`))
    const hash = ripemd160(sha256(message)).toString()
    const address = Buffer.from(hash, `hex`)
    const cosmosAddress = bech32ify(address, Meteor.settings.public.bech32PrefixAccAddr)
    return cosmosAddress
}

export class Ledger {
    constructor({ testModeAllowed }) {
        this.testModeAllowed = testModeAllowed
    }

    // test connection and compatibility
    async testDevice() {
        // poll device with low timeout to check if the device is connected
        const secondsTimeout = 3 // a lower value always timeouts
        await this.connect(secondsTimeout)
    }
    async isSendingData() {
        // check if the device is connected or on screensaver mode
        const response = await this.cosmosApp.publicKey(HDPATH)
        this.checkLedgerErrors(response, {
            timeoutMessag: "Could not find a connected and unlocked Ledger device"
        })
    }
    async isReady() {
    // check if the version is supported
        const version = await this.getCosmosAppVersion()

        if (!semver.gte(version, REQUIRED_COSMOS_APP_VERSION)) {
        const msg = `Outdated version: Please update Ledger Cosmos App to the latest version.`
        throw new Error(msg)
        }

        // throws if not open
        await this.isCosmosAppOpen()
    }
    // connects to the device and checks for compatibility
    async connect(timeout = INTERACTION_TIMEOUT) {
        // assume well connection if connected once
        if (this.cosmosApp) return

        const transport = await TransportWebUSB.create(timeout)
        const cosmosLedgerApp = new CosmosApp(transport)

        this.cosmosApp = cosmosLedgerApp

        await this.isSendingData()
        await this.isReady()
    }
    async getCosmosAppVersion() {
        await this.connect()

        const response = await this.cosmosApp.getVersion()
        this.checkLedgerErrors(response)
        const { major, minor, patch, test_mode } = response
        checkAppMode(this.testModeAllowed, test_mode)
        const version = versionString({ major, minor, patch })

        return version
    }
    async isCosmosAppOpen() {
        await this.connect()

        const response = await this.cosmosApp.appInfo()
        this.checkLedgerErrors(response)
        const { appName } = response

        if (appName.toLowerCase() !== `cosmos`) {
            throw new Error(`Close ${appName} and open the Cosmos app`)
        }
      }
    async getPubKey() {
        await this.connect()

        const response = await this.cosmosApp.publicKey(HDPATH)
        this.checkLedgerErrors(response)
        return response.compressed_pk
      }
    async getCosmosAddress() {
        await this.connect()

        const pubKey = await this.getPubKey(this.cosmosApp)
        return {pubKey, address:createCosmosAddress(pubKey)}
      }
    async confirmLedgerAddress() {
        await this.connect()
        const cosmosAppVersion = await this.getCosmosAppVersion()

        if (semver.lt(cosmosAppVersion, REQUIRED_COSMOS_APP_VERSION)) {
            // we can't check the address on an old cosmos app
            return
        }

        const response = await this.cosmosApp.getAddressAndPubKey(
            HDPATH,
            BECH32PREFIX,
        )
        this.checkLedgerErrors(response, {
            rejectionMessage: "Displayed address was rejected"
        })
      }

    async sign(signMessage) {
        await this.connect()

        const response = await this.cosmosApp.sign(HDPATH, signMessage)
        this.checkLedgerErrors(response)
        // we have to parse the signature from Ledger as it's in DER format
        const parsedSignature = signatureImport(response.signature)
        return parsedSignature
      }

  /* istanbul ignore next: maps a bunch of errors */
  checkLedgerErrors(
    { error_message, device_locked },
    {
      timeoutMessag = "Connection timed out. Please try again.",
      rejectionMessage = "User rejected the transaction"
    } = {}
  ) {
    if (device_locked) {
      throw new Error(`Ledger's screensaver mode is on`)
    }
    switch (error_message) {
        case `U2F: Timeout`:
            throw new Error(timeoutMessag)
        case `Cosmos app does not seem to be open`:
            // hack:
            // It seems that when switching app in Ledger, WebUSB will disconnect, disabling further action.
            // So we clean up here, and re-initialize this.cosmosApp next time when calling `connect`
            this.cosmosApp.transport.close()
            this.cosmosApp = undefined
            throw new Error(`Cosmos app is not open`)
        case `Command not allowed`:
            throw new Error(`Transaction rejected`)
        case `Transaction rejected`:
            throw new Error(rejectionMessage)
        case `Unknown error code`:
            throw new Error(`Ledger's screensaver mode is on`)
        case `Instruction not supported`:
            throw new Error(
                `Your Cosmos Ledger App is not up to date. ` +
                `Please update to version ${REQUIRED_COSMOS_APP_VERSION}.`
            )
        case `No errors`:
            // do nothing
            break
        default:
            throw new Error(error_message)
        }
    }

    static getBytesToSign(tx, txContext) {
        if (typeof txContext === 'undefined') {
            throw new Error('txContext is not defined');
        }
        if (typeof txContext.chainId === 'undefined') {
            throw new Error('txContext does not contain the chainId');
        }
        if (typeof txContext.accountNumber === 'undefined') {
            throw new Error('txContext does not contain the accountNumber');
        }
        if (typeof txContext.sequence === 'undefined') {
            throw new Error('txContext does not contain the sequence value');
        }

        const txFieldsToSign = {
            account_number: txContext.accountNumber.toString(),
            chain_id: txContext.chainId,
            fee: tx.value.fee,
            memo: tx.value.memo,
            msgs: tx.value.msg,
            sequence: txContext.sequence.toString(),
        };

        return JSON.stringify(canonicalizeJson(txFieldsToSign));
    }

    static applyGas(unsignedTx, gas, gasPrice=DEFAULT_GAS_PRICE, denom=DEFAULT_DENOM) {
        if (typeof unsignedTx === 'undefined') {
            throw new Error('undefined unsignedTx');
        }
        if (typeof gas === 'undefined') {
            throw new Error('undefined gas');
        }

        // eslint-disable-next-line no-param-reassign
        unsignedTx.value.fee = {
            amount: [{
                 amount: Math.round(gas * gasPrice).toString(),
                 denom: denom,
            }],
            gas: gas.toString(),
        };

        return unsignedTx;
    }

    static applySignature(unsignedTx, txContext, secp256k1Sig) {
        if (typeof unsignedTx === 'undefined') {
            throw new Error('undefined unsignedTx');
        }
        if (typeof txContext === 'undefined') {
            throw new Error('undefined txContext');
        }
        if (typeof txContext.pk === 'undefined') {
            throw new Error('txContext does not contain the public key (pk)');
        }
        if (typeof txContext.accountNumber === 'undefined') {
            throw new Error('txContext does not contain the accountNumber');
        }
        if (typeof txContext.sequence === 'undefined') {
            throw new Error('txContext does not contain the sequence value');
        }

        const tmpCopy = Object.assign({}, unsignedTx, {});

        tmpCopy.value.signatures = [
            {
                signature: secp256k1Sig.toString('base64'),
                account_number: txContext.accountNumber.toString(),
                sequence: txContext.sequence.toString(),
                pub_key: {
                    type: 'tendermint/PubKeySecp256k1',
                    value: txContext.pk//Buffer.from(txContext.pk, 'hex').toString('base64'),
                },
            },
        ];
        return tmpCopy;
    }

    // Creates a new tx skeleton
    static createSkeleton(txContext, msgs=[]) {
        if (typeof txContext === 'undefined') {
            throw new Error('undefined txContext');
        }
        if (typeof txContext.accountNumber === 'undefined') {
            throw new Error('txContext does not contain the accountNumber');
        }
        if (typeof txContext.sequence === 'undefined') {
            throw new Error('txContext does not contain the sequence value');
        }
        const txSkeleton = {
            type: 'auth/StdTx',
            value: {
                msg: msgs,
                fee: '',
                memo: txContext.memo || DEFAULT_MEMO,
                signatures: [{
                    signature: 'N/A',
                    account_number: txContext.accountNumber.toString(),
                    sequence: txContext.sequence.toString(),
                    pub_key: {
                        type: 'tendermint/PubKeySecp256k1',
                        value: txContext.pk || 'PK',
                    },
                }],
            },
        };
        //return Ledger.applyGas(txSkeleton, DEFAULT_GAS);
        return txSkeleton
    }

    // Creates a new delegation tx based on the input parameters
    // the function expects a complete txContext
    static createDelegate(
        txContext,
        validatorBech32,
        uatomAmount
    ) {
        const txMsg = {
            type: 'cosmos-sdk/MsgDelegate',
            value: {
                amount: {
                    amount: uatomAmount.toString(),
                    denom: txContext.denom,
                },
                delegator_address: txContext.bech32,
                validator_address: validatorBech32,
            },
        };

        return Ledger.createSkeleton(txContext, [txMsg]);
    }

    // Creates a new undelegation tx based on the input parameters
    // the function expects a complete txContext
    static createUndelegate(
        txContext,
        validatorBech32,
        uatomAmount
    ) {
        const txMsg = {
            type: 'cosmos-sdk/MsgUndelegate',
            value: {
                amount: {
                    amount: uatomAmount.toString(),
                    denom: txContext.denom,
                },
                delegator_address: txContext.bech32,
                validator_address: validatorBech32,
            },
        };

        return Ledger.createSkeleton(txContext, [txMsg]);
    }

    // Creates a new redelegation tx based on the input parameters
    // the function expects a complete txContext
    static createRedelegate(
        txContext,
        validatorSourceBech32,
        validatorDestBech32,
        uatomAmount
    ) {
        const txMsg = {
            type: 'cosmos-sdk/MsgBeginRedelegate',
            value: {
                amount: {
                    amount: uatomAmount.toString(),
                    denom: txContext.denom,
                },
                delegator_address: txContext.bech32,
                validator_dst_address: validatorDestBech32,
                validator_src_address: validatorSourceBech32,
            },
        };

        return Ledger.createSkeleton(txContext, [txMsg]);
    }

    // Creates a new transfer tx based on the input parameters
    // the function expects a complete txContext
    static createTransfer(
        txContext,
        toAddress,
        amount
    ) {
        const txMsg = {
            type: 'cosmos-sdk/MsgSend',
            value: {
                amount: [{
                    amount: amount.toString(),
                    denom: txContext.denom
                }],
                from_address: txContext.bech32,
                to_address: toAddress
            }
        };

        return Ledger.createSkeleton(txContext, [txMsg]);
    }

    static createSubmitProposal(
        txContext,
        title,
        description,
        deposit
    ) {
        const txMsg = {
            type: 'cosmos-sdk/MsgSubmitProposal',
            value: {
                content: {
                    type: "cosmos-sdk/TextProposal",
                    value: {
                        description: description,
                        title: title
                    }
                },
                initial_deposit: [{
                    amount: deposit.toString(),
                    denom: txContext.denom
                }],
                proposer: txContext.bech32
            }
        };

        return Ledger.createSkeleton(txContext, [txMsg]);
    }

    static createVote(
        txContext,
        proposalId,
        option,
    ) {
        const txMsg = {
            type: 'cosmos-sdk/MsgVote',
            value: {
                option,
                proposal_id: proposalId.toString(),
                voter: txContext.bech32
            }
        };

        return Ledger.createSkeleton(txContext, [txMsg]);
    }

    static createDeposit(
        txContext,
        proposalId,
        amount,
    ) {
        const txMsg = {
            type: 'cosmos-sdk/MsgDeposit',
            value: {
                amount: [{
                    amount: amount.toString(),
                    denom: txContext.denom
                }],
                depositor: txContext.bech32,
                proposal_id: proposalId.toString()
            }
        };

        return Ledger.createSkeleton(txContext, [txMsg]);
    }

}

function versionString({ major, minor, patch }) {
    return `${major}.${minor}.${patch}`
}

export const checkAppMode = (testModeAllowed, testMode) => {
    if (testMode && !testModeAllowed) {
        throw new Error(
            `DANGER: The Cosmos Ledger app is in test mode and shouldn't be used on mainnet!`
        )
    }
}

function canonicalizeJson(jsonTx) {
    if (Array.isArray(jsonTx)) {
        return jsonTx.map(canonicalizeJson);
    }
    if (typeof jsonTx !== 'object') {
        return jsonTx;
    }
    const tmp = {};
    Object.keys(jsonTx).sort().forEach((key) => {
        // eslint-disable-next-line no-unused-expressions
        jsonTx[key] != null && (tmp[key] = jsonTx[key]);
    });

    return tmp;
}
