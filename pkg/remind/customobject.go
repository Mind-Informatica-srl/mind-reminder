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
	// custom event della stessa sezione dell'oggetto. Servono per esempio in AfterCreate
	CustomEvents []*CustomEvent `gorm:"-"`
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

// AfterCreate di CustomObject
func (c *CustomObject) AfterCreate(tx *gorm.DB) (err error) {
	// si ricavano i prototipi di eventi obbligatori
	// collegati alla sezione dell'oggetto creato
	var list []CustomEventPrototype
	err = tx.Table("custom_event_prototypes as cep").
		Joins("custom_section_custom_event_prototypes as cs on cs.custom_event_prototype_id = cep.id").
		Where("cs.custom_section_id = ? and cep.required_on_object_creation = true").Find(&list).Error
	if err != nil {
		return
	}
	// se mancano degli eventi obbligatori, si lancia errore
	for _, l := range list {
		founded := false
		for _, e := range c.CustomEvents {
			if e.CustomEventPrototypeID == l.ID {
				founded = true
			}
		}
		if !founded {
			return NewRemindGenericError("missing required event prototype: " + l.Name)
		}
	}
	// si inseriscono anche gli eventi collegati
	for _, e := range c.CustomEvents {
		if err = tx.Create(&e).Error; err != nil {
			return
		}
	}
	return
}
