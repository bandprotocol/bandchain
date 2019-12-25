import React from 'react'
import { Flex } from 'rebass'

export default ({ msg }) => (
  <Flex width="100%" justifyContent="center" alignItems="center" height="100%">
    {msg}
  </Flex>
)
