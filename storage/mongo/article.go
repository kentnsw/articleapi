package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Article struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Title string
	Date  primitive.DateTime
	Body  string
	Tags  []string
}

type Filter struct {
	ID    primitive.ObjectID `bson:"_id,omitempty"`
	Date  primitive.DateTime `bson:"date,omitempty"`
	Tags  string             `bson:"tags,omitempty"`
	Limit int64              `bson:"-"`
}

func (article *Article) Save() (id string, err error) {
	ctx, cencel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cencel()

	res, err := ArticleCollection.InsertOne(ctx, article)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	article.ID = res.InsertedID.(primitive.ObjectID)
	log.Println("saved ", res.InsertedID, article.Title)
	return article.ID.Hex(), nil
}

func InsertMany(ctx context.Context, arts []*Article) (count int, err error) {
	var artsInterface []interface{} = make([]interface{}, len(arts))
	for i, a := range arts {
		artsInterface[i] = a
	}
	res, err := ArticleCollection.InsertMany(ctx, artsInterface)
	if err != nil {
		log.Fatal(err)
		return 0, err
	}
	return len(res.InsertedIDs), nil
}

func FindById(ctx context.Context, id string) (Article, error) {
	var art Article
	artId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Article{}, err
	}

	err = ArticleCollection.FindOne(ctx, bson.M{"_id": artId}).Decode(&art)
	if err == mongo.ErrNoDocuments {
		log.Println("artitle does not exist with id ", id)
		return Article{}, err
	} else if err != nil {
		log.Fatal(err)
		return Article{}, err
	}
	return art, nil
}

func Find(ctx context.Context, f *Filter) ([]Article, error) {
	log.Println("Fine() with filter", f)
	if f == nil {
		f = &Filter{}
	}

	var arts []Article
	opt := options.Find()
	opt.SetLimit(f.Limit)
	opt.SetSort(bson.M{"_id": -1})

	cursor, err := ArticleCollection.Find(ctx, f, opt)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &arts); err != nil {
		log.Fatal(err)
		return nil, err
	}
	return arts, nil
}

func FindRelatedTags(ctx context.Context, f *Filter) ([]Article, error) {
	log.Println("FindRelatedTags() with filter", f)
	var arts []Article

	opt := options.Find()
	opt.SetProjection(bson.M{"tags": 1})
	opt.SetSort(bson.M{"_id": -1})

	cursor, err := ArticleCollection.Find(ctx, f, opt)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	if err = cursor.All(ctx, &arts); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return arts, nil
}
