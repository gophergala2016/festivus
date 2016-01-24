package holidays

import "testing"

func TestParseFile(t *testing.T) {
	path := "testdata/hr.txt"
	d, err := parseFile(path)
	if len(d) != 29 {
		t.Errorf("len(parseFile(valid)) = %v, want 29", len(d))
	}
	if err != nil {
		t.Errorf("parseFile(valid) = _, %v; want _, nil", err)
	}
}
