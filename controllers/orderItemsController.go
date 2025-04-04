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

var orderCollection *mongo.Collection = database.OpenCollection(database.Client, "orders")

type OrderItemPack struct {
	Table_id    *string
	Order_items []models.OrderItem
}

var orderItemCollection *mongo.Collection = database.OpenCollection(database.Client, "order_items")

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

func GetOrderItemsByOrder() gin.HandlerFunc {

	return func(c *gin.Context) {

		orderId := c.Param("order_id")

		allOrderItems, err := ItemsByOrder(orderId)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occurered while listing order items by order ID."})
			return
		}

		c.JSON(http.StatusOK, allOrderItems)
	}
}

func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{
		{Key: "order_id", Value: id},
	}}}

	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "food"},
		{Key: "localField", Value: "food_id"},
		{Key: "foreignField", Value: "food_id"},
		{Key: "as", Value: "food"},
	}}}

	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$food"}, {Key: "preserveNullAndEmptyArrays", Value: true}}}}

	lookupOrderStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "orders"},
		{Key: "localField", Value: "order_id"},
		{Key: "foreignField", Value: "order_id"},
		{Key: "as", Value: "order"},
	}}}

	unwindOrderStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$order"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}

	lookupTableStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "table"},
		{Key: "localField", Value: "order.table_id"},
		{Key: "foreignField", Value: "table_id"},
		{Key: "as", Value: "table"},
	}}}

	unwindTableStage := bson.D{{Key: "$unwind", Value: bson.D{
		{Key: "path", Value: "$table"},
		{Key: "preserveNullAndEmptyArrays", Value: true},
	}}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{
		{Key: "amount", Value: "$food.price"},
		{Key: "total_count", Value: 1},
		{Key: "food_name", Value: "$food.name"},
		{Key: "food_image", Value: "$food.food_image"},
		{Key: "table_number", Value: "$table.table_number"},
		{Key: "table_id", Value: "$table.table_id"},
		{Key: "order_id", Value: "$order.order_id"},
		{Key: "price", Value: "$food.price"},
		{Key: "quantity", Value: 1},
	}}}

	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{
			{Key: "order_id", Value: "$order_id"},
		}},
		{Key: "table_number", Value: bson.D{
			{Key: "$first", Value: "$table_number"},
		}},
		{Key: "table_id", Value: bson.D{
			{Key: "$first", Value: "$table_id"},
		}},
		{Key: "payment_due", Value: bson.D{
			{Key: "$sum", Value: "$amount"},
		}},
		{Key: "total_count", Value: bson.D{
			{Key: "$sum", Value: 1},
		}},
		{Key: "order_items", Value: bson.D{
			{Key: "$push", Value: "$$ROOT"},
		}},
	}}}

	projectStage2 := bson.D{{Key: "$project", Value: bson.D{
		{Key: "payment_due", Value: 1},
		{Key: "total_count", Value: 1},
		{Key: "table_number", Value: 1},
		{Key: "table_id", Value: 1},
		{Key: "order_items", Value: 1},
	}}}

	res, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2,
	})

	if err != nil {
		panic(err)
	}

	if res.All(ctx, &OrderItems); err != nil {
		panic(err)
	}

	return OrderItems, err

}

func GetOrderItem() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderItemId := c.Param("orderItem_id")

		var orderItem models.OrderItem

		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error occured while listing the order item."})
		}

		c.JSON(http.StatusOK, orderItem)
	}
}

func CreateOrderItem() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItemPack OrderItemPack
		var order models.Order

		if err := c.BindJSON(&orderItemPack); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		order.Table_id = orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)

		orderItemsToBeInserted := []interface{}{}

		for _, orderItem := range orderItemPack.Order_items {
			orderItem.Order_id = order_id

			validationErr := validate.Struct(orderItem)

			if validationErr != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}

			orderItem.ID = primitive.NewObjectID()
			orderItem.Order_item_id = orderItem.ID.Hex()

			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

			var num = toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num

			orderItemsToBeInserted = append(orderItemsToBeInserted, orderItem)

		}

		insertedOrderItems, insertErr := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)

		if insertErr != nil {
			log.Fatal(insertErr)
		}

		c.JSON(http.StatusOK, insertedOrderItems)
	}
}

func UpdateOrderItem() gin.HandlerFunc {

	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var orderItem models.OrderItem

		orderItemId := c.Param("orderItem_id")

		filter := bson.M{"orderItem_id": orderItemId}

		var updatedObj primitive.D

		if orderItem.Unit_price != nil {
			updatedObj = append(updatedObj, bson.E{Key: "unit_price", Value: orderItem.Unit_price})
		}

		if orderItem.Quantity != nil {
			updatedObj = append(updatedObj, bson.E{Key: "quantity", Value: orderItem.Quantity})
		}

		if orderItem.Food_id != nil {
			updatedObj = append(updatedObj, bson.E{Key: "food_id", Value: orderItem.Food_id})
		}

		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

		updatedObj = append(updatedObj, bson.E{Key: "updated_at", Value: orderItem.Updated_at})

		upsert := true

		opt := options.UpdateOptions{
			Upsert: &upsert,
		}

		res, insertErr := orderItemCollection.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: updatedObj}}, &opt)

		if insertErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update the order item!"})
			return
		}

		c.JSON(http.StatusOK, res)
	}
}

func DeleteOrderItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		orderItemId := c.Param("orderItem_id")

		var orderItem models.OrderItem
		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order item not found"})
			return
		}

		filter := bson.M{"order_item_id": orderItemId}
		res, err := orderItemCollection.DeleteOne(ctx, filter)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order item"})
			return
		}

		count, err := orderItemCollection.CountDocuments(ctx, bson.M{"order_id": orderItem.Order_id})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check remaining order items"})
			return
		}

		if count == 0 {
			orderFilter := bson.M{"order_id": orderItem.Order_id}
			_, err := orderCollection.DeleteOne(ctx, orderFilter)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete order"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"message": "Order and all items deleted"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Order item deleted", "deleted_count": res.DeletedCount})
	}
}

func OrderItemOrderCreator(order models.Order) string {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	order.ID = primitive.NewObjectID()
	order.Order_id = order.ID.Hex()

	_, err := orderCollection.InsertOne(ctx, order)

	if err != nil {
		log.Fatal(err)
	}

	return order.Order_id

}
