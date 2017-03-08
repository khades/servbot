package httpbackend

import (
	"encoding/json"
	"net/http"
)

type optionResponse struct {
	Status string
}

func options(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(optionResponse{"OK"})

}
