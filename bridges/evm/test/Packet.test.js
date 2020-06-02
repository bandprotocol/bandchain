const Packets = artifacts.require("PacketsMock");

contract("Packets", () => {
  beforeEach(async () => {
    this.contract = await Packets.new();
  });
  context("marshalRequestPacket", () => {
    it("should marshal a request packet correctly", async () => {
      (
        await this.contract.marshalRequestPacket([
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
  context("marshalResponsePacket", () => {
    it("should marshal a response packet correctly", async () => {
      (
        await this.contract.marshalResponsePacket([
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

  context("getResultHash", () => {
    it("should calculate result hash correctly", async () => {
      (
        await this.contract.getResultHash(
          ["beeb", 1, "0x030000004254436400000000000000", 1, 1],
          ["beeb", 1, 1, 1589535020, 1589535022, 1, "0x4bb10e0000000000"],
        )
      )
        .toString()
        .should.eq(
          "0x29bcc52d59b39c61a9616365dcd39dbff8d1aebc88a6a7e2b53dff67841dbc06",
        );
    });
    it("should calculate result hash with empty client id correctly", async () => {
      (
        await this.contract.getResultHash(
          ["", 1, "0x030000004254436400000000000000", 1, 1],
          ["", 1, 1, 1590490752, 1590490756, 1, "0x568c0d0000000000"],
        )
      )
        .toString()
        .should.eq(
          "0xa506eb6a23931d1130bfced8b10ec41674a6eaefb888063847e3605bffbdd5ba",
        );
    });
  });
});
