package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/GabrieldeFreire/clean-architecture/configs"
	"github.com/GabrieldeFreire/clean-architecture/internal/event/handler"
	"github.com/GabrieldeFreire/clean-architecture/internal/infra/graph"
	"github.com/GabrieldeFreire/clean-architecture/internal/infra/grpc/pb"
	"github.com/GabrieldeFreire/clean-architecture/internal/infra/grpc/service"
	"github.com/GabrieldeFreire/clean-architecture/internal/infra/web/webserver"
	"github.com/GabrieldeFreire/clean-architecture/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

var (
	config     *configs.Conf
	configsErr error
)

func main() {
	config, configsErr = configs.LoadConfig(".")
	if configsErr != nil {
		panic(configsErr)
	}
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err := sql.Open(config.DBDriver, dbUrl)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	createOrderTable(db)
	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrdersUseCase := NewListOrderUseCase(db)

	webserver := webserver.NewWebServer(fmt.Sprintf(":%s", config.WebServerPort))
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	webserver.AddHandler("POST", "/order", webOrderHandler.Create)
	webserver.AddHandler("GET", "/order", webOrderHandler.ListOrders)

	fmt.Println("Starting web server on port", config.WebServerPort)
	go webserver.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrdersUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", config.GRPCServerPort)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{CreateOrderUseCase: *createOrderUseCase, ListOrdersUseCase: *listOrdersUseCase},
			},
		),
	)

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", config.GraphQLServerPort)
	http.ListenAndServe(":"+config.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://guest:guest@%s:5672/", config.RABBITMQ_HOST))
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}

func createOrderTable(db *sql.DB) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS orders (
		id VARCHAR(36) PRIMARY KEY,
		price NUMERIC,
		tax NUMERIC,
		final_price NUMERIC
	)`)
	if err != nil {
		panic(err)
	}
}
