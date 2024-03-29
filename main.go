package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func process(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/metrics" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	switch r.Method {
	case "POST":
		var root string
		var buf bytes.Buffer
		file, header, err := r.FormFile("file")
		if err != nil {
			panic(err)
		}
		defer file.Close()
		tempOutFile, _ := os.Create("tempOutFile.openmetrics")
		// defer tempOutFile.Close()
		fmt.Printf("File name %s\n", header.Filename)
		io.Copy(&buf, file)
		fmt.Printf("%s\n", &buf)
		scanner := bufio.NewScanner(&buf)
		var strbuf string
		for scanner.Scan() {

			if s1 := scanner.Text()[0]; s1 != 32 {
				root = scanner.Text()[:len(scanner.Text())-1]
			}
			if s2 := scanner.Text()[1]; s2 == 45 {

				a := strings.Split(scanner.Text(), ":")[0]
				strbuf += root + "{" + strings.Split(a, "- ")[1] + "=" + "\"" + strings.Split(scanner.Text(), ": ")[1] + "\"" + "}"
			}
			if s3, s4 := scanner.Text()[0], scanner.Text()[1]; s3 == 32 && s4 != 45 {
				strbuf += " " + strings.Split(scanner.Text(), ": ")[1] + "\n"
				println(strbuf)
				if _, err := tempOutFile.Write([]byte(strbuf)); err != nil {
					fmt.Println(err)
				}
				strbuf = ""
			}
		}

		tempOutFile.Close()
		renderFile(w, "tempOutFile.openmetrics")
		buf.Reset()
		return
	default:
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
	}
}
// comment

func renderFile(w http.ResponseWriter, filename string) {
	fmt.Println("Read request: " + filename)
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Cann't open file: " + filename)
	} else {
		w.Write(file)
	}
}

func main() {
	http.HandleFunc("/metrics", process)

	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println(err)
	}
}
