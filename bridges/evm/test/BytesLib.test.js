const BytesLibMock = artifacts.require("BytesLibMock");

contract("BytesLib", () => {
  beforeEach(async () => {
    this.bytesLib = await BytesLibMock.new();
  });

  context("leafHash", () => {
    it("should calculate leaf hash correctly", async () => {
      (
        await this.bytesLib.leafHash(
          "0x08d1082cc8d85a0833da8815ff1574675c415760e0aff7fb4e32de6de27faf86"
        )
      )
        .toString()
        .should.eq(
          "0x35b401b2a74452d2252df60574e0a6c029885965ae48f006ebddc18e53427e26"
        );
    });
  });

  context("innerHash", () => {
    it("should calculate inner hash correctly", async () => {
      (
        await this.bytesLib.innerHash(
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

  context("decodeVarint", () => {
    it("should decode 1 byte varint correctly", async () => {
      (await this.bytesLib.decodeVarint("0x74")).toString().should.eq("116");
    });
    it("should decode 2 bytes varint correctly", async () => {
      (await this.bytesLib.decodeVarint("0xe374"))
        .toString()
        .should.eq("14947");
    });
    it("should decode >2 bytes varint correctly", async () => {
      (await this.bytesLib.decodeVarint("0xa3f2e574"))
        .toString()
        .should.eq("244939043");
    });
  });

  context("getBytes", () => {
    it("should get 3 bytes from prefix", async () => {
      (
        await this.bytesLib.getBytes(
          "1356938545749799165119972480570561420155507632800475359837393562592742450709"
          // hex format is "0x3000000000000000000000000000000000000000000000000000000009fa615"
        )
      )
        .toString()
        .should.eq("0x9fa615");
    });

    it("should get 3 bytes from prefix with padding zero", async () => {
      (
        await this.bytesLib.getBytes(
          "1356938545749799165119972480570561420155507632800475359837393562592731987993"
          // hex format is "0x300000000000000000000000000000000000000000000000000000000000019"
        )
      )
        .toString()
        .should.eq("0x000019");
    });

    it("should get only last 3 bytes from prefix", async () => {
      (
        await this.bytesLib.getBytes(
          "1356938545749799165119972480570561420155507632800475359837393811163738369617"
          // hex format is "0x300000000000000000000000000000000000000000000000000e212f2896e51"
        )
      )
        .toString()
        .should.eq("0x896e51");
    });
  });
});
