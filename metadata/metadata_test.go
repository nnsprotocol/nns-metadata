package metadata

import "testing"

func TestFormatHex(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{"0xABCDF0123", "0xabcdf0123"},
		{"71445694856899356118833591695777802838200464758455892088604082371692842540648", "0x9df4d48c0891fcf62c436d7d609abf13b418d638d808bf4de94b31f88e5d5e68"},
	}
	for _, tC := range testCases {
		t.Run(tC.in, func(t *testing.T) {
			out, ok := formatHashToHex(tC.in)
			if !ok {
				t.Fatal("formatting failed")
			}
			if out != tC.out {
				t.Fatalf("unexpected value: '%v' vs '%v'", out, tC.out)
			}
		})
	}
}
