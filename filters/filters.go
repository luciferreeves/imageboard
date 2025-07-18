package filters

import (
	"log"

	"github.com/flosch/pongo2/v6"
)

func Initialize() {
	filters := map[string]func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error){
		"naturaltime": naturaltimeFilter,
	}

	for name, filterFunc := range filters {
		if err := pongo2.RegisterFilter(name, filterFunc); err != nil {
			log.Println("Failed to register filter:", name, "Error:", err)
		}
	}
}
