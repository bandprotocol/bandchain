const ValidatorsManager = artifacts.require("ValidatorsManager");

contract("ValidatorsManager", accounts => {
  console.log(accounts);
  it("should work", async () => {
    assert.equal(1 + 2, 3);
  });
});
