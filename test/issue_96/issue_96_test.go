package pegmatch_test

import (
	"bytes"
	"testing"

	pegmatch "github.com/mna/pigeon/test/issue_96"
)

func TestMatchSimple(t *testing.T) {
	m := [1]string{"{test}"}
	pegmatch.ContentString = "my {test}"
	_, err := pegmatch.ParseReader("", bytes.NewBufferString(m[0]))
	if err != nil {
		t.Error("failed")
	}
}
