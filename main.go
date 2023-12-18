package main

import (
	"context"
	"flag"     // Get arguments
	"fmt"      // Formatting
	"log"      // Logging responses and requests
	"net/http" // Build connection between client and server
	"regexp"   // Creating regular expression for path routing
	"strings"
)

type route struct {
	method  string
	regex   *regexp.Regexp
	handler http.HandlerFunc
}

// routines used to get Path parameters
type ctxKey struct{}

func getField(r *http.Request, index int) string {
	fields := r.Context().Value(ctxKey{}).([]string)
	return fields[index]
}

func newRoute(method, pattern string, handler http.HandlerFunc) route {
	return route{method, regexp.MustCompile("^" + pattern + "$"), handler}
}

var routes = []route{
	newRoute("GET", "/", Home),
	newRoute("GET", "/about", About),
	newRoute("GET", "/projects", Projects),
}

// Handler function
func Home(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Home at URL: " + r.URL.Path)

	http.ServeFile(w, r, "./html/index.html")
}

func About(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving About at URL: " + r.URL.Path)

}

func Projects(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving Projects at URL: " + r.URL.Path)
}

func MainServer(w http.ResponseWriter, r *http.Request) {
	var allow []string
	for _, route := range routes {
		matches := route.regex.FindStringSubmatch(r.URL.Path)
		if len(matches) > 0 {
			if r.Method != route.method {
				allow = append(allow, route.method)
				continue
			}

			ctx := context.WithValue(r.Context(), ctxKey{}, matches[1:])
			route.handler(w, r.WithContext(ctx))
			return
		}
	}
	if len(allow) > 0 {
		w.Header().Set("Allow", strings.Join(allow, ", "))
		http.Error(w, "405 method now allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

func main() {

	// Set a port to listen on
	var portnum = flag.Int("p", 8080, "Specify the port number the server will listen on")
	flag.Parse()

	log.Printf("Starting up http server on port :%d\n", *portnum)

	// Registering handler functions
	http.HandleFunc("/", MainServer)
	http.HandleFunc("/about", MainServer)
	http.HandleFunc("/projects", MainServer)

	log.Printf("Close server connection by Ctrl-C\n")

	// Spinning up the server to start listening and serving
	err := http.ListenAndServe(fmt.Sprintf(":%d", *portnum), nil)
	if err != nil {
		log.Fatal(err)
	}
}
