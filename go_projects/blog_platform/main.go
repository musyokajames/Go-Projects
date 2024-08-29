package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func userRegistration(w http.ResponseWriter, r *http.Request) {
	// Connect to the database
	db, err := sql.Open("mysql", "Musyoka:mysqlpassword@tcp(127.0.0.1:3306)/blogPlatformDB")
	if err != nil {
		log.Println("Error connecting to database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("Database connection failed!", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	log.Println("Succesfully connected to database")

	// link to the HTML files
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/registration.html"))
		tmpl.Execute(w, nil)

	} else if r.Method == http.MethodPost {
		err := r.ParseForm()

		if err != nil {
			log.Println("Error parsing form:", err)
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		//Collect form data
		name := r.FormValue("name")
		username := r.FormValue("username")
		password := r.FormValue("password")
		passwordConfirm := r.FormValue("confirmpassword")

		// Check if any required fields are empty
		if name == "" || username == "" || password == "" || passwordConfirm == "" {
			log.Println("All fields are required")
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		// Check if passwords match
		if password != passwordConfirm {
			log.Println("Passwords do not match")
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
			return
		}

		// Check if username already exists
		var existingUsername string
		err = db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&existingUsername)
		if err == nil {
			log.Println("Username already taken")
			http.Error(w, "Username is already taken", http.StatusConflict)
			return
		} else if err != sql.ErrNoRows {
			log.Println("Database error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Generate a hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error generating hash password", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// insert into the database
		query := "INSERT INTO users(name, username, password_hash) VALUES (?, ?, ?)"
		results, err := db.Exec(query, name, username, hashedPassword)
		if err != nil {
			log.Println("Error inserting values into database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		//Check how many rows in the database were affected
		rowsAffected, err := results.RowsAffected()
		if err != nil {
			log.Println("Error getting affected rows:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		log.Printf("Data inserted succesfully, %d rows affected", rowsAffected)

		http.Redirect(w, r, "/login", http.StatusSeeOther)

	}
}

var jwtSecret = []byte("mySecretKey")

func userLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/login.html"))
		tmpl.Execute(w, nil)
	} else if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			log.Println("Error parsing form:", err)
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Connect to database
		datasourcename := "Musyoka:mysqlpassword@tcp(127.0.0.1:3306)/blogPlatformDB"
		db, err := sql.Open("mysql", datasourcename)
		if err != nil {
			log.Println("Error connecting to database", err)
			return
		}
		defer db.Close()

		// query the database
		var storedHashedPassword string
		query := "SELECT password_hash FROM users WHERE username = ?"
		err = db.QueryRow(query, username).Scan(&storedHashedPassword)
		if err == sql.ErrNoRows {
			fmt.Println("User not found")
			http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
			return
		} else if err != nil {
			fmt.Println("Error querying database:", err)
			http.Error(w, "Error querying database", http.StatusInternalServerError)
			return
		}

		// Compare hashedPassword to the input password
		err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
		if err != nil {
			log.Println("Password incorrect", err)
			http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
			return
		}

		fmt.Println("Login succesful!")

		// If login is succesful create a JWT token
		claims := jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(12 * time.Hour).Unix(),
		}

		// Generating the token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// Signing the JWT
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			log.Println("Error generating token:", err)
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			MaxAge:   3600 * 24, //24 hours
		})

		//Redirect to the tasks view page
		http.Redirect(w, r, "/view", http.StatusSeeOther)

	}
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				http.Error(w, "Missing token", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		tokenString := cookie.Value

		//Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		//Continue to next handler
		next.ServeHTTP(w, r)
	})
}

type Post struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Author  string `json:"username"`
}

