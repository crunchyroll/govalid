// chris 081015 Struct field tag processing.

package main

import (
	"log"
	"reflect"
	"strconv"
	"strings"
)

var maxInt = int((^uint(0)) >> 1)

type fieldMetadata struct {
	// Maximum and minimum value or length.  If the empty string,
	// there is no max/min.
	max string
	min string

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
			if kv[0] == "max" {
				meta.max = kv[1]
			} else if kv[0] == "min" {
				meta.min = kv[1]
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
