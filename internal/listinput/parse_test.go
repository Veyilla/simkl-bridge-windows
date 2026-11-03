package listinput

import "testing"

func TestParseListLinkOrID(t *testing.T) {
	for _, v := range []string{"123", " https://simkl.com/50/list/123/test?sort=position ", "https://simkl.com/lists/123/test"} {
		id, err := Parse(v)
		if err != nil || id != 123 {
			t.Fatalf("%q: %d %v", v, id, err)
		}
	}
}
func TestRejectUnsafeOrMalformedInput(t *testing.T) {
	for _, v := range []string{"0", "-1", "1.5", "https://other.test/50/list/123/test", "https://simkl.com.evil.test/50/list/123/test", "https://name@simkl.com/50/list/123/test", "http://simkl.com/50/list/123/test", "https://simkl.com:444/50/list/123/test", "https://simkl.com/50/list/0/test", "https://simkl.com/50/list/123/test/extra", "9999999999999999999999999"} {
		if _, err := Parse(v); err == nil {
			t.Fatalf("Ungültiger Link akzeptiert: %q", v)
		}
	}
}
