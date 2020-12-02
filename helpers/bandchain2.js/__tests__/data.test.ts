import { Coin } from '../src/data'

describe('Coin', () => {
  it('success', () => {
    let coin = new Coin(1000, 'uband')

    expect(coin.validate()).toBeTruthy()

    const exepected = {
      amount: '1000',
      denom: 'uband',
    }

    expect(coin.asJson()).toEqual(exepected)
  })

  it('error integer', () => {
    let coin = new Coin(1000.2, 'uband')
    expect(() => {
      coin.validate()
    }).toThrowError('amount is not an integer')
  })

  it('error amount less than 0', () => {
    let coin = new Coin(-50, 'uband')
    expect(() => {
      coin.validate()
    }).toThrowError('Expect amount more than 0')
  })

  it('error expect denom', () => {
    let coin = new Coin(1000, '')
    expect(() => {
      coin.validate()
    }).toThrowError('Expect denom')
  })
})
