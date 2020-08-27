from __future__ import annotations

import hashlib

from bech32 import bech32_encode, convertbits
from bip32 import BIP32
from ecdsa import SigningKey, VerifyingKey, SECP256k1
from ecdsa.util import sigencode_string_canonize
from mnemonic import Mnemonic

BECH32_PUBKEY_ACC_PREFIX = "bandpub"
BECH32_PUBKEY_VAL_PREFIX = "bandvalpub"
BECH32_PUBKEY_CONS_PREFIX = "bandvalconspub"

BECH32_ADDR_ACC_PREFIX = "band"
BECH32_ADDR_VAL_PREFIX = "bandvaloper"
BECH32_ADDR_CONS_PREFIX = "bandvalcons"


class PrivateKey:
    """
    Class for wrapping SigningKey using for signature creation and public key derivation.

    :ivar signing_key: the ecdsa SigningKey instance
    :vartype signing_key: ecdsa.SigningKey
    """

    def __init__(self, _error_do_not_use_init_directly=None) -> None:
        """Unsupported, please use from_mnemonic to initialize."""
        if not _error_do_not_use_init_directly:
            raise TypeError("Please use PrivateKey.from_mnemonic() to construct me")
        self.signing_key = None

    @classmethod
    def from_mnemonic(cls, words: str, path="m/44'/494'/0'/0/0") -> PrivateKey:
        """
        Create a PrivateKey instance from a given mnemonic phrase and a HD derivation path.
        If path is not given, default to Band's HD prefix 494 and all other indexes being zeroes.

        :param words: the mnemonic phrase for recover private key
        :param path: the HD path that follows the BIP32 standard

        :return: Initialized PrivateKey object
        """
        seed = Mnemonic("english").to_seed(words)
        self = cls(_error_do_not_use_init_directly=True)
        self.signing_key = SigningKey.from_string(
            BIP32.from_seed(seed).get_privkey_from_path(path),
            curve=SECP256k1,
            hashfunc=hashlib.sha256,
        )
        return self

    def to_pubkey(self) -> PublicKey:
        """
        Return the PublicKey associated with this private key.

        :return: a PublicKey that can be used to verify the signatures made with this PrivateKey
        """
        pubkey = PublicKey(_error_do_not_use_init_directly=True)
        pubkey.verify_key = self.signing_key.get_verifying_key()
        return pubkey

    def sign(self, msg: bytes) -> bytes:
        """
        Sign and the given message using the edcsa sign_deterministic function.

        :param msg: the message that will be hashed and signed

        :return: a signature of this private key over the given message
        """
        return self.signing_key.sign_deterministic(
            msg, hashfunc=hashlib.sha256, sigencode=sigencode_string_canonize,
        )


class PublicKey:
    """
    Class for wraping VerifyKey using for signature verification. Adding method to encode/decode
    to Bech32 format.

    :ivar verify_key: the ecdsa VerifyingKey instance
    :vartype verify_key: ecdsa.VerifyingKey
    """

    def __init__(self, _error_do_not_use_init_directly=None) -> None:
        """Unsupported, please do not contruct it directly."""
        if not _error_do_not_use_init_directly:
            raise TypeError("Please use PublicKey's factory methods to construct me")
        self.verify_key = None

    @classmethod
    def _from_bech32(cls, bech: str, prefix: str) -> PublicKey:
        pass

    @classmethod
    def from_acc_pub(cls, bech: str) -> PublicKey:
        return cls._from_bech32(bech, BECH32_PUBKEY_ACC_PREFIX)

    @classmethod
    def from_val_pub(cls, bech: str) -> PublicKey:
        return cls._from_bech32(bech, BECH32_PUBKEY_VAL_PREFIX)

    @classmethod
    def from_val_cons_pub(cls, bech: str) -> PublicKey:
        return cls._from_bech32(bech, BECH32_PUBKEY_CONS_PREFIX)

    def __str__(self) -> str:
        """Return hex-format of compressed pubkey"""
        return self.verify_key.to_string("compressed").hex()

    def _to_bech32(self, prefix: str) -> str:

        five_bit_r = convertbits(self.verify_key.to_string("compressed"), 8, 5)
        assert five_bit_r is not None, "Unsuccessful bech32.convertbits call"
        return bech32_encode(prefix, five_bit_r)

    def to_acc_bech32(self) -> str:
        """Return bech32-encoded with account public key prefix"""
        return self._to_bech32(BECH32_PUBKEY_ACC_PREFIX)

    def to_val_bech32(self) -> str:
        """Return bech32-encoded with validator public key prefix"""
        return self._to_bech32(BECH32_PUBKEY_VAL_PREFIX)

    def to_cons_bech32(self) -> str:
        """Return bech32-encoded with validator consensus public key prefix"""
        return self._to_bech32(BECH32_PUBKEY_CONS_PREFIX)

    def to_address(self) -> Address:
        """Return address instance from this public key"""
        hash = hashlib.new("sha256", self.verify_key.to_string("compressed")).digest()
        return Address(hashlib.new("ripemd160", hash).digest())

    def verify(self, msg: bytes, sig: bytes) -> bool:
        """
        Verify a signature made over provided data.

        :param msg: data signed by the `signature`, will be hashed using sha256 function
        :param sig: encoding of the signature

        :raises BadSignatureError: if the signature is invalid or malformed
        :return: True if the verification was successful
        """
        return self.verify_key.verify(sig, msg, hashfunc=hashlib.sha256)


class Address:
    def __init__(self, addr: bytes) -> None:
        self.addr = addr

    # TODO: Split into acc, val, cons address.
    @classmethod
    def from_bech32(cls, bech: str, prefix="band") -> Address:
        pass

    def _to_bech32(self, prefix: str) -> str:
        five_bit_r = convertbits(self.addr, 8, 5)
        assert five_bit_r is not None, "Unsuccessful bech32.convertbits call"
        return bech32_encode(prefix, five_bit_r)

    def to_acc_bech32(self) -> str:
        """Return bech32-encoded with account address prefix"""
        return self._to_bech32(BECH32_ADDR_ACC_PREFIX)

    def to_val_bech32(self) -> str:
        """Return bech32-encoded with validator address prefix"""
        return self._to_bech32(BECH32_ADDR_VAL_PREFIX)

    def to_cons_bech32(self) -> str:
        """Return bech32-encoded with validator consensus address key prefix"""
        return self._to_bech32(BECH32_ADDR_CONS_PREFIX)

    def __str__(self) -> str:
        return self.addr.hex()
