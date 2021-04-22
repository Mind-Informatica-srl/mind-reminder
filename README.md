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

4. Registra il plugin usando `mindreminder.Register(db)`.

```go
plugin, err := Register(database) // database è *gorm.DB
if err != nil {
	panic(err)
}
```

5. Aggiungere (embed) `mindreminder.Remind` alla model interessata.

```go
type User struct{
    Id        uint
    CreatedAt time.Time
    // ...

    mindreminder.Remind
}
```

6. Per ogni model definire i criteri per generare le scadenze

```go
func (l Remind) GetReminders(db *gorm.DB) (toInsert []ToRemind, toDelete []ToRemind, err error) {
	toInsert = []ToRemind{}
	toDelete = []ToRemind{}
	err = nil
	return
}
```

7. Dopo le chiamate Create, Save, Update, Delete vengono avviati i criteri per generare nuove scadenze.
