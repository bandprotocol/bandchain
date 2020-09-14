from ..utils import sha256
from tbears.libs.scoretest.score_test_case import ScoreTestCase


class TestSha256(ScoreTestCase):

    def setUp(self):
        super().setUp()

    # https://emn178.github.io/online-tools/sha256.html
    def test_sha256(self):
        self.assertEqual(
            bytes.fromhex(
                "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"),
            sha256.digest(b'')
        )

        self.assertEqual(
            bytes.fromhex(
                "a64089b957b8a60cf458d81d1fe35ef56909aaf8deb8f595d365c7bda4edab02"),
            sha256.digest(b'beeb')
        )

        self.assertEqual(
            bytes.fromhex(
                "d7b553c6f09ac85d142415f857c5310f3bbbe7cdd787cce4b985acedd585266f"),
            sha256.digest(b'just a test string')
        )

        self.assertEqual(
            bytes.fromhex(
                "8113ebf33c97daa9998762aacafe750c7cefc2b2f173c90c59663a57fe626f21"),
            sha256.digest(b'just a test string' * 7)
        )

        self.assertEqual(
            bytes.fromhex(
                "4bc3a0e290537de3f6cf9203fc632bea3aa6468d690deb1a5444f2d2856a7086"),
            sha256.digest(b'mumu lulu fufu momo toto coco bobo photo')
        )

        self.assertEqual(
            bytes.fromhex(
                "d3eb3ec84a2ee7f0a18f3c1e0bb95c71d57032ce4cf3715939d923bca1510a6b"),
            sha256.digest(
                bytes('▁ ▂ ▄ ▅ ▆ ▇ █   🎀  ¸„.-•~¹°”ˆ˜¨ 𝔡𝔣𝔰𝔡𝔧𝔬𝔣𝔞 ¨˜ˆ”°¹~•-.„¸  🎀   █', 'utf-8'))
        )

        self.assertEqual(
            bytes.fromhex(
                "e560a602b5f089d690011f196fd67319f9d0b50039bf31d04e35fa32dd2d026f"),
            sha256.digest(
                bytes('¸҉„҉.҉-҉•҉~҉¹҉°҉”҉ˆ҉˜҉¨҉ ҉�҉�҉�҉�҉�҉�҉�҉�҉�҉�҉�҉�҉�҉�҉�҉�҉ ҉¨҉˜҉ˆ҉”҉°҉¹҉~҉•҉-҉.҉„҉¸҉', 'utf-8'))
        )

        self.assertEqual(
            bytes.fromhex(
                "8e69d4fa7fe33ea3f33e112bcb2f57fb48e15aff09c505c220164f07654ef12c"),
            sha256.digest(b'Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum')
        )
