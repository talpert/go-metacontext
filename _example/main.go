package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/talpert/go-metacontext"
	"github.com/InVisionApp/rye"
	"github.com/talpert/go-metacontext/middleware/ryemc"
)

// The struct representing the metadata that the handler
// is expecting to be passed along with the body
type Metadata struct {
	Name string `json:"name"`
	Size int    `json:"size"`
}

// The body struct representing the JSON body that
// the handler is expecting
type Body struct {
	Status string `json:"status"`
	Value  int    `json:"value"`
}

// An example showing use of metacontext with the provided rye middleware
func main() {
	mwh := rye.NewMWHandler(rye.Config{})

	// setup the test server
	s := httptest.NewServer(mwh.Handle([]rye.Handler{
		// middleware to parse metadata and save it to the request context
		ryemc.ParseMetaContextMiddleware,
		// handler that prints metadata and body from context
		exampleHandler,
	}))

	// make a request to the test server
	makeRequestToServer(s.URL)
}

// Fetch the metadata and request body that was added to the request context
// by the middleware and print them
func exampleHandler(rw http.ResponseWriter, r *http.Request) *rye.Response {
	// fetch metadata from context and marshal it to the appropriate struct
	meta := Metadata{}
	if err := metacontext.GetMetadata(r.Context(), &meta); err != nil {
		return &rye.Response{Err: err}
	}

	// fetch body from context and marshal it to the appropriate struct
	body := Body{}
	if err := metacontext.GetBody(r.Context(), &body); err != nil {
		return &rye.Response{Err: err}
	}

	fmt.Printf("Metadata: %+v\n", meta)
	fmt.Printf("Body: %+v\n", body)

	return nil
}

// Make a POST request to the provided URL with metadata and body
// The JSON body is marshaled using the helper which will append
// the metadata to the body in the correct format
func makeRequestToServer(url string) {
	// marshal post body to JSON with metadata
	jsonBody, err := metacontext.MarshalWithMetadata(
		&Metadata{
			Name: "my metadata",
			Size: 24,
		},
		&Body{
			Status: "good to go",
			Value:  32,
		})
	if err != nil {
		log.Printf("failed to marshal json body: %v", err)
		return
	}

	// build the request
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
	if err != nil {
		log.Printf("failed to build request: %v", err)

		return
	}

	// make the request
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Printf("failed to make request: %v", err)

		return
	}

	// print the response code for confirmation
	fmt.Println("response code:", resp.StatusCode)
}
