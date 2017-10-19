package main

import (
    "fmt"
    "io/ioutil"
    "flag"
    "log"
    "mime"
    "errors"
    "net/http"
    "os"
    "strconv"
    "strings"
)


const (
    FormatHeader = "X-Format"
    CodeHeader = "X-Code"
    ContentType = "Content-Type"
)

var port = flag.Int("port", 8080, "Port number to serve.")
var cache = make(map[string][]byte)


func main() {
    http.HandleFunc("/", handler)
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


// Expected headers:
// X-Code           indicates the HTTP code to be returned to the client
// X-Format         the value of the Accept header
// X-Original-URI
// X-Namespace
// X-Ingress-Name
// X-Service-Name
func handler(w http.ResponseWriter, r *http.Request) {
    // debug:
    log.Printf("Request headers: %v\n", r.Header)

    ext, format := getExtAndFormat(acceptHeader2contentType(r.Header.Get(FormatHeader)))
    w.Header().Set(ContentType, format)

    errCode := r.Header.Get(CodeHeader)
    code, err := strconv.Atoi(errCode)
    if err != nil {
        code = 404
        log.Printf("unexpected error reading return code: %v. Using %v\n", err, code)
    }
    w.WriteHeader(code)

    cacheKey := fmt.Sprintf("%v:%v", format, code)
    data := cache[cacheKey]
    if data == nil {
        data = getBody(ext, format, code)
        cache[cacheKey] = data
    }

    log.Printf("serving page for code %v and format %v\n", code, format)
    _, err = w.Write(data)
    if err != nil {
        log.Printf("unexpected error: %v\n", err)
    }
}


func acceptHeader2contentType(header string) string {
    return strings.TrimSpace(strings.Split(header, ",")[0])
}


func getExtAndFormat(requestedFormat string) (string, string) {
    format := requestedFormat

    if format == "" {
        format = "text/html"
        log.Printf("format not specified. Using %v\n", format)
    }

    mediaType, _, _ := mime.ParseMediaType(format)
    var ext string

    cext, err := mime.ExtensionsByType(mediaType)
    if err != nil || len(cext) == 0 {
        if err == nil {
            err = errors.New("no known extensions")
        }

        format = "text/html"
        ext = ".html"
        log.Printf("unexpected error reading media type extension: %v. Using %v\n", err, ext)
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


func getBody(ext string, format string, code int) []byte {
    path := "/www"

    file := fmt.Sprintf("%v/%v%v", path, code, ext)
    data, err := ioutil.ReadFile(file)
    if err != nil {
        log.Printf("error reading file: %v\n", err)
        scode := strconv.Itoa(code)

        file := fmt.Sprintf("%v/%cxx%v", path, scode[0], ext)
        data, err = ioutil.ReadFile(file)
        if err != nil {
            log.Printf("unexpected error reading file: %v\n", err)
            log.Printf("using fallback error response for code %v and format %v\n", code, format)
            return []byte("\"Unknown error\"")
        }
    }

    log.Printf("using custom error response for code %v and format %v from file %v\n", code, format, file)
    return data
}
