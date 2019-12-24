import React from 'react'
import PageContainer from './PageContainer'
import { Flex, Image, Text } from 'rebass'
import Search from './Search'

const Navlink = ({ text, to }) => (
  /* TODO: add NavLink */
  <Text
    fontSize="14px"
    marginLeft="0.5em"
    onClick={() => alert(to)}
    style={{ cursor: 'pointer' }}
  >
    {text}
  </Text>
)

export default () => (
  <Flex
    bg="white"
    justifyContent="center"
    style={{
      height: '120px',
      boxShadow: '0 2px 6px 0 rgba(0, 0, 0, 0.05)',
    }}
  >
    <PageContainer>
      <Flex
        flexDirection="column"
        width="100%"
        padding="1.4em 0 1.4em"
        justifyContent="space-between"
      >
        <Flex flexDirection="row" justifyContent="space-between">
          <Flex flexDirection="row" alignItems="flex-end">
            <Image src="/img/d3n.svg" marginRight="12px" />
            <Flex
              bg="#ececec"
              color="#505050"
              px="7px"
              fontSize="13px"
              height="25px"
              justifyContent="center"
              alignItems="center"
              style={{ borderRadius: '3px' }}
            >
              TESTNET v1.0
            </Flex>
          </Flex>
          <Search />
        </Flex>
        <Flex flexDirection="row" justifyContent="space-between">
          <Text fontSize="18px" color="#5b5b5b">
            DATA REQUEST EXPLORER
          </Text>
          <Flex>
            <Navlink text="Validators" to="/validators" />
            <Navlink text="Blocks" to="/blocks" />
            <Navlink text="Transactions" to="/transactions" />
            <Navlink text="Request Scripts" to="/request-scripts" />
            <Navlink text="Data Providers" to="/dataproviders" />
            <Navlink text="OWASM Studio" to="/owasm-studio" />
          </Flex>
        </Flex>
      </Flex>
    </PageContainer>
  </Flex>
)
