package rust2go

import (
	"fmt"
	"os"
	"strings"
	"sync"
)

// Filecache is a cache for known input output
// cache write to file in line mode, each line is a input output pair(K, V)
// cache file is append only, each time new cache appended to the end of the file
type FileCache[K comparable, V any] struct {
	file  *os.File
	cache map[K]V
	mu    sync.RWMutex
}

func LoadFileCache[K comparable, V any](path string, lineParser func(string) (K, V, error)) (*FileCache[K, V], error) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	lines := strings.Split(string(content), "\n")

	cache := make(map[K]V, len(lines))
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		k, v, err := lineParser(line)
		if err != nil {
			return nil, err
		}
		cache[k] = v
	}

	return &FileCache[K, V]{
		file:  f,
		cache: cache,
	}, nil

}

// Load load value from cache, if not found, call getter to get value and save to cache
func (fc *FileCache[K, V]) Load(k K, getter func(K) (V, error), lineFmt func(k K, v V) string) (V, error) {
	fc.mu.RLock()

	if v, ok := fc.cache[k]; ok {
		defer fc.mu.RUnlock()
		return v, nil
	}

	fc.mu.RUnlock()

	v, err := getter(k)

	fc.mu.Lock()
	defer fc.mu.Unlock()
	if err != nil {
		return v, err
	}

	fc.cache[k] = v
	_, err = fc.file.WriteString(strings.TrimSpace(lineFmt(k, v)) + "\n") // write line
	if err != nil {
		fmt.Println("[FileCache] write file cache failed", err) // TODO now just ignore the error, since it's not critical
	}

	return v, nil
}

func (fc *FileCache[K, V]) Close() error {
	return fc.file.Close()
}
