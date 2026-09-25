package schema

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// TestTypeName pins the detail that replaces a value
func TestTypeName(t *testing.T) {
	require.Equal(t, "string", typeName("s3cr3t-value"))
	secretMap := map[string]any{"secret": "s3cr3t-value"}
	require.Equal(t, "*map[string]interface {}", typeName(&secretMap))
	require.Equal(t, "<nil>", typeName(nil))
}
