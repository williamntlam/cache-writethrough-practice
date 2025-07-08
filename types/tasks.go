package types

type Todo struct {
	ID string `json:"id"`
	Title string `json:"title"`
}

type TaskRequest struct {
	Title string `json:"title" binding:"required"`
}

type TaskResponse struct {
	ID string `json:"id"`
	Title string `json:"title"`
}