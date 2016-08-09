package main

import (
	"io"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kolide/kolide-ose/sessions"
)

// ServerError is a helper which accepts a string error and returns a map in
// format that is required by gin.Context.JSON
func ServerError(e string) *map[string]interface{} {
	return &map[string]interface{}{
		"error": e,
	}
}

// DatabaseError emits a response that is appropriate in the event that a
// database failure occurs, a record is not found in the database, etc
func DatabaseError(c *gin.Context) {
	c.JSON(500, ServerError("Database error"))
}

// UnauthorizedError emits a response that is appropriate in the event that a
// request is received by a user which is not authorized to carry out the
// requested action
func UnauthorizedError(c *gin.Context) {
	c.JSON(401, ServerError("Unauthorized"))
}

// MalformedRequestError emits a response that is appropriate in the event that
// a request is received by a user which does not have required fields or is in
// some way malformed
func MalformedRequestError(c *gin.Context) {
	c.JSON(400, ServerError("Malformed request"))
}

func createEmptyTestServer(db *gorm.DB) *gin.Engine {
	server := gin.New()
	server.Use(DatabaseMiddleware(db))
	server.Use(SessionBackendMiddleware)
	return server
}

// Adapted from https://goo.gl/03Qxiy
func DatabaseMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("DB", db)
		c.Next()
	}
}

// NewSessionManager allows you to get a SessionManager instance for a given
// web request. Unless you're interacting with login, logout, or core auth
// code, this should be abstracted by the ViewerContext pattern.
func NewSessionManager(c *gin.Context) *sessions.SessionManager {
	return &sessions.SessionManager{
		Request: c.Request,
		Backend: GetSessionBackend(c),
		Writer:  c.Writer,
	}
}

// CreateServer creates a gin.Engine HTTP server and configures it to be in a
// state such that it is ready to serve HTTP requests for the kolide application
func CreateServer(db *gorm.DB, w io.Writer) *gin.Engine {
	server := gin.New()
	server.Use(DatabaseMiddleware(db))
	server.Use(SessionBackendMiddleware)

	sessions.Configure(&sessions.SessionConfiguration{
		CookieName:     "KolideSession",
		JWTKey:         config.App.JWTKey,
		SessionKeySize: config.App.SessionKeySize,
		Lifespan:       config.App.SessionExpirationSeconds,
	})

	// TODO: The following loggers are not synchronized with each other or
	// logrus.StandardLogger() used through the rest of the codebase. As
	// such, their output may become intermingled.
	// See https://github.com/Sirupsen/logrus/issues/391

	// Ginrus middleware logs details about incoming requests using the
	// logrus WithFields
	requestLogger := logrus.New()
	requestLogger.Out = w
	server.Use(ginrus.Ginrus(requestLogger, time.RFC3339, false))

	// Recovery middleware recovers from panic(), returning a 500 response
	// code and printing the panic information to the log
	recoveryLogger := logrus.New()
	recoveryLogger.WriterLevel(logrus.ErrorLevel)
	recoveryLogger.Out = w
	server.Use(gin.RecoveryWithWriter(recoveryLogger.Writer()))

	// Kolide react entrypoint
	server.HTMLRender = loadTemplates("react.tmpl")
	server.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "react.tmpl", gin.H{})
	})
	// Kolide assets
	server.Use(static.Serve("/assets", BinaryFileSystem("/build")))

	v1 := server.Group("/api/v1")

	// Kolide application API endpoints
	kolide := v1.Group("/kolide")

	kolide.POST("/login", Login)
	kolide.GET("/logout", Logout)

	kolide.POST("/user", GetUser)
	kolide.PUT("/user", CreateUser)
	kolide.PATCH("/user", ModifyUser)
	kolide.DELETE("/user", DeleteUser)

	kolide.PATCH("/user/password", ChangeUserPassword)
	kolide.PATCH("/user/admin", SetUserAdminState)
	kolide.PATCH("/user/enabled", SetUserEnabledState)

	kolide.POST("/user/sessions", GetInfoAboutSessionsForUser)
	kolide.DELETE("/user/sessions", DeleteSessionsForUser)

	kolide.DELETE("/session", DeleteSession)
	kolide.POST("/session", GetInfoAboutSession)

	// osquery API endpoints
	osquery := v1.Group("/osquery")
	osquery.POST("/enroll", OsqueryEnroll)
	osquery.POST("/config", OsqueryConfig)
	osquery.POST("/log", OsqueryLog)
	osquery.POST("/distributed/read", OsqueryDistributedRead)
	osquery.POST("/distributed/write", OsqueryDistributedWrite)

	return server
}
