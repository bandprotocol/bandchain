from iconservice import *
from . import utils, sha256

# Computes Tendermint's application state hash at this given block. AppHash is actually a
# Merkle hash on muliple stores.
#                         ________________[AppHash]_______________
#                        /                                        \
#             _______[I9]______                          ________[I10]________
#            /                  \                       /                     \
#       __[I5]__             __[I6]__              __[I7]__               __[I8]__
#      /         \          /         \           /         \            /         \
#    [I1]       [I2]     [I3]        [I4]       [8]        [9]          [A]        [B]
#   /   \      /   \    /    \      /    \
# [0]   [1]  [2]   [3] [4]   [5]  [6]    [7]
# [0] - acc      [1] - distr   [2] - evidence  [3] - gov
# [4] - main     [5] - mint    [6] - oracle    [7] - params
# [8] - slashing [9] - staking [A] - supply    [D] - upgrade
# Notice that NOT all leaves of the Merkle tree are needed in order to compute the Merkle
# root hash, since we only want to validate the correctness of [6] In fact, only
# [7], [I3], [I5], and [I10] are needed in order to compute [AppHash].


def get_app_hash(multi_store: bytes) -> bytes:
    acc_to_gov_stores_merkle_hash = multi_store[0:32]  # [I5]
    main_and_mint_stores_merkle_hash = multi_store[32:64]  # [I3]
    oracle_iavl_state_hash = multi_store[64:96]  # [6]
    params_stores_merkle_hash = multi_store[96:128]  # [7]
    slashing_to_upgrade_stores_merkle_hash = multi_store[128:160]  # [I10]
    return utils.merkle_inner_hash(  # [AppHash]
        utils.merkle_inner_hash(  # [I9]
            acc_to_gov_stores_merkle_hash,  # [I5]
            utils.merkle_inner_hash(  # [I6]
                main_and_mint_stores_merkle_hash,  # [I3]
                utils.merkle_inner_hash(
                    utils.merkle_leaf_hash(  # [I4]
                        # oracle prefix (uint8(6) + "oracle" + uint8(32))
                        bytes.fromhex("066f7261636c6520")
                        + sha256.digest(sha256.digest(oracle_iavl_state_hash))
                    ),  # [6]
                    params_stores_merkle_hash,  # [7]
                ),
            ),
        ),
        slashing_to_upgrade_stores_merkle_hash,  # [I10]
    )
