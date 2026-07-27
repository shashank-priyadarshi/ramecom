package internal

import (
	"fmt"

	"github.com/rajabhishekmaurya/ecom/internal/handler"
	"github.com/rajabhishekmaurya/ecom/internal/repo"
	"github.com/rajabhishekmaurya/ecom/internal/server"
	"github.com/rajabhishekmaurya/ecom/internal/service"
	"github.com/rajabhishekmaurya/ecom/libs/config"
	"github.com/rajabhishekmaurya/ecom/libs/db"
	"github.com/rajabhishekmaurya/ecom/libs/schema"
)

type App interface {
	Start() error
}

// Apps are:
// 	gateway
// 	auth
// 	user
// 	order
// 	payment
// 	notification

func New(app *string) (App, error) {
	if app == nil || *app == "" || *app == "unknown" {
		return nil, fmt.Errorf("app flag is required (gateway|auth|user|order|payment|notification)")
	}

	cfg := config.LoadForApp(*app)

	switch *app {
	case "gateway":
		h := handler.NewGatewayHandler(cfg)
		return server.New(cfg, h)

	case "auth":
		database, err := db.NewDB(cfg)
		if err != nil {
			return nil, err
		}
		if err := schema.CreateTables(database); err != nil {
			return nil, err
		}
		userRepo := repo.NewUserRepository(database)
		svc := service.NewAuthService(userRepo, cfg)
		h := handler.NewAuthHandler(svc)
		return server.New(cfg, h)

	case "user":
		database, err := db.NewDB(cfg)
		if err != nil {
			return nil, err
		}
		if err := schema.CreateTables(database); err != nil {
			return nil, err
		}
		userRepo := repo.NewUserRepository(database)
		svc := service.NewUserService(userRepo)
		h := handler.NewUserHandler(svc)
		return server.New(cfg, h)

	case "order":
		svc, err := service.NewOrderService(cfg)
		if err != nil {
			return nil, err
		}
		h := handler.NewOrderHandler(svc)
		return server.New(cfg, h)

	case "payment":
		svc := service.NewPaymentService(cfg)
		h := handler.NewPaymentHandler(svc)
		return server.NewPayment(cfg, h), nil

	case "notification":
		consumer := service.NewKafkaConsumer()
		return server.NewNotification(consumer), nil

	default:
		return nil, fmt.Errorf("unknown app %q", *app)
	}
}
