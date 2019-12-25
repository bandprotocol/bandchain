import React from 'react'
import { Flex, Text, Box } from 'rebass'
import styled from 'styled-components'
import YSeperatorLine from './SeperatorLine'

const Card = styled(Flex)`
  padding: 25px;
  border-radius: 3px;
  border: solid 1px #dadada;
  background-color: #ffffff;
  width: 100%;
  margin-top: 19px;
  justtify-content: space-between;
`

const GreyText = styled(Text)`
  font-size: 0.8em;
  color: #747474;
  white-space: nowrap;
`

const NormalText = styled(Text)`
  font-size: 0.8em;
  color: black;
`

const Badge = styled.div`
  min-height: 40px;
  max-height: 40px;
  min-width: 40px;
  max-width: 40px;
  border-radius: 50%;
  background-color: #d8d8d8;
  border: 1px solid #979797;
`

const XSeperatorLine = () => (
  <Box height="auto" bg="#d8d8d8" width="1px" mx="22px" />
)

export default () => {
  return (
    <Card>
      <Flex flexDirection="column" flex={1}>
        <Flex flexDirection="row" alignItems="center">
          <Badge />
          <Flex width="100%">
            <Flex flexDirection="column" ml="10px" flex={1}>
              <GreyText>BAND Price</GreyText>
              <Flex flexDirection="row">
                <NormalText>$0.642</NormalText>
                <GreyText>@0.012 BTC</GreyText>
                <Text fontSize="14px" color="#3ca88b" ml="5px">
                  (+1.2%)
                </Text>
              </Flex>
            </Flex>
          </Flex>
        </Flex>
        <YSeperatorLine color="#dadada" height="1px" my="10px" />
        <Flex flexDirection="row" alignItems="center">
          <Badge />
          <Flex width="100%">
            <Flex flexDirection="column" ml="10px" flex={1}>
              <GreyText>Marketcap</GreyText>
              <NormalText>$8,428,380.55</NormalText>
            </Flex>
          </Flex>
        </Flex>
      </Flex>
      <XSeperatorLine />
      <Flex flexDirection="column" flex={1}>
        <Flex flexDirection="row" alignItems="center">
          <Badge />
          <Flex width="100%">
            <Flex flexDirection="column" ml="10px" flex={1}>
              <GreyText>Latest Block</GreyText>
              <NormalText>13,325</NormalText>
            </Flex>
            <Flex
              flexDirection="column"
              ml="10px"
              alignItems="flex-end"
              flex={1}
            >
              <GreyText>Transactions</GreyText>
              <NormalText>20</NormalText>
            </Flex>
          </Flex>
        </Flex>
        <YSeperatorLine color="#dadada" height="1px" my="10px" />
        <Flex flexDirection="row" alignItems="center">
          <Badge />
          <Flex width="100%">
            <Flex flexDirection="column" ml="10px" flex={1}>
              <GreyText>Active Validators</GreyText>
              <NormalText>4 Nodes</NormalText>
            </Flex>
            <Flex
              flexDirection="column"
              ml="10px"
              alignItems="flex-end"
              flex={1}
            >
              <GreyText>Online Voting Power</GreyText>
              <NormalText>431,324.98 BAND</NormalText>
            </Flex>
          </Flex>
        </Flex>
      </Flex>
      <XSeperatorLine />
      <Flex flexDirection="column" flex={1}>
        <GreyText>Data Request Transaction History in 14 days</GreyText>
        <Flex justifyContent="center" alignItems="center" height="100%">
          <GreyText>Graph Visualization</GreyText>
        </Flex>
      </Flex>
    </Card>
  )
}
