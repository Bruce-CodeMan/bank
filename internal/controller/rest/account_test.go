package rest_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	db "github.com/BruceCompiler/bank/db/sqlc"
	"github.com/BruceCompiler/bank/internal/mocks"
	"github.com/BruceCompiler/bank/internal/server"
	"github.com/BruceCompiler/bank/utils"
)

func TestGetAccountAPI(t *testing.T) {
	account := randomAccount()

	testCases := []struct {
		name          string
		publicID      string
		buildStubs    func(store *mocks.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			publicID: uuid.UUID(account.PublicID.Bytes).String(),
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().
					GetAccountByUUID(gomock.Any(), gomock.Eq(account.PublicID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:     "NotFound",
			publicID: uuid.UUID(account.PublicID.Bytes).String(),
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().
					GetAccountByUUID(gomock.Any(), gomock.Eq(account.PublicID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:     "InternalError",
			publicID: uuid.UUID(account.PublicID.Bytes).String(),
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().
					GetAccountByUUID(gomock.Any(), gomock.Eq(account.PublicID)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "InvalidPublicID",
			publicID: "0",
			buildStubs: func(store *mocks.MockStore) {
				store.EXPECT().
					GetAccountByUUID(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mocks.NewMockStore(ctrl)

			tc.buildStubs(store)

			// start test server and send request
			srv := server.NewHTTPServer(store)

			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/api/v1/account/%s", tc.publicID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			srv.Router().ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)
		})

	}

}

func randomAccount() db.Account {
	return db.Account{
		ID: utils.RandomInt(1, 1000),
		PublicID: pgtype.UUID{
			Bytes: uuid.New(),
			Valid: true,
		},
		Owner:    utils.RandomOwner(),
		Currency: utils.RandomCurrency(),
		Balance:  utils.RandomMoney(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, expected db.Account) {
	t.Helper()

	var actual db.Account
	err := json.NewDecoder(body).Decode(&actual)

	require.NoError(t, err)
	require.Equal(t, expected, actual)
}
