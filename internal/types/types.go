package types

type User struct {
	ID        string
	FirstName string
	LastName  string
	Password  string
	Username  string
	Email     string
}

type RegisterUser struct {
	FirstName string `json:"firstname" binding:"min=3,max=100,required"`
	LastName  string `json:"lastname" binding:"min=3,max=100,required"`
	Password  string `json:"password" binding:"min=3,max=100,required"`
	Username  string `json:"username" binding:"min=3,max=100,required"`
	Email     string `json:"email" binding:"email,required"`
}
