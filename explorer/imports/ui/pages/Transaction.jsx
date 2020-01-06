import React from 'react'
import { Flex } from 'rebass'
import TxDetail from '/imports/ui/components/TxDetail'
import MsgDetail from '/imports/ui/components/msg/MsgDetail'

export default props => {
  const txHash = props.match.params.txHash
  return (
    <Flex flexDirection="column" width="100%" py="20px">
      <TxDetail txHash={txHash} />
      <MsgDetail txHash={txHash} />
    </Flex>
  )
}
