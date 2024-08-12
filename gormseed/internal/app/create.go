package app

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/raspiantoro/gormseeder/gormseed/internal/templates"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Config struct {
	Name             string
	Filename         string
	Path             string
	SeedFuncName     string
	RollbackFuncName string
	WithCli          bool
	GormseedDir      string
}

func CreateSeed(config Config) {
	name := strings.ReplaceAll(config.Name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")

	config.Filename = "seed_" + time.Now().Format("20060102150405") + "_" + name + ".go"

	var baseFuncName string

	for _, n := range strings.Split(name, "_") {
		baseFuncName += cases.Title(language.English, cases.Compact).String(n)
	}

	config.SeedFuncName = "Seed" + baseFuncName
	config.RollbackFuncName = "Rollback" + baseFuncName

	if config.Path == "" {
		config.Path = "db"
	}

	path, err := scanGormseed(".")
	if err != nil {
		log.Fatalln(err)
	}

	baseDir := filepath.Dir(path)

	mod := getModule()

	t, err := template.New("seed").Parse(string(templates.SeederTemplate()))
	if err != nil {
		log.Fatalln(err)
	}

	var b bytes.Buffer
	err = t.Execute(&b, &config)
	if err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile(mod.Dir+"/"+baseDir+"/"+config.Filename, b.Bytes(), os.ModePerm); err != nil {
		log.Fatalln(err)
	}
}

func scanGormseed(path string) (string, error) {
	var dir string
	var filenames []string

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}

		if filepath.Ext(path) != ".go" {
			return nil
		}

		var filename string
		filename, errs := findSeedStruct(path)
		if errs != nil {
			return errs
		}

		if filename != "" {
			filenames = append(filenames, filename)
		}

		return nil
	})

	if len(filenames) == 0 {
		return dir, errors.New("gormseed files not found. You need to initialize first")
	}

	if len(filenames) > 1 {
		return dir, errors.New("multiple seeds.go files with a Seeds struct found. provide the -d flag to specify your gormseed folder")
	}

	dir = filenames[0]

	if dir == "" && err == nil {
		err = io.EOF
	}

	return dir, err
}

func findSeedStruct(path string) (string, error) {
	var filename string
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return "", err
	}

	ast.Inspect(f, func(n ast.Node) bool {
		if filepath.Base(path) != "seeds.go" {
			return true
		}
		typeSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		_, ok = typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}

		if typeSpec.Name.Name == "Seeds" {
			filename = path
		}

		return true
	})

	return filename, nil
}
