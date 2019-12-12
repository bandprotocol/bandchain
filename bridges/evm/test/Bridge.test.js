const Bridge = artifacts.require("Bridge");
require("chai").should();

contract("Bridge", ([owner, alice, bob]) => {
  beforeEach(async () => {
    this.bridge = await Bridge.new([], { from: owner });
  });
  it("basics", async () => {
    (await this.bridge.numberOfValidators()).toString().should.eq("0");
  });
});
