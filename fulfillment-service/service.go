package main

import (
	"context"
	"log"
	"time"

	"github.com/preslavmihaylov/ordertocubby"
	"github.com/sandio/sort/gen"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const cubbyCount = 10

func newFulfillmentService(client gen.SortingRobotClient) gen.FulfillmentServer {
	f := &fulfillmentService{
		fulfillmentStatuses: make([]*gen.FulfillmentStatus, 0),
		ordersQueue:         make(chan []*gen.Order),
		client:              client}
	go f.controlRobotArm()
	return f
}

type fulfillmentService struct {
	fulfillmentStatuses []*gen.FulfillmentStatus
	ordersQueue         chan []*gen.Order
	client              gen.SortingRobotClient
}

func (f *fulfillmentService) controlRobotArm() {
	for {
		orders := <-f.ordersQueue
		log.Printf("Received on chan %v", orders)
		orderToCubby := make(map[string]string)
		orderToMovedItems := make(map[string][]string)
		var i uint32 = 1
		// map order to cubby
		for _, order := range orders {
			for {
				cubbyId := ordertocubby.Map(order.Id, i, cubbyCount)
				if cubbyId == orderToCubby[order.Id] {
					i++
					continue
				}
				orderToCubby[order.Id] = cubbyId
				fulfillmentStatus := &gen.FulfillmentStatus{
					Order: order,
					Cubby: &gen.Cubby{Id: cubbyId}} // state is pending by default
				f.fulfillmentStatuses = append(f.fulfillmentStatuses, fulfillmentStatus)
				break
			}
		}
	FulfillmentLoop:
		for {
			time.Sleep(100 * time.Millisecond)
			selected, err := f.client.SelectItem(context.Background(), &gen.SelectItemRequest{})
			if err != nil {
				// status code error handling
				if st, ok := status.FromError(err); ok {
					switch st.Code() {
					case codes.AlreadyExists:
						break FulfillmentLoop
					case codes.NotFound:
						break FulfillmentLoop
					}
				}
			}
			for _, order := range orders {
				// move item
				for _, item := range order.Items {
					if selected.Item.Code == item.Code {
						_, err := f.client.MoveItem(context.Background(), &gen.MoveItemRequest{Cubby: &gen.Cubby{Id: orderToCubby[order.Id]}})
						if err != nil {
							log.Printf("%v", err)
							return
						}
						orderToMovedItems[order.Id] = append(orderToMovedItems[order.Id], item.Code) // track what is moved
					}
				}
				// ready an order
				if len(order.Items) == len(orderToMovedItems[order.Id]) {
					for _, fulfillmentStatus := range f.fulfillmentStatuses {
						if fulfillmentStatus.Order.Id == order.Id {
							fulfillmentStatus.State = gen.OrderState_READY
							delete(orderToMovedItems, order.Id) // ready it once
							log.Printf("order %v, status %v", fulfillmentStatus.Order.Id, fulfillmentStatus.State)
						}
					}
				}
			}
		}
	}
}

func (f *fulfillmentService) LoadOrders(ctx context.Context, in *gen.LoadOrdersRequest) (out *gen.LoadOrdersResponse, err error) {
	if in.Orders == nil {
		return nil, status.Errorf(codes.NotFound, "orders not given")
	}
	go func() {
		f.ordersQueue <- in.Orders
	}()
	var preparedOrders []*gen.PreparedOrder

	return &gen.LoadOrdersResponse{PreparedOrders: preparedOrders}, nil
}

func (f *fulfillmentService) GetOrderStatusById(ctx context.Context, in *gen.OrderIdRequest) (*gen.OrdersStatusResponse, error) {
	for _, fulfillmentStatus := range f.fulfillmentStatuses {
		if fulfillmentStatus.Order.Id == in.OrderId {
			return &gen.OrdersStatusResponse{Status: []*gen.FulfillmentStatus{fulfillmentStatus}}, nil
		}
	}
	return nil, status.Errorf(codes.NotFound, "order not found")
}
func (f *fulfillmentService) GetAllOrdersStatus(context.Context, *gen.Empty) (*gen.OrdersStatusResponse, error) {
	return &gen.OrdersStatusResponse{Status: f.fulfillmentStatuses}, nil
}
func (f *fulfillmentService) MarkFullfilled(context.Context, *gen.OrderIdRequest) (*gen.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MarkFullfilled not implemented")
}
