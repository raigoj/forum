package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"text/template"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}
func testHandler(w http.ResponseWriter, r *http.Request) {
	tpl.ExecuteTemplate(w, "test.html", nil)
	// if err != nil {
	// 	panic(err)
	// }
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// http.Redirect(w, r, "/posts", 303)
	// http.Redirect(w, r, "/posts", 303)

	fmt.Println("*****loginHandler running*****")
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "login.html", nil)
		return
	}
	fmt.Println("*****loginAuthHandler running*****")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("username:", username, "password:", password)

	// retrieve password from db to compare (hash) with user supplied password's hash
	var hash string
	var userID int
	stmt := "SELECT Hash, user_id FROM Users WHERE username = ?"
	row := db.QueryRow(stmt, username)
	err := row.Scan(&hash, &userID)
	fmt.Println("hash from db:", hash, "")
	if err != nil {
		fmt.Println("error selecting Hash in db by Username")
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}

	// func CompareHashAndPassword(hashedPassword, password []byte) error
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// returns nil on success
	if err == nil {
		sessionId := inMemorySession.Init(username)

		cookie := &http.Cookie{
			Name:     COOKIE_NAME,
			Value:    sessionId,
			Expires:  time.Now().Add(10 * time.Minute),
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteStrictMode,
			Secure:   false,
		}
		http.SetCookie(w, cookie) //.Redirect(w, r, "/posts", 303)
		if err := AddSessionToDB(cookie.Value, userID); err != nil {
			http.Error(w, "something wrong", http.StatusInternalServerError)
			fmt.Println("Error: ", err)
			return
		}
		authenticated = true
		//fmt.Fprint(w, "You have successfully logged in :)")
		http.Redirect(w, r, "/posts", 303)
		//w.Redirect
		return
	}
	fmt.Println("incorrect password")
	tpl.ExecuteTemplate(w, "login.html", "check username and password")

}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	// isAuthenticated(w, r)
	// if r.Method != "GET" {
	// 	http.Error(w, "Method not allowed", http.StatusBadRequest)
	// }
	rows, err := db.Query("SELECT * FROM Posts")
	checkInternalServerError(err, w)
	var funcMap = template.FuncMap{
		"idToUser": func(n int64) string {
			stmt := "SELECT username FROM Users WHERE user_id = ?"
			row := db.QueryRow(stmt, n)
			var user string
			row.Scan(&user)
			return user
		},
		"idToCategory": func(n int64) string {
			stmt := "SELECT category FROM Categories WHERE category_id = ?"
			row := db.QueryRow(stmt, n)
			var category string
			row.Scan(&category)
			return category
		},
	}

	var post Posts
	var posts []Posts
	for rows.Next() {
		err = rows.Scan(&post.Post_id, &post.Post_name, &post.Post_content,
			&post.Post_date, &post.User_id, &post.Category_id)
		checkInternalServerError(err, w)
		posts = append(posts, post)
	}
	t, err := template.New("posts.html").Funcs(funcMap).ParseFiles(filepath.Join("static/posts.html"))
	checkInternalServerError(err, w)
	err = t.Execute(w, posts)
	checkInternalServerError(err, w)
}

func commentsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM Comments")
	checkInternalServerError(err, w)
	var comments []Comments
	var comment Comments
	for rows.Next() {
		err = rows.Scan(&comment.Comment_id, &comment.Comment_content,
			&comment.Comment_date, &comment.User_id, &comment.Post_id)
		checkInternalServerError(err, w)
		comments = append(comments, comment)
	}
	t, err := template.New("comments.html").ParseFiles(filepath.Join("static/comments.html"))
	checkInternalServerError(err, w)
	err = t.Execute(w, comments)
	checkInternalServerError(err, w)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	c := http.Cookie{
		Name:   COOKIE_NAME,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	http.SetCookie(w, &c)

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	authenticated = false
	isAuthenticated(w, r)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerHandler running*****")
	if r.Method == "GET" {
		tpl.ExecuteTemplate(w, "register.html", nil)
		return
	}

	fmt.Println("*****registerAuthHandler running*****")
	r.ParseForm()
	email := r.FormValue("email")
	validEmail := IsValidEmail(email)
	username := r.FormValue("username")
	validUsername := IsValidName(username)
	password := r.FormValue("password")
	validPassword := IsValidPassword(password)
	fmt.Println("ValidPassword:", validPassword, "\nValidname:", validUsername, "\nValidEmail:", validEmail)
	if !validEmail || !validPassword || !validUsername {
		tpl.ExecuteTemplate(w, "register.html", "please check username and password criteria")
		return
	}
	stmt := "SELECT user_id FROM Users WHERE email = ?"
	row := db.QueryRow(stmt, email)
	var uID string
	err := row.Scan(&uID)
	if err != sql.ErrNoRows {
		fmt.Println("email already exists, err:", err)
		tpl.ExecuteTemplate(w, "register.html", "email already used")
		return
	}
	stmt = "SELECT user_id FROM users WHERE username = ?"
	row = db.QueryRow(stmt, username)
	err = row.Scan(&uID)
	if err != sql.ErrNoRows {
		fmt.Println("username already exists, err:", err)
		tpl.ExecuteTemplate(w, "register.html", "username already taken")
		return
	}
	// create hash from password
	var hash []byte
	// func GenerateFromPassword(password []byte, cost int) ([]byte, error)
	hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println("bcrypt err:", err)
		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}
	fmt.Println("hash:", hash)
	fmt.Println("string(hash):", string(hash))

	// func (db *DB) Prepare(query string) (*Stmt, error)
	var insertStmt *sql.Stmt
	insertStmt, err = db.Prepare("INSERT INTO Users (username, Hash, email) VALUES (?, ?, ?);")
	if err != nil {
		fmt.Println("error preparing statement:", err)
		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}
	defer insertStmt.Close()

	var result sql.Result
	//  func (s *Stmt) Exec(args ...interface{}) (Result, error)
	result, err = insertStmt.Exec(username, hash, email)
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rowsAff:", rowsAff)
	fmt.Println("lastIns:", lastIns)
	fmt.Println("err:", err)
	if err != nil {
		fmt.Println("error inserting new user")
		tpl.ExecuteTemplate(w, "register.html", "there was a problem registering account")
		return
	}
	fmt.Fprint(w, "congrats, your account has been successfully created")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "mainpage.html", struct{}{})
	if err != nil {
		panic(err)
	}

}

func createHandler(w http.ResponseWriter, r *http.Request) {
	isAuthenticated(w, r)
	if r.Method != "POST" {
		http.Redirect(w, r, "/", 301)
	}
	var post Posts
	uuidS, err := r.Cookie("sessionID")
	if err != nil {
		panic(err)
	}
	uuidVal := uuidS.Value
	smtt := "SELECT user_id FROM session WHERE uuid = ?"
	fmt.Println(smtt)
	row := db.QueryRow(smtt, uuidVal)
	fmt.Println(row)
	var userID int64
	row.Scan(&userID)
	fmt.Println(userID)
	post.Post_name = r.FormValue("PostName")
	post.Post_content = r.FormValue("PostContent")
	post.Post_date = time.Now().Format("02 Jan 2006, 15:04")
	post.User_id = userID
	post.Category_id, _ = strconv.ParseInt(r.FormValue("CategoryId"), 10, 64)
	fmt.Println(post)

	// Save to database
	stmt, err := db.Prepare(`
		INSERT INTO Posts(post_name, post_content, post_date, user_id, category_id)
		VALUES(?, ?, ?, ?, ?)
	`)
	if err != nil {
		fmt.Println("Prepare query error")
		panic(err)
	}
	_, err = stmt.Exec(post.Post_name, post.Post_content, post.Post_date,
		post.User_id, post.Category_id)
	if err != nil {
		fmt.Println("Execute query error")
		panic(err)
	}
	http.Redirect(w, r, "/posts", 301)
}
