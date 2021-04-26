package v1

import (
	"reflect"
	"time"

	mindlogger "github.com/Mind-Informatica-srl/mind-logger"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/api/v1/controllers"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/config"
	"gorm.io/gorm"
)

func StartService(structList []interface{}) error {
	mindlogger.CreateLogFolder()
	appLog := mindlogger.CreateLogger()
	timeStart := time.Now()
	appLog.AppendLn("SERVICE START: " + timeStart.Format("01-02-2006"))

	RegisterTypes(structList)
	if err := RicalcolaScadenze(appLog); err != nil {
		appLog.Prepend(err.Error())
		appLog.PrependLn("ERROR")
		timeEnd := time.Now()
		appLog.AppendLn("SERVICE END: " + timeEnd.Format("01-02-2006"))
		//si invia email per comunicare errore
		//msg := appLog.String()
		//sendMail("destinatario@mail.com", msg, "MIND REMINDER ERROR")
		appLog.WriteLog()
		return err
	}
	timeEnd := time.Now()
	appLog.AppendLn("SERVICE END: " + timeEnd.Format("01-02-2006"))
	appLog.WriteLog()
	return nil
}

//avvia il ricalcolo delle scadenze a partire dalla tabella remind_to_calculate per le righe in cui la data di lavorazione è null
func RicalcolaScadenze(appLog *mindlogger.AppLogger) error {
	if err := config.Env.GetDb(appLog).Transaction(func(db *gorm.DB) error {
		//si ricava l'orario attuale (servirà poi per l'update su remind_to_calculate)
		timeStart := time.Now()
		//si ricavano le righe da remind_to_calculate che non sono ancora state lavorate
		toCalculateList, err := controllers.GetRemindToCalculate(db)
		if err != nil {
			return err
		}
		//si lavora ogni rigo della lista
		for _, toCalculate := range toCalculateList {
			var errorString *string
			if err := controllers.UpdateReminders(db, &toCalculate, typeRegistry); err != nil {
				//si scrive l'eventuale errore in errorString
				e := err.Error()
				errorString = &e
			}
			//si aggiornano tutte le righe di remind_to_calculate che hanno stesso object_id, object_type, created_at precedente a timeStart e elaborated_at null
			if err := controllers.UpdateCorrelatedRemindToCalculate(db, &toCalculate, timeStart, errorString); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

//mappa che associa stringa a struct
//ogni struct che ha una scadenza, va aggiunta all'interno di questa mappa
//serve poi per recuperare la struct a partire da una stringa
var typeRegistry = make(map[string]reflect.Type)

// typeRegistry["MyStruct"] = reflect.TypeOf(MyStruct{})
func RegisterTypes(myTypes []interface{}) {
	for _, v := range myTypes {
		typeRegistry[reflect.TypeOf(v).Name()] = reflect.TypeOf(v)
	}
}

/*
func sendMail(mailto string, body string, subj string) error {
	from := mail.Address{
		Address: mrs.MittenteEmail.Username,
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
	servername := mrs.MittenteEmail.ServerName //"smtps.aruba.it:465"

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

	auth := smtp.PlainAuth("", mrs.MittenteEmail.Username, mrs.MittenteEmail.Password, host)
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
*/
