package cmd

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	// SQL server driver
	_ "github.com/denisenkom/go-mssqldb"

	"github.com/Philipid3s/go-rpc/pkg/protocol/grpc"
	"github.com/Philipid3s/go-rpc/pkg/service/v1"
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort string

	// DB Datastore parameters section
	// DatastoreDBHost is host of database
	DatastoreDBHost string
	// DatastoreDBPort is host port of database
	DatastoreDBPort string
	// DatastoreDBUser is username to connect to database
	DatastoreDBUser string
	// DatastoreDBPassword password to connect to database
	DatastoreDBPassword string
	// DatastoreDBSchema is schema of database
	DatastoreDBSchema string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "", "gRPC port to bind")
	flag.StringVar(&cfg.DatastoreDBHost, "db-host", "", "Database host")
	flag.StringVar(&cfg.DatastoreDBPort, "db-port", "", "Database port")
	flag.StringVar(&cfg.DatastoreDBUser, "db-user", "", "Database user")
	flag.StringVar(&cfg.DatastoreDBPassword, "db-password", "", "Database password")
	flag.StringVar(&cfg.DatastoreDBSchema, "db-schema", "", "Database schema")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}

	// add SQL Server driver
	// from denisenkom/go-mssqldb

	dsn := fmt.Sprintf("server=%s;user id=%s;password=%s;port=%s;database=%s;",
		cfg.DatastoreDBHost,
		cfg.DatastoreDBUser,
		cfg.DatastoreDBPassword,
		cfg.DatastoreDBPort,
		cfg.DatastoreDBSchema)

	db, err := sql.Open("mssql", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	v1API := v1.NewContactServiceServer(db)

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort)
}
