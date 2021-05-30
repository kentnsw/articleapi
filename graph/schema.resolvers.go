package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"log"
	"time"

	"github.com/kentnsw/artical-api/graph/generated"
	"github.com/kentnsw/artical-api/graph/model"
	articleStore "github.com/kentnsw/artical-api/storage/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (r *mutationResolver) CreateArticle(ctx context.Context, article model.NewArticle) (*model.Article, error) {
	log.Println("CreateArticle() with ", article)
	art := articleStore.Article{
		Title: article.Title,
		Date:  primitive.NewDateTimeFromTime(article.Date),
		Body:  article.Body,
		Tags:  article.Tags,
	}
	if _, err := art.Save(); err != nil {
		return nil, err
	}

	return &model.Article{
		ID:    art.ID.Hex(),
		Title: art.Title,
		Date:  art.Date.Time(),
		Body:  art.Body,
		Tags:  art.Tags,
	}, nil
}

func (r *queryResolver) Article(ctx context.Context, id string) (*model.Article, error) {
	log.Println("Article() find article by id ", id)
	ctx, cencel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cencel()
	art, err := articleStore.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &model.Article{
		ID:    art.ID.Hex(),
		Title: art.Title,
		Date:  art.Date.Time(),
		Body:  art.Body,
		Tags:  art.Tags,
	}, nil
}

func (r *queryResolver) Articles(ctx context.Context) ([]*model.Article, error) {
	ctx, cencel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cencel()

	arts, err := articleStore.Find(ctx, &articleStore.Filter{})
	if err != nil {
		return nil, err
	}

	var res []*model.Article
	for _, art := range arts {
		gqlArt := model.Article{
			ID:    art.ID.Hex(),
			Title: art.Title,
			Date:  art.Date.Time(),
			Body:  art.Body,
			Tags:  art.Tags,
		}
		res = append(res, &gqlArt)
	}
	log.Printf("Articles() found %d artitles", len(res))
	return res, nil
}

func (r *queryResolver) ArticlesByTag(ctx context.Context, filter model.ArticleFilter) (*model.ArticlesByTag, error) {
	var default_limit int = 10

	log.Println("ArticlesByTag() with filter ", filter)
	artFilter := &articleStore.Filter{Tags: filter.Tag}
	if filter.Date != nil {
		artFilter.Date = primitive.NewDateTimeFromTime(*filter.Date)
	}
	if filter.Limit == nil {
		filter.Limit = &default_limit
	}

	var res model.ArticlesByTag
	res.Tag = filter.Tag

	ctx, cencel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cencel()
	arts, err := articleStore.FindRelatedTags(ctx, artFilter)
	if err != nil {
		return nil, err
	}

	res.Count = len(arts)

	tagsMap := map[string]bool{}
	for i, art := range arts {
		if i < *filter.Limit {
			res.Articles = append(res.Articles, art.ID.Hex())
		}

		for _, tag := range art.Tags {
			if _, ok := tagsMap[tag]; !ok {
				tagsMap[tag] = true
			}
		}
	}
	delete(tagsMap, filter.Tag)

	res.RelatedTags = make([]string, len(tagsMap))
	i := 0
	for k := range tagsMap {
		res.RelatedTags[i] = k
		i++
	}

	return &res, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
