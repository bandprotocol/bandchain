from ..utils import tm_signature
from tbears.libs.scoretest.score_test_case import ScoreTestCase


class TestTMSignature(ScoreTestCase):

    def setUp(self):
        super().setUp()

    def test_recover_signer(self):
        block_hash: bytes = bytes.fromhex(
            "8D1897D04B6B4746021EAF4BF80F0FF9E5A4DDDD451A2FA814DA1C21380D69F8")

        self.assertEqual(
            tm_signature.recover_signer(
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
            tm_signature.recover_signer(
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
            tm_signature.recover_signer(
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
