package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//"text/template"
)

// Compile templates on start of the application
//var templates = template.Must(template.ParseFiles("static/upload.html"))

const (
	AUTH_KEY = "97fd1e27-69cb-4a54-ad43-df4c78a851ff"
)

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

// func display(w http.ResponseWriter, page string, data interface{}) {
// 	templates.ExecuteTemplate(w, page+".html", data)
// }

func main() {
	//http.HandleFunc("/", handler)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.Handle("/data/", http.StripPrefix(
		"/data/",
		http.FileServer(http.Dir("./static")),
	))

	http.Handle("/cam/", http.StripPrefix(
		"/cam/",
		http.FileServer(http.Dir("./cam")),
	))
	http.HandleFunc("/upload", uploadHandler)
	fmt.Printf("Starting server at port 51062\n")
	log.Fatal(http.ListenAndServe(":51062", nil))
}
