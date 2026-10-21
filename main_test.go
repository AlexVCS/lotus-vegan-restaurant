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
		"/contact":                      "Talk to the Lotus team.",
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

// TestContactPage holds the contact page to what it can deliver: there is
// no online form (nothing would receive the message), so the page offers
// the pickup phones and the old form endpoint no longer answers a POST
// with a success redirect.
func TestContactPage(t *testing.T) {
	srv := httptest.NewServer(buildApp().Router())
	defer srv.Close()

	code, body := get(t, srv, "/contact")
	if code != http.StatusOK {
		t.Fatalf("GET /contact = %d, want 200", code)
	}
	body = html.UnescapeString(body)
	for _, want := range []string{
		`href="tel:+17275491688"`,
		`href="tel:+17275440103"`,
		"We can’t take messages online just yet.",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("/contact: body missing %q", want)
		}
	}
	for _, unwanted := range []string{"<form", "/contact/send", "Message sent"} {
		if strings.Contains(body, unwanted) {
			t.Errorf("/contact: body still has %q", unwanted)
		}
	}

	client := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	res, err := client.PostForm(srv.URL+"/contact/send", url.Values{"first_name": {"Ada"}, "last_name": {"L"}, "email": {"ada@example.com"}, "message": {"Hi"}})
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()
	if res.StatusCode < 400 {
		t.Errorf("POST /contact/send = %d, want an error status", res.StatusCode)
	}
}
