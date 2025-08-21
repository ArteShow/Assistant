package money_database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/ArteShow/Assistant/Server/pkg/configloader"
	_ "modernc.org/sqlite"
)

type Stats struct {
	MoneyLeft int
	Sum       int
	MoneyToGo int
}

func OpenDataBase() (*sql.DB, error) {
	logFile, err := os.OpenFile("Server/log/database.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	log.SetOutput(logFile)

	dbPath, err := configloader.GetMoneyDatabasePath()
	if err != nil {
		log.Println("Config error:", err)
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		log.Println("DB open error:", err)
		return nil, err
	}

	return db, nil
}

func SetupDatabase(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS moneyData (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		money_left INTEGER NOT NULL,
		sum INTEGER NOT NULL,
		current_money INTEGER NOT NULL
	);`

	if _, err := db.Exec(query); err != nil {
		return err
	}

	_, err := db.Exec(`INSERT INTO moneyData (money_left, sum, current_money) SELECT 0, 0, 0 WHERE NOT EXISTS (SELECT 1 FROM userdata);`)
	if err != nil {
		return err
	}

	log.Println("Table created successfully.")
	return nil
}

func SetSum(db *sql.DB, sum int) error {
	_, err := db.Exec(`UPDATE moneyData SET sum = ?, money_left = ? WHERE id = 1;`, sum, sum)
	return err
}

func EditValue(db *sql.DB, field string, newValue int) error {
	if field != "money_left" && field != "sum" && field != "current_money" {
		return fmt.Errorf("invalid field: %s", field)
	}
	_, err := db.Exec(fmt.Sprintf(`UPDATE moneyData SET %s = ? WHERE id = 1;`, field), newValue)
	return err
}

func AddMoney(db *sql.DB, money int) error {
	_, err := db.Exec(`UPDATE moneyData
		SET current_money = current_money + ?,
		    money_left = money_left - ?
		WHERE id = 1;`, money, money)
	return err
}

func GetStats(db *sql.DB) (Stats, error) {
	row := db.QueryRow(`SELECT money_left, sum, current_money FROM moneyData WHERE id = 1;`)
	var s Stats
	err := row.Scan(&s.MoneyLeft, &s.Sum, &s.MoneyToGo)
	return s, err
}
