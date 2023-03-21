package grpc

import (
	"fmt"
	"github.com/weizhuxuanzx/go-tool/grpc/proto/order"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// 客户端
func client() {
	// 连接
	// client
	conn, err := grpc.Dial(Address, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalln(err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)
	// 初始化客户端
	c := order.NewOrderClient(conn)
	// 调用方法
	req := &order.OrderRequest{OriginId: "50412100454210"}
	ctx := context.Background()
	getOrder, err := c.GetOrder(ctx, req)
	if err != nil {
		return
	}
	fmt.Println(getOrder)
}
