package main

import (
	"github.com/sandio/sort/gen"
	"testing"
)

func TestLoadItems(t *testing.T) {
	s := newSortingService()
	tests := map[string]struct {
		input *gen.LoadItemsRequest
		want  nil
	}{
		"empty": {input: &gen.LoadItemsRequest{}, want: nil},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := s.LoadItems(nil, tc.input)
			if err != want {
				t.Errorf("got %v, want %v", err, want)
			}
		})
	}
}
