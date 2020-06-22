const MultiStore = artifacts.require("MultiStoreMock");

contract("MultiStore", () => {
    beforeEach(async () => {
        this.contract = await MultiStore.new();
    });

    context("getAppHash", () => {
        it("should get app hash correctly", async () => {
            (
                await this.contract.getAppHash(
                    [
                        "0xBD3021997D04B26F04056F040DDAB5EB0249862398617FF69A477D836619E2B8",
                        "0xF8231CF7FDFEDC7A8BC63140965D6A707F6F302C11DB997FC44C0595053479A0",
                        "0x153AE305B0C588EB5AAD3E685FB740B270FD287ECE08BCC85D668FCB10E3C1F5",
                        "0xB1F2FD852E790E735CA2D3014F96A2A53C60393E9C6BBF941B9A6DD6A05CF6F9",
                        "0xBDFCCB8D10C48AEFE3FE6872BC558C3C17CA7927A82D48E00C1F7199758227AC"
                    ]
                )
            )
                .toString()
                .should.eq("0x0def6341481c4370d561d546c268fb6afed8520689b70abc615c47cfb2a0eee8");
        });
    });
});
