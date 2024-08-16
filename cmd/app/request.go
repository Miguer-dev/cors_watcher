package main

type request struct {
	url     string            `json:"url"`
	method  string            `json:"method"`
	headers map[string]string `json:"headers,omitempty"`
	data    string            `json:"data,omitempty"`
}

func (r request) addRequestsByOrigins(origins []string) []request {
	var requests []request

	for _, origin := range origins {
		copyRequest := r

		copyRequest.headers = make(map[string]string)
		for key, value := range r.headers {
			copyRequest.headers[key] = value
		}

		copyRequest.headers["Origin"] = origin

		requests = append(requests, copyRequest)
	}

	return requests
}
