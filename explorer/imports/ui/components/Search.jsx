import React from 'react'
import { Image, Flex } from 'rebass'
import styled from 'styled-components'

const Input = styled.input`
  width: 385px;
  height: 40px;
  font-size: 14px;
  padding-left: 11px;
  padding-top: 0;
  padding-bottom: 0;
  border-radius: 2px;
  border: solid 1px #dbdbdb;
`

export default () => {
  /* TODO: handle input and search */

  return (
    <Flex>
      <Input placeholder="Search by Address / Txn Hash / Block / Token / Ens" />
      <Flex
        width="40px"
        height="42px"
        justifyContent="center"
        alignItems="center"
        bg="black"
      >
        <Image src="/img/search-input-icon.svg" width={18} />
      </Flex>
    </Flex>
  )
}
