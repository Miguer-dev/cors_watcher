package main

type request struct {
	url     string
	method  string
	headers map[string]string
	data    string
}
