package main

import (
	"net/http"

	"github.com/ffss92/qrpix"
)

func main() {
	static := qrpix.NewStatic(
		"123e4567-e12b-12d1-a456-426655440000",
		"Fulano de Tal",
		"BRASILIA",
		"***",
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := static.Serve(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":8000", nil)
}
