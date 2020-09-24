package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/bxcodec/go-clean-arch/domain/article"
	author2 "github.com/bxcodec/go-clean-arch/domain/author"

	"github.com/sirupsen/logrus"

	"github.com/bxcodec/go-clean-arch/article/repository"
	"github.com/bxcodec/go-clean-arch/domain"
)

type mysqlArticleRepository struct {
	Conn *sql.DB
}

// NewMysqlArticleRepository will create an object that represent the article.Repository interface
func NewMysqlArticleRepository(Conn *sql.DB) article.ArticleRepository {
	return &mysqlArticleRepository{Conn}
}

func (m *mysqlArticleRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []article.Article, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]article.Article, 0)
	for rows.Next() {
		t := article.Article{}
		author := author2.Author{}
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&author.ID,
			&t.UpdatedAt,
			&t.CreatedAt,
			&author.Name,
			&author.CreatedAt,
			&author.UpdatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		t.Author = author
		result = append(result, t)
	}

	return result, nil
}

func (m *mysqlArticleRepository) Fetch(ctx context.Context, cursor string, num int64) (res []article.Article, nextCursor string, err error) {
	query := `SELECT ar.id,ar.title,ar.content, ar.author_id, ar.created_at, ar.updated_at, a.name, a.created_at, a.updated_at
  						FROM article ar JOIN author a ON a.id = ar.author_id WHERE created_at > ? ORDER BY created_at LIMIT ? `

	decodedCursor, err := repository.DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", domain.ErrBadParamInput
	}

	res, err = m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}

	if len(res) == int(num) {
		nextCursor = repository.EncodeCursor(res[len(res)-1].CreatedAt)
	}

	return
}
func (m *mysqlArticleRepository) GetByID(ctx context.Context, id int64) (res article.Article, err error) {
	query := `SELECT ar.id,ar.title,ar.content, ar.author_id, ar.created_at, ar.updated_at, a.name, a.created_at, a.updated_at
  						FROM article ar JOIN author a ON a.id = ar.author_id WHERE ID = ?`

	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return article.Article{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

func (m *mysqlArticleRepository) GetByTitle(ctx context.Context, title string) (res article.Article, err error) {
	query := `SELECT ar.id,ar.title,ar.content, ar.author_id, ar.created_at, ar.updated_at, a.name, a.created_at, a.updated_at
  						FROM article ar JOIN author a on a.id = ar.author_id WHERE title = ?`

	list, err := m.fetch(ctx, query, title)
	if err != nil {
		return
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}
	return
}

func (m *mysqlArticleRepository) Store(ctx context.Context, a *article.Article) (err error) {
	query := `INSERT  article SET title=? , content=? , author_id=?, updated_at=? , created_at=?`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, a.Title, a.Content, a.Author.ID, a.UpdatedAt, a.CreatedAt)
	if err != nil {
		return
	}
	lastID, err := res.LastInsertId()
	if err != nil {
		return
	}
	a.ID = lastID
	return
}

func (m *mysqlArticleRepository) Delete(ctx context.Context, id int64) (err error) {
	query := "DELETE FROM article WHERE id = ?"

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return
	}

	rowsAfected, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rowsAfected != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", rowsAfected)
		return
	}

	return
}
func (m *mysqlArticleRepository) Update(ctx context.Context, ar *article.Article) (err error) {
	query := `UPDATE article set title=?, content=?, author_id=?, updated_at=? WHERE ID = ?`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return
	}

	res, err := stmt.ExecContext(ctx, ar.Title, ar.Content, ar.Author.ID, ar.UpdatedAt, ar.ID)
	if err != nil {
		return
	}
	affect, err := res.RowsAffected()
	if err != nil {
		return
	}
	if affect != 1 {
		err = fmt.Errorf("Weird  Behavior. Total Affected: %d", affect)
		return
	}

	return
}
