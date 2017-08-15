# Golang-API
The backend to Swap built with Golang. Interfaces with MySQL to perform general backend operations.


Todo:

*Create Sign Up process
	1. username, email, position, password, verify password
	2. query user database for username
		3. if username doesn't exist, verify name
			else return username exists
		4. if email doesn't exist, verify email
			else return email already exists
		5. if don't match, return false
			else return true



Code I Probably Don't Need But Wanna Keep Just In Case:

// func deleteUser(w http.ResponseWriter, r *http.Request) {
//      fmt.Fprintf(w, "Hello John")
// }

// func listUsers(w http.ResponseWriter, r *http.Request) {
//      fmt.Fprintf(w, "Hello Owen")
//      i := 0
//      for i <= 10 {
// 	           fmt.Fprintf(w,strconv.Itoa( i ))
// 	           i = i + 1
//      }    
// }

    // GET MUX URL VARS /{username}
	    // vars := mux.Vars(r)
	    // for k, v := range vars { 
	    //      fmt.Fprintf(w, "key[%s] value[%s]\n", k, v)
	    // } 
     
         
    /*
    
      curl -X POST -d "{
    	\"swaptag\": \"jackcrawford\", 
    	\"first_name\": \"Jack\",
    	\"last_name\" : \"Crawford\",
    	\"email\" : \"jackcrawford@gmail.com\",
    	\"position\" : \"Frontend Engineer\",
    	\"password\" : \"pwd1234\"
    	}" http://ec2-52-37-43-157.us-west-2.compute.amazonaws.com:3456/create_new_user


    curl -X POST -d "{
    	\"swaptag\": \"visaals2\", 
    	\"first_name\": \"Visaal\",
    	\"last_name\" : \"Ambalam\",
    	\"email\" : \"visaal.sivakumar@gmail.com\",
    	\"position\" : \"Backend Engineer\",
    	\"password\" : \"pwd123\"
    	}" http://ec2-52-37-43-157.us-west-2.compute.amazonaws.com:3456/create_new_user



    curl -X POST -d "{
    	\"swaptag\": \"visaals\", 
    	\"password\" : \"pwd123\"
    	}" http://ec2-52-37-43-157.us-west-2.compute.amazonaws.com:3456/login_handler


     */
    
    // Converting the request body to a string 
	    // if b, err := ioutil.ReadAll(r.Body); err == nil {
	    //     fmt.Println(string(b))
	    // }

    // Encoding Struct to JSON
	    // output, err := json.Marshal(us)
	    // if err != nil {
	    //     fmt.Println("error: ", err)
	    // }


    //fmt.Fprintf(w, reflect.TypeOf(vars["username"]).String() + "\n" )
    //fmt.Fprintf(w, "Adding " + vars["username"] + " to database. \n" )
    //fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))


