package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"

	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"log-service/dao"
)

// external dependencies:
// go get github.com/go-chi/chi/v5
// go get github.com/go-chi/cors
// go get go.mongodb.org/mongo-driver/mongo

// adding ginkgo
// go get github.com/onsi/ginkgo/v2/ginkgo

const (
	webPort  = "80"
	rpcPort  = "5001"
	mongoURL = "mongodb://mongo:27017"
	//mongoURL = "mongodb://localhost:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type App struct {
	DataStore dao.DataStore
}

func main() {

	// connect to mongo"
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	// package level variable used by the RPC
	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := App{
		DataStore: dao.NewDataStore(mongoClient),
	}

	err = rpc.Register(new(RPCServer))
	// check the error
	go app.rpcListen()

	// TODO move into go function
	//app.serve()

	log.Println("Starting logging service: listening on port:", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Panic(err)
	}
}

func (app *App) rpcListen() error {
	log.Println("Starting RPC server on port:", rpcPort)

	listen, err := net.Listen("tcp", "0.0.0.0:"+rpcPort)
	if err != nil {
		return err
	}
	defer listen.Close()

	for {
		rpcConn, err := listen.Accept()
		if err != nil {
			continue
		}
		log.Println("Accepted RPC connection...")
		go rpc.ServeConn(rpcConn)
	}
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: getenv("mongo_user", "admin"),
		Password: getenv("mongo_password", "password"),
	})

	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}
	log.Println("Connected to mongoDB")
	return c, nil
}
