from iconservice import *
from bridge.utils import *

TAG = 'BRIDGE'


class BRIDGE(IconScoreBase):

    def __init__(self, db: IconScoreDatabase) -> None:
        super().__init__(db)
        # address => voting_power
        self.validator_powers = DictDB("validator_powers", db, value_type=int)
        # total validator power
        self.total_validator_power = VarDB(
            "total_validator_power", db, value_type=int)
        # oracle state
        self.oracle_state = DictDB("oracle_state", db, value_type=bytes)

    # For testing
    def set_oracle_state(self, block_height: int, oracle_state: bytes) -> None:
        self.oracle_state[block_height] = oracle_state

    def on_install(self, validators_bytes: bytes) -> None:
        super().on_install()
        # set validators
        (n, remaining) = obi.decode_int(validators_bytes, 32)
        sum_power = 0
        for i in range(n):
            (pub_key, remaining) = obi.decode_bytes(remaining)
            if len(pub_key) != 64:
                revert(
                    f'PUBKEY_SHOULD_BE_64_BYTES_BUT_GOT_{len(pub_key)}_BYTES')

            (power, remaining) = obi.decode_int(remaining, 64)

            self.validator_powers[pub_key] = power
            sum_power += power

        self.total_validator_power.set(sum_power)

    def on_update(self) -> None:
        super().on_update()

    @external(readonly=True)
    def get_total_validator_power(self) -> int:
        return self.total_validator_power.get()

    @external(readonly=True)
    def get_oracle_state(self, block_height: int) -> bytes:
        return self.oracle_state[block_height]

    @external(readonly=True)
    def get_validator_power(self, pub_key: bytes) -> int:
        return self.validator_powers[pub_key]

    @external
    def update_validator_powers(self, validators_bytes: bytes):
        if self.msg.sender != self.owner:
            revert("NOT_AUTHORIZED")
        (n, remaining) = obi.decode_int(validators_bytes, 32)
        total_validator_power = self.total_validator_power.get()
        for i in range(n):
            (pub_key, remaining) = obi.decode_bytes(remaining)
            if len(pub_key) != 64:
                revert(
                    f'PUBKEY_SHOULD_BE_64_BYTES_BUT_GOT_{len(pub_key)}_BYTES')

            (power, remaining) = obi.decode_int(remaining, 64)

            total_validator_power -= self.validator_powers[pub_key]
            total_validator_power += power

            self.validator_powers[pub_key] = power

        self.total_validator_power.set(total_validator_power)

    @external
    def relay_oracle_state(
        self,
        block_height: int,
        multi_store_bytes: bytes,
        merkle_part_bytes: bytes,
        signatures_bytes: bytes,
    ) -> None:
        app_hash = multi_store.get_app_hash(multi_store_bytes)
        block_hash = merkle_part.get_block_header(
            merkle_part_bytes,
            app_hash,
            block_height
        )
        recover_signers = tm_signature.recover_signers(
            signatures_bytes,
            block_hash
        )
        sum_voting_power = 0
        signers_checking = set()
        for signer in recover_signers:
            if signer in signers_checking:
                revert(f'REPEATED_PUBKEY_FOUND: {signer}')

            signers_checking.add(signer)
            sum_voting_power += self.validator_powers[signer]

        if sum_voting_power * 3 <= self.total_validator_power.get() * 2:
            revert("INSUFFICIENT_VALIDATOR_SIGNATURES")

        self.oracle_state[block_height] = multi_store_bytes[64:96]

    @external(readonly=True)
    def verify_oracle_data(
        self,
        block_height: int,
        encode_packet: bytes,
        version: int,
        merkle_paths: bytes
    ) -> dict:
        oracle_state_root = self.oracle_state[block_height]
        if oracle_state_root == None:
            revert("NO_ORACLE_ROOT_STATE_DATA")

        # request packet
        req = {}
        (req["client_id"], remaining) = obi.decode_str(encode_packet)
        (req["oracle_script_id"], remaining) = obi.decode_int(remaining, 64)
        (req["calldata"], remaining) = obi.decode_bytes(remaining)
        (req["ask_count"], remaining) = obi.decode_int(remaining, 64)
        (req["min_count"], remaining) = obi.decode_int(remaining, 64)

        # response packet
        res = {}
        (_, remaining) = obi.decode_str(remaining)
        (res["request_id"], remaining) = obi.decode_int(remaining, 64)
        (res["ans_count"], remaining) = obi.decode_int(remaining, 64)
        (res["request_time"], remaining) = obi.decode_int(remaining, 64)
        (res["resolve_time"], remaining) = obi.decode_int(remaining, 64)
        (res["resolve_status"], remaining) = obi.decode_int(remaining, 8)
        (res["result"], remaining) = obi.decode_bytes(remaining)

        current_merkle_hash = sha256.digest(
            # Height of tree (only leaf node) is 0 (signed-varint encode)
            b'\x00' +
            b'\x02' +  # Size of subtree is 1 (signed-varint encode)
            utils.encode_varint_signed(version) +
            # Size of data key (1-byte constant 0x01 + 8-byte request ID)
            b'\x09' +
            b'\xff' +  # Constant 0xff prefix data request info storage key
            res["request_id"].to_bytes(8, "big") +
            b'\x20' +  # Size of data hash
            sha256.digest(encode_packet)
        )

        # Goes step-by-step computing hash of parent nodes until reaching root node.
        len_merkle_paths, remaining = obi.decode_int(merkle_paths, 32)
        for i in range(len_merkle_paths):
            is_data_on_right, remaining = obi.decode_bool(remaining)
            subtree_height, remaining = obi.decode_int(remaining, 8)
            subtree_size, remaining = obi.decode_int(remaining, 64)
            subtree_version, remaining = obi.decode_int(remaining, 64)
            sibling_hash, remaining = obi.decode_bytes(remaining)
            current_merkle_hash = iavl_merkle_path.get_parent_hash(
                is_data_on_right,
                subtree_height,
                subtree_size,
                subtree_version,
                sibling_hash,
                current_merkle_hash
            )

        # Verifies that the computed Merkle root matches what currently exists.
        if current_merkle_hash != oracle_state_root:
            revert("INVALID_ORACLE_DATA_PROOF")

        return {"req": req, "res": res}

    @external
    def relay_and_verify(self, proof: bytes) -> dict:
        block_height, remaining = obi.decode_int(proof, 64)
        multi_store, remaining = obi.decode_bytes(remaining)
        merkle_parts, remaining = obi.decode_bytes(remaining)
        signatures, remaining = obi.decode_bytes(remaining)
        self.relay_oracle_state(
            block_height, multi_store, merkle_parts, signatures)

        encode_packet, remaining = obi.decode_bytes(remaining)
        version, remaining = obi.decode_int(remaining, 64)
        merkle_paths, remaining = obi.decode_bytes(remaining)
        return self.verify_oracle_data(
            block_height,
            encode_packet,
            version,
            merkle_paths
        )
