package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kartik7120/booking_moviedb_service/cmd/api"
	movie "github.com/kartik7120/booking_moviedb_service/cmd/grpcServer"
	"github.com/kartik7120/booking_moviedb_service/cmd/helper"
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
