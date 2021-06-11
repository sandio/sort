package main

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/sandio/sort/gen"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func newSortingService() gen.SortingRobotServer {
	return &sortingService{cubbies: make(map[string][]*gen.Item)}
}

type sortingService struct {
	items     []*gen.Item
	selection *gen.Item
	cubbies   map[string][]*gen.Item
}

func (s *sortingService) LoadItems(ctx context.Context, in *gen.LoadItemsRequest) (*gen.LoadItemsResponse, error) {
	log.Printf("Received: %v", in.Items)
	s.items = append(s.items, in.Items...)
	log.Printf("Stored: %v", s.items)

	return &gen.LoadItemsResponse{}, nil
}

func (s *sortingService) SelectItem(context.Context, *gen.SelectItemRequest) (*gen.SelectItemResponse, error) {
	if s.selection != nil {
		return nil, status.Errorf(codes.AlreadyExists, "item already selected")
	}
	if len(s.items) == 0 {
		return nil, status.Errorf(codes.NotFound, "items is empty")
	}
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)
	index := random.Intn(len(s.items))
	s.selection = s.items[index]
	log.Printf("Selected: %v, %v", index, s.selection)

	s.items = append(s.items[:index], s.items[index+1:]...)
	log.Printf("Left: %v", s.items)

	return &gen.SelectItemResponse{Item: s.selection}, nil
}

func (s *sortingService) MoveItem(ctx context.Context, in *gen.MoveItemRequest) (*gen.MoveItemResponse, error) {
	if in.Cubby == nil {
		return nil, status.Errorf(codes.NotFound, "cubby not given")
	}
	if s.selection == nil {
		return nil, status.Errorf(codes.NotFound, "item is not selected")
	}
	s.cubbies[in.Cubby.Id] = append(s.cubbies[in.Cubby.Id], s.selection)
	log.Printf("Cubbies: %v", s.cubbies)
	s.selection = nil
	return &gen.MoveItemResponse{}, nil
}
