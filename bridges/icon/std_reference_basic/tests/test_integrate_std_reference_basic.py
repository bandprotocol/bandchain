import os
import json
from iconsdk.builder.call_builder import CallBuilder
from iconsdk.builder.transaction_builder import DeployTransactionBuilder, CallTransactionBuilder
from iconsdk.libs.in_memory_zip import gen_deploy_data_content
from iconsdk.signed_transaction import SignedTransaction
from tbears.libs.icon_integrate_test import IconIntegrateTestBase, SCORE_INSTALL_ADDRESS
from iconservice.base.exception import IconScoreException


DIR_PATH = os.path.abspath(os.path.dirname(__file__))


class TestTest(IconIntegrateTestBase):
    TEST_HTTP_ENDPOINT_URI_V3 = "http://127.0.0.1:9000/api/v3"
    SCORE_PROJECT = os.path.abspath(os.path.join(DIR_PATH, ".."))

    def setUp(self):
        super().setUp()

        self.icon_service = None
        # if you want to send request to network, uncomment next line and set self.TEST_HTTP_ENDPOINT_URI_V3
        # self.icon_service = IconService(HTTPProvider(self.TEST_HTTP_ENDPOINT_URI_V3))

        # install SCORE
        self._std_reference_basic = self._deploy_score()["scoreAddress"]

    def _deploy_score(self, to: str = SCORE_INSTALL_ADDRESS) -> dict:
        # Generates an instance of transaction for deploying SCORE.
        transaction = (
            DeployTransactionBuilder()
            .from_(self._test1.get_address())
            .to(to)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .content_type("application/zip")
            .content(gen_deploy_data_content(self.SCORE_PROJECT))
            .build()
        )

        # Returns the signed transaction object having a signature
        signed_transaction = SignedTransaction(transaction, self._test1)

        # process the transaction in local
        tx_result = self.process_transaction(signed_transaction, self.icon_service)

        self.assertEqual(True, tx_result["status"])
        self.assertTrue("scoreAddress" in tx_result)

        return tx_result

    def test_relay_fail(self):
        cases = [
            {
                "_symbols": json.dumps(["BTC"]),
                "_rates": json.dumps([11, 12]),
                "_resolve_times": json.dumps([21]),
                "_request_ids": json.dumps([31]),
            },
            {
                "_symbols": json.dumps(["BTC"]),
                "_rates": json.dumps([11]),
                "_resolve_times": json.dumps([21, 22]),
                "_request_ids": json.dumps([31]),
            },
            {
                "_symbols": json.dumps(["BTC"]),
                "_rates": json.dumps([11]),
                "_resolve_times": json.dumps([21]),
                "_request_ids": json.dumps([31, 32]),
            },
        ]

        for case in cases:
            transaction = (
                CallTransactionBuilder()
                .from_(self._test1.get_address())
                .to(self._std_reference_basic)
                .step_limit(100_000_000_000)
                .nid(3)
                .nonce(100)
                .method("relay")
                .params(case)
                .build()
            )
            tx_result = self.process_transaction(
                SignedTransaction(transaction, self._test1), self.icon_service
            )
            self.assertEqual(False, tx_result["status"])

    def set_refs_by_using_relay(self, _symbols, _rates, _resolve_times, _request_ids):
        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._std_reference_basic)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay")
            .params(
                {
                    "_symbols": json.dumps(_symbols),
                    "_rates": json.dumps(_rates),
                    "_resolve_times": json.dumps(_resolve_times),
                    "_request_ids": json.dumps(_request_ids),
                }
            )
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

    def test_relay_and_get_reference_data(self):
        bases = ["BTC", "ETH", "BAND", "ICX", "SUSD"]
        quotes = ["ETH", "BAND", "ICX", "SUSD", "BTC"]

        relay_data = [
            ["BTC", "ETH", "BAND", "ICX", "SUSD"],
            [16000000000, 8000000000, 4000000000, 2000000000, 1000000000],
            [1600419716, 1600419717, 1600419718, 1600419719, 1600419720],
            [1, 2, 3, 4, 5],
        ]

        for (base, quote) in zip(bases, quotes):
            with self.assertRaises(IconScoreException) as e:
                call = (
                    CallBuilder()
                    .from_(self._test1.get_address())
                    .to(self._std_reference_basic)
                    .method("get_reference_data")
                    .params({"_base": base, "_quote": quote})
                    .build()
                )
                self.process_call(call, self.icon_service)

            self.assertEqual(e.exception.code, 32)
            self.assertEqual(e.exception.message, "REF_DATA_NOT_AVAILABLE")

        self.set_refs_by_using_relay(*relay_data)

        expected_rates = [
            (b * int(1e18) // q)
            for (b, q) in list(zip(relay_data[1], relay_data[1][1:] + relay_data[1][:1]))
        ]
        last_update_base_list = [b * int(1e6) for b in relay_data[2]]
        last_update_quote_list = [b * int(1e6) for b in relay_data[2][1:] + relay_data[2][:1]]
        for (base, quote, rate, last_update_base, last_update_quote) in zip(
            bases, quotes, expected_rates, last_update_base_list, last_update_quote_list
        ):
            call = (
                CallBuilder()
                .from_(self._test1.get_address())
                .to(self._std_reference_basic)
                .method("get_reference_data")
                .params({"_base": base, "_quote": quote})
                .build()
            )
            response = self.process_call(call, self.icon_service)
            self.assertEqual(rate, response["rate"])
            self.assertEqual(last_update_base, response["last_update_base"])
            self.assertEqual(last_update_quote, response["last_update_quote"])

    def test_relay_and_get_reference_data_bulk(self):
        bases = ["BTC", "ETH", "BAND", "ICX", "SUSD"]
        quotes = ["ETH", "BAND", "ICX", "SUSD", "BTC"]

        relay_data = [
            ["BTC", "ETH", "BAND", "ICX", "SUSD"],
            [16000000000, 8000000000, 4000000000, 2000000000, 1000000000],
            [1600419716, 1600419717, 1600419718, 1600419719, 1600419720],
            [1, 2, 3, 4, 5],
        ]

        with self.assertRaises(IconScoreException) as e:
            call = (
                CallBuilder()
                .from_(self._test1.get_address())
                .to(self._std_reference_basic)
                .method("get_reference_data_bulk")
                .params({"_bases": json.dumps(bases), "_quotes": json.dumps(quotes)})
                .build()
            )
            self.process_call(call, self.icon_service)

        self.assertEqual(e.exception.code, 32)
        self.assertEqual(e.exception.message, "REF_DATA_NOT_AVAILABLE")

        self.set_refs_by_using_relay(*relay_data)

        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._std_reference_basic)
            .method("get_reference_data_bulk")
            .params({"_bases": json.dumps(bases), "_quotes": json.dumps(quotes)})
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(
            [
                {
                    "rate": 2000000000000000000,
                    "last_update_base": 1600419716000000,
                    "last_update_quote": 1600419717000000,
                },
                {
                    "rate": 2000000000000000000,
                    "last_update_base": 1600419717000000,
                    "last_update_quote": 1600419718000000,
                },
                {
                    "rate": 2000000000000000000,
                    "last_update_base": 1600419718000000,
                    "last_update_quote": 1600419719000000,
                },
                {
                    "rate": 2000000000000000000,
                    "last_update_base": 1600419719000000,
                    "last_update_quote": 1600419720000000,
                },
                {
                    "rate": 62500000000000000,
                    "last_update_base": 1600419720000000,
                    "last_update_quote": 1600419716000000,
                },
            ],
            response,
        )

