package main

func main() {
	server := MakeServer()
	server.ListenAndServe()
}
