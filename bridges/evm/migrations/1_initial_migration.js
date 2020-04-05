const Migrations = artifacts.require("Migrations");
//Testsssssss
module.exports = function(deployer) {
  deployer.deploy(Migrations);
};
