package indent

import (
	"fmt"
	"testing"
)

func TestIndent(t *testing.T) {
	it := NewIndenter(nil)
	output, err := it.Process("a+b*c", "")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(output)
}
