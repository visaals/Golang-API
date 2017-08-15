package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
)

type LogoutStruct struct {
	Swaptag string `json:"swaptag"`
}

// TODO destroy session token
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	var ls LogoutStruct

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&ls)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Println(ls)
	ret, msg := LogoutUser(&ls)
	fmt.Println(ret, msg)

	response := LoginLogoutResponse{Message: msg, Success: ret}
	jsonResponse, err := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

// TODO destroy session token
func LogoutUser(ls *LogoutStruct) (bool, string) {

	if swaptagExists, msg := doesSwaptagExist(&ls.Swaptag); !swaptagExists {
		return false, msg
	}
	// Prepare statement for updating isLoggedIn to false
	stmtIns, err := db.Prepare("UPDATE usr_basic_info SET isLoggedIn=? WHERE swaptag=?") // ? = placeholder
	if err != nil {
		panic(err.Error())
		return false, "Failed to prepare login statement."
	}
	defer stmtIns.Close()
	fmt.Println("Setting isLoggedIn for", ls.Swaptag, "to ", false)

	// execute insert statement
	res, err := stmtIns.Exec(false, ls.Swaptag)
	if driverErr, failed := err.(*mysql.MySQLError); failed && driverErr.Number == 1062 {
		// 1062 is duplicate keys
		return false, driverErr.Message

	}
	rowsUpdated, err := res.RowsAffected()
	if err != nil {
		return false, err.Error()
	}
	if rowsUpdated == 1 {
		return true, "Logout Successful"
	} else if rowsUpdated == 0 {
		return true, "User is already logged out"
	} else {
		return false, "Error, more than one user has been logged out. This shouldn't happen."
	}

}

type LoginStruct struct {
	Swaptag  string `json:"swaptag"`
	Password string `json:"password"`
}

type LoginLogoutResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

//TODO create session token
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// init LoginStruct
	var ls LoginStruct

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&ls)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	fmt.Println(ls)
	ret, msg := LoginUser(&ls)
	fmt.Println(ret, msg)

	response := LoginLogoutResponse{Message: msg, Success: ret}
	jsonResponse, err := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

}

// create session token
func LoginUser(ls *LoginStruct) (bool, string) {
	// check if username password match

	if swaptagExists, msg := doesSwaptagExist(&ls.Swaptag); !swaptagExists {
		return false, msg
	}
	// First, retrieve password of inputted swaptag
	stmtIns, err := db.Prepare("SELECT swaptag, password, isLoggedIn FROM usr_basic_info WHERE (swaptag=?)") // ? = placeholder
	if err != nil {
		panic(err.Error())
		return false, "Failed to prepare login statement."
	}
	defer stmtIns.Close()

	// execute insert statement
	rows, err := stmtIns.Query(ls.Swaptag)
	if driverErr, failed := err.(*mysql.MySQLError); failed && driverErr.Number == 1062 {
		// 1062 is duplicate key
		return false, driverErr.Message

	}
	fmt.Println("ls: ", ls)
	defer rows.Close()

	var swaptag string
	var password string
	var isLoggedIn bool
	// scanning resulting rows, should only be one row
	for rows.Next() {
		err := rows.Scan(&swaptag, &password, &isLoggedIn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Scanning rows result:", swaptag, password, isLoggedIn)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	theymatch := ls.Password == password

	// if they match, set isLoggedIn to true
	if theymatch {
		//$sql = "UPDATE MyGuests SET lastname='Doe' WHERE id=2";

		// Prepare statement for inserting data
		stmtIns, err := db.Prepare("UPDATE usr_basic_info SET isLoggedIn=? WHERE swaptag=?") // ? = placeholder
		if err != nil {
			panic(err.Error())
			return false, "Failed to prepare login statement."
		}
		defer stmtIns.Close()

		// execute insert statement
		res, err := stmtIns.Exec(true, &swaptag)
		if driverErr, failed := err.(*mysql.MySQLError); failed && driverErr.Number == 1062 {
			// 1062 is duplicate key
			return false, driverErr.Message

		}

		rowsUpdated, err := res.RowsAffected()
		if err != nil {
			return false, err.Error()
		}
		if rowsUpdated == 1 {
			return true, "Login Successful"
		} else if rowsUpdated == 0 {
			return true, "User is already logged in"
		} else {
			return false, "Error, more than one user has been logged in. This shouldn't happen."
		}

	}
	// else, login failed, they did not match
	return false, "Login Failed, username and password did not match"

}

func doesSwaptagExist(st *string) (bool, string) {
	var exists bool

	stmtIns, err := db.Prepare("SELECT EXISTS( SELECT swaptag FROM usr_basic_info WHERE (swaptag=?))") // ? = placeholder
	if err != nil {
		panic(err.Error())
		return false, "Failed to prepare swaptag exists statement."
	}
	defer stmtIns.Close()

	// execute insert statement
	rows, err := stmtIns.Query(st)
	if driverErr, failed := err.(*mysql.MySQLError); failed && driverErr.Number == 1062 {
		// 1062 is duplicate key
		return false, driverErr.Message
	}
	fmt.Println("SWAPTAG: ", *st)
	defer rows.Close()

	// scanning resulting rows, should only be one row
	for rows.Next() {
		err := rows.Scan(&exists)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Scanning rows result:", exists)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	var msg string
	if exists {
		msg = *st + " Swaptag exists!"
	} else {
		msg = *st + " Swaptag doesn't exist."
	}
	return exists, msg

}
