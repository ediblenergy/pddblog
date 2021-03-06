package main

import (
	"bufio"
	"bytes"
	"github.com/gorilla/mux"
	"github.com/knieriem/markdown"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type BlogEngine struct {
    wrapper_filepath string
	_wrapper *template.Template
}

func (b *BlogEngine) _build_wrapper() *template.Template {
	wrapper, err := template.New("wrapper").Parse(fileToString(b.wrapper_filepath))
	if err != nil {
		log.Fatal(err)
	}
    return wrapper
}

func (b *BlogEngine) wrapper() *template.Template {
    if b._wrapper != nil {
        return b._wrapper
    }
    b._wrapper = b._build_wrapper()
    return b._wrapper
}

func (b *BlogEngine) HomeHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	file,err := os.Open(entry_filename(params["entry_id"]))
    if err != nil {
        if os.IsNotExist( err ) {
            http.NotFound(w,r)
            return
        }
        log.Fatal( err )
    }
	defer file.Close()
	var opt markdown.Extensions
	var buff bytes.Buffer
	p := markdown.NewParser(&opt)
	writer := bufio.NewWriter(&buff)
	p.Markdown(file, markdown.ToHTML(writer))
	writer.Flush()
	b.render(w, buff.String())
}

type blogEntry struct {
	Content template.HTML
}

func (b *BlogEngine) render(out io.Writer, content string) {
	entry := blogEntry{Content: template.HTML(content)}
	b.wrapper().Execute(out, entry)
}

func fileToString(filename string) string {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}
	return string(bytes)
}

func entry_filename(title string) string {
	return "entries/" + title + ".html"
}


func main() {
	b := BlogEngine{wrapper_filepath:"root/templates/html/wrapper.html"}
	r := mux.NewRouter()
	r.HandleFunc("/blog/entry/{title}-{entry_id:[0-9]+}", b.HomeHandler)
	http.Handle("/", r)
    log.Println("listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
