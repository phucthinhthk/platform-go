package types

// LoginAdministrator is the authenticated principal placed in Gin's context.
type LoginAdministrator struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	AccountType string `json:"accountType"`
}
