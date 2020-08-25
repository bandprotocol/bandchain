from .wallet import PrivateKey

TEST_MNEMONIC = "coach pond canoe lake solution empty vacuum term pave toe burst top violin purpose umbrella color disease thrive diamond found track need filter wait"
TEST_ADDRESS = "band13eznuehmqzd3r84fkxu8wklxl22r2qfmtlth8c"


def test_ez():
    assert 2 + 2 == 4


def test_private_key_from_mnemonic():
    priv_key = PrivateKey.from_mnemonic(TEST_MNEMONIC)
    sig = priv_key.sign(b"hello")
    pub_key = priv_key.to_pubkey()
    print(priv_key)
    print(priv_key.sign(b"hello"))
    print(pub_key.verify(b"hellox", sig))


def test_address_from_mnemonic():
    pass
