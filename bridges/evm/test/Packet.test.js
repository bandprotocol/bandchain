const Packets = artifacts.require("PacketsMock");

contract("Packets", () => {
  beforeEach(async () => {
    this.contract = await Packets.new();
  });
  context("marshalRequestPacket", () => {
    it("should marshal a request packet correctly", async () => {
      (
        await this.contract.marshalRequestPacket([
          "band test",
          1,
          "030000004254436400000000000000",
          4,
          4,
        ])
      )
        .toString()
        .should.eq(
          "0xd9c589270a0962616e64207465737410011a1e30333030303030303432353434333634303030303030303030303030303020042804",
        );
    });
  });
  context("marshalResponsePacket", () => {
    it("should marshal a response packet correctly", async () => {
      (
        await this.contract.marshalResponsePacket([
          "band test",
          1,
          4,
          1587734008,
          1587734012,
          1,
          "d8720b0000000000",
        ])
      )
        .toString()
        .should.eq(
          "0x79b5957c0a0962616e6420746573741001180420f8cb8bf50528fccb8bf50530013a1064383732306230303030303030303030",
        );
    });
  });

  context("getResultHash", () => {
    it("should calculate result hash correctly", async () => {
      (
        await this.contract.getResultHash(
          ["band test", 1, "030000004254436400000000000000", 4, 4],
          ["band test", 1, 4, 1587734008, 1587734012, 1, "d8720b0000000000"],
        )
      )
        .toString()
        .should.eq(
          "0x63d30f34c4b3439a95386912ec9ee2e9c01666685b6a25b11c96d46d47f37a42",
        );
    });
  });
});
