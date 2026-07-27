package service

import (
	"context"
	"fmt"
	"time"

	"github.com/rajabhishekmaurya/ecom/internal/model"
	"github.com/rajabhishekmaurya/ecom/libs/config"
	pb "github.com/rajabhishekmaurya/ecom/libs/proto/payment"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PaymentService struct {
	cfg *config.Config
}

func NewPaymentService(cfg *config.Config) *PaymentService {
	return &PaymentService{cfg: cfg}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	txnID := fmt.Sprintf("TXN-%d", time.Now().Unix())

	producer := NewKafkaProducer(s.cfg)
	defer producer.Close()

	event := &model.EventPayment{
		OrderID:       req.OrderId,
		TransactionID: txnID,
		Amount:        req.Amount,
		Status:        "SUCCESS",
	}

	if err := producer.Publish(event); err != nil {
		return nil, err
	}

	time.Sleep(500 * time.Millisecond)

	return &pb.PaymentResponse{
		Success:       true,
		TransactionId: txnID,
	}, nil
}

type PaymentClient struct {
	client pb.PaymentServiceClient
	conn   *grpc.ClientConn

	cfg *config.Config
}

func NewPaymentClient(cfg *config.Config) (*PaymentClient, error) {
	conn, err := grpc.NewClient(
		cfg.Service.Payment,
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		return nil, err
	}

	client := pb.NewPaymentServiceClient(conn)

	return &PaymentClient{
		client: client,
		conn:   conn,
		cfg:    cfg,
	}, nil
}

func (p *PaymentClient) ProcessPayment(orderID string, amount float64) (*pb.PaymentResponse, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return p.client.ProcessPayment(ctx, &pb.PaymentRequest{
		OrderId: orderID,
		Amount:  amount,
	})
}

func (p *PaymentClient) Close() error {
	return p.conn.Close()
}
