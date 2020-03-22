package app

import (
	"github.com/FRahimov84/ProductService/pkg/core/token"
	"github.com/FRahimov84/ProductService/pkg/mux/middleware/authenticated"
	"github.com/FRahimov84/ProductService/pkg/mux/middleware/authorized"
	"github.com/FRahimov84/ProductService/pkg/mux/middleware/jwt"
	"github.com/FRahimov84/ProductService/pkg/mux/middleware/logger"
	"reflect"
)

func (s Server) InitRoutes() {


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
		"/api/products/{id}",
		s.handProduct(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		authorized.Authorized([]string{"Admin"}, jwt.FromContext),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("post product"),
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