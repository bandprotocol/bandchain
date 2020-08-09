from iconservice import *
from . import utils

# @dev Library for computing Tendermint's block header hash from app hash, time, and height.
#
# In Tendermint, a block header hash is the Merkle hash of a binary tree with 14 leaf nodes.
# Each node encodes a data piece of the blockchain. The notable data leaves are: [A] app_hash,
# [2] height, and [3] - time. All data pieces are combined into one 32-byte hash to be signed
# by block validators. The structure of the Merkle tree is shown below.
#
#                                   [BlockHeader]
#                                /                \
#                   [3A]                                    [3B]
#                 /      \                                /      \
#         [2A]                [2B]                [2C]                [2D]
#        /    \              /    \              /    \              /    \
#    [1A]      [1B]      [1C]      [1D]      [1E]      [1F]        [C]    [D]
#    /  \      /  \      /  \      /  \      /  \      /  \
#  [0]  [1]  [2]  [3]  [4]  [5]  [6]  [7]  [8]  [9]  [A]  [B]
#
#  [0] - version               [1] - chain_id            [2] - height        [3] - time
#  [4] - last_block_id         [5] - last_commit_hash    [6] - data_hash     [7] - validators_hash
#  [8] - next_validators_hash  [9] - consensus_hash      [A] - app_hash      [B] - last_results_hash
#  [C] - evidence_hash         [D] - proposer_address
#
# Notice that NOT all leaves of the Merkle tree are needed in order to compute the Merkle
# root hash, since we only want to validate the correctness of [A] and [2]. In fact, only
# [1A], [3], [2B], [1E], [B], and [2D] are needed in order to compute [BlockHeader].


def get_block_header(data: bytes, app_hash: bytes, block_height: int):
    return utils.merkle_inner_hash(  # [BlockHeader]
        utils.merkle_inner_hash(  # [3A]
            utils.merkle_inner_hash(  # [2A]
                data[0:32],  # [1A]
                utils.merkle_inner_hash(  # [1B]
                    utils.merkle_leaf_hash(utils.encode_varint_unsigned(block_height)),  # [2]
                    data[32:64],  # [3]
                ),
            ),
            data[64:96],  # [2B]
        ),
        utils.merkle_inner_hash(  # [3B]
            utils.merkle_inner_hash(  # [2C]
                data[96:128],  # [1E]
                utils.merkle_inner_hash(  # [1F]
                    utils.merkle_leaf_hash(bytes([32]) + app_hash), data[128:160]  # [A]  # [B]
                ),
            ),
            data[160:192],  # [2D]
        ),
    )
