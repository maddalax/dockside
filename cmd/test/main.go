package main

import "net/http"

func main() {
	http.ListenAndServe(":3001", nil)
}
