package model

type TokenInfo struct {
	ID    int64  `json:"id"`
	Login string `json:"login"`
}

type LoginDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type RegisterDTO struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type GetSecretBody struct {
	Title string `json:"title"`
}
