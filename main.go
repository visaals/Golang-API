package main

import (
	"fmt"
	// "html"
	"log"
	"net/http"
	//"strconv"
	"io"
	"runtime"
	//"strings"
	"github.com/gorilla/mux"
	// "io/ioutil"
	"encoding/json"
	"github.com/go-sql-driver/mysql"
	"reflect"
)

//(remote mysqldb) var db, err = sql.Open("mysql", "user:pwd@tcp(ip:3306)/db)
//(local mysql db) var db, err = sql.Open("mysql", "user:pwd/db")

func main() {
	if err != nil {
		panic(err.Error())

	}
	defer db.Close()

	fmt.Println("Starting API on "+runtime.Version(), reflect.TypeOf(db).String()+"\n")
	defer fmt.Println("main.go has finished running.")

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		fmt.Println("Error connecting to server: " + err.Error())
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	fmt.Println("Database connection successful.")

	router := mux.NewRouter()
	router.HandleFunc("/create_new_user", CreateNewUser)
	router.HandleFunc("/login_handler", LoginHandler)
	router.HandleFunc("/logout_handler", LogoutHandler)
	router.HandleFunc("/", hello)
	http.Handle("/", router)
	http.ListenAndServe(":8000", router)

	//log.Fatal(http.ListenAndServe(":8000", router))
}

// type myHandler struct{}
// func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//     if h, ok := mux[r.URL.String()]; ok {
//         h(w, r)
//         return
//     }

//     io.WriteString(w, "My server: "+r.URL.String())
// }

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world")
	fmt.Fprintf(w, "Yo server")
}

type UserBasicStruct struct {
	Swaptag   string `json:"swaptag"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	Email     string `json:"email"`
	Position  string `json:"position"`
	Password  string `json:"password"`
}

func CreateNewUser(w http.ResponseWriter, r *http.Request) {

	// init user
	var us UserBasicStruct

	// decoding json POST user data into UserBasicStructs
	decoder := json.NewDecoder(r.Body)
	for {
		// if error is EOF, break, else log the fatal error if it's not nil
		if err := decoder.Decode(&us); err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}
	}

	ret, msg := InsertUserBasicInfo(&us)
	fmt.Println(ret, msg)

}

func InsertUserBasicInfo(user *UserBasicStruct) (bool, string) {

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT usr_basic_info SET swaptag=?, first_name=?, last_name=?, email=?, position=?, password=?, isLoggedIn=?") // ? = placeholder
	if err != nil {
		panic(err.Error())
		return false, "Failed to prepare statement."
	}
	defer stmtIns.Close() // Close the statement when we leave main() / the program terminates

	// execute insert statement
	res, err := stmtIns.Exec(user.Swaptag, user.Firstname, user.Lastname, user.Email, user.Position, user.Password, true)
	// return error message if failed
	if driverErr, failed := err.(*mysql.MySQLError); failed && driverErr.Number == 1062 {
		// 1062 is duplicate key
		return false, driverErr.Message

	}

	// get the last inserted id just cause
	id, err := res.LastInsertId()
	if err != nil {
		panic(err)
	}

	// print the id
	fmt.Println(id)
	return true, "Successfully inserted " + string(id)

}
