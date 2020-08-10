import os

from iconsdk.builder.call_builder import CallBuilder
from iconsdk.builder.transaction_builder import DeployTransactionBuilder, CallTransactionBuilder
from iconsdk.wallet.wallet import KeyWallet
from iconsdk.libs.in_memory_zip import gen_deploy_data_content
from iconsdk.signed_transaction import SignedTransaction
from tbears.libs.icon_integrate_test import IconIntegrateTestBase, SCORE_INSTALL_ADDRESS
from ..pyobi import *

DIR_PATH = os.path.abspath(os.path.dirname(__file__))


class TestIntegrationBRIDGE(IconIntegrateTestBase):
    TEST_HTTP_ENDPOINT_URI_V3 = "http://127.0.0.1:9000/api/v3"
    BRIDGE_PROJECT = os.path.abspath(os.path.join(DIR_PATH, ".."))
    RECEIVER_MOCK_PROJECT = os.path.abspath(os.path.join(DIR_PATH, "../../receiver_mock"))
    CACHE_CONSUMER_MOCK_PROJECT = os.path.abspath(
        os.path.join(DIR_PATH, "../../cache_consumer_mock")
    )

    def setUp(self):
        super().setUp()

        self.icon_service = None
        # if you want to send request to network, uncomment next line and set self.TEST_HTTP_ENDPOINT_URI_V3
        # self.icon_service = IconService(HTTPProvider(self.TEST_HTTP_ENDPOINT_URI_V3))

        # install SCORE
        self._bridge_address = self._deploy_bridge()["scoreAddress"]
        self._receiver_mock_address = self._deploy_receiver_mock(self._bridge_address)[
            "scoreAddress"
        ]
        self._cache_consumer_mock = self._deploy_cache_consumer_mock(
            self._bridge_address,
            bytes.fromhex(
                "0000000966726f6d5f7363616e00000000000000010000000f00000003425443000000000000006400000000000000040000000000000004"
            ),
        )["scoreAddress"]

    def _deploy_bridge(self, to: str = SCORE_INSTALL_ADDRESS) -> dict:
        params = {}
        params[
            "validators_bytes"
        ] = "0000000400000040a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba331db530b76775beb0348c52bb8fc1fdac207525e5689caa01c0af8d2f8f371ec5000000000000006400000040724ae29cfeb7497051d09edfd8e822352c4c8361b757647645b78c8cc74ce885f04c26ee07ff6ada08587a4037363838b1dda6e306091ee0690caa8fe0e6fcd2000000000000006400000040f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479000000000000006400000040d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f2c67aaeee2bc8d00ff253a326eb04de8ce88eea1224c2bb8ca69296a1c753f530000000000000064"

        # Generates an instance of transaction for deploying SCORE.
        transaction = (
            DeployTransactionBuilder()
            .from_(self._test1.get_address())
            .to(to)
            .params(params)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .content_type("application/zip")
            .content(gen_deploy_data_content(self.BRIDGE_PROJECT))
            .build()
        )

        # Returns the signed transaction object having a signature
        signed_transaction = SignedTransaction(transaction, self._test1)

        # process the transaction in local
        tx_result = self.process_transaction(signed_transaction, self.icon_service)

        self.assertEqual(True, tx_result["status"])
        self.assertTrue("scoreAddress" in tx_result)

        return tx_result

    def _deploy_receiver_mock(self, bridge_address: str, to: str = SCORE_INSTALL_ADDRESS) -> dict:
        params = {}
        params["bridge_address"] = bridge_address

        # Generates an instance of transaction for deploying SCORE.
        transaction = (
            DeployTransactionBuilder()
            .from_(self._test1.get_address())
            .to(to)
            .params(params)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .content_type("application/zip")
            .content(gen_deploy_data_content(self.RECEIVER_MOCK_PROJECT))
            .build()
        )

        # Returns the signed transaction object having a signature
        signed_transaction = SignedTransaction(transaction, self._test1)

        # process the transaction in local
        tx_result = self.process_transaction(signed_transaction, self.icon_service)

        self.assertEqual(True, tx_result["status"])
        self.assertTrue("scoreAddress" in tx_result)

        return tx_result

    def _deploy_cache_consumer_mock(
        self, bridge_address: str, req_template: bytes, to: str = SCORE_INSTALL_ADDRESS
    ) -> dict:
        params = {}
        params["bridge_address"] = bridge_address
        params["req_template"] = req_template

        # Generates an instance of transaction for deploying SCORE.
        transaction = (
            DeployTransactionBuilder()
            .from_(self._test1.get_address())
            .to(to)
            .params(params)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .content_type("application/zip")
            .content(gen_deploy_data_content(self.CACHE_CONSUMER_MOCK_PROJECT))
            .build()
        )

        # Returns the signed transaction object having a signature
        signed_transaction = SignedTransaction(transaction, self._test1)

        # process the transaction in local
        tx_result = self.process_transaction(signed_transaction, self.icon_service)

        self.assertEqual(True, tx_result["status"])
        self.assertTrue("scoreAddress" in tx_result)

        return tx_result

    def test_receiver_get_bridge_address(self):
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._receiver_mock_address)
            .method("get_bridge_address")
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(self._bridge_address, response)

    def test_cache_consumer_get_bridge_address(self):
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._cache_consumer_mock)
            .method("get_bridge_address")
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(self._bridge_address, response)

    def test_cache_consumer_get_req_template(self):
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._cache_consumer_mock)
            .method("get_request_key_template")
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(
            {
                "client_id": "from_scan",
                "oracle_script_id": 1,
                "calldata": b"\x00\x00\x00\x03BTC\x00\x00\x00\x00\x00\x00\x00d",
                "ask_count": 4,
                "min_count": 4,
            },
            response,
        )

    def test_update_validator_powers_success(self):
        params = {}
        params[
            "pub_key"
        ] = "f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479"
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_validator_power")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual("0x64", response)

        params = {}
        params[
            "validators_bytes"
        ] = "0000000100000040f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b704790000000000000096"

        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("update_validator_powers")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        params = {}
        params[
            "pub_key"
        ] = "f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479"
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_validator_power")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual("0x96", response)

    def test_update_validator_powers_fail_not_authorized(self):
        new_owner = KeyWallet.create()
        params = {}
        params[
            "pub_key"
        ] = "f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479"
        call = (
            CallBuilder()
            .from_(new_owner.get_address())
            .to(self._bridge_address)
            .method("get_validator_power")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual("0x64", response)

        params = {}
        params[
            "validators_bytes"
        ] = "0000000100000040f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b704790000000000000096"

        transaction = (
            CallTransactionBuilder()
            .from_(new_owner.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("update_validator_powers")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, new_owner)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(False, tx_result["status"])

    def test_relay_oracle_state_success_1(self):
        params = {}
        params["block_height"] = 191
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_oracle_state")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(None, response)

        params = {}
        params["block_height"] = 191
        params[
            "multi_store_bytes"
        ] = "685430546d23a44e6b8034eaafbc2f4cd7fef54b54d5b66528cb4e5225bd74fb4f8d0bb0cd3eb9dc70b4dbbea5f0cbd5b523195f7bae02bb401bb00a93aba08e1912057fff0b3e85abf1a319f75d37b21b430f3da7db9e486a5041de47c686d3b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f991cc906286235b676ad402fc04ed768eb2bfeca664e8d595c286571da1433c60"
        params[
            "block_merkle_part_bytes"
        ] = "32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869ed9f175396c0e2d0e77f69856abf5d8e69283cb915eb8886262feda1d519b30054f4d548668a3986db253689234b9cac96303a128b7081b16e53b79ab9e65b887004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01dd991da4d4e69473cc75a4b819f9e07d4956671a6f4a74df4cc16596fcbe68137"
        params[
            "signatures_bytes"
        ] = "00000003000000209279257914bfee6faec46df086e9a673a4c1576fb299094d45f77e12fa3728e70000002004dbf23f5ebb07ba8fd7ca06843dfe363b5e86596930e1889d9bd5bd3e98fc531c000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510979fbd9801320962616e64636861696e00000020826bb17b714ebcd8199ee2a01334102f19248087cfdeee42ebd406b3991c389500000020621281eedf97f3a9ec121224cef9f8c07872d2c81cce44dd846bb3678dd505231b000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510dd89e98a01320962616e64636861696e00000020d775fd0e1580499ef16a4ab1998dcb4cbd47cf6f342cdc51a9552d1434552ed8000000205a1a2075ae97a6bbd07a5a40efe287eceab37c7a7caae82e6c997c51c62334701b000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510ab94a59201320962616e64636861696e"
        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay_oracle_state")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        params = {}
        params["block_height"] = 191
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_oracle_state")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(
            "1912057fff0b3e85abf1a319f75d37b21b430f3da7db9e486a5041de47c686d3", response.hex()
        )

    def test_relay_oracle_state_success_2(self):
        params = {}
        params["block_height"] = 446
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_oracle_state")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(None, response)

        params = {}
        params["block_height"] = 446
        params[
            "multi_store_bytes"
        ] = "b0cc616f4e76eba2c1845fd2321ad881991a2990f17cd278eb5f6d17da032faf9459616c154a23567a7281fb577bbfa1aa8d36382c4f64362968202fd32a0b65cbadb1694a5152c2b03f4960e1229745b39d7c41b32dc54e0c207799ae471981b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f93d913423122f2237a9a046d7ebfd13c3a6f1df93a1517a9f0b60a492dbc369d8"
        params[
            "block_merkle_part_bytes"
        ] = "32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869eeea5d164b32b494765eb256279014366747db8941c23f23b2fb09a9bbbf721762c4fd636adf805490129912f7fc18a2766aabd9fd4807f2ff18ed4ccd23f2583004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d7f4be7e5a1eb872ad44103360ddc190410331280c42a54d829a5d752c796685d"
        params[
            "signatures_bytes"
        ] = "0000000300000020b0eff36d1214c35f785d82a7947589dba3f7d6ae35dcb92b1b3aca1df441e1e7000000207be19dd7f54c0ab0726cc5074fe65a887c2f2cd736cdab1ade1b7bd5c378abe11b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f705108c9fa29b01320962616e64636861696e00000020fc6a4aae096d2606050269fde7ac02f58c23fe10dd219f8805696c8aa3b1bb65000000206c324293f560b082923f9f1f1cc0307a78633095b87e4f4047a57ad650fd971b1b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f70510f6a5fd9c01320962616e64636861696e00000020d039e1d13e1837731a5719626075f1ef92ff48c6fbbec68b2ef9b3b34a6086db00000020278405f72ca62f7ff67c28b896ca38ab0534974d96a24adcefc509e0525301871b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f70510d3d5dd9b01320962616e64636861696e"
        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay_oracle_state")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        params = {}
        params["block_height"] = 446
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_oracle_state")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(
            "cbadb1694a5152c2b03f4960e1229745b39d7c41b32dc54e0c207799ae471981", response.hex()
        )

    def test_relay_oracle_state_success_3(self):
        params = {}
        params["block_height"] = 1412
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_oracle_state")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(None, response)

        params = {}
        params = {
            "block_height": 1412,
            "multi_store_bytes": bytes.fromhex(
                "811032f41a5918cfc860ebcea1b7678564c622b4e3f04e28c9d8b195f04661e11e06ea167a89de60abf52c60350793c3823a0a4b9681f59daeb8f20942b8bcf3ff2b888760ffa9efa09e35763d8a5ba3ec794d00bcd1f9a901fa007f1f16e2ebb1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f9c6507de937080730702a2c39ad8aa4666fd24e12a9496e4546a0bd4031e6d4ee"
            ),
            "block_merkle_part_bytes": bytes.fromhex(
                "32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e2a448c96300f1cb90567d815c6b2f6ec0bba29615fd88deab57299c76f0b30fa3ebf2264c2c941df347eb97f55b4adcd4784f4f0a8e66d9f2e3fda413115e085004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d0efe3e12f46363c7779140d4ce659925db52f19053e114d7cc4efd666b37f79f"
            ),
            "signatures_bytes": bytes.fromhex(
                "0000000300000020628716ac49023de84adddddcbef8007c2e41e5b58306ce87a0afad5447bc6210000000202f520db2bff3003d5612e03b7aaa99472164c73922a977af95e1ffc2a67c53b41b000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510cbe1cfbe02320962616e64636861696e00000020ff2ba7e2bd2175827997c706451b5da768b6873d7ba4129fc6ee54e62ba9c593000000203c7f7e5b08d1733d430658431545c9a2f57e4641b3b4cd52e567f27be9485e601c000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510a3f0e5be02320962616e64636861696e000000205a2f66b4d62d905b98277cd2807a324f0651340e80ae0249e500beb5ddcdce11000000203c1ed3d960b19e0ca7d321215874c6e91407ae3d2748f2e3b617fad833c30b6d1b000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510a394a6bf02320962616e64636861696e"
            ),
        }

        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay_oracle_state")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        params = {}
        params["block_height"] = 1412
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_oracle_state")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(
            "ff2b888760ffa9efa09e35763d8a5ba3ec794d00bcd1f9a901fa007f1f16e2eb", response.hex()
        )

    # fail because sum of validator powers is less than 2/3
    def test_relay_oracle_state_fail_1(self):
        params = {}
        params["block_height"] = 191
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_oracle_state")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(None, response)

        params = {}
        params["block_height"] = 191
        params[
            "multi_store_bytes"
        ] = "685430546d23a44e6b8034eaafbc2f4cd7fef54b54d5b66528cb4e5225bd74fb4f8d0bb0cd3eb9dc70b4dbbea5f0cbd5b523195f7bae02bb401bb00a93aba08e1912057fff0b3e85abf1a319f75d37b21b430f3da7db9e486a5041de47c686d3b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f991cc906286235b676ad402fc04ed768eb2bfeca664e8d595c286571da1433c60"
        params[
            "block_merkle_part_bytes"
        ] = "32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869ed9f175396c0e2d0e77f69856abf5d8e69283cb915eb8886262feda1d519b30054f4d548668a3986db253689234b9cac96303a128b7081b16e53b79ab9e65b887004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01dd991da4d4e69473cc75a4b819f9e07d4956671a6f4a74df4cc16596fcbe68137"
        params[
            "signatures_bytes"
        ] = "00000002000000209279257914bfee6faec46df086e9a673a4c1576fb299094d45f77e12fa3728e70000002004dbf23f5ebb07ba8fd7ca06843dfe363b5e86596930e1889d9bd5bd3e98fc531c000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510979fbd9801320962616e64636861696e00000020d775fd0e1580499ef16a4ab1998dcb4cbd47cf6f342cdc51a9552d1434552ed8000000205a1a2075ae97a6bbd07a5a40efe287eceab37c7a7caae82e6c997c51c62334701b000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510ab94a59201320962616e64636861696e"
        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay_oracle_state")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(False, tx_result["status"])

        params = {}
        params["block_height"] = 191
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_oracle_state")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(None, response)

    # fail because repeated pubkey
    def test_relay_oracle_state_fail_2(self):
        params = {}
        params["block_height"] = 191
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_oracle_state")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(None, response)

        params = {}
        params["block_height"] = 191
        params[
            "multi_store_bytes"
        ] = "685430546d23a44e6b8034eaafbc2f4cd7fef54b54d5b66528cb4e5225bd74fb4f8d0bb0cd3eb9dc70b4dbbea5f0cbd5b523195f7bae02bb401bb00a93aba08e1912057fff0b3e85abf1a319f75d37b21b430f3da7db9e486a5041de47c686d3b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f991cc906286235b676ad402fc04ed768eb2bfeca664e8d595c286571da1433c60"
        params[
            "block_merkle_part_bytes"
        ] = "32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869ed9f175396c0e2d0e77f69856abf5d8e69283cb915eb8886262feda1d519b30054f4d548668a3986db253689234b9cac96303a128b7081b16e53b79ab9e65b887004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01dd991da4d4e69473cc75a4b819f9e07d4956671a6f4a74df4cc16596fcbe68137"
        params[
            "signatures_bytes"
        ] = "00000004000000209279257914bfee6faec46df086e9a673a4c1576fb299094d45f77e12fa3728e70000002004dbf23f5ebb07ba8fd7ca06843dfe363b5e86596930e1889d9bd5bd3e98fc531c000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510979fbd9801320962616e64636861696e00000020826bb17b714ebcd8199ee2a01334102f19248087cfdeee42ebd406b3991c389500000020621281eedf97f3a9ec121224cef9f8c07872d2c81cce44dd846bb3678dd505231b000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510dd89e98a01320962616e64636861696e00000020d775fd0e1580499ef16a4ab1998dcb4cbd47cf6f342cdc51a9552d1434552ed8000000205a1a2075ae97a6bbd07a5a40efe287eceab37c7a7caae82e6c997c51c62334701b000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510ab94a59201320962616e64636861696e000000209279257914bfee6faec46df086e9a673a4c1576fb299094d45f77e12fa3728e70000002004dbf23f5ebb07ba8fd7ca06843dfe363b5e86596930e1889d9bd5bd3e98fc531c000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510979fbd9801320962616e64636861696e"
        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay_oracle_state")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(False, tx_result["status"])

        params = {}
        params["block_height"] = 191
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_oracle_state")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(None, response)

    def test_relay_and_verify_success_1(self):
        params = {}
        params["block_height"] = 446
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_oracle_state")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(None, response)

        params = {}
        params["proof"] = bytes.fromhex(
            "00000000000001be000000a0b0cc616f4e76eba2c1845fd2321ad881991a2990f17cd278eb5f6d17da032faf9459616c154a23567a7281fb577bbfa1aa8d36382c4f64362968202fd32a0b65cbadb1694a5152c2b03f4960e1229745b39d7c41b32dc54e0c207799ae471981b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f93d913423122f2237a9a046d7ebfd13c3a6f1df93a1517a9f0b60a492dbc369d8000000c032fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869eeea5d164b32b494765eb256279014366747db8941c23f23b2fb09a9bbbf721762c4fd636adf805490129912f7fc18a2766aabd9fd4807f2ff18ed4ccd23f2583004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d7f4be7e5a1eb872ad44103360ddc190410331280c42a54d829a5d752c796685d000001e40000000300000020b0eff36d1214c35f785d82a7947589dba3f7d6ae35dcb92b1b3aca1df441e1e7000000207be19dd7f54c0ab0726cc5074fe65a887c2f2cd736cdab1ade1b7bd5c378abe11b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f705108c9fa29b01320962616e64636861696e00000020fc6a4aae096d2606050269fde7ac02f58c23fe10dd219f8805696c8aa3b1bb65000000206c324293f560b082923f9f1f1cc0307a78633095b87e4f4047a57ad650fd971b1b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f70510f6a5fd9c01320962616e64636861696e00000020d039e1d13e1837731a5719626075f1ef92ff48c6fbbec68b2ef9b3b34a6086db00000020278405f72ca62f7ff67c28b896ca38ab0534974d96a24adcefc509e0525301871b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f70510d3d5dd9b01320962616e64636861696e0000006b000000047465737400000000000000010000000f00000003425443000000000000006400000000000000040000000000000004000000047465737400000000000000010000000000000004000000005ef1c3c2000000005ef1c3c9000000010000000800000000000eaae600000000000001bc00000148000000060101000000000000000200000000000001bc000000201d3927941aec0602e1fcab94b4d3e57da16a2455165cb097adfdfdf34a3e16880102000000000000000400000000000001bc0000002017480c9a0faca2c2ecef5fdc7627b11f898f5652f3022e8d0b2da40356d6d8e40103000000000000000700000000000001bc00000020d90fcf3847ccd9ec012100f251f6e258d1e6b4e5e03a94fb06927dbb11f386230104000000000000000d00000000000001bc000000205332e6182a047d2b63df9c4a5edbb34d0760a5eea2816bfcc50ef0a316267a490105000000000000001200000000000001bc00000020fdd7908f23ba8fbe99c7f91c660a94a96231e807566e89fe524e739b1f6252c40107000000000000003800000000000001bd00000020efdf574e7d863cf4d8b8b08ed1558a6e6f47d5091667f29b7b2d37dc44d5f156"
        )
        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay_and_verify")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        params = {}
        params["block_height"] = 446
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_oracle_state")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(
            "cbadb1694a5152c2b03f4960e1229745b39d7c41b32dc54e0c207799ae471981", response.hex()
        )

    def test_relay_and_verify_success_2(self):
        params = {}
        params["block_height"] = 1412
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_oracle_state")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(None, response)

        params = {}
        params["proof"] = bytes.fromhex(
            "0000000000000584000000a0811032f41a5918cfc860ebcea1b7678564c622b4e3f04e28c9d8b195f04661e11e06ea167a89de60abf52c60350793c3823a0a4b9681f59daeb8f20942b8bcf3ff2b888760ffa9efa09e35763d8a5ba3ec794d00bcd1f9a901fa007f1f16e2ebb1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f9c6507de937080730702a2c39ad8aa4666fd24e12a9496e4546a0bd4031e6d4ee000000c032fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e2a448c96300f1cb90567d815c6b2f6ec0bba29615fd88deab57299c76f0b30fa3ebf2264c2c941df347eb97f55b4adcd4784f4f0a8e66d9f2e3fda413115e085004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d0efe3e12f46363c7779140d4ce659925db52f19053e114d7cc4efd666b37f79f000001e40000000300000020628716ac49023de84adddddcbef8007c2e41e5b58306ce87a0afad5447bc6210000000202f520db2bff3003d5612e03b7aaa99472164c73922a977af95e1ffc2a67c53b41b000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510cbe1cfbe02320962616e64636861696e00000020ff2ba7e2bd2175827997c706451b5da768b6873d7ba4129fc6ee54e62ba9c593000000203c7f7e5b08d1733d430658431545c9a2f57e4641b3b4cd52e567f27be9485e601c000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510a3f0e5be02320962616e64636861696e000000205a2f66b4d62d905b98277cd2807a324f0651340e80ae0249e500beb5ddcdce11000000203c1ed3d960b19e0ca7d321215874c6e91407ae3d2748f2e3b617fad833c30b6d1b000000106e080211840500000000000022480a200000003f12240a20b8aac9c5f107c71eacda8ccddd2506f30ecfa75685e12e403ddcc6411f6822ff10012a0c088d92c7f70510a394a6bf02320962616e64636861696e000000630000000000000000000000010000000f0000000342544300000000000003e8000000000000000400000000000000040000000000000000000000020000000000000004000000005ef1c903000000005ef1c907000000010000000800000000009269b300000000000005810000014800000006010100000000000000020000000000000581000000208d09a67a6aebb1498f3d104ae86c23d0c8e2e86dcb0a96f1088aca50b6d5e4030102000000000000000300000000000005810000002021c6cca538be5e9673479ba5694886335cdbbd0fcccd8d62376728c1133f00270103000000000000000500000000000005810000002057c9a0c6d2089b3911dc3bfd173ccd88922ed48030686ba3b3f8e4325b5f3cac0104000000000000000c00000000000005810000002014dc21a48e352ec5d08b52473da16cc6ca515e9e94124ea68b84a011b5da400c010500000000000000170000000000000581000000200fcacf91a581cf180802122abb9d02a67590fe38c7aea5baaa53416e91bf30e201070000000000000043000000000000058300000020c0062b4ff613ecbcf16e37a4423961afebf0302a208702d67e583622ad4a3dda"
        )

        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay_and_verify")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        params = {}
        params["block_height"] = 1412
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_oracle_state")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(
            "ff2b888760ffa9efa09e35763d8a5ba3ec794d00bcd1f9a901fa007f1f16e2eb", response.hex()
        )

    def test_relay_success_1(self):
        request_key = PyObi(
            """
            {
                client_id: string,
                oracle_script_id: u64,
                calldata: bytes,
                ask_count: u64,
                min_count: u64
            }
            """
        ).encode(
            {
                "client_id": "test",
                "oracle_script_id": 1,
                "calldata": bytes.fromhex("000000034254430000000000000064"),
                "ask_count": 4,
                "min_count": 4,
            }
        )
        params = {}
        params["encoded_request"] = request_key
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_latest_response")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(None, response)

        params = {}
        params["proof"] = bytes.fromhex(
            "00000000000001be000000a0b0cc616f4e76eba2c1845fd2321ad881991a2990f17cd278eb5f6d17da032faf9459616c154a23567a7281fb577bbfa1aa8d36382c4f64362968202fd32a0b65cbadb1694a5152c2b03f4960e1229745b39d7c41b32dc54e0c207799ae471981b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f93d913423122f2237a9a046d7ebfd13c3a6f1df93a1517a9f0b60a492dbc369d8000000c032fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869eeea5d164b32b494765eb256279014366747db8941c23f23b2fb09a9bbbf721762c4fd636adf805490129912f7fc18a2766aabd9fd4807f2ff18ed4ccd23f2583004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d7f4be7e5a1eb872ad44103360ddc190410331280c42a54d829a5d752c796685d000001e40000000300000020b0eff36d1214c35f785d82a7947589dba3f7d6ae35dcb92b1b3aca1df441e1e7000000207be19dd7f54c0ab0726cc5074fe65a887c2f2cd736cdab1ade1b7bd5c378abe11b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f705108c9fa29b01320962616e64636861696e00000020fc6a4aae096d2606050269fde7ac02f58c23fe10dd219f8805696c8aa3b1bb65000000206c324293f560b082923f9f1f1cc0307a78633095b87e4f4047a57ad650fd971b1b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f70510f6a5fd9c01320962616e64636861696e00000020d039e1d13e1837731a5719626075f1ef92ff48c6fbbec68b2ef9b3b34a6086db00000020278405f72ca62f7ff67c28b896ca38ab0534974d96a24adcefc509e0525301871b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f70510d3d5dd9b01320962616e64636861696e0000006b000000047465737400000000000000010000000f00000003425443000000000000006400000000000000040000000000000004000000047465737400000000000000010000000000000004000000005ef1c3c2000000005ef1c3c9000000010000000800000000000eaae600000000000001bc00000148000000060101000000000000000200000000000001bc000000201d3927941aec0602e1fcab94b4d3e57da16a2455165cb097adfdfdf34a3e16880102000000000000000400000000000001bc0000002017480c9a0faca2c2ecef5fdc7627b11f898f5652f3022e8d0b2da40356d6d8e40103000000000000000700000000000001bc00000020d90fcf3847ccd9ec012100f251f6e258d1e6b4e5e03a94fb06927dbb11f386230104000000000000000d00000000000001bc000000205332e6182a047d2b63df9c4a5edbb34d0760a5eea2816bfcc50ef0a316267a490105000000000000001200000000000001bc00000020fdd7908f23ba8fbe99c7f91c660a94a96231e807566e89fe524e739b1f6252c40107000000000000003800000000000001bd00000020efdf574e7d863cf4d8b8b08ed1558a6e6f47d5091667f29b7b2d37dc44d5f156"
        )

        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        params = {}
        params["encoded_request"] = request_key
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_latest_response")
            .params(params)
            .build()
        )
        res = self.process_call(call, self.icon_service)
        self.assertEqual("test", res["client_id"])
        self.assertEqual(1, res["request_id"])
        self.assertEqual(4, res["ans_count"])
        self.assertEqual(1592902594, res["request_time"])
        self.assertEqual(1592902601, res["resolve_time"])
        self.assertEqual(1, res["resolve_status"])
        self.assertEqual(b"\x00\x00\x00\x00\x00\x0e\xaa\xe6", res["result"])

    def test_relay_success_2(self):
        request_key = PyObi(
            """
            {
                client_id: string,
                oracle_script_id: u64,
                calldata: bytes,
                ask_count: u64,
                min_count: u64
            }
            """
        ).encode(
            {
                "client_id": "from_scan",
                "oracle_script_id": 1,
                "calldata": bytes.fromhex("000000034254430000000000000064"),
                "ask_count": 4,
                "min_count": 4,
            }
        )
        params = {}
        params["encoded_request"] = request_key
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_latest_response")
            .params(params)
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(None, response)

        params = {}
        params["proof"] = bytes.fromhex(
            "0000000000000a0e000000a05b14e631db5cdbb2c068be853242e45a091e9c178a5ca180656a12f6afde184472628869fbd308d01fcbd3411df82b10329b57df6c0c44626024e1bb0465b9486fdb04f6871c5d4bd5cb74a4c5a13a537bdd1159f7c38e797d297bc5f19770fbb1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f97e9cc402179a6714d7a4e34d8ddd21e853780a131455be3cdf89760b972c4b25000000c032fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869ef5529da80f2d395d551db9bce5ae15f89dd95314ca25124ec2456cf273c56fccd44c55c16e5d04367ddc30ca19f062d7bcd995c66a9f7b7437044391946ff1aa004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d7f4be7e5a1eb872ad44103360ddc190410331280c42a54d829a5d752c796685d000001e400000003000000205625b835f405a26be5befe6f5f84e3078e8c4776ad18f98e746878b6770749ed0000002044096ce96669a26351476f69665ba7ecc86372625a51fe5e6c49a8bd509169f71c000000106e0802110e0a00000000000022480a200000003f12240a20bc2ee65d649f9ac73b1c9da9af321df3058202a8873f4547bb1558ac6e39201810012a0c08f295c8f7051090ab958701320962616e64636861696e000000204add47680935a19b6ec9cedd40423b6729bc989de76cfc6debe51bb5be9735c800000020278a48fb4cd26f8b3735646a6bd4d234c459b8e6e1543b21ef324d0d3f80999c1c000000106e0802110e0a00000000000022480a200000003f12240a20bc2ee65d649f9ac73b1c9da9af321df3058202a8873f4547bb1558ac6e39201810012a0c08f295c8f70510b5c9938301320962616e64636861696e000000203892a4f6e094129190e9a6626de612731333b1efd797efe3325903e6350862400000002073f41504a93434ac2ca26e93d5a32ab40edbb870711992f20621b7d6a2f543351c000000106e0802110e0a00000000000022480a200000003f12240a20bc2ee65d649f9ac73b1c9da9af321df3058202a8873f4547bb1558ac6e39201810012a0c08f295c8f70510c2edf88301320962616e64636861696e000000750000000966726f6d5f7363616e00000000000000010000000f000000034254430000000000000064000000000000000400000000000000040000000966726f6d5f7363616e00000000000000030000000000000004000000005ef20816000000005ef2081a000000010000000800000000000eb71b00000000000007ea0000017e000000070001000000000000000200000000000009b2000000207e34da46bfd69cee4ffcbb512395fa7fa5fc7d5016aabe242f8e99cb82e806af0102000000000000000400000000000009b200000020bbd0ad34229408383dc3e8f4eaa1c44481ae7198a8ad5dfa234ed0559f09183c0103000000000000000600000000000009c300000020d15486b84f62b5535c0fd46c437489d212911d95501211f1c3a730e32f630eb60104000000000000000d00000000000009c3000000203cdb5cc721bc11d70d6eb33377a26c387ef3c45c064916d6a36659402252ff560105000000000000001b00000000000009c300000020121889d461b385a0d77fa5e484f58ad580f126d307526e7ef487c6742a8da9d90106000000000000003200000000000009c300000020562e86bb47132ccc0223eea4c35ef16d0d512086a975f1b473706dd23a571354010700000000000000500000000000000a0d000000200920a260853246a4ffb39a139a1ac98e3590879f521e03514a38448107996c50"
        )

        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        params = {}
        params["encoded_request"] = request_key
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_latest_response")
            .params(params)
            .build()
        )
        res = self.process_call(call, self.icon_service)
        self.assertEqual("from_scan", res["client_id"])
        self.assertEqual(3, res["request_id"])
        self.assertEqual(4, res["ans_count"])
        self.assertEqual(1592920086, res["request_time"])
        self.assertEqual(1592920090, res["resolve_time"])
        self.assertEqual(1, res["resolve_status"])
        self.assertEqual(b"\x00\x00\x00\x00\x00\x0e\xb7\x1b", res["result"])

    def test_relay_and_then_update_success(self):
        request_key = PyObi(
            """
            {
                client_id: string,
                oracle_script_id: u64,
                calldata: bytes,
                ask_count: u64,
                min_count: u64
            }
            """
        ).encode(
            {
                "client_id": "from_scan",
                "oracle_script_id": 1,
                "calldata": bytes.fromhex("000000034254430000000000000064"),
                "ask_count": 4,
                "min_count": 4,
            }
        )

        # relay round 1
        self.test_relay_success_2()

        params = {}
        params["encoded_request"] = request_key
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_latest_response")
            .params(params)
            .build()
        )
        res = self.process_call(call, self.icon_service)
        self.assertEqual("from_scan", res["client_id"])
        self.assertEqual(3, res["request_id"])
        self.assertEqual(4, res["ans_count"])
        self.assertEqual(1592920086, res["request_time"])
        self.assertEqual(1592920090, res["resolve_time"])
        self.assertEqual(1, res["resolve_status"])
        self.assertEqual(b"\x00\x00\x00\x00\x00\x0e\xb7\x1b", res["result"])

        params = {}
        params["proof"] = bytes.fromhex(
            "0000000000000a4f000000a02fb6b7c75032af08bb37d3df59722d54f400bec1f9f14bf840e7711cfc95e12b8e5deda960b2dc8f53406382441e13fc0821322c2b3dae3d8bbb5a4c6326ef087c80a176184bffaf33dcb3618c15e3371727c5b047598647ca3f9aecc4bd305cb1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f995e24c5fd917b4277cfac3e06fde09ae89231b58b1db51dc69f57b98ec18d9e7000000c032fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e4c68ced6b571ac169f65570fd568a2bf2821ca8880ec1d7b1b8052beb3bb4adf15798391d1c7d642e97b1a8fd91c11ab50cb91aa780a859fc5d62ca96c338602004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01dd991da4d4e69473cc75a4b819f9e07d4956671a6f4a74df4cc16596fcbe68137000001e400000003000000201297951f2c4fb97690848ce759ce1ab8ad8d972a951a869818828663159a8a6f000000200871114919ac2b6882b26da7728dd52d150ade92c4f6bca45791f005eb1bbbfe1c000000106e0802114f0a00000000000022480a200000003f12240a203a69db2cf26a980e8f51b3fcb05662eb72389ba98b5a6eec7ab7d51d18f8f04d10012a0c08c896c8f70510b39089e702320962616e64636861696e00000020d9efd0734b37e70e917d92c83e9b3ec9ab1cd3b1b5c91147310588bae960d25c000000202c85bb019f10afb9d2062ff880f0fca633a31e9e191f17c33006233281b131851c000000106e0802114f0a00000000000022480a200000003f12240a203a69db2cf26a980e8f51b3fcb05662eb72389ba98b5a6eec7ab7d51d18f8f04d10012a0c08c896c8f7051086f3c0e602320962616e64636861696e0000002026be63ade11a107dbef75d5e506f1c64173d8bf1d9508bba19b74d7da222df56000000201d1a861a5541af331fb13a5def77e5c3fbfe70378b134292f5380854817a9e9c1b000000106e0802114f0a00000000000022480a200000003f12240a203a69db2cf26a980e8f51b3fcb05662eb72389ba98b5a6eec7ab7d51d18f8f04d10012a0c08c896c8f7051094b592e602320962616e64636861696e000000750000000966726f6d5f7363616e00000000000000010000000f000000034254430000000000000064000000000000000400000000000000040000000966726f6d5f7363616e00000000000000040000000000000004000000005ef20a73000000005ef20a77000000010000000800000000000eb4bf00000000000009b20000017e000000070101000000000000000200000000000009b200000020680ad6cf6e554cee42ff66ecbd120e66009b10284cabc5e02241552c818aae7d0102000000000000000400000000000009b200000020bbd0ad34229408383dc3e8f4eaa1c44481ae7198a8ad5dfa234ed0559f09183c0103000000000000000600000000000009c300000020d15486b84f62b5535c0fd46c437489d212911d95501211f1c3a730e32f630eb60104000000000000000d00000000000009c3000000203cdb5cc721bc11d70d6eb33377a26c387ef3c45c064916d6a36659402252ff560105000000000000001b00000000000009c300000020121889d461b385a0d77fa5e484f58ad580f126d307526e7ef487c6742a8da9d90106000000000000003200000000000009c300000020562e86bb47132ccc0223eea4c35ef16d0d512086a975f1b473706dd23a571354010700000000000000500000000000000a4e00000020240608de388ec1d1100d05d87f2766cc22f98376f955791546074aa8ae6b6d8e"
        )

        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        params = {}
        params["encoded_request"] = request_key
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_latest_response")
            .params(params)
            .build()
        )
        res = self.process_call(call, self.icon_service)
        self.assertEqual("from_scan", res["client_id"])
        self.assertEqual(4, res["request_id"])
        self.assertEqual(4, res["ans_count"])
        self.assertEqual(1592920691, res["request_time"])
        self.assertEqual(1592920695, res["resolve_time"])
        self.assertEqual(1, res["resolve_status"])
        self.assertEqual(b"\x00\x00\x00\x00\x00\x0e\xb4\xbf", res["result"])

    def test_relay_and_then_update_fail_because_update_with_older_state(self):
        request_key = PyObi(
            """
            {
                client_id: string,
                oracle_script_id: u64,
                calldata: bytes,
                ask_count: u64,
                min_count: u64
            }
            """
        ).encode(
            {
                "client_id": "from_scan",
                "oracle_script_id": 1,
                "calldata": bytes.fromhex("000000034254430000000000000064"),
                "ask_count": 4,
                "min_count": 4,
            }
        )

        params = {}
        params["proof"] = bytes.fromhex(
            "0000000000000a4f000000a02fb6b7c75032af08bb37d3df59722d54f400bec1f9f14bf840e7711cfc95e12b8e5deda960b2dc8f53406382441e13fc0821322c2b3dae3d8bbb5a4c6326ef087c80a176184bffaf33dcb3618c15e3371727c5b047598647ca3f9aecc4bd305cb1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f995e24c5fd917b4277cfac3e06fde09ae89231b58b1db51dc69f57b98ec18d9e7000000c032fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e4c68ced6b571ac169f65570fd568a2bf2821ca8880ec1d7b1b8052beb3bb4adf15798391d1c7d642e97b1a8fd91c11ab50cb91aa780a859fc5d62ca96c338602004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01dd991da4d4e69473cc75a4b819f9e07d4956671a6f4a74df4cc16596fcbe68137000001e400000003000000201297951f2c4fb97690848ce759ce1ab8ad8d972a951a869818828663159a8a6f000000200871114919ac2b6882b26da7728dd52d150ade92c4f6bca45791f005eb1bbbfe1c000000106e0802114f0a00000000000022480a200000003f12240a203a69db2cf26a980e8f51b3fcb05662eb72389ba98b5a6eec7ab7d51d18f8f04d10012a0c08c896c8f70510b39089e702320962616e64636861696e00000020d9efd0734b37e70e917d92c83e9b3ec9ab1cd3b1b5c91147310588bae960d25c000000202c85bb019f10afb9d2062ff880f0fca633a31e9e191f17c33006233281b131851c000000106e0802114f0a00000000000022480a200000003f12240a203a69db2cf26a980e8f51b3fcb05662eb72389ba98b5a6eec7ab7d51d18f8f04d10012a0c08c896c8f7051086f3c0e602320962616e64636861696e0000002026be63ade11a107dbef75d5e506f1c64173d8bf1d9508bba19b74d7da222df56000000201d1a861a5541af331fb13a5def77e5c3fbfe70378b134292f5380854817a9e9c1b000000106e0802114f0a00000000000022480a200000003f12240a203a69db2cf26a980e8f51b3fcb05662eb72389ba98b5a6eec7ab7d51d18f8f04d10012a0c08c896c8f7051094b592e602320962616e64636861696e000000750000000966726f6d5f7363616e00000000000000010000000f000000034254430000000000000064000000000000000400000000000000040000000966726f6d5f7363616e00000000000000040000000000000004000000005ef20a73000000005ef20a77000000010000000800000000000eb4bf00000000000009b20000017e000000070101000000000000000200000000000009b200000020680ad6cf6e554cee42ff66ecbd120e66009b10284cabc5e02241552c818aae7d0102000000000000000400000000000009b200000020bbd0ad34229408383dc3e8f4eaa1c44481ae7198a8ad5dfa234ed0559f09183c0103000000000000000600000000000009c300000020d15486b84f62b5535c0fd46c437489d212911d95501211f1c3a730e32f630eb60104000000000000000d00000000000009c3000000203cdb5cc721bc11d70d6eb33377a26c387ef3c45c064916d6a36659402252ff560105000000000000001b00000000000009c300000020121889d461b385a0d77fa5e484f58ad580f126d307526e7ef487c6742a8da9d90106000000000000003200000000000009c300000020562e86bb47132ccc0223eea4c35ef16d0d512086a975f1b473706dd23a571354010700000000000000500000000000000a4e00000020240608de388ec1d1100d05d87f2766cc22f98376f955791546074aa8ae6b6d8e"
        )

        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        params = {}
        params["encoded_request"] = request_key
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .method("get_latest_response")
            .params(params)
            .build()
        )
        res = self.process_call(call, self.icon_service)
        self.assertEqual("from_scan", res["client_id"])
        self.assertEqual(4, res["request_id"])
        self.assertEqual(4, res["ans_count"])
        self.assertEqual(1592920691, res["request_time"])
        self.assertEqual(1592920695, res["resolve_time"])
        self.assertEqual(1, res["resolve_status"])
        self.assertEqual(b"\x00\x00\x00\x00\x00\x0e\xb4\xbf", res["result"])

        params = {}
        params["proof"] = bytes.fromhex(
            "0000000000000a0e000000a05b14e631db5cdbb2c068be853242e45a091e9c178a5ca180656a12f6afde184472628869fbd308d01fcbd3411df82b10329b57df6c0c44626024e1bb0465b9486fdb04f6871c5d4bd5cb74a4c5a13a537bdd1159f7c38e797d297bc5f19770fbb1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f97e9cc402179a6714d7a4e34d8ddd21e853780a131455be3cdf89760b972c4b25000000c032fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869ef5529da80f2d395d551db9bce5ae15f89dd95314ca25124ec2456cf273c56fccd44c55c16e5d04367ddc30ca19f062d7bcd995c66a9f7b7437044391946ff1aa004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d7f4be7e5a1eb872ad44103360ddc190410331280c42a54d829a5d752c796685d000001e400000003000000205625b835f405a26be5befe6f5f84e3078e8c4776ad18f98e746878b6770749ed0000002044096ce96669a26351476f69665ba7ecc86372625a51fe5e6c49a8bd509169f71c000000106e0802110e0a00000000000022480a200000003f12240a20bc2ee65d649f9ac73b1c9da9af321df3058202a8873f4547bb1558ac6e39201810012a0c08f295c8f7051090ab958701320962616e64636861696e000000204add47680935a19b6ec9cedd40423b6729bc989de76cfc6debe51bb5be9735c800000020278a48fb4cd26f8b3735646a6bd4d234c459b8e6e1543b21ef324d0d3f80999c1c000000106e0802110e0a00000000000022480a200000003f12240a20bc2ee65d649f9ac73b1c9da9af321df3058202a8873f4547bb1558ac6e39201810012a0c08f295c8f70510b5c9938301320962616e64636861696e000000203892a4f6e094129190e9a6626de612731333b1efd797efe3325903e6350862400000002073f41504a93434ac2ca26e93d5a32ab40edbb870711992f20621b7d6a2f543351c000000106e0802110e0a00000000000022480a200000003f12240a20bc2ee65d649f9ac73b1c9da9af321df3058202a8873f4547bb1558ac6e39201810012a0c08f295c8f70510c2edf88301320962616e64636861696e000000750000000966726f6d5f7363616e00000000000000010000000f000000034254430000000000000064000000000000000400000000000000040000000966726f6d5f7363616e00000000000000030000000000000004000000005ef20816000000005ef2081a000000010000000800000000000eb71b00000000000007ea0000017e000000070001000000000000000200000000000009b2000000207e34da46bfd69cee4ffcbb512395fa7fa5fc7d5016aabe242f8e99cb82e806af0102000000000000000400000000000009b200000020bbd0ad34229408383dc3e8f4eaa1c44481ae7198a8ad5dfa234ed0559f09183c0103000000000000000600000000000009c300000020d15486b84f62b5535c0fd46c437489d212911d95501211f1c3a730e32f630eb60104000000000000000d00000000000009c3000000203cdb5cc721bc11d70d6eb33377a26c387ef3c45c064916d6a36659402252ff560105000000000000001b00000000000009c300000020121889d461b385a0d77fa5e484f58ad580f126d307526e7ef487c6742a8da9d90106000000000000003200000000000009c300000020562e86bb47132ccc0223eea4c35ef16d0d512086a975f1b473706dd23a571354010700000000000000500000000000000a0d000000200920a260853246a4ffb39a139a1ac98e3590879f521e03514a38448107996c50"
        )

        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(False, tx_result["status"])

    def test_relay_and_then_update_fail_because_submit_an_unresolved_request(self):
        request_key = PyObi(
            """
            {
                client_id: string,
                oracle_script_id: u64,
                calldata: bytes,
                ask_count: u64,
                min_count: u64
            }
            """
        ).encode(
            {
                "client_id": "from_scan",
                "oracle_script_id": 13,
                "calldata": bytes.fromhex(
                    "0000000342544300000003555344000000046d65616e0000000000000064"
                ),
                "ask_count": 4,
                "min_count": 4,
            }
        )

        params = {}
        params["proof"] = bytes.fromhex(
            "0000000000019e8c000000a010d1bb5a23a0ed3668df0b2a4257b46225db7391c56e940ad95536d0e5a75d79f4be0eeff9a2e52e11da7591c029cee08ecfe59435dbd731679905a94c2bb53b8fdacf87d1d46d85ee2585ca5a304a2a499080b8c7460f5efac7d053911fbee6b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f9df898ce21b5981ca9994759bb2b77204496e9d062a080b52c2b38b8d8df0f552000000c032fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e8f7beb334b2ff0c5c176144edbbfe9426cdb41fea466d43d0b89e2515f9e5a5170e961a346803fa2198490f4ded2683e31384b2a37725e719ae7ca164eb22ceb004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d0efe3e12f46363c7779140d4ce659925db52f19053e114d7cc4efd666b37f79f000001e400000003000000204ed98be18e8daa17dface914d1a3cc1dfa9f750e49288eb0f174e5807d789cf400000020552ea28fd59148a5885bdc2d703e2899754a7ebb08e6599fdfc43845c714c8671b000000106e0802118c9e01000000000022480a200000003f12240a20ee95cbd137a9b459e382b317a6342087da4ddd9e4215fbbb69cec3fe4b82028b10012a0c08b3c2d0f70510e3f7da9c01320962616e64636861696e0000002074e9db581a49405c8982e3ecaa81e5ef0a7e2ebbfa72c758778c3070380de2da000000206ace14389fd984502d1ab1273d3f678dbb5a988ceeb4912aa9e93e776ec0c1101b000000106e0802118c9e01000000000022480a200000003f12240a20ee95cbd137a9b459e382b317a6342087da4ddd9e4215fbbb69cec3fe4b82028b10012a0c08b3c2d0f70510cfd6959e01320962616e64636861696e000000207a45a72854b0a9a6a615c55efafc7daa6366aee86700e15f1f45e82c5644e4dc000000207545bdd1dd0d6b5efadd5816b400be55b61f5b9c7127be2521009f43b092112d1b000000106e0802118c9e01000000000022480a200000003f12240a20ee95cbd137a9b459e382b317a6342087da4ddd9e4215fbbb69cec3fe4b82028b10012a0c08b3c2d0f70510d9fa9da001320962616e64636861696e0000007c0000000966726f6d5f7363616e000000000000000d0000001e0000000342544300000003555344000000046d65616e0000000000000064000000000000000400000000000000040000000966726f6d5f7363616e00000000000000090000000000000004000000005ef34682000000005ef346860000000200000000000000000000f8d60000017e0000000700020000000000000003000000000000f93700000020a181077d74a00d1c4888d3957b673b7ac5265cef6642b398ee81564e0a02a86601030000000000000005000000000000f9370000002092dbf6c1d5523fedfec04d0d8cf33992bc2f5a748d935cafaa022d14da1170190104000000000000000d000000000000f94800000020bea043e9a22ced8e9149097d6f7c9c3d719280d05fcf127cb7727fdb8a5c4f7f01050000000000000014000000000000f948000000203be0e603fd2d50e9adff89e384c1586933110ef6320849658961f119a6fb731e01060000000000000022000000000000f94800000020121889d461b385a0d77fa5e484f58ad580f126d307526e7ef487c6742a8da9d901070000000000000039000000000000f94800000020562e86bb47132ccc0223eea4c35ef16d0d512086a975f1b473706dd23a5713540108000000000000007a0000000000019e8b00000020cc921a62033387ad80d703b9fccf04e1d501b4e7d8e006a79d9d91de032092f3"
        )

        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(False, tx_result["status"])

    def test_relay_and_verify_via_receiver_mock_success(self):
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._receiver_mock_address)
            .method("get_req")
            .build()
        )
        req = self.process_call(call, self.icon_service)
        self.assertEqual(None, req)

        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._receiver_mock_address)
            .method("get_res")
            .build()
        )
        res = self.process_call(call, self.icon_service)
        self.assertEqual(None, res)

        params = {}
        params["proof"] = bytes.fromhex(
            "00000000000001be000000a0b0cc616f4e76eba2c1845fd2321ad881991a2990f17cd278eb5f6d17da032faf9459616c154a23567a7281fb577bbfa1aa8d36382c4f64362968202fd32a0b65cbadb1694a5152c2b03f4960e1229745b39d7c41b32dc54e0c207799ae471981b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f93d913423122f2237a9a046d7ebfd13c3a6f1df93a1517a9f0b60a492dbc369d8000000c032fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869eeea5d164b32b494765eb256279014366747db8941c23f23b2fb09a9bbbf721762c4fd636adf805490129912f7fc18a2766aabd9fd4807f2ff18ed4ccd23f2583004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d7f4be7e5a1eb872ad44103360ddc190410331280c42a54d829a5d752c796685d000001e40000000300000020b0eff36d1214c35f785d82a7947589dba3f7d6ae35dcb92b1b3aca1df441e1e7000000207be19dd7f54c0ab0726cc5074fe65a887c2f2cd736cdab1ade1b7bd5c378abe11b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f705108c9fa29b01320962616e64636861696e00000020fc6a4aae096d2606050269fde7ac02f58c23fe10dd219f8805696c8aa3b1bb65000000206c324293f560b082923f9f1f1cc0307a78633095b87e4f4047a57ad650fd971b1b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f70510f6a5fd9c01320962616e64636861696e00000020d039e1d13e1837731a5719626075f1ef92ff48c6fbbec68b2ef9b3b34a6086db00000020278405f72ca62f7ff67c28b896ca38ab0534974d96a24adcefc509e0525301871b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f70510d3d5dd9b01320962616e64636861696e0000006b000000047465737400000000000000010000000f00000003425443000000000000006400000000000000040000000000000004000000047465737400000000000000010000000000000004000000005ef1c3c2000000005ef1c3c9000000010000000800000000000eaae600000000000001bc00000148000000060101000000000000000200000000000001bc000000201d3927941aec0602e1fcab94b4d3e57da16a2455165cb097adfdfdf34a3e16880102000000000000000400000000000001bc0000002017480c9a0faca2c2ecef5fdc7627b11f898f5652f3022e8d0b2da40356d6d8e40103000000000000000700000000000001bc00000020d90fcf3847ccd9ec012100f251f6e258d1e6b4e5e03a94fb06927dbb11f386230104000000000000000d00000000000001bc000000205332e6182a047d2b63df9c4a5edbb34d0760a5eea2816bfcc50ef0a316267a490105000000000000001200000000000001bc00000020fdd7908f23ba8fbe99c7f91c660a94a96231e807566e89fe524e739b1f6252c40107000000000000003800000000000001bd00000020efdf574e7d863cf4d8b8b08ed1558a6e6f47d5091667f29b7b2d37dc44d5f156"
        )
        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._receiver_mock_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay_and_safe")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._receiver_mock_address)
            .method("get_req")
            .build()
        )
        req = self.process_call(call, self.icon_service)
        self.assertEqual("test", req["client_id"])
        self.assertEqual(1, req["oracle_script_id"])
        self.assertEqual(bytes.fromhex("000000034254430000000000000064"), req["calldata"])
        self.assertEqual(4, req["ask_count"])
        self.assertEqual(4, req["min_count"])

        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._receiver_mock_address)
            .method("get_res")
            .build()
        )
        res = self.process_call(call, self.icon_service)
        self.assertEqual("test", res["client_id"])
        self.assertEqual(1, res["request_id"])
        self.assertEqual(4, res["ans_count"])
        self.assertEqual(1592902594, res["request_time"])
        self.assertEqual(1592902601, res["resolve_time"])
        self.assertEqual(1, res["resolve_status"])
        self.assertEqual(b"\x00\x00\x00\x00\x00\x0e\xaa\xe6", res["result"])

    def test_consume_data_fail_because_no_data(self):
        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._cache_consumer_mock)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("consume_cache")
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(False, tx_result["status"])

    def test_relay_and_consume_data_by_consumer_success(self):
        # should get None because haven't relay and comsume yet
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._cache_consumer_mock)
            .method("get_res")
            .build()
        )
        res = self.process_call(call, self.icon_service)
        self.assertEqual(None, res)

        # relay
        params = {}
        params["proof"] = bytes.fromhex(
            "0000000000000a4f000000a02fb6b7c75032af08bb37d3df59722d54f400bec1f9f14bf840e7711cfc95e12b8e5deda960b2dc8f53406382441e13fc0821322c2b3dae3d8bbb5a4c6326ef087c80a176184bffaf33dcb3618c15e3371727c5b047598647ca3f9aecc4bd305cb1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f995e24c5fd917b4277cfac3e06fde09ae89231b58b1db51dc69f57b98ec18d9e7000000c032fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e4c68ced6b571ac169f65570fd568a2bf2821ca8880ec1d7b1b8052beb3bb4adf15798391d1c7d642e97b1a8fd91c11ab50cb91aa780a859fc5d62ca96c338602004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01dd991da4d4e69473cc75a4b819f9e07d4956671a6f4a74df4cc16596fcbe68137000001e400000003000000201297951f2c4fb97690848ce759ce1ab8ad8d972a951a869818828663159a8a6f000000200871114919ac2b6882b26da7728dd52d150ade92c4f6bca45791f005eb1bbbfe1c000000106e0802114f0a00000000000022480a200000003f12240a203a69db2cf26a980e8f51b3fcb05662eb72389ba98b5a6eec7ab7d51d18f8f04d10012a0c08c896c8f70510b39089e702320962616e64636861696e00000020d9efd0734b37e70e917d92c83e9b3ec9ab1cd3b1b5c91147310588bae960d25c000000202c85bb019f10afb9d2062ff880f0fca633a31e9e191f17c33006233281b131851c000000106e0802114f0a00000000000022480a200000003f12240a203a69db2cf26a980e8f51b3fcb05662eb72389ba98b5a6eec7ab7d51d18f8f04d10012a0c08c896c8f7051086f3c0e602320962616e64636861696e0000002026be63ade11a107dbef75d5e506f1c64173d8bf1d9508bba19b74d7da222df56000000201d1a861a5541af331fb13a5def77e5c3fbfe70378b134292f5380854817a9e9c1b000000106e0802114f0a00000000000022480a200000003f12240a203a69db2cf26a980e8f51b3fcb05662eb72389ba98b5a6eec7ab7d51d18f8f04d10012a0c08c896c8f7051094b592e602320962616e64636861696e000000750000000966726f6d5f7363616e00000000000000010000000f000000034254430000000000000064000000000000000400000000000000040000000966726f6d5f7363616e00000000000000040000000000000004000000005ef20a73000000005ef20a77000000010000000800000000000eb4bf00000000000009b20000017e000000070101000000000000000200000000000009b200000020680ad6cf6e554cee42ff66ecbd120e66009b10284cabc5e02241552c818aae7d0102000000000000000400000000000009b200000020bbd0ad34229408383dc3e8f4eaa1c44481ae7198a8ad5dfa234ed0559f09183c0103000000000000000600000000000009c300000020d15486b84f62b5535c0fd46c437489d212911d95501211f1c3a730e32f630eb60104000000000000000d00000000000009c3000000203cdb5cc721bc11d70d6eb33377a26c387ef3c45c064916d6a36659402252ff560105000000000000001b00000000000009c300000020121889d461b385a0d77fa5e484f58ad580f126d307526e7ef487c6742a8da9d90106000000000000003200000000000009c300000020562e86bb47132ccc0223eea4c35ef16d0d512086a975f1b473706dd23a571354010700000000000000500000000000000a4e00000020240608de388ec1d1100d05d87f2766cc22f98376f955791546074aa8ae6b6d8e"
        )

        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._bridge_address)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay")
            .params(params)
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        # should get None because haven't comsume yet
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._cache_consumer_mock)
            .method("get_res")
            .build()
        )
        res = self.process_call(call, self.icon_service)
        self.assertEqual(None, res)

        # consume cache
        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._cache_consumer_mock)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("consume_cache")
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        # should get a correct res
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._cache_consumer_mock)
            .method("get_res")
            .build()
        )
        res = self.process_call(call, self.icon_service)
        self.assertEqual(
            {
                "client_id": "from_scan",
                "request_id": 4,
                "ans_count": 4,
                "request_time": 1592920691,
                "resolve_time": 1592920695,
                "resolve_status": 1,
                "result": b"\x00\x00\x00\x00\x00\x0e\xb4\xbf",
            },
            res,
        )

