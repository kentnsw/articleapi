package model

import (
	store "github.com/kentnsw/artical-api/storage/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (article *NewArticle) StoreArtical() *store.Article {
	if article == nil {
		return &store.Article{}
	}
	return &store.Article{
		Title: article.Title,
		Date:  primitive.NewDateTimeFromTime(article.Date),
		Body:  article.Body,
		Tags:  article.Tags,
	}
}

func NewGqlArticle(art *store.Article) *Article {
	if art == nil {
		return &Article{}
	}
	return &Article{
		ID:    art.ID.Hex(),
		Title: art.Title,
		Date:  art.Date.Time(),
		Body:  art.Body,
		Tags:  art.Tags,
	}
}
