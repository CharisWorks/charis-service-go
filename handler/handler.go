package handler

import (
	"fmt"
	"log"
	"strconv"

	"github.com/charisworks/charisworks-service-go/images"
	"github.com/charisworks/charisworks-service-go/mail"
	"github.com/charisworks/charisworks-service-go/strapi"
	"github.com/charisworks/charisworks-service-go/util"

	_stripe "github.com/charisworks/charisworks-service-go/stripe"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v76"
)

var r2conns = images.R2Conns{}

type Handler struct {
	Router *gin.Engine
}

func NewHandler(router *gin.Engine) *Handler {
	r2conns.Init()
	return &Handler{
		Router: router,
	}
}

func CheckoutSessionCompleteHandler(event stripe.Event) (err error) {
	billing := sanitizeNill(event.Data.Object["customer_details"].(map[string]interface{}))
	address := sanitizeNill(event.Data.Object["customer_details"].(map[string]interface{})["address"].(map[string]interface{}))

	// 構造体をJSONにエンコード
	transaction, err := strapi.GetTransactionById(event.Data.Object["id"].(string))
	fmt.Printf(`
*************************************************
CheckoutSession was completed!
transactionId: %s
****Customer Infomation****
state: %s
city: %s
line1: %s
line2: %s
postal_code: %s
email: %s
name: %s
phone: %s
****Transaction Information****
ItemId: %d
Item Name: %s
Quantity: %d
*************************************************
		`,
		strconv.Itoa(transaction.Data[0].ID),
		address["state"],
		address["city"],
		address["line1"],
		address["line2"],
		address["postal_code"],
		billing["email"],
		billing["name"],
		billing["phone"],
		transaction.Data[0].Attributes.Item.Data.Id,
		transaction.Data[0].Attributes.Item.Data.Attributes.Name,
		transaction.Data[0].Attributes.Quantity,
	)
	if err != nil {
		return err
	}
	if err := strapi.CheckoutSessionDetailRegister(
		strconv.Itoa(transaction.Data[0].ID),
		address["state"],
		address["city"],
		address["line1"],
		address["line2"],
		address["postal_code"],
		billing["email"],
		billing["name"],
		billing["phone"],
		event.Data.Object["payment_intent"].(string),
	); err != nil {
		return err
	}
	if err := strapi.ReducePreStock(transaction.Data[0].Attributes.Item.Data.Id, transaction.Data[0].Attributes.Quantity); err != nil {
		return err
	}
	item, err := strapi.GetItem(transaction.Data[0].Attributes.Item.Data.Id)
	if err != nil {
		return err
	}
	trId, err := _stripe.Transfer(float64(transaction.Data[0].Attributes.Item.Data.Attributes.Price*transaction.Data[0].Attributes.Quantity), item.Data.Attributes.Worker.Data.Attributes.StripeAccountID, transaction.Data[0].Attributes.TransactionID)
	if err != nil {
		return err
	}
	if err := strapi.CheckoutSessionTransferRegister(strconv.Itoa(transaction.Data[0].ID), trId); err != nil {
		return err
	}
	if err := mail.SendEmail(util.ADMIN_EMAIL, "購入通知",
		mail.PurchasedAdminEmailFactory(
			billing["name"],
			billing["email"],
			address["postal_code"],
			address["state"],
			address["city"],
			address["line1"],
			address["line2"],
			transaction.Data[0].Attributes.Item.Data.Attributes.Name,
			transaction.Data[0].Attributes.Item.Data.Attributes.Price,
			transaction.Data[0].Attributes.Quantity,
			transaction.Data[0].Attributes.CreatedAt,
		)); err != nil {
		return err
	}
	if err := mail.SendEmail(billing["email"], "ご購入ありがとうございます",
		mail.PurchasedCustomerEmailFactory(
			billing["name"],
			trId,
			transaction.Data[0].Attributes.Item.Data.Attributes.Name,
			transaction.Data[0].Attributes.Item.Data.Attributes.Price,
			transaction.Data[0].Attributes.Quantity,
			util.SHIPPING_FEE,
			transaction.Data[0].Attributes.CreatedAt,
			address["postal_code"],
			address["state"],
			address["city"],
			address["line1"],
			address["line2"],
			billing["email"],
		)); err != nil {
		return err
	}
	if err := mail.SendEmail(item.Data.Attributes.Worker.Data.Attributes.Email, "出品された商品が購入されました",
		mail.PurchasedWorkerEmailFactory(
			item.Data.Attributes.Worker.Data.Attributes.UserName,
			address["postal_code"],
			address["state"],
			address["city"],
			address["line1"],
			address["line2"],
			transaction.Data[0].Attributes.Item.Data.Attributes.Name,
			transaction.Data[0].Attributes.Item.Data.Attributes.Price,
			transaction.Data[0].Attributes.Quantity,
			int(float64(transaction.Data[0].Attributes.Item.Data.Attributes.Price)*(1-util.MARGIN)),
			transaction.Data[0].Attributes.CreatedAt,
		)); err != nil {
		return err
	}
	return nil
}

