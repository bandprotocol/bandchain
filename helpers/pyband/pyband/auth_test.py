import base64

from .auth import Auth

VALIDATOR_TEST = "bandvaloper1p40yh3zkmhcv0ecqp3mcazy83sa57rgjde6wec"


def test_get_msg_sign_bytes():
    assert Auth.get_msg_sign_bytes("bandchain", VALIDATOR_TEST, "1", "1") == bytes.fromhex(
        "7b22636861696e5f6964223a2262616e64636861696e222c2265787465726e616c5f6964223a2231222c22726571756573745f6964223a2231222c2276616c696461746f72223a2262616e6476616c6f706572317034307968337a6b6d6863763065637170336d63617a7938337361353772676a646536776563227d"
    )


def test_verify_verification_message():
    assert (
        Auth.verify_verification_message(
            "bandchain",
            VALIDATOR_TEST,
            "1",
            "1",
            "bandpub1addwnpepqgugvxy0ueqwfmlzh2ta5at2lumcy4wpzzjs4hjz8j44lrdcryqs66wh3rp",
            base64.b64decode(
                "IsgagGxxSVHOPyzProTYBW9sFNMjLGkuDm+JvLgBH8Ux6GMpj3p6e5YGY8KRVWV3fdYWm/UBZdpVqsMbnpV6PQ=="
            ),
        )
        == True
    )

