const ValidatorsManager = artifacts.require("ValidatorsManager");

contract(
  "ValidatorsManager",
  ([owner, alice, v1, v2, v3, v4, v5, v6, v7, v8]) => {
    beforeEach(async () => {
      this.vm = await ValidatorsManager.new([v1], { from: owner });
    });
    it("should be able to verify validator correctly", async () => {
      (await this.vm.validators(v1)).should.eq(true);
      for (const v of [v2, v3, v4, v5, v6, v7, v8]) {
        (await this.vm.validators(v)).should.eq(false);
      }
    });
  }
);
