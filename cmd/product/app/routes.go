package app

import (
	"context"
	"errors"
	"github.com/FRahimov84/ProductService/pkg/core/token"
	"github.com/FRahimov84/ProductService/pkg/mux/middleware/authenticated"
	"github.com/FRahimov84/ProductService/pkg/mux/middleware/authorized"
	"github.com/FRahimov84/ProductService/pkg/mux/middleware/jwt"
	"github.com/FRahimov84/ProductService/pkg/mux/middleware/logger"
	"reflect"
)

func (s Server) InitRoutes() {

	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		panic(errors.New("can't create database"))
	}
	defer conn.Release()
	_, err = conn.Exec(context.Background(), `
CREATE TABLE if not exists products (
             id BIGSERIAL PRIMARY KEY,
             name TEXT NOT NULL unique,
             description TEXT NOT NULL,
             price Integer check ( price>=0 ) NOT NULL,
             pic varchar NOT NULL,
             removed BOOLEAN DEFAULT FALSE
);
`)
	if err != nil {
		panic(errors.New("can't create database"))
	}

	s.router.GET(
		"/api/products",
		s.handleProductList(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("get list"),
	)

	s.router.GET(
		"/api/products/{id}",
		s.handleProductByID(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("get product by id"),
	)

	s.router.POST(
		"/api/products/new",
		s.handleNewProduct(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		authorized.Authorized([]string{"Admin"}, jwt.FromContext),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("post new product"),
	)

	s.router.DELETE(
		"/api/products/{id}",
		s.handleDeleteProduct(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		authorized.Authorized([]string{"Admin"}, jwt.FromContext),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("delete product"),
	)


}