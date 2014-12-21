package asciify

import (
	"golang.org/x/text/unicode/norm"
	"unicode"
)

// Clobber attempts to convert a string into plain ASCII by dropping accents
// and unrepresentable characters.
func Clobber(s string) string {
	// XXX This is a bit over-engineered... perhaps we should just drop any
	// message that contains non-ASCII?
	var out []byte
	for _, r := range norm.NFKD.String(s) {
		if unicode.IsSpace(r) {
			out = append(out, 32)
		} else if r < unicode.MaxASCII {
			out = append(out, byte(r))
		}
	}
	return string(out)
}
