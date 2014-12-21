package asciify

import "testing"

var asciiMapping = map[string]string{
	"càché": "cache",
	"\t":    " ",
}

func TestClobber(t *testing.T) {
	for input, output := range asciiMapping {
		result := Clobber(input)
		if result != output {
			t.Errorf("%#v -> %#v != %#v", input, result, output)
		}
	}
}
