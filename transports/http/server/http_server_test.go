package httpserver

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
)

func TestNewHttpServer(t *testing.T) {

	testLogger := zerolog.Nop()

	server := NewHttpServer(8180, &testLogger)
	server.ApplyConfiguration(WithPanicEndpoint)

	// Create a test request
	req, err := http.NewRequest("GET", "/panic", nil)
	require.NoError(t, err)

	// Create a response recorderPanic to record the response
	recorderPanic := httptest.NewRecorder()

	// Perform the request to the server
	require.NotPanics(t, func() {
		server.handler.ServeHTTP(recorderPanic, req)
	}, "The server should recover from panic")

	// Check that the response has a 500 Internal Server Error status code
	require.Equal(t, http.StatusInternalServerError, recorderPanic.Code, "Status code should be 500")

	// Additional verification: Issue another request to verify that the server is still up after a panic
	recorderOK := httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/ok", nil)
	require.NoError(t, err)
	require.NotPanics(t, func() {
		server.handler.ServeHTTP(recorderPanic, req)
	}, "The server should should replay OK")

	require.Equal(t, http.StatusOK, recorderOK.Code, "Status code should be 200")
}

func TestDebugWriter(t *testing.T) {
	// Arrange
	testLogger := zerolog.Nop()
	writer := debugWriter(&testLogger)

	// Act
	_, err := writer.Write([]byte("Debug log message"))

	// Assert
	require.NoError(t, err, "Write should not return an error")
}

func TestApplyConfiguration(t *testing.T) {
	// Arrange
	testLogger := zerolog.Nop()

	server := NewHttpServer(8180, &testLogger)
	require.NotNil(t, server, "Server should be created")

	// Act
	server.ApplyConfiguration(WithOK)

	// Assert
	require.NotNil(t, server.handler, "Gin engine should be configured")
}

func TestShutdownWithContext(t *testing.T) {
	// Arrange
	testLogger := zerolog.Nop()

	server := NewHttpServer(8180, &testLogger)
	require.NotNil(t, server, "Server should be created")

	// Act
	err := server.ShutdownWithContext(context.Background())

	// Assert
	require.NoError(t, err, "ShutdownWithContext should not return an error")
}

func WithPanicEndpoint(engine *gin.Engine) {
	engine.GET("/panic", func(c *gin.Context) {
		panic("Simulated panic")
	})
}

func WithOK(engine *gin.Engine) {
	engine.GET("/ok", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "OK"})
	})
}
