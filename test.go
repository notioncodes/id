package id

var TestIDs = map[string]string{
	"valid-dashed":   "550e8400-e29b-41d4-a716-446655440000",
	"valid-undashed": "550e8400e29b41d4a716446655440000",
	"invalid-char":   "xyzzy123e29b41d4a716446655440000",
	"invalid-id":     "invalid-id",
}
