package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"log"
)

type message struct {
	WalletId      string  `json:"walletId"`
	OperationType string  `json:"operationType"`
	Amount        float32 `json:"amount"`
}

type createWalletmessage struct {
	WalletId string `json:"walletId"`
}

type WalletStorage interface {
	Get(uuid string) (float32, error)
	Check(uuid string) (bool, error)
	ChangeBalance(sum float32, uuid string) error
	CreateWallet(uuid string) error
}

type Server struct {
	storage WalletStorage
}

func newGetBalanceHandler(ds WalletStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {

			log.Println("Get request '", r.URL.Path, "'")

			parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/") //разделяем адрес на части

			if len(parts) != 4 || r.URL.Path != "/api/v1/wallets/"+parts[3] { // проверяем, что запрос имеет вид /api/v1/wallets/{WALLET_UUID}
				log.Print("wrong path: " + r.URL.Path)
				http.NotFound(w, r)
				return
			}

			uuid := parts[3]

			log.Println("uuid:", uuid)

			sum, err := ds.Get(uuid)

			if err != nil {
				log.Println("error get request:", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Println("Operation is done")
			fmt.Fprintln(w, sum)

		} else {
			log.Println("wrong method on path:", r.URL.Path)
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}
}

func newChangeBalanceHandler(ds WalletStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("wrong method on path:", r.URL.Path)
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		} else {

			log.Println("change request'", r.URL.Path, "'")
			if r.URL.Path != "/api/v1/wallets/wallet" { // проверяем, что запрос имеет вид /api/v1/wallets/{WALLET_UUID}
				log.Println("wrong path:", r.URL.Path)
				http.NotFound(w, r)
				return
			}

			var msg message
			err := json.NewDecoder(r.Body).Decode(&msg)
			if err != nil {
				log.Println("wrong json")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if msg.Amount <= 0 {
				log.Println("wrong json among:", msg.Amount)
				http.Error(w, "sum must be more 0", http.StatusBadRequest)
				return
			}

			check, err := ds.Check(msg.WalletId)

			if err != nil {
				log.Println("wrong serwer bihaviour")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if !check {
				log.Println("wrong uuid:", msg.WalletId)
				http.Error(w, "UUID is wrong", http.StatusBadRequest)
				return
			}

			switch msg.OperationType {
			case "DEPOSIT":
				err = ds.ChangeBalance(msg.Amount, msg.WalletId)
			case "WITHDRAW":
				err = ds.ChangeBalance(-msg.Amount, msg.WalletId)
			default:
				log.Println("wrong operation type")
				http.Error(w, err.Error(), http.StatusBadRequest)
			}

			if err != nil {
				log.Println("wrong serwer bihaviour")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				log.Println("Operation complit")
				fmt.Fprintln(w, "Operation complit")
				return
			}

		}
	}
}

func newCreateWalletHandler(ds WalletStorage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			log.Println("wrong method on path:", r.URL.Path)
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		} else {

			log.Println("create wallet start")

			if r.URL.Path != "/api/v1/wallets/wallet/create" { // проверяем, что запрос имеет вид /api/v1/wallets/{WALLET_UUID}
				http.NotFound(w, r)
				return
			}

			var msg createWalletmessage

			err := json.NewDecoder(r.Body).Decode(&msg)

			if err != nil {
				log.Println("wrong json")
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			log.Println("uuid:", msg.WalletId)
			check, err := ds.Check(msg.WalletId)

			if err != nil {
				log.Println("error in check method")
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if check {
				log.Println("UUID is actually exist:", msg.WalletId)
				http.Error(w, "UUID is actually exist", http.StatusBadRequest)
				return
			}

			err = ds.CreateWallet(msg.WalletId)

			if err != nil {
				log.Println("error in create method: ", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return

			} else {
				log.Println("Wallet created")
				fmt.Fprintln(w, "Wallet created")
				return
			}

		}
	}
}

func (server *Server) Start(ds WalletStorage) {

	server.storage = ds

	http.HandleFunc("/api/v1/wallets/", newGetBalanceHandler(server.storage))

	http.HandleFunc("/api/v1/wallets/wallet/create", newCreateWalletHandler(server.storage))

	http.HandleFunc("/api/v1/wallets/wallet", newChangeBalanceHandler(server.storage))

	fmt.Println("Starting server at port 80")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
