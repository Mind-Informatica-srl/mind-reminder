# README

pacchetto per gestire le scadenze

### Linguaggio

- Go

### Librerie usate

- Gorm

## Utilizzo

1. In console eseguire: export GOPRIVATE=github.com/Mind-Informatica-srl/mind-reminder

2. In console eseguire: go get github.com/Mind-Informatica-srl/mind-reminder

3. Registra il plugin usando `mindreminder.Register(db)`.

```go
plugin, err := Register(database) // database è *gorm.DB
if err != nil {
	panic(err)
}
```

4. Aggiungere (embed) `mindreminder.RemindableModel` alla model interessata.

```go
type User struct{
    Id        uint
    CreatedAt time.Time
    // ...

    mindreminder.RemindableModel
}
```

5. Per ogni model definire i criteri per generare le scadenze
6. Dopo le chiamate Create, Save, Update, Delete vengono avviati i criteri per generare nuove scadenze.
