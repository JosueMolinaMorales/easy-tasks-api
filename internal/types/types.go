package types

type User struct {
	ID        string
	FirstName string
	LastName  string
	Password  string
	Username  string
	Email     string
}

type AuthUser struct {
	ID        string `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
}

type RegisterUser struct {
	FirstName string `json:"firstname" binding:"min=3,max=100,required"`
	LastName  string `json:"lastname" binding:"min=3,max=100,required"`
	Password  string `json:"password" binding:"min=3,max=100,required"`
	Username  string `json:"username" binding:"min=3,max=100,required"`
	Email     string `json:"email" binding:"email,required"`
}
