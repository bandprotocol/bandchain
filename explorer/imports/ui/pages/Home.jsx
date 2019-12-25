import React from 'react'
import { Flex } from 'rebass'
import styled from 'styled-components'
import Summary from '/imports/ui/components/Summary'
import LatestBlocks from '/imports/ui/components/LatestBlocks'
import LatestTxs from '/imports/ui/components/LatestTxs'

const TableContainer = styled(Flex).attrs(() => ({
  mt: '25px',
  width: '100%',
}))``

export default () => (
  <Flex width="100%" alignItems="center" flexDirection="column" py="20px">
    <Summary />
    <TableContainer>
      <Flex flex={1}>
        <LatestBlocks />
      </Flex>
      <Flex width="20px" />
      <Flex flex={1}>
        <LatestTxs />
      </Flex>
    </TableContainer>
  </Flex>
)
