const Utils = artifacts.require("UtilsMock");

require("chai").should();

contract("Utils", () => {
  beforeEach(async () => {
    this.utils = await Utils.new();
  });

  context("merkleLeafHash", () => {
    it("should calculate leaf hash correctly", async () => {
      (
        await this.utils.merkleLeafHash(
          "0x08d1082cc8d85a0833da8815ff1574675c415760e0aff7fb4e32de6de27faf86"
        )
      )
        .toString()
        .should.eq(
          "0x35b401b2a74452d2252df60574e0a6c029885965ae48f006ebddc18e53427e26"
        );
    });
  });

  context("merkleInnerHash", () => {
    it("should calculate inner hash correctly", async () => {
      (
        await this.utils.merkleInnerHash(
          "0x08d1082cc8d85a0833da8815ff1574675c415760e0aff7fb4e32de6de27faf86",
          "0x789411d15a12768a9c3eb99d3453d6ebb4481c2a03ab59cc262a97e25757afe6"
        )
      )
        .toString()
        .should.eq(
          "0xca48b611419f0848bf0fce9750caca6fd4fb8ef96ba8d7d3ccd4f05bf2af1661"
        );
    });
  });

  context("encodeVarintUnsigned", () => {
    it("should encode 1 byte varint correctly", async () => {
      (await this.utils.encodeVarintUnsigned("116"))
        .toString()
        .should.eq("0x74");
    });
    it("should encode 2 bytes varint correctly", async () => {
      (await this.utils.encodeVarintUnsigned("14947"))
        .toString()
        .should.eq("0xe374");
    });
    it("should encode >2 bytes varint correctly", async () => {
      (await this.utils.encodeVarintUnsigned("244939043"))
        .toString()
        .should.eq("0xa3f2e574");
    });
  });

  context("encodeVarintSigned", () => {
    it("should encode 1 byte varint correctly", async () => {
      (await this.utils.encodeVarintSigned("58")).toString().should.eq("0x74");
    });
    it("should encode 2 bytes varint correctly", async () => {
      (await this.utils.encodeVarintSigned("7473"))
        .toString()
        .should.eq("0xe274");
    });
    it("should encode >2 bytes varint correctly", async () => {
      (await this.utils.encodeVarintSigned("122469521"))
        .toString()
        .should.eq("0xa2f2e574");
    });
  });
});
