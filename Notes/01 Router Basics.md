# Router Basics
There are this chapter's objectives
- Learn the mechanism of a HTTP web server.
- Handler Functions and the Request Multiplexer (or Router).
- Request and Response Headers.
- Marshaling data structures to JSON and XML.

## Hello World

In the heart of a HTTP web server lies a multiplexer, this component is responsible for matching (the request's patterns with) all the registered routes and then it passed the control to the apropiate handler function.

![[Request Multiplexer.png]]

We will not be using any external libraries in this course, the default http package already provides a router package and we only need two functions: one for registering routes and another one for starting the server.

Create a new `banking` folder, with a `main.go`-

```go
package main

import (
"fmt"
"net/http"
)

func main() {

// Creates a new /greet endpoint
// takes a function as a seccond parameter with response and request params

http.HandleFunc("/greet", func(w http.ResponseWriter, r \*http.Request) {

fmt.Fprint(w, "Hello World")

})

// Listen on port 8000, we can pass a router as a seccond param but we are using the default one

http.ListenAndServe("localhost:8000", nil)

}
```

Run `go run main.go` to run this program, visit this endpoint located at `localhost:8000/api/greet` and the "Hello World!" message should show up.

## JSON Encoding

Let's now return some sample customer data, create a struct outside the `main` function.

```go
type Customer struct {
	Name string \`json:"full\_name"\`
	City string \`json:"city"\`
	Zipcode string \`json:"zip\_code"\`
}
```

Create a new `http.HandleFunc` for a `/Customers` endppoint which should return all registered customers.

```go
http.HandleFunc("/customers", func(w http.ResponseWriter, r \*http.Request) {
	// Create a slice of type Customer
	customers := \[\]Customer{
		{Name: "Jeremy", City: "San Jose", Zipcode: "11301"},
		{Name: "Eren", City: "Heredia", Zipcode: "12064"},
		{Name: "Levi", City: "Alajuela", Zipcode: "19804"},
	}

	// Add a Content-type: application/json header for better format
	w.Header().Add("Content-type", "application/json")
	// write customers in json format
	json.NewEncoder(w).Encode(customers)
})
```

## XML Encoding

Switching to `xml` encoding is pretty easy, we only need to alter our `header` and `encoder`, we can also alter the `struct`.

```go
type Customer struct {
	Name string \`json:"full\_name" xml:"full\_name"\`
	City string \`json:"city" xml:"city"\`
	Zipcode string \`json:"zip\_code" xml:"zip\_code"\`
}
```

```go
http.HandleFunc("/customers", func(w http.ResponseWriter, r \*http.Request) {
	// Create a slice of type Customer
	customers := \[\]Customer{
		{Name: "Jeremy", City: "San Jose", Zipcode: "11301"},
		{Name: "Eren", City: "Heredia", Zipcode: "12064"},
		{Name: "Levi", City: "Alajuela", Zipcode: "19804"},
	}
	// Add a Content-type: application/json header for better format

	//w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Content-Type", "application/xml")

	// write customers with format
	//json.NewEncoder(w).Encode(customers)
	xml.NewEncoder(w).Encode(customers)
})
```

Though we should be able to return both `json` and `xml` formats based on our request `Content-Type` header. We can use the request object and extract the `Content-Type` header and return the apropiate format.

```go
http.HandleFunc("/customers", func(w http.ResponseWriter, r \*http.Request) {

// Create a slice of type Customer

customers := \[\]Customer{

{Name: "Jeremy", City: "San Jose", Zipcode: "11301"},

{Name: "Eren", City: "Heredia", Zipcode: "12064"},

{Name: "Levi", City: "Alajuela", Zipcode: "19804"},

}

  

// if the Content-type header is application/json

if r.Header.Get("Content-Type") \== "application/xml" {

// Return a xml response

w.Header().Add("Content-Type", "application/xml")

xml.NewEncoder(w).Encode(customers)

} else {

// Return a json response

w.Header().Add("Content-Type", "application/json")

json.NewEncoder(w).Encode(customers)

}

})
```

## Refractoring & Go Mudules

Let's separate our code into modules, we'll move our server's code into `app.go` and the handler functions into `handlers.go` both inside a `app` package.

```go
package app
// app.go

  

import (

"net/http"

)

  

func Start() {

// Creates a new /greet endpoint

// takes a function as a seccond parameter with response and request params

http.HandleFunc("/greet", greet)

  

http.HandleFunc("/customers", getAllCustomers)

  

// Listen on port 8000, we can pass a router as a seccond param but we are using the default one

http.ListenAndServe("localhost:8000", nil)

}
```

```go
package app
// handlers.go
 

import (

"encoding/json"

"encoding/xml"

"fmt"

"net/http"

)

  

type Customer struct {

Name string \`json:"full\_name" xml:"full\_name"\`

City string \`json:"city" xml:"city"\`

Zipcode string \`json:"zip\_code" xml:"zip\_code"\`

}

  

func getAllCustomers(w http.ResponseWriter, r \*http.Request) {

// Create a slice of type Customer

customers := \[\]Customer{

{Name: "Jeremy", City: "San Jose", Zipcode: "11301"},

{Name: "Eren", City: "Heredia", Zipcode: "12064"},

{Name: "Levi", City: "Alajuela", Zipcode: "19804"},

}

  

// if the Content-type header is application/json

if r.Header.Get("Content-Type") \== "application/xml" {

// Return a xml response

w.Header().Add("Content-Type", "application/xml")

xml.NewEncoder(w).Encode(customers)

} else {

// Return a json response

w.Header().Add("Content-Type", "application/json")

json.NewEncoder(w).Encode(customers)

}

}

  

func greet(w http.ResponseWriter, r \*http.Request) {

fmt.Fprint(w, "Hello World")

}
```

Import this package into `main.go` and run `app.Start()`

```go
package main

  

import "banking/app"

  

func main() {

app.Start()

}
```

If you are not working in your `$GOPATH` you need to start a new Go module with `go mod init <name>` which will create a `go.mod` file. Our app should work the same as before.

### Creating our own multiplexer

We can create a new multiplexer to ground our patters into

```go
func Start() {

	// Create a new multiplexer

	mux := http.NewServeMux()
	// Creates a new /greet endpoint
	// takes a function as a seccond parameter with response and request params

	mux.HandleFunc("/greet", greet)
	mux.HandleFunc("/customers", getAllCustomers)

	// Listen on port 8000, we can pass a router as a seccond param but we are using the default one

	http.ListenAndServe("localhost:8000", mux)

}
```

## Gorilla / Mux

[Gorilla / Mux](https://github.com/gorilla/mux) is a powerfull multiplexer built for Go. We can add it as a dependency in `app.go`

```go
package app

import (
	"net/http"
	"github.com/gorilla/mux"
)
```

This dependency still need to be downloaded run: `go run main.go`, we don't need to change much in our app in order to use `Gorilla / Mux`.

```go
func Start() {

	// Create a new gorilla multiplexer
	mux := mux.NewRouter()
}
```

Still let's change this `mux` variable into `router` to avoid any namming clashes.

These external libraries simply our workflow and allow us to create our app at a faster speed.

Let's create a endpoint to retrieve a single customer, using route parameters.

```go
func Start() {
	router.HandleFunc("/customers/{customer_id}", getCustomer)
}
```

And the `getCustomer` function, here `gorilla/mux` provides a `Vars(http.Request)` object which returns a [`map`](https://golangdocs.com/maps-in-golang) with all route parameters like a client' id.

```go
func getCustomer(w http.ResponseWriter, r *http.Request) {
	// Get the route parameters with mux.Vars, passing the request object
	vars := mux.Vars(r)

	// print the customer_id key's value inside the map
	fmt.Fprint(w, vars["customer_id"])
}
```

We can restrict the parameter type using a [regular expression](https://en.wikipedia.org/wiki/Regular_expression) in our endpoint, maybe accept only numeric values.

```go
func Start() {
	// Only accept values from 0 through 9 and onward.
	router.HandleFunc("/customers/{customer_id:[0-9]+}", getCustomer)
}
```

We can also specify the http method to be used in each endpoint.

```go
func Start() {
	// Create a new gorilla multiplexer
	router := mux.NewRouter()

	// Creates a new /greet endpoint

	// takes a function as a seccond parameter with response and request params

	router.HandleFunc("/greet", greet).Methods(http.MethodGet)

	router.HandleFunc("/customers", getAllCustomers).Methods(http.MethodGet)

	// Only accept values from 0 through 9 and onward.
	router.HandleFunc("/customers/{customer\_id:\[0-9\]+}", getCustomer).Methods(http.MethodGet)
}
```