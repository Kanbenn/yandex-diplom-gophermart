package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Kanbenn/gophermart/internal/app"
	"github.com/Kanbenn/gophermart/internal/config"
	"github.com/Kanbenn/gophermart/internal/mocks"
	"github.com/Kanbenn/gophermart/internal/models"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockThisStorerResult func(s *mocks.Mockstorer)

func TestRegisterNewUser(t *testing.T) {
	type mockResult struct {
		user int
		err  error
	}
	tests := []struct {
		name               string
		input              string
		res                mockResult
		expectedStatusCode int
	}{
		{
			name:  "OK",
			input: `{"login":"Lilu Dallas", "password":"Multipass"}`,
			res: mockResult{
				user: 1,
				err:  nil},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:  "Err: Login is already taken",
			input: `{"login":"Lilu Dallas", "password":"Multipass"}`,
			res: mockResult{
				user: 0,
				err:  models.ErrLoginNotUnique},
			expectedStatusCode: http.StatusConflict,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := buildRequestWithBodyFromInput(test.input)

			mockingAction := func(mockStorer *mocks.Mockstorer) {
				mockStorer.EXPECT().InsertNewUser(gomock.Any()).Return(test.res.user, test.res.err)
			}
			h, ctrl := buildHandler(t, mockingAction)
			defer ctrl.Finish()

			h.RegisterNewUser(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

func TestLoginUser(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		mockReturn         models.User
		expectedStatusCode int
	}{
		{
			name:  "OK",
			input: `{"login":"Lilu Dallas", "password":"Multipass"}`,
			mockReturn: models.User{
				ID:       1,
				Login:    "Lilu Dallas",
				Password: "Multipass"},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Err: User Unknown",
			input:              `{"login":"Corban Dallas", "password":"Fake Pass"}`,
			mockReturn:         models.User{},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := buildRequestWithBodyFromInput(test.input)

			mockingAction := func(mockStorer *mocks.Mockstorer) {
				mockStorer.EXPECT().SelectUserAuth(gomock.Any()).Return(test.mockReturn, nil)
			}
			h, ctrl := buildHandler(t, mockingAction)
			defer ctrl.Finish()

			h.LoginUser(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

func TestGetUserOrders(t *testing.T) {
	tests := []struct {
		name               string
		user               int
		mockReturn         []models.UserOrder
		expectedStatusCode int
	}{
		{
			name: "OK",
			user: 1,
			mockReturn: []models.UserOrder{{
				Number: "56167261530448",
				Status: "PROCESSED",
				Bonus:  729.99,
				Time:   "2023-11-01T19:13:09+07:00"},
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Err: NoContent",
			user:               4,
			mockReturn:         []models.UserOrder{},
			expectedStatusCode: http.StatusNoContent,
		}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			req := &http.Request{}
			ctx := context.WithValue(req.Context(), models.CtxKeyUser, test.user)
			req = req.WithContext(ctx)

			mockingAction := func(mockStorer *mocks.Mockstorer) {
				mockStorer.EXPECT().SelectUserOrders(test.user).Return(test.mockReturn, nil)
			}

			h, ctrl := buildHandler(t, mockingAction)
			defer ctrl.Finish()

			h.GetUserOrders(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)

			if len(test.mockReturn) > 0 {
				resp := w.Result()
				defer resp.Body.Close()
				result, _ := io.ReadAll(resp.Body)
				expectedJsn, _ := json.Marshal(test.mockReturn)

				assert.Equal(t, expectedJsn, result)
			}
		})
	}
}

func TestGetUserBalance(t *testing.T) {
	tests := []struct {
		name string
		user int
		ub   models.UserBalance
	}{
		{
			name: "OK wealthy balance",
			user: 2,
			ub: models.UserBalance{
				ID:        2,
				Balance:   360000.0,
				Withdrawn: 720.0,
			},
		},
		{
			name: "OK zero balance",
			user: 4,
			ub: models.UserBalance{
				ID:        4,
				Balance:   0,
				Withdrawn: 0,
			},
		}}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := &http.Request{}
			ctx := context.WithValue(req.Context(), models.CtxKeyUser, test.user)
			req = req.WithContext(ctx)

			mockingAction := func(mockStorer *mocks.Mockstorer) {
				mockStorer.EXPECT().SelectUserBalance(test.user).Return(test.ub, nil)
			}
			h, ctrl := buildHandler(t, mockingAction)
			defer ctrl.Finish()

			h.GetUserBalance(w, req)

			assert.Equal(t, http.StatusOK, w.Code)

			resp := w.Result()
			defer resp.Body.Close()
			result, _ := io.ReadAll(resp.Body)
			expectedJsn, _ := json.Marshal(test.ub)

			assert.Equal(t, expectedJsn, result)
		})
	}
}

func TestGetUserWithdrawHistory(t *testing.T) {
	tests := []struct {
		name               string
		user               int
		mockReturn         []models.Order
		expectedStatusCode int
	}{
		{
			name:               "Err: NoContent",
			user:               4,
			mockReturn:         []models.Order{},
			expectedStatusCode: http.StatusNoContent,
		},
		{
			name: "OK",
			user: 2,
			mockReturn: []models.Order{
				{
					Number: "56167261530448",
					Sum:    729.99,
					Time:   "2023-11-01T19:13:09+07:00",
				},
				{
					Number: "059",
					Sum:    927.99,
					Time:   "2023-11-01T19:13:09+07:20",
				},
			},
			expectedStatusCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			w := httptest.NewRecorder()

			req := &http.Request{}
			ctx := context.WithValue(req.Context(), models.CtxKeyUser, test.user)
			req = req.WithContext(ctx)

			mockingAction := func(mockStorer *mocks.Mockstorer) {
				mockStorer.EXPECT().SelectUserWithdrawHistory(test.user).Return(test.mockReturn, nil)
			}

			h, ctrl := buildHandler(t, mockingAction)
			defer ctrl.Finish()

			h.GetUserWithdrawHistory(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)

			if len(test.mockReturn) > 0 {
				resp := w.Result()
				defer resp.Body.Close()
				result, _ := io.ReadAll(resp.Body)
				expectedJsn, _ := json.Marshal(test.mockReturn)

				assert.Equal(t, expectedJsn, result)
			}
		})
	}
}

func TestPostNewOrderWithBonus(t *testing.T) {
	tests := []struct {
		name               string
		user               int
		input              string
		mockReturn         error
		expectedStatusCode int
	}{
		{
			name:               "OK",
			user:               2,
			input:              `{"order": "2377225624", "sum": 720.99}`,
			mockReturn:         nil,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Err Conflict",
			user:               2,
			input:              `{"order": "2377225624", "sum": 444.44}`,
			mockReturn:         models.ErrOrderWasPostedByThisUser,
			expectedStatusCode: http.StatusConflict,
		},
		{
			name:               "Err Not Enough Minerals",
			user:               3,
			input:              `{"order": "2377225624", "sum": 9999.99}`,
			mockReturn:         models.ErrNotEnoughMinerals,
			expectedStatusCode: http.StatusPaymentRequired,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := buildRequestWithBodyFromInput(test.input)
			ctx := context.WithValue(req.Context(), models.CtxKeyUser, test.user)
			req = req.WithContext(ctx)

			mockingAction := func(mockStorer *mocks.Mockstorer) {
				mockStorer.EXPECT().InsertOrderWithdrawal(gomock.Any()).Return(test.mockReturn)
			}

			h, ctrl := buildHandler(t, mockingAction)
			defer ctrl.Finish()

			h.PostNewOrderWithBonus(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

func TestPostNewOrder(t *testing.T) {
	tests := []struct {
		name               string
		user               int
		input              string
		mockReturn         error
		expectedStatusCode int
	}{
		{
			name:               "Accepted",
			user:               2,
			input:              "2377225624",
			mockReturn:         nil,
			expectedStatusCode: http.StatusAccepted,
		},
		{
			name:               "Err Conflict",
			user:               4,
			input:              "2377225624",
			mockReturn:         models.ErrOrderWasPostedByAnotherUser,
			expectedStatusCode: http.StatusConflict,
		},
		{
			name:               "OK",
			user:               2,
			input:              "2377225624",
			mockReturn:         models.ErrOrderWasPostedByThisUser,
			expectedStatusCode: http.StatusOK,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := buildRequestWithBodyFromInput(test.input)
			req.Header.Set("Content-Type", "text/plain")

			ctx := context.WithValue(req.Context(), models.CtxKeyUser, test.user)
			req = req.WithContext(ctx)

			mockingAction := func(mockStorer *mocks.Mockstorer) {
				mockStorer.EXPECT().InsertOrder(gomock.All()).Return(test.mockReturn)
			}

			h, ctrl := buildHandler(t, mockingAction)
			defer ctrl.Finish()

			h.PostNewOrder(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)
		})
	}
}

func buildRequestWithBodyFromInput(in string) *http.Request {
	// req := &http.Request{}
	// req.Body = io.NopCloser(bytes.NewBufferString(in))
	// return req
	return httptest.NewRequest("POST", "/", bytes.NewBufferString(in))
}

func buildHandler(t gomock.TestReporter, mockThis mockThisStorerResult) (*Handler, *gomock.Controller) {
	cfg := config.New()
	ctrl := gomock.NewController(t)
	mockStorer := mocks.NewMockstorer(ctrl)
	mockThis(mockStorer)
	app := app.New(cfg, mockStorer)
	h := New(cfg, app)
	return h, ctrl
}
