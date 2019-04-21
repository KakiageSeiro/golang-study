package main_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	baseURL = "http://localhost:8888"
	userNum = 3
)

type Transaction struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Amount      int    `json:"amount"`
	Description string `json:"description"`
}

func TestCreate(t *testing.T) {
	//ユーザーごとの値段上限
	limitMap := make(map[int]int, userNum)
	for uID := 1; uID <= userNum; uID++ {
		limitMap[uID] = uID * 1000
	}

	////テスト実行前にユーザーの取引を削除
	////buffer := bytes.NewBuffer(make([]byte, 0, 128))
	//for i := 0; i < 3; i++{
	//	//削除リクエスト作成
	//	req2, err := http.NewRequest(
	//		http.MethodDelete,
	//		baseURL+"/transactions/delete" + "/?id=" + strconv.Itoa(i + 1),
	//		nil,
	//	)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//
	//	//パラメータ作成
	//	values := url.Values{} //url.Valuesオブジェクト生成
	//	values.Add("id", strconv.Itoa(i + 1)) //key-valueを追加
	//	//リクエストに付加
	//	req2.URL.RawQuery = values.Encode()
	//
	//	//結果確認
	//	resp2, err := http.DefaultClient.Do(req2)
	//	if err != nil {
	//		t.Fatal(err)
	//	}
	//	fmt.Println("削除処理" + strconv.Itoa(resp2.StatusCode))
	//}

	// Create transactions
	var wg sync.WaitGroup
	wg.Add(userNum)

	for uID := 1; uID <= userNum; uID++ {
		go func(uID int) {
			defer wg.Done()
			var total int
			for j := 0; j <20; j++ {
				buffer := bytes.NewBuffer(make([]byte, 0, 128))
				amount := uID * 100
				total += amount
				if err := json.NewEncoder(buffer).Encode(Transaction{
					UserID:      uID,
					Amount:      amount,
					Description: fmt.Sprintf("商品%d", uID),
				}); err != nil {
					t.Fatal(err)
				}
				req, err := http.NewRequest(
					http.MethodPost,
					baseURL+"/transactions",
					buffer,
				)
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("apikey", fmt.Sprintf("secure-api-key-%d", uID))

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Fatal(err)
				}

				limit := limitMap[uID]
				if total > limit {
					want := http.StatusPaymentRequired
					if resp.StatusCode != want {
						t.Errorf("POST /transactions status %d != %d total:%d limit:%d", resp.StatusCode, want, total, limit)
					}
				} else {
					want := http.StatusCreated
					if resp.StatusCode != want {
						t.Errorf("POST /transactions status %d != %d total:%d limit:%d", resp.StatusCode, want, total, limit)
					}
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					t.Fatal(err)
				}
				t.Log(string(body))

				if err := resp.Body.Close(); err != nil {
					t.Fatal(err)
				}
			}
		}(uID)
	}
	wg.Wait()

	// Check limit
	conn, err := sql.Open("mysql", "root@tcp(127.0.0.1:43306)/codetest")
	if err != nil {
		t.Fatal(err)
	}
	for uID := 1; uID <= userNum; uID++ {
		var amount int
		if err := conn.QueryRow(
			"select sum(amount) from transactions where user_id=?",
			uID,
		).Scan(&amount); err != nil {
			t.Fatal(err)
		}
		limit := limitMap[uID]
		//値段制限を超えてしまった場合テストエラー
		if amount > limit {
			t.Errorf("User %d amount %d over the limit %d", uID, amount, limit)
		}
	}

	//func () recordDelete(uId string){
	//
	//	r := mux.NewRouter()
	//
	//	db, err := infrastructure.NewDB()
	//	if err != nil {
	//		log.Panic(err)
	//	}
	//
	//	r.Use(middleware.Authenticate(db))
	//	r.Use(middleware.RequestLogger)
	//
	//	base := controller.Base{
	//		DB: db,
	//	}
	//
	//	//トランザクション作成
	//	transaction := controller.TransactionController{Base: base}
	//	ctx := r.Context()
	//	//id := mux.Vars(r)["id"]
	//	//uID := r.Context().Value(middleware.UserIDKey).(int)
	//	result, err := transaction.DB.Open().ExecContext(
	//		ctx,
	//		//"delete from transactions where id=? and user_id=?",
	//		"delete from transactions where user_id=?",
	//		uID,
	//	)
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//
	//	affected, err := result.RowsAffected()
	//	if err != nil {
	//		http.Error(w, err.Error(), http.StatusInternalServerError)
	//		return
	//	}
	//	if affected == 0 {
	//		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	//		return
	//	}
	//
	//
	//
	//}
}
