package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	pdfDir   = "./pdfs"
	webDir   = "./web"
	buildDir = "./build"
	addr     = ":8000"
)

// 模板用于列出 PDF 文件
var indexTemplate = `
<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>PDF book reader</title>
</head>
<body>
	<h1>PDF book reader</h1>
	<ul>
		{{range .}}
			<li><a href="/web/viewer.html?file=/pdfs/{{.}}">{{.}}</a></li>
		{{else}}
			<li>The PDF file was not found</li>
		{{end}}
	</ul>
</body>
</html>
`

func main() {

	http.Handle("/pdfs/", http.StripPrefix("/pdfs/", http.FileServer(http.Dir(pdfDir))))
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir(webDir))))
	http.HandleFunc("/build/", func(w http.ResponseWriter, r *http.Request) {
		filePath := filepath.Join(buildDir, strings.TrimPrefix(r.URL.Path, "/build/"))
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.Error(w, "文件未找到", http.StatusNotFound)
			return
		}
		switch {
		case strings.HasSuffix(r.URL.Path, ".mjs"):
			w.Header().Set("Content-Type", "application/javascript")
		case strings.HasSuffix(r.URL.Path, ".js"):
			w.Header().Set("Content-Type", "application/javascript")
		case strings.HasSuffix(r.URL.Path, ".css"):
			w.Header().Set("Content-Type", "text/css")
		default:
			http.Error(w, "不支持的文件类型", http.StatusBadRequest)
			return
		}
		http.ServeFile(w, r, filePath)
	})

	// 列出 pdfs 目录下的所有 .pdf 文件
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		files, err := os.ReadDir(pdfDir)
		if err != nil {
			http.Error(w, "无法读取 pdf 目录", http.StatusInternalServerError)
			return
		}

		var pdfs []string
		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".pdf") {
				pdfs = append(pdfs, file.Name())
			}
		}

		tmpl := template.Must(template.New("index").Parse(indexTemplate))
		err = tmpl.Execute(w, pdfs)
		if err != nil {
			http.Error(w, "模板渲染错误", http.StatusInternalServerError)
		}
	})

	fmt.Printf("📖 服务已启动：访问 http://localhost%s 查看 PDF 文件\n", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
