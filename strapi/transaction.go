package strapi

import (
	"encoding/json"
)

func CheckoutSessionCreate(transactionId string, itemId string, quantity string) error {
	return nil
}
func GetTransactionById(transactionId string) (*Transaction, error) {
	// Get the price id.
	res, err := requestToStrapi(GET, "/transactions?filters[transaction_id]="+transactionId+"&populate[0]=item", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var transaction Transaction
	if err := json.NewDecoder(res.Body).Decode(&transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}
func GetTransactionByPaymentIntent(paymentIntent string) (*Transaction, error) {
	// Get the price id.
	res, err := requestToStrapi(GET, "/transactions?filters[payment_intent]="+paymentIntent+"&populate[0]=item", nil)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var transaction Transaction
	if err := json.NewDecoder(res.Body).Decode(&transaction); err != nil {
		return nil, err
	}

	return &transaction, nil
}
func RegisterTransferId(transactionId string, transferId string) error {
	putRest, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"transfer_id": transferId,
		},
	})
	res, err := requestToStrapi(PUT, "/transactions/"+transactionId, putRest)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

type status string

const (
	Pending   status = "pending"
	Cancelled status = "cancelled"
	Shipped   status = "shipped"
	Completed status = "completed"
	Refunded  status = "refunded"
)

func TransactionRegister(transactionId string, itemId int, quantity int, status status) error {
	putRest, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"transaction_id": transactionId,
			"item":           itemId,
			"quantity":       quantity,
			"status":         status,
		},
	})
	res, err := requestToStrapi(POST, "/transactions", putRest)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}
func CheckoutSessionDetailRegister(transactionId string, state string, city string, line1 string, line2 string, postalCode string, email string, name string, phone string, paymentIntent string) error {
	putRest, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"status":         "paid",
			"state":          state,
			"city":           city,
			"line1":          line1,
			"line2":          line2,
			"postal_code":    postalCode,
			"email":          email,
			"name":           name,
			"phone":          phone,
			"payment_intent": paymentIntent,
		},
	})
	res, err := requestToStrapi(PUT, "/transactions/"+transactionId, putRest)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
func CheckoutSessionTransferRegister(transactionId string, transferId string) error {
	putRest, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"transfer_id": transferId,
		},
	})
	res, err := requestToStrapi(PUT, "/transactions/"+transactionId, putRest)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
func CheckoutSessionStatusRegister(transactionId string, status status) error {
	putRest, _ := json.Marshal(map[string]interface{}{
		"data": map[string]interface{}{
			"status": status,
		},
	})
	res, err := requestToStrapi(PUT, "/transactions/"+transactionId, putRest)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
