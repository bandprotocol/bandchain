const MultiStoreMerkleParts = artifacts.require("MultiStoreMerklePartsMock");

contract("MultiStoreMerkleParts", () => {
    beforeEach(async () => {
        this.contract = await MultiStoreMerkleParts.new();
    });

    context("getAppHash", () => {
        it("should get app hash correctly", async () => {
            (
                await this.contract.getAppHash(
                    [
                        "0x9F08CCC91501E4FB16E3836B1E2BA03468308E62CFF1755D523585FE478BD9B9",
                        "0x7406A3DD093D0B2EDAE2815B981323119FD77AD26A6BCEF6EA8DFA9CBC747D34",
                        "0xD0EE29EDB1A80F80B6DC2C058B07E85846E2A1D4EC49FCE1DD0CF1B946CCF456",
                        "0x7C42CE4D440B05F9B1019577C43BE61C2E00FCEE065CC04945DA2B967B70F501",
                        "0x1F1F381B147866B258B3F9B57B79BCEA8C8B42E66EA8FE02A663C7D768B4DFC8"
                    ],
                    "0x066f7261636c6520"
                )
            )
                .toString()
                .should.eq("0x2f3beac1586c205052b74e1cf3d284cd022f739200b74cb51b910d0f3d0bf13d");
        });
    });
});
