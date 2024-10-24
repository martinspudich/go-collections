package gocollections

// CopyMap function will creates hard copy of the source map to destination map.
func CopyMap[K comparable, V any](dst, src map[K]V) {
	for k, v := range src {
		dst[k] = v
	}
}
