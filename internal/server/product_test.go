package server

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/lahnasti/go-market/internal/models"
	"github.com/lahnasti/go-market/mocks"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/golang/mock/gomock"

)

//mockery --dir=internal/repository --name=Repository --output=mocks/ --outpkg=mocks

func TestGetAllProductsHandler(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mocks.NewRepository(ctrl)
	srv := &Server{
		Db:    m,
		log:   zerolog.New(os.Stdout),
		Valid: validator.New(),
	}
	r := gin.Default()
	r.GET("/products", srv.GetAllProductsHandler)
	httpSrv := httptest.NewServer(r)
	defer httpSrv.Close()

	type want struct {
		code    int
		products string
	}
	type test struct {
		name    string
		request string
		method  string
		product  []models.Product
		err     error
		want    want
	}
	tests := []test{
		{
			name: "Test 'GetAllProductsHandler' #1; Default call",
			request: "/products",
			method: http.MethodGet,
			product: []models.Product {
				{UID: 1, Name: "apple", Description: "fruit", Price: 100, Quantity: 10},
				{UID: 2, Name: "banana", Description: "fruit", Price: 50, Quantity: 5},
			},
			want: want {
				code:    http.StatusOK,
                products: `{"List of product","products"[{"uid":1,"name":"apple","description":"fruit","price":100,"quantity":10,delete:"false"},{"uid":2,"name":"banana","description":"fruit","price":50,"quantity":5,"delete":"false"}]}`,
			},
		},
		{
			name:    "Test 'GetAllProductsHandler' #2; Error call",
            request: "/products",
            method:  http.MethodGet,
            err:     errors.New("db error"),
            want:    want{
				code: http.StatusInternalServerError,
			},
		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err != nil {
				// Устанавливаем поведение мока для ошибки
				m.EXPECT().GetAllProducts().Return(nil, tt.err)
			} else {
				// Устанавливаем поведение мока для успешного вызова
				m.EXPECT().GetAllProducts().Return(tt.product, nil)
			}

			req, _ := http.NewRequest(tt.method, tt.request, nil)
			resp := httptest.NewRecorder()

			r.ServeHTTP(resp, req)

			assert.Equal(t, tt.want.code, resp.Code)
			if tt.want.products != "" {
				assert.JSONEq(t, tt.want.products, resp.Body.String())
			}
		})
	}
}
