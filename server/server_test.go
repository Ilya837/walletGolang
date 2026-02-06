package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGoodGetMethod(t *testing.T) {
	ds := NewMockWalletStorage(t)

	uuid := "1"

	ds.EXPECT().
		Get(uuid).
		Return(true, 3, nil).
		Once()

	handler := newGetBalanceHandler(ds)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/wallets/"+uuid,
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	body := rec.Body.String()

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.Equal(t, string(body), "3\n")

}

func TestUUIDUndefinedGetMethod(t *testing.T) {
	ds := NewMockWalletStorage(t)

	uuid := "1"

	ds.EXPECT().
		Get(uuid).
		Return(false, 0, nil).
		Once()

	handler := newGetBalanceHandler(ds)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/wallets/"+uuid,
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	body := rec.Body.String()

	assert.Equal(t, http.StatusBadRequest, res.StatusCode)
	assert.Equal(t, string(body), "uuid undefined\n")

}

func TestBDErrorGetMethod(t *testing.T) {
	ds := NewMockWalletStorage(t)

	uuid := "1"
	errorText := "some error text"

	ds.EXPECT().
		Get(uuid).
		Return(false, 0, errors.New(errorText)).
		Once()

	handler := newGetBalanceHandler(ds)

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/wallets/"+uuid,
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	body := rec.Body.String()

	assert.Equal(t, http.StatusInternalServerError, res.StatusCode)
	assert.Equal(t, string(body), errorText+"\n")

}

func TestWrongMethodGetMethod(t *testing.T) {
	ds := NewMockWalletStorage(t)

	uuid := "1"

	handler := newGetBalanceHandler(ds)

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/wallets/"+uuid,
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	body := rec.Body.String()

	assert.Equal(t, http.StatusMethodNotAllowed, res.StatusCode)
	assert.Equal(t, string(body), "Invalid request method\n")

}

func TestWrongLongURLGetMethod(t *testing.T) {
	ds := NewMockWalletStorage(t)

	uuid := "1"

	handler := newGetBalanceHandler(ds)

	url := "/api/v1/wallets/" + uuid + "/get"

	req := httptest.NewRequest(
		http.MethodGet,
		url,
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)

}

func TestWrongTextURLGetMethod(t *testing.T) {
	ds := NewMockWalletStorage(t)

	uuid := "1"

	handler := newGetBalanceHandler(ds)

	url := "/api/v1/wallet/" + uuid

	req := httptest.NewRequest(
		http.MethodGet,
		url,
		nil,
	)

	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	assert.Equal(t, http.StatusNotFound, res.StatusCode)

}
