package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/str1ngs/util/file"
)

var (
	files = map[string]string{}
	dir   = "."
)

type Config struct {
	env string
	pkg string
}

func init() {
	log.SetPrefix("testit: ")
	log.SetFlags(log.Lshortfile)
	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	tick := time.Tick(time.Second)

	// getting env vars
	env := flag.String("env", "dev", "environment, if none specified, it will be dev")
	pkg := flag.String("pkg", "main", "package name to be specified i.e. catalog-api, if none specified, it will be main")
	flag.Parse()

	conf := Config{env: *env, pkg: *pkg}
	log.Println(conf)

	for _ = range tick {
		dirty, err := update_files(conf)
		if err != nil {
			fmt.Println(err)
		}
		if dirty {
			doMetaLinter()
			doTests(conf)
		}
	}
}

func doMetaLinter() {
	gorun := exec.Command("gometalinter")
	gorun.Stderr = os.Stderr
	gorun.Stdout = os.Stdout
	if err := gorun.Start(); err != nil {
		log.Println(err)
	}
}

//
// exec command gometalinter
func doTests(conf Config) {
	exec.Command("killall", conf.pkg).Run()
	gobuild := exec.Command("go", "build")
	gobuild.Stderr = os.Stderr
	gobuild.Stdout = os.Stdout
	if err := gobuild.Run(); err != nil {
		log.Println(err)
	}

	command := "./" + conf.pkg
	commandEnv := "-env=" + conf.env
	gorun := exec.Command(command, commandEnv)
	gorun.Stderr = os.Stderr
	gorun.Stdout = os.Stdout
	if err := gorun.Start(); err != nil {
		log.Println(err)
	}
}

func update_files(conf Config) (changed bool, err error) {
	markFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if path[:1] == "." || info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		hash, err := file.Md5(path)
		if err != nil {
			return err
		}
		if _, exists := files[path]; !exists {
			changed = true
			fmt.Println("adding", hash, path)
			files[path] = hash
			return nil
		}
		if files[path] != hash {
			changed = true
			fmt.Println("changed", path)
			files[path] = hash
		}
		return nil
	}
	if file.Exists(".testit") {
		doTests(conf)
		os.Remove(".testit")
	}
	return changed, filepath.Walk(".", markFn)
}
