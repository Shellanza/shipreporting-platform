package webserver

import (
	"log"
	"net/http"
	"strings"

	"github.com/deeper-x/shipreporting-platform/auth"
	"github.com/deeper-x/shipreporting-platform/db"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var sessionKey = "shipreporting-session"
var portinformer = "portinformer"
var token string
var ptoken *string

// AuthRequired middleware for restricted content
var AuthRequired = func(c *gin.Context) {
	session := sessions.Default(c)
	sessionIn := session.Get(sessionKey)

	if sessionIn == nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}

	return

	// Continue down the chain to handler etc
	// c.Next()
}

// ProcessAuth auth processing
var ProcessAuth = func(c *gin.Context) (bool, string) {
	username := c.PostForm("user")
	password := c.PostForm("pass")

	lenUsername := len(strings.Trim(username, " "))
	lenPassword := len(strings.Trim(password, " "))

	if lenUsername == 0 || lenPassword == 0 {
		return false, "Empty credentials"
	}

	resp, err := auth.SignOn(username, password)

	if err != nil {
		log.Println(err)
	}

	if resp.StatusCode != http.StatusOK {
		return false, "not authorized"
	}

	token, err := auth.ReadTokenAuth(resp)

	if err != nil {
		return false, "Credentials error - Not authorized"
	}

	if len(token) > 0 {
		return true, token
	}

	return false, ""

}

// CreateSession build session
var CreateSession = func(c *gin.Context, token string) bool {
	connector := db.Connector()
	repository := db.Repository{Conn: connector}

	session := sessions.Default(c)

	userID, exists := GetUserID(repository, token)

	if exists {
		portinformer := GetManagedPortinformer(repository, userID)

		session.Set(portinformer, userID)
		session.Set(sessionKey, token)
		session.Set("managedPortinformer", portinformer)

		if err := session.Save(); err != nil {
			log.Println(err)
		}

		repository.Close()

		return true
	}

	return false
}

// GetUserID retrive user id from auth token
var GetUserID = func(repo db.Repository, authToken string) (string, bool) {
	return repo.SelectUserID(authToken)
}

// GetManagedPortinformer retrieve managed portinformer
var GetManagedPortinformer = func(repo db.Repository, userID string) string {
	return repo.SelectUserPortinformer(userID)
}

// DestroySession destroy session
var DestroySession = func(c *gin.Context) (int, string) {
	session := sessions.Default(c)
	user := session.Get(sessionKey)

	if user == nil {
		return http.StatusBadRequest, "Invalid session token"
	}

	session.Delete(sessionKey)
	session.Delete(portinformer)
	session.Delete("managedPortinformer")

	if err := session.Save(); err != nil {
		return http.StatusInternalServerError, "Failed to store session"
	}

	return http.StatusOK, "Successfully logged out"
}
