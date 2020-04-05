const Migrations = artifacts.require("Migrations");
//Testssssss
module.exports = function(deployer) {
  deployer.deploy(Migrations);
};
