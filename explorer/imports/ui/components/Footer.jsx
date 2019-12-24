import React from 'react'
import { Image, Text, Flex, Box } from 'rebass'
import PageContainer from './PageContainer'

export default () => (
  <Flex bg="#444444" justifyContent="center" style={{ height: '260px' }}>
    <PageContainer>
      <Flex
        width="100%"
        color="white"
        pt="34px"
        pb="16px"
        flexDirection="column"
      >
        {/* Upper */}
        <Flex pb="27px" flex={2} justifyContent="space-evenly">
          {/* First Column */}
          <Flex flexDirection="column" mr="70px">
            <Flex flexDirection="row">
              <Box>
                <Image src="/img/band-logo.png" width="47px" />
              </Box>
              <Flex flexDirection="column" ml="12px">
                <Text fontSize="16px">POWERED BY</Text>
                <Text fontSize="25px">Band Protocol</Text>
              </Flex>
            </Flex>
            <Flex maxWidth="270px" mt="20px" fontSize="14px">
              D3NScan.com is a block explorer and data analytic platform for
              Band Protocol, a decentralized platform for reliable and provably
              secure data on the blockchains
            </Flex>
          </Flex>

          {/* Second Column */}
          <Flex flexDirection="column" flex={1}>
            <Text>D3N Project</Text>
            <Box width="100%" bg="white" height="1px" my="11px" />
            <Text mb="10px">Band Protocol Website</Text>
            <Text mb="10px">D3N Blockchain</Text>
            <Text mb="10px">D3N Wallet</Text>
          </Flex>

          {/* Third Column */}
          <Flex flexDirection="column" flex={1} mx="20px">
            <Text>Community</Text>
            <Box width="100%" bg="white" height="1px" my="11px" />
            <Text mb="10px">Open Source Repositories</Text>
            <Text mb="10px">Developer Docs</Text>
            <Text mb="10px">Network Status</Text>
          </Flex>

          {/* Fourth Column */}
          <Flex flexDirection="column" flex={1}>
            <Text>Social Links</Text>
            <Box width="100%" bg="white" height="1px" my="11px" />
            <Text mb="10px">Twitter</Text>
            <Text mb="10px">Telegram</Text>
            <Text mb="10px">Medium</Text>
          </Flex>
        </Flex>
        <Box width="100%" bg="white" height="1.5px" />
        {/* Lower */}
        <Flex mt="16px" width="100%" justifyContent="space-between">
          <Text>World Data Foundation © 2019-2020 (C)</Text>
          <Text>connect@bandprotocol.com</Text>
        </Flex>
      </Flex>
    </PageContainer>
  </Flex>
)
