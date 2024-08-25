package auth

type ClientResponse struct {
	Status  bool        `json:"status"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Data struct {
	IsFirstLogin int      `json:"isFirstLogin"`
	UserID       string   `json:"userId"`
	IsBod        string   `json:"isBod"`
	RoleID       int      `json:"roleId"`
	RoleName     string   `json:"roleName"`
	Permissions  []string `json:"permissions"`
	Token        string   `json:"token"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	FcmToken string `json:"fcmToken"`
}

type LogoutRequest struct {
	UserId string `json:"userId"`
	Token  string `json:"token"`
}
