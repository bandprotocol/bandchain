import React from 'react'
import { Box } from 'rebass'
import Jazzicon, { jsNumberForAddress } from 'react-jazzicon'

export default ({
  size = 40,
  address = '0xBEEFF494F77294C7CBAea03ca7246C071497870B', // TODO: remove it later
}) => (
  <Box>
    <Jazzicon diameter={size} seed={jsNumberForAddress(address)} />
  </Box>
)
