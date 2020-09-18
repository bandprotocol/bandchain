import os

from iconsdk.builder.call_builder import CallBuilder
from iconsdk.builder.transaction_builder import DeployTransactionBuilder, CallTransactionBuilder
from iconsdk.libs.in_memory_zip import gen_deploy_data_content
from iconsdk.signed_transaction import SignedTransaction
from tbears.libs.icon_integrate_test import IconIntegrateTestBase, SCORE_INSTALL_ADDRESS
from iconservice.base.exception import IconScoreException
from ..pyobi import *

DIR_PATH = os.path.abspath(os.path.dirname(__file__))


class TestTest(IconIntegrateTestBase):
    TEST_HTTP_ENDPOINT_URI_V3 = "http://127.0.0.1:9000/api/v3"
    STD_REFERENCE_PROXY_PROJECT = os.path.abspath(os.path.join(DIR_PATH, ".."))
    STD_BASIC_PROJECT = os.path.abspath(os.path.join(DIR_PATH, "../../std_reference_basic"))

    def setUp(self):
        super().setUp()

        self.icon_service = None
        # if you want to send request to network, uncomment next line and set self.TEST_HTTP_ENDPOINT_URI_V3
        # self.icon_service = IconService(HTTPProvider(self.TEST_HTTP_ENDPOINT_URI_V3))

        # install SCORE
        self._std_basic = self._deploy_std_basic()["scoreAddress"]
        self._std_reference_proxy = self._deploy_std_reference_proxy(self._std_basic)[
            "scoreAddress"
        ]

    def _deploy_std_basic(self, to: str = SCORE_INSTALL_ADDRESS) -> dict:
        # Generates an instance of transaction for deploying SCORE.
        transaction = (
            DeployTransactionBuilder()
            .from_(self._test1.get_address())
            .to(to)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .content_type("application/zip")
            .content(gen_deploy_data_content(self.STD_BASIC_PROJECT))
            .build()
        )

        # Returns the signed transaction object having a signature
        signed_transaction = SignedTransaction(transaction, self._test1)

        # process the transaction in local
        tx_result = self.process_transaction(signed_transaction, self.icon_service)

        self.assertEqual(True, tx_result["status"])
        self.assertTrue("scoreAddress" in tx_result)

        return tx_result

    def _deploy_std_reference_proxy(self, _ref: str, to: str = SCORE_INSTALL_ADDRESS) -> dict:
        # Generates an instance of transaction for deploying SCORE.
        transaction = (
            DeployTransactionBuilder()
            .from_(self._test1.get_address())
            .to(to)
            .params({"_ref": _ref})
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .content_type("application/zip")
            .content(gen_deploy_data_content(self.STD_REFERENCE_PROXY_PROJECT))
            .build()
        )

        # Returns the signed transaction object having a signature
        signed_transaction = SignedTransaction(transaction, self._test1)

        # process the transaction in local
        tx_result = self.process_transaction(signed_transaction, self.icon_service)

        self.assertEqual(True, tx_result["status"])
        self.assertTrue("scoreAddress" in tx_result)

        return tx_result

    def test_std_reference_proxy_get_ref_address(self):
        call = (
            CallBuilder()
            .from_(self._std_basic)
            .to(self._std_reference_proxy)
            .method("get_ref")
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(self._std_basic, response)

    def test_std_reference_proxy_set_std_basic(self):
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._std_reference_proxy)
            .method("get_ref")
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(self._std_basic, response)

        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._std_reference_proxy)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("set_ref")
            .params({"_ref": self._test1.get_address()})
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._std_reference_proxy)
            .method("get_ref")
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(self._test1.get_address(), response)

    def set_initial_refs(self, symbols, rates, resolve_times):
        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._std_basic)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("relay")
            .params(
                {
                    "_encoded_data_list": PyObi(
                        "[{symbol:string,rate:u64,resolve_time:u64}]"
                    ).encode(
                        [
                            {"symbol": s, "rate": r, "resolve_time": rt}
                            for (s, r, rt) in list(zip(symbols, rates, resolve_times))
                        ]
                    )
                }
            )
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

    def test_get_reference_data(self):
        symbols = ["BTC", "ETH", "BAND", "ICX", "SUSD"]
        rates = [int(16e9), int(8e9), int(4e9), int(2e9), int(1e9)]
        resolve_times = [1600419716, 1600419717, 1600419718, 1600419719, 1600419720]
        for symbol in symbols:
            with self.assertRaises(IconScoreException) as e:
                call = (
                    CallBuilder()
                    .from_(self._test1.get_address())
                    .to(self._std_reference_proxy)
                    .method("get_reference_data")
                    .params({"_base": symbol, "_quote": "USD"})
                    .build()
                )
                self.process_call(call, self.icon_service)

            self.assertEqual(e.exception.code, 32)
            self.assertEqual(e.exception.message, "REF_DATA_NOT_AVAILABLE")

        self.set_initial_refs(symbols, rates, resolve_times)

        pairs = ["BTC/ETH", "ETH/BAND", "BAND/ICX", "ICX/SUSD", "SUSD/BTC"]
        expected_rates = [
            (b * int(1e18) // q) for (b, q) in list(zip(rates, rates[1:] + rates[:1]))
        ]
        last_update_base_list = [b * int(1e6) for b in resolve_times]
        last_update_quote_list = [b * int(1e6) for b in resolve_times[1:] + resolve_times[:1]]
        for (pair, rate, last_update_base, last_update_quote) in list(
            zip(pairs, expected_rates, last_update_base_list, last_update_quote_list)
        ):
            [base, quote] = pair.split("/")
            call = (
                CallBuilder()
                .from_(self._test1.get_address())
                .to(self._std_reference_proxy)
                .method("get_reference_data")
                .params({"_base": base, "_quote": quote})
                .build()
            )
            response = self.process_call(call, self.icon_service)
            self.assertEqual(rate, response["rate"])
            self.assertEqual(last_update_base, response["last_update_base"])
            self.assertEqual(last_update_quote, response["last_update_quote"])

    def test_relay_and_get_reference_data_bulk(self):
        symbols = ["BTC", "ETH", "BAND", "ICX", "SUSD"]
        rates = [int(16e9), int(8e9), int(4e9), int(2e9), int(1e9)]
        resolve_times = [1600419716, 1600419717, 1600419718, 1600419719, 1600419720]

        pairs = ["BTC/ETH", "ETH/BAND", "BAND/ICX", "ICX/SUSD", "SUSD/BTC"]

        with self.assertRaises(IconScoreException) as e:
            call = (
                CallBuilder()
                .from_(self._test1.get_address())
                .to(self._std_reference_proxy)
                .method("get_reference_data_bulk")
                .params(
                    {
                        "_encoded_pairs": PyObi("[{base:string,quote:string}]").encode(
                            [
                                {"base": b, "quote": q}
                                for (b, q) in list(zip(symbols, symbols[1:] + symbols[:1]))
                            ]
                        )
                    }
                )
                .build()
            )
            self.process_call(call, self.icon_service)

        self.assertEqual(e.exception.code, 32)
        self.assertEqual(e.exception.message, "REF_DATA_NOT_AVAILABLE")

        self.set_initial_refs(symbols, rates, resolve_times)

        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._std_reference_proxy)
            .method("get_reference_data_bulk")
            .params(
                {
                    "_encoded_pairs": PyObi("[{base:string,quote:string}]").encode(
                        [
                            {"base": b, "quote": q}
                            for (b, q) in list(zip(symbols, symbols[1:] + symbols[:1]))
                        ]
                    )
                }
            )
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

