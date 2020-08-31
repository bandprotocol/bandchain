from ..utils import utils
from tbears.libs.scoretest.score_test_case import ScoreTestCase


class TestUtils(ScoreTestCase):

    def setUp(self):
        super().setUp()

    def test_merkle_leaf_hash(self):
        self.assertEqual(
            utils.merkle_leaf_hash(
                bytes.fromhex('08d1082cc8d85a0833da8815ff1574675c415760e0aff7fb4e32de6de27faf86')).hex(),
            "35b401b2a74452d2252df60574e0a6c029885965ae48f006ebddc18e53427e26"
        )

    def test_merkle_inner_hash(self):
        self.assertEqual(
            utils.merkle_inner_hash(
                bytes.fromhex(
                    'aa58d2b0be1dbadf21619645c9245062a283c19a6701bfe07d07cc2441687c85'),
                bytes.fromhex(
                    'ca978112ca1bbdcafac231b39a23dc4da786eff8147c4e72b9807785afee48bb')
            ).hex(),
            "da32ca60aaddf58c20ff66ba1bb73ef126fea41822fedbfc10fbbedc8587ffe2"
        )

        self.assertEqual(
            utils.merkle_inner_hash(
                bytes.fromhex(
                    '08d1082cc8d85a0833da8815ff1574675c415760e0aff7fb4e32de6de27faf86'),
                bytes.fromhex(
                    '789411d15a12768a9c3eb99d3453d6ebb4481c2a03ab59cc262a97e25757afe6')
            ).hex(),
            "ca48b611419f0848bf0fce9750caca6fd4fb8ef96ba8d7d3ccd4f05bf2af1661"
        )

    def test_encode_varint_unsigned(self):
        self.assertEqual(utils.encode_varint_unsigned(116).hex(), "74")
        self.assertEqual(utils.encode_varint_unsigned(14947).hex(), "e374")
        self.assertEqual(utils.encode_varint_unsigned(
            244939043).hex(), "a3f2e574")

    def test_encode_varint_signed(self):
        self.assertEqual(utils.encode_varint_signed(58).hex(), "74")
        self.assertEqual(utils.encode_varint_signed(7473).hex(), "e274")
        self.assertEqual(utils.encode_varint_signed(
            122469521).hex(), "a2f2e574")
