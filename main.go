package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"os"
	"github.com/golang/gddo/httputil/header"
)

// Creating expected data structure
type manageFile struct {
    Action string `json:"action"`
}

// manage_file is called by the Mux handler in main.  Logic is done within this handler to ensure we have
// the expected payload.  We'll do our best to provide input to the end-user if the payload malformed. 
// Only two actions are actionable at this time
func manage_file(w http.ResponseWriter, r *http.Request) {
	
	// Making variable so we can compare the values
	readAction := manageFile{
		Action: "read",
	}
	downloadAction := manageFile{
		Action: "download",
	}

	// application/json header is mandatory
	if r.Header.Get("Content-Type") != "" {
		value, _ := header.ParseValueAndParams(r.Header, "Content-Type")
        if value != "application/json" {
            msg := "Content-Type header must be application/json"
            http.Error(w, msg, http.StatusUnsupportedMediaType)
            return
        }
    }

	// Preventing accidental or suspicisoulsy large requests
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)

	dec := json.NewDecoder(r.Body)
	// Error out if field is not found or defiened
    dec.DisallowUnknownFields()

	var mf manageFile

	err := dec.Decode(&mf)
    if err != nil {
        var syntaxError *json.SyntaxError
        var unmarshalTypeError *json.UnmarshalTypeError
        
		// Beginning of analyzing and providing responses to malformed payloads
		switch {

		case errors.As(err, &syntaxError):
            msg := fmt.Sprintf("JSON formatting error in body (at position %d)", syntaxError.Offset)
            http.Error(w, msg, http.StatusBadRequest)
			fmt.Fprint(w, "HTTP Status: ", http.StatusBadRequest, "\n")

		case errors.Is(err, io.ErrUnexpectedEOF):
            msg := fmt.Sprintf("JSON body badly formed")
            http.Error(w, msg, http.StatusBadRequest)
			fmt.Fprint(w, "HTTP Status: ", http.StatusBadRequest, "\n")

		case errors.As(err, &unmarshalTypeError):
            msg := fmt.Sprintf("Unexpected value type for the %q field (at position %d), string expected", unmarshalTypeError.Field, unmarshalTypeError.Offset)
            http.Error(w, msg, http.StatusBadRequest)
			fmt.Fprint(w, "HTTP Status: ", http.StatusBadRequest, "\n")

		case strings.HasPrefix(err.Error(), "json: unknown field "):
            fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
            msg := fmt.Sprintf("Field %s doesn't exist, please check documentation", fieldName)
            http.Error(w, msg, http.StatusBadRequest)
			fmt.Fprint(w, "HTTP Status: ", http.StatusBadRequest, "\n")

		case errors.Is(err, io.EOF):
            msg := "Request body must not be empty, action field is required"
            http.Error(w, msg, http.StatusBadRequest)
			fmt.Fprint(w, "HTTP Status: ", http.StatusBadRequest, "\n")

		case err.Error() == "http: request body too large":
            msg := "Request body must not be larger than 1MB"
            http.Error(w, msg, http.StatusRequestEntityTooLarge)
			fmt.Fprint(w, "HTTP Status: ", http.StatusRequestEntityTooLarge, "\n")

		default:
            log.Println(err.Error())
            http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        }
        return
    }
    
	err = dec.Decode(&struct{}{})
	if err != io.EOF {
        msg := "Request body must only contain a single JSON object"
        http.Error(w, msg, http.StatusBadRequest)
        return
    }

	// Provide confirmation of what was parsed
    fmt.Fprint(w, "HTTP Status: ", http.StatusOK, "\n")
	fmt.Fprint(w, "action: ", mf, "\n")

	// Comparing action payload to determine if action needs to be made
	if mf == downloadAction {
		fmt.Fprint(w, "Download executed")
		// Call func found in download.go to download file
		learningContainerDownload()
	}

	if mf == readAction {
		fmt.Fprint(w, "Downloading latest version..\n\n")
		// Call func found in download.go to download file.  This guarentees the file will always exist and is not outdated
		learningContainerDownload()
		contents, err := os.ReadFile("downloads/sample-text-file.txt")
		if err != nil {
			log.Fatalf("unable to read file: %v", err)
		}
		fmt.Fprint(w, string(contents))

	}

}

// main calls upon a handler for the registered route. At this time the application supports one route.
func main() {
	
	mux := http.NewServeMux()
	mux.HandleFunc("/manage_file", manage_file)

	err := http.ListenAndServe("0.0.0.0:8080", mux)
	log.Fatal(err)

}