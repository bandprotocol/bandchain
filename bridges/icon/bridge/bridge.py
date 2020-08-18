from iconservice import *
from .utils import *
from .pyobi import *

TAG = "BRIDGE"


class BRIDGE(IconScoreBase):
    def __init__(self, db: IconScoreDatabase) -> None:
        super().__init__(db)
        # a single integer
        self.total_validator_power = VarDB("total_validator_power", db, value_type=int)

        # pubkey:bytes => voting_power:int
        self.validator_powers = DictDB("validator_powers", db, value_type=int)

        # block_number:int => oracle_state_hash:bytes
        self.oracle_state = DictDB("oracle_state", db, value_type=bytes)

	# encoded_request:bytes => encoded_response:bytes
        self.requests_cache = DictDB("requests_cache", db, value_type=bytes)

    # For unit tests
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
                revert(f"PUBKEY_SHOULD_BE_64_BYTES_BUT_GOT_{len(pubkey)}_BYTES")

            self.validator_powers[pubkey] = power
            sum_power += power

        self.total_validator_power.set(sum_power)

    def on_update(self) -> None:
        super().on_update()

    # Get the total voting power of active validators currently on duty.
    @external(readonly=True)
    def get_total_validator_power(self) -> int:
        return self.total_validator_power.get()

    # Get the hash of "oracle" iAVL Merkle tree from block height.
    # @param block_height The height of block that the hash was relayed.
    @external(readonly=True)
    def get_oracle_state(self, block_height: int) -> bytes:
        return self.oracle_state[block_height]

    # Get voting power of a validator from public key.
    # @param pub_key Public key of the validator
    @external(readonly=True)
    def get_validator_power(self, pub_key: bytes) -> int:
        return self.validator_powers[pub_key]

    # Update validator powers by owner.
    # @param validators_bytes OBI encoded of the changed set of BandChain validators.
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
                revert(f"PUBKEY_SHOULD_BE_64_BYTES_BUT_GOT_{len(pubkey)}_BYTES")

            total_validator_power -= self.validator_powers[pubkey]
            total_validator_power += power

            self.validator_powers[pubkey] = power

        self.total_validator_power.set(total_validator_power)

    # Relays a new oracle state to the bridge contract.
    # @param block_height The height of block to relay to this bridge contract.
    # @param multi_store_bytes OBI encoded of extra multi store to compute app hash. See multi_store lib.
    # @param block_merkle_part_bytes OBI encoded of extra merkle parts to compute block hash. See merkle_part lib.
    # @param signatures_bytes OBI encoded of the signatures signed on this block.
    @external
    def relay_oracle_state(
        self,
        block_height: int,
        multi_store_bytes: bytes,
        block_merkle_part_bytes: bytes,
        signatures_bytes: bytes,
    ) -> None:
        app_hash = multi_store.get_app_hash(multi_store_bytes)
        block_hash = merkle_part.get_block_header(block_merkle_part_bytes, app_hash, block_height)
        recover_signers = tm_signature.recover_signers(signatures_bytes, block_hash)
        sum_voting_power = 0
        signers_checking = set()
        for signer in recover_signers:
            if signer in signers_checking:
                revert(f"REPEATED_PUBKEY_FOUND: {signer.hex()}")

            signers_checking.add(signer)
            sum_voting_power += self.validator_powers[signer]

        if sum_voting_power * 3 <= self.total_validator_power.get() * 2:
            revert("INSUFFICIENT_VALIDATOR_SIGNATURES")

        self.oracle_state[block_height] = multi_store_bytes[64:96]

    # Verifies that the given data is a valid data on BandChain as of the given block height.
    # @param block_height The block height. Someone must already relay this block.
    # @param encode_packet The OBI encoded of a request packet and a response packet of this request.
    # @param version Lastest block height that the data node was updated.
    # @param iavl_merkle_paths The OBI encoded of merkle proof that shows how the data leave is part of the oracle iAVL.
    @external(readonly=True)
    def verify_oracle_data(
        self, block_height: int, encode_packet: bytes, version: int, iavl_merkle_paths: bytes
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
                    resolve_status: u32,
                    result: bytes
                }
            }
            """
        ).decode(encode_packet)

        current_merkle_hash = sha256.digest(
            # Height of tree (only leaf node) is 0 (signed-varint encode)
            bytes([0])
            + bytes([2])
            + utils.encode_varint_signed(version)  # Size of subtree is 1 (signed-varint encode)
            +
            # Size of data key (1-byte constant 0x01 + 8-byte request ID)
            bytes([9])
            + b"\xff"
            + packet["res"][  # Constant 0xff prefix data request info storage key
                "request_id"
            ].to_bytes(8, "big")
            + bytes([32])
            + sha256.digest(encode_packet)  # Size of data hash
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
        ).decode(iavl_merkle_paths)

        # Goes step-by-step computing hash of parent nodes until reaching root node.
        for path in len_merkle_paths:
            current_merkle_hash = iavl_merkle_path.get_parent_hash(
                path["is_data_on_right"],
                path["subtree_height"],
                path["subtree_size"],
                path["subtree_version"],
                path["sibling_hash"],
                current_merkle_hash,
            )

        # Verifies that the computed Merkle root matches what currently exists.
        if current_merkle_hash != oracle_state_root:
            revert("INVALID_ORACLE_DATA_PROOF")

        return packet

    # Performs oracle state relay and oracle data verification in one go. The caller submits
    # the OBI encoded proof and receives back the decoded data, ready to be validated and used.
    # @param `proof` The OBI encoded data for oracle state relay and data verification.
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
                iavl_merkle_paths: bytes
            }
            """
        ).decode(proof)

        self.relay_oracle_state(
            proof_dict["block_height"],
            proof_dict["multi_store"],
            proof_dict["merkle_parts"],
            proof_dict["signatures"],
        )

        return self.verify_oracle_data(
            proof_dict["block_height"],
            proof_dict["encoded_packet"],
            proof_dict["version"],
            proof_dict["iavl_merkle_paths"],
        )

    # Returns the response packet for a given request packet.
    # @param `encoded_request` The OBI encoded of request packet struct.
    @external(readonly=True)
    def get_latest_response(self, encoded_request: bytes) -> dict:
        res = self.requests_cache[encoded_request]
        if not res:
            return None
        return PyObi(
            """
            {
                client_id: string,
                request_id: u64,
                ans_count: u64,
                request_time: u64,
                resolve_time: u64,
                resolve_status: u32,
                result: bytes
            }
            """
        ).decode(res)

    # Performs oracle state relay and oracle data verification in one go.
    # After that, the results will be recorded to the state by using the OBI encoded of RequestPacket as key.
    # @param `proof` The OBI encoded data for oracle state relay and data verification.
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

        prev_res = self.get_latest_response(req_key)
        prev_resolve_time = (
            (
                prev_res.get("resolve_time", 0)
                if isinstance(prev_res.get("resolve_time", 0), int)
                else 0
            )
            if isinstance(prev_res, dict)
            else 0
        )

        if prev_resolve_time >= res["resolve_time"]:
            revert("FAIL_LATEST_REQUEST_SHOULD_BE_NEWEST")

        if res["resolve_status"] != 1:
            revert("FAIL_REQUEST_IS_NOT_SUCCESSFULLY_RESOLVED")

        self.requests_cache[req_key] = PyObi(
            """
            {
                client_id: string,
                request_id: u64,
                ans_count: u64,
                request_time: u64,
                resolve_time: u64,
                resolve_status: u32,
                result: bytes
            }
            """
        ).encode(res)
