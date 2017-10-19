package main

import (
	"fmt"
	"io"
	"flag"
	"log"
	"mime"
	"net/http"
	"os"
	"strconv"
)


const (
	FormatHeader = "X-Format"
	CodeHeader = "X-Code"
	ContentType = "Content-Type"
)

var port = flag.Int("port", 8080, "Port number to serve.")


func main() {
	http.HandleFunc("/", errorHandler("/www"))
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	})

	err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not start http server: %s\n", err)
		os.Exit(1)
	}
}


func errorHandler(path string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ext, format := getExtAndFormat(r)
		w.Header().Set(ContentType, format)

		errCode := r.Header.Get(CodeHeader)
		code, err := strconv.Atoi(errCode)
		if err != nil {
			code = 404
			log.Printf("unexpected error reading return code: %v. Using %v\n", err, code)
		}
		w.WriteHeader(code)

		file := fmt.Sprintf("%v/%v%v", path, code, ext)
		f, err := os.Open(file)
		if err != nil {
			log.Printf("error opening file: %v\n", err)
			scode := strconv.Itoa(code)
			
			file := fmt.Sprintf("%v/%cxx%v", path, scode[0], ext)
			f, err = os.Open(file)
			if err != nil {
				log.Printf("unexpected error opening file: %v\n", err)
				fmt.Fprint(w, "Unknown error")
				return
			}
		}

		defer f.Close()
		log.Printf("serving custom error response for code %v and format %v from file %v\n", code, format, file)
		io.Copy(w, f)
	}
}


func getExtAndFormat(r *http.Request) (string, string) {
	format := r.Header.Get(FormatHeader)
	if format == "" {
		format = "text/html"
		log.Printf("format not specified. Using %v\n", format)
	}

	mediaType, _, _ := mime.ParseMediaType(format)
	var ext string

	cext, err := mime.ExtensionsByType(mediaType)
	if err != nil {
		log.Printf("unexpected error reading media type extension: %v. Using %v\n", err, ext)
		format = "text/html"
		ext = "html"
	} else {
		ext = getLongest(cext)
	}
	
	return ext, format
}


func getLongest(vs []string) string {
	longest := vs[0]

    for _, v := range vs {
        if len(v) > len(longest) {
            longest = v
        }
    }
    
    return longest
}
