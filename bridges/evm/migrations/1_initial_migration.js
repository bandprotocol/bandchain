const Migrations = artifacts.require("Migrations");
//Testssssssss
module.exports = function(deployer) {
  deployer.deploy(Migrations);
};
