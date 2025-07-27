package id

type ID string

func (id ID) String() string {
	return string(id)
}

func (id ID) MarshalText() ([]byte, error) {
	return []byte(id), nil
}
