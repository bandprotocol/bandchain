import React from 'react'
import { Flex, Text } from 'rebass'
import SeperatorLine from '/imports/ui/components/SeperatorLine'
import MsgTable from './MsgTable'

const mockData = [
  {
    scriptHash: '0x9q391923912392139231923912323',
    creator: '0x01230013012023021301230',
    name: 'ETH/USD Median Price',
    sources: ['CoinMarketcap', 'Crypto Compare', 'Binance'],
    tag: ['Crypto Price'],
  },
  {
    scriptHash: '0x834473473473467347347347347',
    creator: '0x1348421874384384383484383483',
    name: 'BTC/USD Median Price',
    sources: ['CoinMarketcap', 'Crypto Compare'],
    tag: ['Crypto Price'],
  },
  {
    scriptHash: '0x834473473473467347347347347',
    creator: '0x1348421874384384383484383483',
    name: 'BTC/USD Median Price',
    sources: ['CoinMarketcap', 'Crypto Compare'],
    tag: ['Crypto Price'],
  },
]

export default () => {
  return (
    <>
      <Flex flexDirection="row" alignItems="center" mt="20px">
        <Text fontSize="1.2em" color="#5b5b5b">
          Messages ({mockData.length})
        </Text>
      </Flex>
      <SeperatorLine color="#dadada" height="1px" my="15px" />
      {mockData.map((each, i) => (
        <MsgTable msg={each} key={i} />
      ))}
    </>
  )
}
