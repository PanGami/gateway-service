package user

type Response struct {
	Status  bool   `json:"status"`
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ListFilter struct {
	Paging int    `json:"paging"`
	Page   int    `json:"page"`
	Size   int    `json:"size"`
	Name   string `json:"name"`
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username" validate:"required"`
	FullName string `json:"full_name" validate:"required"`
	Password string `json:"password" validate:"required"`
}
