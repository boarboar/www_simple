package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	//"text/template"
	"time"
)

type CAMInfo struct {
	Id      int    `json:"id"`
	Path    string `json:"path"`
	Updated string `json:"updated"`
}

// type MainHandler struct {
// 	loc *time.Location
// }

type tplParams struct {
	LastDate string
}

var loc *time.Location

// Compile templates on start of the application
//var mainTmpl = template.Must(template.ParseFiles("static/index.html"))

const (
	AUTH_KEY = "97fd1e27-69cb-4a54-ad43-df4c78a851ff"
	DEF_TZ   = "Europe/Moscow"
)

func getFileModdate(path string) string {
	res := "--:--"
	if fileInfo, err := os.Lstat("./cam/CAM-1.jpg"); err == nil {
		if loc != nil {
			res = fileInfo.ModTime().In(loc).Format("2006-01-02 15:04:05")
		}
	} else {
		log.Println(err)
	}
	return res
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	// case "GET":
	// 	display(w, "upload", nil)
	case "POST":
		uploadFile(w, r)
	}
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	log.Println("upload request")
	key := r.Header.Get("X-API-Key")
	log.Println("Auth:", key)
	if key != AUTH_KEY {
		log.Println("Bad key")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("imageFile")
	if err != nil {
		log.Println("Error Retrieving the File")
		log.Println(err)
		return
	}
	defer file.Close()
	log.Printf("Uploaded File: %+v\n", handler.Filename)
	log.Printf("File Size: %+v\n", handler.Size)
	log.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	//tempFile, err := ioutil.TempFile("uploads", "upload-*.png")
	tempFile, err := os.Create("./cam/" + handler.Filename)
	if err != nil {
		log.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File to %s\n", tempFile.Name())
	log.Println("Successfully Uploaded File to ", tempFile.Name())
}

// func mainHandler(w http.ResponseWriter, r *http.Request) {
// 	params := tplParams{
// 		LastDate: getFileModdate("/cam/CAM-1.jpg"),
// 	}

// 	mainTmpl.Execute(w, params)
// }

func apiHandler(w http.ResponseWriter, r *http.Request) {
	cam := CAMInfo{
		Id:      1,
		Path:    "/cam/CAM-1.jpg",
		Updated: getFileModdate("/cam/CAM-1.jpg"),
	}
	json.NewEncoder(w).Encode(cam)
}

// func (h *MainHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	params := tplParams{
// 		LastDate: "--:--",
// 	}

// 	if fileInfo, err := os.Lstat("./cam/CAM-1.jpg"); err == nil {
// 		if h.loc != nil {
// 			params.LastDate = fileInfo.ModTime().In(h.loc).Format("2006-01-02 15:04:05")
// 		}
// 	} else {
// 		log.Println(err)
// 	}

// 	mainTmpl.Execute(w, params)
// }

func main() {

	loc, _ = time.LoadLocation(DEF_TZ)
	//mainHandler := &MainHandler{loc: loc}
	//http.Handle("/", mainHandler)

	//http.HandleFunc("/", mainHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/api", apiHandler)

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.Handle("/data/", http.StripPrefix(
		"/data/",
		http.FileServer(http.Dir("./static")),
	))

	http.Handle("/cam/", http.StripPrefix(
		"/cam/",
		http.FileServer(http.Dir("./cam")),
	))

	fmt.Printf("Starting server at port 51062\n")
	log.Fatal(http.ListenAndServe(":51062", nil))
}
