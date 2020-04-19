package mysqldb

import "github.com/kevinwylder/sbvision"

func (sb *SBDatabase) prepareAddUser() (err error) {
	sb.addUser, err = sb.db.Prepare(`INSERT INTO users (email, username) VALUES (?, ?);`)
	return
}

// AddUser adds the user to the database
func (sb *SBDatabase) AddUser(user *sbvision.User) error {
	result, err := sb.addUser.Exec(user.Email, user.Username)
	if err != nil {
		return err
	}
	user.ID, err = result.LastInsertId()
	return err
}

func (sb *SBDatabase) prepareGetUser() (err error) {
	sb.getUser, err = sb.db.Prepare(`
SELECT 
	id, 
	email,
	username
FROM users
WHERE email = ?
	`)
	return
}

// GetUser looks up the user in the database
func (sb *SBDatabase) GetUser(email string) (*sbvision.User, error) {
	result := sb.getUser.QueryRow(email)
	var user sbvision.User
	err := result.Scan(&user.ID, &user.Email, &user.Username)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
