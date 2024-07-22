// server/server.go
package server

import (
    "archive/zip"
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "io"
    "log"
    "mime"
    "net"
    "net/http"
    "os"
    "path/filepath"
    "strings"
)

type HttpRequestHandler struct {
    zipFilePath string
    zipReader   *zip.ReadCloser
}

func NewHttpRequestHandler(zipFilePath string) (*HttpRequestHandler, error) {
    zipReader, err := zip.OpenReader(zipFilePath)
    if err != nil {
        return nil, err
    }
    return &HttpRequestHandler{zipFilePath: zipFilePath, zipReader: zipReader}, nil
}

func (h *HttpRequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    defer h.zipReader.Close()

    filePath := strings.TrimLeft(r.URL.Path, "/")
    if filePath == "" {
        filePath = "index.html"
    }

    var file *zip.File
    for _, f := range h.zipReader.File {
        if f.Name == filePath {
            file = f
            break
        }
    }

    if file == nil {
        http.Error(w, "File Not Found", http.StatusNotFound)
        return
    }

    fileReader, err := file.Open()
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    defer fileReader.Close()

    contentType := mime.TypeByExtension(filepath.Ext(filePath))
    if contentType == "" {
        contentType = "application/octet-stream"
    }

    w.Header().Set("Content-Type", contentType)
    w.Header().Set("Content-Length", fmt.Sprintf("%d", file.UncompressedSize64))
    w.Header().Set("Access-Control-Allow-Origin", "*")

    _, err = io.Copy(w, fileReader)
    if err != nil {
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
    }
}

func Serve(zipFilePath, ip string, port int, tlsConfig *tls.Config) {
    handler, err := NewHttpRequestHandler(zipFilePath)
    if err != nil {
        log.Fatalf("Failed to create handler: %v", err)
    }

    server := &http.Server{
        Addr:      fmt.Sprintf("%s:%d", ip, port),
        Handler:   handler,
        TLSConfig: tlsConfig,
    }

    listener, err := net.Listen("tcp", server.Addr)
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    fmt.Printf("Serving at http://%s:%d\n", ip, port)

    if tlsConfig != nil {
        err = server.ServeTLS(listener, "", "")
    } else {
        err = server.Serve(listener)
    }

    if err != nil && err != http.ErrServerClosed {
        log.Fatalf("Server error: %v", err)
    }
}
