import re


class PyObiSpec(object):
    impls = []

    def __init_subclass__(cls, **kwargs):
        super().__init_subclass__(**kwargs)
        cls.impls.append(cls)

    @classmethod
    def from_spec(cls, spec):
        for impl in cls.impls:
            if re.match(impl.REGEX, spec):
                return impl(spec)
        raise ValueError("Cannot parse spec: {}".format(spec))

    def __init__(self, spec):
        raise NotImplementedError()

    def encode(self, value):
        raise NotImplementedError()

    def decode(self, data):
        raise NotImplementedError()


class PyObiInteger(PyObiSpec):
    REGEX = re.compile(r"^(u|i)(8|16|32|64|128|256)$")

    def __init__(self, spec):
        self.is_signed = spec[0] == "i"
        self.size_in_bytes = int(spec[1:]) // 8

    def encode(self, value):
        return value.to_bytes(self.size_in_bytes, byteorder="big", signed=self.is_signed)

    def decode(self, data):
        return (
            int.from_bytes(data[: self.size_in_bytes], byteorder="big", signed=self.is_signed,),
            data[self.size_in_bytes :],
        )


class PyObiVector(PyObiSpec):
    REGEX = re.compile(r"^\[.*\]$")

    def __init__(self, spec):
        self.intl_obi = self.from_spec(spec[1:-1])

    def encode(self, value):
        result = PyObiInteger("u32").encode(len(value))
        for each in value:
            result = result + self.intl_obi.encode(each)
        return result

    def decode(self, data):
        length, remaining = PyObiInteger("u32").decode(data)
        result = []
        for _ in range(length):
            each, remaining = self.intl_obi.decode(remaining)
            result.append(each)
        return result, remaining


class PyObiStruct(PyObiSpec):
    REGEX = re.compile(r"^{.*}$")

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

    def encode(self, value):
        result = b""
        for key, spec in self.intl_obi_kvs:
            result = result + spec.encode(value[key])
        return result

    def decode(self, data):
        result = {}
        for key, spec in self.intl_obi_kvs:
            result[key], data = spec.decode(data)
        return result, data


class PyObiString(PyObiSpec):
    REGEX = re.compile(r"^string$")

    def __init__(self, spec):
        pass

    def encode(self, value):
        return PyObiInteger("u32").encode(len(value)) + value.encode()

    def decode(self, data):
        length, remaining = PyObiInteger("u32").decode(data)
        return remaining[:length].decode(), remaining[length:]


class PyObiBytes(PyObiSpec):
    REGEX = re.compile(r"^bytes$")

    def __init__(self, spec):
        pass

    def encode(self, value):
        return PyObiInteger("u32").encode(len(value)) + value

    def decode(self, data):
        length, remaining = PyObiInteger("u32").decode(data)
        return remaining[:length], remaining[length:]


class PyObi(object):
    def __init__(self, schema):
        normalized_schema = re.sub(r"\s+", "", schema)
        tokens = normalized_schema.split("/")
        if len(tokens) != 2:
            raise ValueError("Expect one forward slash in OBI schema")
        self.input_schema = PyObiSpec.from_spec(tokens[0])
        self.output_schema = PyObiSpec.from_spec(tokens[1])

    def encode_input(self, value):
        return self.input_schema.encode(value)

    def decode_input(self, data):
        result, remaining = self.input_schema.decode(data)
        if remaining:
            raise ValueError("Not all data is consumed after decoding input")
        return result

    def encode_output(self, value):
        return self.output_schema.encode(value)

    def decode_output(self, data):
        result, remaining = self.output_schema.decode(data)
        if remaining:
            raise ValueError("Not all data is consumed after decoding output")
        return result

