import React from 'react'
import Table from './Table'
import { Flex, Text } from 'rebass'
import styled from 'styled-components'
import EnchancedList from './EnchancedList'
import NoData from './NoData'
import Jazzicon from './Jazzicon'

/* Latest Transactions specification */
/**
 * [
 *  {
 *    txHash: string,
 *    ts: moment,
 *    action: string,
 *    pair: string,
 *    fee: int,
 *  }
 * ]
 */

const mockData = [
  {
    txHash:
      '0xa1799f83c2f98ddf4f73f5176bb0ce54e66fa224330cfc407c9ea35fec9da67e',
    ts: 22,
    action: 'Request data',
    pair: 'ETH/USD Price',
    fee: 0.85,
  },
  {
    txHash:
      '0xa1799f83c2f98ddf4f73f5176bb0ce54e66fa224330cfc407c9ea35fec9da67e',
    ts: 22,
    action: 'Create new script',
    pair: 'ETH/USD Price',
    fee: 0.85,
  },
  {
    txHash:
      '0xa1799f83c2f98ddf4f73f5176bb0ce54e66fa224330cfc407c9ea35fec9da67e',
    ts: 22,
    action: 'Request data',
    pair: 'ETH/USD Price',
    fee: 0.85,
  },
  {
    txHash:
      '0xa1799f83c2f98ddf4f73f5176bb0ce54e66fa224330cfc407c9ea35fec9da67e',
    ts: 22,
    action: 'Request data',
    pair: 'ETH/USD Price',
    fee: 0.85,
  },
  {
    txHash:
      '0xa1799f83c2f98ddf4f73f5176bb0ce54e66fa224330cfc407c9ea35fec9da67e',
    ts: 22,
    action: 'Create new script',
    pair: 'ETH/USD Price',
    fee: 0.85,
  },
]

const BodyContainer = styled(Flex).attrs(() => ({
  px: '25px',
  width: '100%',
  flexDirection: 'column',
}))`
  '&:nth-last-child': {
    border-bottom: 0;
  },
`

const Body = ({ data }) => {
  if (!data)
    return <NoData msg="Somethings error when getting latest transations" />

  return (
    <BodyContainer>
      {data.map(({ txHash, ts, action, pair, fee }, i) => (
        <EnchancedList
          key={i}
          logo={<Jazzicon />}
          leftHeader={
            <Text
              fontSize="0.8em"
              color="black"
              width="140px"
              style={{
                whiteSpace: 'nowrap',
                textOverflow: 'elipsis',
                overflow: 'hidden',
              }}
            >
              {txHash}
            </Text>
          }
          leftMsg={
            <Text fontSize="0.5em" color="#747474">
              {ts}
            </Text>
          }
          rightHeader={
            <Flex flexDirection="row">
              <Text fontSize="0.8em" color="black" mr="5px">
                {action}
              </Text>
              <Flex
                justifyContent="center"
                alignItems="center"
                bg="#e1e1e1"
                fontSize="0.6em"
                p="3px 5px"
                style={{ borderRadius: '10px' }}
              >
                {pair}
              </Flex>
            </Flex>
          }
          rightMsg={
            <Text fontSize="0.5em" color="#747474">
              Fee: {fee} BAND
            </Text>
          }
        />
      ))}
    </BodyContainer>
  )
}

export default () => (
  <Table header="Latest Transactions" body={<Body data={mockData} />} />
)
