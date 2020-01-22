package types

type RawJson []byte

var EmptyMap = []byte("{}")

func (j RawJson) MarshalJSON() ([]byte, error) {
	return []byte(j), nil
}

func (j *RawJson) UnmarshalJSON(b []byte) error {
	*j = RawJson(b)
	return nil
}
