import React from 'react'
import { Flex } from 'rebass'
import styled from 'styled-components'

const Container = styled(Flex)`
  border-bottom: 1px solid #d7d7d7;
  &:last-child {
    border-bottom: 0;
  }
`

export default ({ logo, leftHeader, leftMsg, rightHeader, rightMsg }) => (
  <Container flexDirection="row" alignItems="center" width="100%" py="10px">
    {logo}
    <Flex width="100%">
      <Flex flexDirection="column" ml="10px" flex={1}>
        {leftHeader}
        {leftMsg}
      </Flex>
      <Flex flexDirection="column" ml="10px" alignItems="flex-end" flex={1}>
        {rightHeader}
        {rightMsg}
      </Flex>
    </Flex>
  </Container>
)
