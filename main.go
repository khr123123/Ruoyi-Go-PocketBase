// main.go
package main

import (
	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"log"
)

func main() {
	app := pocketbase.New()
	app.OnRecordViewRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		return e.Next()
	})
	app.OnRecordCreateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		return e.Next()
	})
	app.OnRecordUpdateRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		return e.Next()
	})
	app.OnRecordDeleteRequest().BindFunc(func(e *core.RecordRequestEvent) error {
		return e.Next()
	})
	app.OnServe().BindFunc(func(se *core.ServeEvent) error {
		return se.Next()
	})
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
