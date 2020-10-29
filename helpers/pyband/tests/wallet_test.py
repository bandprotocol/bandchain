from pyband.wallet import PrivateKey, PublicKey, Address

TEST_MNEMONIC = "coach pond canoe lake solution empty vacuum term pave toe burst top violin purpose umbrella color disease thrive diamond found track need filter wait"

TEST_ACC_PUBKEY = (
    "bandpub1addwnpepqdg7nrsmuztj2re07svgcz4vuzn3de56nykdwlualepkk05txs5q6mw8s9v"
)
TEST_VAL_PUBKEY = (
    "bandvaloperpub1addwnpepqdg7nrsmuztj2re07svgcz4vuzn3de56nykdwlualepkk05txs5q69gsm29"
)
TEST_CONS_PUBKEY = (
    "bandvalconspub1addwnpepqdg7nrsmuztj2re07svgcz4vuzn3de56nykdwlualepkk05txs5q6r8ytws"
)

TEST_ACC_ADDRESS = "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"
TEST_VAL_ADDRESS = "bandvaloper13eznuehmqzd3r84fkxu8wklxl22r2qfm8f05zn"
TEST_CONS_ADDRESS = "bandvalcons13eznuehmqzd3r84fkxu8wklxl22r2qfmn6ugwj"


def test_private_key_from_mnemonic():
    assert (
        PrivateKey.from_mnemonic(TEST_MNEMONIC).to_hex()
        == "2159f40dda9e4c9d8ed9d6f8c353e247a2658993a9d53a94ff17410cd0ea471d"
    )
    assert (
        PrivateKey.from_mnemonic(TEST_MNEMONIC, "m/44'/494'/0'/0/1").to_hex()
        == "987af53f91a09926274e5a2ef86223356112056f61b35a57df345d7b14176bb3"
    )


def test_private_key_from_hex():
    assert (
        PrivateKey.from_hex(
            "2159f40dda9e4c9d8ed9d6f8c353e247a2658993a9d53a94ff17410cd0ea471d"
        ).to_hex()
        == "2159f40dda9e4c9d8ed9d6f8c353e247a2658993a9d53a94ff17410cd0ea471d"
    )


def test_private_key_to_public_key():
    pub_key = PrivateKey.from_mnemonic(TEST_MNEMONIC).to_pubkey()
    assert (
        pub_key.to_hex()
        == "0351e98e1be097250f2ff4188c0aace0a716e69a992cd77f9dfe436b3e8b34280d"
    )


def test_private_key_sign():
    priv_key = PrivateKey.from_mnemonic(TEST_MNEMONIC)
    sig = priv_key.sign(b"test msg")

    assert (
        sig.hex()
        == "42a1e41012155ae2daa9b9a2e038f76463da4662564b4989f236ecb4d97f592c1190d92319363e2d1eb78fb98f0dac30c5e2a850f45bb4c44f1f6203ebe6efbe"
    )


def test_public_key_from_bech32():
    assert (
        PublicKey.from_acc_bech32(TEST_ACC_PUBKEY).to_hex()
        == "0351e98e1be097250f2ff4188c0aace0a716e69a992cd77f9dfe436b3e8b34280d"
    )
    assert (
        PublicKey.from_val_bech32(TEST_VAL_PUBKEY).to_hex()
        == "0351e98e1be097250f2ff4188c0aace0a716e69a992cd77f9dfe436b3e8b34280d"
    )
    assert (
        PublicKey.from_cons_bech32(TEST_CONS_PUBKEY).to_hex()
        == "0351e98e1be097250f2ff4188c0aace0a716e69a992cd77f9dfe436b3e8b34280d"
    )


def test_public_key_to_bech32():
    priv_key = PrivateKey.from_mnemonic(TEST_MNEMONIC)

    assert priv_key.to_pubkey().to_acc_bech32() == TEST_ACC_PUBKEY
    assert priv_key.to_pubkey().to_val_bech32() == TEST_VAL_PUBKEY
    assert priv_key.to_pubkey().to_cons_bech32() == TEST_CONS_PUBKEY


def test_public_key_verify():
    priv_key = PrivateKey.from_mnemonic(TEST_MNEMONIC)
    sig = priv_key.sign(b"test msg")

    pub_key = priv_key.to_pubkey()
    assert pub_key.verify(b"test msg", sig) == True
    # Invalid message
    assert pub_key.verify(b"another msg", sig) == False
    # Invalid signature
    assert pub_key.verify(b"test msg", sig[1:]) == False


def test_address_from_bech32():
    assert (
        Address.from_acc_bech32(TEST_ACC_ADDRESS).to_hex()
        == "8e453e66fb009b119ea9b1b8775be6fa9435013b"
    )
    assert (
        Address.from_val_bech32(TEST_VAL_ADDRESS).to_hex()
        == "8e453e66fb009b119ea9b1b8775be6fa9435013b"
    )
    assert (
        Address.from_cons_bech32(TEST_CONS_ADDRESS).to_hex()
        == "8e453e66fb009b119ea9b1b8775be6fa9435013b"
    )


def test_address_to_bech32():
    addr = PrivateKey.from_mnemonic(TEST_MNEMONIC).to_pubkey().to_address()

    assert addr.to_acc_bech32() == TEST_ACC_ADDRESS
    assert addr.to_val_bech32() == TEST_VAL_ADDRESS
    assert addr.to_cons_bech32() == TEST_CONS_ADDRESS
