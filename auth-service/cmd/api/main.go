package main

import (
	"auth-service/data"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// external dependencies:
// go get github.com/go-chi/chi/v5
// go get github.com/go-chi/cors
// go get github.com/jackc/pgconn
// go get github.com/jackc/pgx/v4

// adding grpc:
// go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.27
// go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@1.2

// protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative auth.proto
// go get google.golang.org/grpc

// how to access postgres in the container:
// docker exec -it project_postgres_1 psql -U postgres

const webPort = "80"

type App struct {
	DB       *sql.DB
	UserRepo *data.PostgresUserRepo
}

func main() {
	fmt.Println("starting authentication server")

	conn := connectToDB()
	if conn == nil {
		log.Panic("Cannot connect to Postgres")
	}

	app := App{
		DB:       conn,
		UserRepo: data.NewPostgresUserRepo(conn),
	}

	// start the gRPC server
	go app.gRPCListen()

	fmt.Printf("Listening on port %s\n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")
	retryCount := 0

	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not ready yet...")
			retryCount++
		} else {
			log.Println("Connected to Postgres.")
			return conn
		}

		// give up
		if retryCount > 10 {
			log.Println(err)
			return nil
		}

		//
		time.Sleep(2 * time.Second)
	}
}
