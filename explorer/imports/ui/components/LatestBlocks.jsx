import React from 'react'
import Table from './Table'
import { Flex, Text } from 'rebass'
import styled from 'styled-components'
import EnchancedList from './EnchancedList'
import NoData from './NoData'

/* Latest Blocks specification */
/**
 * [
 *  {
 *    blkNumber: int,
 *    ts: moment,
 *    proposer: string,
 *    txAmount: int
 *  }
 * ]
 */

const mockData = [
  {
    blkNumber: 23042,
    ts: 22,
    proposer: 'Validator1',
    txAmount: 30,
  },
  {
    blkNumber: 23042,
    ts: 22,
    proposer: 'Validator1',
    txAmount: 30,
  },
  {
    blkNumber: 23042,
    ts: 22,
    proposer: 'Validator1',
    txAmount: 30,
  },
  {
    blkNumber: 23042,
    ts: 22,
    proposer: 'Validator1',
    txAmount: 30,
  },
  {
    blkNumber: 23042,
    ts: 22,
    proposer: 'Validator1',
    txAmount: 30,
  },
  {
    blkNumber: 23042,
    ts: 22,
    proposer: 'Validator1',
    txAmount: 30,
  },
  {
    blkNumber: 23042,
    ts: 22,
    proposer: 'Validator1',
    txAmount: 30,
  },
  {
    blkNumber: 23042,
    ts: 22,
    proposer: 'Validator1',
    txAmount: 30,
  },
]

const Logo = styled.div`
  max-width: 40px;
  min-width: 40px;
  min-height: 40px;
  max-height: 40px;
  border: solid 1px #979797;
  background-color: #d8d8d8;
  border-radius: 3px;
`

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
  if (!data) return <NoData msg="Somethings error when getting latest blocks" />

  return (
    <BodyContainer>
      {data.map(({ blkNumber, ts, proposer, txAmount }, i) => (
        <EnchancedList
          key={i}
          logo={<Logo />}
          leftHeader={
            <Text fontSize="0.8em" color="black">
              {blkNumber}
            </Text>
          }
          leftMsg={
            <Text fontSize="0.5em" color="#747474">
              {ts}
            </Text>
          }
          rightHeader={
            <Text fontSize="0.8em" color="black">
              Proposed by {proposer}
            </Text>
          }
          rightMsg={
            <Text fontSize="0.5em" color="#747474">
              {txAmount} transactions
            </Text>
          }
        />
      ))}
    </BodyContainer>
  )
}

export default () => (
  <Table header="Latest Blocks" body={<Body data={mockData} />} />
)
