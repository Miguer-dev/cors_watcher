package main

type request struct {
	URL     string            `json:"url"`
	Method  string            `json:"method"`
	Headers map[string]string `json:"headers,omitempty"`
	Data    string            `json:"data,omitempty"`
}

// foreach origin duplicate request
func (r request) addRequestsByOrigins(origins []string) []request {
	var requests []request

	for _, origin := range origins {
		copyRequest := r

		copyRequest.Headers = make(map[string]string)
		for key, value := range r.Headers {
			copyRequest.Headers[key] = value
		}

		copyRequest.Headers["Origin"] = origin

		requests = append(requests, copyRequest)
	}

	return requests
}
