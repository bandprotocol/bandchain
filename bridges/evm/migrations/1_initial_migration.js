const Migrations = artifacts.require("Migrations");
//Testssss
module.exports = function(deployer) {
  deployer.deploy(Migrations);
};
