package main

func main() {
	if ok, server := MakeServer(); !ok {
		return
	} else {
		server.ListenAndServe()
	}
}
