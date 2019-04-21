package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/shinofara/golang-study/rest/app/middleware"
)

// Transaction transaction entity.
type Transaction struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

// TransactionController controller.
type TransactionController struct {
	Base
}

func (t *TransactionController) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uID := r.Context().Value(middleware.UserIDKey).(int)
	//ユーザーに紐づく取引を取得
	rows, err := t.DB.Open().QueryContext(
		ctx,
		"select id, user_id, amount, description from transactions where user_id=?",
		uID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	list := make([]Transaction, 0)
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Description); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		list = append(list, t)
	}

	if err := json.NewEncoder(w).Encode(list); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *TransactionController) Show(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := mux.Vars(r)["id"]
	uID := r.Context().Value(middleware.UserIDKey).(int)
	var transaction Transaction
	if err := t.DB.Open().QueryRowContext(
		ctx,
		"select id, user_id, amount, description from transactions where id=? and user_id=?",
		id,
		uID,
	).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.Amount,
		&transaction.Description,
	); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *TransactionController) Create(w http.ResponseWriter, r *http.Request) {
	//パラメータを取得
	var transaction Transaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Println(err)
			return
		}
	}()

	ctx := r.Context()
	uID := r.Context().Value(middleware.UserIDKey).(int)



	//利用金額の合計を取得
	rows, err := t.DB.Open().QueryContext(
		ctx,
		"select id, user_id, amount, description from transactions where user_id=?",
		uID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Println(err)
		}
	}()

	total := 0
	for rows.Next() {
		var t Transaction
		if err := rows.Scan(&t.ID, &t.UserID, &t.Amount, &t.Description); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		total = total + t.Amount
	}

	//利用金額の合計が、（uId * 1000)を上回っている場合httpエラー(402)
	println("1■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■")
	println("取引を実行した場合の合計：%d" ,total + transaction.Amount)
	println("ユーザーが使ってもいい最大値：%d", uID * 1000)
	if total + transaction.Amount > uID * 1000 {
		//StatusPaymentRequired
		//http.Error(w,err.Error(), http.StatusPaymentRequired)
		http.Error(w, fmt.Sprintf("limit:%d < total:%d", total, uID * 1000), http.StatusPaymentRequired)
		return
	}

	//取引実行（取引レコードの作成）
	result, err := t.DB.Open().ExecContext(
		ctx,
		"insert into transactions (user_id, amount, description) values (?,?,?)",
		uID,
		transaction.Amount,
		transaction.Description,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	transaction.ID = int(id)



	//レスポンス
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(transaction); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (t *TransactionController) Delete(w http.ResponseWriter, r *http.Request) {
	println("■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■")
	ctx := r.Context()
	//id := mux.Vars(r)["id"]
	//uID := r.Context().Value(middleware.UserIDKey).(int)
	//uID := mux.Vars(r)["id"]
	result, err := t.DB.Open().ExecContext(
		ctx,
		//"delete from transactions where id=? and user_id=?",
		"delete from transactions",
		//id,
		//uID,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	affected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if affected == 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
}
