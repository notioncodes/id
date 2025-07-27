package id

import (
	"errors"
)

// hexTable provides O(1) hex character validation as using a lookup table
// is faster than range checks.
var hexTable = [256]bool{
	'0': true, '1': true, '2': true, '3': true, '4': true,
	'5': true, '6': true, '7': true, '8': true, '9': true,
	'a': true, 'b': true, 'c': true, 'd': true, 'e': true, 'f': true,
	'A': true, 'B': true, 'C': true, 'D': true, 'E': true, 'F': true,
}

func isHex(c byte) bool {
	return hexTable[c]
}

type IDParser struct {
	cache Cache
}

func NewIDParser(c Cache) *IDParser {
	return &IDParser{cache: c}
}

func (p *IDParser) Parse(input string) (string, error) {
	if val, ok := p.cache.Get(input); ok {
		return val, nil
	}

	// Parse and format the ID
	result, err := parse(input)
	if err != nil {
		// Cache the error to avoid re-parsing invalid IDs
		// uuidCache.StoreError(id)
		return "", err
	}

	p.cache.Set(input, result)
	return result, nil
}

// Pre-allocate common error values to avoid allocations.
var (
	errInvalidLength = errors.New("invalid ID length")
	errInvalidChar   = errors.New("invalid ID: non-hex character")
	errMisplacedDash = errors.New("invalid ID: misplaced dash")
	errInvalidStruct = errors.New("invalid ID structure")
)

// parse parses a Notion ID and returns a normalized formatted ID as a UUID string.
//
// Arguments:
// - id: The ID to parse
//
// Returns:
// - The normalized formatted ID as a UUID string.
// - An error if the ID is invalid.
func parse(id string) (string, error) {
	// Early length validation - fastest possible bailout.
	if len(id) != 32 && len(id) != 36 {
		return "", errInvalidLength
	}

	var buf [36]byte
	var hexCount, dashCount int
	var i, j int

	switch len(id) {
	case 32:
		for i = 0; i < 32; i++ {
			c := id[i]
			if !isHex(c) {
				return "", errInvalidChar
			}
			buf[j] = c
			j++
			hexCount++
			switch hexCount {
			case 8, 12, 16, 20:
				buf[j] = '-'
				j++
				dashCount++
			}
		}
	case 36:
		for i = 0; i < 36; i++ {
			c := id[i]
			switch i {
			case 8, 13, 18, 23:
				if c != '-' {
					return "", errMisplacedDash
				}
				buf[j] = c
				j++
				dashCount++
			default:
				if !isHex(c) {
					return "", errInvalidChar
				}
				buf[j] = c
				j++
				hexCount++
			}
		}
	}

	if hexCount != 32 || dashCount != 4 {
		return "", errInvalidStruct
	}

	// Avoid string allocation by using unsafe conversion when possible.
	// For now, we'll keep it safe but consider this could be optimized further.
	return string(buf[:]), nil
}
