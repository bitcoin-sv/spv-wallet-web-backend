package auth_test

import (
	"net/http/httptest"
	"testing"

	"github.com/bitcoin-sv/spv-wallet-web-backend/domain/users"
	"github.com/bitcoin-sv/spv-wallet-web-backend/transports/http/auth"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestTerminateSession(t *testing.T) {
	// Arrange
	ctx := setupTest()
	session := sessions.Default(ctx)

	session.Set(auth.SessionAccessKeyID, gofakeit.HexUint256())
	session.Set(auth.SessionAccessKey, gofakeit.HexUint256())
	session.Set("random-key", gofakeit.HipsterWord())
	session.Save()

	// Act
	auth.TerminateSession(ctx)

	// Assert
	session = sessions.Default(ctx)

	assert.Nil(t, session.Get(auth.SessionAccessKeyID))
	assert.Nil(t, session.Get(auth.SessionAccessKey))
	assert.Nil(t, session.Get("random-key"))
}

func TestUpdateSession(t *testing.T) {
	// Arrange
	ctx := setupTest()

	user := users.AuthenticatedUser{
		AccessKey: users.AccessKey{
			ID:  gofakeit.HexUint256(),
			Key: gofakeit.HexUint256(),
		},
		User: &users.User{
			ID:      gofakeit.IntRange(0, 1000),
			Paymail: gofakeit.HexUint256(),
		},
		Xpriv: "xprivtest",
	}

	// Act
	auth.UpdateSession(ctx, &user)

	// Assert
	session := sessions.Default(ctx)

	assert.Equal(t, user.AccessKey.ID, session.Get(auth.SessionAccessKeyID))
	assert.Equal(t, user.AccessKey.Key, session.Get(auth.SessionAccessKey))
	assert.Equal(t, user.User.ID, session.Get(auth.SessionUserID))
	assert.Equal(t, user.User.Paymail, session.Get(auth.SessionUserPaymail))
	assert.Equal(t, user.Xpriv, session.Get(auth.SessionXPriv))
}

func setupTest() (ctx *gin.Context) {
	gin.SetMode(gin.TestMode)

	ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
	store := memstore.NewStore([]byte("secret"))

	ctx.Request = httptest.NewRequest("GET", "/test", nil)
	sessions.Sessions("test", store)(ctx)

	return
}
