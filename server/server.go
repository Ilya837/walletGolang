package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type message struct {
	walletId      string  `json:"walletId"`
	operationType string  `json:"operationType"`
	amount        float32 `json:"amount"`
}

type DataStorage interface {
	Get(uuid string) (int64, error)
	Check(uuid string) bool
	ChangeBalance(sum float32, uuid string) error
}

type Server struct {
	storage DataStorage
}

func createGetBalanceHandler(ds DataStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/") //разделяем адрес на части

			if len(parts) != 4 || r.URL.Path != "//api/v1/wallets/"+parts[3]+"/" { // проверяем, что запрос имеет вид /api/v1/wallets/{WALLET_UUID}
				http.NotFound(w, r)
				return
			}

			uuid := parts[1]

			sum, err := ds.Get(uuid)

			if err != nil {
				fmt.Fprintln(w, err)
			}

			fmt.Fprintln(w, sum)

		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func createChangeBalanceHandler(ds DataStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		} else {

			if r.URL.Path != "//api/v1/wallets/wallet/" { // проверяем, что запрос имеет вид /api/v1/wallets/{WALLET_UUID}
				http.NotFound(w, r)
				return
			}

			var msg message
			err := json.NewDecoder(r.Body).Decode(&msg)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if msg.amount <= 0 {
				fmt.Fprintln(w, "sum must be more 0")
				return
			}

			if !ds.Check(msg.walletId) {
				fmt.Fprintln(w, "UUID is wrong")
				return
			}

			switch msg.operationType {
			case "DEPOSIT":
				ds.ChangeBalance(msg.amount, msg.walletId)
			case "WITHDRAW":
				ds.ChangeBalance(-msg.amount, msg.walletId)
			}

		}
	}
}

func (server Server) Start() {

	http.HandleFunc("/api/v1/wallets/", createGetBalanceHandler(server.storage))

	http.HandleFunc("/api/v1/wallets/wallet", createChangeBalanceHandler(server.storage))

	fmt.Println("Starting server at port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
