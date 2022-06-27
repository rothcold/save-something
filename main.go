package main

import (
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var path string
var listen string

func init() {
	tempDir := os.TempDir()
	flag.StringVar(&path, "path", tempDir, "file save path, default is os temp dir (typically /tmp in *nix)")
	flag.StringVar(&listen, "listen", ":8080", "listen address, default is :8080")
}

type SaveHandler struct{}

func (SaveHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	bytes, _ := ioutil.ReadAll(request.Body)
	filename := strings.ReplaceAll(request.RequestURI, "/", " ")

	fullFilePath := fmt.Sprintf("%s/%s/%s", path, request.Method, filename)
	dir := filepath.Dir(fullFilePath)
	err := os.MkdirAll(dir, fs.ModePerm)
	if err != nil {
		response.Write([]byte(err.Error()))
		return
	}
	err = ioutil.WriteFile(fullFilePath, bytes, fs.ModePerm)
	if err != nil {
		response.Write([]byte(err.Error()))
		return
	}
	response.Write([]byte("Saved"))
}

func main() {
	flag.Parse()
	log.Printf("Save to %s", path)
	log.Printf("Server listening on %s", listen)
	s := &http.Server{
		Addr:           listen,
		Handler:        SaveHandler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
