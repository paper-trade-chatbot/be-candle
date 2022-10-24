package service

import (
	"context"
	"fmt"

	"github.com/paper-trade-chatbot/be-candle/logging"

	"github.com/paper-trade-chatbot/be-candle/config"
	"github.com/paper-trade-chatbot/be-candle/service/product"
	productGrpc "github.com/paper-trade-chatbot/be-proto/product"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var Impl ServiceImpl
var (
	ProductServiceHost    = config.GetString("PRODUCT_GRPC_HOST")
	ProductServerGRpcPort = config.GetString("PRODUCT_GRPC_PORT")
	productServiceConn    *grpc.ClientConn
)

type ServiceImpl struct {
	ProductIntf product.ProductIntf
}

func GrpcDial(addr string) (*grpc.ClientConn, error) {
	return grpc.Dial(addr, grpc.WithInsecure(), grpc.WithDefaultCallOptions(
		grpc.MaxCallRecvMsgSize(20*1024*1024),
		grpc.MaxCallSendMsgSize(20*1024*1024)), grpc.WithUnaryInterceptor(clientInterceptor))
}

func Initialize(ctx context.Context) {

	var err error

	addr := ProductServiceHost + ":" + ProductServerGRpcPort
	fmt.Println("dial to order grpc server...", addr)
	productServiceConn, err = GrpcDial(addr)
	if err != nil {
		fmt.Println("Can not connect to gRPC server:", err)
	}
	fmt.Println("dial done")
	productConn := productGrpc.NewProductServiceClient(productServiceConn)
	Impl.ProductIntf = product.New(productConn)
}

func Finalize(ctx context.Context) {
	productServiceConn.Close()
}

func clientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	requestId, _ := ctx.Value(logging.ContextKeyRequestId).(string)
	account, _ := ctx.Value(logging.ContextKeyAccount).(string)

	ctx = metadata.NewOutgoingContext(ctx, metadata.New(map[string]string{
		logging.ContextKeyRequestId: requestId,
		logging.ContextKeyAccount:   account,
	}))

	err := invoker(ctx, method, req, reply, cc, opts...)
	if err != nil {
		fmt.Println("clientInterceptor err:", err.Error())
	}

	return err
}