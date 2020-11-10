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
    cy.get('[id="sendToken"] > button', {timeout: 500}).contains("Send").click();
    cy.get('[id="recipientAddressInput"]')
      .wait(2000)
      .type("band1jrhuqrymzt4mnvgw8cvy3s9zhx3jj0dq30qpte")
      .get('[id="sendAmountInput"]')
      .type("2");
    cy.get('[id="nextButtonContainer"] > button').contains("Next").click();
    cy.get('[id="broadcastButtonContainer"] > button').wait(1000).click();
    cy.get('[id="successMsgContainer"] > span').should(
      "contain",
      "Broadcast transaction success"
    );
    cy.get('[id="closeModal"]').click();
  });
});

describe("Delegation", () => {
  it("Should be able to delegate with Carol", () => {
    cy.get('[id="navigationBar"] > div > a').contains("Validator").click();
    cy.get('[id="validatorsSection"] > div > div > div > a > span', {timeout: 1000})
      .contains("Carol")
      .should('be.visible')
      .click();
    cy.get('[id="validatorDelegationinfoDlegate"] button:nth-of-type(1)', {timeout: 1000})
      .click();
    cy.get('[id="nextButtonContainer"] > button', {timeout: 1000}).should("be.disabled");
    cy.get('[id="delegateAmountInput').type("1");
    cy.get('[id="memoInput"]').type("cypress");
    cy.get('[id="nextButtonContainer"] > button', {timeout: 1000}).click();
    cy.get('[id="broadcastButtonContainer"] > button', {timeout: 1000}).click();
    cy.get('[id="successMsgContainer"] > span').should(
      "contain",
      "Broadcast transaction success"
    );
    cy.get('[id="closeModal"]').click();
  });

  it("Should be able to undelegate with Carol", () => {
    cy.get('[id="validatorDelegationinfoDlegate"] button:nth-of-type(2)', {timeout: 1000})
      .should('be.visible')
      .click();
    cy.get('[id="undelegateAmountInput').type("0.5");
    cy.get('[id="memoInput"]').type("cypress");
    cy.get('[id="nextButtonContainer"] > button', {timeout: 1000}).click();
    cy.get('[id="broadcastButtonContainer"] > button', {timeout: 1000}).click();
    cy.get('[id="successMsgContainer"] > span').should(
      "contain",
      "Broadcast transaction success"
    );
    cy.get('[id="closeModal"]').click();
  });

  it("Should be able to redelegate with Carol", () => {
    cy.get('[id="validatorDelegationinfoDlegate"] button:nth-of-type(3)', {timeout: 1000})
      .should('be.visible')
      .click();
    cy.get('[id="redelegateContainer"] > div:nth-of-type(1)').click();
    cy.get('[id="redelegateContainer"] input').wait(500).type("Bobby.fish 🐡");
    cy.get(
      '[id="redelegateContainer"] > div:nth-of-type(1) > div:nth-of-type(2) > div:nth-of-type(1) > div:nth-of-type(1)'
    ).click();
    cy.get('[id="redelegateAmountInput').type("0.5");
    cy.get('[id="memoInput"]').type("cypress");
    cy.get('[id="nextButtonContainer"] > button', {timeout: 1000}).click();
    cy.get('[id="broadcastButtonContainer"] > button', {timeout: 1000}).click();
    cy.get('[id="successMsgContainer"] > span').should(
      "contain",
      "Broadcast transaction success"
    );
    cy.get('[id="closeModal"]').click();
  });

  it("Should be able to withdraw reward with Carol", () => {
    cy.get('[id="withdrawRewardContainer"] > button', {timeout: 1000})
    .should('be.visible')
    .click();
    cy.get('[id="memoInput"]').type("cypress");
    cy.get('[id="nextButtonContainer"] > button', {timeout: 1000}).click();
    cy.get('[id="broadcastButtonContainer"] > button', {timeout: 1000}).click();
    cy.get('[id="successMsgContainer"] > span').should(
      "contain",
      "Broadcast transaction success"
    );
    cy.get('[id="closeModal"]').click();
  });
});
