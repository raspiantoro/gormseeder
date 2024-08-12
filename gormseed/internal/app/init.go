package app

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/raspiantoro/gormseeder/gormseed/internal/templates"
)

func Init(config Config) {
	if err := os.Setenv("GOWORK", "off"); err != nil {
		log.Fatalln(err)
	}
	defer os.Setenv("GOWORK", "on")

	mod := getModule()

	if mod.Path == "command-line-arguments" {
		fmt.Println("ERROR: `go.mod` file not found in the current directory")
		return
	}

	if config.GormseedDir == "" {
		config.GormseedDir = "db"
	}

	if config.WithCli {
		config.GormseedDir = config.GormseedDir + "/seeder/seeds"
	} else {
		config.GormseedDir = config.GormseedDir + "/seeds"
	}

	if _, err := os.Stat(config.GormseedDir); err != nil {
		err = os.MkdirAll(config.GormseedDir, os.ModePerm)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if err := os.WriteFile(config.GormseedDir+"/seeds.go", templates.SeedTemplate(), os.ModePerm); err != nil {
		log.Fatalln(err)
	}

	if config.WithCli {

		t, err := template.New("cli").Parse(string(templates.CliTemplate()))
		if err != nil {
			log.Fatalln(err)
		}

		templateProps := map[string]string{
			"SeedModuleName": mod.Path + "/" + config.GormseedDir,
		}

		var b bytes.Buffer
		err = t.Execute(&b, templateProps)
		if err != nil {
			log.Fatalln(err)
		}

		if err := os.WriteFile(mod.Dir+"/"+filepath.Clean(filepath.Join(config.GormseedDir, ".."))+"/seeder.go", b.Bytes(), os.ModePerm); err != nil {
			log.Fatalln(err)
		}
	}

	cmd := exec.Command("go", "mod", "tidy")

	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
