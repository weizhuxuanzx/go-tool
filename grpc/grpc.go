package grpc

import (
	"fmt"
	"github.com/weizhuxuanzx/go-tool/grpc/proto/order"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"net"
)

const (
	Address = "127.0.0.1:50052"
)

// 定义orderService并实现约定的接口
type orderService struct {
	order.UnimplementedOrderServer
	service order.OrderServer
}

// GetOrder 实现Hello服务接口
func (h orderService) GetOrder(ctx context.Context, in *order.OrderRequest) (*order.OrderResponse, error) {
	resp := new(order.OrderResponse)
	resp.OriginId = fmt.Sprintf("Hello 恭喜你调通GRPC客户端 %s.", in.OriginId)
	return resp, nil
}

func (h orderService) mustEmbedUnimplementedOrderServer() {

}

// 服务端
func server() {
	listen, err := net.Listen("tcp", Address)
	if err != nil {
		grpclog.Fatalf("Failed to listen: %v", err)
	}
	// 实例化grpc Server
	s := grpc.NewServer()
	var service orderService
	// 注册OrderService
	order.RegisterOrderServer(s, service)
	err = s.Serve(listen)
	if err != nil {
		panic(err)
	}
}
