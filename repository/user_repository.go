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
	GetUser(ctx *gin.Context, userId string) (db_model.User, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return userRepository{db: db}
}

// UserExists TODO "select exists" query returning an entry in an array with ID = 0 even if user does not exists. Need to check later.
func (ur userRepository) UserExists(ctx *gin.Context, userId string) (bool, error) {
	query := "SELECT user_id from users WHERE user_id= ?"
	rows := make([]int, 0)

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

func (ur userRepository) GetUser(ctx *gin.Context, userId string) (db_model.User, error) {
	var user db_model.User

	query := "select * from users where userId = ?"

	dbErr := ur.db.Select(&user, query, userId)
	if dbErr != nil {
		fmt.Println(fmt.Sprintf("error when fetching user for user id %s", userId))
		return db_model.User{}, dbErr
	}

	return user, nil
}
