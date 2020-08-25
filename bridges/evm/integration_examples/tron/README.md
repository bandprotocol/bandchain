# Tron Deployment And Interaction

### Official Documentation

[https://developers.tron.network/docs/deploying](https://developers.tron.network/docs/deploying)

### RPC Enpoint

```text
https://api.shasta.trongrid.io
```

### IDE

[tronide](http://www.tronide.io/)

### Steps

1. install chrome extension [TronLink](https://chrome.google.com/webstore/detail/tronlink%EF%BC%88%E6%B3%A2%E5%AE%9D%E9%92%B1%E5%8C%85%EF%BC%89/ibnejdfjmmkpcnlpebklmnkoeoihofec)

2. create your account with `TronLink`

3. fund your account with [https://www.trongrid.io/shasta/#request](https://www.trongrid.io/shasta/#request)

4. use [tronide](http://www.tronide.io/) for contract creation, compiling and deployment

5. enable plugins using plugin manager

   - ![img](https://user-images.githubusercontent.com/12705423/91192951-c3502280-e720-11ea-8e81-9cc151d322d9.png)

6. currently we still have a problem with encoding tuple as a parameter. So we have to modify our bridge by replace tuple input to bytes instead.

   - ![img](https://user-images.githubusercontent.com/12705423/91195479-3b1f4c80-e723-11ea-8c4f-c444ffbd9e32.png)

7. deploy our bridge with this following bytes

   ```text
   0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000070000000000000000000000004705e5935acb7c6eb23a49b17767e01fb7d8ab7300000000000000000000000000000000000000000000000000000000000f42410000000000000000000000001c18c1ed4ec5d1dd39b15c26b3346fc0b4adb13800000000000000000000000000000000000000000000000000000000000f4241000000000000000000000000a704058b50414c47cb2dff32929d0830bbcb7f3900000000000000000000000000000000000000000000000000000000000f424100000000000000000000000003899178b4994a3e8ea88c69b0014fdf71c656c600000000000000000000000000000000000000000000000000000000000f4241000000000000000000000000eb666c2755aaa9c7e69d6e29f57d8416938c7ee600000000000000000000000000000000000000000000000000000000000f4241000000000000000000000000975e60e8544bbb2945d6630553c269c5c8e7298000000000000000000000000000000000000000000000000000000000000f424100000000000000000000000053256a45877e8cc3a868d53dc4a4493256a0170100000000000000000000000000000000000000000000000000000000000f4241
   ```

   which is an abi-encode of

   ```
   [["0x4705e5935aCb7c6eb23A49b17767e01fB7D8Ab73","1000001"],["0x1c18c1Ed4ec5d1Dd39b15c26b3346Fc0b4ADb138","1000001"],["0xa704058b50414C47cb2DFf32929D0830bbcb7f39","1000001"],["0x03899178B4994a3E8EA88c69b0014FDf71c656C6","1000001"],["0xeb666C2755AaA9C7e69D6e29f57d8416938c7ee6","1000001"],["0x975E60E8544bbb2945d6630553C269c5c8E72980","1000001"],["0x53256a45877e8Cc3A868D53Dc4a4493256a01701","1000001"]]
   ```

8. test our contract by calling `relay` function with a proof that copy from [https://guanyu-poa.cosmoscan.io/request/53465](https://guanyu-poa.cosmoscan.io/request/53465)

9. try to read the of the contract by calling `requestCache` function with this `0x701c27c511f035a3aacf880d8eae7b7b51a8a3dcd4aeaafc77baf34e339b5937` as a parameter.

   - `0x701c27c511f035a3aacf880d8eae7b7b51a8a3dcd4aeaafc77baf34e339b5937` is a keccak256 of `["from_scan", 5, "0x00000003494358000000035553440000000000000064", 4, 4]`
