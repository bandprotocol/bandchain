import { Blockscon } from '../../api/blocks/blocks.js';
import { Proposals } from '../../api/proposals/proposals.js';
import { ValidatorRecords, Analytics, MissedBlocksStats, MissedBlocks, AverageData, AverageValidatorData } from '../../api/records/records.js';
// import { Status } from '../../api/status/status.js';
import { Transactions } from '../../api/transactions/transactions.js';
import { ValidatorSets } from '../../api/validator-sets/validator-sets.js';
import { Validators } from '../../api/validators/validators.js';
import { VotingPowerHistory } from '../../api/voting-power/history.js';
import { Evidences } from '../../api/evidences/evidences.js';
import { CoinStats } from '../../api/coin-stats/coin-stats.js';
import { ChainStates } from '../../api/chain/chain.js';

ChainStates.rawCollection().createIndex({height: -1},{unique:true});

Blockscon.rawCollection().createIndex({height: -1},{unique:true});
Blockscon.rawCollection().createIndex({proposerAddress:1});

Evidences.rawCollection().createIndex({height: -1});

Proposals.rawCollection().createIndex({proposalId: 1}, {unique:true});

ValidatorRecords.rawCollection().createIndex({address:1,height: -1}, {unique:1});
ValidatorRecords.rawCollection().createIndex({address:1,exists:1, height: -1});

Analytics.rawCollection().createIndex({height: -1}, {unique:true})

MissedBlocks.rawCollection().createIndex({proposer:1, voter:1, updatedAt: -1});
MissedBlocks.rawCollection().createIndex({proposer:1, blockHeight:-1});
MissedBlocks.rawCollection().createIndex({voter:1, blockHeight:-1});
MissedBlocks.rawCollection().createIndex({voter:1, proposer:1, blockHeight:-1}, {unique:true});

MissedBlocksStats.rawCollection().createIndex({proposer:1});
MissedBlocksStats.rawCollection().createIndex({voter:1});
MissedBlocksStats.rawCollection().createIndex({proposer:1, voter:1},{unique:true});

AverageData.rawCollection().createIndex({type:1, createdAt:-1},{unique:true});
AverageValidatorData.rawCollection().createIndex({proposerAddress:1,createdAt:-1},{unique:true});
// Status.rawCollection.createIndex({})

Transactions.rawCollection().createIndex({txhash:1},{unique:true});
Transactions.rawCollection().createIndex({height:-1});
// Transactions.rawCollection().createIndex({action:1});
Transactions.rawCollection().createIndex({"events.attributes.key":1});
Transactions.rawCollection().createIndex({"events.attributes.value":1});

ValidatorSets.rawCollection().createIndex({block_height:-1});

Validators.rawCollection().createIndex({address:1},{unique:true, partialFilterExpression: { address: { $exists: true } } });
Validators.rawCollection().createIndex({consensus_pubkey:1},{unique:true});
Validators.rawCollection().createIndex({"pub_key.value":1},{unique:true, partialFilterExpression: { "pub_key.value": { $exists: true } }});

VotingPowerHistory.rawCollection().createIndex({address:1,height:-1});
VotingPowerHistory.rawCollection().createIndex({type:1});

CoinStats.rawCollection().createIndex({last_updated_at:-1},{unique:true});
