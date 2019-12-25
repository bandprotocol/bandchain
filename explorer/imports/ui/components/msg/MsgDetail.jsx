import React from 'react'
import { Flex, Text } from 'rebass'
import SeperatorLine from '/imports/ui/components/SeperatorLine'
import MsgTable from './MsgTable'

const mockData = [
  {
    action: 'CREATE_SCRIPT',
    scriptHash: '0x9q391923912392139231923912323',
    creator: '0x01230013012023021301230',
    name: 'ETH/USD Median Price',
    sources: ['CoinMarketcap', 'Crypto Compare', 'Binance'],
    tag: ['Crypto Price'],
    code: 'You will see the code right here in the future.',
  },
  {
    action: 'REPORT',
    reqID: '32',
    status: 'Succes',
    from: '0x1348421874384384383484383483',
    data: ['0x0000008332', '0x0000008332'],
  },
  {
    action: 'REQUEST',
    reqID: '32',
    status: 'pending',
    ts: '1,329',
    from: '0x1348421874384384383484383483',
    script: '0x1348421874384384383484383483',
    result: '0x0000008332',
    proof: 'asdkaksdkasdkaskdk',
    reports: [
      {
        txHash: '0x9123912923192139213912391239329',
        block: '1,324',
        ts: '22',
        from: '0x85334473743743743734734734734',
        value: ['0x0000008332', '0x0000008332', '0x0000008332'],
      },
      {
        txHash: '0x9123912923192139213912391239329',
        block: '1,324',
        ts: '22',
        from: '0x85334473743743743734734734734',
        value: ['0x0000008332', '0x0000008332', '0x0000008332'],
      },
    ],
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
