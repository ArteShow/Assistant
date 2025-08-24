package authorization

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/ArteShow/Assistant/Server/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SaveJwtKey(jwtKey []byte, db *sql.DB) error {
	// Ensure table exists
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS jwt_token (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		jwt_key BLOB NOT NULL
	);`)
	if err != nil {
		return err
	}

	// Insert the key (overwrite old key by clearing first)
	_, err = db.Exec(`DELETE FROM jwt_token;`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO jwt_token (jwt_key) VALUES (?)`, jwtKey)
	return err
}

func GetJwtKey(db *sql.DB) ([]byte, error) {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS jwt_token (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		jwt_key BLOB NOT NULL
	);`)
	if err != nil {
		return nil, err
	}

	var key []byte
	err = db.QueryRow(`SELECT jwt_key FROM jwt_token ORDER BY id DESC LIMIT 1;`).Scan(&key)
	if err != nil {
		return nil, err
	}

	return key, nil
}

func SetupUserTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		user_id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	);`
	_, err := db.Exec(query)
	return err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func SaveUser(db *sql.DB, username, password string) error {
	hashed, err := HashPassword(password)
	if err != nil {
		return err
	}
	_, err = db.Exec(`INSERT INTO users (username, password) VALUES (?, ?)`, username, hashed)
	return err
}

func LoginUser(db *sql.DB, username, password string) (string, error) {
	var id int64
	var hashed string
	err := db.QueryRow(`SELECT user_id, password FROM users WHERE username = ?`, username).Scan(&id, &hashed)
	if err != nil {
		log.Println(err)
		return "", errors.New("user not found")
	}
	if !CheckPasswordHash(password, hashed) {
		return "", errors.New("invalid password")
	}

	expiration := time.Now().Add(time.Hour)
	claims := &models.Claims{
		UserID: id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiration),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtKey, err := GetJwtKey(db)
	if err != nil {
		return "", errors.New("Failed to get the key")
	}
	return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string, db *sql.DB) (*models.Claims, error) {
	claims := &models.Claims{}
	jwtKey, err := GetJwtKey(db)
	if err != nil {
		return nil, errors.New("Failed to get the key")
	}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
