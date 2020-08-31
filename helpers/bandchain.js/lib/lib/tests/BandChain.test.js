"use strict";

var _interopRequireDefault = require("@babel/runtime/helpers/interopRequireDefault");

var _regenerator = _interopRequireDefault(require("@babel/runtime/regenerator"));

var _asyncToGenerator2 = _interopRequireDefault(require("@babel/runtime/helpers/asyncToGenerator"));

var BandChain = require('../Bandchain');

jest.setTimeout(30000);
var endpoint = 'http://guanyu-devnet.bandchain.org/rest';
var mnemonic = 'final little loud vicious door hope differ lucky alpha morning clog oval milk repair off course indicate stumble remove nest position journey throw crane';
var testRequestID = 1;
it('Test BandChain constructor', function () {
  var bandchain = new BandChain(endpoint);
  expect(bandchain.endpoint).toBe(endpoint);
});
it('Test BandChain getOracleScript success', /*#__PURE__*/(0, _asyncToGenerator2["default"])( /*#__PURE__*/_regenerator["default"].mark(function _callee() {
  var oracleScriptID, bandchain, oracleScript;
  return _regenerator["default"].wrap(function _callee$(_context) {
    while (1) {
      switch (_context.prev = _context.next) {
        case 0:
          oracleScriptID = 1;
          bandchain = new BandChain(endpoint);
          _context.next = 4;
          return bandchain.getOracleScript(oracleScriptID);

        case 4:
          oracleScript = _context.sent;
          expect(JSON.stringify(oracleScript)).toBe(JSON.stringify({
            owner: 'band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs',
            name: 'Cryptocurrency Price in USD',
            description: 'Oracle script that queries the average cryptocurrency price using current price data from CoinGecko, CryptoCompare, and Binance',
            filename: 'ff48f3b9876cddb41e2371fe0fc5cc516619944adec75e91dfda85ace561dd9c',
            schema: '{symbol:string,multiplier:u64}/{px:u64}',
            source_code_url: 'https://ipfs.io/ipfs/QmY3S4dYuWMX4L7RMioEbUcLLZxc3tRoDMJxVQMthd7Amy',
            id: oracleScriptID
          }));

        case 6:
        case "end":
          return _context.stop();
      }
    }
  }, _callee);
})));
it('Test BandChain getOracleScript error', function () {
  var oracleScriptID = 1e18;
  var bandchain = new BandChain(endpoint);
  expect(bandchain.getOracleScript(oracleScriptID)).rejects.toThrow('No oracle script found with the given ID');
});
it('Test BandChain submitRequestTx', /*#__PURE__*/(0, _asyncToGenerator2["default"])( /*#__PURE__*/_regenerator["default"].mark(function _callee2() {
  var oracleScriptID, bandchain, oracleScript, requestID;
  return _regenerator["default"].wrap(function _callee2$(_context2) {
    while (1) {
      switch (_context2.prev = _context2.next) {
        case 0:
          oracleScriptID = 1;
          bandchain = new BandChain(endpoint);
          _context2.next = 4;
          return bandchain.getOracleScript(oracleScriptID);

        case 4:
          oracleScript = _context2.sent;
          _context2.next = 7;
          return bandchain.submitRequestTx(oracleScript, {
            symbol: 'BAND',
            multiplier: BigInt('1000000')
          }, {
            minCount: 2,
            askCount: 4
          }, mnemonic);

        case 7:
          requestID = _context2.sent;
          testRequestID = requestID;
          expect(requestID).toBeDefined();

        case 10:
        case "end":
          return _context2.stop();
      }
    }
  }, _callee2);
})));
it('Test BandChain getRequestID error', /*#__PURE__*/(0, _asyncToGenerator2["default"])( /*#__PURE__*/_regenerator["default"].mark(function _callee3() {
  var bandchain;
  return _regenerator["default"].wrap(function _callee3$(_context3) {
    while (1) {
      switch (_context3.prev = _context3.next) {
        case 0:
          bandchain = new BandChain(endpoint);
          expect(bandchain.getRequestID('13DEADCF273FCE723B809DDD6F29E5D0B5FD397256FD872D602676094061F20D' // Not a request tx
          )).rejects.toThrow('Not a request tx');

        case 2:
        case "end":
          return _context3.stop();
      }
    }
  }, _callee3);
})));
it('Test BandChain getRequestEVMProof', /*#__PURE__*/(0, _asyncToGenerator2["default"])( /*#__PURE__*/_regenerator["default"].mark(function _callee4() {
  var bandchain, requestProof;
  return _regenerator["default"].wrap(function _callee4$(_context4) {
    while (1) {
      switch (_context4.prev = _context4.next) {
        case 0:
          bandchain = new BandChain(endpoint);
          _context4.next = 3;
          return bandchain.getRequestEVMProof(testRequestID);

        case 3:
          requestProof = _context4.sent;
          expect(requestProof).toBeDefined();

        case 5:
        case "end":
          return _context4.stop();
      }
    }
  }, _callee4);
})));
it('Test BandChain getRequestNonEVMProof', /*#__PURE__*/(0, _asyncToGenerator2["default"])( /*#__PURE__*/_regenerator["default"].mark(function _callee5() {
  var bandchain, requestProof;
  return _regenerator["default"].wrap(function _callee5$(_context5) {
    while (1) {
      switch (_context5.prev = _context5.next) {
        case 0:
          bandchain = new BandChain(endpoint);
          _context5.next = 3;
          return bandchain.getRequestNonEVMProof(testRequestID);

        case 3:
          requestProof = _context5.sent;
          expect(requestProof).toBeDefined();

        case 5:
        case "end":
          return _context5.stop();
      }
    }
  }, _callee5);
})));
it('Test BandChain getRequestResult', /*#__PURE__*/(0, _asyncToGenerator2["default"])( /*#__PURE__*/_regenerator["default"].mark(function _callee6() {
  var bandchain, requestID, requestResult;
  return _regenerator["default"].wrap(function _callee6$(_context6) {
    while (1) {
      switch (_context6.prev = _context6.next) {
        case 0:
          bandchain = new BandChain(endpoint);
          requestID = 1;
          _context6.next = 4;
          return bandchain.getRequestResult(requestID);

        case 4:
          requestResult = _context6.sent;
          expect(requestResult).toBeDefined();

        case 6:
        case "end":
          return _context6.stop();
      }
    }
  }, _callee6);
})));
it('Test BandChain getLastMatchingRequestResult', /*#__PURE__*/(0, _asyncToGenerator2["default"])( /*#__PURE__*/_regenerator["default"].mark(function _callee7() {
  var oracleScriptID, bandchain, oracleScript, lastRequestResult;
  return _regenerator["default"].wrap(function _callee7$(_context7) {
    while (1) {
      switch (_context7.prev = _context7.next) {
        case 0:
          oracleScriptID = 1;
          bandchain = new BandChain(endpoint);
          _context7.next = 4;
          return bandchain.getOracleScript(oracleScriptID);

        case 4:
          oracleScript = _context7.sent;
          _context7.next = 7;
          return bandchain.getLastMatchingRequestResult(oracleScript, {
            symbol: 'BAND',
            multiplier: BigInt('1000000')
          }, {
            minCount: 2,
            askCount: 4
          });

        case 7:
          lastRequestResult = _context7.sent;
          expect(lastRequestResult).toBeDefined();
          expect(lastRequestResult.result).toBeDefined();
          expect(lastRequestResult.result.px).toBeGreaterThan(0);

        case 11:
        case "end":
          return _context7.stop();
      }
    }
  }, _callee7);
})));