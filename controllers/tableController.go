package controllers

import (
	"context"
	"golang-restaurant-management/database"
	"golang-restaurant-management/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var tableCollection *mongo.Collection = database.OpenCollection(database.Client, "table")

func GetTables() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		res, err := tableCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list all the tables!"})
			return
		}

		var allTables []bson.M

		if err = res.All(ctx, &allTables); err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, allTables)
	}
}

func GetTable() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tableId := c.Param("table_id")

		var table models.Table

		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load the table!"})
			return
		}

		c.JSON(http.StatusOK, table)
	}
}

func CreateTable() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(table)

		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		table.ID = primitive.NewObjectID()
		table.Table_id = table.ID.Hex()

		res, insertErr := tableCollection.InsertOne(ctx, table)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Table item was not created!"})
			return
		}

		c.JSON(http.StatusOK, res)

	}
}

func UpdateTable() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var table models.Table
		tableId := c.Param("table_id")

		if err := c.BindJSON(&table); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var updatedObj primitive.D

		if table.Number_of_guests != nil {
			updatedObj = append(updatedObj, bson.E{Key: "number_of_guests", Value: table.Number_of_guests})
		}

		if table.Table_number != nil {
			updatedObj = append(updatedObj, bson.E{Key: "table_number", Value: table.Table_number})
		}

		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		upsert := true
		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		filter := bson.M{"table_id": tableId}

		res, err := tableCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updatedObj}}, &opt)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Table item update failed!"})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func DeleteTable() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		tableId := c.Param("table_id")

		filter := bson.M{"table_id": tableId}

		res, err := tableCollection.DeleteOne(ctx, filter)

		if err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, res)

	}
}
