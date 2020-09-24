package mysql

import (
	"context"
	"database/sql"
	"github.com/bxcodec/go-clean-arch/domain/author"
)

type mysqlAuthorRepo struct {
	DB *sql.DB
}

// NewMysqlAuthorRepository will create an implementation of author.Repository
func NewMysqlAuthorRepository(db *sql.DB) author.AuthorRepository {
	return &mysqlAuthorRepo{
		DB: db,
	}
}

func (m *mysqlAuthorRepo) getOne(ctx context.Context, query string, args ...interface{}) (res author.Author, err error) {
	stmt, err := m.DB.PrepareContext(ctx, query)
	if err != nil {
		return author.Author{}, err
	}
	row := stmt.QueryRowContext(ctx, args...)
	res = author.Author{}

	err = row.Scan(
		&res.ID,
		&res.Name,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	return
}

func (m *mysqlAuthorRepo) GetByID(ctx context.Context, id int64) (author.Author, error) {
	query := `SELECT id, name, created_at, updated_at FROM author WHERE id=?`
	return m.getOne(ctx, query, id)
}
