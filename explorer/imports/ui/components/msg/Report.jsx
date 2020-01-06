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
`

const TBody = styled(Flex)`
  padding: 0 18px;
  flex-direction: column;
`

const Title = styled(Flex)`
  font-size: 0.8em;
  color: black;
  width: 200px;
`

const Tab = styled(Flex)`
  justify-content: center;
  align-items: center;
  padding: 10px 15px;
  cursor: pointer;

  ${props =>
    props.active &&
    `
    border-bottom: 1px solid black;
  `}
`

export default ({ msg }) => {
  const { reqID, status, from, data } = msg
  return (
    <Table mb="20px">
      <Flex
        width="100%"
        flexDirection="row"
        style={{ borderBottom: '1px solid #ececec' }}
      >
        <Title style={{ marginRight: '18px', padding: '5px' }}>
          <Flex
            bg="#ececec"
            justifyContent="flex-start"
            pl="21px"
            alignItems="center"
            width="100%"
          >
            Data Report
          </Flex>
        </Title>
        <Tab active>Overview</Tab>
      </Flex>
      <TBody>
        <Flex flexDirection="row" mx="10px" my="13px">
          <Title>Request ID:</Title>
          <Text fontSize="0.8em">{reqID}</Text>
        </Flex>
        <SeperatorLine color="#dadada" height="1px" my="0px" />
        <Flex flexDirection="row" mx="10px" my="13px">
          <Title>Status:</Title>
          <Text fontSize="0.8em">{status}</Text>
        </Flex>
        <SeperatorLine color="#dadada" height="1px" my="0px" />
        <Flex flexDirection="row" mx="10px" my="13px">
          <Title>From:</Title>
          <Text fontSize="0.8em">{from}</Text>
        </Flex>
        <SeperatorLine color="#dadada" height="1px" my="0px" />
        <Flex flexDirection="row" mx="10px" my="13px">
          <Title>Data:</Title>
          <Text fontSize="0.8em">{data}</Text>
        </Flex>
      </TBody>
    </Table>
  )
}
