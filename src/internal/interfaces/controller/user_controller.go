package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "my-project/generated/api"
	"my-project/internal/common"
	"my-project/internal/infrastructure/queryservice"
	"my-project/internal/usecases"
)

type UserController struct {
	usecase      *usecases.UserUsecase
	queryService *queryservice.UserQueryService
}

func NewUserController(usecase *usecases.UserUsecase, queryService *queryservice.UserQueryService) *UserController {
	return &UserController{
		usecase:      usecase,
		queryService: queryService,
	}
}

// Kiểm tra xem UserController có thỏa mãn interface ServerInterface không
var _ api.ServerInterface = (*UserController)(nil)

// Get users
// (GET /users)
func (uc *UserController) GetUsers(c *gin.Context) {
	// Call Query Service directly for Read operation (CQRS)
	users, err := uc.queryService.FetchUsers(c.Request.Context())
	if err != nil {
		common.IsErrorWithMessage(c, err, "Failed to fetch users")
		return
	}

	res := make([]api.User, len(users))
	for i, u := range users {
		id := u.ID
		name := u.Name
		createdAt := u.CreatedAt
		updatedAt := u.UpdatedAt
		res[i] = api.User{
			Id:        &id,
			Name:      &name,
			CreatedAt: &createdAt,
			UpdatedAt: &updatedAt,
		}
	}

	c.JSON(http.StatusOK, res)
}

// Create a user
// (POST /users)
func (uc *UserController) CreateUser(c *gin.Context) {
	var body api.CreateUserJSONRequestBody
	if err := common.BindJSON(c, &body); err != nil {
		return
	}

	if body.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	// Call Usecase for Command operation
	u, err := uc.usecase.CreateUser(c.Request.Context(), body.Name)
	if err != nil {
		common.IsErrorWithMessage(c, err, "Failed to create user")
		return
	}

	id := u.ID
	name := u.Name
	createdAt := u.CreatedAt
	updatedAt := u.UpdatedAt

	c.JSON(http.StatusCreated, api.User{
		Id:        &id,
		Name:      &name,
		CreatedAt: &createdAt,
		UpdatedAt: &updatedAt,
	})
}
