package usecase

import (
	"context"
	"github.com/bxcodec/go-clean-arch/domain/article"
	"github.com/bxcodec/go-clean-arch/domain/author"
	"time"

	"github.com/bxcodec/go-clean-arch/domain"
)

type articleUsecase struct {
	articleRepo    article.ArticleRepository
	authorRepo     author.AuthorRepository
	contextTimeout time.Duration
}

// NewArticleUsecase will create new an articleUsecase object representation of domain.ArticleUsecase interface
func NewArticleUsecase(a article.ArticleRepository, ar author.AuthorRepository, timeout time.Duration) article.ArticleUsecase {
	return &articleUsecase{
		articleRepo:    a,
		authorRepo:     ar,
		contextTimeout: timeout,
	}
}

func (a *articleUsecase) Fetch(c context.Context, cursor string, num int64) (res []article.Article, nextCursor string, err error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, nextCursor, err = a.articleRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	return
}

func (a *articleUsecase) GetByID(c context.Context, id int64) (res article.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	res, err = a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}

	return
}

func (a *articleUsecase) Update(c context.Context, ar *article.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	ar.UpdatedAt = time.Now()
	return a.articleRepo.Update(ctx, ar)
}

func (a *articleUsecase) GetByTitle(c context.Context, title string) (res article.Article, err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	res, err = a.articleRepo.GetByTitle(ctx, title)
	if err != nil {
		return
	}

	return
}

func (a *articleUsecase) Store(c context.Context, m *article.Article) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, _ := a.GetByTitle(ctx, m.Title)
	if existedArticle != (article.Article{}) {
		return domain.ErrConflict
	}

	err = a.articleRepo.Store(ctx, m)
	return
}

func (a *articleUsecase) Delete(c context.Context, id int64) (err error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()
	existedArticle, err := a.articleRepo.GetByID(ctx, id)
	if err != nil {
		return
	}
	if existedArticle == (article.Article{}) {
		return domain.ErrNotFound
	}
	return a.articleRepo.Delete(ctx, id)
}
