package main

import (
	"html"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/AlexVCS/lotus-vegan/internal/menu"
)

func get(t *testing.T, srv *httptest.Server, path string) (int, string) {
	t.Helper()
	res, err := http.Get(srv.URL + path)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	return res.StatusCode, string(body)
}

// TestRoutes checks the pages and static files. /healthz is registered by
// App.Start, so it is checked against the running server, not here.
func TestRoutes(t *testing.T) {
	srv := httptest.NewServer(buildApp().Router())
	defer srv.Close()

	for path, want := range map[string]string{
		"/":                             "Food with peace at heart.",
		"/menu":                         "Find your favorite.",
		"/contact":                      "A note for the Lotus team.",
		"/assets/avocado-salad.jpg":     "",
		"/assets/dm-sans-regular.woff2": "",
		"/favicon.ico":                  "",
		"/favicon.svg":                  "<svg",
		"/lotus.js":                     "gofastr.colorScheme",
	} {
		code, body := get(t, srv, path)
		if code != http.StatusOK {
			t.Errorf("GET %s = %d, want 200", path, code)
			continue
		}
		if want != "" && !strings.Contains(body, want) {
			t.Errorf("GET %s: body missing %q", path, want)
		}
	}
}

// TestEveryMenuItemRenders holds the menu pages to the menu data: every
// item's title and price (with its $ prefix), and each kind's name and
// ingredients, is on both / and /menu.
func TestEveryMenuItemRenders(t *testing.T) {
	srv := httptest.NewServer(buildApp().Router())
	defer srv.Close()

	for _, path := range []string{"/", "/menu"} {
		_, body := get(t, srv, path)
		body = html.UnescapeString(body)
		for _, g := range menu.Data.Groups {
			if !strings.Contains(body, g.Title) {
				t.Errorf("%s: group %q missing", path, g.Title)
			}
			for _, it := range g.Items {
				if !strings.Contains(body, it.Title) {
					t.Errorf("%s: item %q missing", path, it.Title)
				}
				if it.Price != "" && !strings.Contains(body, it.Price) {
					t.Errorf("%s: price %q of %q missing", path, it.Price, it.Title)
				}
				if it.Description != "" && !strings.Contains(body, it.Description) {
					t.Errorf("%s: description of %q missing", path, it.Title)
				}
				for _, k := range it.Kinds {
					row := ">" + k.Name + "</span><span class=\"kind-ingredients\">" + k.Ingredients + "<"
					if !strings.Contains(body, row) {
						t.Errorf("%s: kind %q (%q) of %q missing", path, k.Name, k.Ingredients, it.Title)
					}
				}
			}
		}
		if !strings.Contains(body, menu.Data.AllergyNote.Text) {
			t.Errorf("%s: allergy note missing", path)
		}
	}
}

func TestContactForm(t *testing.T) {
	srv := httptest.NewServer(buildApp().Router())
	defer srv.Close()
	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}

	for name, tc := range map[string]struct {
		form url.Values
		want string
	}{
		"valid":     {url.Values{"first_name": {"Ada"}, "last_name": {"L"}, "email": {"ada@example.com"}, "message": {"Hi"}}, "/contact?sent=1"},
		"missing":   {url.Values{"first_name": {"Ada"}, "email": {"ada@example.com"}, "message": {"Hi"}}, "/contact?error=1"},
		"bad email": {url.Values{"first_name": {"Ada"}, "last_name": {"L"}, "email": {"nope"}, "message": {"Hi"}}, "/contact?error=1"},
	} {
		res, err := client.PostForm(srv.URL+"/contact/send", tc.form)
		if err != nil {
			t.Fatal(err)
		}
		res.Body.Close()
		if res.StatusCode != http.StatusSeeOther || !strings.HasPrefix(res.Header.Get("Location"), tc.want) {
			t.Errorf("%s: %d %q, want 303 %s", name, res.StatusCode, res.Header.Get("Location"), tc.want)
		}
	}

	_, body := get(t, srv, "/contact?sent=1")
	if !strings.Contains(html.UnescapeString(body), "ready to connect to the Lotus inbox") {
		t.Error("/contact?sent=1 does not show the thank-you message")
	}
}