func CheckoutSessionExpiredHandler(event stripe.Event) (err error) {
	transaction, err := strapi.GetTransactionById(event.Data.Object["id"].(string))
	if err != nil {
		return err
	}
	if err := strapi.ReturnPreStock(transaction.Data[0].Attributes.Item.Data.Id, transaction.Data[0].Attributes.Quantity); err != nil {
		return err
	}
	if err = strapi.CheckoutSessionStatusRegister(strconv.Itoa(transaction.Data[0].ID), strapi.Cancelled); err != nil {
		return err
	}
	return nil
}

func ChargeRefundedHandler(event stripe.Event) (err error) {
	transaction, err := strapi.GetTransactionByPaymentIntent(event.Data.Object["payment_intent"].(string))
	if err != nil {
		return err
	}
	if err := _stripe.ReverseTransfer(transaction.Data[0].Attributes.TransferID.(string)); err != nil {
		return err
	}
	if err := strapi.CheckoutSessionStatusRegister(strconv.Itoa(transaction.Data[0].ID), strapi.Refunded); err != nil {
		return err
	}
	item, err := strapi.GetItem(transaction.Data[0].Attributes.Item.Data.Id)
	if err != nil {
		return err
	}
	log.Print(transaction.Data[0].Attributes.Email)
	log.Print(item.Data.Attributes.Worker.Data.Attributes.Email)
	log.Print(transaction.Data[0].Attributes.TransactionID)

	if err := mail.SendEmail(item.Data.Attributes.Worker.Data.Attributes.Email, "商品がキャンセルされました",
		mail.RefundedWorkerEmailFactory(
			transaction.Data[0].Attributes.Name,
			transaction.Data[0].Attributes.TransactionID,
			transaction.Data[0].Attributes.Item.Data.Attributes.Name,
			transaction.Data[0].Attributes.Item.Data.Attributes.Price,
			transaction.Data[0].Attributes.Quantity,
			transaction.Data[0].Attributes.CreatedAt,
		),
	); err != nil {
		return err
	}
	if err := mail.SendEmail(transaction.Data[0].Attributes.Email, "返金が完了しました。",
		mail.RefundedCustomerEmailFactory(
			transaction.Data[0].Attributes.Name,
			transaction.Data[0].Attributes.TransactionID,
			transaction.Data[0].Attributes.Item.Data.Attributes.Name,
			transaction.Data[0].Attributes.Item.Data.Attributes.Price,
			transaction.Data[0].Attributes.Quantity,
			transaction.Data[0].Attributes.CreatedAt,
		),
	); err != nil {
		return err
	}
	return nil
}

func ShippingHandler(transaction *strapi.Transaction) (err error) {
	if transaction.Data[0].Attributes.TrackingID == nil {
		return fmt.Errorf("TrackingID is empty")
	}
	if err := mail.SendEmail(transaction.Data[0].Attributes.Email, "商品が発送されました",
		mail.ShippingCustomerEmailFactory(
			transaction.Data[0].Attributes.TransactionID,
			transaction.Data[0].Attributes.TrackingID.(string),
			transaction.Data[0].Attributes.Item.Data.Attributes.Name,
			transaction.Data[0].Attributes.Item.Data.Attributes.Price,
			transaction.Data[0].Attributes.Quantity,
			transaction.Data[0].Attributes.CreatedAt,
			transaction.Data[0].Attributes.Email,
			transaction.Data[0].Attributes.Name,
			transaction.Data[0].Attributes.PostalCode,
			transaction.Data[0].Attributes.State,
			transaction.Data[0].Attributes.City,
			transaction.Data[0].Attributes.Line1,
			transaction.Data[0].Attributes.Line2,
		),
	); err != nil {
		return err
	}
	item, err := strapi.GetItem(transaction.Data[0].Attributes.Item.Data.Id)
	if err != nil {
		return err
	}
	if err := mail.SendEmail(item.Data.Attributes.Worker.Data.Attributes.Email, "商品が発送されました",
		mail.ShippingAdminEmailFactory(
			transaction.Data[0].Attributes.TransactionID,
			transaction.Data[0].Attributes.Name,
			transaction.Data[0].Attributes.Email,
			transaction.Data[0].Attributes.PostalCode,
			transaction.Data[0].Attributes.State,
			transaction.Data[0].Attributes.City,
			transaction.Data[0].Attributes.Line1,
			transaction.Data[0].Attributes.Line2,
			transaction.Data[0].Attributes.Item.Data.Attributes.Name,
			transaction.Data[0].Attributes.Item.Data.Attributes.Price,
			transaction.Data[0].Attributes.Quantity,
			transaction.Data[0].Attributes.CreatedAt,
		),
	); err != nil {
		return err
	}
	return nil
}