func viewBlog(w http.ResponseWriter, r *http.Request) {
	// connect to the database
	db, err := sql.Open("mysql", "Musyoka:mysqlpassword@tcp(127.0.0.1:3306)/blogPlatformDB")
	if err != nil {
		log.Println("Error connecting to database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	//Retrieve the JWT token from the cookie
	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	//Define a struct to hold the JWT claims
	type Claims struct {
		Username string `json:"username"`
		Exp      int64  `json:"exp"`
		jwt.RegisteredClaims
	}

	//Parse the token and extract the username
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	username := claims.Username

	var posts []Post

	// retrieve the files from the database
	query := "SELECT id, title FROM posts"
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Error querrying database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var post Post
		if err := rows.Scan(&post.ID, &post.Title); err != nil {
			log.Println("Error scanning row:", err)
		}
		posts = append(posts, post)
	}

	if len(posts) == 0 {
		log.Println("No posts found")
		http.Error(w, "No posts available", http.StatusNotFound)
		return
	}

	//Prepare the data to pass to the template
	type TemplateData struct {
		Username string
		Posts    []Post
	}
	data := TemplateData{
		Username: username,
		Posts:    posts,
	}

	//Execute the template
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	tmpl.Execute(w, data)
}

func viewPost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr, exists := params["id"]
	if !exists || idStr == "" {
		log.Println("Error: ID parameter is missing")
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Error converting to integer:", err)
		http.Error(w, "Invalid Post ID", http.StatusBadRequest)
		return
	}

	// connect to the database
	db, err := sql.Open("mysql", "Musyoka:mysqlpassword@tcp(127.0.0.1:3306)/blogPlatformDB")
	if err != nil {
		log.Println("Error connecting to database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// query the database
	var post Post
	query := "SELECT title, content, author_username FROM posts WHERE id = ?"
	err = db.QueryRow(query, id).Scan(&post.Title, &post.Content, &post.Author)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Post not found", http.StatusNotFound)
		} else {
			log.Println("Error querying database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	type Comment struct {
		ID        int
		PostID    int
		Author    string
		Content   string
		CreatedAt time.Time
	}

	var comments []Comment
	query = "SELECT author, content, created_at FROM comments WHERE post_id = ?"
	rows, err := db.Query(query, id)
	if err != nil {
		log.Println("Error querying comments:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var comment Comment
		if err := rows.Scan(&comment.Author, &comment.Content, &comment.CreatedAt); err != nil {
			log.Println("Error scanning comment:", err)
			continue
		}
		comments = append(comments, comment)
	}

	data := struct {
		Post     Post
		Comments []Comment
	}{
		Post:     post,
		Comments: comments,
	}

	// Render the template
	tmpl, err := template.ParseFiles("templates/post.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("Error executing template:", err)
	}
}

func addComments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cookie, err := r.Cookie("token")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type Claims struct {
		Username string `json:"username"`
		jwt.RegisteredClaims
	}

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	username := claims.Username

	if err := r.ParseForm(); err != nil {
		log.Println("Error parsing form:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	log.Println("Form Values:", r.Form)

	postID := r.FormValue("post_id")
	author := username
	content := r.FormValue("content")

	if postID == "" || author == "" || content == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	postIDInt, err := strconv.Atoi(postID)
	if err != nil {
		log.Println("Error converting post_id to integer:", err)
		http.Error(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	// connect to the database
	db, err := sql.Open("mysql", "Musyoka:mysqlpassword@tcp(127.0.0.1:3306)/blogPlatformDB")
	if err != nil {
		log.Println("Error connecting to database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var postExists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM posts WHERE id = ?)", postID).Scan(&postExists)
	if err != nil {
		log.Println("Error checking post existence:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !postExists {
		http.Error(w, "Post not found", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO comments(post_id, author, content) VALUES(?, ?, ?)"
	results, err := db.Exec(query, postIDInt, author, content)
	if err != nil {
		log.Println("Error inserting comment:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	affectedRows, err := results.RowsAffected()
	if err != nil {
		log.Println("Error getting affected rows:", err)
	} else {
		log.Printf("Inserted comment. Affected rows: %d", affectedRows)
	}

	http.Redirect(w, r, fmt.Sprintf("/view/%s", postID), http.StatusSeeOther)

}

func createPost(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/new_post.html"))
		tmpl.Execute(w, nil)

	} else if r.Method == http.MethodPost {
		db, err := sql.Open("mysql", "Musyoka:mysqlpassword@tcp(127.0.0.1:3306)/blogPlatformDB")
		if err != nil {
			log.Println("Error connecting to database:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		type Claims struct {
			Username string `json:"username"`
			jwt.RegisteredClaims
		}

		claims := &Claims{}

		token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		username := claims.Username

		r.ParseForm()

		title := r.FormValue("title")
		content := r.FormValue("content")
		author := username

		// Check if any field is empty
		if title == "" || content == "" {
			log.Println("Error: All fields are required")
			http.Error(w, "All fields are required", http.StatusBadRequest)
			return
		}

		query := "INSERT INTO posts(title, content, author_username) VALUES (?, ?, ?)"
		results, err := db.Exec(query, title, content, author)
		if err != nil {
			log.Println("Error inserting to database:", err)
			return
		}

		//Check how many rows in the database were affected
		rowsAffected, err := results.RowsAffected()
		if err != nil {
			log.Println("Error getting affected rows:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}

		log.Printf("Data inserted succesfully, %d rows affected", rowsAffected)

		var existingUsername string
		err = db.QueryRow("SELECT username FROM users WHERE username = ?", author).Scan(&existingUsername)
		if err != nil {
			log.Println("Username does not exist", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/view", http.StatusSeeOther)
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/register", userRegistration).Methods("GET", "POST")
	r.HandleFunc("/login", userLogin).Methods("GET", "POST")

	r.Handle("/view", authMiddleware(http.HandlerFunc(viewBlog))).Methods("GET")
	r.Handle("/view/{id:[0-9]+}", authMiddleware(http.HandlerFunc(viewPost))).Methods("GET", "POST")
	r.Handle("/comment", authMiddleware(http.HandlerFunc(addComments))).Methods("POST")
	r.Handle("/add", authMiddleware(http.HandlerFunc(createPost))).Methods("GET", "POST")
	fmt.Println("Server starting on port 9000...")
	log.Fatal(http.ListenAndServe(":9000", r))
}
