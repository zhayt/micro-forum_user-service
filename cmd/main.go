package main

import (
	"fmt"
	"github.com/zhayt/user-service/config"
	"github.com/zhayt/user-service/logger"
	pb "github.com/zhayt/user-service/proto"
	"github.com/zhayt/user-service/service"
	"github.com/zhayt/user-service/storage"
	"github.com/zhayt/user-service/storage/postgre"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"sync"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// config
	var once sync.Once
	once.Do(config.PrepareENV)

	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}

	// logger
	l, err := logger.NewLogger(cfg)
	if err != nil {
		return err
	}

	defer func(logger2 *zap.Logger) {
		err := l.Sync()
		if err != nil {
			log.Fatalln(err)
		}
	}(l)

	// repo

	db, err := postgre.Dial(makeDSN(cfg))
	if err != nil {
		return err
	}
	defer db.Close()

	repo := storage.NewStorage(db, l)

	// usecases
	validate := service.NewValidateService()
	userService := service.NewUserService(repo, validate, l)

	// init
	lis, err := net.Listen("tcp", net.JoinHostPort("", cfg.AppPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)

	reflection.Register(grpcServer)

	// register

	pb.RegisterUserServiceServer(grpcServer, userService)

	//Start GRPC server
	log.Printf("Start GRPC server on address: %s", cfg.AppPort)
	err = grpcServer.Serve(lis)
	if err != nil {
		l.Fatal("Serve", zap.Error(err))
	}
	return nil
}

func makeDSN(cfg *config.Config) string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.TZ)
}
