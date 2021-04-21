# README

pacchetto per gestire le scadenze

### Linguaggio

- Go

### Librerie usate

- Gorm

## Utilizzo

1. Registra il plugin usando `mindreminder.Register(db)`.

```go
plugin, err := Register(database) // database è *gorm.DB
if err != nil {
	panic(err)
}
```

2. Aggiungere (embed) `mindreminder.RemindableModel` alla model interessata.

```go
type User struct{
    Id        uint
    CreatedAt time.Time
    // ...

    mindreminder.RemindableModel
}
```

3. Per ogni model definire i criteri per generare le scadenze
4. Dopo le chiamate Create, Save, Update, Delete vengono avviati i criteri per generare nuove scadenze.
