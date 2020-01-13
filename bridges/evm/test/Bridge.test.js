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
        "1160", // _blockHeight,
        "0xACAB016EC9FB3AA28A6A4BE8A364AEDAA9A42866E2957C5C267E340CE67C55EE", // _oracleIAVLStateHash
        "0x6ABB9CA0E0AC77A3B7C7F94D56E181DE954B92A19389829CE0E5A95B74BE0B7D", // _otherStoresMerkleHash
        "0x2726178BFFB0D462C15AB546DE7B4CA86588A98FF0F629DB7CA7E318AA61A846", // _supplyStoresMerkleHash
        [
          "0x32FA694879095840619F5E49380612BD296FF7E950EAFB66FF654D99CA70869E", // subtreeVersionAndChainIdHash
          "0x2157913D927F4249C52FAB326E9E0E83FACFAF167FA038A88173FA42ADF2452C", // timeHash
          "0x7EECC6A0EE0136DE143C92370E4BE8FA6F545C02C23DAFA62CC4AA0A14701787", // txCountAndLastBlockInfoHash
          "0xDEF482CDA986470C27374601EC716E9853DE47D72828AE0131CF8EF98E2972C5", // consensusDataHash
          "0x6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D", // lastResultsHash
          "0x0EFE3E12F46363C7779140D4CE659925DB52F19053E114D7CC4EFD666B37F79F" // evidenceAndProposerHash
        ],
        "0x6E080211880400000000000022480A20", // _signedDataPrefix
        [
          [
            "0xB88E0A2054A96A6775A9F5D1FA23B6FFA41274DD35C6431DAB0977F8CE4FB480", // r
            "0x3D759EFF85E17601624D560A8ACD70E782EA23B58C2E718FAC98EBF488750A86", // s
            28, // v
            "0x12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E" // _signedDataSuffix
          ],
          [
            "0x17A66FF70C81C6A9C3040C1037CCC4EE9319E184D40956DC0DC30C1318901D36", // r
            "0x4A4C0C9BF150967CE25C724E613DF6BE0C401B84AA29DE8599963F52A7DFA940", // s
            27, // v
            "0x12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E" // _signedDataSuffix
          ],
          [
            "0x453498042685AB34C627B5652E2F1FAD839C21DB3CEE4E01822F00885F1E0321", // r
            "0x679781C8F2E3597ED3DF15E8E44B9CAF17D74894F0D4E22CD0F8C7CC1CB43963", // s
            27, // v
            "0x12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E" // _signedDataSuffix
          ],
          [
            "0x174505557E61260C06C7FD8962FF485BEBAD68E91B00C225452962B1FCBF1114", // r
            "0x39B37ACD1759D47B09D18C6C3144EAB5B2D2CA34347DC60A4D58B369730C0DB9", // s
            27, // v
            "0x12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E" // _signedDataSuffix
          ]
        ]
      );
    });

    it("should accept correct state relay (3 signatures)", async () => {
      await this.bridge.relayOracleState(
        "1160", // _blockHeight,
        "0xACAB016EC9FB3AA28A6A4BE8A364AEDAA9A42866E2957C5C267E340CE67C55EE", // _oracleIAVLStateHash
        "0x6ABB9CA0E0AC77A3B7C7F94D56E181DE954B92A19389829CE0E5A95B74BE0B7D", // _otherStoresMerkleHash
        "0x2726178BFFB0D462C15AB546DE7B4CA86588A98FF0F629DB7CA7E318AA61A846", // _supplyStoresMerkleHash
        [
          "0x32FA694879095840619F5E49380612BD296FF7E950EAFB66FF654D99CA70869E", // subtreeVersionAndChainIdHash
          "0x2157913D927F4249C52FAB326E9E0E83FACFAF167FA038A88173FA42ADF2452C", // timeHash
          "0x7EECC6A0EE0136DE143C92370E4BE8FA6F545C02C23DAFA62CC4AA0A14701787", // txCountAndLastBlockInfoHash
          "0xDEF482CDA986470C27374601EC716E9853DE47D72828AE0131CF8EF98E2972C5", // consensusDataHash
          "0x6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D", // lastResultsHash
          "0x0EFE3E12F46363C7779140D4CE659925DB52F19053E114D7CC4EFD666B37F79F" // evidenceAndProposerHash
        ],
        "0x6E080211880400000000000022480A20", // _signedDataPrefix
        [
          [
            "0xB88E0A2054A96A6775A9F5D1FA23B6FFA41274DD35C6431DAB0977F8CE4FB480", // r
            "0x3D759EFF85E17601624D560A8ACD70E782EA23B58C2E718FAC98EBF488750A86", // s
            28, // v
            "0x12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E" // _signedDataSuffix
          ],
          [
            "0x17A66FF70C81C6A9C3040C1037CCC4EE9319E184D40956DC0DC30C1318901D36", // r
            "0x4A4C0C9BF150967CE25C724E613DF6BE0C401B84AA29DE8599963F52A7DFA940", // s
            27, // v
            "0x12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E" // _signedDataSuffix
          ],
          [
            "0x174505557E61260C06C7FD8962FF485BEBAD68E91B00C225452962B1FCBF1114", // r
            "0x39B37ACD1759D47B09D18C6C3144EAB5B2D2CA34347DC60A4D58B369730C0DB9", // s
            27, // v
            "0x12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E" // _signedDataSuffix
          ]
        ]
      );
    });

    it("should not accept out-of-order signatures", async () => {
      await expectRevert(
        this.bridge.relayOracleState(
          "1160", // _blockHeight,
          "0xACAB016EC9FB3AA28A6A4BE8A364AEDAA9A42866E2957C5C267E340CE67C55EE", // _oracleIAVLStateHash
          "0x6ABB9CA0E0AC77A3B7C7F94D56E181DE954B92A19389829CE0E5A95B74BE0B7D", // _otherStoresMerkleHash
          "0x2726178BFFB0D462C15AB546DE7B4CA86588A98FF0F629DB7CA7E318AA61A846", // _supplyStoresMerkleHash
          [
            "0x32FA694879095840619F5E49380612BD296FF7E950EAFB66FF654D99CA70869E", // subtreeVersionAndChainIdHash
            "0x2157913D927F4249C52FAB326E9E0E83FACFAF167FA038A88173FA42ADF2452C", // timeHash
            "0x7EECC6A0EE0136DE143C92370E4BE8FA6F545C02C23DAFA62CC4AA0A14701787", // txCountAndLastBlockInfoHash
            "0xDEF482CDA986470C27374601EC716E9853DE47D72828AE0131CF8EF98E2972C5", // consensusDataHash
            "0x6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D", // lastResultsHash
            "0x0EFE3E12F46363C7779140D4CE659925DB52F19053E114D7CC4EFD666B37F79F" // evidenceAndProposerHash
          ],
          "0x6E080211880400000000000022480A20", // _signedDataPrefix
          [
            [
              "0xB88E0A2054A96A6775A9F5D1FA23B6FFA41274DD35C6431DAB0977F8CE4FB480", // r
              "0x3D759EFF85E17601624D560A8ACD70E782EA23B58C2E718FAC98EBF488750A86", // s
              28, // v
              "0x12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E" // _signedDataSuffix
            ],
            [
              "0x174505557E61260C06C7FD8962FF485BEBAD68E91B00C225452962B1FCBF1114", // r
              "0x39B37ACD1759D47B09D18C6C3144EAB5B2D2CA34347DC60A4D58B369730C0DB9", // s
              27, // v
              "0x12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E" // _signedDataSuffix
            ],
            [
              "0x17A66FF70C81C6A9C3040C1037CCC4EE9319E184D40956DC0DC30C1318901D36", // r
              "0x4A4C0C9BF150967CE25C724E613DF6BE0C401B84AA29DE8599963F52A7DFA940", // s
              27, // v
              "0x12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E" // _signedDataSuffix
            ],
            [
              "0x453498042685AB34C627B5652E2F1FAD839C21DB3CEE4E01822F00885F1E0321", // r
              "0x679781C8F2E3597ED3DF15E8E44B9CAF17D74894F0D4E22CD0F8C7CC1CB43963", // s
              27, // v
              "0x12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E" // _signedDataSuffix
            ]
          ]
        ),
        "INVALID_SIGNATURE_SIGNER_ORDER"
      );
    });

    it("should not accept invalid signature", async () => {
      await expectRevert(
        this.bridge.relayOracleState(
          "1160", // _blockHeight,
          "0xACAB016EC9FB3AA28A6A4BE8A364AEDAA9A42866E2957C5C267E340CE67C55EE", // _oracleIAVLStateHash
          "0x6ABB9CA0E0AC77A3B7C7F94D56E181DE954B92A19389829CE0E5A95B74BE0B7D", // _otherStoresMerkleHash
          "0x2726178BFFB0D462C15AB546DE7B4CA86588A98FF0F629DB7CA7E318AA61A846", // _supplyStoresMerkleHash
          [
            "0x32FA694879095840619F5E49380612BD296FF7E950EAFB66FF654D99CA70869E", // subtreeVersionAndChainIdHash
            "0x2157913D927F4249C52FAB326E9E0E83FACFAF167FA038A88173FA42ADF2452C", // timeHash
            "0x7EECC6A0EE0136DE143C92370E4BE8FA6F545C02C23DAFA62CC4AA0A14701787", // txCountAndLastBlockInfoHash
            "0xDEF482CDA986470C27374601EC716E9853DE47D72828AE0131CF8EF98E2972C5", // consensusDataHash
            "0x6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D", // lastResultsHash
            "0x0EFE3E12F46363C7779140D4CE659925DB52F19053E114D7CC4EFD666B37F79F" // evidenceAndProposerHash
          ],
          "0x6E080211880400000000000022480A20", // _signedDataPrefix
          [
            [
              "0xE88E0A2054A96A6775A9F5D1FA23B6FFA41274DD35C6431DAB0977F8CE4FB480", // r INVALID HERE
              "0x3D759EFF85E17601624D560A8ACD70E782EA23B58C2E718FAC98EBF488750A86", // s
              28, // v
              "0x12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E" // _signedDataSuffix
            ],
            [
              "0x17A66FF70C81C6A9C3040C1037CCC4EE9319E184D40956DC0DC30C1318901D36", // r
              "0x4A4C0C9BF150967CE25C724E613DF6BE0C401B84AA29DE8599963F52A7DFA940", // s
              27, // v
              "0x12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E" // _signedDataSuffix
            ],
            [
              "0x453498042685AB34C627B5652E2F1FAD839C21DB3CEE4E01822F00885F1E0321", // r
              "0x679781C8F2E3597ED3DF15E8E44B9CAF17D74894F0D4E22CD0F8C7CC1CB43963", // s
              27, // v
              "0x12240A20B0E9D07640EE2E758D01EA69E0733276D90946B0E7D11FE86067F97BAB1CC11B10012A0C08CBBAB8F00510B8C48F8902320962616E64636861696E" // _signedDataSuffix
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
        "1456", // _blockHeight
        "0x93C06B3C6357C6CD299552B4CD683B95F99695F1285C877D2CC481CE92FF0ACA" // _oracleIAVLStateHash
      );
    });

    it("should not accept unrelayed block", async () => {
      await expectRevert(
        this.bridge.verifyOracleData(
          "1450", // _blockHeight
          "0x00000000000AEAD8", // _data
          "1", // _requestId
          "0x0874EE3E5ABA7AE0EB8CB43BFEBED358826D111CECE0EF0F804E99EEA9264060", // _codeHash
          "0x706172616D73", // _params
          "51", // _version
          [
            [
              true, // isDataOnRight
              "1", // subtreeHeight
              "2", // subtreeSize
              "51", // subtreeVersion
              "0x81C8ADD6C13611F5416CA722C98815846504D7D2A77777B15D4D4212902F7785" // siblingHash
            ],
            [
              true, // isDataOnRight
              "2", // subtreeHeight
              "3", // subtreeSize
              "51", // subtreeVersion
              "0x29BB2C7201C2ECC2B0D1524D7E51827BC6C79E1E6E53249A630485A19876A173" // siblingHash
            ],
            [
              true, // isDataOnRight
              "3", // subtreeHeight
              "5", // subtreeSize
              "51", // subtreeVersion
              "0x8ED79B88CF56E86D5FA467274468ADE8FCA0BCCD5C797DE745AEB3D04F6B94DD" // siblingHash
            ],
            [
              true, // isDataOnRight
              "4", // subtreeHeight
              "9", // subtreeSize
              "1455", // subtreeVersion
              "0x3E3734FD33F6C1DC10C2A0979FF3E36E31D1039331B19041BCE34BD4155E6757" // siblingHash
            ]
          ]
        ),
        "NO_ORACLE_ROOT_STATE_DATA"
      );
    });

    it("should accept correct data verification", async () => {
      (
        await this.bridge.verifyOracleData(
          "1456", // _blockHeight
          "0x00000000000AEAD8", // _data
          "1", // _requestId
          "0x0874EE3E5ABA7AE0EB8CB43BFEBED358826D111CECE0EF0F804E99EEA9264060", // _codeHash
          "0x706172616D73", // _params
          "51", // _version
          [
            [
              true, // isDataOnRight
              "1", // subtreeHeight
              "2", // subtreeSize
              "51", // subtreeVersion
              "0x81C8ADD6C13611F5416CA722C98815846504D7D2A77777B15D4D4212902F7785" // siblingHash
            ],
            [
              true, // isDataOnRight
              "2", // subtreeHeight
              "3", // subtreeSize
              "51", // subtreeVersion
              "0x29BB2C7201C2ECC2B0D1524D7E51827BC6C79E1E6E53249A630485A19876A173" // siblingHash
            ],
            [
              true, // isDataOnRight
              "3", // subtreeHeight
              "5", // subtreeSize
              "51", // subtreeVersion
              "0x8ED79B88CF56E86D5FA467274468ADE8FCA0BCCD5C797DE745AEB3D04F6B94DD" // siblingHash
            ],
            [
              true, // isDataOnRight
              "4", // subtreeHeight
              "9", // subtreeSize
              "1455", // subtreeVersion
              "0x3E3734FD33F6C1DC10C2A0979FF3E36E31D1039331B19041BCE34BD4155E6757" // siblingHash
            ]
          ]
        )
      )
        .toString()
        .should.eq(
          [
            "0x00000000000aead8",
            "0x0874ee3e5aba7ae0eb8cb43bfebed358826d111cece0ef0f804e99eea9264060",
            "0x706172616d73"
          ].toString()
        );
    });

    it("should not accept invalid data verification", async () => {
      await expectRevert(
        this.bridge.verifyOracleData(
          "1456", // _blockHeight
          "0x00000000000AEAD9", // _data WRONG HERE
          "1", // _requestId
          "0x0874EE3E5ABA7AE0EB8CB43BFEBED358826D111CECE0EF0F804E99EEA9264060", // _codeHash
          "0x706172616D73", // _params
          "51", // _version
          [
            [
              true, // isDataOnRight
              "1", // subtreeHeight
              "2", // subtreeSize
              "51", // subtreeVersion
              "0x81C8ADD6C13611F5416CA722C98815846504D7D2A77777B15D4D4212902F7785" // siblingHash
            ],
            [
              true, // isDataOnRight
              "2", // subtreeHeight
              "3", // subtreeSize
              "51", // subtreeVersion
              "0x29BB2C7201C2ECC2B0D1524D7E51827BC6C79E1E6E53249A630485A19876A173" // siblingHash
            ],
            [
              true, // isDataOnRight
              "3", // subtreeHeight
              "5", // subtreeSize
              "51", // subtreeVersion
              "0x8ED79B88CF56E86D5FA467274468ADE8FCA0BCCD5C797DE745AEB3D04F6B94DD" // siblingHash
            ],
            [
              true, // isDataOnRight
              "4", // subtreeHeight
              "9", // subtreeSize
              "1455", // subtreeVersion
              "0x3E3734FD33F6C1DC10C2A0979FF3E36E31D1039331B19041BCE34BD4155E6757" // siblingHash
            ]
          ]
        ),
        "INVALID_ORACLE_DATA_PROOF"
      );
    });

    it("should not accept incomplete proof", async () => {
      await expectRevert(
        this.bridge.verifyOracleData(
          "1456", // _blockHeight
          "0x00000000000AEAD8", // _data
          "1", // _requestId
          "0x0874EE3E5ABA7AE0EB8CB43BFEBED358826D111CECE0EF0F804E99EEA9264060", // _codeHash
          "0x706172616D73", // _params
          "51", // _version
          [
            [
              true, // isDataOnRight
              "1", // subtreeHeight
              "2", // subtreeSize
              "51", // subtreeVersion
              "0x81C8ADD6C13611F5416CA722C98815846504D7D2A77777B15D4D4212902F7785" // siblingHash
            ],
            [
              true, // isDataOnRight
              "2", // subtreeHeight
              "3", // subtreeSize
              "51", // subtreeVersion
              "0x29BB2C7201C2ECC2B0D1524D7E51827BC6C79E1E6E53249A630485A19876A173" // siblingHash
            ],
            [
              true, // isDataOnRight
              "3", // subtreeHeight
              "5", // subtreeSize
              "51", // subtreeVersion
              "0x8ED79B88CF56E86D5FA467274468ADE8FCA0BCCD5C797DE745AEB3D04F6B94DD" // siblingHash
            ]
          ]
        ),
        "INVALID_ORACLE_DATA_PROOF"
      );
    });
  });
});
