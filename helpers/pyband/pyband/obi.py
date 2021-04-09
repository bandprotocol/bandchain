from typing import Any, Tuple
from .exceptions import SchemaError, DecodeError


class PyObiSpec(object):
    impls = []

    def __init_subclass__(cls, **kwargs):
        super().__init_subclass__(**kwargs)
        cls.impls.append(cls)

    @classmethod
    def from_spec(cls, spec: str) -> Any:
        for impl in cls.impls:
            if impl.match_schema(spec):
                return impl(spec)
        raise SchemaError("Cannot parse spec: {}".format(spec))

    def __init__(self, spec):
        raise NotImplementedError()

    @classmethod
    def match_schema(cls, schema):
        raise NotImplementedError()

    def encode(self, value):
        raise NotImplementedError()

    def decode(self, data):
        raise NotImplementedError()


class PyObiInteger(PyObiSpec):
    def __init__(self, spec):
        self.is_signed = spec[0] == "i"
        self.size_in_bytes = int(spec[1:]) // 8

    @classmethod
    def match_schema(cls, schema: str) -> bool:
        return schema[:1] in ["i", "u"] and schema[1:] in ["8", "16", "32", "64", "128", "256"]

    def encode(self, value: int) -> bytes:
        return value.to_bytes(self.size_in_bytes, byteorder="big", signed=self.is_signed)

    def decode(self, data: bytes) -> Tuple[int, bytes]:
        return (
            int.from_bytes(data[: self.size_in_bytes], byteorder="big", signed=self.is_signed),
            data[self.size_in_bytes :],
        )


class PyObiBool(PyObiSpec):
    def __init__(self, spec="bool"):
        pass

    @classmethod
    def match_schema(cls, schema: str) -> bool:
        return schema == "bool"

    def encode(self, value: bool) -> bytes:
        return PyObiInteger("u8").encode(1 if value else 0)

    def decode(self, data: bytes) -> Tuple[bool, bytes]:
        u8, remaining = PyObiInteger("u8").decode(data)
        if u8 == 1:
            return True, remaining
        elif u8 == 0:
            return False, remaining
        raise ValueError("Boolean value must be 1 or 0 but got {}".format(u8))


class PyObiArray(PyObiSpec):
    def __init__(self, _spec):
        [spec, size] = _spec[1:-1].rsplit(";", 1)
        self.size = int(size, 10)
        self.intl_obi = self.from_spec(spec)

    @classmethod
    def match_schema(cls, schema: str) -> bool:
        if not (schema[0] == "[" and schema[-1] == "]"):
            return False
        try:
            [spec, size] = schema[1:-1].rsplit(";", 1)
        except:
            return False

        if not size.isdigit():
            return False

        for impl in cls.impls:
            if impl.match_schema(spec):
                return True

        return False

    def encode(self, value: list) -> bytes:
        if len(value) != self.size:
            raise ValueError("array size should be {} but got {}".format(self.size, len(value)))
        result = b""
        for each in value:
            result = result + self.intl_obi.encode(each)
        return result

    def decode(self, data: bytes) -> Tuple[list, bytes]:
        remaining = data[:]
        result = []
        for _ in range(self.size):
            each, remaining = self.intl_obi.decode(remaining)
            result.append(each)
        return result, remaining


class PyObiVector(PyObiSpec):
    def __init__(self, spec):
        self.intl_obi = self.from_spec(spec[1:-1])

    @classmethod
    def match_schema(cls, schema: str) -> bool:
        return schema[0] == "[" and schema[-1] == "]"

    def encode(self, value: list) -> bytes:
        result = PyObiInteger("u32").encode(len(value))
        for each in value:
            result = result + self.intl_obi.encode(each)
        return result

    def decode(self, data: bytes) -> Tuple[list, bytes]:
        length, remaining = PyObiInteger("u32").decode(data)
        result = []
        for _ in range(length):
            each, remaining = self.intl_obi.decode(remaining)
            result.append(each)
        return result, remaining


class PyObiStruct(PyObiSpec):
    def __init__(self, spec):
        self.intl_obi_kvs = []
        fields = [""]
        curly_count = 0
        for c in spec[1:-1]:
            if c == "," and curly_count == 0:
                fields.append("")
            else:
                fields[-1] = fields[-1] + c
                if c == "{":
                    curly_count += 1
                if c == "}":
                    curly_count -= 1
        for each in fields:
            tokens = each.split(":", 1)
            if len(tokens) != 2:
                raise ValueError("Expect at least one colon for each struct field")
            self.intl_obi_kvs.append((tokens[0], self.from_spec(tokens[1])))

    @classmethod
    def match_schema(cls, schema: str) -> bool:
        return schema[0] == "{" and schema[-1] == "}"

    def encode(self, value: dict) -> bytes:
        result = b""
        for key, spec in self.intl_obi_kvs:
            result = result + spec.encode(value[key])
        return result

    def decode(self, data: bytes) -> Tuple[dict, bytes]:
        result = {}
        for key, spec in self.intl_obi_kvs:
            result[key], data = spec.decode(data)
        return result, data


class PyObiString(PyObiSpec):
    def __init__(self, spec="string"):
        pass

    @classmethod
    def match_schema(cls, schema: str) -> bool:
        return schema == "string"

    def encode(self, value: str) -> bytes:
        return PyObiInteger("u32").encode(len(value)) + value.encode()

    def decode(self, data: bytes) -> Tuple[str, bytes]:
        length, remaining = PyObiInteger("u32").decode(data)
        return remaining[:length].decode(), remaining[length:]


class PyObiBytes(PyObiSpec):
    def __init__(self, spec="bytes"):
        pass

    @classmethod
    def match_schema(cls, schema: str) -> bool:
        return schema == "bytes"

    def encode(self, value: bytes) -> bytes:
        return PyObiInteger("u32").encode(len(value)) + value

    def decode(self, data: bytes) -> Tuple[bytes, bytes]:
        length, remaining = PyObiInteger("u32").decode(data)
        return remaining[:length], remaining[length:]


class PyObi(object):
    def __init__(self, schema):
        normalized_schema = "".join(schema.split())
        tokens = normalized_schema.split("/")
        self.schemas = [PyObiSpec.from_spec(token) for token in tokens]

    def encode(self, data: Any, index=0) -> bytes:
        return self.schemas[index].encode(data)

    def decode(self, data: bytes, index=0) -> Any:
        result, remaining = self.schemas[index].decode(data)
        if remaining:
            raise DecodeError("Not all data is consumed after decoding input")
        return result

    def encode_input(self, data: Any) -> bytes:
        return self.encode(data, index=0)

    def encode_output(self, data: Any) -> bytes:
        return self.encode(data, index=1)

    def decode_input(self, data: bytes) -> Any:
        return self.decode(data, index=0)

    def decode_output(self, data: bytes) -> Any:
        return self.decode(data, index=1)
