import axios from 'axios'
import { Client } from '../../src'
import { Address } from '../../src/wallet'

jest.mock('axios')
const mockedAxios = axios as jest.Mocked<typeof axios>

const TEST_RPC = 'https://api-mock.bandprotocol.com/rest'

const client = new Client(TEST_RPC)

describe('Client get data', () => {
  it('get chain ID', () => {
    const resp = { data: { chain_id: 'bandchain' } }
    mockedAxios.get.mockResolvedValue(resp)

    const response = client.getChainID()
    response.then((e) => expect(e).toEqual('bandchain'))
  })

  it('get data source by ID', () => {
    const resp = {
      data: {
        height: '651093',
        result: {
          owner: 'band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs',
          name: 'CoinGecko Cryptocurrency Price',
          description:
            'Retrieves current price of a cryptocurrency from https://www.coingecko.com',
          filename:
            'c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0',
        },
      },
    }
    mockedAxios.get.mockResolvedValue(resp)

    const expected = {
      owner: Address.fromAccBech32(
        'band1m5lq9u533qaya4q3nfyl6ulzqkpkhge9q8tpzs',
      ),
      name: 'CoinGecko Cryptocurrency Price',
      description:
        'Retrieves current price of a cryptocurrency from https://www.coingecko.com',
      fileName:
        'c56de9061a78ac96748c83e8a22330accf6ee8ebb499c8525613149a70ec49d0',
    }
    const response = client.getDataSource(1)
    response.then((e) => expect(e).toEqual(expected))
  })
})

