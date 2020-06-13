/*!
 *
 *   bandchain.js v1.0.0
 *   https://github.com/bandprotocol/bandchain
 *
 *   Copyright (c) Band Protocol (https://github.com/bandprotocol)
 *
 *   This source code is licensed under the MIT license found in the
 *   LICENSE file in the root directory of this source tree.
 *
 */
!(function (e, r) {
  'object' === typeof exports && 'object' === typeof module
    ? (module.exports = r())
    : 'function' === typeof define && define.amd
    ? define('BandChain', [], r)
    : 'object' === typeof exports
    ? (exports.BandChain = r())
    : (e.BandChain = r())
})(window, function () {
  return (function (e) {
    var r = {}
    function __webpack_require__(t) {
      if (r[t]) return r[t].exports
      var n = (r[t] = { i: t, l: !1, exports: {} })
      return (
        e[t].call(n.exports, n, n.exports, __webpack_require__),
        (n.l = !0),
        n.exports
      )
    }
    return (
      (__webpack_require__.m = e),
      (__webpack_require__.c = r),
      (__webpack_require__.d = function (e, r, t) {
        __webpack_require__.o(e, r) ||
          Object.defineProperty(e, r, { enumerable: !0, get: t })
      }),
      (__webpack_require__.r = function (e) {
        'undefined' !== typeof Symbol &&
          Symbol.toStringTag &&
          Object.defineProperty(e, Symbol.toStringTag, { value: 'Module' }),
          Object.defineProperty(e, '__esModule', { value: !0 })
      }),
      (__webpack_require__.t = function (e, r) {
        if ((1 & r && (e = __webpack_require__(e)), 8 & r)) return e
        if (4 & r && 'object' === typeof e && e && e.__esModule) return e
        var t = Object.create(null)
        if (
          (__webpack_require__.r(t),
          Object.defineProperty(t, 'default', { enumerable: !0, value: e }),
          2 & r && 'string' != typeof e)
        )
          for (var n in e)
            __webpack_require__.d(
              t,
              n,
              function (r) {
                return e[r]
              }.bind(null, n),
            )
        return t
      }),
      (__webpack_require__.n = function (e) {
        var r =
          e && e.__esModule
            ? function () {
                return e.default
              }
            : function () {
                return e
              }
        return __webpack_require__.d(r, 'a', r), r
      }),
      (__webpack_require__.o = function (e, r) {
        return Object.prototype.hasOwnProperty.call(e, r)
      }),
      (__webpack_require__.p = ''),
      __webpack_require__((__webpack_require__.s = 0))
    )
  })([
    function (e, r, t) {
      e.exports = t(2)
    },
    function (e, r, t) {},
    function (e, r, t) {
      'use strict'
      t.r(r)
      t(1)
      function _defineProperty(e, r, t) {
        return (
          r in e
            ? Object.defineProperty(e, r, {
                value: t,
                enumerable: !0,
                configurable: !0,
                writable: !0,
              })
            : (e[r] = t),
          e
        )
      }
      var n = function App() {
        !(function (e, r) {
          if (!(e instanceof r))
            throw new TypeError('Cannot call a class as a function')
        })(this, App),
          _defineProperty(this, 'myVar', !0),
          _defineProperty(this, 'myArrowMethod', function () {
            console.log('Arrow method fired')
          })
        var e = this.myArrowMethod,
          r = this.myVar
        console.log('Lib constructor called', r), e()
      }
      r.default = n
    },
  ])
})
//# sourceMappingURL=index.js.map
