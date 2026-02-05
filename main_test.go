package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

type createWallet struct {
	Id string `json:"walletId"`
}

type UpdateWallet struct {
	WalletId      string  `json:"walletId"`
	OperationType string  `json:"operationType"`
	Amount        float32 `json:"amount"`
}

func createUser(name string) error {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	createWallet := createWallet{
		Id: name,
	}

	body, err := json.Marshal(createWallet)
	if err != nil {
		return err
	}

	url := "http://localhost:80/api/v1/wallets/wallet/create"

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	resp.Body.Close()

	return nil
}

func TestMain(m *testing.M) {

	log.SetOutput(io.Discard)

	err := godotenv.Load("config.env")

	go func() {
		startServer()
	}()

	time.Sleep(5 * time.Second)

	err = createUser("asd1")

	if err != nil {

	}

	code := m.Run()

	os.Exit(code)

}

func TestManyGetRequest(t *testing.T) {

	requests := 2000
	var wg sync.WaitGroup

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	start := time.Now()

	sem := make(chan struct{}, 100)

	for i := 0; i < requests; i++ {
		sem <- struct{}{}
		wg.Go(func() {
			defer func() { <-sem }()

			resp, err := client.Get("http://localhost" + os.Getenv("SERVER_PORT") + "/api/v1/wallets/asd1")
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("status: %d", resp.StatusCode)
			}
		})
	}

	wg.Wait()
	fmt.Println(t.Name(), "Time:", time.Since(start))
	if time.Since(start) > 20*time.Second {
		t.Fatal("too long")
	} else {
		fmt.Println(t.Name(), "PASS")
	}
}

func TestManyPostRequest(t *testing.T) {

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	requests := 2000
	var wg sync.WaitGroup

	updateWallet := UpdateWallet{
		WalletId:      "asd1",
		OperationType: "DEPOSIT",
		Amount:        1,
	}

	body, err := json.Marshal(updateWallet)
	if err != nil {
		panic(err)
	}

	start := time.Now()

	sem := make(chan struct{}, 100)

	for i := 0; i < requests; i++ {
		sem <- struct{}{}
		wg.Go(func() {
			defer func() { <-sem }()

			resp, err := client.Post("http://localhost"+os.Getenv("SERVER_PORT")+"/api/v1/wallets/wallet",
				"application/json", bytes.NewReader(body))
			if err != nil {
				t.Error(err)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				t.Errorf("status: %d", resp.StatusCode)
			}
		})
	}

	wg.Wait()
	fmt.Println(t.Name(), "Time:", time.Since(start))
	if time.Since(start) > 5*time.Second {
		t.Fatal("too long")
	} else {
		fmt.Println(t.Name(), "PASS")
	}
}
