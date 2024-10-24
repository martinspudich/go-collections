package gocollections

import (
	"reflect"
	"testing"
)

func TestCopyMap(t *testing.T) {
	tests := []struct {
		name string
		src  map[string]any
		want map[string]any
	}{
		{
			name: "simple copy - value int",
			src:  map[string]any{"a": 1, "b": 2},
			want: map[string]any{"a": 1, "b": 2},
		},
		{
			name: "simple copy - value string",
			src:  map[string]any{"a": "1", "b": "2"},
			want: map[string]any{"a": "1", "b": "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := map[string]any{}
			CopyMap(got, tt.src)
			result := reflect.DeepEqual(got, tt.want)
			if result != true {
				t.Errorf("got: %v, want: %v", got, tt.want)
			}
		})
	}
}
