const Migrations = artifacts.require("Migrations");
//Testsss
module.exports = function(deployer) {
  deployer.deploy(Migrations);
};
