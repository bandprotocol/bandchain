/// <reference types="cypress" />

describe("Navigation", () => {
  describe("Home", () => {
    describe("Should be able to render Home page", () => {
      before(() => {
        cy.visit("/").wait(1000);
      });

      it("Validate highlight info", () => {
        cy.get('[id="highlight-Band Price"] > span:nth-of-type(1)').should(
          "contain",
          "Band Price"
        );
        cy.get('[id="highlight-Market Cap"] > span:nth-of-type(1)').should(
          "contain",
          "Market Cap"
        );
        cy.get('[id="highlight-Latest Block"] > span:nth-of-type(1)').should(
          "contain",
          "Latest Block"
        );
        cy.get(
          '[id="highlight-Active Validators"] > span:nth-of-type(1)'
        ).should("contain", "Active Validators");
      });

      it("Validate total request graph", () => {
        cy.get('[id="totalRequestsGraphSection"] > h4:nth-of-type(1)').should(
          "contain",
          "Total Requests"
        );
      });

      it("Validate latest requests", () => {
        cy.get(
          '[id="latestRequestsSectionHeader"] > div:nth-of-type(1) > span:nth-of-type(1)'
        ).should("contain", "Latest Requests");
      });

      it("Validate latest blocks", () => {
        cy.get(
          '[id="latestBlocksSectionHeader"] > div:nth-of-type(1) > span:nth-of-type(1)'
        ).should("contain", "Latest Blocks");
      });

      it("Validate latest transactions", () => {
        cy.get(
          '[id="latestTransactionsSectionHeader"] > div:nth-of-type(1) > span:nth-of-type(1)'
        ).should("contain", "Latest Transactions");
      });
    });
  });

  describe("Validators", () => {
    before(() => {
      cy.visit("/validators").wait(1000);
    });

    describe("Should be able to render Validators page", () => {
      it("Validate title", () => {
        cy.get('[id="validatorsSection"] h2').should(
          "contain",
          "All Validators"
        );
      });
    });
  });

  describe("Blocks", () => {
    before(() => {
      cy.visit("/blocks").wait(1000);
    });

    describe("Should be able to render Blocks page", () => {
      it("Validate title", () => {
        cy.get('[id="blocksSection"] h2').should("contain", "All Blocks");
      });
    });
  });

  describe("Transactions", () => {
    before(() => {
      cy.visit("/txs").wait(1000);
    });

    describe("Should be able to render Transactions page", () => {
      it("Validate title", () => {
        cy.get('[id="transactionsSection"] h2').should(
          "contain",
          "All Transactions"
        );
      });
    });
  });

  describe("Proposals", () => {
    before(() => {
      cy.visit("/proposals").wait(1000);
    });

    describe("Should be able to render Transactions page", () => {
      it("Validate title", () => {
        cy.get('[id="proposalsSection"] h2').should("contain", "All Proposals");
      });
    });
  });

  describe("Data Sources", () => {
    before(() => {
      cy.visit("/data-sources").wait(1000);
    });

    describe("Should be able to render Data Sources page", () => {
      it("Validate title", () => {
        cy.get('[id="datasourcesSection"] h2').should(
          "contain",
          "All Data Sources"
        );
      });
    });
  });

  describe("Oracle Scripts", () => {
    before(() => {
      cy.visit("/oracle-scripts").wait(1000);
    });

    describe("Should be able to render Oracle Scripts page", () => {
      it("Validate title", () => {
        cy.get('[id="oraclescriptsSection"] h2').should(
          "contain",
          "All Oracle Scripts"
        );
      });
    });
  });

  describe("Requests", () => {
    before(() => {
      cy.visit("/requests").wait(1000);
    });

    describe("Should be able to render Requests page", () => {
      it("Validate title", () => {
        cy.get('[id="requestsSection"] h2').should("contain", "All Requests");
      });
    });
  });
});
