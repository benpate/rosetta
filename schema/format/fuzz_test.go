package format

import (
	"net/url"
	"strings"
	"testing"
)

// FuzzURI feeds arbitrary strings to the URI validator to confirm that it never panics,
// and that its accept/reject decision matches its contract: an accepted value is returned
// unchanged and (unless empty) parses via url.ParseRequestURI with a scheme and fits the
// length limit; a rejected value returns the empty string.
func FuzzURI(f *testing.F) {

	f.Add("")
	f.Add("https://example.com")
	f.Add("mailto:a@b.com")
	f.Add("https://example.com/path?q=1")
	f.Add("example.com")
	f.Add("/path/only")
	f.Add("//example.com/path")
	f.Add("not a uri")
	f.Add("https://example.com/%zz")
	f.Add("https://example.com\x00")
	f.Add("https://example.com/" + strings.Repeat("a", 2048))

	validate := URI("")

	f.Fuzz(func(t *testing.T, value string) {

		result, err := validate(value)

		// A rejected value must return the empty string alongside the error.
		if err != nil {
			if result != "" {
				t.Fatalf("URI rejected %q but returned non-empty result %q", value, result)
			}
			return
		}

		// An accepted value must pass through unchanged.
		if result != value {
			t.Fatalf("URI accepted %q but returned altered result %q", value, result)
		}

		// The empty string is accepted by design, with no further requirements.
		if value == "" {
			return
		}

		// Every other accepted value must fit the length limit and re-derive as an
		// absolute URI: parseable by ParseRequestURI and carrying a scheme.
		if len(value) > 2048 {
			t.Fatalf("URI accepted %q, which exceeds the 2048-byte limit", value)
		}

		if parsed, parseErr := url.ParseRequestURI(value); parseErr != nil || parsed.Scheme == "" {
			t.Fatalf("URI accepted %q, which is not an absolute URI (err=%v)", value, parseErr)
		}
	})
}

// FuzzURL feeds arbitrary strings to the URL validator to confirm that it never panics,
// and that its accept/reject decision matches its contract: an accepted value is returned
// unchanged and (unless empty) parses via url.Parse with both a scheme and a host and fits
// the length limit; a rejected value returns the empty string.
func FuzzURL(f *testing.F) {

	f.Add("")
	f.Add("https://example.com")
	f.Add("http://example.com/path?q=1#fragment")
	f.Add("https://user:pass@example.com:8080/path")
	f.Add("https://[::1]:8080/")
	f.Add("https://example.com/a b")
	f.Add("mailto:user@example.com")
	f.Add("javascript:alert(1)")
	f.Add("file:///etc/passwd")
	f.Add("//example.com/path")
	f.Add("http://exa mple.com")
	f.Add("https://example.com/%zz")
	f.Add("https://example.com\x00")
	f.Add("https://example.com/" + strings.Repeat("a", 2048))

	validate := URL("")

	f.Fuzz(func(t *testing.T, value string) {

		result, err := validate(value)

		// A rejected value must return the empty string alongside the error.
		if err != nil {
			if result != "" {
				t.Fatalf("URL rejected %q but returned non-empty result %q", value, result)
			}
			return
		}

		// An accepted value must pass through unchanged.
		if result != value {
			t.Fatalf("URL accepted %q but returned altered result %q", value, result)
		}

		// The empty string is accepted by design, with no further requirements.
		if value == "" {
			return
		}

		// Every other accepted value must fit the length limit and re-derive as an
		// absolute URL: parseable by url.Parse with both a scheme and a host.
		if len(value) > 2048 {
			t.Fatalf("URL accepted %q, which exceeds the 2048-byte limit", value)
		}

		parsed, parseErr := url.Parse(value)

		if parseErr != nil {
			t.Fatalf("URL accepted %q, which does not parse: %v", value, parseErr)
		}

		if !parsed.IsAbs() {
			t.Fatalf("URL accepted %q, which has no scheme", value)
		}

		if parsed.Host == "" {
			t.Fatalf("URL accepted %q, which has no host", value)
		}
	})
}
