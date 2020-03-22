package product

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
)

type Service struct {
	pool *pgxpool.Pool
}

func NewService(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}



func (s *Service) Start() {

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
}

type Product struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	Pic         string `json:"pic"`
}

func (s *Service) AddNewProduct(ctx context.Context, prod Product) (err error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return
	}
	defer conn.Release()
	_, err = conn.Exec(ctx, `INSERT INTO products(name, description, price, pic)
VALUES ($1, $2, $3, $4);`, prod.Name, prod.Description, prod.Price, prod.Pic)
	if err != nil {
		return
	}
	return nil
}

func (s *Service) ProductList(ctx context.Context) (list []Product, err error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()
	rows, err := conn.Query(ctx,
		`select id, name, description, price, pic from products where removed=false;`)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		item := Product{}
		err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Price, &item.Pic)
		if err != nil {
			return nil, errors.New("can't scan row from rows")
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.New("rows error!")
	}
	return
}

func (s *Service) RemoveByID(ctx context.Context, id int64) (err error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return errors.New("can't connect to database!")
	}
	defer conn.Release()
	_, err = conn.Exec(ctx, `update products set removed = true where id = $1`, id)
	if err != nil {
		return errors.New(fmt.Sprintf("can't remove from database product (id: %d)!", id))
	}
	return nil
}

func (s *Service) ProductByID(ctx context.Context, id int64) (prod Product, err error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return Product{}, errors.New("can't connect to database!")
	}
	defer conn.Release()
	err = conn.QueryRow(ctx, `select id, name, description, price, pic from products where id=$1`,
		id).Scan(&prod.ID, &prod.Name, &prod.Description, &prod.Price, &prod.Pic)
	if err != nil {
		return Product{}, errors.New(fmt.Sprintf("can't remove from database burger (id: %d)!", id))
	}
	return
}

func (s *Service) UpdateProduct(ctx context.Context,id int64, prod Product) (err error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return errors.New("can't connect to database!")
	}
	defer conn.Release()
	begin, err := conn.Begin(ctx)
	if err != nil {
		return errors.New("can't connect to database!")
	}
	defer func() {
		if err != nil {
			err2 := begin.Rollback(ctx)
			log.Printf("can't rollback some err %v", err2)
			return
		}
		err2 := begin.Commit(ctx)
		log.Printf("can't commit tranzaction err %v", err2)
	}()
	if prod.Name != "" {
		_, err = begin.Exec(ctx, `update products set name = $2 where id = $1`, id, prod.Name)
		return
	}
	if prod.Description !="" {
		_, err = begin.Exec(ctx, `update products set description = $2 where id = $1`, id, prod.Description)
		return
	}
	if prod.Price != -1 {
		_, err = begin.Exec(ctx, `update products set price = $2 where id = $1`, id, prod.Price)
		return
	}
	if prod.Pic != "" {
		_, err = begin.Exec(ctx, `update products set pic = $2 where id = $1`, id, prod.Pic)
		return
	}
	return nil
}
