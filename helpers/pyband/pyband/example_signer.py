from .transaction import Transaction


class Signer(object):
    def sign(self, msg):
        raise NotImplementedError()


class PrivateKey(Signer):
    self.signingkey

    def sign(self, data):
        # same as in wallet

class Ledger(Signer):


    def sign(msg):
        # sign by ledger