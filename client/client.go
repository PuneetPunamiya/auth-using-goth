package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"

	"github.com/gorilla/pat"
)

func main() {
	m := make(map[string]string)
	m["github"] = "Github"
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	providerIndex := &ProviderIndex{Providers: keys, ProvidersMap: m}

	p := pat.New()
	p.Get("/auth/show", func(res http.ResponseWriter, req *http.Request) {
		fmt.Println("----------------------------")
		fmt.Println("Ayo ayo")
		fmt.Println("----------------------------")
		user := req.Body
		fmt.Println(user)
		// fmt.Println(req)
		t, _ := template.New("foo").Parse(userTemplate)
		t.Execute(res, user)
	})

	p.Get("/", func(res http.ResponseWriter, req *http.Request) {
		t, _ := template.New("foo").Parse(indexTemplate)
		t.Execute(res, providerIndex)
	})

	log.Println("listening on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", p))
}

type ProviderIndex struct {
	Providers    []string
	ProvidersMap map[string]string
}

var indexTemplate = `{{range $key,$value:=.Providers}}
    <p><a href="http://localhost:4200/auth/{{$value}}">Log in with {{index $.ProvidersMap $value}}</a></p>
{{end}}`
var userTemplate = `
<p><a href="/logout/{{.Provider}}">logout</a></p>
<p>Name: {{.Name}} [{{.LastName}}, {{.FirstName}}]</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
<p>ExpiresAt: {{.ExpiresAt}}</p>
<p>RefreshToken: {{.RefreshToken}}</p>
`
