const { expectRevert } = require("openzeppelin-test-helpers");
const Bridge = artifacts.require("BridgeMock");

require("chai").should();

contract("Bridge", () => {
  context("Checking oracle state relay (4 validators)", () => {
    beforeEach(async () => {
      this.bridge = await Bridge.new([
        "0x652D89a66Eb4eA55366c45b1f9ACfc8e2179E1c5",
        "0x88e1cd00710495EEB93D4f522d16bC8B87Cb00FE",
        "0xaAA22E077492CbaD414098EBD98AA8dc1C7AE8D9",
        "0xB956589b6fC5523eeD0d9eEcfF06262Ce84ff260"
      ]);
    });

    it("should accept correct state relay (4 signatures)", async () => {
      await this.bridge.relayOracleState(
        "3210", // _blockHeight,
        "0xac0b4d71daabb289e21514b004628c7236fb351cd950f67c93c449c5d06b35d1", // _oracleIAVLStateHash
        "0xcf27679618b88c2adad50ed008d40bf5fe631a7b3718e766e547d203253bec6a", // _otherStoresMerkleHash
        "0xcf27679618b88c2adad50ed008d40bf5fe631a7b3718e766e547d203253bec6a", // _supplyStoresMerkleHash
        [
          "0x32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e", // subtreeVersionAndChainIdHash
          "0x3dddc70f2fb58681b19a6e7bb5d2361ef0a50370940d249d6290da9f086b99e8", // timeHash
          "0x942fc7a26da18816b277812839ac4f413654ac70fbfb64f27baab7f88fa6ec25", // txCountAndLastBlockInfoHash
          "0xdef482cda986470c27374601ec716e9853de47d72828ae0131cf8ef98e2972c5", // consensusDataHash
          "0x6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d", // lastResultsHash
          "0x7f4be7e5a1eb872ad44103360ddc190410331280c42a54d829a5d752c796685d" // evidenceAndProposerHash
        ],
        "0x6e0802118a0c00000000000022480a20", // _signedDataPrefix
        [
          [
            "0x1581824b505ee5977c87d3d1a5279544bb68d29e06aa541c068054383d7acffd", // r
            "0x34251b479113d5359e2bcc49055336bd5dfc1c02a0865c822414ff96c5cd3b19", // s
            28, // v
            "0x12240a20e8a6f45e19dd488df492c7b82f476d3691e46ff176ad4a13ff03a61e60eddfff10012a0c08abcdf2ef0510c8cb9c8a01320962616e64636861696e" // _signedDataSuffix
          ],
          [
            "0xbcf871c64d1f92d92c712535c3b15e3d8ed1dfe79c28090360d687a9e019952b", // r
            "0x5032abeb19ad210a93bff410efd2944762487dc49a3ef1682ac59007cd7b21dd", // s
            28, // v
            "0x12240a20e8a6f45e19dd488df492c7b82f476d3691e46ff176ad4a13ff03a61e60eddfff10012a0c08abcdf2ef0510c8afb68a01320962616e64636861696e" // _signedDataSuffix
          ],
          [
            "0xa50af0e18b03076d6f070d9b9321278c668a8d2e247fc54abc1dc8c777b98369", // r
            "0x25ea236f9d0dd2270fe8d9b6b6499de05e7e4a6b12b8a0e6849653cec64c3e61", // s
            28, // v
            "0x12240a20e8a6f45e19dd488df492c7b82f476d3691e46ff176ad4a13ff03a61e60eddfff10012a0c08abcdf2ef051080fafa8f01320962616e64636861696e" // _signedDataSuffix
          ],
          [
            "0xddeea6cf5643a100e8004590296c56ebd1231c774313699b8bc774ea86b59b96", // r
            "0x09779c784dc684d08f328bde36b91bd827dc86bfc47b68c45121089b3c1f6d89", // s
            27, // v
            "0x12240a20e8a6f45e19dd488df492c7b82f476d3691e46ff176ad4a13ff03a61e60eddfff10012a0c08abcdf2ef0510a0dfd68e01320962616e64636861696e" // _signedDataSuffix
          ]
        ]
      );
    });

    it("should accept correct state relay (3 signatures)", async () => {
      await this.bridge.relayOracleState(
        "3210", // _blockHeight,
        "0xac0b4d71daabb289e21514b004628c7236fb351cd950f67c93c449c5d06b35d1", // _oracleIAVLStateHash
        "0xcf27679618b88c2adad50ed008d40bf5fe631a7b3718e766e547d203253bec6a", // _otherStoresMerkleHash
        "0xcf27679618b88c2adad50ed008d40bf5fe631a7b3718e766e547d203253bec6a", // _supplyStoresMerkleHash
        [
          "0x32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e", // subtreeVersionAndChainIdHash
          "0x3dddc70f2fb58681b19a6e7bb5d2361ef0a50370940d249d6290da9f086b99e8", // timeHash
          "0x942fc7a26da18816b277812839ac4f413654ac70fbfb64f27baab7f88fa6ec25", // txCountAndLastBlockInfoHash
          "0xdef482cda986470c27374601ec716e9853de47d72828ae0131cf8ef98e2972c5", // consensusDataHash
          "0x6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d", // lastResultsHash
          "0x7f4be7e5a1eb872ad44103360ddc190410331280c42a54d829a5d752c796685d" // evidenceAndProposerHash
        ],
        "0x6e0802118a0c00000000000022480a20", // _signedDataPrefix
        [
          [
            "0xbcf871c64d1f92d92c712535c3b15e3d8ed1dfe79c28090360d687a9e019952b", // r
            "0x5032abeb19ad210a93bff410efd2944762487dc49a3ef1682ac59007cd7b21dd", // s
            28, // v
            "0x12240a20e8a6f45e19dd488df492c7b82f476d3691e46ff176ad4a13ff03a61e60eddfff10012a0c08abcdf2ef0510c8afb68a01320962616e64636861696e" // _signedDataSuffix
          ],
          [
            "0xa50af0e18b03076d6f070d9b9321278c668a8d2e247fc54abc1dc8c777b98369", // r
            "0x25ea236f9d0dd2270fe8d9b6b6499de05e7e4a6b12b8a0e6849653cec64c3e61", // s
            28, // v
            "0x12240a20e8a6f45e19dd488df492c7b82f476d3691e46ff176ad4a13ff03a61e60eddfff10012a0c08abcdf2ef051080fafa8f01320962616e64636861696e" // _signedDataSuffix
          ],
          [
            "0xddeea6cf5643a100e8004590296c56ebd1231c774313699b8bc774ea86b59b96", // r
            "0x09779c784dc684d08f328bde36b91bd827dc86bfc47b68c45121089b3c1f6d89", // s
            27, // v
            "0x12240a20e8a6f45e19dd488df492c7b82f476d3691e46ff176ad4a13ff03a61e60eddfff10012a0c08abcdf2ef0510a0dfd68e01320962616e64636861696e" // _signedDataSuffix
          ]
        ]
      );
    });

    it("should not accept out-of-order signatures", async () => {
      await expectRevert(
        this.bridge.relayOracleState(
          "3210", // _blockHeight,
          "0xac0b4d71daabb289e21514b004628c7236fb351cd950f67c93c449c5d06b35d1", // _oracleIAVLStateHash
          "0xcf27679618b88c2adad50ed008d40bf5fe631a7b3718e766e547d203253bec6a", // _otherStoresMerkleHash
          "0xcf27679618b88c2adad50ed008d40bf5fe631a7b3718e766e547d203253bec6a", // _supplyStoresMerkleHash
          [
            "0x32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e", // subtreeVersionAndChainIdHash
            "0x3dddc70f2fb58681b19a6e7bb5d2361ef0a50370940d249d6290da9f086b99e8", // timeHash
            "0x942fc7a26da18816b277812839ac4f413654ac70fbfb64f27baab7f88fa6ec25", // txCountAndLastBlockInfoHash
            "0xdef482cda986470c27374601ec716e9853de47d72828ae0131cf8ef98e2972c5", // consensusDataHash
            "0x6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d", // lastResultsHash
            "0x7f4be7e5a1eb872ad44103360ddc190410331280c42a54d829a5d752c796685d" // evidenceAndProposerHash
          ],
          "0x6e0802118a0c00000000000022480a20", // _signedDataPrefix
          [
            [
              "0xbcf871c64d1f92d92c712535c3b15e3d8ed1dfe79c28090360d687a9e019952b", // r
              "0x5032abeb19ad210a93bff410efd2944762487dc49a3ef1682ac59007cd7b21dd", // s
              28, // v
              "0x12240a20e8a6f45e19dd488df492c7b82f476d3691e46ff176ad4a13ff03a61e60eddfff10012a0c08abcdf2ef0510c8afb68a01320962616e64636861696e" // _signedDataSuffix
            ],
            [
              "0x1581824b505ee5977c87d3d1a5279544bb68d29e06aa541c068054383d7acffd", // r
              "0x34251b479113d5359e2bcc49055336bd5dfc1c02a0865c822414ff96c5cd3b19", // s
              28, // v
              "0x12240a20e8a6f45e19dd488df492c7b82f476d3691e46ff176ad4a13ff03a61e60eddfff10012a0c08abcdf2ef0510c8cb9c8a01320962616e64636861696e" // _signedDataSuffix
            ],
            [
              "0xa50af0e18b03076d6f070d9b9321278c668a8d2e247fc54abc1dc8c777b98369", // r
              "0x25ea236f9d0dd2270fe8d9b6b6499de05e7e4a6b12b8a0e6849653cec64c3e61", // s
              28, // v
              "0x12240a20e8a6f45e19dd488df492c7b82f476d3691e46ff176ad4a13ff03a61e60eddfff10012a0c08abcdf2ef051080fafa8f01320962616e64636861696e" // _signedDataSuffix
            ],
            [
              "0xddeea6cf5643a100e8004590296c56ebd1231c774313699b8bc774ea86b59b96", // r
              "0x09779c784dc684d08f328bde36b91bd827dc86bfc47b68c45121089b3c1f6d89", // s
              27, // v
              "0x12240a20e8a6f45e19dd488df492c7b82f476d3691e46ff176ad4a13ff03a61e60eddfff10012a0c08abcdf2ef0510a0dfd68e01320962616e64636861696e" // _signedDataSuffix
            ]
          ]
        ),
        "INVALID_SIGNATURE_SIGNER_ORDER"
      );
    });

    it("should not accept invalid signature", async () => {
      await expectRevert(
        this.bridge.relayOracleState(
          "3210", // _blockHeight,
          "0xac0b4d71daabb289e21514b004628c7236fb351cd950f67c93c449c5d06b35d1", // _oracleIAVLStateHash
          "0xcf27679618b88c2adad50ed008d40bf5fe631a7b3718e766e547d203253bec6a", // _otherStoresMerkleHash
          "0xcf27679618b88c2adad50ed008d40bf5fe631a7b3718e766e547d203253bec6a", // _supplyStoresMerkleHash
          [
            "0x32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e", // subtreeVersionAndChainIdHash
            "0x3dddc70f2fb58681b19a6e7bb5d2361ef0a50370940d249d6290da9f086b99e8", // timeHash
            "0x942fc7a26da18816b277812839ac4f413654ac70fbfb64f27baab7f88fa6ec25", // txCountAndLastBlockInfoHash
            "0xdef482cda986470c27374601ec716e9853de47d72828ae0131cf8ef98e2972c5", // consensusDataHash
            "0x6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d", // lastResultsHash
            "0x7f4be7e5a1eb872ad44103360ddc190410331280c42a54d829a5d752c796685d" // evidenceAndProposerHash
          ],
          "0x6e0802118a0c00000000000022480a20", // _signedDataPrefix
          [
            [
              "0x1581824b505ee5977c87d3d1a5279544bb68d29e06aa541c068054383d7acffe", // r INVALID HERE
              "0x34251b479113d5359e2bcc49055336bd5dfc1c02a0865c822414ff96c5cd3b19", // s
              28, // v
              "0x12240a20e8a6f45e19dd488df492c7b82f476d3691e46ff176ad4a13ff03a61e60eddfff10012a0c08abcdf2ef0510c8cb9c8a01320962616e64636861696e" // _signedDataSuffix
            ],
            [
              "0xa50af0e18b03076d6f070d9b9321278c668a8d2e247fc54abc1dc8c777b98369", // r
              "0x25ea236f9d0dd2270fe8d9b6b6499de05e7e4a6b12b8a0e6849653cec64c3e61", // s
              28, // v
              "0x12240a20e8a6f45e19dd488df492c7b82f476d3691e46ff176ad4a13ff03a61e60eddfff10012a0c08abcdf2ef051080fafa8f01320962616e64636861696e" // _signedDataSuffix
            ],
            [
              "0xddeea6cf5643a100e8004590296c56ebd1231c774313699b8bc774ea86b59b96", // r
              "0x09779c784dc684d08f328bde36b91bd827dc86bfc47b68c45121089b3c1f6d89", // s
              27, // v
              "0x12240a20e8a6f45e19dd488df492c7b82f476d3691e46ff176ad4a13ff03a61e60eddfff10012a0c08abcdf2ef0510a0dfd68e01320962616e64636861696e" // _signedDataSuffix
            ]
          ]
        ),
        "INSUFFICIENT_VALIDATOR_SIGNATURES"
      );
    });
  });

  context("Checking data verification", () => {
    beforeEach(async () => {
      this.bridge = await Bridge.new([]);
      await this.bridge.setOracleState(
        "55", // _blockHeight
        "0x7148d7db351b9b4624449801053d45fcf6a90edd64deeb3159ffe813c182f013" // _oracleIAVLStateHash
      );
    });

    it("should not accept unrelayed block", async () => {
      await expectRevert(
        this.bridge.verifyOracleData(
          "49", // _blockHeight
          "0x08011220bd86d5649a5b9218d5b96e009463dc91f7cf9d974f6227eb3a5b6d684db70361180f220800000000000add88", // _data
          "1", // _requestId
          "0x25ea236f9d0dd2270fe8d9b6b6499de05e7e4a6b12b8a0e6849653cec64c3e61", // _codeHash
          "0x", // _params
          "15", // _version
          [
            [
              false, // isDataOnRight
              "1", // subtreeHeight
              "2", // subtreeSize
              "15", // subtreeVersion
              "0xc3f1960b397caec50db04da19e1d674ccf9dd4b2af5dac010de6b143d5368693" // siblingHash
            ],
            [
              true, // isDataOnRight
              "2", // subtreeHeight
              "4", // subtreeSize
              "17", // subtreeVersion
              "0xc07cba57d6a2df08444c68500177a57837a8be3f9dd59e0885d7cf497ecf3c99" // siblingHash
            ],
            [
              false, // isDataOnRight
              "3", // subtreeHeight
              "8", // subtreeSize
              "17", // subtreeVersion
              "0xda058f8d239a1261661552f5bc0801acfca8a90d08c834df845d6ea1e34e2cfc" // siblingHash
            ]
          ]
        ),
        "NO_ORACLE_ROOT_STATE_DATA"
      );
    });

    it("should accept correct data verification", async () => {
      (
        await this.bridge.verifyOracleData(
          "55", // _blockHeight
          "0x08011220bd86d5649a5b9218d5b96e009463dc91f7cf9d974f6227eb3a5b6d684db70361180f220800000000000add88", // _data
          "1", // _requestId
          "0x25ea236f9d0dd2270fe8d9b6b6499de05e7e4a6b12b8a0e6849653cec64c3e61", // _codeHash
          "0x", // _params
          "15", // _version
          [
            [
              false, // isDataOnRight
              "1", // subtreeHeight
              "2", // subtreeSize
              "15", // subtreeVersion
              "0xc3f1960b397caec50db04da19e1d674ccf9dd4b2af5dac010de6b143d5368693" // siblingHash
            ],
            [
              true, // isDataOnRight
              "2", // subtreeHeight
              "4", // subtreeSize
              "17", // subtreeVersion
              "0xc07cba57d6a2df08444c68500177a57837a8be3f9dd59e0885d7cf497ecf3c99" // siblingHash
            ],
            [
              false, // isDataOnRight
              "3", // subtreeHeight
              "8", // subtreeSize
              "17", // subtreeVersion
              "0xda058f8d239a1261661552f5bc0801acfca8a90d08c834df845d6ea1e34e2cfc" // siblingHash
            ]
          ]
        )
      ).should.eq(true);
    });

    it("should not accept invalid data verification", async () => {
      await expectRevert(
        this.bridge.verifyOracleData(
          "55", // _blockHeight
          "0x08011220bd86d5649a5b9218d5b96e009463dc91f7cf9d974f6227eb3a5b6d684db70361180f220800000000000add89", // _data INVALID HERE
          "1", // _requestId
          "0x25ea236f9d0dd2270fe8d9b6b6499de05e7e4a6b12b8a0e6849653cec64c3e61", // _codeHash
          "0x", // _params
          "15", // _version
          [
            [
              false, // isDataOnRight
              "1", // subtreeHeight
              "2", // subtreeSize
              "15", // subtreeVersion
              "0xc3f1960b397caec50db04da19e1d674ccf9dd4b2af5dac010de6b143d5368693" // siblingHash
            ],
            [
              true, // isDataOnRight
              "2", // subtreeHeight
              "4", // subtreeSize
              "17", // subtreeVersion
              "0xc07cba57d6a2df08444c68500177a57837a8be3f9dd59e0885d7cf497ecf3c99" // siblingHash
            ],
            [
              false, // isDataOnRight
              "3", // subtreeHeight
              "8", // subtreeSize
              "17", // subtreeVersion
              "0xda058f8d239a1261661552f5bc0801acfca8a90d08c834df845d6ea1e34e2cfc" // siblingHash
            ]
          ]
        ),
        "INVALID_ORACLE_DATA_PROOF"
      );
    });

    it("should not accept incomplete proof", async () => {
      await expectRevert(
        this.bridge.verifyOracleData(
          "55", // _blockHeight
          "0x08011220bd86d5649a5b9218d5b96e009463dc91f7cf9d974f6227eb3a5b6d684db70361180f220800000000000add88", // _data
          "1", // _requestId
          "0x25ea236f9d0dd2270fe8d9b6b6499de05e7e4a6b12b8a0e6849653cec64c3e61", // _codeHash
          "0x", // _params
          "15", // _version
          [
            [
              false, // isDataOnRight
              "1", // subtreeHeight
              "2", // subtreeSize
              "15", // subtreeVersion
              "0xc3f1960b397caec50db04da19e1d674ccf9dd4b2af5dac010de6b143d5368693" // siblingHash
            ],
            [
              true, // isDataOnRight
              "2", // subtreeHeight
              "4", // subtreeSize
              "17", // subtreeVersion
              "0xc07cba57d6a2df08444c68500177a57837a8be3f9dd59e0885d7cf497ecf3c99" // siblingHash
            ]
          ]
        ),
        "INVALID_ORACLE_DATA_PROOF"
      );
    });
  });
});
