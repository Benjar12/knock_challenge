package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Benjar12/knock_challenge/model/totable"
	"github.com/Benjar12/knock_challenge/util/httpresp"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	file, fileHeaders, err := r.FormFile("file")
	// TODO: don't panic
	if err != nil {
		httpresp.WriteJSON(w, "", 400, errors.New("No file"))
		return
	}

	tableName := r.FormValue("tablename")
	if tableName == "" {
		httpresp.WriteJSON(w, "", 400, errors.New("Missing tablename"))
		return
	}

	ext := filepath.Ext(fileHeaders.Filename)
	randID := uuid.Must(uuid.NewV4()).String()

	fPath := fmt.Sprintf("/tmp/%s%s", randID, ext)

	f, err := os.Create(fPath)
	if err != nil {
		httpresp.WriteJSON(w, "", 500, errors.New("Server failure"))
		return
	}
	defer f.Close()

	if _, err := io.Copy(f, file); err != nil {
		httpresp.WriteJSON(w, "", 500, errors.New("Server failure"))
		return
	}

	if err := f.Sync(); err != nil {
		httpresp.WriteJSON(w, "", 500, errors.New("Server failure"))
		return
	}

	err = totable.ProcessFile(tableName, fPath)
	// TODO: don't panic
	if err := f.Sync(); err != nil {
		httpresp.WriteJSON(w, "", 500, errors.New("Server railure or bad file"))
		return
	}

	httpresp.WriteJSON(w, "Processing started", 202, nil)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/load_file", handleFileUpload).Methods("POST")
	log.Fatal(http.ListenAndServe(":3302", r))
}
