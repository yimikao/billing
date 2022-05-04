package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	fluidlog "github.com/fluidcoins/log"
	"github.com/go-redis/redismock/v8"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/yimikao/billing/mocks"
)

func TestRegisterUserHandler(t *testing.T) {
	testCases := []struct {
		name           string
		expectedStatus int
		buildStubs     func(store *mocks.MockUserRepository)
		// buildCacheStubs func(cache *mocks.MockCache)
		requestBody userRegistrationRequest
	}{
		{
			name:           "successful registration",
			expectedStatus: http.StatusOK,
			buildStubs: func(store *mocks.MockUserRepository) {
				store.EXPECT().Create(gomock.Any()).Times(1).Return(nil)
			},
			// buildCacheStubs: func(cache *mocks.MockCache) {

			// },

			requestBody: userRegistrationRequest{
				FirstName:       "firstname",
				LastName:        "lastname",
				Email:           "email@email.com",
				Tag:             "moneyman",
				TransactionCode: "0000",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var b = new(bytes.Buffer)

			require.NoError(t, json.NewEncoder(b).Encode(tc.requestBody))

			req := httptest.NewRequest(http.MethodPost, "/register", b)
			req.Header.Add("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			logger := fluidlog.New(fluidlog.LevelDebug, 4)

			controller := gomock.NewController(t)
			defer controller.Finish()

			userRepo := mocks.NewMockUserRepository(controller)
			tc.buildStubs(userRepo)

			redis, mock := redismock.NewClientMock()

			mc := &mocks.MockCache{
				Inner: redis,
				Mock:  mock,
			}
			// bb, _ := json.Marshal(billing.User{
			// 	FirstName:       "firstname",
			// 	LastName:        "lastname",
			// 	Email:           "email@email.com",
			// 	Tag:             "moneyman",
			// 	TransactionCode: "0000",
			// })
			mc.Mock.Regexp().ExpectSet("userdata", `[a-z]+`, 0).SetErr(errors.New("FAIL"))
			// tc.buildCacheStubs(mc)

			userRegistrationHandler := NewUserRegistrationHandler(userRepo, logger, mc, context.Background())

			userRegistrationHandler.registerUser(recorder, req)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Error(err)
			}
			require.Equal(t, tc.expectedStatus, recorder.Code)
			verifyMatch(t, recorder.Body)
		})
	}
}
