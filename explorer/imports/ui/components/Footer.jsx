import React from 'react'
import { Image, Text, Flex } from 'rebass'
import PageContainer from './PageContainer'

export default () => (
  <Flex bg="#444444" justifyContent="center" style={{ height: '260px' }}>
    <PageContainer>
      <Flex
        width="100%"
        justifyContent="center"
        color="white"
        fontSize="25px"
        alignItems="center"
      >
        THIS IS FOOTER
      </Flex>
    </PageContainer>
  </Flex>
)
