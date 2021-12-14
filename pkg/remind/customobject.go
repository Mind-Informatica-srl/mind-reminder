package remind

import (
	"strconv"

	"github.com/Mind-Informatica-srl/restapi/pkg/models"
	"gorm.io/gorm"
)

// CustomObject oggetto di tipo custom
type CustomObject struct {
	ID int
	// id del prototipo
	CustomObjectPrototypeID int
	CustomSectionID         int
	// jsonb con i dati dell'oggetto
	Data models.JSONB
	// prototipo dell'oggetto
	CustomObjectPrototype *CustomObjectPrototype `gorm:"association_autoupdate:false;"`
}

// SetPK set the pk for the model
func (c *CustomObject) SetPK(pk interface{}) error {
	c.ID = pk.(int)
	return nil
}

// VerifyPK check the pk value
func (c *CustomObject) VerifyPK(pk interface{}) (bool, error) {
	id := pk.(int)
	return c.ID == id, nil
}

// BeforeDelete di CustomObject
func (c *CustomObject) BeforeDelete(tx *gorm.DB) (err error) {
	if c.CustomSectionID == 0 {
		// se non abbiamo CustomSectionID si recupera da db
		var obj CustomObject
		if err = tx.First(&obj).Error; err != nil {
			return
		}
		c.CustomSectionID = obj.CustomSectionID
	}
	var evts []CustomEvent
	// si ricavano gli evento collegati
	if err = tx.Model(&evts).
		Where("custom_section_id=? and object_reference_id = ?", c.CustomSectionID, strconv.Itoa(c.ID)).
		Find(&evts).Error; err != nil {
		return
	}
	// si eliminano gli eventi collegati uno ad uno
	// in modo da avviare il ricalcolo delle scadenze
	for _, e := range evts {
		if err = tx.Delete(&e).Error; err != nil {
			return
		}
	}
	return
}

// AfterUpdate di CustomObject
func (c *CustomObject) AfterUpdate(tx *gorm.DB) (err error) {
	// si aggiornano gli eventi associati in modo da avviare il ricalcolo delle scadenze
	var evts []CustomEvent
	if err = tx.Model(&evts).
		Where("custom_section_id=? and object_reference_id = ?", c.CustomSectionID, strconv.Itoa(c.ID)).
		Find(&evts).Error; err != nil {
		return
	}
	for _, e := range evts {
		if err = tx.Updates(&e).Error; err != nil {
			return
		}
	}
	return
}
