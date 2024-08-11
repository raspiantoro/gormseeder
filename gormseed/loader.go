package gormseed

import (
	"os"
	"reflect"
	"strings"

	"github.com/raspiantoro/gormseeder"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

type Seeds interface {
	Path() string
}

func Load(seeds Seeds) []*gormseeder.Seed {
	seeders := []*gormseeder.Seed{}
	items, _ := os.ReadDir(seeds.Path())

	for _, item := range items {
		if !item.IsDir() {
			if item.Name() == "seeds.go" {
				continue
			}

			nameSplit := strings.Split(item.Name(), "_")
			key := nameSplit[1]
			name := strings.Replace(strings.Join(nameSplit[2:], " "), ".go", "", -1)

			var baseFuncName string

			for _, n := range nameSplit[2:] {
				baseFuncName += cases.Title(language.English, cases.Compact).String(n)
			}

			baseFuncName = strings.Replace(baseFuncName, ".go", "", -1)

			seedFuncName := "Seed" + baseFuncName
			rollbackFuncName := "Rollback" + baseFuncName

			value := reflect.ValueOf(seeds)

			seeder := &gormseeder.Seed{
				Key:      key,
				Name:     name,
				Seed:     intoSeederFunc(value.MethodByName(seedFuncName)),
				Rollback: intoSeederFunc(value.MethodByName(rollbackFuncName)),
			}

			seeders = append(seeders, seeder)
		}
	}

	return seeders
}

func intoSeederFunc(value reflect.Value) gormseeder.SeederFunc {
	f := value.Interface().(func(*gorm.DB) error)
	return gormseeder.SeederFunc(f)
}
