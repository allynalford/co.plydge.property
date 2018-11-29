package model

import (
	"app/shared/database"
	"encoding/json"
	"log"
)

// *****************************************************************************
// Bcpa
// *****************************************************************************

// Bcpa table contains the information for each user
type Bcpa struct {
	Siteaddress    string `json:"siteaddress"`
	Owner          string `json:"owner"`
	MailingAddress string `json:"mailingAddress"`
	ID             string `json:"id"`
	Milage         string `json:"milage"`
	Use            string `json:"use"`
	Legal          string `json:"legal"`
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
