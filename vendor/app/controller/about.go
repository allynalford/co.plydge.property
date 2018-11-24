package controller

import (
	"net/http"
	"os"

	"app/shared/view"
)

// AboutGET displays the About page
func AboutGET(w http.ResponseWriter, r *http.Request) {
	// Display the view
	v := view.New(r)
	v.Name = "about/about"
	v.Vars["commit_sha"] = os.Getenv("COMMIT_SHA")
	v.Render(w)
}
