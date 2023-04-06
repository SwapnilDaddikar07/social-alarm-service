package repository

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"social-alarm-service/db_model"
)

type UserRepository interface {
	GetProfiles(ctx *gin.Context, phoneNumbers []string) ([]db_model.User, error)
	UserExists(ctx *gin.Context, userId string) (bool, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return userRepository{db: db}
}

func (ur userRepository) UserExists(ctx *gin.Context, userId string) (bool, error) {
	query := "SELECT EXISTS(SELECT user_id from users WHERE user_id= ?)"
	var rows []int

	dbFetchError := ur.db.Select(&rows, query, userId)
	if dbFetchError != nil {
		fmt.Println("db fetch error when checking if user id exists in the db", dbFetchError)
		return false, dbFetchError
	}
	return len(rows) == 1, nil
}

func (ur userRepository) GetProfiles(ctx *gin.Context, phoneNumbers []string) ([]db_model.User, error) {
	query, args, err := sqlx.In("select user_id , display_name , phone_number from users where phone_number IN (?);", phoneNumbers)
	if err != nil {
		fmt.Println("error when creating query", err)
		return []db_model.User{}, err
	}

	query = ur.db.Rebind(query)

	var users []db_model.User
	dbErr := ur.db.Select(&users, query, args...)
	if dbErr != nil {
		fmt.Println("error when getting user profiles for mobile numbers", dbErr)
		return []db_model.User{}, err
	}
	return users, nil
}
