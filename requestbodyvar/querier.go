package requestbodyvar

import (
	"bytes"
	"fmt"
	"mime"
	"net/url"
	"strings"

	"github.com/basgys/goxml2json"
	"github.com/tidwall/gjson"
)

// Querier interface is implemented by types that can query structured data.
type Querier interface {
	Query(string) string
}

// JSON struct to handle JSON data
type JSON struct {
	buf *bytes.Buffer
}

func (j JSON) Query(key string) string {
	return getJSONField(j.buf, key)
}

// XML struct to handle XML data
type XML struct {
	buf *bytes.Buffer
}

func (x XML) Query(key string) string {
	json, err := xml2json.Convert(x.buf)
	if err != nil {
		return ""
	}
	return getJSONField(json, key)
}

// Form struct to handle application/x-www-form-urlencoded form data
type Form struct {
	buf *bytes.Buffer
}

func (f Form) Query(key string) string {
	// Parse the form data
	values, err := url.ParseQuery(f.buf.String())
	if err != nil {
		return ""
	}
	// Return the value associated with the given key
	return values.Get(key)
}

// newQuerier creates a new Querier based on the content type.
func newQuerier(buf *bytes.Buffer, contentType string) (Querier, error) {
	mediaType := "application/json"
	if contentType != "" {
		var err error
		mediaType, _, err = mime.ParseMediaType(contentType)
		if err != nil {
			return nil, err
		}
	}

	switch {
	case mediaType == "application/json":
		return JSON{buf: buf}, nil
	case strings.HasSuffix(mediaType, "/xml"):
		// application/xml or text/xml
		return XML{buf: buf}, nil
	case mediaType == "application/x-www-form-urlencoded":
		// form data
		return Form{buf: buf}, nil
	default:
		return nil, fmt.Errorf("unsupported Media Type: %q", mediaType)
	}
}

// getJSONField gets the value of the given field from the JSON body,
// which is buffered in buf.
func getJSONField(buf *bytes.Buffer, key string) string {
	if buf == nil {
		return ""
	}
	value := gjson.GetBytes(buf.Bytes(), key)
	return value.String()
}
