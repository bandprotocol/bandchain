import React from 'react'
import styled from 'styled-components'
import { Flex } from 'rebass'

const Container = styled.div`
  display: flex;
  max-width: 1180px;
  width: 100%;
`

export default ({ children, content }) => (
  <Flex
    justifyContent="center"
    width="100%"
    style={{ minHeight: content ? 'calc(100vh - 380px)' : 'unset' }}
  >
    <Container>{children}</Container>
  </Flex>
)
