import React from 'react'
import { Box } from 'rebass'

export default ({ size = '40px', isCircle = false }) => (
  <Box
    width={size}
    height={size}
    style={{
      borderRadius: isCircle ? '50%' : '3px',
      border: 'solid 1px #979797',
      backgroundColor: '#d8d8d8',
    }}
  />
)
