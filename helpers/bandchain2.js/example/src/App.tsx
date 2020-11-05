import React from 'react'
import logo from './logo.svg'
import './App.css'
import { Message, Data } from 'bandchain2.js'

function App() {
  const { MsgSend } = Message
  const { Coin } = Data
  const msgSend = new MsgSend('aaa', 'aaa', [new Coin(100000, 'uband')])
  console.log(msgSend)
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
