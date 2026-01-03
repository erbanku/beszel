package migrations

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
)

func init() {
	m.Register(func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("alerts")
		if err != nil {
			return err
		}

		// Add excluded_partitions field to alerts collection
		collection.Fields.Add(&core.JSONField{
			Id:           "json_excluded_partitions",
			Name:         "excluded_partitions",
			Required:     false,
			Presentable:  false,
			System:       false,
			Hidden:       false,
			MaxSize:      2000,
		})

		return app.Save(collection)
	}, func(app core.App) error {
		// Rollback - remove the field
		collection, err := app.FindCollectionByNameOrId("alerts")
		if err != nil {
			return err
		}

		collection.Fields.RemoveById("json_excluded_partitions")

		return app.Save(collection)
	})
}
