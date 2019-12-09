import { Blockscon } from '../../api/blocks/blocks.js';

let blocksCount = Blockscon.find({}).count();
console.log(blocksCount);
