package main

import (
	"html/template"
	"io"
	"net/http"
	"strconv"

	"github.com/shurcooL/htmlg"
	"github.com/shurcooL/reactions"
	"github.com/shurcooL/resume"
	"github.com/shurcooL/users"
)

var resumeHTML = template.Must(template.New("").Funcs(template.FuncMap{"noescape": func(s string) template.HTML { return template.HTML(s) }}).Parse(`<html>
	<head>
		<title>Dmitri Shuralyov - Resume</title>
		<link href="/blog/assets/octicons/octicons.min.css" rel="stylesheet" type="text/css">
		<link href="resume.css" rel="stylesheet" type="text/css">

		{{noescape "<!-- Unminified source is at https://github.com/shurcooL/resume. -->"}}
		<script async src="resume.js"></script>

		{{if .Production}}` + googleAnalytics + `{{end}}
	</head>
	<body>`))

const googleAnalytics = `<script>
		  (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
		  (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
		  m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
		  })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

		  ga('create', 'UA-56541369-3', 'auto');
		  ga('send', 'pageview');

		</script>`

// fileServer contains /resume.{js,css}.
func initResume(fileServer http.Handler, reactions reactions.Service, users users.Service) {
	http.Handle("/resume", errorHandler{func(w http.ResponseWriter, req *http.Request) error {
		if req.Method != "GET" {
			return MethodError{Allowed: []string{"GET"}}
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		data := struct{ Production bool }{*productionFlag}
		err := resumeHTML.Execute(w, data)
		if err != nil {
			return err
		}

		// Optional (still experimental) server-side rendering.
		if ok, _ := strconv.ParseBool(req.URL.Query().Get("prerender")); ok {
			authenticatedUser, err := users.GetAuthenticated(req.Context())
			if err != nil {
				return err
			}

			// THINK.
			resume.CurrentUser = authenticatedUser
			resume.Reactions = reactions

			returnURL := req.URL.String()
			_, err = io.WriteString(w, string(htmlg.Render(resume.Header{ReturnURL: returnURL}.Render()...)))
			if err != nil {
				return err
			}
			_, err = io.WriteString(w, string(htmlg.Render(resume.DmitriShuralyov{}.Render()...)))
			if err != nil {
				return err
			}
		}

		_, err = io.WriteString(w, `</body></html>`)
		return err
	}})
	http.Handle("/resume.js", fileServer)
	http.Handle("/resume.css", fileServer)
}
