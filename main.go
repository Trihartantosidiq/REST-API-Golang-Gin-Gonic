package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	_ "github.com/lib/pq"
)

type Student struct {
	NIM       uint64
	FirstName string
	LastName  string
}

var students = []Student{}

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

// products adalah slice untuk menyimpan produk-produk
var (
	products   []Product
	nextID     = 1
	productsMu sync.Mutex
)

var db *sql.DB

func initDB() {
	var err error
	connStr := "user=postgres dbname=products sslmode=disable password=admin123"
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	if err = db.Ping(); err != nil {
		fmt.Println("Error pinging database:", err)
		return
	}
	fmt.Println("Connected to database")
}

func main() {
	initDB()
	port := "8080"

	http.HandleFunc("/products", getProductsHandler)
	http.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProductByIDHandler(w, r)
		case http.MethodPost:
			createProductHandler(w, r)
		case http.MethodPut:
			updateProductHandler(w, r)
		case http.MethodDelete:
			deleteProductHandler(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	fmt.Println("Server started at :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}

	http.ListenAndServe(":"+port, nil)
}

func getProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	productsMu.Lock()
	defer productsMu.Unlock()

	for _, product := range products {
		if product.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	http.Error(w, "Product not found", http.StatusNotFound)
}

func getProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	rows, err := db.Query("SELECT id, name, price FROM products")
	if err != nil {
		http.Error(w, "Error fetching products", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price); err != nil {
			http.Error(w, "Error scanning product", http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	response := struct {
		Status string    `json:"status"`
		Data   []Product `json:"data"`
	}{
		Status: "success",
		Data:   products,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func createProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var newProduct Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := db.QueryRow(
		"INSERT INTO products (name, price) VALUES ($1, $2) RETURNING id",
		newProduct.Name, newProduct.Price,
	).Scan(&newProduct.ID)
	if err != nil {
		http.Error(w, "Error creating product", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string  `json:"status"`
		Data   Product `json:"data"`
	}{
		Status: "success",
		Data:   newProduct,
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func updateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	var updatedProduct Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update product in the database
	result, err := db.Exec("UPDATE products SET name=$1, price=$2 WHERE id=$3", updatedProduct.Name, updatedProduct.Price, id)
	if err != nil {
		http.Error(w, "Error updating product", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking rows affected", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	response := struct {
		Status string  `json:"status"`
		Data   Product `json:"data"`
	}{
		Status: "success",
		Data:   updatedProduct,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// test123
func deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Path[len("/products/"):]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		return
	}

	// Delete product from the database
	result, err := db.Exec("DELETE FROM products WHERE id=$1", id)
	if err != nil {
		http.Error(w, "Error deleting product", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking rows affected", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(struct {
		Status string      `json:"status"`
		Data   interface{} `json:"data"`
	}{
		Status: "success",
		Data:   nil,
	})
}

func createapi() {
	// membuat aplikasi server:
	// API WEB SERVER

	// GOLANG: support untuk membuat API
	// dengan menggunakan NATIVE PACKAGE (net/http)
	// dengan menggunakan WEB FRAMEWORK (gin gonic)

	port := "8080"
	version := "/api/v1"

	http.HandleFunc(version+"/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello From My Awesome WEB SERVER API")
	})
	// student api
	http.HandleFunc(version+"/students", func(w http.ResponseWriter, r *http.Request) {
		// METHOD GET
		if r.Method == http.MethodGet {
			// check apakah ada param ID / NIM
			path := strings.Split(r.URL.Path, "/students")
			if len(path) > 1 {
				if path[1] != "" {
					id := path[1]
					// search student with same nim
					for _, student := range students {
						if fmt.Sprintf("%v", student.NIM) == id {
							json.NewEncoder(w).Encode(student)
							return
						}
					}
					http.Error(w, "Student Not Found", http.StatusNotFound)
					return
				}
			}
			// return seluruh data student
			// dalam bentuk / format JSON
			json.NewEncoder(w).Encode(students)
			return
		}

		// METHOD POST
		if r.Method == http.MethodPost {
			// akan menambahkan data student
			// dari masukan API (form)
			nim := r.FormValue("nim")
			firstName := r.FormValue("first_name")
			lastName := r.FormValue("last_name")

			// error handling
			nimNumber, err := strconv.Atoi(nim)
			if err != nil {
				http.Error(w, "invalid nim", http.StatusBadRequest)
				return
			}
			// append to slice
			students = append(students, Student{
				NIM:       uint64(nimNumber),
				FirstName: firstName,
				LastName:  lastName,
			})
			fmt.Fprintf(w, "ok")
			return
		}
	})
	// student api with id
	http.HandleFunc(version+"/students/", func(w http.ResponseWriter, r *http.Request) {
		// METHOD GET
		if r.Method == http.MethodGet {
			// check apakah ada param ID / NIM
			path := strings.Split(r.URL.Path, "/students/")
			if len(path) > 1 {
				if path[1] != "" {
					id := path[1]
					// search student with same nim
					for _, student := range students {
						if fmt.Sprintf("%v", student.NIM) == id {
							json.NewEncoder(w).Encode(student)
							return
						}
					}
					http.Error(w, "Student Not Found", http.StatusNotFound)
					return
				}
			}
			// return seluruh data student
			// dalam bentuk / format JSON
			http.Error(w, "Student Not Found", http.StatusNotFound)
			return
		}
	})

	// html template
	// /students
	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		// SERVER SIDE RENDERING (SSR)
		tpl, err := template.ParseFiles("./student.html")
		if err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
		err = tpl.Execute(w, students)
		if err != nil {
			http.Error(w, "something went wrong", http.StatusInternalServerError)
			return
		}
	})

	http.ListenAndServe(":"+port, nil)
}

func otherfunction() {

	defer fmt.Println("defer from other function")
}

func errorpanic() {
	str := "10"
	num, err := strconv.Atoi(str)

	if err != nil {
		fmt.Println("something wrong")
	} else {
		fmt.Println("got some input num", num)
	}

	var err2 error

	userinput := 100

	if userinput < 0 {
		err2 = fmt.Errorf("invalid user input", userinput)

	}
	if err2 != nil {
		fmt.Println(err2.Error())
	}
}

func deferexit() {
	// num := 100
	// if num == 100 {
	// 	return
	// }

	defer fmt.Println("end of block")
	defer fmt.Println("end of block2")
	otherfunction()
	// fmt.Println("first line")
	// fmt.Println("second line")

	for i := 0; i < 10; i++ {
		fmt.Println("number :", i)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {

		defer wg.Done()
		for i := 0; i < 9; i++ {

		}
		return
	}()
	wg.Wait()

	os.Exit(1)
}
