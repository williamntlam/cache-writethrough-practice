package types

type Todo struct {
	ID string `json:"id"`
	Title string `json:"title"`
}

type CreateTaskRequest struct {
	Title string `json:"title" binding:"required"`
}