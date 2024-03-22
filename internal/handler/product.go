package handler

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ProductAlreadyExists = errors.New("product already exists")

	sqlInsert = "INSERT INTO products (name, price) VALUES ($1, $2) RETURNING *"
	sqlFetch  = "SELECT * FROM products WHERE id = $1"
	sqlExists = "SELECT count(id) FROM products WHERE name = $1"
)

type Product struct {
	ID    int64   `db:"id"`
	Name  string  `db:"name"`
	Price float64 `db:"price"`
}

type NewProduct struct {
	Name  string
	Price float64
}

type ProductService interface {
	Add(ctx context.Context, product *NewProduct) (*Product, error)
	Get(ctx context.Context, id int64) (*Product, error)
}

type productService struct {
	db *pgxpool.Pool
}

func newProductService(db *pgxpool.Pool) ProductService {
	return &productService{db: db}
}

func (ps productService) Add(ctx context.Context, p *NewProduct) (*Product, error) {
	//var count int
	//_ = ps.db.QueryRow(ctx, sqlExists, p.Name).Scan(&count)
	//if count > 0 {
	//	return nil, ProductAlreadyExists
	//}

	tx, err := ps.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	rows, _ := tx.Query(ctx, sqlInsert, p.Name, p.Price)
	product, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Product])
	var pge *pgconn.PgError
	if err != nil {
		if errors.As(err, &pge) {
			if pge.Code == "23505" {
				return nil, ProductAlreadyExists
			}
		}
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (ps productService) Get(ctx context.Context, id int64) (*Product, error) {
	rows, err := ps.db.Query(ctx, sqlFetch, id)
	if err != nil {
		return nil, err
	}

	product, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[Product])
	if err != nil {
		return nil, err
	}

	return product, nil
}
