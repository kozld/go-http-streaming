# go-http-streaming

Simple Go client for streaming data to an HTTP server via a POST multipart/form-data request

## Features

* Using I/O streams to avoid storing entire objects in memory
* Sending multiple files within one HTTP connection using the `SendUnclose()` method
* Zero-dependencies project

## Installation

```
go get github.com/kozld/go-http-streaming
```

## Example

The example below initializes the client and sends content from a file

```
c := client.New(http.MethodPost, endpoint)

file, _ := os.Open(filepath)

resp, err := c.Send(file)
if err != nil {
    log.Fatalln(err)
}

io.Copy(os.Stdout, resp.Body)
```

The following example shows sending multiple files within a single *HTTP* connection.  
**Note:** we need to close the connection ourselves by calling the `c.Close()` method to return `http.Response` or an `error`

```
var err error

c := client.New(http.MethodPost, endpoint)

file1, _ := os.Open(filepath1)
file2, _ := os.Open(filepath2)

err = c.SendUnclose(file1)
if err != nil {
    log.Fatalln(err)
}

err = c.SendUnclose(file2)
if err != nil {
    log.Fatalln(err)
}

resp, err := c.Close()
if err != nil {
    log.Fatalln(err)
}

io.Copy(os.Stdout, resp.Body)
```
