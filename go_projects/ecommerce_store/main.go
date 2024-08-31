package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type Product struct {
	ID       int     `json:"productID"`
	Name     string  `json:"productName"`
	Price    float32 `json:"price"`
	Category Category
	Quantity int    `json:"quantity_left"`
	Image    string `json:"ImageURL"`
}

type Category struct {
	ID           int    `json:"categoryID"`
	CategoryName string `json:"categoryName"`
}

func viewProducts(w http.ResponseWriter, r *http.Request) {
	//Connect to the database
	dsn := "Musyoka:mysqlpassword@tcp(127.0.0.1:3306)/myStoreDB"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Database not reachable:", err)
	}

	var products []Product

	//query the database
	query := `SELECT products.productName, products.image_url, products.price, products.quantity_left, categories.CategoryName
			FROM products
			INNER JOIN categories ON products.CategoryID = categories.categoryId;`
	rows, err := db.Query(query)
	if err != nil {
		log.Println("Error querrying database:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var product Product
		var category Category
		if err := rows.Scan(&product.Name, &product.Image, &product.Price, &product.Quantity, &category.CategoryName); err != nil {
			log.Println("Error scanning row:", err)
		}
		product.Category = category
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error with rows iteration:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//Execute the template
	tmpl := template.Must(template.ParseFiles("templates/home.html"))
	err = tmpl.Execute(w, products)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Error displaying page", http.StatusInternalServerError)
	}

}

func searchProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		//Connect to the database
		dsn := "Musyoka:mysqlpassword@tcp(127.0.0.1:3306)/myStoreDB"
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Println("Error connecting to database:", err)
			return
		}
		defer db.Close()

		if err = db.Ping(); err != nil {
			log.Fatal("Database not reachable:", err)
		}

		var products []Product
		r.ParseForm()

		productName := r.FormValue("productName")
		categoryName := r.FormValue("categoryName")

		var product Product
		// var category Category

		if productName != "" && categoryName == "" {

			query := `SELECT p.productName, p.price, p.quantity_left, p.image_url , c.categoryName
						FROM products p
						JOIN categories c ON p.categoryID = c.categoryId
						WHERE productName = ?`

			err = db.QueryRow(query, productName).Scan(&product.Name, &product.Price, &product.Quantity, &product.Image, &product.Category.CategoryName)
			if err != nil {
				if err == sql.ErrNoRows {
					http.Error(w, "Product not found", http.StatusNotFound)
				} else {
					log.Println("Error querying database:", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
				return
			}

			products = append(products, product)

		} else if productName == "" && categoryName != "" {
			query := `SELECT p.productName, p.price, p.quantity_left, p.image_url , c.categoryName
						FROM products p
						JOIN categories c ON p.categoryID = c.categoryId
						WHERE categoryName = ?`

			rows, err := db.Query(query, categoryName)
			if err != nil {
				log.Println("Error querying database:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			for rows.Next() {
				var product Product
				err := rows.Scan(&product.Name, &product.Price, &product.Quantity, &product.Image, &product.Category.CategoryName)
				if err != nil {
					log.Println("Error scanning row:", err)
					continue
				}
				products = append(products, product)
			}

			if err := rows.Err(); err != nil {
				log.Println("Error with rows iteration:", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

		} else {
			// Handle case where neither or both fields are filled in
			http.Error(w, "Please provide either a product name or a category name", http.StatusBadRequest)
			return
		}

		//Execute the template
		tmpl := template.Must(template.ParseFiles("templates/product.html"))
		err = tmpl.Execute(w, products)
		if err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Error displaying page", http.StatusInternalServerError)
		}
	}
}

func addToCart(w http.ResponseWriter, r *http.Request) {

}

func main() {

	r := mux.NewRouter()

	imageDir := "/home/musyoka/go_projects/ecommerce_store/templates/images"
	fs := http.FileServer(http.Dir(imageDir))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", fs))

	r.HandleFunc("/view", viewProducts).Methods("GET", "POST")
	r.HandleFunc("/search", searchProduct).Methods("GET", "POST")
	r.HandleFunc("/add", addToCart).Methods("GET", "POST")

	fmt.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
