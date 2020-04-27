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
          1,
          1,
        ])
      )
        .toString()
        .should.eq(
          "0xd9c589270a0962616e64207465737410011a1e30333030303030303432353434333634303030303030303030303030303020012801"
        );
    });
  });
});
