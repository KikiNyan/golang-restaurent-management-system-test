package controllers

import (
	"context"
	"golang-restaurent-management/database"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type InvoiceViewFormat struct {
	Invoice_Id       string
	Payment_method   string
	Order_Id         string
	Payment_status   *string
	Payment_due      interface{}
	Table_number     interface{}
	Payment_due_date time.Time
	Order_details    interface{}
}

var invoiceCollection *mongo.Collection = database.OpenCollection(database.Client, "invoice")

// this func will get all invoices by specific ID and display it in web
func GetInvoices() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		result, err := orderCollection.Find(context.TODO(), bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing invoices"})
		}

		var allInvoices []bson.M
		if err = result.All(ctx, &allInvoices); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusOK, allInvoices)
	}
}

// this func will get one invoice only by getinvoice by certain invoice id
func GetInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		// this invoice id will be created everytime when user want to see single invoice
		invoiceId := c.Param("invoice_id")

		//call invoice from invoicemodels
		var invoice models.Invoice

		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		defer cancel()

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurde while fetching invoice"})
		}

		var invoiceView InvoiceViewFormat

		// this code will call all orderitem by order id and display all orderitem inside the invoice
		allOrdersItems, err := ItemsByOrder(invoice.Order_Id)

		// the function collects information about all the checkout items to be displayed,
		// such as checkout ID, payment due date, payment method, payment status, table number, and order details.
		invoiceView.Order_Id = invoice.Order_Id
		invoiceView.Payment_due = invoice.Payment_due_date
		invoiceView.Payment_method = "null"

		if invoice.Payment_method != nil {
			invoiceView.Payment_method = *invoice.Payment_method
		}

		invoiceView.Invoice_Id = invoice.Invoice_Id
		invoiceView.Payment_status = *&invoice.Payment_status
		invoiceView.Payment_due = allOrdersItems[0]["payment_due"]
		invoiceView.Table_number = allOrdersItems[0]["table_number"]
		invoiceView.Order_details = allOrdersItems[0]["order_items"]
		//invoiceView.Order_details = allOrdersItems[0]["order_items"]<it means collect all orderitem that stored in array and display in invoiceview orderdetails
		//  after gathering all the checkout invoice information, it provides the information to display the checkout invoice
		c.JSON(http.StatusOK, invoiceView)
	}
}
func CreateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
func UpdateInvoice() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var invoice models.Invoice
		var order models.Order

		invoiceId := c.Param("invoice_id")

		if err := c.BindJSON(&invoice); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		var updateObj primitive.D

		if invoice.Payment_method != nil {
			updateObj = append(updateObj, bson.E{"payment_method": invoice.Payment_method})
		}
	}
}
