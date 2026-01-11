package main

import(
    "fmt"
    "net/http"
    "io"
)

func SimpleCall() {
    // Calling the endpoint.
    fmt.Println("Calling the API...")
    resp, err := http.Get("http://httpbin.org/uuid")
    if err != nil {
        fmt.Println("There was an Error!")
        fmt.Println(err)
        return
    }
    fmt.Println("No errors found.")
    
    // Printing response.
    fmt.Println(resp.Status)
    var body []byte
    body, err = io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("There was an Error!")
        fmt.Println(err)
        return
    }
    fmt.Printf("Body: %s\n", body)
}

