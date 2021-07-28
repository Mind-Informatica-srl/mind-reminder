package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"

	mindlogger "github.com/Mind-Informatica-srl/mind-logger"
	"github.com/Mind-Informatica-srl/mind-reminder/internal/config"
	models "github.com/Mind-Informatica-srl/mind-reminder/pkg/models"
	"gorm.io/gorm"
)

func Paginate(r *http.Request) (paginateScope func(db *gorm.DB) *gorm.DB, page int, pageSize int) {
	var pageString, pageSizeString string
	if params, ok := r.URL.Query()["page"]; ok {
		pageString = params[0]
	}
	if params, ok := r.URL.Query()["pageSize"]; ok {
		pageSizeString = params[0]
	}
	if pageString != "" && pageSizeString != "" {
		pageSize, _ = strconv.Atoi(pageSizeString)
		page, _ = strconv.Atoi(pageString)
	}
	paginateScope = func(db *gorm.DB) *gorm.DB {
		if page > 0 && pageSize > 0 {
			offset := (page - 1) * pageSize
			return db.Offset(offset).Limit(pageSize)
		}
		return db
	}
	return
}

type pagedResultReminder struct {
	TotalCount int64             `json:"totalCount"`
	Items      []models.Reminder `json:"items"`
}

func GetAllReminders(w http.ResponseWriter, r *http.Request) {
	var log = r.Context().Value(mindlogger.LoggerContextKey).(*mindlogger.AppLogger)
	db := config.Env.GetDb(log)
	paginationScope, page, pageSize := Paginate(r)
	var list []models.Reminder
	if err := db.Scopes(paginationScope).Find(&list).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var count int64
	if page > 0 && pageSize > 0 {
		if err := db.Model(&models.Reminder{}).Count(&count).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		res := pagedResultReminder{
			TotalCount: count,
			Items:      list,
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if err := json.NewEncoder(w).Encode(list); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func UpdateReminders(db *gorm.DB, el *models.RemindToCalculate, typeRegistry map[string]reflect.Type) error {
	//si ricava object da RemindToCalculate
	obj, err := GetObjectFromRemindToCalculate(el, typeRegistry)
	if err != nil {
		return err
	}
	//si ricavano le scadenze da cancellare e quelle da inserire
	toInsertList, toDeleteList, err := obj.Reminders(db, el.Action)
	if err != nil {
		return err
	}
	//si apre transazione: se una sola insert o una sola delete ha sollevato errore, si fa rollback
	db.Transaction(func(tx2 *gorm.DB) error {
		//si cancellano ed inseriscono le scadenze ottenute
		if toDeleteList != nil && len(toDeleteList) > 0 {
			if err := tx2.Delete(toDeleteList).Error; err != nil {
				return err
			}
		}
		if toInsertList != nil && len(toInsertList) > 0 {
			if err := tx2.Create(toInsertList).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}
