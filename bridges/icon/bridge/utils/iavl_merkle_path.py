from iconservice import *
from . import utils, sha256


def get_parent_hash(
    is_data_on_right: bool,
    subtree_height: int,
    subtree_size: int,
    subtree_version: int,
    sibling_hash: bytes,
    data_subtree_hash: bytes
) -> bytes:
    left_subtree = sibling_hash if is_data_on_right else data_subtree_hash
    right_subtree = data_subtree_hash if is_data_on_right else sibling_hash
    return sha256.digest(
        bytes([(subtree_height << 1) & 255]) +
        utils.encode_varint_signed(subtree_size) +
        utils.encode_varint_signed(subtree_version) +
        bytes([32]) +
        left_subtree +
        bytes([32]) +
        right_subtree
    )
