package main

import (
	"log"
	"net/http"

	"github.com/el-mike/restrict"
	"github.com/el-mike/restrict/adapters"
	"github.com/gin-gonic/gin"
)

const (
	TestUserId = "test-user"

	UserContextKey     = "User"
	ResourceContextKey = "Resource"
)

var accessManager *restrict.AccessManager

// Authenticate and return user (Subject) that made the request.
func authenticateUser(c *gin.Context) (*User, error) {
	// ... Authentication logic (For example JWT parsing).

	return &User{
		ID:   TestUserId,
		Role: "User",
	}, nil
}

// Get Conversation by id, for example from database.
func getConversationById(id string) (*Conversation, error) {
	// ... retrieve Conversation
	return &Conversation{
		CreatedBy: "other-test-user",
		ID:        "test-conversation-1",
	}, nil
}

func WithAuth(
	actions []string,
	resource restrict.Resource,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := authenticateUser(c)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(UserContextKey, user)

		// If any previous handler fetched concrete Resource, we want to
		// override it here, so we can test it against Conditions if necessary.
		if contextResource, ok := c.Get(ResourceContextKey); ok {
			if res, ok := contextResource.(restrict.Resource); ok {
				resource = res
			}
		}

		err = accessManager.Authorize(&restrict.AccessRequest{
			Subject:  user,
			Resource: resource,
			Actions:  actions,
		})

		// If error is related to insufficient privileges, we want to
		// return appropriate response.
		if accessError, ok := err.(*restrict.AccessDeniedError); ok {
			log.Print(accessError)

			c.AbortWithError(http.StatusForbidden, accessError)
			return
		}

		// If not, some other error occurred.
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}
}

func main() {
	router := gin.Default()

	policyManager, err := restrict.NewPolicyManager(adapters.NewInMemoryAdapter(Policy), true)
	if err != nil {
		panic(err)
	}

	accessManager = restrict.NewAccessManager(policyManager)

	router.GET("/users", WithAuth(
		[]string{"read"},
		// We don't need actual data to validate if User can
		// read other Users, therefore we can use empty instance.
		restrict.UseResource("User"),
	), func(c *gin.Context) {
		users := []*User{}

		c.JSON(http.StatusOK, users)
		return
	})

	router.PATCH("/conversations/:id", func(c *gin.Context) {
		conversation, err := getConversationById(c.Param("id"))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		c.Set(ResourceContextKey, conversation)
	}, WithAuth(
		[]string{"update"},
		restrict.UseResource("Conversation"),
	), func(c *gin.Context) {
		c.String(http.StatusOK, "Conversation updated")
	})

	router.Run()
}
