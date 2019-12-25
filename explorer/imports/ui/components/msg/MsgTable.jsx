import React from 'react'
import CreateReqScript from './CreateReqScript'
import Request from './Request'
import Report from './Report'

export default ({ msg }) => {
  switch (msg.action) {
    case 'CREATE_SCRIPT':
      return <CreateReqScript msg={msg} />
    case 'REQUEST':
      return <Request msg={msg} />
    case 'REPORT':
      return <Report msg={msg} />
    default:
      return <div>{JSON.stringify(msg)}</div>
  }
}
