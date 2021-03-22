package project

import (
	"database/sql"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

// Config приложения
type Database struct {
	conn *sql.DB
}

// NewConfig - создание новой конфигурации приложения на основе файла .env
func Open(config *DatabaseConfig) (*Database, error) {
	db, err := sql.Open("mysql", config.Login+":@("+config.Host+":"+strconv.Itoa(int(config.Port))+")/"+config.Name)
	if err != nil {
		return nil, err
	}
	return &Database{
		conn: db,
	}, nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}

// Добавляет курс валюты.
func (db *Database) Insert(currency *Currency) error {
	_, err := db.conn.Exec("INSERT INTO currency (`name`, rate, insert_dt) VALUES (?, ?, NOW())", currency.Name, currency.Rate)
	return err
}

// Изменяет курс валюты.
func (db *Database) Update(currency *Currency) error {
	_, err := db.conn.Exec("UPDATE currency SET rate=?, insert_dt=NOW() WHERE name=?", currency.Rate, currency.Name)
	return err
}

// Удаляет курс валюты.
func (db *Database) DeleteByID(id int) error {
	_, err := db.conn.Exec("DELETE FROM currency WHERE id=?", id)
	return err
}

// Возвращает валюту по ID.
func (db *Database) GetByID(id int) (currency *Currency, err error) {
	c := Currency{}
	err = db.conn.QueryRow("SELECT c.id, c.name, c.rate, c.insert_dt FROM currency c WHERE id=?", id).Scan(
		&c.ID,
		&c.Name,
		&c.Rate,
		&c.InsertDt,
	)
	currency = &c
	return
}

// Возвращает все курсы из БД.
func (db *Database) GetAll() ([]Currency, error) {
	var res []Currency
	rows, err := db.conn.Query("SELECT id, name, rate FROM currency")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		c := Currency{}
		err := rows.Scan(&c.ID, &c.Name, &c.Rate)
		if err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

// Возвращает все курсы из БД.
func (db *Database) GetAllForPage(page, itemsPerPage int) ([]Currency, error) {
	var res []Currency
	rows, err := db.conn.Query("SELECT id, name FROM currency ORDER BY name LIMIT ?,?", (page-1)*itemsPerPage, itemsPerPage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		c := Currency{}
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		res = append(res, c)
	}
	return res, nil
}

func (db *Database) UpdateCountTotalCache() error {
	count := 0
	err := db.conn.QueryRow("SELECT COUNT(*) FROM currency").Scan(&count)
	if err != nil {
		return err
	}

	cacheCount := 0
	err = db.conn.QueryRow("SELECT COUNT(*) FROM cache").Scan(&cacheCount)
	if err != nil {
		return err
	}

	if cacheCount > 0 {
		_, err = db.conn.Exec("UPDATE cache SET count_total=?", count)
	} else {
		_, err = db.conn.Exec("INSERT INTO cache (count_total) VALUES (?)", count)
	}
	return err
}

// Возвращает количество валют.
func (db *Database) GetCountAll() (int, error) {
	v, err := db.getCachedValue("count_total")
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(v)
}

func (db *Database) getCachedValue(fieldName string) (v string, err error) {
	err = db.conn.QueryRow("SELECT " + fieldName + " FROM cache LIMIT 0, 1").Scan(&v)
	if err != nil {
		return
	}
	return
}
