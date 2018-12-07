package model

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

//OfficialRecords used to hold data from records
type OfficialRecords struct {
	ParcelID string `json:"parcelid"`
	Records  []Parcel
}

//Parcel used to hold each parcel record
type Parcel struct {
	ParcelID          string `json:"parcelid"`
	FirstDirectName   string `json:"firstdirectname"`
	FirstInDirectName string `json:"firstindirectname"`
	BookType          string `json:"booktype"`
	BookPage          string `json:"bookpage"`
	DateRecorded      string `json:"daterecorded"`
	DocType           string `json:"doctype"`
	InstrumentNumber  string `json:"instrumentnumber"`
	Legal             string `json:"legal"`
}

// Marshal Convert BCPA	to string
func Marshal(or interface{}) string {
	b, err := json.Marshal(or)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return "0"
	}
	//fmt.Println(string(b))

	return string(b)
}

// Scrape Pull Data
func Scrape() OfficialRecords {
	var err error

	// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome instance
	//c, err := chromedp.New(ctxt, chromedp.WithLog(log.Printf))
	c, err := chromedp.New(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// run task list
	var res, html string
	err = c.Run(ctxt, submit(`https://officialrecords.broward.org/AcclaimWeb/search/SearchTypeParcel`, `504209450170`, &res, &html))
	if err != nil {
		log.Fatal(err)
	}

	// shutdown chrome
	err = c.Shutdown(ctxt)
	if err != nil {
		log.Fatal(err)
	}

	// wait for chrome to finish
	err = c.Wait()
	if err != nil {
		log.Fatal(err)
	}

	p := strings.NewReader(html)
	doc, _ := goquery.NewDocumentFromReader(p)

	or := ParseRecordTable(doc)

	fmt.Println(Marshal(or))

	//file, err := os.Create("C:\\gowork\\testFiles\\parcel.txt")
	//if err != nil {
	//	log.Fatal("Cannot create file", err)
	//}
	//defer file.Close()

	//file.Write([]byte(Marshal(or)))

	return or
}

// ParseRecordTable Parse the records table of official records
func ParseRecordTable(doc *goquery.Document) OfficialRecords {

	//The Parent Struct returned
	or := OfficialRecords{}

	//Loop the Table Rows
	doc.Find("tr").Each(func(i int, tr *goquery.Selection) {
		//Init the Struct
		p := Parcel{}
		//Set the values by grabing each Table Data's contents
		p.ParcelID = tr.Find("td:nth-child(3)").Contents().Text()
		p.FirstDirectName = tr.Find("td:nth-child(4)").Contents().Text()
		p.FirstInDirectName = tr.Find("td:nth-child(5)").Contents().Text()
		p.BookType = tr.Find("td:nth-child(6)").Contents().Text()
		p.BookPage = tr.Find("td:nth-child(7)").Contents().Text()
		p.DateRecorded = tr.Find("td:nth-child(8)").Contents().Text()
		p.DocType = tr.Find("td:nth-child(9)").Contents().Text()
		p.InstrumentNumber = tr.Find("td:nth-child(10)").Contents().Text()
		p.Legal = tr.Find("td:nth-child(11)").Contents().Text()

		//Append to the array
		or.Records = append(or.Records, p)
	})

	return or
}

func submit(urlstr string, parcelid string, res *string, html *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(`/html/body/div[2]/div/div/div/div[2]/form`),
		chromedp.Sleep(2 * time.Second),
		chromedp.Submit(`#btnButton`),
		chromedp.WaitVisible(`//input[@name="ParcelId"]`),
		chromedp.Sleep(2 * time.Second),
		chromedp.SendKeys(`//input[@name="ParcelId"]`, parcelid),
		chromedp.Click(`//*[@id="btnSearch"]`),
		chromedp.Sleep(2 * time.Second),
		chromedp.WaitReady(`//*[@id="0"]`),
		chromedp.OuterHTML(`//*[@id="RsltsGrid"]/div[4]/table`, html),
		chromedp.Text(`//*[@id="RsltsGrid"]/div[4]/table/tbody/tr[1]/td[3]`, res),
	}
}
