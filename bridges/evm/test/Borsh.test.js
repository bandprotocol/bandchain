const { expectRevert } = require("openzeppelin-test-helpers");
const BorshUser = artifacts.require("BorshUser");

require("chai").should();

contract("Borsh", () => {

  context("Borsh decoder should work correctly", () => {
    beforeEach(async () => {
      this.forTest = await BorshUser.new();
    });

    it("should decode correctly", async () => {
      result = await this.forTest.decode("0x00000003425443000000000000003264");
      result[0].toString().should.eq("BTC");
      result[1].toString().should.eq("50");
      result[2].toString().should.eq("100");
    });

    it("should decode correctly 2", async () => {
      result = await this.forTest.decode("0x0000000462616e64000000000000019064");
      result[0].toString().should.eq("band");
      result[1].toString().should.eq("400");
      result[2].toString().should.eq("100");
    });

    it("should revert if invalid bytes", async () => {
      await expectRevert(
        this.forTest.decode("0x000000034254433200000000000064"),
        "Borsh: Out of range"
      );
    });
  });
});
