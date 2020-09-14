import os

from iconsdk.builder.call_builder import CallBuilder
from iconsdk.builder.transaction_builder import DeployTransactionBuilder, CallTransactionBuilder
from iconsdk.libs.in_memory_zip import gen_deploy_data_content
from iconsdk.signed_transaction import SignedTransaction
from tbears.libs.icon_integrate_test import IconIntegrateTestBase, SCORE_INSTALL_ADDRESS
from ..pyobi import *

DIR_PATH = os.path.abspath(os.path.dirname(__file__))


class TestTest(IconIntegrateTestBase):
    TEST_HTTP_ENDPOINT_URI_V3 = "http://127.0.0.1:9000/api/v3"
    AGGREGATOR_PROJECT = os.path.abspath(os.path.join(DIR_PATH, ".."))
    BRIDGE_PROJECT = os.path.abspath(os.path.join(DIR_PATH, "../../bridge"))
    EXPECTED_PRICES = {
        "USD": 1000000000,
        "RENBTC": 10260210000000,
        "WBTC": 10327140000000,
        "DIA": 2244000000,
        "BTM": 73443500,
        "IOTX": 11043595,
        "FET": 85140000,
        "JST": 42351025,
        "MCO": 4190500000,
        "KMD": 621285000,
        "BTS": 27380329,
        "QKC": 30134499,
        "YAMV2": 12047944700,
        "XZC": 14707500000,
        "UOS": 5230000000,
        "AKRO": 98833000,
        "HNT": 24615310,
        "HOT": 654846000,
        "KAI": 11536359,
        "OGN": 141293435,
        "WRX": 183095500,
        "KDA": 305222999,
        "ORN": 488868500,
        "FOR": 1251150000,
        "AST": 110447000,
        "STORJ": 295666500,
        "BTC": 10308200000000,
        "ETH": 364790000000,
        "USDT": 1000200000,
        "XRP": 241500000,
        "LINK": 12307000000,
        "DOT": 4419650000,
        "BCH": 224290000000,
        "LTC": 48474500000,
        "ADA": 95000000,
        "BSV": 165440000000,
        "CRO": 154100000,
        "BNB": 23627600000,
        "EOS": 2795000000,
        "XTZ": 2521700000,
        "TRX": 32700000,
        "XLM": 81300000,
        "ATOM": 5265800000,
        "XMR": 83807600000,
        "OKB": 3220050000,
        "USDC": 1000000000,
        "NEO": 19075000000,
        "XEM": 130300000,
        "LEO": 1148000000,
        "HT": 4733000000,
        "VET": 12900000,
    }

    def setUp(self):
        super().setUp()

        self.icon_service = None
        # if you want to send request to network, uncomment next line and set self.TEST_HTTP_ENDPOINT_URI_V3
        # self.icon_service = IconService(HTTPProvider(self.TEST_HTTP_ENDPOINT_URI_V3))

        # install SCORE
        self._bridge_address = self._deploy_bridge()["scoreAddress"]
        self._aggregator = self._deploy_aggregator(self._bridge_address)["scoreAddress"]

    def _deploy_bridge(self, to: str = SCORE_INSTALL_ADDRESS) -> dict:
        params = {}
        params[
            "validators_bytes"
        ] = "0000000700000040c528e467394d305f7420206669bf096096dc5ba1149ec8082d87674d58336cd4c03443d06d014dbd9acbf86311f2946020bfd982429196b2d43b3afc7a97a27f00000000000f42410000004004379b986845684d64d07afb724253499e683df3e1149176086893757d69e12e224fc93ee1432e3b041e2c228030f6e87125726b5776a509721a033fa025230400000000000f4241000000405b0846aeeeeee76b12ea40692205333c4ee8638ef4e534b83e522b8e555a95cb1dabfa6672a1d6dbb46f5885a0b04039dc8ce03b91594ce9d8ebefc8352a956d00000000000f424100000040d50a83a038973637a94cf9d0dc49f4bbb15f283cd765dfb4a91d665920c3d67f566a1e79437b93a4a567efd8a60d6e08406ac27bcf257ca3d9ee1fa8e702f14a00000000000f424100000040022e3ddcfadb8036725c14bc85fe421694b64bf72fad3d4303782d6e11b13e2158d0491518dc149d1562fbeb5d0bea60da2299dd808ee9dcf3faf077a3b1e49700000000000f4241000000400408e635dd9e2fb17bdb98210e9063e27ceae95998fcde39cd4621cf9a66cddaab84a9737f526a2611a72d3cc59cf1483125bb038cb1225b8a979ce02f811ba500000000000f424100000040a1a5e3b177b2e766a01559fe495767f859ec1caf068f88b9b3bc9e51a0f2838d432911f4b123cd1f371b011fa05fbea18770d8289ff24cbe121e6a050ad7e2ae00000000000f4241"

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

    def _deploy_aggregator(self, bridge_address: str, to: str = SCORE_INSTALL_ADDRESS) -> dict:
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
            .content(gen_deploy_data_content(self.AGGREGATOR_PROJECT))
            .build()
        )

        # Returns the signed transaction object having a signature
        signed_transaction = SignedTransaction(transaction, self._test1)

        # process the transaction in local
        tx_result = self.process_transaction(signed_transaction, self.icon_service)

        self.assertEqual(True, tx_result["status"])
        self.assertTrue("scoreAddress" in tx_result)

        return tx_result

    def test_aggregator_get_bridge_address(self):
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._aggregator)
            .method("get_bridge_address")
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(self._bridge_address, response)

    def test_aggregator_set_bridge(self):
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._aggregator)
            .method("get_bridge_address")
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(self._bridge_address, response)

        transaction = (
            CallTransactionBuilder()
            .from_(self._test1.get_address())
            .to(self._aggregator)
            .step_limit(100_000_000_000)
            .nid(3)
            .nonce(100)
            .method("set_bridge")
            .params({"bridge_address": self._test1.get_address()})
            .build()
        )
        signed_transaction = SignedTransaction(transaction, self._test1)
        tx_result = self.process_transaction(signed_transaction, self.icon_service)
        self.assertEqual(True, tx_result["status"])

        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._aggregator)
            .method("get_bridge_address")
            .build()
        )
        response = self.process_call(call, self.icon_service)
        self.assertEqual(self._test1.get_address(), response)

    def test_relay_1(self):
        request_key = bytes.fromhex(
            "0000000862616e647465616d0000000000000008000000c5000000190000000652454e4254430000000457425443000000034449410000000342544d00000004494f545800000003464554000000034a5354000000034d434f000000034b4d440000000342545300000003514b430000000559414d563200000003585a4300000003554f5300000004414b524f00000003484e5400000003484f54000000034b4149000000034f474e00000003575258000000034b4441000000034f524e00000003464f52000000034153540000000553544f524a000000003b9aca0000000000000000040000000000000003"
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
        self.assertEqual(None, self.process_call(call, self.icon_service))

        params = {}
        params["proof"] = bytes.fromhex(
            "00000000001010fb000000a0057295b19ea9b639ec51f078949305f1867fa42f2021e18130afd3153f183757f9ed4517d45aea1319218fecb6e4453b57da44748be82f48467d94a7c0eecb9383cb0f5ab95240b32bc622097a6c5666a338f1f1dcc1f495b0fa18d250d02e242b6a7e0f44ed9c179a47a40f93d5824189a5426d6c3f77692de28e50e20a33dd02f0ca5e4bae557663910f3f4a2d9ffb1c3f6314a5eaa3a3fb696c7155a929ba000000c03561783e9c3bdf932a16580fc0c9ceffec4c509073fff65a42871bfab64408ae0fe00a766f2eb8a5a3e1055c4e8971a64c7e5af7b479f9ed19701f5079062bf2a54ee3f75ccd1ad2b277eb8f63c87eb3a5d94f8b4f0f7069eec2e7f617da3066ea01cd62e714b603323a21a4a7382f8ab04788c867a0c99ade687d00e7d5fe626e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d992b65122f7d4795f8e84204b71a395050f83c11603cd69ec25dc41f1818c7230000034200000005000000209f36d6e91fb9de9424da2d7285e4cb47a599a5445d8c74a537f082d0990b24ca0000002007c4e71363e1f7fed66057e1db8cba51ffd2e3924c36b8a37ef3c6f0495419011c0000001074080211fb1010000000000022480a200000004512240a2072b0f5d879f0a77a94ae215564661dc37a93d59032ed34d8cf707bc15d1141a610012a0c08fefae7fa0510a1ccf1c101320f62616e642d6775616e79752d706f6100000020abdcf08a3b96e3bb3923a139d850e60ecd8fcc43fe8291f6ddf6666128b1c4fa000000201624404b9a6001a8c7a8e53d43d2cb11a9d76a730332f193951189e14b738f341b0000001074080211fb1010000000000022480a200000004512240a2072b0f5d879f0a77a94ae215564661dc37a93d59032ed34d8cf707bc15d1141a610012a0c08fefae7fa0510ddd1a8ae01320f62616e642d6775616e79752d706f61000000200d4b48e064dbb553de7ee045b26069d018d8333fa7a22f34df4c28fa88967100000000207261ea03388705cc9e0e91c9acb1ad258b29be2facf8e405a23aa1e6527c55721b0000001074080211fb1010000000000022480a200000004512240a2072b0f5d879f0a77a94ae215564661dc37a93d59032ed34d8cf707bc15d1141a610012a0c08fefae7fa051098a6b5c201320f62616e642d6775616e79752d706f6100000020c000fa880e0c16ee7a586f3953f330079da3c290b604eca72a68a650acfd9c6d0000002029bd362027bd238f365d47beba26e4910ea05c46c2ac9e23fb02c2a2841b86671c0000001074080211fb1010000000000022480a200000004512240a2072b0f5d879f0a77a94ae215564661dc37a93d59032ed34d8cf707bc15d1141a610012a0c08fefae7fa0510f2d4c39f01320f62616e642d6775616e79752d706f6100000020f9824317ee9c376102e9e8d65ccade531eb25f9318bee522ce60fe729dff740c000000205e73e883509cdc9c5f65b1b93b72601ff535baeae2627bddfe106f1866464ce31b0000001074080211fb1010000000000022480a200000004512240a2072b0f5d879f0a77a94ae215564661dc37a93d59032ed34d8cf707bc15d1141a610012a0c08fefae7fa0510a5ca93a301320f62616e642d6775616e79752d706f61000001ed0000000862616e647465616d0000000000000008000000c5000000190000000652454e4254430000000457425443000000034449410000000342544d00000004494f545800000003464554000000034a5354000000034d434f000000034b4d440000000342545300000003514b430000000559414d563200000003585a4300000003554f5300000004414b524f00000003484e5400000003484f54000000034b4149000000034f474e00000003575258000000034b4441000000034f524e00000003464f52000000034153540000000553544f524a000000003b9aca00000000000000000400000000000000030000000862616e647465616d00000000000249270000000000000004000000005f59fd1d000000005f59fd2000000001000000cc0000001900000954e42c2080000009647982e1000000000085c0b900000000000460a8ac0000000000a8830b000000000513222000000000028639b100000000f9c5f4a00000000025080e880000000001a1ca690000000001cbd0e300000002ce1d0bfc000000036ca2a3e00000000137bb77800000000005e41268000000000177998e00000000270828300000000000b007e700000000086bf77b000000000ae9d0cc0000000012315557000000001d238a94000000004a9308b0000000000695499800000000119f834400000000001010dd00000472000000150101000000000000000200000000001010e50000002027191320f20d7d1d49298adf147c81e93d6f75368d47e76f06dd6cd159cfdfd50002000000000000000400000000001010e500000020e61b88d094cbee0b6072b68be24b2cdf12179f713c02426d9a8bed736a9f28840103000000000000000600000000001010e5000000208e419f70b2739881779b9724c8cc555aa51f7e47cc2851cad3491ef5c3f66fa70104000000000000000a00000000001010e500000020639e1d346c0e032821e6ebf564ded5b4e47d81506080cfa10fcd14b540f5180e0105000000000000001a00000000001010e5000000204da2e82153c12c8416f1e0da077b3d0117a5aedb6726ed9e19672d16e3da981f0106000000000000003a00000000001010e500000020758de4fcfa8a519f3572bd0abb68ba905a96c81f0618be22a16aca1f3949c81e0107000000000000007a00000000001010e500000020f51c43e0daee501d7574f3fd50ff6a52a60eb4684cfed386d7318ac2f5d48977010800000000000000ba00000000001010e50000002052f853fd38afc7080b3cb9f5766e561334c063b747f443965e0713014157d4e5010900000000000001ba00000000001010e5000000208a547560966f7b998eaa36840df760701d41e8379e63e369c1c96f48ab18165d010a00000000000003ba00000000001010e500000020431fd93b5c3f6f08e9d535930702c1bf4b3da9b272b80f52bbc31af1aed1c7fe010b00000000000005ba00000000001010e500000020b58d803307e5928ca764feac86123a14375d1192996bab45dc2e392e4f06cada010c00000000000009ba00000000001010e500000020bb8019bf65ea8a24a033270728fc5e7669928e9f67253b69bb5429abb10430c9010d00000000000019b900000000001010e50000002062d6b7db3b4eb0378b46404cae011d0ae5cd798f6b0697a3ec75310b5b42a63b010e00000000000039b900000000001010e50000002007e8397f6ac70ea8dfc18db69ea006f26493ebf186ac6aa46570d24d81f1fde6010f00000000000059b900000000001010e500000020d612111a304e5007d7e9d9af21034ef494c1f0d1203848098b2774a84adce9f40110000000000000d9b600000000001010e50000002036be7502c0b68cfec2860c6507354a089a3e58ee36552052337e83fdd56ab5eb011100000000000159b300000000001010e50000002003016bbbcb9037402b22fe5f61a486208295501a8919adfb74694b2dc332ee120112000000000002497700000000001010e5000000204f2abf0d75f8f681c53ee01590b5139bcc22167dde76e7468cd5f265bd7333f701140000000000061f7a00000000001010e50000002039e7a064c1897a83b4b02fc1f52e5fbaa74fe75ebd614040030f6edb7ee4e429011500000000000b5d5400000000001010e5000000201963acbd3922fb0f34bb83dba2c961742f0e1b33945e0ad9da6e5eb98741112f011600000000000e39f600000000001010fa00000020d6cef64943fe0ee0afd1d30d4138ee45dbd3986faf7aff945e3166ebf0944b80"
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
        self.assertEqual(
            {
                "client_id": "bandteam",
                "request_id": 149799,
                "ans_count": 4,
                "request_time": 1599733021,
                "resolve_time": 1599733024,
                "resolve_status": 1,
                "result": b"\x00\x00\x00\x19\x00\x00\tT\xe4, \x80\x00\x00\tdy\x82\xe1\x00\x00\x00\x00\x00\x85\xc0\xb9\x00\x00\x00\x00\x00\x04`\xa8\xac\x00\x00\x00\x00\x00\xa8\x83\x0b\x00\x00\x00\x00\x05\x13\" \x00\x00\x00\x00\x02\x869\xb1\x00\x00\x00\x00\xf9\xc5\xf4\xa0\x00\x00\x00\x00%\x08\x0e\x88\x00\x00\x00\x00\x01\xa1\xcai\x00\x00\x00\x00\x01\xcb\xd0\xe3\x00\x00\x00\x02\xce\x1d\x0b\xfc\x00\x00\x00\x03l\xa2\xa3\xe0\x00\x00\x00\x017\xbbw\x80\x00\x00\x00\x00\x05\xe4\x12h\x00\x00\x00\x00\x01w\x99\x8e\x00\x00\x00\x00'\x08(0\x00\x00\x00\x00\x00\xb0\x07\xe7\x00\x00\x00\x00\x08k\xf7{\x00\x00\x00\x00\n\xe9\xd0\xcc\x00\x00\x00\x00\x121UW\x00\x00\x00\x00\x1d#\x8a\x94\x00\x00\x00\x00J\x93\x08\xb0\x00\x00\x00\x00\x06\x95I\x98\x00\x00\x00\x00\x11\x9f\x83D",
            },
            self.process_call(call, self.icon_service),
        )

    def test_relay_2(self):
        request_key = bytes.fromhex(
            "0000000862616e647465616d0000000000000008000000be000000190000000342544300000003455448000000045553445400000003585250000000044c494e4b00000003444f5400000003424348000000034c544300000003414441000000034253560000000343524f00000003424e4200000003454f530000000358545a0000000354525800000003584c4d0000000441544f4d00000003584d52000000034f4b420000000455534443000000034e454f0000000358454d000000034c454f00000002485400000003564554000000003b9aca0000000000000000040000000000000003"
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
        self.assertEqual(None, self.process_call(call, self.icon_service))

        params = {}
        params["proof"] = bytes.fromhex(
            "0000000000108029000000a0d0bfdf1590ebd8039b07664d5ac2365fb3e072a6b1262d1469a8f259fbd883cbfecfa66f4d6f2e25bc58a93ff47f3a04ddb6bd78c3224f2046d2237a4cddf0d4740addcedaeee44fa921ef488bda8476e02fa8d09598b4574d6e2dc36d84d1e82b6a7e0f44ed9c179a47a40f93d5824189a5426d6c3f77692de28e50e20a33dd9944f27f0adccdf18da69a4459699ec807ace5a46f34cba1dff347a70a71ff0c000000c03561783e9c3bdf932a16580fc0c9ceffec4c509073fff65a42871bfab64408ae51fa188390b5725a3f086ef5e14143b975811b207a0028d46a96388ade69f6c0fbdda7f560010b35c90b6d9a601cb1a8fdf5baf5a1c6ca54d260ab857cea73a3ea01cd62e714b603323a21a4a7382f8ab04788c867a0c99ade687d00e7d5fe626e340b9cffb37a989ca544e6bb780a2c78901d3fb33738768511a30617afa01d992b65122f7d4795f8e84204b71a395050f83c11603cd69ec25dc41f1818c723000003420000000500000020a68c145b6c458d8f7c47508fff2ec08d1c7e2922aa7490985e4919e860b265f7000000204b01cb266870df0b6f3d01ecdf1007e97a329e1364413871c83e00ca5253b7171b0000001074080211298010000000000022480a200000004512240a20a90ab8858148df88fcee18ff6f17340f095a75004424aed9dd888a472b0a349d10012a0c08da9eedfa0510f9cacab901320f62616e642d6775616e79752d706f6100000020f6eb8298880142c59b6cff4eb4fe600c35af602b4842ee6a27e5d0e52cfb2c77000000200ceea4eb0a9b432b9fb14fe6aadd102449df2b3a5017e50927507dedb70247151c0000001074080211298010000000000022480a200000004512240a20a90ab8858148df88fcee18ff6f17340f095a75004424aed9dd888a472b0a349d10012a0c08da9eedfa051085f7b69d01320f62616e642d6775616e79752d706f6100000020ec0f0f8a9a8c5d70a0a809df2abab106d40c973959e2394085548d3728a8651a000000205e068bacba6ce6c8d9695739a9cedea180ec3741c59476b7daf78958a8eb190e1b0000001074080211298010000000000022480a200000004512240a20a90ab8858148df88fcee18ff6f17340f095a75004424aed9dd888a472b0a349d10012a0c08da9eedfa0510e591d89901320f62616e642d6775616e79752d706f6100000020d8f409d8ea3e0ebbb56d9d8ea64e7aed57d902c92585d7512fc734c5660c92ea000000204eabdf25e64c69ebe9ecf406a7a8c74e217b873f7a543c98ebd17510f3900b341b0000001074080211298010000000000022480a200000004512240a20a90ab8858148df88fcee18ff6f17340f095a75004424aed9dd888a472b0a349d10012a0c08da9eedfa051091f4ccb601320f62616e642d6775616e79752d706f61000000207edc09c892d2297b2476514c3b422b3e70c09752e309800097017edcec0b0375000000204bc9fde7eefa751e978cb95aa85710a800d41b43729da5d361c30f097b0dd59f1c0000001074080211298010000000000022480a200000004512240a20a90ab8858148df88fcee18ff6f17340f095a75004424aed9dd888a472b0a349d10012a0c08da9eedfa0510c5a1d4b701320f62616e642d6775616e79752d706f61000001e60000000862616e647465616d0000000000000008000000be000000190000000342544300000003455448000000045553445400000003585250000000044c494e4b00000003444f5400000003424348000000034c544300000003414441000000034253560000000343524f00000003424e4200000003454f530000000358545a0000000354525800000003584c4d0000000441544f4d00000003584d52000000034f4b420000000455534443000000034e454f0000000358454d000000034c454f00000002485400000003564554000000003b9aca00000000000000000400000000000000030000000862616e647465616d00000000000261e70000000000000003000000005f5b4e73000000005f5b4e7900000001000000cc000000190000096010996a0000000054ef2da980000000003b9dd740000000000e64ff6000000002dd8deac000000001076e81d00000003438b9cc800000000b494e2ba00000000005a995c00000002684fe100000000000092f6120000000058050908000000000a69850c000000000964e16a00000000001f2f6600000000004d88a200000000139ddbb400000001383527d8000000000bfee1050000000003b9aca000000000470f566c00000000007c4386000000000446d1700000000011a1bd9400000000000c4d6a00000000000107fe00000047200000015010100000000000000020000000000107fe700000020cc89e51472a95233993a50d3ed72522fbc64697dbfc263c7280626a43e2bbfbd010200000000000000040000000000107fea00000020dbe1cd2f8eace946ee2d3b5cfe5b6640cbb84d97931b7a64c1c332a7b828a2d1010300000000000000080000000000107ff1000000206cf577315faa740914d99a559d285c5e3a6b4cd315b4ece543e572a158ddfbe00004000000000000000e0000000000107ff300000020b512faed452949ed30044b98062c217d859b9eb7874f43a0a562f90b139545dc0105000000000000001e0000000000107ff300000020ee741d94fc9d31ab9711d78310dfb9530efaba9b6c40aa25cf4e30cd4aeae0720106000000000000003e0000000000107ff300000020fe65397542d1056533427d7b78d0f41e331ec9e34c3479fb6770eab8071ec81f0107000000000000007e0000000000107ff300000020ac934cd2882ec63c11ecce393cb3f03bd8c472094476a1f40610c5c62854c99a010800000000000000fe0000000000107ff300000020a7f213cb1f41f86430af1ad66a545e2f27e7601918010e1f4fb55a9b826287880109000000000000017e0000000000107ff300000020af9cb5f494b6b110660539ce6274e427752b9814a079a0cda7cb8791f12f29a0010a000000000000027e0000000000107ff30000002019b6b94acf33480d8cc0d634bc12b816f38e7d772a0ab14cff5a424311500c43010b000000000000067e0000000000107ff300000020fe5626ef4032dc1d187d38e73ea9ed5498a0d75eb7a31265e3f93ce99937c542010c0000000000000a7e0000000000107ff3000000205197af0b3dcc7a1cde645f6815234c6c36956729a3997f4fd3e82b3970aef6de010d000000000000127e0000000000107ff3000000200b6a198005e415ba12737c47e77df605179632dad5c482e08a184a73ceb8439d010e000000000000327d0000000000107ff30000002018e0dd7765d604ab59f865f0e4d23c93216ea53f121172776a1b0f2516ab6d6b010f000000000000727d0000000000107ff300000020e826f091e1fedf776cc2fe12d1b7fea7256839cc097d0c1d4772aa5933729b7f0110000000000000f27a0000000000107ff30000002036be7502c0b68cfec2860c6507354a089a3e58ee36552052337e83fdd56ab5eb011100000000000172770000000000107ff30000002003016bbbcb9037402b22fe5f61a486208295501a8919adfb74694b2dc332ee120112000000000002623c0000000000107ff3000000207469c0028d137e248b3190c1f17b46022317a975f5a277be188cf19d8a79bf1201140000000000069c680000000000107ff40000002019d02d9a701ef232a5b60a9a71650ec4f3f32cd40a0a2a58eee7630330479f45011500000000000bda420000000000107ff4000000201963acbd3922fb0f34bb83dba2c961742f0e1b33945e0ad9da6e5eb98741112f011600000000000ecfa9000000000010802800000020f3d76f18d212b9470d10862802ce13aab8b78ddd28f00d1d82a7556f90ebeb60"
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
        self.assertEqual(
            {
                "client_id": "bandteam",
                "request_id": 156135,
                "ans_count": 3,
                "request_time": 1599819379,
                "resolve_time": 1599819385,
                "resolve_status": 1,
                "result": b"\x00\x00\x00\x19\x00\x00\t`\x10\x99j\x00\x00\x00\x00T\xef-\xa9\x80\x00\x00\x00\x00;\x9d\xd7@\x00\x00\x00\x00\x0ed\xff`\x00\x00\x00\x02\xdd\x8d\xea\xc0\x00\x00\x00\x01\x07n\x81\xd0\x00\x00\x0048\xb9\xcc\x80\x00\x00\x00\x0bIN+\xa0\x00\x00\x00\x00\x05\xa9\x95\xc0\x00\x00\x00&\x84\xfe\x10\x00\x00\x00\x00\x00\t/a \x00\x00\x00\x05\x80P\x90\x80\x00\x00\x00\x00\xa6\x98P\xc0\x00\x00\x00\x00\x96N\x16\xa0\x00\x00\x00\x00\x01\xf2\xf6`\x00\x00\x00\x00\x04\xd8\x8a \x00\x00\x00\x019\xdd\xbb@\x00\x00\x00\x13\x83R}\x80\x00\x00\x00\x00\xbf\xee\x10P\x00\x00\x00\x00;\x9a\xca\x00\x00\x00\x00\x04p\xf5f\xc0\x00\x00\x00\x00\x07\xc48`\x00\x00\x00\x00Dm\x17\x00\x00\x00\x00\x01\x1a\x1b\xd9@\x00\x00\x00\x00\x00\xc4\xd6\xa0",
            },
            self.process_call(call, self.icon_service),
        )

    def test_get_rate(self):
        self.test_relay_1()
        self.test_relay_2()

        for symbol, price_in_usd in self.EXPECTED_PRICES.items():
            call = (
                CallBuilder()
                .from_(self._test1.get_address())
                .to(self._aggregator)
                .method("get_rate")
                .params({"symbol": symbol})
                .build()
            )
            self.assertEqual(hex(price_in_usd), self.process_call(call, self.icon_service))

    def test_get_reference_data(self):
        self.test_relay_1()
        self.test_relay_2()

        test_pairs = ["BTC/LTC", "FET/VET", "YAMV2/BNB", "DOT/USDC"]
        call = (
            CallBuilder()
            .from_(self._test1.get_address())
            .to(self._aggregator)
            .method("get_reference_data")
            .params({"encoded_pairs": PyObi("[string]").encode(test_pairs)})
            .build()
        )

        self.assertEqual(
            [
                (self.EXPECTED_PRICES["BTC"] * 1000000000000000000) // self.EXPECTED_PRICES["LTC"],
                (self.EXPECTED_PRICES["FET"] * 1000000000000000000) // self.EXPECTED_PRICES["VET"],
                (self.EXPECTED_PRICES["YAMV2"] * 1000000000000000000)
                // self.EXPECTED_PRICES["BNB"],
                (self.EXPECTED_PRICES["DOT"] * 1000000000000000000) // self.EXPECTED_PRICES["USDC"],
            ],
            self.process_call(call, self.icon_service),
        )

