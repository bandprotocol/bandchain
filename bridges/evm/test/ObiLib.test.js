const { expectRevert } = require("openzeppelin-test-helpers");
const ObiMock = artifacts.require("ObiMock");

require("chai").should();

contract("Obi", () => {

    context("Obi decoder should work correctly", () => {
        beforeEach(async () => {
            this.forTest = await ObiMock.new();
        });

        it("should decodeU8 correctly", async () => {
            result = await this.forTest.decodeU8([0, "0x00"]);
            result.toString().should.eq("0");

            result = await this.forTest.decodeU8([0, "0x01"]);
            result.toString().should.eq("1");

            result = await this.forTest.decodeU8([0, "0x0a"]);
            result.toString().should.eq("10");

            result = await this.forTest.decodeU8([0, "0x0b"]);
            result.toString().should.eq("11");

            result = await this.forTest.decodeU8([1, "0x000b"]);
            result.toString().should.eq("11");

            result = await this.forTest.decodeU8([0, "0x31"]);
            result.toString().should.eq("49");

            result = await this.forTest.decodeU8([0, "0xff"]);
            result.toString().should.eq("255");

            result = await this.forTest.decodeU8([2, "0x0000ff"]);
            result.toString().should.eq("255");

            await expectRevert(this.forTest.decodeU8([3, "0x0000ff"]), 'Obi: Out of range');

        });

        it("should decodeI8 correctly", async () => {
            result = await this.forTest.decodeI8([0, "0x00"]);
            result.toString().should.eq("0");

            result = await this.forTest.decodeI8([0, "0x01"]);
            result.toString().should.eq("1");

            result = await this.forTest.decodeI8([0, "0x0a"]);
            result.toString().should.eq("10");

            result = await this.forTest.decodeI8([0, "0x0b"]);
            result.toString().should.eq("11");

            result = await this.forTest.decodeI8([1, "0x000b"]);
            result.toString().should.eq("11");

            result = await this.forTest.decodeI8([0, "0x31"]);
            result.toString().should.eq("49");

            result = await this.forTest.decodeI8([0, "0xff"]);
            result.toString().should.eq("-1");

            result = await this.forTest.decodeI8([2, "0x0000ff"]);
            result.toString().should.eq("-1");

            await expectRevert(this.forTest.decodeI8([3, "0x0000ff"]), 'Obi: Out of range');
        });

        it("should decodeU16 correctly", async () => {
            result = await this.forTest.decodeU16([0, "0xffff"]);
            result.toString().should.eq("65535");

            result = await this.forTest.decodeU16([0, "0xabcd"]);
            result.toString().should.eq("43981");

            result = await this.forTest.decodeU16([1, "0x00abcd"]);
            result.toString().should.eq("43981");

            await expectRevert(this.forTest.decodeU16([2, "0x00abcd"]), 'Obi: Out of range');
        });

        it("should decodeI16 correctly", async () => {
            result = await this.forTest.decodeI16([0, "0xffff"]);
            result.toString().should.eq("-1");

            result = await this.forTest.decodeI16([0, "0x0112"]);
            result.toString().should.eq("274");

            result = await this.forTest.decodeI16([1, "0x00011200"]);
            result.toString().should.eq("274");

            await expectRevert(this.forTest.decodeI16([2, "0x011200"]), 'Obi: Out of range');
        });

        it("should decodeU32 correctly", async () => {
            result = await this.forTest.decodeU32([0, "0xffff0000"]);
            result.toString().should.eq("4294901760");

            result = await this.forTest.decodeU32([0, "0xabcd0000"]);
            result.toString().should.eq("2882338816");

            result = await this.forTest.decodeU32([1, "0x00abcd0000"]);
            result.toString().should.eq("2882338816");

            await expectRevert(this.forTest.decodeU32([2, "0x00abcd0000"]), 'Obi: Out of range');
        });

        it("should decodeI32 correctly", async () => {
            result = await this.forTest.decodeI32([0, "0xffffffff"]);
            result.toString().should.eq("-1");

            result = await this.forTest.decodeI32([0, "0x01120000"]);
            result.toString().should.eq("17956864");

            result = await this.forTest.decodeI32([1, "0x0001120000"]);
            result.toString().should.eq("17956864");

            await expectRevert(this.forTest.decodeI32([2, "0x0001120000"]), 'Obi: Out of range');
        });

        it("should decodeU64 correctly", async () => {
            result = await this.forTest.decodeU64([0, "0xffff000000000000"]);
            result.toString().should.eq("18446462598732840960");

            result = await this.forTest.decodeU64([0, "0xabcd000000000000"]);
            result.toString().should.eq("12379550950711361536");

            result = await this.forTest.decodeU64([1, "0x00abcd000000000000"]);
            result.toString().should.eq("12379550950711361536");

            await expectRevert(this.forTest.decodeU64([2, "0x00abcd000000000000"]), 'Obi: Out of range');
        });

        it("should decodeI64 correctly", async () => {
            result = await this.forTest.decodeI64([0, "0x0011000000001111"]);
            result.toString().should.eq("4785074604085521");

            result = await this.forTest.decodeI64([1, "0x000011000000001111"]);
            result.toString().should.eq("4785074604085521");

            await expectRevert(this.forTest.decodeI64([2, "0x000011000000001111"]), 'Obi: Out of range');
        });

        it("should decodeU128 correctly", async () => {
            result = await this.forTest.decodeU128([0, "0x01000000000000000000000000000000"]);
            result.toString().should.eq("1329227995784915872903807060280344576");

            result = await this.forTest.decodeU128([1, "0x0001000000000000000000000000000000"]);
            result.toString().should.eq("1329227995784915872903807060280344576");

            await expectRevert(this.forTest.decodeU128([2, "0x0001000000000000000000000000000000"]), 'Obi: Out of range');
        });

        it("should decodeI128 correctly", async () => {
            result = await this.forTest.decodeI128([0, "0x00010000000000000000000000000000"]);
            result.toString().should.eq("5192296858534827628530496329220096");

            result = await this.forTest.decodeI128([1, "0x0000010000000000000000000000000000"]);
            result.toString().should.eq("5192296858534827628530496329220096");

            await expectRevert(this.forTest.decodeI128([2, "0x0000010000000000000000000000000000"]), 'Obi: Out of range');
        });

        it("should decodeU256 correctly", async () => {
            result = await this.forTest.decodeU256([0, "0x0100000000000000000000000000000000000000000000000000000000000000"]);
            result.toString().should.eq("452312848583266388373324160190187140051835877600158453279131187530910662656");

            result = await this.forTest.decodeU256([1, "0x000100000000000000000000000000000000000000000000000000000000000000"]);
            result.toString().should.eq("452312848583266388373324160190187140051835877600158453279131187530910662656");

            await expectRevert(this.forTest.decodeU256([2, "0x000100000000000000000000000000000000000000000000000000000000000000"]), 'Obi: Out of range');
        });

        it("should decodeI256 correctly", async () => {
            result = await this.forTest.decodeI256([0, "0x0001000000000000000000000000000000000000000000000000000000000000"]);
            result.toString().should.eq("1766847064778384329583297500742918515827483896875618958121606201292619776");

            result = await this.forTest.decodeI256([1, "0x000001000000000000000000000000000000000000000000000000000000000000"]);
            result.toString().should.eq("1766847064778384329583297500742918515827483896875618958121606201292619776");

            await expectRevert(this.forTest.decodeI256([2, "0x000100000000000000000000000000000000000000000000000000000000000000"]), 'Obi: Out of range');
        });

        it("should decodeBool correctly", async () => {
            result = await this.forTest.decodeBool([0, "0x00"]);
            result.toString().should.eq("false");

            result = await this.forTest.decodeBool([0, "0x01"]);
            result.toString().should.eq("true");

            result = await this.forTest.decodeBool([1, "0x0000"]);
            result.toString().should.eq("false");

            result = await this.forTest.decodeBool([1, "0x0001"]);
            result.toString().should.eq("true");

            await expectRevert(this.forTest.decodeBool([2, "0x0001"]), 'Obi: Out of range');

        });

        it("should decodeBytes correctly", async () => {
            result = await this.forTest.decodeBytes([0, "0x00000003425443"]);
            result.toString().should.eq("0x425443");

            result = await this.forTest.decodeBytes([1, "0x0000000003425443"]);
            result.toString().should.eq("0x425443");

            await expectRevert(this.forTest.decodeBytes([2, "0x0000000003425443"]), 'Obi: Out of range');
        });
    });
});
