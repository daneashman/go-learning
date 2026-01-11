package main

import(
    "fmt"
    "net/http"
    "io"
)

type handler struct {
    handlerFunc func(resW http.ResponseWriter, req *http.Request)
}
func (c handler) ServeHTTP(resW http.ResponseWriter, req *http.Request) {
    defer c.handlerFunc(resW, req)

    // Log blocking middlewear
    if req.URL.Path == "/favicon.ico" {
        return;
    }

    // Logging middlewear
    fmt.Print("# REQUEST:\n")
    fmt.Printf("%v\n\n",*req)
    fmt.Print("# HEADERS:\n")
    for h := range req.Header {
        fmt.Printf("%v: %v\n",h,req.Header[h])
    }
    fmt.Print("\n")
    body, _ := io.ReadAll(req.Body)
    fmt.Print("# BODY:\n")
    fmt.Printf("%s\n",string(body))
    fmt.Print("\n\n")
}

func main() {
    h1 := handler{handlerFunc: rootHandler}
    h2 := handler{handlerFunc: messageHandler}

    mux := http.NewServeMux()
    mux.Handle("GET /", h1)
    mux.Handle("GET /h2", h2)

    err := http.ListenAndServeTLS(":8080", "./crypto/cert.pem", "./crypto/key.pem", mux)
    if err != nil {
        fmt.Println(err)
    }
}

func rootHandler(resW http.ResponseWriter, req *http.Request) {
    http.ServeFile(resW, req, "./dist/index.html")
}

func messageHandler (resW http.ResponseWriter, req *http.Request) {
    io.WriteString(resW, "This is h2")
}
