package main

import (
	"fmt"
	"net/http"

	"github.com/NamanJain8/distributed-tracing-golang-sample/datastore"
	"github.com/NamanJain8/distributed-tracing-golang-sample/utils"
	"github.com/gorilla/mux"
)

type user struct {
	ID       int64  `json:"id" validate:"-"`
	UserName string `json:"user_name" validate:"required"`
	Account  string `json:"account" validate:"required"`
	Amount   int
}

type paymentData struct {
	Amount int `json:"amount" validate:"required"`
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var u user
	if err := utils.ReadBody(w, r, &u); err != nil {
		return
	}

	ctx, span := tracer.Start(r.Context(), "create user")
	defer span.End()
	id, err := db.InsertOne(ctx, datastore.InsertParams{
		Query: `INSERT INTO USERS(USER_NAME, ACCOUNT) VALUES (?, ?)`,
		Vars:  []interface{}{u.UserName, u.Account},
	})
	if err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("create user error: %w", err))
		return
	}

	u.ID = id
	utils.WriteResponse(w, http.StatusCreated, u)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userID"]
	var u user

	ctx, span := tracer.Start(r.Context(), "get user")
	defer span.End()
	if err := db.SelectOne(ctx, datastore.SelectParams{
		Query:   `select ID, USER_NAME, ACCOUNT, AMOUNT from USERS where ID = ?`,
		Filters: []interface{}{userID},
		Result:  []interface{}{&u.ID, &u.UserName, &u.Account, &u.Amount},
	}); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("get user error: %w", err))
		return
	}

	utils.WriteResponse(w, http.StatusOK, u)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userID"]
	var data paymentData
	if err := utils.ReadBody(w, r, &data); err != nil {
		return
	}

	ctx, span := tracer.Start(r.Context(), "update user amount")
	defer span.End()
	if err := db.UpdateOne(ctx, datastore.UpdateParams{
		Query: `update USERS set AMOUNT = AMOUNT + ? where ID = ?`,
		Vars:  []interface{}{data.Amount, userID},
	}); err != nil {
		utils.WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("get user error: %w", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}
