package site

import (
	"context"
	"net/http"
	"net/mail"
	"strings"

	uiapp "github.com/DonaldMurillo/gofastr/core-ui/app"
	"github.com/DonaldMurillo/gofastr/core-ui/component"
	"github.com/DonaldMurillo/gofastr/core-ui/html"
	"github.com/DonaldMurillo/gofastr/core/render"
	"github.com/DonaldMurillo/gofastr/core/router"
	"github.com/DonaldMurillo/gofastr/framework/headless"
	ui "github.com/DonaldMurillo/gofastr/framework/ui"
)

// ContactPath is where the contact form posts.
const ContactPath = "/contact/send"

// contactThanks is what a visitor sees after sending the form. The form
// is not wired to an inbox yet (it never was): the message is checked and
// then dropped, exactly as the old client-side handler did.
const contactThanks = "Thanks — this form is ready to connect to the Lotus inbox."

// ContactScreen is the contact page: the note on the left, the form in a
// forest card on the right. ?sent=1 and ?error=1 (set by the form's
// handler) show the outcome above the form.
type ContactScreen struct {
	component.ContextOnly
}

func (s *ContactScreen) RenderCtx(ctx context.Context) render.HTML {
	q := uiapp.QueryFromContext(ctx)

	var status render.HTML
	switch {
	case q.Get("sent") == "1":
		status = ui.Callout(ui.CalloutConfig{Title: "Message sent", Variant: ui.StatusSuccess, ID: "contact-status",
			ExtraAttrs: html.Attrs{"role": "status"}}, render.Text(contactThanks))
	case q.Get("error") == "1":
		status = ui.Callout(ui.CalloutConfig{Title: "Please check the form", Variant: ui.StatusDanger, ID: "contact-status",
			ExtraAttrs: html.Attrs{"role": "alert"}}, render.Text("Every field is required, with a valid email address."))
	}

	intro := ui.Stack(ui.StackConfig{Gap: ui.GapLG},
		ui.Hero(ui.HeroConfig{
			Eyebrow:  "Contact us",
			Title:    "A note for the Lotus team.",
			Subtitle: "Questions about the menu, dietary needs, or your next visit? Send us a message and we’ll be happy to help.",
			Footer: html.Paragraph(html.TextConfig{},
				render.Text("Need us sooner? "),
				ui.Link(ui.LinkConfig{Href: phones[0].Tel, Text: "Call " + phones[0].Display})),
		}),
	)

	form := ui.Form(ui.FormConfig{
		Action:      ContactPath,
		Method:      "POST",
		ID:          "contact-form",
		SubmitLabel: "Send message ↗",
	},
		ui.Grid(ui.GridConfig{Min: "12rem", Gap: ui.GapMD},
			ui.TextField(ui.TextFieldConfig{Name: "first_name", Label: "First name", ID: "contact-first-name", AutoComplete: "given-name", Required: true, MaxLength: 100}),
			ui.TextField(ui.TextFieldConfig{Name: "last_name", Label: "Last name", ID: "contact-last-name", AutoComplete: "family-name", Required: true, MaxLength: 100}),
		),
		ui.FormField(ui.FormFieldConfig{
			Label: "Email", For: "contact-email", Required: true,
			Input: func(fc headless.FieldControl) render.HTML {
				return ui.Control(ui.ControlConfig{
					Field: fc, Type: "email", Name: "email", AutoComplete: "email", MaxLength: 254,
					ExtraAttrs: html.Attrs{"inputmode": "email"},
				})
			},
		}),
		ui.TextArea(ui.TextAreaConfig{Name: "message", Label: "Comment", ID: "contact-message", Rows: 5, Required: true, MaxLength: 5000}),
		ui.Muted(render.Text("All fields are required.")),
	)

	// The form card is drawn in the forest band's palette; the frame
	// clips the band's square backdrop to the card's corners.
	card := html.Div(html.DivConfig{Class: ContactStyle.Frame(), ID: "contact-card"},
		ui.Themed(forestBand, ui.Card(ui.CardConfig{Heading: "Say hello", HeadingLevel: 2, ID: "contact"},
			ui.Stack(ui.StackConfig{Gap: ui.GapLG}, status, form))))

	return ui.Container(ui.ContainerConfig{Width: ui.ContainerPage, Pad: ui.ContainerPadPage},
		ui.Grid(ui.GridConfig{Min: "20rem", Gap: ui.Gap2XL}, intro, card))
}

// MountRoutes registers the site's non-page routes on the framework
// router. Call it before mounting the UI host.
func MountRoutes(r *router.Router) {
	r.Post(ContactPath, http.HandlerFunc(handleContact))
}

// handleContact checks a contact submission and answers with a 303 back
// to the page, which shows the outcome. No message is stored or sent.
func handleContact(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	if err := r.ParseForm(); err != nil {
		http.Redirect(w, r, "/contact?error=1#contact-card", http.StatusSeeOther)
		return
	}
	if !validContact(r.PostForm.Get("first_name"), r.PostForm.Get("last_name"), r.PostForm.Get("email"), r.PostForm.Get("message")) {
		http.Redirect(w, r, "/contact?error=1#contact-card", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/contact?sent=1#contact-card", http.StatusSeeOther)
}

func validContact(first, last, email, message string) bool {
	for _, v := range []string{first, last, email, message} {
		if strings.TrimSpace(v) == "" {
			return false
		}
	}
	addr, err := mail.ParseAddress(strings.TrimSpace(email))
	return err == nil && addr.Name == "" && strings.Contains(addr.Address, "@")
}
