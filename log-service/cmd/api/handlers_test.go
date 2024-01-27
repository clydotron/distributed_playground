package main_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"

	. "log-service/cmd/api"
	"log-service/dao"
	mocks "log-service/dao/mocks"
	"log-service/models"
)

var _ = Describe("Handlers", func() {

	var (
		logger      App
		mockCtl     *gomock.Controller
		mockLogRepo *mocks.MockLogRepo
	)

	BeforeEach(func() {

		mockCtl = gomock.NewController(GinkgoT())
		mockLogRepo = mocks.NewMockLogRepo(mockCtl)

		dataStore := dao.DataStore{
			LogRepo: mockLogRepo,
		}
		logger = App{DataStore: dataStore}
	})

	// is this a describe?
	Context("#WriteLog", func() {
		var (
			res  *httptest.ResponseRecorder
			body bytes.Buffer
		)

		JustBeforeEach(func() {
			res = httptest.NewRecorder()

			req := httptest.NewRequest("POST", "/", &body)
			handler := http.HandlerFunc(logger.WriteLog)
			handler.ServeHTTP(res, req)
		})

		Context("valid body", func() {
			BeforeEach(func() {

				// @TODO check the values:
				expectedLog := models.LogEntry{Name: "name", Data: "data"}
				mockLogRepo.EXPECT().Insert(expectedLog).Return(nil)
				type jsonPayload struct {
					Name string `json:"name"`
					Data string `json:"data"`
				}
				payload := jsonPayload{
					Name: "name",
					Data: "data",
				}

				json.NewEncoder(&body).Encode(payload)

			})

			It("succeeds", func() {
				Expect(res.Result().StatusCode).To(Equal(http.StatusAccepted))
			})
		})
	})
})
