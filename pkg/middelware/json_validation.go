package middelware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	responses "github.com/IbrahimHYoussef/shared-utils-go/pkg/responses"
	"github.com/xeipuuv/gojsonschema"
)

func ValidationMiddelWare(schema string) func(http.Handler) http.Handler {
	// make sure we have a schema to load
	if len(schema) == 0 {
		log.Println("No Schema String was added Reverting to Default")
		schema = "{}"
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// read body
			var body map[string]interface{}
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Error("Error no body to read in ValidationMiddelWare")
				responses.RespondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
				return
			}

			// marshal response to the body
			err = json.Unmarshal(bodyBytes, &body)
			if err != nil {
				log.Printf("Failed to Unmarshal in ValidationMiddelWare")
				responses.RespondWithError(w, http.StatusBadRequest, "Invalid Request Payload")
				return
			}
			schemaLoader := gojsonschema.NewStringLoader(schema)

			documentLoader := gojsonschema.NewGoLoader(body)
			result, err := gojsonschema.Validate(schemaLoader, documentLoader)
			if err != nil {
				responses.RespondWithError(w, http.StatusInternalServerError, "Error Validating json")
				return
			}
			if !result.Valid() {
				var errs []string
				for _, err := range result.Errors() {
					errs = append(errs, err.String())
				}
				responses.RespondWithError(w, http.StatusBadRequest, strings.Join(errs, ", "))
				return
			}
			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
			next.ServeHTTP(w, r)
		})
	}
}
