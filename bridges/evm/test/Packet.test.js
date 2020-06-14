const Packets = artifacts.require("PacketsMock");

contract("Packets", () => {
  beforeEach(async () => {
    this.contract = await Packets.new();
  });
  context("encodeRequestPacket", () => {
    it("should marshal a request packet correctly", async () => {
      (
        await this.contract.encodeRequestPacket([
          "beeb",
          1,
          "0x030000004254436400000000000000",
          1,
          1,
        ])
      )
        .toString()
        .should.eq(
          "0x000000046265656200000000000000010000000f03000000425443640000000000000000000000000000010000000000000001",
        );
    });
  });
  context("encodeResponsePacket", () => {
    it("should marshal a response packet correctly", async () => {
      (
        await this.contract.encodeResponsePacket([
          "beeb",
          1,
          1,
          1589535020,
          1589535022,
          1,
          "0x4bb10e0000000000",
        ])
      )
        .toString()
        .should.eq(
          "0x000000046265656200000000000000010000000000000001000000005ebe612c000000005ebe612e00000001000000084bb10e0000000000",
        );
    });
  });

  context("getEncodedResult", () => {
    it("should calculate result hash correctly", async () => {
      (
        await this.contract.getEncodedResult(
          ["beeb", 1, "0x030000004254436400000000000000", 1, 1],
          ["beeb", 1, 1, 1589535020, 1589535022, 1, "0x4bb10e0000000000"],
        )
      )
        .toString()
        .should.eq(
          "0x000000046265656200000000000000010000000f03000000425443640000000000000000000000000000010000000000000001000000046265656200000000000000010000000000000001000000005ebe612c000000005ebe612e00000001000000084bb10e0000000000",
        );
    });
    it("should calculate result hash with empty client id correctly", async () => {
      (
        await this.contract.getEncodedResult(
          ["", 1, "0x030000004254436400000000000000", 1, 1],
          ["", 1, 1, 1590490752, 1590490756, 1, "0x568c0d0000000000"],
        )
      )
        .toString()
        .should.eq(
          "0x0000000000000000000000010000000f030000004254436400000000000000000000000000000100000000000000010000000000000000000000010000000000000001000000005eccf680000000005eccf6840000000100000008568c0d0000000000",
        );
    });
  });
});
