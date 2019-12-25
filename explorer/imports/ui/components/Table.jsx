import React from 'react'
import { Flex } from 'rebass'
import styled from 'styled-components'

const Table = styled(Flex)`
  border-radius: 3px;
  border: solid 1px #dadada;
  background-color: #ffffff;
  flex-direction: column;
  width: 100%;
`

const THead = styled(Flex)`
  padding: 12px 30px;
  font-size: 0.8em;
  color: black;
  font-weight: 600;
  border-bottom: solid 1px #dadada;
`

const TBody = styled(Flex)`
  width: 100%;
`

export default ({ header, body }) => {
  return (
    <Table>
      <THead>{header}</THead>
      <TBody>{body}</TBody>
    </Table>
  )
}
