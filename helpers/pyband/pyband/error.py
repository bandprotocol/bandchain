class EmptyRequestMsgError(Exception):
    def __init__(self, message: str):
        super().__init__(message)


class QueryError(Exception):
    def __init__(self, message: str):
        super().__init__(message)


class NegativeIntegerError(Exception):
    def __init__(self, message: str):
        super().__init__(message)


class ValueTooLargeError(Exception):
    def __init__(self, message: str):
        super().__init__(message)


class InsufficientCoinError(Exception):
    def __init__(self, message: str):
        super().__init__(message)


class EmptyMsgError(Exception):
    def __init__(self, message: str):
        super().__init__(message)


class NotFoundError(Exception):
    def __init__(self, message: str):
        super().__init__(message)


class UndefinedError(Exception):
    def __init__(self, message: str):
        super().__init__(message)


class DecodeError(Exception):
    def __init__(self, message: str):
        super().__init__(message)


class ConvertError(Exception):
    def __init__(self, message: str):
        super().__init__(message)


class UnsuccessfulCallError(Exception):
    def __init__(self, message: str):
        super().__init__(message)


class CreateError(Exception):
    def __init__(self, message: str):
        super().__init__(message)


class SchemaError(Exception):
    def __init__(self, message: str):
        super().__init__(message)
