package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	fluidlog "github.com/fluidcoins/log"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"github.com/yimikao/billing/mocks"
)

func TestRegisterUserHandler(t *testing.T) {
	testCases := []struct {
		name           string
		expectedStatus int
		buildStubs     func(store *mocks.MockUserRepository)
		requestBody    userRegistrationRequest
	}{
		{
			name:           "successful registration",
			expectedStatus: http.StatusPermanentRedirect,
			buildStubs: func(store *mocks.MockUserRepository) {
				store.EXPECT().Create(gomock.Any()).Times(1).Return(nil)
			},
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

			userRegistrationHandler := NewUserRegistrationHandler(userRepo, logger)

			userRegistrationHandler.registerUser(recorder, req)

			require.Equal(t, tc.expectedStatus, recorder.Code)
			verifyMatch(t, recorder.Body)
		})
	}
}
