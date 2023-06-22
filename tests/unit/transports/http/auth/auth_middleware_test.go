package auth_test

import (
	"bux-wallet/domain"
	mock "bux-wallet/tests/mocks"
	"bux-wallet/transports/http/auth"
	"net/http/httptest"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
)

func TestApplyToApi(t *testing.T) {
	cases := []struct {
		name string
	}{
		{
			name: "Authorize session .1: valid session",
			// accessKey:   "TBD",
			// accessKeyId: "TBD",
			// userId:      "TBD",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := setupGinContext()

			buxClientMq := mock.NewMockUserBuxClient(ctrl)
			clientFctrMq := mock.NewMockBuxClientFactory(ctrl)

			clientFctrMq.EXPECT().
				CreateWithAccessKey(gomock.Any()).
				Return(buxClientMq, nil).
				AnyTimes()

			sut := auth.NewAuthMiddleware(&domain.Services{BuxClientFactory: clientFctrMq})

			// Act
			sut.ApplyToApi(ctx)

			// Assert
			//ctx.Request.
		})
	}
}

func setupGinContext() (ctx *gin.Context) {
	gin.SetMode(gin.TestMode)

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	store := memstore.NewStore([]byte("secret"))

	ctx.Request = httptest.NewRequest("GET", "/test", nil)
	sessions.Sessions("test", store)(ctx)

	return
}
