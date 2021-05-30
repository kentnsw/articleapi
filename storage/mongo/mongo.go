package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client
var ArticleCollection *mongo.Collection

const ARTICLE_DB_NAME = "articledb"
const ARTICLE_COLLECTION_NAME = "articles"

func Connect(uri string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// ping test
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connected to", uri)

	db := client.Database(ARTICLE_DB_NAME)
	db.CreateCollection(ctx, ARTICLE_COLLECTION_NAME)
	ArticleCollection = db.Collection(ARTICLE_COLLECTION_NAME)
	createIndex(ctx, ArticleCollection, "tags", false, 1)
	createIndex(ctx, ArticleCollection, "date", false, -1)

	return client
}

func createIndex(ctx context.Context, collection *mongo.Collection, field string, unique bool, order int) {
	mod := mongo.IndexModel{
		Keys:    bson.M{field: order}, // index in ascending order or -1 for descending order
		Options: options.Index().SetUnique(unique),
	}
	_, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		log.Fatal(err)
	}
}

func Disconnect() {
	if err := client.Disconnect(context.Background()); err != nil {
		log.Panic(err)
	}
}
