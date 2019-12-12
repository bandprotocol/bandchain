const BytesLibMock = artifacts.require("BytesLibMock");

const { expect } = require("chai");

describe("BytesLib", () => {
  beforeEach(async () => {
    this.bytesLib = await BytesLibMock.new();
  });

  describe("leafHash", () => {
    it("calculate leaf hash correctly", async () => {
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

  describe("innerHash", () => {
    it("calculate inner hash correctly", async () => {
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

  describe("decodeVarint", () => {
    it("decode 1 byte varint correctly", async () => {
      (await this.bytesLib.decodeVarint("0x74")).toString().should.eq("116");
    });
    it("decode 2 bytes varint correctly", async () => {
      (await this.bytesLib.decodeVarint("0xe374"))
        .toString()
        .should.eq("14947");
    });
    it("decode >2 bytes varint correctly", async () => {
      (await this.bytesLib.decodeVarint("0xa3f2e574"))
        .toString()
        .should.eq("244939043");
    });
  });

  describe("getBytes", () => {
    it("calculate (3) bytes from prefix", async () => {
      (
        await this.bytesLib.getBytes(
          "1356938545749799165119972480570561420155507632800475359837393562592742450709"
        )
      )
        .toString()
        .should.eq("0x9fa615");
    });

    it("calculate (3) bytes from prefix must padding zero", async () => {
      (
        await this.bytesLib.getBytes(
          "1356938545749799165119972480570561420155507632800475359837393562592731987993"
        )
      )
        .toString()
        .should.eq("0x000019");
    });

    it("calculate (3) bytes from prefix must not use exceed bytes", async () => {
      (
        await this.bytesLib.getBytes(
          "1356938545749799165119972480570561420155507632800475359837393811163738369617"
        )
      )
        .toString()
        .should.eq("0x896e51");
    });
  });
});
