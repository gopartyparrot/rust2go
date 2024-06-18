package rust2go

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadFileCache(t *testing.T) {
	var fc *FileCache[string, string]

	lineParser := func(line string) (string, string, error) {
		arr := strings.Split(strings.TrimSpace(line), ":")
		return arr[0], arr[1], nil

	}
	getter := func(k string) (string, error) {
		t.Log("getter called")
		return "b", nil
	}
	lineFmt := func(k string, v string) string {
		return k + ":" + v
	}
	load := func() {
		c, err := LoadFileCache("./test.cache", lineParser)
		require.NoError(t, err)
		fc = c
	}

	load()
	cached, err := fc.Load("a", getter, lineFmt)
	require.NoError(t, err)
	require.Equal(t, "b", cached)
	_ = fc.Close()

	load()
	cached, err = fc.Load("a", getter, lineFmt)
	require.NoError(t, err)
	require.Equal(t, "b", cached)

}
