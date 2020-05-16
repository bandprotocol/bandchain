const Packets = artifacts.require('PacketsMock');

contract('Packets', () => {
  beforeEach(async () => {
    this.contract = await Packets.new();
  });
  context('marshalRequestPacket', () => {
    it('should marshal a request packet correctly', async () => {
      (await this.contract.marshalRequestPacket(
           ['beeb', 1, '0x030000004254436400000000000000', 1, 1]))
          .toString()
          .should.eq(
              '0xd9c589270a046265656210011a0f03000000425443640000000000000020012801',
          );
    });
  });
  context('marshalResponsePacket', () => {
    it('should marshal a response packet correctly', async () => {
      (await this.contract.marshalResponsePacket(
           ['beeb', 3, 1, 1589535020, 1589535022, 1, '0x4bb10e0000000000']))
          .toString()
          .should.eq(
              '0x79b5957c0a04626565621003180120acc2f9f50528aec2f9f50530013a084bb10e0000000000',
          );
    });
  });

  context('getResultHash', () => {
    it('should calculate result hash correctly', async () => {
      (await this.contract.getResultHash(
           ['beeb', 1, '0x030000004254436400000000000000', 1, 1],
           ['beeb', 3, 1, 1589535020, 1589535022, 1, '0x4bb10e0000000000'],
           ))
          .toString()
          .should.eq(
              '0xdbbbf5596a975c50c601bdd6ae26a5007e8483344afd7d2ae41e37891cb81b86',
          );
    });
  });
});
