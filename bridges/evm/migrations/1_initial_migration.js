const Migrations = artifacts.require("Migrations");
//Testsssss
module.exports = function(deployer) {
  deployer.deploy(Migrations);
};
