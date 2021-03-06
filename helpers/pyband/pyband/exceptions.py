class PybandError(Exception):
    pass


class EmptyRequestMsgError(PybandError):
    pass


class QueryError(PybandError):
    pass


class NegativeIntegerError(PybandError):
    pass


class ValueTooLargeError(PybandError):
    pass


class InsufficientCoinError(PybandError):
    pass


class EmptyMsgError(PybandError):
    pass


class NotFoundError(PybandError):
    pass


class UndefinedError(PybandError):
    pass


class DecodeError(PybandError):
    pass


class ConvertError(PybandError):
    pass


class UnsuccessfulCallError(PybandError):
    pass


class CreateError(PybandError):
    pass


class SchemaError(PybandError):
    pass
