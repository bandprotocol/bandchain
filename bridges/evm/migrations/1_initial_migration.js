const Migrations = artifacts.require("Migrations");
//Testss
module.exports = function(deployer) {
  deployer.deploy(Migrations);
};
