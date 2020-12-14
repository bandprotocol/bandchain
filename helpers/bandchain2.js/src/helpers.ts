// import * as lg from 'ledger-cosmos-js'
// import {LedgerResponse, AppInfo, AddressAndPublicKey, Version, Sign} from 'ledger-cosmos-js'

export function isBip44(path: string): boolean {
  let result = path.match(/m\/\d+'\/\d+'\/\d+'\/\d+\/\d$/)
  return !!result
}

export function bip44ToArray(path: string): number[] {
  if (!isBip44(path)) throw Error('Not BIP 44')
  let result = path.match(/\d+/g)
  if ((result?.length ?? 0) !== 5) throw Error('Not BIP 44')
  return result!.map((x:string) => parseInt(x))
}

// type Race = 'LedgerResponse' | 'AppInfo' | 'AddressAndPublicKey' | 'Version' | 'Sign'
// export type Race = LedgerResponse | AppInfo | AddressAndPublicKey | Version | Sign | undefined
export async function race<T>(promise: Promise<T | undefined>, timeout: number): Promise<T | undefined> {
  let timer: NodeJS.Timeout

  return Promise.race([
    promise.then((value?: T) => {
      clearTimeout(timer)
      return value
    }),
    new Promise((_, rj): T | undefined => {
      timer = setTimeout(rj, timeout)
      return undefined
    }) as Promise<T | undefined>
  ])
}