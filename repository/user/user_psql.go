package userRepository

import (
	"database/sql"
	"user-authentication/models"
)

type UserRepository struct{}

func (u UserRepository) SignUp(db *sql.DB, user models.User) (int, error) {
	stmt := "insert into users1 (email, password, name, status) values($1, $2, $3, $4) RETURNING id;"
	
		err := db.QueryRow(stmt, user.Email, user.Password, user.Name, user.Status).Scan(&user.ID)

		if err != nil {
			return 0, err
		}
	
		return user.ID, nil
}

func (u UserRepository) Login(db *sql.DB, user models.User) (models.User, error) {
	
	row := db.QueryRow("select * from users1 where email=$1", user.Email)
		err := row.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Status)

		if err != nil {
			return user, err
		}
	
		return user, nil
}

func (u UserRepository) GetUsers(db *sql.DB, user models.User, users []models.User) ([]models.User, error) {

	rows, err := db.Query("select * from users1")
		
		if err != nil {
			return []models.User{}, err
		}

		for rows.Next() {
			err = rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Status,)
			
			user.Password= ""
	
			users = append(users, user)
		}

		if err != nil {
			return []models.User{}, err
		}

		return users, nil
}

func (u UserRepository) GetUser(db *sql.DB, user models.User, id int) (models.User, error) {
	rows := db.QueryRow("select * from users1 where id=$1", id)
	err := rows.Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.Status)

	user.Password=""

	return user, err
}

func (u UserRepository) AddUser(db *sql.DB, user models.User) (int, error) {
	err := db.QueryRow("insert into users1 (email, password, name, status) values($1, $2, $3, $4) RETURNING id;", user.Email, user.Password, user.Name, user.Status).Scan(&user.ID)

	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

func (u UserRepository) UpdateUser(db *sql.DB, user models.User) (int64, error) {
	result, err := db.Exec("update users1 set email=$1, name=$2, status=$3 where id=$4 RETURNING id", &user.Email, &user.Name, &user.Status, &user.ID)

	if err != nil {
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsUpdated, nil
}

func (u UserRepository) RemoveUser(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("delete from users1 where id=$1", id)
		
	
		
	if err != nil {
		return 0, err
	}

	rowsDeleted, err := result.RowsAffected()
		

	if err != nil {
		return 0, err
	}

	return rowsDeleted, nil
}
