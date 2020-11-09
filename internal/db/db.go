package db

import (
	"database/sql"
	"fmt"
	"msuxrun-bot/internal/logs"
	"os"
	"strconv"

	_ "github.com/lib/pq" // ..
)

// DB for bot
var DB *sql.DB

// User bot
type User struct {
	ID int

	MO int //1
	TU int //2
	WE int //3
	TH int //4
	FR int //5
	SU int //6
}

// CreateUserTable - bot_users (id INT PRIMARY KEY,)
func CreateUserTable() {
	str := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s INT PRIMARY KEY, %s INT, %s INT, %s INT, %s INT, %s INT, %s INT);",
		"bot_users",
		"id",
		"mo",
		"tu",
		"we",
		"th",
		"fr",
		"su",
	)
	_, err := DB.Exec(str)
	if err != nil {
		logs.DBErr("could`t create main table")
		os.Exit(1)
	} else {
		logs.DB("create user table")
	}
}

// DropTable by name
func DropTable(name string) {
	_, err := DB.Exec("DROP TABLE " + name + ";")
	if err != nil {
		logs.DBErr("could`t not drop %s table. Reason: %s", name, err.Error())
		os.Exit(1)
	} else {
		logs.DB("drop %s table", name)
	}
}

func checkID(userID int) bool {
	rows, err := DB.Query("SELECT id FROM bot_users WHERE id = " + strconv.Itoa(userID) + ";")
	defer rows.Close()

	if err != nil {
		logs.DBErr("could`t not select id. Reason: %s", err.Error())
		os.Exit(1)
	} else {
		for rows.Next() {
			return true
		}
		return false
	}
	return false
}

func createNewUser(userID int) {
	_, err := DB.Exec("INSERT INTO bot_users (id,mo,tu,we,th,fr,su) VALUES (" + strconv.Itoa(userID) + ", 0, 0, 0, 0 ,0 ,0);")
	if err != nil {
		logs.DBErr("could`t not init new user. Reason: %s", err.Error())
		os.Exit(1)
	} else {
		logs.DB("new user init id = %d", userID)
	}
}

// SetInt ..
func SetInt(userID int, column string, value int) {
	_, err := DB.Exec("UPDATE bot_users SET " + column + " = " + strconv.Itoa(value) + " WHERE id = " + strconv.Itoa(userID) + ";")

	if err != nil {
		logs.DBErr("could`t not update %d. Reason: %s", userID, err.Error())
		os.Exit(1)
	}
}

// GetInt ..
func GetInt(userID int, column string) (value int) {
	rows, err := DB.Query("SELECT " + column + " FROM bot_users WHERE id = " + strconv.Itoa(userID) + ";")
	defer rows.Close()
	if err != nil {
		logs.DBErr("could`t not select %d. Reason: %s", userID, err.Error())
		os.Exit(1)
	} else {
		for rows.Next() {
			rows.Scan(&value)
		}
	}
	return value
}

func getFullUser(userID int) *User {
	return &User{
		ID: userID,
		MO: GetInt(userID, "mo"),
		TU: GetInt(userID, "tu"),
		WE: GetInt(userID, "we"),
		TH: GetInt(userID, "th"),
		FR: GetInt(userID, "fr"),
		SU: GetInt(userID, "su"),
	}
}

// CheckUserByID if not exist create, and return user template
func (user *User) CheckUserByID() *User {
	if !checkID(user.ID) {
		createNewUser(user.ID)
		return user
	}
	return getFullUser(user.ID)
}
