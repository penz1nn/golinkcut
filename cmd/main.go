package main

import (
	"flag"
	"golinkcut/internal/config"
	"golinkcut/internal/link"
)

// TODO: comments for constants and functions

func main() {

	cfg := buildConfig()
	repo := link.NewStorage(cfg)
	uc := link.NewUseCase(repo, cfg)
	go runGrpcServer(uc, cfg)
	runRestApi(uc, cfg)
}

// buildConfig uses the flag package to create a config.Config object from the
// command-line arguments
func buildConfig() config.Config {
	// define flags
	enableDebugFlag := flag.Bool("debug", false, "enable debug logs")
	memoryDbFlag := flag.Bool("memory", false, "use in-memory storage")
	validateFlag := flag.Bool("validate", true, "validate submitted URLs")
	grpcPort := flag.String("grpc-port", "50051", "listen port for GRPC server")
	httpPort := flag.String("http-port", "8080", "listen port for REST API")
	httpHost := flag.String("http-host", "localhost", "what to use as REST API hostname")
	redirect := flag.Bool("redirect", false, "whether or not to redirect on GET calls to short links (REST API)")
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
		cfg["memory"] = true
	} else {
		cfg["memory"] = false
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
	cfg["httpHost"] = *httpHost
	cfg["redirect"] = *redirect
	return cfg
}
