package convert

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestEpoch(t *testing.T) {

	JuneFirstInt64 := int64(1685602800)
	JuneFirstTime := time.Unix(JuneFirstInt64, 0)
	JuneFirstString := "2023-06-01T07:00:00Z"

	require.Equal(t, int64(0), EpochDate(0))
	require.Equal(t, JuneFirstInt64, EpochDate(JuneFirstInt64))
	require.Equal(t, JuneFirstInt64, EpochDate(JuneFirstTime))
	require.Equal(t, JuneFirstInt64, EpochDate(JuneFirstString))

	require.Equal(t, int64(0), EpochDate("not-a-date"))
}
