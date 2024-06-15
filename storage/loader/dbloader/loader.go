package dbloader

import (
	"reflect"
	"time"

	"gorm.io/gorm"

	"geniusrabbit.dev/adcorelib/storage/loader"
)

type merger interface {
	Merge(interface{})
}

type targeter interface {
	Target() interface{}
}

type dbQueryPrepare interface {
	PrepareQuery(db *gorm.DB) *gorm.DB
}

// SelectLoader returns new DB loader for select query
func SelectLoader(db *gorm.DB, query string, args ...interface{}) loader.LoaderFnk {
	return func(objectTarget interface{}, lastUpdate *time.Time) error {
		if lastUpdate == nil {
			lastUpdate = &time.Time{}
			*lastUpdate = time.Now().AddDate(-50, 0, 0)
		}
		query := db
		realTarget := objectTarget
		if m, _ := objectTarget.(merger); m != nil {
			realTarget = reflect.New(reflect.TypeOf(objectTarget).Elem()).Interface()
		}
		if t, _ := realTarget.(targeter); t != nil {
			realTarget = t.Target()
		}
		if p, _ := realTarget.(dbQueryPrepare); p != nil {
			query = p.PrepareQuery(query)
		}
		if p, _ := objectTarget.(dbQueryPrepare); p != nil {
			query = p.PrepareQuery(query)
		}
		res := query.Select(query, append(args, *lastUpdate)...).Find(realTarget)
		if m, _ := objectTarget.(merger); m != nil {
			m.Merge(reflect.ValueOf(realTarget).Elem().Interface())
		}
		return res.Error
	}
}

// Loader returns new DB loader for select query
func Loader(db *gorm.DB, args ...interface{}) loader.LoaderFnk {
	if len(args) == 0 {
		args = append(args, `updated_at>=?`)
	}
	return func(objectTarget interface{}, lastUpdate *time.Time) error {
		if lastUpdate == nil {
			lastUpdate = &time.Time{}
			*lastUpdate = time.Now().AddDate(-50, 0, 0)
		}
		query := db
		realTarget := objectTarget
		if m, _ := objectTarget.(merger); m != nil {
			realTarget = reflect.New(reflect.TypeOf(objectTarget).Elem()).Interface()
		}
		if t, _ := realTarget.(targeter); t != nil {
			realTarget = t.Target()
		}
		if p, _ := realTarget.(dbQueryPrepare); p != nil {
			query = p.PrepareQuery(query)
		}
		if p, _ := objectTarget.(dbQueryPrepare); p != nil {
			query = p.PrepareQuery(query)
		}
		res := query.Find(realTarget, append(args, *lastUpdate)...)
		if m, _ := objectTarget.(merger); m != nil {
			m.Merge(reflect.ValueOf(realTarget).Elem().Interface())
		}
		return res.Error
	}
}
