package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path"
	"time"

	"dmitri.shuralyov.com/service/change"
	"github.com/shurcooL/home/component"
	"github.com/shurcooL/home/exp/vec"
	"github.com/shurcooL/home/exp/vec/attr"
	"github.com/shurcooL/home/exp/vec/elem"
	"github.com/shurcooL/htmlg"
	"github.com/shurcooL/httperror"
	"github.com/shurcooL/issues"
	"github.com/shurcooL/notifications"
	"github.com/shurcooL/users"
	"golang.org/x/net/html"
)

// packageHandler is a handler for a Go package index page.
type packageHandler struct {
	Repo repoInfo
	Pkg  pkgInfo

	issues        issueCounter
	change        changeCounter
	notifications notifications.Service
	users         users.Service
}

var packageHTML = template.Must(template.New("").Parse(`<html>
	<head>
		<title>{{.Title}}</title>
		<link href="/icon.png" rel="icon" type="image/png">
		<meta name="viewport" content="width=device-width">
		<link href="/assets/fonts/fonts.css" rel="stylesheet" type="text/css">
		<link href="/assets/package/style.css" rel="stylesheet" type="text/css">
		{{if .Production}}` + googleAnalytics + `{{end}}
	</head>
	<body>`))

func (h *packageHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) error {
	if req.Method != "GET" {
		return httperror.Method{Allowed: []string{"GET"}}
	}

	t0 := time.Now()
	openIssues, err := h.issues.Count(req.Context(), issues.RepoSpec{URI: h.Repo.Spec}, issues.IssueListOptions{State: issues.StateFilter(issues.OpenState)})
	if err != nil {
		return err
	}
	openChanges, err := h.change.Count(req.Context(), h.Repo.Spec, change.ListOptions{Filter: change.FilterOpen})
	if err != nil {
		return err
	}
	fmt.Println("counting open issues & changes took:", time.Since(t0).Nanoseconds(), "for:", h.Repo.Spec)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var title string
	if h.Pkg.Name == "main" {
		title = "Command " + path.Base(h.Pkg.Spec)
	} else {
		title = "Package " + h.Pkg.Name
	}
	err = packageHTML.Execute(w, struct {
		Production bool
		Title      string
	}{
		Production: *productionFlag,
		Title:      title,
	})
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, `<div style="max-width: 800px; margin: 0 auto 100px auto;">`)
	if err != nil {
		return err
	}

	authenticatedUser, err := h.users.GetAuthenticated(req.Context())
	if err != nil {
		log.Println(err)
		authenticatedUser = users.User{} // THINK: Should it be a fatal error or not? What about on frontend vs backend?
	}
	var nc uint64
	if authenticatedUser.ID != 0 {
		nc, err = h.notifications.Count(req.Context(), nil)
		if err != nil {
			return err
		}
	}

	// Render the header.
	header := component.Header{
		CurrentUser:       authenticatedUser,
		NotificationCount: nc,
		ReturnURL:         req.RequestURI,
	}
	err = htmlg.RenderComponents(w, header)
	if err != nil {
		return err
	}

	err = html.Render(w, htmlg.H2(htmlg.Text(h.Repo.Spec+"/...")))
	if err != nil {
		return err
	}

	// Render the tabnav.
	err = htmlg.RenderComponents(w, repositoryTabnav(noTab, h.Repo, openIssues, openChanges))
	if err != nil {
		return err
	}

	err = vec.RenderHTML(w,
		elem.H1(title),
		elem.P(elem.Code(fmt.Sprintf(`import "%s"`, h.Pkg.Spec))),
	)
	if err != nil {
		return err
	}
	if h.Pkg.DocHTML != "" {
		err = vec.RenderHTML(w, elem.H3("Overview"), vec.UnsafeHTML(h.Pkg.DocHTML))
		if err != nil {
			return err
		}
	}
	err = vec.RenderHTML(w,
		elem.H3("Installation"),
		elem.P(elem.Pre("go get -u "+h.Pkg.Spec)),
		elem.H3(elem.A("Documentation", attr.Href("https://godoc.org/"+h.Pkg.Spec))),
		elem.H3(elem.A("Code", attr.Href("https://gotools.org/"+h.Pkg.Spec))),
		elem.H3(elem.A("License", attr.Href(h.Pkg.LicenseURL))),
	)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, `</div>`)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, `</body></html>`)
	return err
}
