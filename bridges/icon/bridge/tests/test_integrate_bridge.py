import os

from iconsdk.builder.call_builder import CallBuilder
from iconsdk.builder.transaction_builder import (
    DeployTransactionBuilder,
    CallTransactionBuilder,
)
from iconsdk.wallet.wallet import KeyWallet
from iconsdk.libs.in_memory_zip import gen_deploy_data_content
from iconsdk.signed_transaction import SignedTransaction
from tbears.libs.icon_integrate_test import IconIntegrateTestBase, SCORE_INSTALL_ADDRESS

DIR_PATH = os.path.abspath(os.path.dirname(__file__))


class TestTest(IconIntegrateTestBase):
    TEST_HTTP_ENDPOINT_URI_V3 = "http://127.0.0.1:9000/api/v3"
    SCORE_PROJECT = os.path.abspath(os.path.join(DIR_PATH, '..'))

    def setUp(self):
        super().setUp()

        self.icon_service = None
        # if you want to send request to network, uncomment next line and set self.TEST_HTTP_ENDPOINT_URI_V3
        # self.icon_service = IconService(HTTPProvider(self.TEST_HTTP_ENDPOINT_URI_V3))

        # install SCORE
        self._score_address = self._deploy_score()['scoreAddress']

    def _deploy_score(self, to: str = SCORE_INSTALL_ADDRESS) -> dict:
        params = {}
        params["validators_bytes"] = "0000000400000040a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba331db530b76775beb0348c52bb8fc1fdac207525e5689caa01c0af8d2f8f371ec5000000000000006400000040724ae29cfeb7497051d09edfd8e822352c4c8361b757647645b78c8cc74ce885f04c26ee07ff6ada08587a4037363838b1dda6e306091ee0690caa8fe0e6fcd2000000000000006400000040f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479000000000000006400000040d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f2c67aaeee2bc8d00ff253a326eb04de8ce88eea1224c2bb8ca69296a1c753f530000000000000064"

        # Generates an instance of transaction for deploying SCORE.
        transaction = DeployTransactionBuilder().from_(self._test1.get_address()).to(to).params(params).step_limit(100_000_000_000).nid(
            3).nonce(100).content_type("application/zip").content(gen_deploy_data_content(self.SCORE_PROJECT)).build()

        # Returns the signed transaction object having a signature
        signed_transaction = SignedTransaction(transaction, self._test1)

        # process the transaction in local
        tx_result = self.process_transaction(
            signed_transaction, self.icon_service)

        self.assertEqual(True, tx_result['status'])
        self.assertTrue('scoreAddress' in tx_result)

        return tx_result

    def test_update_validator_powers_success(self):
        params = {}
        params["pub_key"] = "f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479"
        call = CallBuilder().from_(self._test1.get_address()).to(
            self._score_address).method("get_validator_power").params(params).build()
        response = self.process_call(call, self.icon_service)
        self.assertEqual("0x64", response)

        params = {}
        params["validators_bytes"] = "0000000100000040f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b704790000000000000096"

        transaction = CallTransactionBuilder().from_(self._test1.get_address()).to(
            self._score_address).step_limit(100_000_000_000).nid(
            3).nonce(100).method("update_validator_powers").params(params).build()
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(
            signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result['status'])

        params = {}
        params["pub_key"] = "f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479"
        call = CallBuilder().from_(self._test1.get_address()).to(
            self._score_address).method("get_validator_power").params(params).build()
        response = self.process_call(call, self.icon_service)
        self.assertEqual("0x96", response)

    def test_update_validator_powers_fail_not_authorized(self):
        new_owner = KeyWallet.create()
        params = {}
        params["pub_key"] = "f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479"
        call = CallBuilder().from_(new_owner.get_address()).to(
            self._score_address).method("get_validator_power").params(params).build()
        response = self.process_call(call, self.icon_service)
        self.assertEqual("0x64", response)

        params = {}
        params["validators_bytes"] = "0000000100000040f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b704790000000000000096"

        transaction = CallTransactionBuilder().from_(new_owner.get_address()).to(
            self._score_address).step_limit(100_000_000_000).nid(
            3).nonce(100).method("update_validator_powers").params(params).build()
        signed_transaction = SignedTransaction(transaction, new_owner)
        tx_result = self.process_transaction(
            signed_transaction, self.icon_service)
        self.assertEqual(False, tx_result['status'])

    def test_relay_oracle_state_1(self):
        params = {}
        params["block_height"] = 191
        call = CallBuilder().from_(self._test1.get_address()).to(
            self._score_address).method("get_oracle_state").params(params).build()
        response = self.process_call(call, self.icon_service)
        self.assertEqual(None, response)

        params = {}
        params["block_height"] = 191
        params["multi_store_bytes"] = "685430546d23a44e6b8034eaafbc2f4cd7fef54b54d5b66528cb4e5225bd74fb4f8d0bb0cd3eb9dc70b4dbbea5f0cbd5b523195f7bae02bb401bb00a93aba08e1912057fff0b3e85abf1a319f75d37b21b430f3da7db9e486a5041de47c686d3b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f991cc906286235b676ad402fc04ed768eb2bfeca664e8d595c286571da1433c60"
        params["merkle_part_bytes"] = "32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869ed9f175396c0e2d0e77f69856abf5d8e69283cb915eb8886262feda1d519b30054f4d548668a3986db253689234b9cac96303a128b7081b16e53b79ab9e65b887004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01dd991da4d4e69473cc75a4b819f9e07d4956671a6f4a74df4cc16596fcbe68137"
        params["signatures_bytes"] = "00000003000000209279257914bfee6faec46df086e9a673a4c1576fb299094d45f77e12fa3728e70000002004dbf23f5ebb07ba8fd7ca06843dfe363b5e86596930e1889d9bd5bd3e98fc531c000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510979fbd9801320962616e64636861696e00000020826bb17b714ebcd8199ee2a01334102f19248087cfdeee42ebd406b3991c389500000020621281eedf97f3a9ec121224cef9f8c07872d2c81cce44dd846bb3678dd505231b000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510dd89e98a01320962616e64636861696e00000020d775fd0e1580499ef16a4ab1998dcb4cbd47cf6f342cdc51a9552d1434552ed8000000205a1a2075ae97a6bbd07a5a40efe287eceab37c7a7caae82e6c997c51c62334701b000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510ab94a59201320962616e64636861696e"
        transaction = CallTransactionBuilder().from_(self._test1.get_address()).to(
            self._score_address).step_limit(100_000_000_000).nid(
            3).nonce(100).method("relay_oracle_state").params(params).build()
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(
            signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result['status'])

        params = {}
        params["block_height"] = 191
        call = CallBuilder().from_(self._test1.get_address()).to(
            self._score_address).method("get_oracle_state").params(params).build()
        response = self.process_call(call, self.icon_service)
        self.assertEqual(
            "1912057fff0b3e85abf1a319f75d37b21b430f3da7db9e486a5041de47c686d3", response.hex())

    def test_relay_oracle_state_2(self):
        params = {}
        params["block_height"] = 446
        call = CallBuilder().from_(self._test1.get_address()).to(
            self._score_address).method("get_oracle_state").params(params).build()
        response = self.process_call(call, self.icon_service)
        self.assertEqual(None, response)

        params = {}
        params["block_height"] = 446
        params["multi_store_bytes"] = "b0cc616f4e76eba2c1845fd2321ad881991a2990f17cd278eb5f6d17da032faf9459616c154a23567a7281fb577bbfa1aa8d36382c4f64362968202fd32a0b65cbadb1694a5152c2b03f4960e1229745b39d7c41b32dc54e0c207799ae471981b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f93d913423122f2237a9a046d7ebfd13c3a6f1df93a1517a9f0b60a492dbc369d8"
        params["merkle_part_bytes"] = "32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869eeea5d164b32b494765eb256279014366747db8941c23f23b2fb09a9bbbf721762c4fd636adf805490129912f7fc18a2766aabd9fd4807f2ff18ed4ccd23f2583004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d7f4be7e5a1eb872ad44103360ddc190410331280c42a54d829a5d752c796685d"
        params["signatures_bytes"] = "0000000300000020b0eff36d1214c35f785d82a7947589dba3f7d6ae35dcb92b1b3aca1df441e1e7000000207be19dd7f54c0ab0726cc5074fe65a887c2f2cd736cdab1ade1b7bd5c378abe11b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f705108c9fa29b01320962616e64636861696e00000020fc6a4aae096d2606050269fde7ac02f58c23fe10dd219f8805696c8aa3b1bb65000000206c324293f560b082923f9f1f1cc0307a78633095b87e4f4047a57ad650fd971b1b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f70510f6a5fd9c01320962616e64636861696e00000020d039e1d13e1837731a5719626075f1ef92ff48c6fbbec68b2ef9b3b34a6086db00000020278405f72ca62f7ff67c28b896ca38ab0534974d96a24adcefc509e0525301871b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f70510d3d5dd9b01320962616e64636861696e"
        transaction = CallTransactionBuilder().from_(self._test1.get_address()).to(
            self._score_address).step_limit(100_000_000_000).nid(
            3).nonce(100).method("relay_oracle_state").params(params).build()
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(
            signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result['status'])

        params = {}
        params["block_height"] = 446
        call = CallBuilder().from_(self._test1.get_address()).to(
            self._score_address).method("get_oracle_state").params(params).build()
        response = self.process_call(call, self.icon_service)
        self.assertEqual(
            "cbadb1694a5152c2b03f4960e1229745b39d7c41b32dc54e0c207799ae471981", response.hex())

    def test_relay_and_verify(self):
        params = {}
        params["block_height"] = 446
        call = CallBuilder().from_(self._test1.get_address()).to(
            self._score_address).method("get_oracle_state").params(params).build()
        response = self.process_call(call, self.icon_service)
        self.assertEqual(None, response)

        params = {}
        params["proof"] = bytes.fromhex("00000000000001be000000a0b0cc616f4e76eba2c1845fd2321ad881991a2990f17cd278eb5f6d17da032faf9459616c154a23567a7281fb577bbfa1aa8d36382c4f64362968202fd32a0b65cbadb1694a5152c2b03f4960e1229745b39d7c41b32dc54e0c207799ae471981b1f2fd852e790e735ca2d3014f96a2a53c60393e9c6bbf941b9a6dd6a05cf6f93d913423122f2237a9a046d7ebfd13c3a6f1df93a1517a9f0b60a492dbc369d8000000c032fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869eeea5d164b32b494765eb256279014366747db8941c23f23b2fb09a9bbbf721762c4fd636adf805490129912f7fc18a2766aabd9fd4807f2ff18ed4ccd23f2583004209a161040ab1778e2f2c00ee482f205b28efba439fcb04ea283f619478d96e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d7f4be7e5a1eb872ad44103360ddc190410331280c42a54d829a5d752c796685d000001e40000000300000020b0eff36d1214c35f785d82a7947589dba3f7d6ae35dcb92b1b3aca1df441e1e7000000207be19dd7f54c0ab0726cc5074fe65a887c2f2cd736cdab1ade1b7bd5c378abe11b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f705108c9fa29b01320962616e64636861696e00000020fc6a4aae096d2606050269fde7ac02f58c23fe10dd219f8805696c8aa3b1bb65000000206c324293f560b082923f9f1f1cc0307a78633095b87e4f4047a57ad650fd971b1b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f70510f6a5fd9c01320962616e64636861696e00000020d039e1d13e1837731a5719626075f1ef92ff48c6fbbec68b2ef9b3b34a6086db00000020278405f72ca62f7ff67c28b896ca38ab0534974d96a24adcefc509e0525301871b000000106e080211be0100000000000022480a200000003f12240a201e0d38ac4548746d650a022463529c1f64777038aefdd2d16bf358daa6e41b8010012a0c08ce87c7f70510d3d5dd9b01320962616e64636861696e0000006b000000047465737400000000000000010000000f00000003425443000000000000006400000000000000040000000000000004000000047465737400000000000000010000000000000004000000005ef1c3c2000000005ef1c3c9000000010000000800000000000eaae600000000000001bc00000148000000060101000000000000000200000000000001bc000000201d3927941aec0602e1fcab94b4d3e57da16a2455165cb097adfdfdf34a3e16880102000000000000000400000000000001bc0000002017480c9a0faca2c2ecef5fdc7627b11f898f5652f3022e8d0b2da40356d6d8e40103000000000000000700000000000001bc00000020d90fcf3847ccd9ec012100f251f6e258d1e6b4e5e03a94fb06927dbb11f386230104000000000000000d00000000000001bc000000205332e6182a047d2b63df9c4a5edbb34d0760a5eea2816bfcc50ef0a316267a490105000000000000001200000000000001bc00000020fdd7908f23ba8fbe99c7f91c660a94a96231e807566e89fe524e739b1f6252c40107000000000000003800000000000001bd00000020efdf574e7d863cf4d8b8b08ed1558a6e6f47d5091667f29b7b2d37dc44d5f156")
        transaction = CallTransactionBuilder().from_(self._test1.get_address()).to(
            self._score_address).step_limit(100_000_000_000).nid(
            3).nonce(100).method("relay_and_verify").params(params).build()
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(
            signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result['status'])

        params = {}
        params["block_height"] = 446
        call = CallBuilder().from_(self._test1.get_address()).to(
            self._score_address).method("get_oracle_state").params(params).build()
        response = self.process_call(call, self.icon_service)
        self.assertEqual(
            "cbadb1694a5152c2b03f4960e1229745b39d7c41b32dc54e0c207799ae471981", response.hex())
