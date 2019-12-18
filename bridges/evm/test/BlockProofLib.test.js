const { BN, expectRevert } = require("openzeppelin-test-helpers");
const BlockProofLibMock = artifacts.require("BlockProofLibMock");

contract("BlockProofLib", () => {
  beforeEach(async () => {
    this.blockProofLib = await BlockProofLibMock.new();
    this.validators = [
      "0xaAA22E077492CbaD414098EBD98AA8dc1C7AE8D9",
      "0xB956589b6fC5523eeD0d9eEcfF06262Ce84ff260"
    ];
    this.data = {
      appHash:
        "0x6fec355d8a3ce024eed694d19f1e41ae2815c6c9609c5ef59d715935d2e76712",
      encodedHeight: "0x37",
      others: [
        "0x32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e",
        "0xd82f0576c09d2dfe5783eea26b3f834c4ce4866b330670dec7b6b97d53ce9687",
        "0xa468e310ffeda3113422e774f1fe7785b53e2bab9eaf231a0e85c3eda4338ecf",
        "0xdef482cda986470c27374601ec716e9853de47d72828ae0131cf8ef98e2972c5",
        "0x6e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d",
        "0xd991da4d4e69473cc75a4b819f9e07d4956671a6f4a74df4cc16596fcbe68137"
      ],
      leftMsg: "0x6e080211370000000000000022480a20",
      rightMsgSeperator:
        "904625697166532776746648320380374280103671755200316906558262375061821341503",
      rightMsg:
        "0x12240a204369248f6ca1f8caa75acdb98560c7c9f015ab5c85283480984e901af6019b5310012a0c0891cce6ef051094f598d503320962616e64636861696e12240a204369248f6ca1f8caa75acdb98560c7c9f015ab5c85283480984e901af6019b5310012a0c0891cce6ef051084e3e5d603320962616e64636861696e",
      signatures:
        "0xe57f9399c976b4aacc2b7ce923edfd1802da637a0aec5eae2e5cce272e893e2c7ec5209a1cc5cec0e57c1849cc2b62816d09c3e95860fd2f10e29dd1afe76d2e00aea1120a7e539fb98d2fb94ce59cfb2ab84435bd3847dba6b0a846be7399817938fe77c8b788933e8ec9ee7bf3c6d41feba6f12ab7b9607ea2662c6c0ea10ef300"
    };
  });

  context("getSignersFromSignatures", () => {
    it("should return correct validators set", async () => {
      (
        await this.blockProofLib.getSignersFromSignatures(
          this.data.appHash,
          this.data.encodedHeight,
          this.data.others,
          this.data.leftMsg,
          this.data.rightMsgSeperator,
          this.data.rightMsg,
          this.data.signatures
        )
      )
        .toString()
        .should.eq(this.validators.toString());
    });

    // it("should revert if there is no prefix", async () => {
    //   await expectRevert(
    //     this.blockProofLib.getLeafHash(
    //       [],
    //       this.data.paths,
    //       this.data.otherMSHashes,
    //       this.data.key,
    //       this.data.value
    //     ),
    //     "FIRST_PREFIX_IS_NEEDED"
    //   );
    // });
  });
});
