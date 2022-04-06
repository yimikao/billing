package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	fluidlog "github.com/fluidcoins/log"
	"github.com/golang/mock/gomock"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
	"github.com/yimikao/billing"
	"github.com/yimikao/billing/mocks"
)

func verifyMatch(t *testing.T, v io.Reader) {
	g := goldie.New(t, goldie.WithFixtureDir("./testdata"))

	b := new(bytes.Buffer)

	_, err := io.Copy(b, v)

	require.NoError(t, err)
	g.Assert(t, t.Name(), b.Bytes())
}

func TestUserHandler_Create(t *testing.T) {

	tt := []struct {
		name           string
		expectedStatus int
		buildStubs     func(store *mocks.MockUserRepository)
		requestBody    userCreateRequest
	}{
		{
			name:           "validation failed",
			expectedStatus: http.StatusBadRequest,
			buildStubs: func(store *mocks.MockUserRepository) {
				store.EXPECT().Create(gomock.Any()).Times(0)
			},
			requestBody: userCreateRequest{},
		},
		{
			name:           "internal error",
			expectedStatus: http.StatusInternalServerError,
			buildStubs: func(store *mocks.MockUserRepository) {
				store.EXPECT().Create(gomock.Any()).Times(1).Return(sql.ErrConnDone)
			},
			requestBody: userCreateRequest{
				Email: billing.Email("a@gmail.com"),
			},
		},
		{
			name:           "successfully created a user",
			expectedStatus: http.StatusOK,
			buildStubs: func(store *mocks.MockUserRepository) {
				store.EXPECT().Create(gomock.Any()).Times(1)
			},
			requestBody: userCreateRequest{
				Email: billing.Email("ayo@gmail.com"),
			},
		},
	}

	for _, v := range tt {

		t.Run(v.name, func(t *testing.T) {

			var b = new(bytes.Buffer)

			require.NoError(t, json.NewEncoder(b).Encode(v.requestBody))

			req := httptest.NewRequest(http.MethodPost, "/oops", b)
			req.Header.Add("Content-Type", "application/json")

			recorder := httptest.NewRecorder()

			logger := fluidlog.New(fluidlog.LevelDebug, 4)

			controller := gomock.NewController(t)
			defer controller.Finish()

			userRepo := mocks.NewMockUserRepository(controller)
			v.buildStubs(userRepo)

			userHandler := NewUserHandler(logger, userRepo)

			userHandler.create(recorder, req)

			require.Equal(t, v.expectedStatus, recorder.Code)
			verifyMatch(t, recorder.Body)
		})
	}
}
