from iconservice import *
from . import utils


def get_block_header(data: bytes, app_hash: bytes, block_height: int):
    return utils.merkle_inner_hash(  # [BlockHeader]
        utils.merkle_inner_hash(  # [3A]
            utils.merkle_inner_hash(  # [2A]
                data[0:32],  # [1A]
                utils.merkle_inner_hash(  # [1B]
                    utils.merkle_leaf_hash(  # [2]
                        utils.encode_varint_unsigned(block_height)
                    ),
                    data[32:64]  # [3]
                )
            ),
            data[64:96]  # [2B]
        ),
        utils.merkle_inner_hash(  # [3B]
            utils.merkle_inner_hash(  # [2C]
                data[96:128],  # [1E]
                utils.merkle_inner_hash(  # [1F]
                    utils.merkle_leaf_hash(  # [A]
                        bytes([32]) + app_hash
                    ),
                    data[128:160]  # [B]
                )
            ),
            data[160:192]  # [2D]
        )
    )
