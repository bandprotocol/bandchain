from .wallet import PrivateKey

TEST_MNEMONIC = "coach pond canoe lake solution empty vacuum term pave toe burst top violin purpose umbrella color disease thrive diamond found track need filter wait"
TEST_ACC_ADDRESS = "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"
TEST_ACC_PUBKEY = "bandpub1addwnpepqdg7nrsmuztj2re07svgcz4vuzn3de56nykdwlualepkk05txs5q6mw8s9v"
TEST_VAL_ADDRESS = "bandvaloper13eznuehmqzd3r84fkxu8wklxl22r2qfm8f05zn"
TEST_VAL_PUBKEY = (
    "bandvaloperpub1addwnpepqdg7nrsmuztj2re07svgcz4vuzn3de56nykdwlualepkk05txs5q69gsm29"
)
TEST_CONS_ADDRESS = "bandvalcons13eznuehmqzd3r84fkxu8wklxl22r2qfmn6ugwj"
TEST_CONS_PUBKEY = (
    "bandvalconspub1addwnpepqdg7nrsmuztj2re07svgcz4vuzn3de56nykdwlualepkk05txs5q6r8ytws"
)


def test_ez():
    assert 2 + 2 == 4


def test_private_key_from_mnemonic():
    priv_key = PrivateKey.from_mnemonic(TEST_MNEMONIC)
    sig = priv_key.sign(b"hello")
    pub_key = priv_key.to_pubkey()

    assert pub_key.verify(b"hello", sig) == True


def test_public_key_from_mnemonic():
    priv_key = PrivateKey.from_mnemonic(TEST_MNEMONIC)

    assert (
        str(priv_key.to_pubkey())
        == "0351e98e1be097250f2ff4188c0aace0a716e69a992cd77f9dfe436b3e8b34280d"
    )

    assert priv_key.to_pubkey().to_acc_bech32() == TEST_ACC_PUBKEY
    assert priv_key.to_pubkey().to_val_bech32() == TEST_VAL_PUBKEY
    assert priv_key.to_pubkey().to_cons_bech32() == TEST_CONS_PUBKEY


def test_address_from_mnemonic():
    priv_key = PrivateKey.from_mnemonic(TEST_MNEMONIC)

    assert priv_key.to_pubkey().to_address().to_acc_bech32() == TEST_ACC_ADDRESS
    assert priv_key.to_pubkey().to_address().to_val_bech32() == TEST_VAL_ADDRESS
    assert priv_key.to_pubkey().to_address().to_cons_bech32() == TEST_CONS_ADDRESS
