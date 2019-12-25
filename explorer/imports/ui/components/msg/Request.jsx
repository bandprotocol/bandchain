import React, { useState } from 'react'
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
  const { reqID, status, ts, from, script, result, proof, reports } = msg
  const [tabId, setTabId] = useState(0)
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
            Data Request #{reqID}
          </Flex>
        </Title>
        <Tab active={tabId == 0} onClick={() => setTabId(0)}>
          Overview
        </Tab>
        <Tab active={tabId == 1} onClick={() => setTabId(1)}>
          Data Report Status
        </Tab>
      </Flex>
      <TBody>
        {tabId == 0 ? (
          <>
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
              <Title>Timestamp:</Title>
              <Text fontSize="0.8em">{ts}</Text>
            </Flex>
            <SeperatorLine color="#dadada" height="1px" my="0px" />
            <Flex flexDirection="row" mx="10px" my="13px">
              <Title>From:</Title>
              <Text fontSize="0.8em">{from}</Text>
            </Flex>
            <SeperatorLine color="#dadada" height="1px" my="0px" />
            <Flex flexDirection="row" mx="10px" my="13px">
              <Title>Request Script:</Title>
              <Text fontSize="0.8em">{script}</Text>
            </Flex>
            <SeperatorLine color="#dadada" height="1px" my="0px" />
            <Flex flexDirection="row" mx="10px" my="13px">
              <Title>Result:</Title>
              <Text fontSize="0.8em">{result}</Text>
            </Flex>
            <SeperatorLine color="#dadada" height="1px" my="0px" />
            <Flex flexDirection="row" mx="10px" my="13px">
              <Title>Proof of validity:</Title>
              <Text fontSize="0.8em">{proof}</Text>
            </Flex>
          </>
        ) : (
          <>
            <Flex flexDirection="row" mx="10px" my="13px">
              <Text fontSize="0.8em">
                Data Reports from 4 validators (Completed 3/4)
              </Text>
            </Flex>
            <Flex>Table</Flex>
          </>
        )}
      </TBody>
    </Table>
  )
}
