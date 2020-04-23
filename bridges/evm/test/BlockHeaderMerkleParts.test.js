const BlockHeaderMerkleParts = artifacts.require("BlockHeaderMerklePartsMock");

contract("BlockHeaderMerkleParts", () => {
    beforeEach(async () => {
        this.contract = await BlockHeaderMerkleParts.new();
    });

    context("getBlockHeader", () => {
        it("should get block header correctly", async () => {
            (
                await this.contract.getBlockHeader(
                    [
                        "0x32FA694879095840619F5E49380612BD296FF7E950EAFB66FF654D99CA70869E",
                        "0x4BAEF831B309C193CC94DCF519657D832563B099A6F62C6FA8B7A043BA4F3B3B",
                        "0x5E1A8142137BDAD33C3875546E42201C050FBCCDCF33FFC15EC5B60D09803A25",
                        "0x004209A161040AB1778E2F2C00EE482F205B28EFBA439FCB04EA283F619478D9",
                        "0x6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D",
                        "0x0CF1E6ECE60E49D19BB57C1A432E805F39BB4F65C366741E4F03FA54FBD90714"
                    ],
                    "0x1CCD765C80D0DC1705BB7B6BE616DAD3CF2E6439BB9A9B776D5BD183F89CA141",
                    381837
                )
            )
                .toString()
                .should.eq("0xa35617a81409ce46f1f820450b8ad4b217d99ae38aaa719b33c4fc52dca99b22");
        });
    });
});
