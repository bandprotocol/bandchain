const Migrations = artifacts.require("Migrations");
//Test
module.exports = function(deployer) {
  deployer.deploy(Migrations);
};
