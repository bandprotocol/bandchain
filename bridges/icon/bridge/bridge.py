from iconservice import *
from .utils import *
from .pyobi import *

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
        # requests cache
        self.requests_cache = DictDB("requests_cache", db, value_type=dict)

    # For testing
    def set_oracle_state(self, block_height: int, oracle_state: bytes) -> None:
        self.oracle_state[block_height] = oracle_state

    def on_install(self, validators_bytes: bytes) -> None:
        super().on_install()
        # set validators
        obi = PyObi("""[{pubkey:bytes, power:u64}]""")
        sum_power = 0
        for vp in obi.decode(validators_bytes):
            pubkey = vp["pubkey"]
            power = vp["power"]
            if len(pubkey) != 64:
                revert(
                    f'PUBKEY_SHOULD_BE_64_BYTES_BUT_GOT_{len(pubkey)}_BYTES')

            self.validator_powers[pubkey] = power
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

        obi = PyObi("""[{pubkey:bytes, power:u64}]""")
        total_validator_power = self.total_validator_power.get()
        for vp in obi.decode(validators_bytes):
            pubkey = vp["pubkey"]
            power = vp["power"]
            if len(pubkey) != 64:
                revert(
                    f'PUBKEY_SHOULD_BE_64_BYTES_BUT_GOT_{len(pubkey)}_BYTES')

            total_validator_power -= self.validator_powers[pubkey]
            total_validator_power += power

            self.validator_powers[pubkey] = power

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
                revert(f'REPEATED_PUBKEY_FOUND: {signer.hex()}')

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

        packet = PyObi(
            """
            {
                req: {
                    client_id: string,
                    oracle_script_id: u64,
                    calldata: bytes,
                    ask_count: u64,
                    min_count: u64
                },
                res: {
                    client_id: string,
                    request_id: u64,
                    ans_count: u64,
                    request_time: u64,
                    resolve_time: u64,
                    resolve_status: u8,
                    result: bytes
                }
            }
            """
        ).decode(encode_packet)

        current_merkle_hash = sha256.digest(
            # Height of tree (only leaf node) is 0 (signed-varint encode)
            bytes([0]) +
            bytes([2]) +  # Size of subtree is 1 (signed-varint encode)
            utils.encode_varint_signed(version) +
            # Size of data key (1-byte constant 0x01 + 8-byte request ID)
            bytes([9]) +
            b'\xff' +  # Constant 0xff prefix data request info storage key
            packet["res"]["request_id"].to_bytes(8, "big") +
            bytes([32]) +  # Size of data hash
            sha256.digest(encode_packet)
        )

        len_merkle_paths = PyObi(
            """
            [
                {
                    is_data_on_right: bool,
                    subtree_height: u8,
                    subtree_size: u64,
                    subtree_version: u64,
                    sibling_hash: bytes
                }
            ]
            """
        ).decode(merkle_paths)

        # Goes step-by-step computing hash of parent nodes until reaching root node.
        for path in len_merkle_paths:
            current_merkle_hash = iavl_merkle_path.get_parent_hash(
                path["is_data_on_right"],
                path["subtree_height"],
                path["subtree_size"],
                path["subtree_version"],
                path["sibling_hash"],
                current_merkle_hash
            )

        # Verifies that the computed Merkle root matches what currently exists.
        if current_merkle_hash != oracle_state_root:
            revert("INVALID_ORACLE_DATA_PROOF")

        return packet

    @external
    def relay_and_verify(self, proof: bytes) -> dict:
        proof_dict = PyObi(
            """
            {
                block_height: u64,
                multi_store: bytes,
                merkle_parts: bytes,
                signatures: bytes,
                encoded_packet: bytes,
                version: u64,
                merkle_paths: bytes
            }
            """
        ).decode(proof)

        self.relay_oracle_state(
            proof_dict["block_height"],
            proof_dict["multi_store"],
            proof_dict["merkle_parts"],
            proof_dict["signatures"]
        )

        return self.verify_oracle_data(
            proof_dict["block_height"],
            proof_dict["encoded_packet"],
            proof_dict["version"],
            proof_dict["merkle_paths"]
        )

    @external(readonly=True)
    def get_latest_response(self, encoded_request: bytes) -> dict:
        return self.requests_cache[encoded_request]

    @external
    def relay(self, proof: bytes) -> None:
        packet = self.relay_and_verify(proof)
        req = packet["req"]
        res = packet["res"]

        req_key = PyObi(
            """
                {
                    client_id: string,
                    oracle_script_id: u64,
                    calldata: bytes,
                    ask_count: u64,
                    min_count: u64
                }
            """
        ).encode(packet["req"])

        prev_res = self.requests_cache[req_key]
        prev_resolve_time = (
            prev_res.get("resolve_time", 0)
            if isinstance(prev_res.get("resolve_time", 0), int)
            else 0
        ) if isinstance(prev_res, dict) else 0

        if prev_resolve_time >= res["resolve_time"]:
            revert("FAIL_LATEST_REQUEST_SHOULD_BE_NEWEST")

        if res["resolve_status"] != 1:
            revert("FAIL_REQUEST_IS_NOT_SUCCESSFULLY_RESOLVED")

        self.requests_cache[req_key] = res
