package model

import (
	"app/shared/database"
	"time"

	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// Folio
// *****************************************************************************

// Folio table contains the information for each user
type Folio struct {
	ObjectID       bson.ObjectId `bson:"_id"`
	ID             string        `db:"folioid" bson:"id,omitempty"` // Don't use Id, use FolioID() instead for consistency with MongoDB
	SiteAddress    string        `db:"site_address" bson:"site_address"`
	MailingAddress string        `db:"mailing_address" bson:"mailing_address"`
	Owner          string        `db:"ownerinfo" bson:"owner"`
	Milage         string        `db:"mailage" bson:"mailage"`
	Use            string        `db:"useage" bson:"use"`
	Legal          string        `db:"legal" bson:"legal"`
	CreatedAt      time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt      time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted        uint8         `db:"deleted" bson:"deleted"`
}

//FolioCreate is used Create a folio propery for the user
func FolioCreate(UserID, ID, SiteAddress, MailingAddress, Owner, Milage, Use, Legal string) error {
	var err error

	now := time.Now()

	switch database.ReadConfig().Type {
	case database.TypeMySQL:
		_, err = database.SQL.Exec("INSERT INTO folio (user_id,id, site_address, mailing_address, owner, mailage, use, legal) VALUES (?,?,?,?)", UserID, ID, SiteAddress, MailingAddress, Owner, Milage, Use, Legal)
	case database.TypeBolt:
		folio := &Folio{
			ObjectID:       bson.NewObjectId(),
			ID:             ID,
			SiteAddress:    SiteAddress,
			MailingAddress: MailingAddress,
			Owner:          Owner,
			Milage:         Milage,
			Use:            Use,
			Legal:          Legal,
			CreatedAt:      now,
			UpdatedAt:      now,
			Deleted:        0,
		}

		err = database.Update("folio", folio.ID, &folio)
	default:
		err = ErrCode
	}

	return standardizeError(err)
}
