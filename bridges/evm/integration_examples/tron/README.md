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

1. Install chrome extension [TronLink](https://chrome.google.com/webstore/detail/tronlink%EF%BC%88%E6%B3%A2%E5%AE%9D%E9%92%B1%E5%8C%85%EF%BC%89/ibnejdfjmmkpcnlpebklmnkoeoihofec)

2. Create your account with `TronLink`

3. Fund your account with [https://www.trongrid.io/shasta/#request](https://www.trongrid.io/shasta/#request)

4. Use [tronide](http://www.tronide.io/) for contract creation, compiling and deployment

5. Enable plugins using plugin manager

   - ![img](https://user-images.githubusercontent.com/12705423/91192951-c3502280-e720-11ea-8e81-9cc151d322d9.png)

6. Currently we still have a problem with encoding tuple as a parameter. So right now we have to modify our bridge by replace tuple input to bytes instead.

   - ![img](https://user-images.githubusercontent.com/12705423/91195479-3b1f4c80-e723-11ea-8c4f-c444ffbd9e32.png)

7. Deploy our bridge with this following bytes

   ```text
   00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000004000000000000000000000000634db8fc70651e63204f636a1c44d92f2689d16900000000000000000000000000000000000000000000000000000000000f4240000000000000000000000000d4d7099737da368d7178df89faa69d24aa9fdb0c00000000000000000000000000000000000000000000000000000000000f42400000000000000000000000003715ecfc3d04f358a18482147fd13c5f44fe8d7e00000000000000000000000000000000000000000000000000000000000f427200000000000000000000000013559c6a3709c25969dc9a5776ae19711c78cad100000000000000000000000000000000000000000000000000000000000f4240
   ```

   which is an abi-encode of

   ```
   [["0x634db8FC70651E63204f636a1c44D92f2689D169","1000000"],["0xd4D7099737dA368d7178DF89faA69d24aa9fdB0C","1000000"],["0x3715eCFC3D04f358A18482147Fd13C5f44FE8D7E","1000050"],["0x13559C6A3709C25969DC9a5776Ae19711C78CAD1","1000000"]]
   ```

8. Test our contract by calling `relay` function with a proof that copy from [https://guanyu-devnet.cosmoscan.io/oracle-script/76](https://guanyu-devnet.cosmoscan.io/oracle-script/76) and then try to read the result by calling `requestCache` function with `0xc108bbcdbdfe335f46444f46cfaa4f623270db239bbd1e2ea05ab15db47d10c0` as a parameter.

   - `0xc108bbcdbdfe335f46444f46cfaa4f623270db239bbd1e2ea05ab15db47d10c0` is a keccak256 of `["from_scan", 76, "0x00000003584147", 4,3]`

### Example deployed contracts

#### Bridge

[TPxsemS7h9rrJPZAPDjP7rmLoA4ErYny69](https://shasta.tronscan.org/#/contract/TPxsemS7h9rrJPZAPDjP7rmLoA4ErYny69)

#### Simple Price DB

A contract that consume data from Bridge contract via function `setPrice`.

[TWp5svfQxLesfbzqX9LREx4DYeqxwB9BTT](https://shasta.tronscan.org/#/contract/TWp5svfQxLesfbzqX9LREx4DYeqxwB9BTT)
