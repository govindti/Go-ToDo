package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type HeathResponse struct {
	Status  string `json:"status`
	Message string `json:"message`
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	health := HeathResponse{
		Status:  "Ok",
		Message: "Api is Up & Running",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(health)
}
func main() {
	http.HandleFunc("/health", healthHandler)

	fmt.Println("App is running on PORT 3000")

	err := http.ListenAndServe(":3000", nil)

	if err != nil {
		fmt.Println("Error Starting the App", err)
	}

}
