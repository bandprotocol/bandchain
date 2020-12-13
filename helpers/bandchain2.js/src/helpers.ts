

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