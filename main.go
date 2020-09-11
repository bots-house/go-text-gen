package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/bots-house/go-text-gen/generator"
	"github.com/bots-house/go-text-gen/loader"
)

func main() {

	if err := run(); err != nil {
		log.Fatalf("failed: %v", err)
	}
}

func run() error {

	var (
		cfgInputDir      string
		cfgInputPattern  string
		cfgOutputDir     string
		cfgDefaultLocale string
		cfgPackageName   string
	)

	flag.StringVar(&cfgInputDir, "input-dir", ".", "input directory")
	flag.StringVar(&cfgOutputDir, "output-dir", "", "output directory")
	flag.StringVar(&cfgInputPattern, "input-pattern", "*.yml", "input file pattern")
	flag.StringVar(&cfgPackageName, "pkg", "", "use custom package name")

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
		return errors.New("default locale required as postional argument")
	}

	cfgDefaultLocale = flag.Arg(0)

	if cfgPackageName == "" {
		cfgPackageName = filepath.Base(cfgInputDir)
	}

	if cfgOutputDir == "" {
		cfgOutputDir = cfgInputDir
	}

	bundle, err := loader.Load(
		cfgInputDir,
		cfgInputPattern,
		cfgDefaultLocale,
	)
	if err != nil {
		return err
	}

	data, err := generator.Generate(bundle, nil, cfgPackageName)
	if err != nil {
		return fmt.Errorf("generate: %w", err)
	}

	if err := ioutil.WriteFile(
		path.Join(cfgOutputDir, cfgPackageName+".go"),
		data,
		os.ModePerm,
	); err != nil {
		return fmt.Errorf("write file: %w", err)
	}

	return nil
}
