package main

import (
	"flag"
	"golinkcut/api/proto"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
	"golinkcut/pkg/log"
	"google.golang.org/grpc"
	"net"
)

func main() {
	cfg := buildConfig()
	repo := link.NewStorage(cfg)
	logger := log.NewWithConfig(cfg)
	uc := link.NewUseCase(repo, logger, cfg)
	runGrpcServer(uc, logger, cfg)
}

func runGrpcServer(uc link.UseCase, logger log.Logger, cfg config.Config) {
	srv := proto.NewGrpcServer(uc)
	lis, err := net.Listen("tcp", ":"+cfg["grpcPort"].(string))
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	proto.RegisterLinkServiceServer(s, &srv)
	logger.Infof("Server started at port %s", cfg["grpcPort"].(string))
	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

// buildConfig uses the flag package to create a config.Config object from the
// command-line arguments
func buildConfig() config.Config {
	// define flags
	enableDebugFlag := flag.Bool("debug", false, "enable debug logs")
	memoryDbFlag := flag.Bool("memory", true, "use in-memory storage")
	validateFlag := flag.Bool("validate", false, "validate submitted URLs")
	grpcPort := flag.String("grpc-port", "50051", "listen port for GRPC server")
	httpPort := flag.String("http-port", "8080", "listen port for HTTP server (REST API)")
	dbHost := flag.String("db-host", "localhost", "PostgreSQL host")
	dbUser := flag.String("db-user", "golinkcut", "PostgreSQL user")
	dbPassword := flag.String("db-password", "example", "PostgreSQL password")
	dbName := flag.String("db-name", "golinkcut", "PostgreSQL database name")
	dbPort := flag.String("db-port", "5432", "PostgreSQL port")
	dbTz := flag.String("db-timezone", "Europe/Moscow", "PostgreSQL timezone")

	flag.Parse()

	cfg := config.Config{}

	if *enableDebugFlag {
		cfg["debug"] = true
	} else {
		cfg["debug"] = false
	}

	if *memoryDbFlag {
		cfg["db"] = "memory"
	} else {
		cfg["db"] = map[string]string{
			"host":     *dbHost,
			"user":     *dbUser,
			"password": *dbPassword,
			"dbname":   *dbName,
			"port":     *dbPort,
			"tz":       *dbTz,
		}
	}

	if *validateFlag {
		cfg["validate"] = true
	} else {
		cfg["validate"] = false
	}

	cfg["grpcPort"] = *grpcPort
	cfg["httpPort"] = *httpPort
	return cfg
}
