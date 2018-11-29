package controller

import (
	"app/model"
	"app/shared/database"
	"app/shared/session"
	"app/shared/view"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/josephspurrier/csrfbanana"
	"gopkg.in/mgo.v2/bson"
)

// PropertyReadGET displays the notes in the notepad
func PropertyReadGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	//userID := fmt.Sprintf("%s", sess.Values["id"])

	//notes, err := model.NotesByUserID(userID)
	//if err != nil {
	//	log.Println(err)
	//	notes = []model.Note{}
	//}

	// Display the view
	v := view.New(r)
	v.Name = "property/read"
	v.Vars["first_name"] = sess.Values["first_name"]
	//v.Vars["notes"] = notes
	v.Render(w)
}

// PropertyCreateGET displays the property creation page
func PropertyCreateGET(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	// Display the view
	v := view.New(r)
	v.Name = "property/create"
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)
}

//PropertyCreatePOST used to create a property
func PropertyCreatePOST(w http.ResponseWriter, r *http.Request) {
	// Get session
	sess := session.Instance(r)

	userID := fmt.Sprintf("%s", sess.Values["id"])

	// Validate with required fields
	if validate, missingField := view.Validate(r, []string{"house_number", "street_direction", "street_name"}); !validate {
		sess.AddFlash(view.Flash{"Field missing: " + missingField, view.FlashError})
		sess.Save(r, w)
		PropertyCreateGET(w, r)
		return
	}

	// Get form values
	houseNumber := r.FormValue("house_number")
	streetDirection := r.FormValue("street_direction")
	streetName := r.FormValue("street_name")
	situsStreetType := r.FormValue("Situs_Street_Type")
	postDirection := r.FormValue("post_direction")
	unitNumber := r.FormValue("unit_number")
	city := r.FormValue("city")

	//Build the request URL
	urlToCall := fmt.Sprintf("https://p1nng7wkbk.execute-api.us-east-1.amazonaws.com/Prod/?SN=%s&UN=%s&SD=%s&HN=%s&ST=%s&PD=%s&CT=%s", houseNumber, unitNumber, streetDirection, streetName, situsStreetType, postDirection, city)

	//Make a HTTPS call to the API service for the data
	var dataStr string
	response, err := http.Get(urlToCall)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {

		data, _ := ioutil.ReadAll(response.Body)
		dataStr = string(data)

		log.Println(dataStr)
	}

	itemInfoBytes := []byte(dataStr)
	//var bcpa model.bcpa
	bcpa := model.ParseBcpa(itemInfoBytes)

	fmt.Println(bcpa.Siteaddress)

	//Save it as a test
	err = model.BcpaCreate(userID, bcpa, database.TypeMongoDB)

	// Will only error if there is a problem with the query
	if err != nil {
		log.Println(err)
		sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
		sess.Save(r, w)
	} else {
		sess.AddFlash(view.Flash{"BCPA added!", view.FlashSuccess})
		sess.Save(r, w)
	}

	//m := map[string]interface{}{}
	//if err := json.Unmarshal([]byte(dataStr), &m); err != nil {
	//	log.Println(err)
	//}

	folio := model.Folio{
		ObjectID:       bson.NewObjectId(),
		ID:             "ID",
		SiteAddress:    "SiteAddress",
		MailingAddress: "MailingAddress",
		Owner:          "Owner",
		Milage:         "1234",
		Use:            "1234",
		Legal:          "Legal",
		Deleted:        0,
	}

	// Get database result
	//err := model.NoteCreate(content, userID)
	// Will only error if there is a problem with the query
	//if err != nil {
	//	log.Println(err)
	//	sess.AddFlash(view.Flash{"An error occurred on the server. Please try again later.", view.FlashError})
	//	sess.Save(r, w)
	//} else {
	//	sess.AddFlash(view.Flash{"Note added!", view.FlashSuccess})
	//	sess.Save(r, w)
	//	http.Redirect(w, r, "/notepad", http.StatusFound)
	//	return
	//}

	//sess.AddFlash(view.Flash{"Property: " + houseNumber, view.FlashSuccess})
	//sess.Save(r, w)

	// Display the view
	v := view.New(r)
	v.Name = "property/claim"
	v.Vars["houseNumber"] = houseNumber
	v.Vars["userID"] = userID
	v.Vars["folio"] = folio
	v.Vars["bcpa"] = bcpa
	v.Vars["token"] = csrfbanana.Token(w, r, sess)
	v.Render(w)

	// Display the same page
	//PropertyReadGET(w, r)
}
