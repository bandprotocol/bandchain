import React from 'react'
import { Flex } from 'rebass'

export default props => (
  <Flex justifyContent="center" alignItems="center" width="100%">
    Transaction Details {props.match.params.txHash}
  </Flex>
)
