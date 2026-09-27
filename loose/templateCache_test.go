package loose

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

// cachedTemplate returns the template cached for a source string, and TRUE if there is one
func cachedTemplate(source string) (Template, bool) {

	cached, ok := templateCache.Load(source)

	if !ok {
		return Template{}, false
	}

	compiled, ok := cached.(Template)
	return compiled, ok
}

// TestCachedTemplate_CachesTemplates requires that a template is compiled into the cache
// once, and served from it after that
func TestCachedTemplate_CachesTemplates(t *testing.T) {

	const source = "/cache/{{.ID}}"

	first := CachedTemplate(source)
	require.True(t, first.IsTemplate())

	cached, ok := cachedTemplate(source)
	require.True(t, ok)
	require.Equal(t, first, cached)

	// A second compile is the cached template itself, and does not grow the cache
	size := templateCacheSize.Load()
	require.Equal(t, first, CachedTemplate(source))
	require.Equal(t, size, templateCacheSize.Load())
}

// TestCachedTemplate_SkipsPlainText requires that plain text and malformed templates are
// returned as plain text, and never cached
func TestCachedTemplate_SkipsPlainText(t *testing.T) {

	for _, source := range []string{"", "off", "/cache/plain", "/cache/{{.ID | nosuchfunc}}"} {

		compiled := CachedTemplate(source)
		require.False(t, compiled.IsTemplate(), source)
		require.Equal(t, source, compiled.String())

		_, ok := cachedTemplate(source)
		require.False(t, ok, source)
	}
}

// TestCachedTemplate_Limit requires that a full cache stores nothing more, while still
// returning a compiled template
func TestCachedTemplate_Limit(t *testing.T) {

	const source = "/cache/limit/{{.ID}}"

	size := templateCacheSize.Load()
	templateCacheSize.Store(templateCacheLimit)
	t.Cleanup(func() { templateCacheSize.Store(size) })

	compiled := CachedTemplate(source)
	require.True(t, compiled.IsTemplate())
	require.Equal(t, int64(templateCacheLimit), templateCacheSize.Load())

	_, ok := cachedTemplate(source)
	require.False(t, ok)
}

// TestStoreTemplate_KeepsFirst requires that storing a source that is already cached returns
// the cached template, and leaves the cache size unchanged: the path two goroutines take when they
// compile one source at once
func TestStoreTemplate_KeepsFirst(t *testing.T) {

	const source = "/cache/first/{{.ID}}"

	first := CachedTemplate(source)
	size := templateCacheSize.Load()

	second := NewTemplate(source)
	require.Equal(t, first, storeTemplate(source, second))
	require.Equal(t, size, templateCacheSize.Load())
}

// TestCachedTemplate_Concurrent compiles and renders one source from many goroutines at once,
// which the race detector reports if the cache is not safe to share
func TestCachedTemplate_Concurrent(t *testing.T) {

	const source = "/cache/concurrent/{{.ID}}"

	results := make(chan bool, 50)

	var wait sync.WaitGroup
	for index := range 50 {
		id := strconv.Itoa(index)

		wait.Go(func() {
			result, err := CachedTemplate(source).Execute(templateData{ID: id})
			results <- (err == nil) && (result == "/cache/concurrent/"+id)
		})
	}

	wait.Wait()
	close(results)

	for ok := range results {
		require.True(t, ok)
	}

	_, ok := cachedTemplate(source)
	require.True(t, ok)
}
