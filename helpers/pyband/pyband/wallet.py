from __future__ import annotations

import hashlib

from bech32 import bech32_encode, convertbits
from bip32 import BIP32
from ecdsa import SigningKey, VerifyingKey, SECP256k1
from mnemonic import Mnemonic


class PrivateKey:
    def __init__(self, signing_key: SigningKey) -> None:
        self.signing_key = signing_key

    @classmethod
    def from_mnemonic(cls, words: str, path="m/44'/494'/0'/0/0") -> PrivateKey:
        seed = Mnemonic("english").to_seed(words)
        return cls(
            SigningKey.from_string(
                BIP32.from_seed(seed).get_privkey_from_path(path),
                curve=SECP256k1,
                hashfunc=hashlib.sha256,
            )
        )

    def to_pubkey(self) -> PublicKey:
        return PublicKey(self.signing_key.get_verifying_key())

    def sign(self, msg: bytes) -> bytes:
        pass


class PublicKey:
    def __init__(self, verify_key: VerifyingKey) -> None:
        self.verify_key = verify_key

    @classmethod
    def from_bech32(cls, bech: str, prefix="bandvalconspub") -> PublicKey:
        pass

    def __str__(self) -> str:
        return self.verify_key.to_string("compressed")

    def to_bech32(self, prefix="bandvalconspub") -> str:
        pass

    def to_address(self) -> Address:
        hash = hashlib.new("sha256", self.__str__()).digest()
        return Address(hashlib.new("ripemd160", hash).digest())

    def verify(self, msg: bytes, sig: bytes) -> bool:
        pass


class Address:
    def __init__(self, addr: bytes) -> None:
        self.addr = addr

    @classmethod
    def from_bech32(cls, bech: str, prefix="band") -> Address:
        pass

    def to_bech32(self, prefix="band") -> str:
        five_bit_r = convertbits(self.addr, 8, 5)
        assert five_bit_r is not None, "Unsuccessful bech32.convertbits call"
        return bech32_encode(prefix, five_bit_r)

    def __str__(self) -> str:
        return self.addr.hex()
