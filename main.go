package main

import (
	"CPEventsBackend/pkgs/codechef"
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var password string = "<Redacted>"
var user string = "<Redacted>"
var url string = "mongodb+srv://" + user + ":" + password + "@<Redacted>/test?retryWrites=true&w=majority"
var results []codechef.CodechefEvent

func updateDatabaseCodechef() {
	ctx := context.TODO()
	clientOpt := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOpt)
	if err != nil {
		log.Fatal(err)
	}
	codechefCollection := client.Database("EventsBase").Collection("codechef-database")
	delAll, err := codechefCollection.DeleteMany(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Documents Deleted from codechef database: %d\n", delAll.DeletedCount)
	docs := make([]interface{}, 0, 2000)
	for _, value := range codechef.GetContestDataFromCodeChef() {
		docs = append(docs, interface{}(value))
	}
	insertAll, err := codechefCollection.InsertMany(ctx, docs)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Documents Inserted to codechef database:", len(insertAll.InsertedIDs))
	client.Disconnect(ctx)
}
func GetCodechefEvents() {
	ctx := context.TODO()
	clientOpt := options.Client().ApplyURI(url)
	client, err := mongo.Connect(ctx, clientOpt)
	if err != nil {
		log.Fatal(err)
	}

	results = make([]codechef.CodechefEvent, 0, 2000)
	codechefCollection := client.Database("EventsBase").Collection("codechef-database")
	cur, err := codechefCollection.Find(ctx, bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(ctx) {
		var e codechef.CodechefEvent
		err = cur.Decode(&e)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, e)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(ctx)
	client.Disconnect(ctx)
}

func updateRoutine() {
	for true {
		GetCodechefEvents()
		time.Sleep(4 * time.Minute)
	}
}
func databaseUpdateRoutine() {
	for true {
		time.Sleep(15 * time.Minute)
		updateDatabaseCodechef()
	}
}
func handleCodechefGET(ctx *gin.Context) {
	if len(results) == 0 {
		ctx.JSON(500, map[string]string{
			"status": "fail",
		})
	} else {
		ctx.JSON(200, map[string]interface{}{
			"status": "success",
			"events": results,
		})
	}
}

func main() {
	go updateRoutine()
	go databaseUpdateRoutine()
	router := gin.New()
	router.Use(cors.Default())
	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/codechef", handleCodechefGET)
		}
	}
	router.Run()
}
