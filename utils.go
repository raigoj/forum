package main

import (
	"database/sql"
	"net/http"
	"regexp"
	"unicode"
)

func checkInternalServerError(err error, w http.ResponseWriter) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		panic(err)
	}
}

func isAuthenticated(w http.ResponseWriter, r *http.Request) {
	if !authenticated {
		http.Redirect(w, r, "/login", 301)
	}
}

func AddSessionToDB(uuid string, user_id int) error {
	var ins *sql.Stmt
	q := `INSERT INTO session (uuid, user_id) VALUES (?, ?)`
	ins, err := db.Prepare(q)
	if err != nil {
		panic(err)
	}
	defer ins.Close()
	if _, err := ins.Exec(uuid, user_id); err != nil {
		return err
	}
	return nil
}

func IsValidName(username string) bool {
	var nameAlphaNumeric = true
	for _, char := range username {
		if unicode.IsLetter(char) == false && unicode.IsNumber(char) == false {
			nameAlphaNumeric = false
		}
	}
	var nameLength bool
	if 5 <= len(username) && len(username) <= 50 {
		nameLength = true
	}
	if nameAlphaNumeric && nameLength {
		return true
	}
	return false

}

func IsValidPassword(password string) bool {
	var pswdLowercase, pswdUppercase, pswdNumber, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true

	for _, char := range password {
		switch {
		// func IsLower(r rune) bool
		case unicode.IsLower(char):
			pswdLowercase = true
		// func IsUpper(r rune) bool
		case unicode.IsUpper(char):
			pswdUppercase = true
		// func IsNumber(r rune) bool
		case unicode.IsNumber(char):
			pswdNumber = true

		// func IsSpace(r rune) bool, type rune = int32
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}
	if 6 < len(password) && len(password) < 60 {
		pswdLength = true
	}
	if pswdLowercase && pswdUppercase && pswdNumber && pswdLength && pswdNoSpaces {
		return true
	}
	return false
}
func IsValidEmail(email string) bool {
	// check email syntax is valid
	emailRegex, err := regexp.Compile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if err != nil {
		//fmt.Println(err)
		return false //errors.New("sorry, something went wrong")
	}
	rg := emailRegex.MatchString(email)
	if !rg {
		return false //errors.New("email address is not a valid syntax, please check again")
	}
	// check email length
	if len(email) < 4 {
		return false //errors.New("email length is too short")
	}
	if len(email) > 253 {
		return false //errors.New("email length is too long")
	}
	return true
}
