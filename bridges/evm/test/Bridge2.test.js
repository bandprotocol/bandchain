const { expectRevert } = require("openzeppelin-test-helpers");
const OracleBridge = artifacts.require("OracleBridgeMock");

require("chai").should();

contract("OracleBridge", () => {
  context("Checking oracle state relay (2 validators)", () => {
    beforeEach(async () => {
      this.bridge = await OracleBridge.new([
        "0xaAA22E077492CbaD414098EBD98AA8dc1C7AE8D9",
        "0xB956589b6fC5523eeD0d9eEcfF06262Ce84ff260"
      ]);
    });

    it("should accept correct state relay", async () => {
      await this.bridge.relayOracleState(
        "55", // _blockHeight,
        "0x7148d7db351b9b4624449801053d45fcf6a90edd64deeb3159ffe813c182f013", // _oracleIAVLStateHash
        "0x406cfc22544f4c74049983f871d44a3cf2be94bbde3961cab7ed4773d2e57ee0", // _otherStoresMerkleHash
        [
          "0x32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e", // subtreeVersionAndChainIdHash
          "0xd82f0576c09d2dfe5783eea26b3f834c4ce4866b330670dec7b6b97d53ce9687", // timeHash
          "0xa468e310ffeda3113422e774f1fe7785b53e2bab9eaf231a0e85c3eda4338ecf", // txCountAndLastBlockInfoHash
          "0xdef482cda986470c27374601ec716e9853de47d72828ae0131cf8ef98e2972c5", // consensusDataHash
          "0x6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d", // lastResultsHash
          "0xd991da4d4e69473cc75a4b819f9e07d4956671a6f4a74df4cc16596fcbe68137" // evidenceAndProposerHash
        ],
        "0x6e080211370000000000000022480a20", // _signedDataPrefix
        [
          [
            "0xe57f9399c976b4aacc2b7ce923edfd1802da637a0aec5eae2e5cce272e893e2c", // r
            "0x7ec5209a1cc5cec0e57c1849cc2b62816d09c3e95860fd2f10e29dd1afe76d2e", // s
            "27", // v
            "0x12240a204369248f6ca1f8caa75acdb98560c7c9f015ab5c85283480984e901af6019b5310012a0c0891cce6ef051094f598d503320962616e64636861696e" // signedDataSuffix
          ],
          [
            "0xaea1120a7e539fb98d2fb94ce59cfb2ab84435bd3847dba6b0a846be73998179", // r
            "0x38fe77c8b788933e8ec9ee7bf3c6d41feba6f12ab7b9607ea2662c6c0ea10ef3", // s
            "27", // v
            "0x12240a204369248f6ca1f8caa75acdb98560c7c9f015ab5c85283480984e901af6019b5310012a0c0891cce6ef051084e3e5d603320962616e64636861696e" // signedDataSuffix
          ]
        ]
      );
      // Another relay should just pass if it contains the same _oracleIAVLStateHash
      await this.bridge.relayOracleState(
        "55", // _blockHeight,
        "0x7148d7db351b9b4624449801053d45fcf6a90edd64deeb3159ffe813c182f013", // _oracleIAVLStateHash
        "0x406cfc22544f4c74049983f871d44a3cf2be94bbde3961cab7ed4773d2e57ee0", // _otherStoresMerkleHash
        [
          "0x32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e", // subtreeVersionAndChainIdHash
          "0xd82f0576c09d2dfe5783eea26b3f834c4ce4866b330670dec7b6b97d53ce9687", // timeHash
          "0xa468e310ffeda3113422e774f1fe7785b53e2bab9eaf231a0e85c3eda4338ecf", // txCountAndLastBlockInfoHash
          "0xdef482cda986470c27374601ec716e9853de47d72828ae0131cf8ef98e2972c5", // consensusDataHash
          "0x6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d", // lastResultsHash
          "0xd991da4d4e69473cc75a4b819f9e07d4956671a6f4a74df4cc16596fcbe68137" // evidenceAndProposerHash
        ],
        "0x6e080211370000000000000022480a20", // _signedDataPrefix
        []
      );
      // Another relay should fail if it contains a different _oracleIAVLStateHash
      await expectRevert(
        this.bridge.relayOracleState(
          "55", // _blockHeight,
          "0x7148d7db351b9b4624449801053d45fcf6a90edd64deeb3159ffe813c182f014", // _oracleIAVLStateHash INVALID HERE
          "0x406cfc22544f4c74049983f871d44a3cf2be94bbde3961cab7ed4773d2e57ee0", // _otherStoresMerkleHash
          [
            "0x32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e", // subtreeVersionAndChainIdHash
            "0xd82f0576c09d2dfe5783eea26b3f834c4ce4866b330670dec7b6b97d53ce9687", // timeHash
            "0xa468e310ffeda3113422e774f1fe7785b53e2bab9eaf231a0e85c3eda4338ecf", // txCountAndLastBlockInfoHash
            "0xdef482cda986470c27374601ec716e9853de47d72828ae0131cf8ef98e2972c5", // consensusDataHash
            "0x6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d", // lastResultsHash
            "0xd991da4d4e69473cc75a4b819f9e07d4956671a6f4a74df4cc16596fcbe68137" // evidenceAndProposerHash
          ],
          "0x6e080211370000000000000022480a20", // _signedDataPrefix
          []
        ),
        "INCONSISTENT_ORACLE_IAVL_STATE"
      );
    });

    it("should accept out-of-order signatures", async () => {
      await expectRevert(
        this.bridge.relayOracleState(
          "55", // _blockHeight,
          "0x7148d7db351b9b4624449801053d45fcf6a90edd64deeb3159ffe813c182f013", // _oracleIAVLStateHash
          "0x406cfc22544f4c74049983f871d44a3cf2be94bbde3961cab7ed4773d2e57ee0", // _otherStoresMerkleHash
          [
            "0x32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e", // subtreeVersionAndChainIdHash
            "0xd82f0576c09d2dfe5783eea26b3f834c4ce4866b330670dec7b6b97d53ce9687", // timeHash
            "0xa468e310ffeda3113422e774f1fe7785b53e2bab9eaf231a0e85c3eda4338ecf", // txCountAndLastBlockInfoHash
            "0xdef482cda986470c27374601ec716e9853de47d72828ae0131cf8ef98e2972c5", // consensusDataHash
            "0x6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d", // lastResultsHash
            "0xd991da4d4e69473cc75a4b819f9e07d4956671a6f4a74df4cc16596fcbe68137" // evidenceAndProposerHash
          ],
          "0x6e080211370000000000000022480a20", // _signedDataPrefix
          [
            [
              "0xaea1120a7e539fb98d2fb94ce59cfb2ab84435bd3847dba6b0a846be73998179", // r
              "0x38fe77c8b788933e8ec9ee7bf3c6d41feba6f12ab7b9607ea2662c6c0ea10ef3", // s
              "27", // v
              "0x12240a204369248f6ca1f8caa75acdb98560c7c9f015ab5c85283480984e901af6019b5310012a0c0891cce6ef051084e3e5d603320962616e64636861696e" // signedDataSuffix
            ],
            [
              "0xe57f9399c976b4aacc2b7ce923edfd1802da637a0aec5eae2e5cce272e893e2c", // r
              "0x7ec5209a1cc5cec0e57c1849cc2b62816d09c3e95860fd2f10e29dd1afe76d2e", // s
              "27", // v
              "0x12240a204369248f6ca1f8caa75acdb98560c7c9f015ab5c85283480984e901af6019b5310012a0c0891cce6ef051094f598d503320962616e64636861696e" // signedDataSuffix
            ]
          ]
        ),
        "INVALID_SIGNATURE_SIGNER_ORDER"
      );
    });

    it("should not accept invalid signature", async () => {
      await expectRevert(
        this.bridge.relayOracleState(
          "55", // _blockHeight,
          "0x7148d7db351b9b4624449801053d45fcf6a90edd64deeb3159ffe813c182f013", // _oracleIAVLStateHash
          "0x406cfc22544f4c74049983f871d44a3cf2be94bbde3961cab7ed4773d2e57ee0", // _otherStoresMerkleHash
          [
            "0x32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e", // subtreeVersionAndChainIdHash
            "0xd82f0576c09d2dfe5783eea26b3f834c4ce4866b330670dec7b6b97d53ce9687", // timeHash
            "0xa468e310ffeda3113422e774f1fe7785b53e2bab9eaf231a0e85c3eda4338ecf", // txCountAndLastBlockInfoHash
            "0xdef482cda986470c27374601ec716e9853de47d72828ae0131cf8ef98e2972c5", // consensusDataHash
            "0x6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d", // lastResultsHash
            "0xd991da4d4e69473cc75a4b819f9e07d4956671a6f4a74df4cc16596fcbe68137" // evidenceAndProposerHash
          ],
          "0x6e080211370000000000000022480a20", // _signedDataPrefix
          [
            [
              "0xaea1120a7e539fb98d2fb94ce59cfb2ab84435bd3847dba6b0a846be73998178", // r INVALID HERE!
              "0x38fe77c8b788933e8ec9ee7bf3c6d41feba6f12ab7b9607ea2662c6c0ea10ef3", // s
              "27", // v
              "0x12240a204369248f6ca1f8caa75acdb98560c7c9f015ab5c85283480984e901af6019b5310012a0c0891cce6ef051084e3e5d603320962616e64636861696e" // signedDataSuffix
            ],
            [
              "0xe57f9399c976b4aacc2b7ce923edfd1802da637a0aec5eae2e5cce272e893e2c", // r
              "0x7ec5209a1cc5cec0e57c1849cc2b62816d09c3e95860fd2f10e29dd1afe76d2e", // s
              "27", // v
              "0x12240a204369248f6ca1f8caa75acdb98560c7c9f015ab5c85283480984e901af6019b5310012a0c0891cce6ef051094f598d503320962616e64636861696e" // signedDataSuffix
            ]
          ]
        ),
        "INSUFFICIENT_VALIDATOR_SIGNATURES"
      );
    });
  });

  context("Checking data verification", () => {
    beforeEach(async () => {
      this.bridge = await OracleBridge.new([]);
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
