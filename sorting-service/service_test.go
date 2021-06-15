package main

import (
	"github.com/sandio/sort/gen"
	"testing"
)

func TestLoadItems(t *testing.T) {
	s := newSortingService()
	tests := map[string]struct {
		input *gen.LoadItemsRequest
		want  *gen.LoadItemsResponse
	}{
		"empty": {input: &gen.LoadItemsRequest{}, want: &gen.LoadItemsResponse{}},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			_, err := s.LoadItems(nil, tc.input)
			if err != nil {
				t.Errorf("got %v, want %v", err, nil)
			}
		})
	}
}
