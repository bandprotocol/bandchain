const { BN, expectRevert } = require("openzeppelin-test-helpers");
const StoreProofLibMock = artifacts.require("StoreProofLibMock");

contract("StoreProofLib", () => {
  beforeEach(async () => {
    this.storeProofLib = await StoreProofLibMock.new();
    this.leafHash =
      "0x27dec3dfda0ba44e9cc8b5bb76a9990073a6c1e19ecda62b970014e974efcc71";
    this.storeHash =
      "0x93645691212499386ff2dae262c62ad0bf0aa491c39839604ebda2e8a6313bc5";
    this.appHash =
      "0x51e3a28648c882a2ab40e3e025ed98d4d25cbf608ac3e7d122ac71b2837a50ca";
    this.data = {
      key: "0x0000000000000001",
      value:
        "0x08011220bd86d5649a5b9218d5b96e009463dc91f7cf9d974f6227eb3a5b6d684db70361180f220800000000000add88",
      prefixes: [
        "1356938545749799165119972480570561420155507632800475359837393562592731988510",
        "1356938545749799165119972480570561420155507632800475359837393562592732120094",
        "59252983164407896876905464984914515346790499965620757379566185566549297072162",
        "1356938545749799165119972480570561420155507632800475359837393562592732385314"
      ],
      paths: [
        "0xc3f1960b397caec50db04da19e1d674ccf9dd4b2af5dac010de6b143d5368693",
        "0xc07cba57d6a2df08444c68500177a57837a8be3f9dd59e0885d7cf497ecf3c99",
        "0xda058f8d239a1261661552f5bc0801acfca8a90d08c834df845d6ea1e34e2cfc"
      ],
      otherMSHashes:
        "0x406cfc22544f4c74049983f871d44a3cf2be94bbde3961cab7ed4773d2e57ee0"
    };
  });

  context("getLeafHash", () => {
    it("should calculate leaf hash correctly", async () => {
      (
        await this.storeProofLib.getLeafHash(
          this.data.prefixes,
          this.data.paths,
          this.data.otherMSHashes,
          this.data.key,
          this.data.value
        )
      )
        .toString()
        .should.eq(this.leafHash);
    });

    it("should revert if there is no prefix", async () => {
      await expectRevert(
        this.storeProofLib.getLeafHash(
          [],
          this.data.paths,
          this.data.otherMSHashes,
          this.data.key,
          this.data.value
        ),
        "FIRST_PREFIX_IS_NEEDED"
      );
    });
  });

  context("getAVLHash", () => {
    it("should calculate store hash correctly", async () => {
      (
        await this.storeProofLib.getAVLHash(
          this.data.prefixes,
          this.data.paths,
          this.data.otherMSHashes,
          this.data.key,
          this.data.value
        )
      )
        .toString()
        .should.eq(this.storeHash);
    });

    it("should revert if prefixs's length != path's length + 1", async () => {
      await expectRevert(
        this.storeProofLib.getAVLHash(
          [],
          this.data.paths.filter((_, i) => i > 0),
          this.data.otherMSHashes,
          this.data.key,
          this.data.value
        ),
        "LENGTH_OF_PREFIXS_AND_PATHS_ARE_INCOMPATIBLE"
      );
    });
  });

  context("getAppHash", () => {
    it("should calculate app hash correctly", async () => {
      (
        await this.storeProofLib.getAppHash(
          this.data.prefixes,
          this.data.paths,
          this.data.otherMSHashes,
          this.data.key,
          this.data.value
        )
      )
        .toString()
        .should.eq(this.appHash);
    });
  });
});
