package auth_test

import (
	"net/http/httptest"
	"testing"

	"bux-wallet/domain/users"
	"bux-wallet/transports/http/auth"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestTerminateSession(t *testing.T) {
	// Arrange
	ctx := setupTest()
	session := sessions.Default(ctx)

	session.Set(auth.SessionAccessKeyId, gofakeit.HexUint256())
	session.Set(auth.SessionAccessKey, gofakeit.HexUint256())
	session.Set("random-key", gofakeit.HipsterWord())
	session.Save()

	// Act
	auth.TerminateSession(ctx)

	// Assert
	session = sessions.Default(ctx)

	assert.Nil(t, session.Get(auth.SessionAccessKeyId))
	assert.Nil(t, session.Get(auth.SessionAccessKey))
	assert.Nil(t, session.Get("random-key"))
}

func TestUpdateSession(t *testing.T) {
	// Arrange
	ctx := setupTest()

	accessKey := users.AccessKey{
		Id:  gofakeit.HexUint256(),
		Key: gofakeit.HexUint256(),
	}

	userId := gofakeit.IntRange(0, 1000)

	// Act
	auth.UpdateSession(ctx, accessKey, int(userId))

	// Assert
	session := sessions.Default(ctx)

	assert.Equal(t, accessKey.Id, session.Get(auth.SessionAccessKeyId))
	assert.Equal(t, accessKey.Key, session.Get(auth.SessionAccessKey))
	assert.Equal(t, userId, session.Get(auth.SessionUserId))
}

func setupTest() (ctx *gin.Context) {
	gin.SetMode(gin.TestMode)

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	store := memstore.NewStore([]byte("secret"))

	ctx.Request = httptest.NewRequest("GET", "/test", nil)
	sessions.Sessions("test", store)(ctx)

	return
}
