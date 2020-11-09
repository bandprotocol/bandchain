import React from 'react'
import logo from './logo.svg'
import './App.css'
import { Message, Data } from 'bandchain2.js'

function App() {
  const { MsgSend, MsgRequest } = Message
  const { Coin } = Data
  const msgSend = new MsgSend('aaa', 'aaa', [new Coin(100000, 'uband')])
  const msgRequest = new MsgRequest(1, "000000034254430000000000000001", 4, 2, "from_bandchain2.js", "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c")
  console.log(msgSend)
  console.log(msgRequest)
  console.log(msgRequest.asJson())
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <a
          className="App-link"
          href="https://reactjs.org"
          target="_blank"
          rel="noopener noreferrer"
        >
          Learn React
        </a>
      </header>
    </div>
  )
}

export default App
