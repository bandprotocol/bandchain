from ..std_reference_proxy import StdReferenceProxy
from tbears.libs.scoretest.score_test_case import ScoreTestCase
from iconservice.base.exception import IconScoreException


class TestStdReferenceProxy(ScoreTestCase):
    def setUp(self):
        super().setUp()
        self.score = self.get_score_instance(
            StdReferenceProxy, self.test_account1, on_install_params={"_ref": self.test_account1}
        )

    def test_get_ref(self):
        self.assertEqual(self.score.get_ref(), self.test_account1)

