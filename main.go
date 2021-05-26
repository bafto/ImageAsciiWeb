package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"sync"
)

var wg sync.WaitGroup //needed to keep the goroutines in sync when the server is shut down

var serverHandler *http.ServeMux //for all the request handler
var staticHandler http.Handler   //serves the static folder (javascript + css and assets)
var server http.Server           //the server itself

var templates *template.Template //the html files

func loadTemplates() {
	templates = template.Must(template.ParseFiles("index.html"))
}

func main() {
	log.SetFlags(log.Lshortfile)

	loadTemplates()

	serverHandler = http.NewServeMux()
	server = http.Server{Addr: ":3000", Handler: serverHandler} //server on Port :3000

	staticHandler = http.FileServer(http.Dir("static")) //serves the static directory for js + css and assets
	serverHandler.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	//add the html handler
	serverHandler.HandleFunc("/", indexHandler)
	serverHandler.HandleFunc("/image", imageHandler)

	log.Println("Starting cmd goroutine")
	wg.Add(1)
	//this goroutine is for the cmd interface (at the moment only the quit command for a gracefull shutdown)
	go func() {
		defer wg.Done() //tell the waiter group that we are finished at the end
		cmdInterface()
		log.Println("cmd goroutine finished")
	}()

	log.Println("server starting on Port :3000")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err.Error())
	} else if err == http.ErrServerClosed {
		log.Println("Server not listening anymore")
	}
}

func cmdInterface() {
	for loop := true; loop; {
		var inp string
		_, err := fmt.Scanln(&inp)
		if err != nil {
			log.Println(err.Error())
		} else {
			switch inp {
			case "quit":
				log.Println("Attempting to shutdown server")
				err := server.Shutdown(context.Background())
				if err != nil {
					log.Fatal("Error while trying to shutdown server: " + err.Error())
				}
				log.Println("Server was shutdown")
				loop = false
			default:
				fmt.Println("cmd not supported")
			}
		}
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	loadTemplates()
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	var (
		status int
		err    error
	)
	defer func() {
		if nil != err {
			http.Error(w, err.Error(), status)
			log.Println(err)
		}
	}()
	if err = r.ParseMultipartForm(32 << 20); nil != err {
		status = http.StatusInternalServerError
		return
	}
	fmt.Println("No memory problem")
	for _, fheaders := range r.MultipartForm.File {
		for _, hdr := range fheaders {
			fmt.Println(hdr.Filename)
			// open uploaded
			var infile multipart.File
			if infile, err = hdr.Open(); nil != err {
				status = http.StatusInternalServerError
				return
			}
			buff := bytes.NewBuffer(nil)
			if _, err = io.Copy(buff, infile); nil != err {
				status = http.StatusInternalServerError
				return
			}
			var text string
			if text, err = ImageToAscii(buff); nil != err {
				status = http.StatusInternalServerError
				return
			}
			w.Write([]byte(text))
		}
	}
}
