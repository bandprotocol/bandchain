from dacite import from_dict

from .auth import Auth

VALIDATOR_TEST = "bandvaloper13eznuehmqzd3r84fkxu8wklxl22r2qfm8f05zn"
# REPORTER_PUBKEY =


def test_get_msg_sign_bytes():
    assert Auth.get_msg_sign_bytes("bandchain", VALIDATOR_TEST, "3", "1") == bytes.fromhex(
        "7b22636861696e5f6964223a2262616e64636861696e222c2265787465726e616c5f6964223a2231222c22726571756573745f6964223a2233222c2276616c696461746f72223a2262616e6476616c6f7065723133657a6e7565686d717a6433723834666b787538776b6c786c3232723271666d386630357a6e227d"
    )


# def test_verify_verification_message():
#     assert Auth.verify_verification_message("bandchain", VALIDATOR_TEST, "3", "1",)

