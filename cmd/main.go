package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kartik7120/booking_moviedb_service/cmd/api"
	"github.com/kartik7120/booking_moviedb_service/cmd/consumers"
	movie "github.com/kartik7120/booking_moviedb_service/cmd/grpcServer"
	"github.com/kartik7120/booking_moviedb_service/cmd/helper"
	"github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	err := godotenv.Load()

	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)

	if os.Getenv("ENV") == "production" {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	if err != nil {
		log.Error("Error loading .env file")
		panic(err)
	}

	lis, err := net.Listen("tcp", ":1102")

	if err != nil {
		log.Error("error starting the server")
		panic(err)
	}

	// Creating a new grpc server

	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	// Connect to database

	DB, err := helper.ConnectToDB()

	if err != nil {
		log.Error("Error connecting to database")
		panic(err)
	}

	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")

	if err != nil {
		log.Error("error connecting to to rabbitmq")
		os.Exit(1)
		return
	}
	defer conn.Close()

	ch, err := conn.Channel()

	if err != nil {
		log.Error("error opening a channel")
		os.Exit(1)
		return
	}

	defer ch.Close()

	consumer := consumers.NewConsumer(ch)

	go func() {
		log.Info("Listening on incoming message from Send_Mail_Consumer")
		err := consumer.Send_Mail_Consumer()
		if err != nil {
			log.Error("failed to consume send mail messages")
			os.Exit(1)
			return
		}

	}()

	moviedbObj := api.NewMovieDB()

	moviedbObj.DB.Conn = DB

	movie.RegisterMovieDBServiceServer(grpcServer, &api.MoviedbService{
		MovieDB: moviedbObj,
	})

	if os.Getenv("ENV") != "production" {
		reflection.Register(grpcServer)
	}

	go func() {
		log.Info("Starting the moviedb service")
		if err := grpcServer.Serve(lis); err != nil {
			log.Error("error starting the server")
			panic(err)
		}
	}()

	<-signalChan

	log.Info("Stopping the server")
	grpcServer.GracefulStop()
}
