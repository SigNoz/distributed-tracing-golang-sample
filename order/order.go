package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/NamanJain8/distributed-tracing-golang-sample/datastore"
	"github.com/NamanJain8/distributed-tracing-golang-sample/utils"
)

type orderData struct {
	ID          int64  `json:"id"`
	UserID      int    `json:"user_id" validate:"required"`
	ProductName string `json:"product_name" validate:"required"`
	Price       int    `json:"price" validate:"required"`
}

type user struct {
	ID       int64  `json:"id"`
	UserName string `json:"user_name"`
	Account  string `json:"account"`
	Amount   int
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var request orderData
	if err := utils.ReadBody(w, r, &request); err != nil {
		return
	}

	// get user details from user service
	url := fmt.Sprintf("http://localhost:%d/users/%d", userServicePort(), request.UserID)
	userResponse, err := utils.SendRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Printf("%v", err)
		utils.WriteResponse(w, http.StatusInternalServerError, err)
		return
	}

	b, err := ioutil.ReadAll(userResponse.Body)
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}
	defer userResponse.Body.Close()

	if userResponse.StatusCode != http.StatusOK {
		utils.WriteErrorResponse(w, userResponse.StatusCode, fmt.Errorf("payment failed. got response: %s", b))
		return
	}

	var user user
	if err := json.Unmarshal(b, &user); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	// basic check for the user balance
	if user.Amount < request.Price {
		utils.WriteErrorResponse(w, http.StatusUnprocessableEntity, fmt.Errorf("insufficient balance. add %d more amount to account", request.Price-user.Amount))
		return
	}

	// insert the order into order table
	id, err := db.InsertOne(datastore.InsertParams{
		Query: `insert into ORDERS(ACCOUNT, PRODUCT_NAME, PRICE, ORDER_STATUS) VALUES (?,?,?, ?)`,
		Vars:  []interface{}{user.Account, request.ProductName, request.Price, "SUCCESS"},
	})
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	// update the pending amount in user table
	if err := db.UpdateOne(datastore.UpdateParams{
		Query: `update USERS set AMOUNT = AMOUNT - ? where ID = ?`,
		Vars:  []interface{}{request.Price, user.ID},
	}); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, err)
		return
	}

	// send response
	response := request
	response.ID = id
	utils.WriteResponse(w, http.StatusCreated, response)
}
