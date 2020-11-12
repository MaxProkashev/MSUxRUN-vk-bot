package db

import (
	"database/sql"
	"fmt"
	"msuxrun-bot/internal/logs"
	"os"

	_ "github.com/lib/pq" // driver for postgreSQL
)

// User of bot
type User struct {
	ID   int
	Sign int

	Text  string
	Train []bool
}

func (user *User) String() string {
	return fmt.Sprintf("id: %d\nsign: %d\ntext: %s\ntrain: %v",
		user.ID,
		user.Sign,
		user.Text,
		user.Train,
	)
}

// DB for bot
var (
	DB *sql.DB

	nameTable  = "bot_users"
	createUser = func(id int) {
		str := fmt.Sprintf("INSERT INTO %s (id,sign) VALUES (%d,0);",
			nameTable,
			id,
		)
		_, err := DB.Exec(str)
		if err != nil {
			logs.DBErr("could`t not create new user. Reason: %s", err.Error())
			os.Exit(1)
		}
	}
)

// InitDB database heroku from url
func InitDB(url string) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		logs.DBErr("could`t not init database by url %s. Reason: %s", url, err.Error())
		os.Exit(1)
	}
	DB = db

	logs.Succes("init database postgres")
}

// CreateUserTable table bot_users (id INT PRIMARY KEY, sign INT)
func CreateUserTable() {
	str := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (%s INT PRIMARY KEY, %s INT);",
		nameTable,
		"id",
		"sign",
	)
	_, err := DB.Exec(str)
	if err != nil {
		logs.DBErr("could`t create %s table. Reason: %s", nameTable, err.Error())
		os.Exit(1)
	}

	logs.DB("CREATE %s TABLE", nameTable)
}

// DropUserTable bot_users
func DropUserTable() {
	_, err := DB.Exec("DROP TABLE " + nameTable + ";")
	if err != nil {
		logs.DBWarn("could`t drop %s table. Reason: %s", nameTable, err.Error())
		//os.Exit(1)
	} else {
		logs.DB("DROP %s TABLE", nameTable)
	}
}

//! User logic

// GetUser get user by init id
func (user *User) GetUser(id int) {
	user.ID = id

	str := fmt.Sprintf("SELECT * FROM %s WHERE id = %d;",
		nameTable,
		id,
	)

	var fl int
	err := DB.QueryRow(str).Scan(&id, &fl)
	switch {
	case err == sql.ErrNoRows:
		logs.DB("NEW user[%d]", user.ID)
		createUser(user.ID)
	case err != nil:
		logs.DBErr("could`t get query row. Reason: %s", err.Error())
		os.Exit(1)
	default:
		logs.DB("user[%d] select", user.ID)
		user.Sign = fl
	}
}

// ParseSign to slice bool len count
func (user *User) ParseSign(count int) {
	fl := make([]bool, count)

	for i := 0; i < count; i++ {
		if user.Sign&1 == 1 {
			fl[i] = true
		} else {
			fl[i] = false
		}
		user.Sign = user.Sign >> 1
	}

	user.Train = fl
}

// SetTrain set new slice train
func (user *User) SetTrain(count int) {
	newSign := 0

	for i := count - 1; i >= 0; i-- {
		if user.Train[i] {
			newSign = newSign<<1 + 1
		} else {
			newSign = newSign << 1
		}
	}

	str := fmt.Sprintf("UPDATE %s SET sign = %d WHERE id = %d;",
		nameTable,
		newSign,
		user.ID,
	)
	_, err := DB.Exec(str)
	if err != nil {
		logs.DBErr("could`t not update %d. Reason: %s", user.ID, err.Error())
		os.Exit(1)
	}
	logs.DB("user[%d] update", user.ID)
}

// GetAllUser get all user from table and send to out ch
func GetAllUser(out chan<- *User) {
	defer close(out)

	str := fmt.Sprintf("SELECT * FROM %s;",
		nameTable,
	)

	rows, err := DB.Query(str)
	defer rows.Close()
	if err != nil {
		logs.DBErr("could`t get query rows. Reason: %s", err.Error())
		os.Exit(1)
	}

	for rows.Next() {
		var id int
		var sign int
		if err := rows.Scan(&id, &sign); err != nil {
			logs.DBErr("could`t scan row. Reason: %s", err.Error())
			os.Exit(1)
		}
		out <- &User{
			ID:   id,
			Sign: sign,
		}
	}

}
