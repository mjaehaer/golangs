package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// type roots []string

func process(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/metrics" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":
		// var roots roots
		var root string
		var buf bytes.Buffer
		file, header, err := r.FormFile("file")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		tempOutFile, _ := os.Open("tempOutFile.txt")
		defer tempOutFile.Close()
		fmt.Printf("File name %s\n", header.Filename)
		io.Copy(&buf, file)
		fmt.Printf("%s\n", &buf)
		scanner := bufio.NewScanner(&buf)
		for scanner.Scan() {
			var strbuf string
			if s1 := scanner.Text()[0]; s1 != 32 {
				root = scanner.Text()[:len(scanner.Text())-1]
			}
			if s2 := scanner.Text()[1]; s2 == 45 {

				a := strings.Split(scanner.Text(), ":")[0]
				strbuf += root + "{" + strings.Split(a, "- ")[1] + "=" + "\"" + strings.Split(scanner.Text(), ": ")[1] + "\"" + "}"
				fmt.Println(strbuf)
			} else {
				if s3 := scanner.Text()[0]; s3 == 32 {
					strbuf += " " + strings.Split(scanner.Text(), ": ")[1]
					println(strbuf)
				}
			}
		}
		// fmt.Println(root)
		buf.Reset()
		return
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}

func main() {
	http.HandleFunc("/metrics", process)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
