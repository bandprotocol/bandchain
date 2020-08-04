from ..bridge import BRIDGE
from tbears.libs.scoretest.score_test_case import ScoreTestCase


class TestBRIDGE(ScoreTestCase):

    def setUp(self):
        super().setUp()
        self.score = self.get_score_instance(BRIDGE, self.test_account1)

    # =-=-=-=-=-=-=-=-=-=-=-=-= \Utils =-=-=-=-=-=-=-=-=-=-=-=-=

    def test_merkle_leaf_hash(self):
        self.assertEqual(
            self.score.merkle_leaf_hash(
                bytes.fromhex('08d1082cc8d85a0833da8815ff1574675c415760e0aff7fb4e32de6de27faf86')).hex(),
            "35b401b2a74452d2252df60574e0a6c029885965ae48f006ebddc18e53427e26"
        )

    def test_merkle_inner_hash(self):
        self.assertEqual(
            self.score.merkle_inner_hash(
                bytes.fromhex(
                    'aa58d2b0be1dbadf21619645c9245062a283c19a6701bfe07d07cc2441687c85'),
                bytes.fromhex(
                    'ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb')
            ).hex(),
            "da32ca60aaddf58c20ff66ba1bb73ef126fea41822fedbfc10fbbedc8587ffe2"
        )

        self.assertEqual(
            self.score.merkle_inner_hash(
                bytes.fromhex(
                    '08d1082cc8d85a0833da8815ff1574675c415760e0aff7fb4e32de6de27faf86'),
                bytes.fromhex(
                    '789411d15a12768a9c3eb99d3453d6ebb4481c2a03ab59cc262a97e25757afe6')
            ).hex(),
            "ca48b611419f0848bf0fce9750caca6fd4fb8ef96ba8d7d3ccd4f05bf2af1661"
        )

    def test_encode_varint_unsigned(self):
        self.assertEqual(self.score.encode_varint_unsigned(116).hex(), "74")
        self.assertEqual(
            self.score.encode_varint_unsigned(14947).hex(), "e374")
        self.assertEqual(self.score.encode_varint_unsigned(
            244939043).hex(), "a3f2e574")

    def test_encode_varint_signed(self):
        self.assertEqual(self.score.encode_varint_signed(58).hex(), "74")
        self.assertEqual(
            self.score.encode_varint_signed(7473).hex(), "e274")
        self.assertEqual(self.score.encode_varint_signed(
            122469521).hex(), "a2f2e574")

    # =-=-=-=-=-=-=-=-=-=-=-=-= /Utils =-=-=-=-=-=-=-=-=-=-=-=-=

    # =-=-=-=-=-=-=-=-=-=-=-=-= \BlockHeaderMerkleParts =-=-=-=-=-=-=-=-=-=-=-=-=

    def test_get_block_header(self):
        data = b''
        data += bytes.fromhex(
            "32fa694879095840619f5e49380612bd296ff7e950eafb66ff654d99ca70869e")
        data += bytes.fromhex(
            "4BAEF831B309C193CC94DCF519657D832563B099A6F62C6FA8B7A043BA4F3B3B")
        data += bytes.fromhex(
            "5E1A8142137BDAD33C3875546E42201C050FBCCDCF33FFC15EC5B60D09803A25")
        data += bytes.fromhex(
            "004209A161040AB1778E2F2C00EE482F205B28EFBA439FCB04EA283F619478D9")
        data += bytes.fromhex(
            "6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D")
        data += bytes.fromhex(
            "0CF1E6ECE60E49D19BB57C1A432E805F39BB4F65C366741E4F03FA54FBD90714")

        self.assertEqual(
            self.score.get_block_header(
                data,
                bytes.fromhex(
                    "1CCD765C80D0DC1705BB7B6BE616DAD3CF2E6439BB9A9B776D5BD183F89CA141"),
                381837
            ).hex(),
            "a35617a81409ce46f1f820450b8ad4b217d99ae38aaa719b33c4fc52dca99b22"
        )

    # =-=-=-=-=-=-=-=-=-=-=-=-= /BlockHeaderMerkleParts =-=-=-=-=-=-=-=-=-=-=-=-=

    # =-=-=-=-=-=-=-=-=-=-=-=-= \TMSignature =-=-=-=-=-=-=-=-=-=-=-=-=

    def test_recover_signer(self):
        block_hash: bytes = bytes.fromhex(
            "8D1897D04B6B4746021EAF4BF80F0FF9E5A4DDDD451A2FA814DA1C21380D69F8")

        self.assertEqual(
            self.score.recover_signer(
                bytes.fromhex(
                    "81AC28C67F636974BDC70D1A694BA050652FBAA9AA83A1F8B7B10F84C6BC9171"),
                bytes.fromhex(
                    "5D022B62644E496504FC6AF1DD138544CD311F26A979DF4B4DAF92E60CA0F762"),
                28,
                bytes.fromhex(
                    "6E080211400200000000000022480A20"),
                bytes.fromhex(
                    "12240A206C5235F345A661B3136AB762F761045D297D582440751DEA68DBD6083403A31D10012A0C08AA9C85F50510D581FBA001320962616E64636861696E"),
                block_hash,
            ).hex(),
            # Ethereum Address 652d89a66eb4ea55366c45b1f9acfc8e2179e1c5
            "a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba331db530b76775beb0348c52bb8fc1fdac207525e5689caa01c0af8d2f8f371ec5"
        )

        self.assertEqual(
            self.score.recover_signer(
                bytes.fromhex(
                    "F4ABDF0CB47604292B9B0D9636692A0D5379B646EA3246180004BFEAD2D7CA8A"),
                bytes.fromhex(
                    "1AF744F61921AA03D5327333F654747C928CDD1D324A27FF79181AC8E1F6841E"),
                27,
                bytes.fromhex(
                    "6E080211400200000000000022480A20"),
                bytes.fromhex(
                    "12240A206C5235F345A661B3136AB762F761045D297D582440751DEA68DBD6083403A31D10012A0C08AA9C85F50510FA86D8A701320962616E64636861696E"),
                block_hash,
            ).hex(),
            # Ethereum Address aAA22E077492CbaD414098EBD98AA8dc1C7AE8D9
            "f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479"
        )

        self.assertEqual(
            self.score.recover_signer(
                bytes.fromhex(
                    "4258784CC9659EEC320EA86AB7DD1C41C7BF8E9F22035B9E50FA8B527A6079BE"),
                bytes.fromhex(
                    "35C1D785DA88F2D0D563E3AA64B15B96E7C53D025E85895D37F25D99AD11CA14"),
                27,
                bytes.fromhex(
                    "6E080211400200000000000022480A20"),
                bytes.fromhex(
                    "12240A206C5235F345A661B3136AB762F761045D297D582440751DEA68DBD6083403A31D10012A0C08AA9C85F5051089FE8FA601320962616E64636861696E"),
                block_hash,
            ).hex(),
            # Ethereum Address B956589b6fC5523eeD0d9eEcfF06262Ce84ff260
            "d03708f161d1583f49e4260a42b2b08d3ba186d7803a23cc3acd12f074d9d76f2c67aaeee2bc8d00ff253a326eb04de8ce88eea1224c2bb8ca69296a1c753f53"
        )

    # =-=-=-=-=-=-=-=-=-=-=-=-= /TMSignature =-=-=-=-=-=-=-=-=-=-=-=-=

    # =-=-=-=-=-=-=-=-=-=-=-=-= \IAVLMerklePath =-=-=-=-=-=-=-=-=-=-=-=-=-=-=

    def test_get_parent_hash(self):
        self.assertEqual(
            self.score.get_parent_hash(
                True,
                1,
                2,
                599,
                bytes.fromhex(
                    "1459BBC7DB7FFCCE3DECEEA3DF062968F4E3D2B206B93D59FA936E334B9EC434"),
                bytes.fromhex(
                    "5c67c4993b78e7900b56f86dfe426e30dbf597e152918e8ebd029103fae32905"),
            ).hex(),
            "3f590cfc6a2568ef992f39ace15ca3256b6ee7bc06b81b847455d4190d67d775"
        )

        self.assertEqual(
            self.score.get_parent_hash(
                True,
                2,
                4,
                599,
                bytes.fromhex(
                    "5B70BFADD16EC95409072FB3686BF2BFEC48113F96CACA88397883110C672F13"),
                bytes.fromhex(
                    "3f590cfc6a2568ef992f39ace15ca3256b6ee7bc06b81b847455d4190d67d775"),
            ).hex(),
            "56e60e1bf390a5045077f34cb6f96b91f3f3d79965b23e4d649d62b3f58dfb0a"
        )

        self.assertEqual(
            self.score.get_parent_hash(
                True,
                3,
                8,
                599,
                bytes.fromhex(
                    "709E1C73511B24EFDD9B8D3CD717A5210BA20E2411A8529E8B642C54FB002DC4"),
                bytes.fromhex(
                    "56e60e1bf390a5045077f34cb6f96b91f3f3d79965b23e4d649d62b3f58dfb0a"),
            ).hex(),
            "320fec0587c84c21c33d8c2361c06e0e51354f5b9d353c3f3dc3a123e97cbded"
        )

        self.assertEqual(
            self.score.get_parent_hash(
                True,
                4,
                13,
                663,
                bytes.fromhex(
                    "D546BD0A2922A50FD184831DACA108E742A14C3D339C92A7467D31BCD7B0DAD2"),
                bytes.fromhex(
                    "320fec0587c84c21c33d8c2361c06e0e51354f5b9d353c3f3dc3a123e97cbded"),
            ).hex(),
            "d0ee29edb1a80f80b6dc2c058b07e85846e2a1d4ec49fce1dd0cf1b946ccf456"
        )

    # =-=-=-=-=-=-=-=-=-=-=-=-= /IAVLMerklePath =-=-=-=-=-=-=-=-=-=-=-=-=-=-=

    # =-=-=-=-=-=-=-=-=-=-=-=-= \MultiStore =-=-=-=-=-=-=-=-=-=-=-=-=-=-=

    def test_get_app_hash(self):
        self.assertEqual(
            self.score.get_app_hash(
                bytes.fromhex(
                    "BD3021997D04B26F04056F040DDAB5EB0249862398617FF69A477D836619E2B8") +
                bytes.fromhex(
                    "F8231CF7FDFEDC7A8BC63140965D6A707F6F302C11DB997FC44C0595053479A0") +
                bytes.fromhex(
                    "153AE305B0C588EB5AAD3E685FB740B270FD287ECE08BCC85D668FCB10E3C1F5") +
                bytes.fromhex(
                    "B1F2FD852E790E735CA2D3014F96A2A53C60393E9C6BBF941B9A6DD6A05CF6F9") +
                bytes.fromhex(
                    "BDFCCB8D10C48AEFE3FE6872BC558C3C17CA7927A82D48E00C1F7199758227AC")
            ).hex(),
            "0def6341481c4370d561d546c268fb6afed8520689b70abc615c47cfb2a0eee8"
        )

    # =-=-=-=-=-=-=-=-=-=-=-=-= /MultiStore =-=-=-=-=-=-=-=-=-=-=-=-=-=-=

    # =-=-=-=-=-=-=-=-=-=-=-=-= \Validator =-=-=-=-=-=-=-=-=-=-=-=-=-=-=

    def test_set_validators(self):
        self.assertEqual(
            self.score.get_validator_power(
                bytes.fromhex(
                    "a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba331db530b76775beb0348c52bb8fc1fdac207525e5689caa01c0af8d2f8f371ec5")
            ),
            0
        )
        self.assertEqual(
            self.score.get_validator_power(
                bytes.fromhex(
                    "724ae29cfeb7497051d09edfd8e822352c4c8361b757647645b78c8cc74ce885f04c26ee07ff6ada08587a4037363838b1dda6e306091ee0690caa8fe0e6fcd2")
            ),
            0
        )
        self.assertEqual(
            self.score.get_validator_power(
                bytes.fromhex(
                    "f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479")
            ),
            0
        )

        self.score.set_validators(
            bytes.fromhex("0000000300000040a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba331db530b76775beb0348c52bb8fc1fdac207525e5689caa01c0af8d2f8f371ec5000000000000006400000040724ae29cfeb7497051d09edfd8e822352c4c8361b757647645b78c8cc74ce885f04c26ee07ff6ada08587a4037363838b1dda6e306091ee0690caa8fe0e6fcd2000000000000006400000040f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b704790000000000000064")
        )

        self.assertEqual(
            self.score.get_validator_power(
                bytes.fromhex(
                    "a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba331db530b76775beb0348c52bb8fc1fdac207525e5689caa01c0af8d2f8f371ec5")
            ),
            100
        )
        self.assertEqual(
            self.score.get_validator_power(
                bytes.fromhex(
                    "724ae29cfeb7497051d09edfd8e822352c4c8361b757647645b78c8cc74ce885f04c26ee07ff6ada08587a4037363838b1dda6e306091ee0690caa8fe0e6fcd2")
            ),
            100
        )
        self.assertEqual(
            self.score.get_validator_power(
                bytes.fromhex(
                    "f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479")
            ),
            100
        )

    # =-=-=-=-=-=-=-=-=-=-=-=-= \Validator =-=-=-=-=-=-=-=-=-=-=-=-=-=-=

    # =-=-=-=-=-=-=-=-=-=-=-=-= \Bridge =-=-=-=-=-=-=-=-=-=-=-=-=-=-=

    def test_relay_oracle_state_by_part(self):
        # derive app_hash
        app_hash = self.score.get_app_hash(
            bytes.fromhex(
                "685430546D23A44E6B8034EAAFBC2F4CD7FEF54B54D5B66528CB4E5225BD74FB") +
            bytes.fromhex(
                "4F8D0BB0CD3EB9DC70B4DBBEA5F0CBD5B523195F7BAE02BB401BB00A93ABA08E") +
            bytes.fromhex(
                "1912057FFF0B3E85ABF1A319F75D37B21B430F3DA7DB9E486A5041DE47C686D3") +
            bytes.fromhex(
                "B1F2FD852E790E735CA2D3014F96A2A53C60393E9C6BBF941B9A6DD6A05CF6F9") +
            bytes.fromhex(
                "91CC906286235B676AD402FC04ED768EB2BFECA664E8D595C286571DA1433C60")
        )

        # derive block_hash
        block_hash = self.score.get_block_header(
            bytes.fromhex(
                "32FA694879095840619F5E49380612BD296FF7E950EAFB66FF654D99CA70869E") +
            bytes.fromhex(
                "D9F175396C0E2D0E77F69856ABF5D8E69283CB915EB8886262FEDA1D519B3005") +
            bytes.fromhex(
                "4F4D548668A3986DB253689234B9CAC96303A128B7081B16E53B79AB9E65B887") +
            bytes.fromhex(
                "004209A161040AB1778E2F2C00EE482F205B28EFBA439FCB04EA283F619478D9") +
            bytes.fromhex(
                "6E340B9CFFB37A989CA544E6BB780A2C78901D3FB33738768511A30617AFA01D") +
            bytes.fromhex(
                "D991DA4D4E69473CC75A4B819F9E07D4956671A6F4A74DF4CC16596FCBE68137"),
            app_hash,
            191
        )

        # derive signers
        signers = self.score.recover_signers(
            bytes.fromhex("00000003000000209279257914bfee6faec46df086e9a673a4c1576fb299094d45f77e12fa3728e70000002004dbf23f5ebb07ba8fd7ca06843dfe363b5e86596930e1889d9bd5bd3e98fc531c000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510979fbd9801320962616e64636861696e00000020826bb17b714ebcd8199ee2a01334102f19248087cfdeee42ebd406b3991c389500000020621281eedf97f3a9ec121224cef9f8c07872d2c81cce44dd846bb3678dd505231b000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510dd89e98a01320962616e64636861696e00000020d775fd0e1580499ef16a4ab1998dcb4cbd47cf6f342cdc51a9552d1434552ed8000000205a1a2075ae97a6bbd07a5a40efe287eceab37c7a7caae82e6c997c51c62334701b000000106e080211bf0000000000000022480a200000003f12240a2066424e8f0417945a71067a55b7121282a90524e3c0709d8f8addb6adc8fd46d110012a0c08e6e9c6f70510ab94a59201320962616e64636861696e"),
            block_hash
        )

        self.assertEqual(signers, [
            # Ethereum Address 652D89a66Eb4eA55366c45b1f9ACfc8e2179E1c5
            bytes.fromhex(
                "a54ffaa84c8f2f798782de8b962a84784e288487a747813a0857243a60e2ba331db530b76775beb0348c52bb8fc1fdac207525e5689caa01c0af8d2f8f371ec5"),
            # Ethereum Address 88e1cd00710495EEB93D4f522d16bC8B87Cb00FE
            bytes.fromhex(
                "724ae29cfeb7497051d09edfd8e822352c4c8361b757647645b78c8cc74ce885f04c26ee07ff6ada08587a4037363838b1dda6e306091ee0690caa8fe0e6fcd2"),
            # Ethereum Address aAA22E077492CbaD414098EBD98AA8dc1C7AE8D9
            bytes.fromhex(
                "f57f3997a4e81d8f321e9710927e22c2e6d30fb6d8f749a9e4a07afb3b3b7909caefd2ec5f359903d492bc45026b6b45baafe5ad67974e75d8d3e0bb44b70479"),
        ])

    def test_verify_oracle_data(self):
        self.score.set_oracle_state(
            182,
            bytes.fromhex(
                "B59AD73DB9147F6AC7C88A64B1BAD51C90F8C48B4487ADA9276A323808E56E3E")
        )
        self.assertEqual(
            self.score.verify_oracle_data(
                182,
                bytes.fromhex("000000046265656200000000000000010000000f0000000342544300000000000003e800000000000000040000000000000004000000046265656200000000000000010000000000000004000000005edddcd3000000005edddcd700000001000000080000000000948e69"),
                163,
                bytes.fromhex(
                    "000000060102000000000000000300000000000000b4000000204d4479f8cf02cba65f95231b06eeaa51e99f75a153c3ed28d9a86b565de593060103000000000000000500000000000000b4000000208861ad25f99a677d4934541a99928cdbe18bbf34ee39d754e93cce090be705020104000000000000000900000000000000b4000000205d162c955ae030390cea63ca7ed8ba72b7a49e0bc73daea23cf6b79e8899c49c0105000000000000001000000000000000b40000002008dd24b96c3c9413ba46a7149232c18ad66db68ac04cf2a770b707b6d29fc8f70106000000000000001800000000000000b40000002032d2e1cc89f3fbd139c35b5d5fe6430799aeefd495c882c419a120b7257e3f5a0107000000000000003e00000000000000b500000020691b2504cd253868bccd43774c8b8dbfa1e2d0d41c07c23f2f1237fcf0d2f0d8")
            ),
            {
                'req': {
                    'ask_count': 4,
                    'calldata': b'\x00\x00\x00\x03BTC\x00\x00\x00\x00\x00\x00\x03\xe8',
                    'client_id': 'beeb',
                    'min_count': 4,
                    'oracle_script_id': 1
                },
                'res': {
                    'ans_count': 4,
                    'request_id': 1,
                    'request_time': 1591598291,
                    'resolve_status': 0,
                    'resolve_time': 1591598295,
                    'result': b'\x00\x00\x08\x00\x00\x00\x00\x00\x94\x8ei'
                }
            }
        )

    # =-=-=-=-=-=-=-=-=-=-=-=-= /Bridge =-=-=-=-=-=-=-=-=-=-=-=-=-=-=
