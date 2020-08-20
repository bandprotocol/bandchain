from __future__ import annotations

import hashlib

from bech32 import bech32_encode, convertbits
from bip32 import BIP32
from ecdsa import SigningKey, VerifyingKey, SECP256k1
from mnemonic import Mnemonic

BECH32_PUBKEY_ACC_PREFIX = "bandaccpub"
BECH32_PUBKEY_VAL_PREFIX = "bandvalpub"
BECH32_PUBKEY_CONS_PREFIX = "bandvalconspub"


class PrivateKey:
    """
    Class for wraping SigningKey using for signature creation.

    :ivar signing_key: the ecdsa SigningKey instance
    :vartype signing_key: ecdsa.SigningKey
    """

    def __init__(self, _error__please_use_generate=None) -> None:
        """Unsupported, please use from_mnemonic to initialise."""
        if not _error__please_use_generate:
            raise TypeError("Please use SigningKey.from_mnemonic() to construct me")
        self.signing_key = None

    @classmethod
    def from_mnemonic(cls, words: str, path="m/44'/494'/0'/0/0") -> PrivateKey:
        """
        Create a Private key instance that wraps SigningKey from a given mnemonic phrase and path.

        :param words: The mnemonic phrase for recover private key
        :type words: str
        :param path: The path using for generate private key follow BIP32 standard the default value
        is Band prefix on index 0
        :type path: str

        :return: Initialised PrivateKey object
        :rtype: PrivateKey
        """
        seed = Mnemonic("english").to_seed(words)
        self = cls(_error__please_use_generate=True)
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
        :rtype: VerifyingKey
        """
        return PublicKey(self.signing_key.get_verifying_key())

    def sign(self, msg: bytes) -> bytes:
        """
        Create signature over data using the ecdsa sign_deterministic function to get the signature.

        :param msg: msg that will be hashed for signing
        :type msg: bytes like object

        :return: encoded signature of the hash of `msg`
        :rtype: bytes
        """
        pass


class PublicKey:
    """
    Class for wraping VerifyKey using for signature verification. Adding method to encode/decode
    to Bech32 format.

    :ivar verify_key: the ecdsa VerifyingKey instance
    :vartype verify_key: ecdsa.VerifyingKey
    """

    def __init__(self, verify_key: VerifyingKey) -> None:
        self.verify_key = verify_key

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
        return self.verify_key.to_string("compressed")

    def _to_bech32(self, prefix: str) -> str:
        pass

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
        hash = hashlib.new("sha256", self.__str__()).digest()
        return Address(hashlib.new("ripemd160", hash).digest())

    def verify(self, msg: bytes, sig: bytes) -> bool:
        """
        Verify a signature made over provided data.

        :param msg: data signed by the `signature`, will be hashed using sha256 function
        :type msg: bytes like object
        :param sig: encoding of the signature
        :type sig: bytes like object

        :raises BadSignatureError: if the signature is invalid or malformed
        :return: True if the verification was successful
        :rtype: bool
        """
        pass


class Address:
    def __init__(self, addr: bytes) -> None:
        self.addr = addr

    # TODO: Split into acc, val, cons address.
    @classmethod
    def from_bech32(cls, bech: str, prefix="band") -> Address:
        pass

    def to_bech32(self, prefix="band") -> str:
        five_bit_r = convertbits(self.addr, 8, 5)
        assert five_bit_r is not None, "Unsuccessful bech32.convertbits call"
        return bech32_encode(prefix, five_bit_r)

    def __str__(self) -> str:
        return self.addr.hex()
