package model

import (
	"app/shared/database"
	"encoding/json"
	"log"
	"time"
)

// *****************************************************************************
// Bcpa
// *****************************************************************************

// Bcpa table contains the information for each Property
type Bcpa struct {
	Siteaddress         string `json:"siteaddress"`
	Owner               string `json:"owner"`
	MailingAddress      string `json:"mailingAddress"`
	ID                  string `json:"id"`
	Milage              string `json:"milage"`
	Use                 string `json:"use"`
	Legal               string `json:"legal"`
	PropertyAssessments []PropertyAssessmentValue
	ExemptionsTaxable   ExemptionsTaxableValuesbyTaxingAuthority
	SalesHistory        []Sale
	LandCalculations    LandCalculations
	SpecialAssessments  []SpecialAssessment
}

// RecBuildingCard Card page Structure
type RecBuildingCard struct {
	CardURL                   string `json:"cardurl"`
	TaxYear                   string `json:"taxyear"`
	Folio                     string `json:"folio"`
	ParcelIDNumber            string `json:"parcelidnumber"`
	UseCode                   string `json:"usecode"`
	NoBedrooms                string `json:"nobedrooms"`
	NoBaths                   string `json:"nobaths"`
	NoUnits                   string `json:"nounits"`
	NoStories                 string `json:"nostories"`
	NoBuildings               string `json:"nobuildings"`
	Foundation                string `json:"foundation"`
	Exterior                  string `json:"exterior"`
	RoofType                  string `json:"rooftype"`
	RoofMaterial              string `json:"roofmaterial"`
	Interior                  string `json:"interior"`
	Floors                    string `json:"floors"`
	Plumbing                  string `json:"plumbing"`
	Electric                  string `json:"electric"`
	Classification            string `json:"classification"`
	CeilingHeights            string `json:"ceilingheights"`
	QualityOfConstruction     string `json:"qualityofconstruction"`
	CurrentConditionStructure string `json:"currentconditionstructure"`
	ConstructionClass         string `json:"constructionclass"`
	Permits                   []Permit
	ExtraFeatures             []ExtraFeature
}

// ExtraFeature card page features
type ExtraFeature struct {
	Feature string `json:"feature"`
}

// Permit card page permits
type Permit struct {
	PermitNo   string `json:"permitco"`
	PermitType string `json:"permittype"`
	EstCost    string `json:"estcost"`
	PermitDate string `json:"permitdate"`
	CODate     string `json:"codate"`
}

// LandCalculations land calculations structure
type LandCalculations struct {
	Calculations    []LandCalculation
	AdjBldgSF       string `json:"adjbldgsf"`
	Units           string `json:"units"`
	Cards           []RecBuildingCard
	SketchURL       string `json:"sketchurl"`
	EffActYearBuilt string `json:"effactyearbuilt"`
}

// LandCalculation land calculation structure
type LandCalculation struct {
	Price  string `json:"price"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

// SpecialAssessment Special Assessment structure
type SpecialAssessment struct {
	Fire  string `json:"fire"`
	Garb  string `json:"garb"`
	Light string `json:"light"`
	Drain string `json:"drain"`
	Impr  string `json:"impr"`
	Safe  string `json:"safe"`
	Storm string `json:"storm"`
	Clean string `json:"clean"`
	Misc  string `json:"misc"`
}

// RecPatriotSketch Sketch structure
type RecPatriotSketch struct {
	Sketch       string `json:"sketch"`
	Building     string `json:"building"`
	URL          string `json:"url"`
	SketchImgURL string `json:"sketchimgurl"`
	Codes        []PatriotSketchCode
	AdjAreaTotal string `json:"adjareatotal"`
}

// PatriotSketchCode Sketch Code structure
type PatriotSketchCode struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Area        string `json:"area"`
	Factor      string `json:"factor"`
	AdjArea     string `json:"adjarea"`
	Stories     string `json:"stories"`
}

// ExemptionsTaxableValuesbyTaxingAuthority table contains the exemptions
type ExemptionsTaxableValuesbyTaxingAuthority struct {
	County      ExemptionsAndTaxableValue
	SchoolBoard ExemptionsAndTaxableValue
	Municipal   ExemptionsAndTaxableValue
	Independent ExemptionsAndTaxableValue
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
}

// ExemptionsAndTaxableValue table contains the exemption values
type ExemptionsAndTaxableValue struct {
	JustValue    string `json:"justvalue"`
	Portability  string `json:"portability"`
	AssessedSOH  string `json:"assessedsoh"`
	Homestead    string `json:"homestead"`
	AddHomestead string `json:"addhomestead"`
	WidVetDis    string `json:"widvetdis"`
	Senior       string `json:"senior"`
	XemptType    string `json:"xempttype"`
	Taxable      string `json:"taxable"`
}

// PropertyAssessmentValue table contains the house values
type PropertyAssessmentValue struct {
	Year                string    `json:"year"`
	Land                string    `json:"land"`
	BuildingImprovement string    `json:"buildingimprovement"`
	JustMarketValue     string    `json:"justmarketvalue"`
	AssessedSOHValue    string    `json:"assessedsohvalue"`
	Tax                 string    `json:"tax"`
	CreatedAt           time.Time `json:"createdat"`
	UpdatedAt           time.Time `json:"updatedat"`
}

// Sale property sales
type Sale struct {
	Date        string `json:"date"`
	Type        string `json:"type"`
	Price       string `json:"price"`
	BookPageCIN string `json:"bookpagecin"`
}

// BcpaCreate creates a BCPA record
func BcpaCreate(userID string, b Bcpa, dbType database.Type) error {
	var err error

	//now := time.Now()
	//swap to the database we are using
	database.SwapConnect(dbType)

	switch dbType {
	case database.TypeMySQL:
		//_, err = database.SQL.Exec("INSERT INTO note (content, user_id) VALUES (?,?)", content, userID)
	case database.TypeMongoDB:
		if database.CheckConnection() {
			// Create a copy of mongo
			session := database.Mongo.Copy()
			defer session.Close()
			c := session.DB(database.ReadConfig().MongoDB.Database).C("bcpa")

			err = c.Insert(b)

		} else {
			err = ErrUnavailable
		}
	case database.TypeBolt:

		err = database.Update("bcpa", userID+b.ID, &b)
	default:
		err = ErrCode
	}

	//swap back to default connect as per config
	database.DefaultConnect()

	return standardizeError(err)
}

//ParseBcpa This functions accepts a byte array containing a JSON
func ParseBcpa(jsonBuffer []byte) Bcpa {

	// We create an empty array
	var b Bcpa

	// Unmarshal the json into it. this will use the struct tag
	err := json.Unmarshal(jsonBuffer, &b)
	if err != nil {
		log.Println(err)
		return b
	}

	// the array is now filled
	return b

}
