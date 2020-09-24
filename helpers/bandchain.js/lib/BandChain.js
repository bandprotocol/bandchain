'use strict'

var _interopRequireDefault = require('@babel/runtime/helpers/interopRequireDefault')

var _toConsumableArray2 = _interopRequireDefault(
  require('@babel/runtime/helpers/toConsumableArray'),
)

var _slicedToArray2 = _interopRequireDefault(
  require('@babel/runtime/helpers/slicedToArray'),
)

var _regenerator = _interopRequireDefault(require('@babel/runtime/regenerator'))

var _classCallCheck2 = _interopRequireDefault(
  require('@babel/runtime/helpers/classCallCheck'),
)

var _createClass2 = _interopRequireDefault(
  require('@babel/runtime/helpers/createClass'),
)

var _asyncToGenerator2 = _interopRequireDefault(
  require('@babel/runtime/helpers/asyncToGenerator'),
)

var axios = require('axios')

var _require = require('@bandprotocol/obi.js'),
  Obi = _require.Obi

var cosmosjs = require('@cosmostation/cosmosjs')

var delay = require('delay')

var _require2 = require('lodash'),
  result = _require2.result

function createRequestMsg(_x, _x2, _x3, _x4, _x5, _x6, _x7) {
  return _createRequestMsg.apply(this, arguments)
}

function _createRequestMsg() {
  _createRequestMsg = (0, _asyncToGenerator2['default'])(
    /*#__PURE__*/ _regenerator['default'].mark(function _callee13(
      cosmos,
      sender,
      oracleScriptID,
      validatorCounts,
      calldata,
      chainID,
      fee,
    ) {
      var account
      return _regenerator['default'].wrap(function _callee13$(_context13) {
        while (1) {
          switch ((_context13.prev = _context13.next)) {
            case 0:
              _context13.next = 2
              return cosmos.getAccounts(sender)

            case 2:
              account = _context13.sent
              return _context13.abrupt(
                'return',
                cosmos.newStdMsg({
                  msgs: [
                    {
                      type: 'oracle/Request',
                      value: {
                        oracle_script_id: String(oracleScriptID),
                        calldata: Buffer.from(calldata).toString('base64'),
                        ask_count: String(validatorCounts.askCount),
                        min_count: String(validatorCounts.minCount),
                        client_id: 'bandchain.js',
                        sender: sender,
                      },
                    },
                  ],
                  chain_id: chainID,
                  fee: fee,
                  memo: '',
                  account_number: String(account.result.value.account_number),
                  sequence: String(account.result.value.sequence || 0),
                }),
              )

            case 4:
            case 'end':
              return _context13.stop()
          }
        }
      }, _callee13)
    }),
  )
  return _createRequestMsg.apply(this, arguments)
}

