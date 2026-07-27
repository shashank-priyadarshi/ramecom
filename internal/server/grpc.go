package server

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/rajabhishekmaurya/ecom/libs/config"
	pb "github.com/rajabhishekmaurya/ecom/libs/proto/payment"
)

// PaymentServer is a gRPC process that serves PaymentService.
type PaymentServer struct {
	cfg     *config.Config
	handler pb.PaymentServiceServer
}

func NewPayment(cfg *config.Config, handler pb.PaymentServiceServer) *PaymentServer {
	return &PaymentServer{
		cfg:     cfg,
		handler: handler,
	}
}

func (s *PaymentServer) Start() error {
	lis, err := net.Listen("tcp", ":"+s.cfg.Srv.Port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	pb.RegisterPaymentServiceServer(grpcServer, s.handler)

	log.Printf("%s started on :%s", s.cfg.Srv.Name, s.cfg.Srv.Port)
	return grpcServer.Serve(lis)
}
