package metacontext

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/mitchellh/mapstructure"
)

type ctxKey struct{}

var (
	cKey ctxKey
)

type schema struct {
	Metadata interface{} `json:"metadata"`
	Body     interface{} `json:"body"`
}

// Parse the request body and return a context based on the request context
// containing the parsed body
func ParseRequest(r *http.Request) (context.Context, error) {
	wrapper := &schema{}

	if err := json.NewDecoder(r.Body).Decode(wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse request: %v", err)
	}

	//TODO: if the post does not fit the schema, then attempt to store it as body instead
	// maybe also create a readcloser with the body and set it on the request so
	// this becomes entirely transparent to the handlers

	return context.WithValue(r.Context(), cKey, wrapper), nil
}

// Parse the response body and unmarshal the components
func ParseResponse(r *http.Response, metadata, body interface{}) error {
	wrapper := &schema{}

	if err := json.NewDecoder(r.Body).Decode(wrapper); err != nil {
		return fmt.Errorf("failed to parse response: %v", err)
	}

	if err := mapstructure.Decode(wrapper.Metadata, metadata); err != nil {
		return fmt.Errorf("failed to decode metadata: %v", err)
	}

	if err := mapstructure.Decode(wrapper.Body, body); err != nil {
		return fmt.Errorf("failed to decode body: %v", err)
	}

	return nil
}

// Get the metadata out of the context and unmarshal it to i
// i must be a pointer to something that the data can unmarshal to
func GetMetadata(ctx context.Context, i interface{}) error {
	wrapper := getWrapperFromContext(ctx)
	if wrapper == nil {
		return errors.New("could not read metadata from context")
	}

	if err := mapstructure.Decode(wrapper.Metadata, i); err != nil {
		return fmt.Errorf("failed to decode metadata: %v", err)
	}

	return nil
}

// Get the body out of the context and unmarshal it to i
// i must be a pointer to something that the data can unmarshal to
func GetBody(ctx context.Context, i interface{}) error {
	wrapper := getWrapperFromContext(ctx)
	if wrapper == nil {
		return errors.New("could not read body from context")
	}

	if err := mapstructure.Decode(wrapper.Body, i); err != nil {
		return fmt.Errorf("failed to decode body: %v", err)
	}

	return nil
}

func AddMetadata(ctx context.Context, metadata interface{}) context.Context {
	wrapper := getWrapperFromContext(ctx)
	if wrapper == nil {
		wrapper = &schema{}
	}

	wrapper.Metadata = metadata

	return context.WithValue(ctx, cKey, wrapper)
}

// Marshal to JSON containing the metadata and body
func MarshalWithMetadata(metadata, body interface{}) ([]byte, error) {
	return json.Marshal(schema{
		Metadata: metadata,
		Body:     body,
	})
}

// Marshal a JSON body containing the metadata and body
func MarshalWithMetadataFromCtx(ctx context.Context, body interface{}) ([]byte, error) {
	wrapper := getWrapperFromContext(ctx)
	if wrapper == nil {
		// in this case the metadata will be empty but that's ok
		wrapper = &schema{}
	}

	wrapper.Body = body

	return json.Marshal(wrapper)
}

func getWrapperFromContext(ctx context.Context) *schema {
	val := ctx.Value(cKey)
	if val == nil {
		return nil
	}

	wrapper, ok := ctx.Value(cKey).(*schema)
	if !ok {
		// if the wrong type is stored under that key treat as non-existent
		return nil
	}

	return wrapper
}

