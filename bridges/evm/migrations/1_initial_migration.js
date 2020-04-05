const Migrations = artifacts.require("Migrations");
//Tests
module.exports = function(deployer) {
  deployer.deploy(Migrations);
};
