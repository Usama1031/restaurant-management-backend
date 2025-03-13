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
)

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "orderItem")

func GetOrderItems() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		res, err := orderItemCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing orders items"})
			return
		}

		var allOrderItems []bson.M

		if err = res.All(ctx, &allOrderItems); err != nil {
			log.Fatal(err)
			return
		}

		c.JSON(http.StatusOK, allOrderItems)

	}
}

func GetOrderItem() gin.HandlerFunc {

	return func(c *gin.Context) {

	}
}

func GetOrderItemsByOrder() gin.HandlerFunc {

	return func(c *gin.Context) {

	}
}

func ItemsByOrder(id string) (OrderITems []primitive.M, err error) {

}

func CreateOrderItem() gin.HandlerFunc {

	return func(c *gin.Context) {

	}
}

func UpdateOrderItem() gin.HandlerFunc {

	return func(c *gin.Context) {

	}
}
