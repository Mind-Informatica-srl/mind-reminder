package logic

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"os"
	"reflect"
	"time"

	"github.com/Mind-Informatica-srl/mind-reminder/internal/config"
	mrmodel "github.com/Mind-Informatica-srl/mind-reminder/pkg/mrnodel"
	"gorm.io/gorm"
)

type B struct{}

func (b B) Doh() {

}

const (
	ToEmailError       = "TO_MAIL_ERROR"
	UsernameEmailError = "USERNAME_MAIL_ERROR"
	PasswordMailError  = "PASSWORD_MAIL_ERROR"
	ServerMailError    = "SERVER_NAME_MAIL_ERROR"
)

//mappa che associa stringa a struct
//ogni struct che ha una scadenza, va aggiunta all'interno di questa mappa
//serve poi per recuperare la struct a partire da una stringa
var typeRegistry = make(map[string]reflect.Type)

// typeRegistry["MyStruct"] = reflect.TypeOf(MyStruct{})
func RegisterTypes(myTypes []interface{}) {
	for _, v := range myTypes {
		structType := reflect.TypeOf(v)
		typeRegistry[structType.Name()] = structType
	}
}

//si aggiornano tutte le righe di remind_to_calculate che hanno stesso object_id, object_type, created_at precedente a timeStart e elaborated_at null
func updateCorrelatedRemindToCalculate(db *gorm.DB, toCalculate *RemindToCalculate, timeStart time.Time, errorString *string) error {
	if err := db.Model(RemindToCalculate{}).Scopes(filterNotElaborated(toCalculate.ObjectID, toCalculate.ObjectType)).Where("created_at < ?", timeStart).Updates(RemindToCalculate{ElaboratedAt: &timeStart, Error: errorString}).Error; err != nil {
		return err
	}
	return nil
}

//avvia il ricalcolo delle scadenze a partire dalla tabella remind_to_calculate per le righe in cui la data di lavorazione è null
func RicalcolaScadenze() error {
	db := config.Current().DB
	if err := db.Transaction(func(db *gorm.DB) error {
		//si ricava l'orario attuale (servirà poi per l'update su remind_to_calculate)
		timeStart := time.Now()
		//si ricavano le righe da remind_to_calculate che non sono ancora state lavorate
		toCalculateList, err := getRemindToCalculate(db)
		if err != nil {
			return err
		}
		//si lavora ogni rigo della lista
		for _, toCalculate := range toCalculateList {
			var errorString *string
			if err := updateReminders(db, &toCalculate, typeRegistry); err != nil {
				//si scrive l'eventuale errore in errorString
				e := err.Error()
				errorString = &e
			}
			//si aggiornano tutte le righe di remind_to_calculate che hanno stesso object_id, object_type, created_at precedente a timeStart e elaborated_at null
			if err := updateCorrelatedRemindToCalculate(db, &toCalculate, timeStart, errorString); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

//resituisce righe da lavorare (elaborated_at is null) in ordine di data creazione decrescente (in modo da avere le struct in object_raw più aggiornate)
//per ogni coppia object_id, object_type restituisce solo la riga inserita più recentemente
//si trascurano quindi le righe più vecchie (verranno comunque aggiornate dal servizio)
//(se abbiamo due righe non ancora lavorate con stesso object_id e object_type non ha infatti senso eseguire due volte il ricalcolo delle scadenze.
//Faremo il ricalcolo prendendo solo il rigo con created_at più recente)
func getRemindToCalculate(db *gorm.DB) ([]RemindToCalculate, error) {
	var list []RemindToCalculate
	if err := db.Select("object_id, object_type").Group("object_id, object_type").Where("elaborated_at is null").Find(&list).Error; err != nil {
		return nil, err
	}
	for i, v := range list {
		if err := db.Order("created_at desc").Scopes(filterNotElaborated(v.ObjectID, v.ObjectType)).First(&v).Error; err != nil {
			return nil, err
		}
		list[i] = v
	}
	return list, nil
}

func filterNotElaborated(objectID string, objectType string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("object_id = ? and object_type = ? and elaborated_at is null", objectID, objectType)
	}
}

func updateReminders(db *gorm.DB, el *RemindToCalculate, typeRegistry map[string]reflect.Type) error {
	//si ricava object da RemindToCalculate
	obj, err := getObjectFromRemindToCalculate(el, typeRegistry)
	if err != nil {
		return err
	}
	//si ricavano le scadenze da cancellare e quelle da inserire
	reminders, err := obj.Reminders(db)
	if err != nil {
		return err
	}
	//si apre transazione: se una sola insert o una sola delete ha sollevato errore, si fa rollback
	db.Transaction(func(tx2 *gorm.DB) error {
		for _, reminder := range reminders {
			if err := reminder.ModifyReminds(tx2, el.Action); err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

//restituisce ObjectRaw (json) sotto forma di struct sfruttando typeRegistry per ricavare il tipo di struct da ObjectType
func getObjectFromRemindToCalculate(el *RemindToCalculate, typeRegistry map[string]reflect.Type) (mrmodel.Event, error) {
	// si ricava il tipo di struct da ObjectType
	if t, ok := typeRegistry[el.ObjectType]; ok {
		//si converte ObjectRaw (json) in struct e si mette dentro Object di el
		if event, err := el.Event(t); err != nil {
			return nil, err
		} else {
			return event, nil
		}
	} else {
		return nil, errors.New("Missing ObjectType " + el.ObjectType + " in typeRegistry")
	}
}

func sendMail(mailto string, body string, subj string) error {
	mailFrom := os.Getenv(UsernameEmailError)
	from := mail.Address{
		Address: mailFrom,
	}
	to := mail.Address{
		Address: mailto,
	}
	// Setup headers
	headers := make(map[string]string)
	headers["From"] = from.String()
	headers["To"] = to.String()
	headers["Subject"] = subj
	headers["MIME-version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\";"
	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Connect to the SMTP Server
	servername := os.Getenv(ServerMailError)

	host, _, err := net.SplitHostPort(servername)
	if err != nil {
		return err
	}

	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", mailFrom, os.Getenv(PasswordMailError), host)
	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(from.Address); err != nil {
		return err
	}

	if err = c.Rcpt(to.Address); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()
	return nil
}

// Writes new reminder row to db.
func AddRecordRemindToCalculate(element interface{}, objectID string, action mrmodel.Action, db *gorm.DB) error {
	r, err := NewRemindToCalculate(element, objectID, action)
	if err != nil {
		return nil
	}
	return db.Model(&r).Create(&r).Error
}