var BandChain = /*#__PURE__*/ (function () {
  function BandChain(endpoint) {
    ;(0, _classCallCheck2['default'])(this, BandChain)
    this.endpoint = endpoint
  }

  ;(0, _createClass2['default'])(BandChain, [
    {
      key: '_getChainID',
      value: (function () {
        var _getChainID2 = (0, _asyncToGenerator2['default'])(
          /*#__PURE__*/ _regenerator['default'].mark(function _callee() {
            var res
            return _regenerator['default'].wrap(
              function _callee$(_context) {
                while (1) {
                  switch ((_context.prev = _context.next)) {
                    case 0:
                      if (this._chainID) {
                        _context.next = 11
                        break
                      }

                      _context.prev = 1
                      _context.next = 4
                      return axios.get(
                        ''.concat(this.endpoint, '/bandchain/genesis'),
                      )

                    case 4:
                      res = _context.sent
                      this._chainID = res.data.chain_id
                      _context.next = 11
                      break

                    case 8:
                      _context.prev = 8
                      _context.t0 = _context['catch'](1)
                      throw new Error('Cannot retrieve chainID')

                    case 11:
                      return _context.abrupt('return', this._chainID)

                    case 12:
                    case 'end':
                      return _context.stop()
                  }
                }
              },
              _callee,
              this,
              [[1, 8]],
            )
          }),
        )

        function _getChainID() {
          return _getChainID2.apply(this, arguments)
        }

        return _getChainID
      })(),
    },
    {
      key: 'getOracleScript',
      value: (function () {
        var _getOracleScript = (0, _asyncToGenerator2['default'])(
          /*#__PURE__*/ _regenerator['default'].mark(function _callee2(
            oracleScriptID,
          ) {
            var res
            return _regenerator['default'].wrap(
              function _callee2$(_context2) {
                while (1) {
                  switch ((_context2.prev = _context2.next)) {
                    case 0:
                      _context2.prev = 0
                      _context2.next = 3
                      return axios.get(
                        ''
                          .concat(this.endpoint, '/oracle/oracle_scripts/')
                          .concat(oracleScriptID),
                      )

                    case 3:
                      res = _context2.sent
                      res.data.result.id = oracleScriptID
                      return _context2.abrupt('return', res.data.result)

                    case 8:
                      _context2.prev = 8
                      _context2.t0 = _context2['catch'](0)
                      throw new Error(
                        'No oracle script found with the given ID',
                      )

                    case 11:
                    case 'end':
                      return _context2.stop()
                  }
                }
              },
              _callee2,
              this,
              [[0, 8]],
            )
          }),
        )

        function getOracleScript(_x8) {
          return _getOracleScript.apply(this, arguments)
        }

        return getOracleScript
      })(),
    },
    {
      key: 'submitRequestTx',
      value: (function () {
        var _submitRequestTx = (0, _asyncToGenerator2['default'])(
          /*#__PURE__*/ _regenerator['default'].mark(function _callee3(
            oracleScript,
            parameters,
            validatorCounts,
            mnemonic,
          ) {
            var gasAmount,
              gasLimit,
              chainID,
              obiObj,
              calldata,
              cosmos,
              ecpairPriv,
              sender,
              requestMsg,
              signedTx,
              broadcastResponse,
              _args3 = arguments
            return _regenerator['default'].wrap(
              function _callee3$(_context3) {
                while (1) {
                  switch ((_context3.prev = _context3.next)) {
                    case 0:
                      gasAmount =
                        _args3.length > 4 && _args3[4] !== undefined
                          ? _args3[4]
                          : 0
                      gasLimit =
                        _args3.length > 5 && _args3[5] !== undefined
                          ? _args3[5]
                          : 1000000
                      _context3.next = 4
                      return this._getChainID()

                    case 4:
                      chainID = _context3.sent
                      obiObj = new Obi(oracleScript.schema)
                      calldata = obiObj.encodeInput(parameters)
                      cosmos = cosmosjs.network(this.endpoint, chainID)
                      cosmos.setPath("m/44'/494'/0'/0/0")
                      cosmos.setBech32MainPrefix('band')
                      ecpairPriv = cosmos.getECPairPriv(mnemonic)
                      sender = cosmos.getAddress(mnemonic)
                      _context3.next = 14
                      return createRequestMsg(
                        cosmos,
                        sender,
                        oracleScript.id,
                        validatorCounts,
                        calldata,
                        chainID,
                        {
                          amount: [
                            {
                              amount: ''.concat(gasAmount),
                              denom: 'uband',
                            },
                          ],
                          gas: ''.concat(gasLimit),
                        },
                      )

                    case 14:
                      requestMsg = _context3.sent
                      signedTx = cosmos.sign(requestMsg, ecpairPriv, 'block')
                      _context3.next = 18
                      return cosmos.broadcast(signedTx)

                    case 18:
                      broadcastResponse = _context3.sent
                      return _context3.abrupt(
                        'return',
                        this.getRequestID(broadcastResponse.txhash),
                      )

                    case 20:
                    case 'end':
                      return _context3.stop()
                  }
                }
              },
              _callee3,
              this,
            )
          }),
        )

        function submitRequestTx(_x9, _x10, _x11, _x12) {
          return _submitRequestTx.apply(this, arguments)
        }

        return submitRequestTx
      })(),
    },
    {
      key: 'getRequestID',
      value: (function () {
        var _getRequestID = (0, _asyncToGenerator2['default'])(
          /*#__PURE__*/ _regenerator['default'].mark(function _callee4(txHash) {
            var retryTimeout,
              requestEndpoint,
              res,
              requestID,
              _args4 = arguments
            return _regenerator['default'].wrap(
              function _callee4$(_context4) {
                while (1) {
                  switch ((_context4.prev = _context4.next)) {
                    case 0:
                      retryTimeout =
                        _args4.length > 1 && _args4[1] !== undefined
                          ? _args4[1]
                          : 200
                      requestEndpoint = ''
                        .concat(this.endpoint, '/txs/')
                        .concat(txHash) // Loop until the txHash is included in the block

                    case 2:
                      if (!true) {
                        _context4.next = 26
                        break
                      }

                      res = void 0
                      _context4.prev = 4
                      _context4.next = 7
                      return axios.get(requestEndpoint)

                    case 7:
                      res = _context4.sent
                      _context4.next = 15
                      break

                    case 10:
                      _context4.prev = 10
                      _context4.t0 = _context4['catch'](4)
                      _context4.next = 14
                      return delay(retryTimeout)

                    case 14:
                      return _context4.abrupt('continue', 2)

                    case 15:
                      if (!(res.status == 200)) {
                        _context4.next = 24
                        break
                      }

                      _context4.prev = 16
                      requestID = res.data.logs[0].events
                        .find(function (_ref) {
                          var type = _ref.type
                          return type === 'request'
                        })
                        .attributes.find(function (_ref2) {
                          var key = _ref2.key
                          return key === 'id'
                        }).value
                      return _context4.abrupt('return', requestID)

                    case 21:
                      _context4.prev = 21
                      _context4.t1 = _context4['catch'](16)
                      throw new Error('Not a request tx')

                    case 24:
                      _context4.next = 2
                      break

                    case 26:
                    case 'end':
                      return _context4.stop()
                  }
                }
              },
              _callee4,
              this,
              [
                [4, 10],
                [16, 21],
              ],
            )
          }),
        )

        function getRequestID(_x13) {
          return _getRequestID.apply(this, arguments)
        }

        return getRequestID
      })(),
    },
    {
      key: 'getRequestProof',
      value: (function () {
        var _getRequestProof = (0, _asyncToGenerator2['default'])(
          /*#__PURE__*/ _regenerator['default'].mark(function _callee5(
            requestID,
          ) {
            var retryTimeout,
              requestEndpoint,
              _requestEndpoint,
              res,
              _result,
              _args5 = arguments

            return _regenerator['default'].wrap(
              function _callee5$(_context5) {
                while (1) {
                  switch ((_context5.prev = _context5.next)) {
                    case 0:
                      retryTimeout =
                        _args5.length > 1 && _args5[1] !== undefined
                          ? _args5[1]
                          : 200
                      _context5.prev = 1
                      requestEndpoint = ''
                        .concat(this.endpoint, '/oracle/requests/')
                        .concat(requestID)
                      _context5.next = 5
                      return axios.get(requestEndpoint)

                    case 5:
                      _context5.next = 10
                      break

                    case 7:
                      _context5.prev = 7
                      _context5.t0 = _context5['catch'](1)
                      throw new Error('Request not found')

                    case 10:
                      if (!true) {
                        _context5.next = 31
                        break
                      }

                      _context5.prev = 11
                      _requestEndpoint = ''
                        .concat(this.endpoint, '/oracle/proof/')
                        .concat(requestID)
                      _context5.next = 15
                      return axios.get(_requestEndpoint)

                    case 15:
                      res = _context5.sent

                      if (
                        !(
                          res.status == 200 &&
                          res.data.result.evmProofBytes != null
                        )
                      ) {
                        _context5.next = 21
                        break
                      }

                      _result = res.data.result
                      return _context5.abrupt('return', _result)

                    case 21:
                      if (
                        !(
                          res.status == 200 &&
                          res.data.result.evmProofBytes == null
                        )
                      ) {
                        _context5.next = 23
                        break
                      }

                      throw new Error(
                        'No proof found for the specified requestID',
                      )

                    case 23:
                      _context5.next = 29
                      break

                    case 25:
                      _context5.prev = 25
                      _context5.t1 = _context5['catch'](11)
                      _context5.next = 29
                      return delay(retryTimeout)

                    case 29:
                      _context5.next = 10
                      break

                    case 31:
                    case 'end':
                      return _context5.stop()
                  }
                }
              },
              _callee5,
              this,
              [
                [1, 7],
                [11, 25],
              ],
            )
          }),
        )

        function getRequestProof(_x14) {
          return _getRequestProof.apply(this, arguments)
        }

        return getRequestProof
      })(),
    },
    {
      key: 'getRequestEVMProof',
      value: (function () {
        var _getRequestEVMProof = (0, _asyncToGenerator2['default'])(
          /*#__PURE__*/ _regenerator['default'].mark(function _callee6(
            requestID,
          ) {
            var retryTimeout,
              result,
              _args6 = arguments
            return _regenerator['default'].wrap(
              function _callee6$(_context6) {
                while (1) {
                  switch ((_context6.prev = _context6.next)) {
                    case 0:
                      retryTimeout =
                        _args6.length > 1 && _args6[1] !== undefined
                          ? _args6[1]
                          : 200
                      _context6.next = 3
                      return this.getRequestProof(requestID, retryTimeout)

                    case 3:
                      result = _context6.sent
                      return _context6.abrupt('return', result.evmProofBytes)

                    case 5:
                    case 'end':
                      return _context6.stop()
                  }
                }
              },
              _callee6,
              this,
            )
          }),
        )

        function getRequestEVMProof(_x15) {
          return _getRequestEVMProof.apply(this, arguments)
        }

        return getRequestEVMProof
      })(),
    },
    {
      key: 'getRequestNonEVMProof',
      value: (function () {
        var _getRequestNonEVMProof = (0, _asyncToGenerator2['default'])(
          /*#__PURE__*/ _regenerator['default'].mark(function _callee7(
            requestID,
          ) {
            var retryTimeout,
              proofSchema,
              result,
              requestPacket,
              responsePacket,
              obiObj,
              proof,
              _args7 = arguments
            return _regenerator['default'].wrap(
              function _callee7$(_context7) {
                while (1) {
                  switch ((_context7.prev = _context7.next)) {
                    case 0:
                      retryTimeout =
                        _args7.length > 1 && _args7[1] !== undefined
                          ? _args7[1]
                          : 200
                      proofSchema =
                        '{client_id:string,oracle_script_id:u64,calldata:bytes,ask_count:u64,min_count:u64}/{client_id:string,request_id:u64,ans_count:u64,request_time:u64,resolve_time:u64,resolve_status:u8,result:bytes}'
                      _context7.next = 4
                      return this.getRequestProof(requestID, retryTimeout)

                    case 4:
                      result = _context7.sent
                      requestPacket =
                        result.jsonProof.oracleDataProof.requestPacket
                      responsePacket =
                        result.jsonProof.oracleDataProof.responsePacket
                      requestPacket.calldata = Buffer.from(
                        requestPacket.calldata,
                        'base64',
                      )
                      responsePacket.result = Buffer.from(
                        responsePacket.result,
                        'base64',
                      )
                      obiObj = new Obi(proofSchema)
                      proof = Buffer.concat([
                        obiObj.encodeInput(requestPacket),
                        obiObj.encodeOutput(responsePacket),
                      ]).toString('hex')
                      return _context7.abrupt('return', proof)

                    case 12:
                    case 'end':
                      return _context7.stop()
                  }
                }
              },
              _callee7,
              this,
            )
          }),
        )

        function getRequestNonEVMProof(_x16) {
          return _getRequestNonEVMProof.apply(this, arguments)
        }

        return getRequestNonEVMProof
      })(),
    },
    {
      key: 'getRequestResult',
      value: (function () {
        var _getRequestResult = (0, _asyncToGenerator2['default'])(
          /*#__PURE__*/ _regenerator['default'].mark(function _callee8(
            requestID,
          ) {
            var requestEndpoint, res
            return _regenerator['default'].wrap(
              function _callee8$(_context8) {
                while (1) {
                  switch ((_context8.prev = _context8.next)) {
                    case 0:
                      if (!true) {
                        _context8.next = 19
                        break
                      }

                      _context8.prev = 1
                      requestEndpoint = ''
                        .concat(this.endpoint, '/oracle/requests/')
                        .concat(requestID)
                      _context8.next = 5
                      return axios.get(requestEndpoint)

                    case 5:
                      res = _context8.sent

                      if (
                        !(res.status == 200 && res.data.result.result != null)
                      ) {
                        _context8.next = 10
                        break
                      }

                      return _context8.abrupt('return', res.data.result.result)

                    case 10:
                      if (
                        !(res.status == 200 && res.data.result.request == null)
                      ) {
                        _context8.next = 12
                        break
                      }

                      throw new Error(
                        'No result found for the specified requestID',
                      )

                    case 12:
                      _context8.next = 17
                      break

                    case 14:
                      _context8.prev = 14
                      _context8.t0 = _context8['catch'](1)
                      throw new Error('Error querying the request result')

                    case 17:
                      _context8.next = 0
                      break

                    case 19:
                    case 'end':
                      return _context8.stop()
                  }
                }
              },
              _callee8,
              this,
              [[1, 14]],
            )
          }),
        )

        function getRequestResult(_x17) {
          return _getRequestResult.apply(this, arguments)
        }

        return getRequestResult
      })(),
    },
    {
      key: 'getLastMatchingRequestResult',
      value: (function () {
        var _getLastMatchingRequestResult = (0, _asyncToGenerator2['default'])(
          /*#__PURE__*/ _regenerator['default'].mark(function _callee9(
            oracleScript,
            parameters,
            validatorCounts,
          ) {
            var obiObj, calldata, requestEndpoint, res, response
            return _regenerator['default'].wrap(
              function _callee9$(_context9) {
                while (1) {
                  switch ((_context9.prev = _context9.next)) {
                    case 0:
                      obiObj = new Obi(oracleScript.schema)
                      calldata = Buffer.from(
                        obiObj.encodeInput(parameters),
                      ).toString('hex')
                      requestEndpoint = ''
                        .concat(this.endpoint, '/oracle/request_search?oid=')
                        .concat(oracleScript.id, '&calldata=')
                        .concat(calldata, '&min_count=')
                        .concat(validatorCounts.minCount, '&ask_count=')
                        .concat(validatorCounts.askCount)
                      _context9.prev = 3
                      _context9.next = 6
                      return axios.get(requestEndpoint)

                    case 6:
                      res = _context9.sent

                      if (
                        !(res.status == 200 && res.data.result.result != null)
                      ) {
                        _context9.next = 13
                        break
                      }

                      response = res.data.result.result.response_packet_data
                      response.result = obiObj.decodeOutput(
                        Buffer.from(response.result, 'base64'),
                      )
                      return _context9.abrupt('return', response)

                    case 13:
                      return _context9.abrupt('return', null)

                    case 14:
                      _context9.next = 19
                      break

                    case 16:
                      _context9.prev = 16
                      _context9.t0 = _context9['catch'](3)
                      throw new Error(
                        'Error querying the latest matching request result',
                      )

                    case 19:
                    case 'end':
                      return _context9.stop()
                  }
                }
              },
              _callee9,
              this,
              [[3, 16]],
            )
          }),
        )

        function getLastMatchingRequestResult(_x18, _x19, _x20) {
          return _getLastMatchingRequestResult.apply(this, arguments)
        }

        return getLastMatchingRequestResult
      })(),
    },
    {
      key: 'getLatestValue',
      value: (function () {
        var _getLatestValue = (0, _asyncToGenerator2['default'])(
          /*#__PURE__*/ _regenerator['default'].mark(function _callee10(
            oracleScriptID,
            parameters,
          ) {
            var minCount, askCount, validatorCounts, oracleScript, latestValue
            return _regenerator['default'].wrap(
              function _callee10$(_context10) {
                while (1) {
                  switch ((_context10.prev = _context10.next)) {
                    case 0:
                      minCount = 3
                      askCount = 4
                      validatorCounts = {
                        minCount: minCount,
                        askCount: askCount,
                      }
                      _context10.next = 5
                      return this.getOracleScript(oracleScriptID)

                    case 5:
                      oracleScript = _context10.sent
                      _context10.next = 8
                      return this.getLastMatchingRequestResult(
                        oracleScript,
                        parameters,
                        validatorCounts,
                      )

                    case 8:
                      latestValue = _context10.sent
                      return _context10.abrupt('return', latestValue)

                    case 10:
                    case 'end':
                      return _context10.stop()
                  }
                }
              },
              _callee10,
              this,
            )
          }),
        )

        function getLatestValue(_x21, _x22) {
          return _getLatestValue.apply(this, arguments)
        }

        return getLatestValue
      })(),
    },
    {
      key: 'getReferenceData',
      value: (function () {
        var _getReferenceData = (0, _asyncToGenerator2['default'])(
          /*#__PURE__*/ _regenerator['default'].mark(function _callee12(pairs) {
            var _this = this

            var sourceList, set, symbolDict, data
            return _regenerator['default'].wrap(function _callee12$(
              _context12,
            ) {
              while (1) {
                switch ((_context12.prev = _context12.next)) {
                  case 0:
                    sourceList = [
                      {
                        id: 8,
                        exponent: 9,
                        symbols: [
                          'BTC',
                          'ETH',
                          'USDT',
                          'XRP',
                          'LINK',
                          'DOT',
                          'BCH',
                          'LTC',
                          'ADA',
                          'BSV',
                          'CRO',
                          'BNB',
                          'EOS',
                          'XTZ',
                          'TRX',
                          'XLM',
                          'ATOM',
                          'XMR',
                          'OKB',
                          'USDC',
                          'NEO',
                          'XEM',
                          'LEO',
                          'HT',
                          'VET',
                        ],
                      },
                      {
                        id: 8,
                        exponent: 9,
                        symbols: [
                          'YFI',
                          'MIOTA',
                          'LEND',
                          'SNX',
                          'DASH',
                          'COMP',
                          'ZEC',
                          'ETC',
                          'OMG',
                          'MKR',
                          'ONT',
                          'NXM',
                          'AMPL',
                          'BAT',
                          'THETA',
                          'DAI',
                          'REN',
                          'ZRX',
                          'ALGO',
                          'FTT',
                          'DOGE',
                          'KSM',
                          'WAVES',
                          'EWT',
                          'DGB',
                        ],
                      },
                      {
                        id: 8,
                        exponent: 9,
                        symbols: [
                          'KNC',
                          'ICX',
                          'TUSD',
                          'SUSHI',
                          'BTT',
                          'BAND',
                          'EGLD',
                          'ANT',
                          'NMR',
                          'PAX',
                          'LSK',
                          'LRC',
                          'HBAR',
                          'BAL',
                          'RUNE',
                          'YFII',
                          'LUNA',
                          'DCR',
                          'SC',
                          'STX',
                          'ENJ',
                          'BUSD',
                          'OCEAN',
                          'RSR',
                          'SXP',
                        ],
                      },
                      {
                        id: 8,
                        exponent: 9,
                        symbols: [
                          'BTG',
                          'BZRX',
                          'SRM',
                          'SNT',
                          'SOL',
                          'CKB',
                          'BNT',
                          'CRV',
                          'MANA',
                          'YFV',
                          'KAVA',
                          'MATIC',
                          'TRB',
                          'REP',
                          'FTM',
                          'TOMO',
                          'ONE',
                          'WNXM',
                          'PAXG',
                          'WAN',
                          'SUSD',
                          'RLC',
                          'OXT',
                          'RVN',
                          'NANO',
                        ],
                      },
                      {
                        id: 9,
                        exponent: 9,
                        symbols: [
                          'EUR',
                          'GBP',
                          'CNY',
                          'SGD',
                          'RMB',
                          'KRW',
                          'JPY',
                          'INR',
                          'RUB',
                          'CHF',
                          'AUD',
                          'BRL',
                          'CAD',
                          'HKD',
                          'XAU',
                          'XAG',
                        ],
                      },
                    ]
                    set = new Set()
                    pairs.forEach(function (pair) {
                      var _pair$split = pair.split('/'),
                        _pair$split2 = (0, _slicedToArray2['default'])(
                          _pair$split,
                          2,
                        ),
                        baseSymbol = _pair$split2[0],
                        quoteSymbol = _pair$split2[1]

                      sourceList.forEach(function (_ref3, index) {
                        var symbols = _ref3.symbols

                        if (
                          symbols.includes(baseSymbol) ||
                          symbols.includes(quoteSymbol)
                        ) {
                          set.add(index)
                        }
                      })
                    })
                    symbolDict = {}
                    _context12.next = 6
                    return Promise.all(
                      (0, _toConsumableArray2['default'])(set).map(
                        /*#__PURE__*/ (function () {
                          var _ref4 = (0, _asyncToGenerator2['default'])(
                            /*#__PURE__*/ _regenerator['default'].mark(
                              function _callee11(index) {
                                var result
                                return _regenerator['default'].wrap(
                                  function _callee11$(_context11) {
                                    while (1) {
                                      switch (
                                        (_context11.prev = _context11.next)
                                      ) {
                                        case 0:
                                          _context11.next = 2
                                          return _this.getLatestValue(
                                            sourceList[index].id,
                                            {
                                              symbols:
                                                sourceList[index].symbols,
                                              multiplier: Math.pow(
                                                10,
                                                sourceList[index].exponent,
                                              ),
                                            },
                                          )

                                        case 2:
                                          result = _context11.sent
                                          return _context11.abrupt(
                                            'return',
                                            sourceList[index].symbols.map(
                                              function (symbol, id) {
                                                symbolDict[symbol] = {
                                                  value:
                                                    result.result.rates[id],
                                                  updated: result.resolve_time,
                                                  decimals:
                                                    sourceList[index].exponent,
                                                }
                                              },
                                            ),
                                          )

                                        case 4:
                                        case 'end':
                                          return _context11.stop()
                                      }
                                    }
                                  },
                                  _callee11,
                                )
                              },
                            ),
                          )

                          return function (_x24) {
                            return _ref4.apply(this, arguments)
                          }
                        })(),
                      ),
                    )

                  case 6:
                    data = []
                    pairs.forEach(function (pair) {
                      var _pair$split3 = pair.split('/'),
                        _pair$split4 = (0, _slicedToArray2['default'])(
                          _pair$split3,
                          2,
                        ),
                        baseSymbol = _pair$split4[0],
                        quoteSymbol = _pair$split4[1]

                      if (baseSymbol == 'USD' && quoteSymbol == 'USD') {
                        data.push({
                          pair: pair,
                          rate: 1.0,
                          updated: {
                            base: 0,
                            quote: 0,
                          },
                          rawRate: {
                            value: BigInt(1e9),
                            decimals: 9,
                          },
                        })
                      } else if (baseSymbol == 'USD') {
                        var rate =
                          Math.pow(10, symbolDict[quoteSymbol].decimals) /
                          Number(symbolDict[quoteSymbol].value)
                        data.push({
                          pair: pair,
                          rate: rate,
                          updated: {
                            base: 0,
                            quote: Number(symbolDict[quoteSymbol].updated),
                          },
                          rawRate: {
                            value: BigInt(
                              BigInt(
                                Math.pow(
                                  10,
                                  symbolDict[quoteSymbol].decimals + 9,
                                ),
                              ) / symbolDict[quoteSymbol].value,
                            ),
                            decimals: 9,
                          },
                        })
                      } else if (quoteSymbol == 'USD') {
                        data.push({
                          pair: pair,
                          rate:
                            Number(symbolDict[baseSymbol].value) /
                            Math.pow(10, symbolDict[baseSymbol].decimals),
                          updated: {
                            base: Number(symbolDict[baseSymbol].updated),
                            quote: 0,
                          },
                          rawRate: {
                            value: symbolDict[baseSymbol].value,
                            decimals: symbolDict[baseSymbol].decimals,
                          },
                        })
                      } else {
                        data.push({
                          pair: pair,
                          rate:
                            Number(symbolDict[baseSymbol].value) /
                            Number(symbolDict[quoteSymbol].value),
                          updated: {
                            base: Number(symbolDict[baseSymbol].updated),
                            quote: Number(symbolDict[quoteSymbol].updated),
                          },
                          rawRate: {
                            value:
                              (symbolDict[baseSymbol].value * BigInt(1e9)) /
                              symbolDict[quoteSymbol].value,
                            decimals: 9,
                          },
                        })
                      }
                    })
                    return _context12.abrupt('return', data)

                  case 9:
                  case 'end':
                    return _context12.stop()
                }
              }
            },
            _callee12)
          }),
        )

        function getReferenceData(_x23) {
          return _getReferenceData.apply(this, arguments)
        }

        return getReferenceData
      })(),
    },
  ])
  return BandChain
})()

module.exports = BandChain
