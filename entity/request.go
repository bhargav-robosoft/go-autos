package entity

type AuthRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type NewCarRequest struct {
	Name    string `json:"name" binding:"required"`
	Company string `json:"company" binding:"required"`
	Model   int    `json:"model" binding:"required"`
}

type ModifyCarRequest struct {
	Id      string `json:"id" binding:"required"`
	Name    string `json:"name"`
	Company string `json:"company"`
	Model   int    `json:"model"`
}
