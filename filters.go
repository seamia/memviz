package memviz

import (
	"fmt"
	"os"
	"strings"
)

type ignoreResponse int

const (
	doNotSkip        ignoreResponse = 0
	ignoreCompletely ignoreResponse = 1
	ignoreValue      ignoreResponse = 2
)

func skipField(kind, collection, field string) ignoreResponse {

	if len(Options().Discard) != 0 {

		collection = strings.Trim(collection, " \"\\")

		full := kind + ":" + collection + "." + field
		if value, found := Options().Discard[full]; found {
			switch value {
			case 0:
				return doNotSkip
			case 1:
				return ignoreCompletely
			case 2:
				return ignoreValue
			default:
				warning("unrecognized value (%v) for key (%v)", value, full)
			}
		}
	}
	return doNotSkip
}

func interpretValueType(val, typ string) string {

	if len(Options().Substitute) != 0 {
		if replace, found := Options().Substitute[typ]; found && len(replace) != 0 {
			if change, exists := replace[val]; exists {
				return change
			}
		}
	}

	if typ == "string" {
		return val
	}

	return val + " (" + typ + ")"
}

func warning(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Warning: "+format+"\n", args...)
}
