package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/FRahimov84/Mux/pkg/mux"
	"github.com/FRahimov84/ProductService/cmd/product/app"
	"github.com/FRahimov84/ProductService/pkg/core/product"
	"github.com/FRahimov84/di/pkg/di"
	"github.com/FRahimov84/myJwt/pkg/jwt"
	"github.com/jackc/pgx/v4/pgxpool"
	"net"
	"net/http"
	"os"
)

var (
	host = flag.String("host", "", "Server host")
	port = flag.String("port", "", "Server port")
	dsn  = flag.String("dsn", "", "Postgres DSN")
)
//-host 0.0.0.0 -port 9999 -dsn postgres://user:pass@localhost:5430/product
const (
	envHost = "HOST"
	envPort = "PORT"
	envDSN  = "DATABASE_URL"
)

type DSN string

func main() {
	flag.Parse()
	serverHost := checkENV(envHost, *host)
	serverPort := checkENV(envPort, *port)
	serverDsn := checkENV(envDSN, *dsn)
	addr := net.JoinHostPort(serverHost, serverPort)
	secret := jwt.Secret("secret")
	start(addr, serverDsn, secret)
}
func checkENV(env string, loc string) string {
	str, ok := os.LookupEnv(env)
	if !ok {
		return loc
	}
	return str
}
func start(addr string, dsn string,  secret jwt.Secret) {
	container := di.NewContainer()
	container.Provide(
		app.NewServer,
		mux.NewExactMux,
		product.NewService,
		func() DSN { return DSN(dsn) },
		func() jwt.Secret { return secret },
		func(dsn DSN) *pgxpool.Pool {
			pool, err := pgxpool.Connect(context.Background(), string(dsn))
			if err != nil {
				panic(fmt.Errorf("can't create pool: %w", err))
			}
			return pool
		},
	)

	container.Start()

	var appServer *app.Server
	container.Component(&appServer)
	panic(http.ListenAndServe(addr, appServer))
}