describe('get latest block', () => {
  it('success', () => {
    const resp = {
      data: {
        block_id: {
          hash:
            '391E99908373F8590C928E0619956DA3D87EB654445DA4F25A185C9718561D53',
          parts: {
            total: '1',
            hash:
              '9DC1DB39B7DDB97DE353DFB2898198BAADEFB7DF8090BF22678714F769D69F42',
          },
        },
        block: {
          header: {
            version: { block: '10', app: '0' },
            chain_id: 'bandchain',
            height: '1032007',
            time: '2020-11-05T09:15:18.445494105Z',
            last_block_id: {
              hash:
                '4BC01E257662B5F9337D615D06D5D19D8D79F7BA5A4022A85B4DC4ED3C59F7CA',
              parts: {
                total: '1',
                hash:
                  '6471C0A51FB7043171EAA76CAFA900B36A4494F47F975980732529D247E8EBA5',
              },
            },
            last_commit_hash:
              '17B2CE4ABA910E85847537F1323DB95C9F16C20C60E9B9BBB04C633C3125BD92',
            data_hash:
              'EFE5E3F549554FEE8EB9B393740C250D74580427A96A175ABB105806039CFFE2',
            validators_hash:
              'E3F0EA129867E1AB4D7D6A97C23771D4D89B9E4DFE0A5B11E03B681244E00151',
            next_validators_hash:
              'E3F0EA129867E1AB4D7D6A97C23771D4D89B9E4DFE0A5B11E03B681244E00151',
            consensus_hash:
              '0EAA6F4F4B8BD1CC222D93BBD391D07F074DE6BE5A52C6964875BB355B7D0B45',
            app_hash:
              '6E2B1ECE9D912D86C25182E8B7419583ABCE978BFC66DC2556BB0D06A8D528EF',
            last_results_hash: '',
            evidence_hash: '',
            proposer_address: 'BDB6A0728C8DFE2124536F16F2BA428FE767A8F9',
          },
          data: {
            txs: [
              'yAEoKBapCj5CcI40CAESDwAAAANCVEMAAAAAAAAAARgEIAMqC2Zyb21fcHliYW5kMhSQ78AMmxLrubEOPhhIwKK5oyk9oBIQCgoKBXViYW5kEgEwEMCEPRpqCibrWumHIQP+cIvaZlJP0sa86QaC44VVqFHgPSruT2KbBd6Q9R7ZvBJANbPqLRAgwwULWWwb5O2/eb6ddtDr65kRFgDcOTTGIqQS5Iz1NvHWHfkPKRoM8egErMsgE9YnuE+pAqoc/YvNfiIEVEVTVA==',
            ],
          },
          evidence: { evidence: null },
          last_commit: {
            height: '1032006',
            round: '0',
            block_id: {
              hash:
                '4BC01E257662B5F9337D615D06D5D19D8D79F7BA5A4022A85B4DC4ED3C59F7CA',
              parts: {
                total: '1',
                hash:
                  '6471C0A51FB7043171EAA76CAFA900B36A4494F47F975980732529D247E8EBA5',
              },
            },
            signatures: [
              {
                block_id_flag: 3,
                validator_address: '5179B0BB203248E03D2A1342896133B5C58E1E44',
                timestamp: '2020-11-05T09:15:18.53815896Z',
                signature:
                  'TZY24CKwZOE8wqfE0NM3qzkQ7qCpCrGEHNZdf8n31L4otZzbKGfOL05kGtBsGkTnZkVv7aJmrJ7XbvIzv0SREQ==',
              },
              {
                block_id_flag: 2,
                validator_address: 'BDB6A0728C8DFE2124536F16F2BA428FE767A8F9',
                timestamp: '2020-11-05T09:15:18.445494105Z',
                signature:
                  'mcUMQtCR38MK69IeUDri0zkfllsXKgnVFTsFwNaO/7cnBaIUUz9U4d3H9EjSH4kANJxWRFO3dSnPy1uOD36K6A==',
              },
              {
                block_id_flag: 3,
                validator_address: 'F0C23921727D869745C4F9703CF33996B1D2B715',
                timestamp: '2020-11-05T09:15:18.537783045Z',
                signature:
                  'fpr26xz+Gg5Rl7Fvx34a0QZpb5yJc5P4t5Z1OctIDQ0VMmh9vEWagsqQGErt1bm+CaKFtkFfZZ4CU0DKN27GbQ==',
              },
              {
                block_id_flag: 3,
                validator_address: 'F23391B5DBF982E37FB7DADEA64AAE21CAE4C172',
                timestamp: '2020-11-05T09:15:18.539946947Z',
                signature:
                  'KGsiIaralMMr1M9A7Ulhbc0THt1pLgNIrNQ6Kx+oGtwjl2w5ke5iivAAtZMduhyIAUMhrp6PN7ZvKgv9TPCNNw==',
              },
            ],
          },
        },
      },
    }

    mockedAxios.get.mockResolvedValue(resp)
    const expected = {
      blockID: {
        hash: Buffer.from(
          '391E99908373F8590C928E0619956DA3D87EB654445DA4F25A185C9718561D53',
          'hex',
        ),
      },
      block: {
        header: {
          chainID: 'bandchain',
          height: 1032007,
          time: new Date('2020-11-05T09:15:18.445494105Z'),
          lastCommitHash: Buffer.from(
            '17B2CE4ABA910E85847537F1323DB95C9F16C20C60E9B9BBB04C633C3125BD92',
            'hex',
          ),
          dataHash: Buffer.from(
            'EFE5E3F549554FEE8EB9B393740C250D74580427A96A175ABB105806039CFFE2',
            'hex',
          ),
          validatorsHash: Buffer.from(
            'E3F0EA129867E1AB4D7D6A97C23771D4D89B9E4DFE0A5B11E03B681244E00151',
            'hex',
          ),
          nextValidatorsHash: Buffer.from(
            'E3F0EA129867E1AB4D7D6A97C23771D4D89B9E4DFE0A5B11E03B681244E00151',
            'hex',
          ),
          consensusHash: Buffer.from(
            '0EAA6F4F4B8BD1CC222D93BBD391D07F074DE6BE5A52C6964875BB355B7D0B45',
            'hex',
          ),
          appHash: Buffer.from(
            '6E2B1ECE9D912D86C25182E8B7419583ABCE978BFC66DC2556BB0D06A8D528EF',
            'hex',
          ),
          lastResultsHash: Buffer.from(''),
          evidenceHash: Buffer.from(''),
          proposerAddress: Buffer.from(
            'BDB6A0728C8DFE2124536F16F2BA428FE767A8F9',
            'hex',
          ),
        },
      },
    }
    const response = client.getLatestBlock()
    response.then((e) => expect(e).toEqual(expected))
  })
})
