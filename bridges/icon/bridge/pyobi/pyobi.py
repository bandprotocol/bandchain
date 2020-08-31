class PyObiSpec(object):
    impls = []

    def __init_subclass__(cls, **kwargs):
        super().__init_subclass__(**kwargs)
        cls.impls.append(cls)

    @classmethod
    def from_spec(cls, spec):
        for impl in cls.impls:
            if impl.match_schema(spec):
                return impl(spec)
        raise ValueError("Cannot parse spec: {}".format(spec))

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
    def match_schema(cls, schema):
        return schema[:1] in ["i", "u"] and schema[1:] in ["8", "16", "32", "64", "128", "256"]

    def encode(self, value):
        return value.to_bytes(self.size_in_bytes, byteorder="big", signed=self.is_signed)

    def decode(self, data):
        return (
            int.from_bytes(data[: self.size_in_bytes], byteorder="big", signed=self.is_signed),
            data[self.size_in_bytes :],
        )


class PyObiBool(PyObiSpec):
    def __init__(self, spec=""):
        pass

    @classmethod
    def match_schema(cls, schema):
        return schema == "bool"

    def encode(self, value):
        return PyObiInteger("u8").encode(1 if value else 0)

    def decode(self, data):
        u8, remaining = PyObiInteger("u8").decode(data)
        if u8 == 1:
            return True, remaining
        elif u8 == 0:
            return False, remaining
        raise ValueError("Boolean value must be 1 or 0 but got {}".format(u8))


class PyObiVector(PyObiSpec):
    def __init__(self, spec):
        self.intl_obi = self.from_spec(spec[1:-1])

    @classmethod
    def match_schema(cls, schema):
        return schema[0] == "[" and schema[-1] == "]"

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
    def __init__(self, spec):
        self.intl_obi_kvs = []
        fields = ['']
        curly_count = 0
        for c in spec[1:-1]:
            if c == ',' and curly_count == 0:
                fields.append('')
            else:
                fields[-1] = fields[-1] + c
                if c == '{':
                    curly_count += 1
                if c == '}':
                    curly_count -= 1
        for each in fields:
            tokens = each.split(":", 1)
            if len(tokens) != 2:
                raise ValueError("Expect at least one colon for each struct field")
            self.intl_obi_kvs.append((tokens[0], self.from_spec(tokens[1])))

    @classmethod
    def match_schema(cls, schema):
        return schema[0] == "{" and schema[-1] == "}"

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
    def __init__(self, spec=""):
        pass

    @classmethod
    def match_schema(cls, schema):
        return schema == "string"

    def encode(self, value):
        return PyObiInteger("u32").encode(len(value)) + value.encode()

    def decode(self, data):
        length, remaining = PyObiInteger("u32").decode(data)
        return remaining[:length].decode(), remaining[length:]


class PyObiBytes(PyObiSpec):
    def __init__(self, spec=""):
        pass

    @classmethod
    def match_schema(cls, schema):
        return schema == "bytes"

    def encode(self, value):
        return PyObiInteger("u32").encode(len(value)) + value

    def decode(self, data):
        length, remaining = PyObiInteger("u32").decode(data)
        return remaining[:length], remaining[length:]


class PyObi(object):
    def __init__(self, schema):
        normalized_schema = "".join(schema.split())
        tokens = normalized_schema.split("/")
        self.schemas = [PyObiSpec.from_spec(token) for token in tokens]

    def encode(self, data, index=0):
        return self.schemas[index].encode(data)

    def decode(self, data, index=0):
        result, remaining = self.schemas[index].decode(data)
        if remaining:
            raise ValueError("Not all data is consumed after decoding input")
        return result

    def encode_input(self, data):
        return self.encode(data, index=0)

    def encode_output(self, data):
        return self.encode(data, index=1)

    def decode_input(self, data):
        return self.decode(data, index=0)

    def decode_output(self, data):
        return self.decode(data, index=1)