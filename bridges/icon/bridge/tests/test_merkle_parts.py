from ..utils import merkle_part
from tbears.libs.scoretest.score_test_case import ScoreTestCase


class TestBlockHeaderMerkleParts(ScoreTestCase):

    def setUp(self):
        super().setUp()

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
            merkle_part.get_block_header(
                data,
                bytes.fromhex(
                    "1CCD765C80D0DC1705BB7B6BE616DAD3CF2E6439BB9A9B776D5BD183F89CA141"),
                381837
            ).hex(),
            "a35617a81409ce46f1f820450b8ad4b217d99ae38aaa719b33c4fc52dca99b22"
        )
