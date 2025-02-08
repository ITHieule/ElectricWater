package types

type BuserType struct {
	UserID   string `json:"USERID" binding:"required"`
	Password string `json:"PWD" binding:"required"`
	DBChoice string `json:"DB_CHOICE" binding:"required"`
}
