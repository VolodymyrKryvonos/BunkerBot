package DB

import (
	"database/sql"
	"sync"
)

var dataBase *sql.DB
var once sync.Once

func GetDataBase() *sql.DB {
	once.Do(func() {
		DB, err := sql.Open("sqlite3", "/home/kryvonos/go/BunkerBot/BunkerBot/Bot/App/src/Game/DB/bunkerBD")
		if err != nil {
			panic("Can't open database " + err.Error())
		}
		dataBase = DB
	})
	return dataBase
}
