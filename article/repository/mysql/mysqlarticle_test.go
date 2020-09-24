package mysql_test

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain/article"
	"github.com/bxcodec/go-clean-arch/domain/author"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/bxcodec/go-clean-arch/article/repository"
	articleMysqlRepo "github.com/bxcodec/go-clean-arch/article/repository/mysql"
)

func TestFetch(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mockArticles := []article.Article{
		article.Article{
			ID: 1, Title: "title 1", Content: "content 1",
			Author: author.Author{ID: 1}, UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
		article.Article{
			ID: 2, Title: "title 2", Content: "content 2",
			Author: author.Author{ID: 1}, UpdatedAt: time.Now(), CreatedAt: time.Now(),
		},
	}

	rows := sqlmock.NewRows([]string{"ar.id", "ar.title", "ar.content", "ar.author_id", "ar.created_at", "ar.updated_at", "a.name", "a.created_at", "a.updated_at"}).
		AddRow(mockArticles[0].ID, mockArticles[0].Title, mockArticles[0].Content,
			mockArticles[0].Author.ID, mockArticles[0].UpdatedAt, mockArticles[0].CreatedAt, mockArticles[0].Author.Name,
			mockArticles[0].Author.CreatedAt, mockArticles[0].Author.UpdatedAt).
		AddRow(mockArticles[1].ID, mockArticles[1].Title, mockArticles[1].Content,
			mockArticles[1].Author.ID, mockArticles[1].UpdatedAt, mockArticles[1].CreatedAt, mockArticles[1].Author.Name,
			mockArticles[1].Author.CreatedAt, mockArticles[1].Author.UpdatedAt)

	query := "SELECT ar.id,ar.title,ar.content, ar.author_id, ar.created_at, ar.updated_at, a.name, a.created_at, a.updated_at FROM article ar JOIN author a ON a.id = ar.author_id WHERE created_at > \\? ORDER BY created_at LIMIT \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := articleMysqlRepo.NewMysqlArticleRepository(db)
	cursor := repository.EncodeCursor(mockArticles[1].CreatedAt)
	num := int64(2)
	list, nextCursor, err := a.Fetch(context.TODO(), cursor, num)
	assert.NotEmpty(t, nextCursor)
	assert.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestGetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at", "name", "created_at", "updated_at"}).
		AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now(), "Author Name", time.Now(), time.Now())

	query := "SELECT ar.id,ar.title,ar.content, ar.author_id, ar.created_at, ar.updated_at, a.name, a.created_at, a.updated_at FROM article ar JOIN author a ON a.id = ar.author_id WHERE ID = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := articleMysqlRepo.NewMysqlArticleRepository(db)

	num := int64(1)
	anArticle, err := a.GetByID(context.TODO(), num)
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}

func TestStore(t *testing.T) {
	now := time.Now()
	ar := &article.Article{
		Title:     "Judul",
		Content:   "Content",
		CreatedAt: now,
		UpdatedAt: now,
		Author: author.Author{
			ID:   1,
			Name: "Iman Tumorang",
		},
	}
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "INSERT  article SET title=\\? , content=\\? , author_id=\\?, updated_at=\\? , created_at=\\?"
	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Title, ar.Content, ar.Author.ID, ar.CreatedAt, ar.UpdatedAt).WillReturnResult(sqlmock.NewResult(12, 1))

	a := articleMysqlRepo.NewMysqlArticleRepository(db)

	err = a.Store(context.TODO(), ar)
	assert.NoError(t, err)
	assert.Equal(t, int64(12), ar.ID)
}

func TestGetByTitle(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	rows := sqlmock.NewRows([]string{"id", "title", "content", "author_id", "created_at", "updated_at", "name", "created_at", "updated_at"}).
		AddRow(1, "title 1", "Content 1", 1, time.Now(), time.Now(), "Author Name", time.Now(), time.Now())

	query := "SELECT ar.id,ar.title,ar.content, ar.author_id, ar.created_at, ar.updated_at, a.name, a.created_at, a.updated_at FROM article ar JOIN author a on a.id = ar.author_id WHERE title = \\?"

	mock.ExpectQuery(query).WillReturnRows(rows)
	a := articleMysqlRepo.NewMysqlArticleRepository(db)

	title := "title 1"
	anArticle, err := a.GetByTitle(context.TODO(), title)
	assert.NoError(t, err)
	assert.NotNil(t, anArticle)
}

func TestDelete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "DELETE FROM article WHERE id = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(12).WillReturnResult(sqlmock.NewResult(12, 1))

	a := articleMysqlRepo.NewMysqlArticleRepository(db)

	num := int64(12)
	err = a.Delete(context.TODO(), num)
	assert.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	now := time.Now()
	ar := &article.Article{
		ID:        12,
		Title:     "Judul",
		Content:   "Content",
		CreatedAt: now,
		UpdatedAt: now,
		Author: author.Author{
			ID:   1,
			Name: "Iman Tumorang",
		},
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	query := "UPDATE article set title=\\?, content=\\?, author_id=\\?, updated_at=\\? WHERE ID = \\?"

	prep := mock.ExpectPrepare(query)
	prep.ExpectExec().WithArgs(ar.Title, ar.Content, ar.Author.ID, ar.UpdatedAt, ar.ID).WillReturnResult(sqlmock.NewResult(12, 1))

	a := articleMysqlRepo.NewMysqlArticleRepository(db)

	err = a.Update(context.TODO(), ar)
	assert.NoError(t, err)
}
