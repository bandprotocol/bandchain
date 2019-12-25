import React from 'react'
import { Flex, Text } from 'rebass'
import styled from 'styled-components'
import SeperatorLine from '/imports/ui/components/SeperatorLine'

const Table = styled(Flex)`
  border-radius: 3px;
  border: solid 1px #dadada;
  background-color: #ffffff;
  flex-direction: column;
  width: 100%;
  padding: 0 18px;
`

const Title = styled(Text)`
  font-size: 0.8em;
  color: black;
  width: 200px;
`

export default ({ txHash }) => (
  <>
    <Flex flexDirection="row" alignItems="center">
      <Text fontSize="1.2em" color="#5b5b5b">
        Transaction Details
      </Text>
      <Text fontSize="0.8em" ml="23px">
        {txHash}
      </Text>
    </Flex>
    <SeperatorLine color="#dadada" height="1px" my="15px" />
    <Table>
      <Flex flexDirection="row" mx="10px" my="13px">
        <Title>Transaction Hash:</Title>
        <Text fontSize="0.8em">{txHash}</Text>
      </Flex>
      <SeperatorLine color="#dadada" height="1px" my="0px" />
      <Flex flexDirection="row" mx="10px" my="13px">
        <Title>Status:</Title>
        <Text fontSize="0.8em">Success</Text>
      </Flex>
      <SeperatorLine color="#dadada" height="1px" my="0px" />
      <Flex flexDirection="row" mx="10px" my="13px">
        <Title>Height:</Title>
        <Text fontSize="0.8em">1,325</Text>
      </Flex>
      <SeperatorLine color="#dadada" height="1px" my="0px" />
      <Flex flexDirection="row" mx="10px" my="13px">
        <Title>Timestamp:</Title>
        <Text fontSize="0.8em">51 secs ago (Dec-17-2019 07:11:54 AM +UTC)</Text>
      </Flex>
    </Table>
  </>
)
