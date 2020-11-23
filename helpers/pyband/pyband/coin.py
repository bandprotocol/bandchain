class Coin(object):
    def __init__(self, amount: int, denom: str) -> None:
        self.amount = amount
        self.denom = denom

    def __eq__(self, o: "Coin") -> bool:
        return self.amount == o.amount and self.denom == o.denom

    @classmethod
    def from_json(cls, coin) -> "Coin":
        return cls(int(coin["amount"]), coin["denom"])

    def as_json(self) -> dict:
        return {"amount": str(self.amount), "denom": self.denom}

    def validate(self) -> bool:
        if self.amount < 0:
            raise ValueError("Expect amount more than 0")

        if len(self.denom) == 0:
            raise ValueError("Expect denom")

        return True