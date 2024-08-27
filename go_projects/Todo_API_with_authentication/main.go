package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

func userRegistration(w http.ResponseWriter, r *http.Request) {
	// connecting to my database
	dsn := "Musyoka:mysqlpassword@tcp(127.0.0.1:3306)/users_db"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Error connecting to database:", err)
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Println("Database connection failed:", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	log.Println("Succesfully connected to the database")

	// link the html form with my code
	if r.Method == http.MethodGet {
		tmpl := template.Must(template.ParseFiles("templates/registration.html"))
		tmpl.Execute(w, nil)

	} else if r.Method == http.MethodPost {
		err := r.ParseForm()

		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}

		name := r.FormValue("name")
		username := r.FormValue("username")
		password := r.FormValue("password")

		// generate a hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error generating hash password:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// now insert these values to the database
		query := "INSERT INTO users(name ,username, password_hash) VALUES (?, ?, ?)"
		result, err := db.Exec(query, name, username, hashedPassword)
		if err != nil {
			log.Println("Error writing to SQL database:", err)
			http.Error(w, "Error writing to SQL database", http.StatusInternalServerError)
			return
		}

		// Check how many rows were affected
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Println("Error getting affected rows:", err)
		} else {
			log.Printf("Data inserted successfully, %d row(s) affected", rowsAffected)
		}

	}

}

var jwtSecret = []byte("mySecretKey")

func userLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		loginTmpl := template.Must(template.ParseFiles("templates/login.html"))
		loginTmpl.Execute(w, nil)

	} else if r.Method == http.MethodPost {
		r.ParseForm()

		username := r.FormValue("username")
		password := r.FormValue("password")

		// Connect to the database
		datasourcename := "Musyoka:mysqlpassword@tcp(127.0.0.1:3306)/users_db"
		db, err := sql.Open("mysql", datasourcename)
		if err != nil {
			http.Error(w, "Error connecting to database", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Query database to fetch stored password
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

		//Compare hash and password
		err = bcrypt.CompareHashAndPassword([]byte(storedHashedPassword), []byte(password))
		if err != nil {
			fmt.Println("Password incorrect")
			http.Error(w, "Invalid login credentials", http.StatusUnauthorized)
			return
		}

		fmt.Println("Login Successful")

		// If Login is successful create a JWT token
		claims := jwt.MapClaims{
			"username": username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(), //token expires in 24 hours
		}

		//Generating the JWT token.
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		//Signing the JWT
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		//Sending the token to the client
		//w.Header().Set("Authorization", "Bearer "+tokenString)
		//fmt.Fprintf(w, "Login successful, token: %s", tokenString)

		//Set the token in a secure HTTP-only cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    tokenString,
			HttpOnly: true,
			Secure:   true, // Set to true if using HTTPS
			SameSite: http.SameSiteStrictMode,
			Path:     "/",
			MaxAge:   3600 * 24, // 24 hours
		})

		// Redirect to the tasks view page
		http.Redirect(w, r, "/view", http.StatusSeeOther)
		return

	}
}

// Middleware to Protect Routes

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// authHeader := r.Header.Get("Authorization")
		// if authHeader == "" {
		// 	http.Error(w, "Missing token", http.StatusUnauthorized)
		// 	return

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
		//tokenString := strings.TrimPrefix(authHeader, "Bearer ")

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

type Task struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
}

var tasks []Task

func viewTasks(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/tasks.html"))
	tmpl.Execute(w, tasks)
}

func viewTask(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	for _, item := range tasks {
		if item.ID == id {
			tmpl := template.Must(template.ParseFiles("templates/task.html"))
			tmpl.Execute(w, item)
			return
		}
	}
	http.Error(w, "Page not found", http.StatusNotFound)
}

func addTasks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		r.ParseForm()

		newTask := Task{
			ID:          len(tasks) + 1,
			Title:       r.FormValue("title"),
			Explanation: r.FormValue("explanation"),
		}

		tasks = append(tasks, newTask)

		http.Redirect(w, r, "/view", http.StatusSeeOther)
	} else {
		http.ServeFile(w, r, "templates/new_task.html")
		return
	}
}

func deleteTasks(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	idStr := params["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	for index, item := range tasks {
		if item.ID == id {
			tasks = append(tasks[:index], tasks[index+1:]...)
			break
		}

	}
	http.Redirect(w, r, "/view", http.StatusSeeOther)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/register", userRegistration).Methods("GET", "POST")
	r.HandleFunc("/login", userLogin).Methods("GET", "POST")

	//Protected routes (Require authentication)
	r.Handle("/view", authMiddleware(http.HandlerFunc(viewTasks))).Methods("GET")
	r.Handle("/view/{id}", authMiddleware(http.HandlerFunc(viewTask))).Methods("GET")
	r.Handle("/add", authMiddleware(http.HandlerFunc(addTasks))).Methods("POST", "GET")
	r.Handle("/delete/{id}", authMiddleware(http.HandlerFunc(deleteTasks))).Methods("POST")

	fmt.Println("Starting port on server 8080....")
	log.Fatal(http.ListenAndServe(":8080", r))
}
