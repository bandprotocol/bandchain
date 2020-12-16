import { isBip44, bip44ToArray } from '../src/helpers'

describe('Bip44', () => {
  it('isBip44', () => {
    expect(isBip44(`m/44'/494'/0'/0/0`)).toBeTruthy()
    expect(isBip44(`m/44'/118'/0'/0/0`)).toBeTruthy()
    expect(isBip44(`m/44'/118'/89'/000/00`)).toBeTruthy()

    expect(!isBip44(`m/'/'/0'/0/0`)).toBeTruthy()
    expect(!isBip44(`m/44'/494'/0'/0/0'`)).toBeTruthy()
    expect(!isBip44(`m/44'/494'/0'/0/0/`)).toBeTruthy()
    expect(!isBip44(`m/44'/494'/0'/0/_0'`)).toBeTruthy()
    expect(!isBip44(`m/44'/494'/0'/0/'0`)).toBeTruthy()
    expect(!isBip44(`m/44'/494'/0'/0`)).toBeTruthy()
  })
  
  it('bip44ToArray', () => {
    expect(bip44ToArray(`m/44'/494'/0'/0/0`)).toEqual([44, 494, 0, 0, 0])
    expect(bip44ToArray(`m/44'/118'/0'/0/0`)).toEqual([44, 118, 0, 0, 0])
    expect(bip44ToArray(`m/44'/118'/0'/88/99`)).toEqual([44, 118, 0, 88, 99])
    expect(bip44ToArray(`m/0'/0'/0'/88/99`)).toEqual([0, 0, 0, 88, 99])
  })
})