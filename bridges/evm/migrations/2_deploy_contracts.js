const Bridge = artifacts.require("Bridge");

module.exports = function(deployer) {
  deployer.deploy(Bridge);
};
