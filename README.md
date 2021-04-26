# README

pacchetto per gestire le scadenze

### Linguaggio

- Go

### Librerie usate

- Gorm

## Utilizzo

1. In .gitconfig deve esserci in fondo una delle due opzioni di seguito (esegui cat ~/.gitconfig per verificarne la presenza)
   [url "git@github.com:"]
   insteadOf = https://github.com/
   [url "https://github"]
   insteadOf = git://github

Per aggiungere una delle opzioni eseguire per esempio:
git config --global --add url."git@github.com:".insteadOf "https://github.com/"

2. In console eseguire: export GOPRIVATE=github.com/Mind-Informatica-srl/mind-reminder

3. In console eseguire: go get github.com/Mind-Informatica-srl/mind-reminder

4. Aggiungere (embed) `mindreminder.Remind` alla model interessata. (Dopo le chiamate Create, Save, Update, Delete vengono avviati i criteri per generare nuove scadenze)

```go
type User struct{
    Id        uint
    CreatedAt time.Time
    // ...

    mindreminder.Remind
}
```

5. Per ogni model definire i criteri per generare le scadenze

```go
func (l Remind) Reminders(db *gorm.DB) (toInsert []Reminder, toDelete []Reminder, err error) {
	return
}
```

7. Avviare il servizio passando l'elenco delle model che implementano l'interfaccia Remindable (ovvero che estendo la struct Remind)
   ed il nome dell'app (serve solo per i log)

```go
    structList := []interface{}{
		User{},
        ExampleModel2{}
	}
	if err := mindreminder.StartService(structList, "[APP_NAME]"); err != nil {
		log.Fatal(err)
	}
```

8. Nelle configurazioni di avvio è possibile specificare le variabili per l'invio di una email in caso di errore bloccante

````json
    "TO_MAIL_ERROR":"email destinatario",
    "SERVER_NAME_MAIL_ERROR":"[server]:[porta]",
    "USERNAME_MAIL_ERROR":"[email del mittente]",
    "PASSWORD_MAIL_ERROR":"[password del mittente]"
                ```
````
