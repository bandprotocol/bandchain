/// <reference types="cypress" />

describe("Login", () => {
  beforeEach(() => {
    cy.visit("/");
  });

  it("Should have s address on account panel", () => {
    cy.get('[id="connectButton"]').click();
    cy.get('[id="mnemonicInput"]').type("s");
    cy.get('[id="mnemonicConnectButton"] > button').click();
    cy.get('[id="userInfoButton"]').click();
    cy.get('[id="addressWrapper"] > a > span').should(
      "contain",
      "band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte"
    );
  });
});

describe("Send", () => {
  it("Status should be Success", () => {
    cy.get('[id="getFreeButton"] > button').click();
    cy.get('[id="sendToken"] > button').wait(500).contains("Send").click();
    cy.get('[id="recipientAddressInput"]')
      .wait(1000)
      .type("band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte")
      .get('[id="sendAmountInput"]')
      .type("2");
    cy.get('[id="nextButton"]').contains("Next").click();
    cy.get('[id="broadcastButton"]').wait(1000).click();
    cy.get('[id="successMsgContainer"] > span').should(
      "contain",
      "Broadcast Transaction Success"
    );
  });
});
