from ..utils import secp256k1
from tbears.libs.scoretest.score_test_case import ScoreTestCase


class TestSecp256k1(ScoreTestCase):

    def setUp(self):
        super().setUp()

    # https://rosettacode.org/wiki/Cipolla%27s_algorithm#Python
    def test_ecc_sqrt(self):
        self.assertEqual((4, 3), secp256k1.ecc_sqrt(2, 7))
        self.assertEqual((9872, 135), secp256k1.ecc_sqrt(8218, 10007))
        self.assertEqual((37, 64), secp256k1.ecc_sqrt(56, 101))
        self.assertEqual((1, 10), secp256k1.ecc_sqrt(1, 11))
        self.assertEqual((), secp256k1.ecc_sqrt(8219, 10007))

    def test_ecrecover(self):
        self.assertEqual(
            "ba9231dbfdaa0c5c62c320402f9136ced720a585244c4daacfade4b15614c4dcd13fcc4241a5dbdb10865791e22a2c639ce5673ff6213c018f7a93808ab5b9b7",
            secp256k1.ecrecover(
                bytes.fromhex(
                    "45298f64c176ea53805cec6a6bff4bdd3fccb40250d62dc8f954aab6ef9037cb"),
                bytes.fromhex(
                    "10ce9c6360ed05342478a0ed2ec6d378adf5a1d437a2be3c8e89740b3cd3bff1"),
                bytes.fromhex(
                    "aa4b0588903d649c20d4d924998cb8be13e1f4ffe7e11ad8f84f80a8b2771fc9"),
                27
            ).hex()
        )

        self.assertEqual(
            "dbefccd70f6d5723d7761b0104038d7a60c0f447d05404ad5f9bd43ef30d47a23497e35f5bf02201aae875b3f37eefa055932260c537d2a9e4d674c284627965",
            secp256k1.ecrecover(
                bytes.fromhex(
                    "966404b77bf1454fcb89da5e1193f9499a0a85d1837fc98584457cbca6f36ed6"),
                bytes.fromhex(
                    "1487239876fde6b58d23cda54c86035ae9e6904c82846b7fb5dcb1ace71a77b3"),
                bytes.fromhex(
                    "88a3fa74e6bc5210204e80084e5fb0dbbc055d996a3af3cc2bca51bdf41d5cc1"),
                28
            ).hex()
        )
