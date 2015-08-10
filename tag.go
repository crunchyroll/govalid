// chris 081015 Struct field tag processing.

package main

import (
	"errors"
	"log"
	"reflect"
	"strconv"
	"strings"

	"math/big"
)

var errBadLength = errors.New("negative length")
var errBadBounds = errors.New("minimum length > maximum length")

type fieldMetadata struct {
	max *big.Int // Maximum value or length.
	min *big.Int // Minimum value or length.

	// Default value; if non-nil, the field is optional.  If empty,
	// the default value is the zero value of the field type.  If
	// non-empty, the value will be injected literally into the
	// generated code as the default value.
	def *string
}

func parseFieldMetadata(tag string) *fieldMetadata {
	meta := new(fieldMetadata)
	uq, err := strconv.Unquote(tag)
	if err != nil {
		// TODO Log line/char where this was encountered in .v
		// file.
		return meta
	}
	v := reflect.StructTag(uq).Get("valid")
	if v == "" {
		return meta
	}
	vs := strings.Split(v, ",")
	for _, item := range vs {
		kv := strings.Split(item, ":")
		if len(kv) == 1 {
			if kv[0] == "def" {
				zuul := ""
				meta.def = &zuul
			} else {
				// TODO Log line/char where this was
				// encountered in .v file.
				log.Printf("warning: unexpected field tag item %q\n", kv[0])
			}
		} else if len(kv) == 2 {
			if kv[0] == "min" {
				meta.min = new(big.Int)
				meta.min.SetString(kv[1], 0)
			} else if kv[0] == "max" {
				meta.max = new(big.Int)
				meta.max.SetString(kv[1], 0)
			} else if kv[0] == "def" {
				meta.def = &kv[1]
			} else {
				// TODO Log line/char where this was
				// encountered in .v file.
				log.Printf("warning: unexpected field tag item %q:%q\n", kv[0], kv[1])
			}
		} else {
			// TODO Log line/char where this was
			// encountered in .v file.
			log.Printf("warning: superfluous : in field tag item %q\n", item)
		}
	}
	return meta
}

// checkFieldMetadata makes sure that the max and min values contained
// in the given metadata are reasonable given that they'll be used for
// measuring lengths.  It returns nil if so, and an error if not.
func checkFieldMetadata(meta *fieldMetadata) error {
	// TODO Make sure these fit into ints (since len(x) is an int).
	if meta.max != nil {
		cmp := meta.max.Cmp(big.NewInt(0))
		if cmp == -1 {
			// TODO Track line/character where the offense
			// occurred.
			return errBadLength
		}
	}
	if meta.min != nil {
		cmp := meta.max.Cmp(big.NewInt(0))
		if cmp == -1 {
			// TODO Track line/character where the offense
			// occurred.
			return errBadLength
		}
	}
	if meta.max != nil && meta.min != nil {
		cmp := meta.max.Cmp(meta.min)
		if cmp == -1 {
			// TODO Track line/character where the offense
			// occurred.
			return errBadBounds
		}
	}
	return nil
}
