package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	pi := filepath.Join(wd, "npm", "molecast", "package.json")
	raw, err := ioutil.ReadFile(pi)
	if err != nil {
		panic(err)
	}

	pkg := map[string]interface{}{}
	if err := json.Unmarshal(raw, &pkg); err != nil {
		panic(err)
	}

	ver := pkg["version"].(string)

	script := filepath.Join(wd, "npm", "molecast", "bootstrap.sh")
	scriptRaw, err := ioutil.ReadFile(script)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile("VERSION=.*")
	scriptStr := re.ReplaceAllString(string(scriptRaw), "VERSION="+ver)
	err = ioutil.WriteFile(script, []byte(scriptStr), 0644)
	if err != nil {
		panic(err)
	}

	// update subs
	dir := filepath.Join(wd, "npm")
	subs, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	for _, sub := range subs {
		if sub.Name() == "molecast" || !sub.IsDir() {
			continue
		}
		pi := filepath.Join(dir, sub.Name(), "package.json")
		raw, err := ioutil.ReadFile(pi)
		if err != nil {
			panic(err)
		}

		pkg := map[string]interface{}{}
		if err := json.Unmarshal(raw, &pkg); err != nil {
			panic(err)
		}

		pkg["version"] = ver

		str, err := json.MarshalIndent(pkg, "", "  ")
		if err != nil {
			panic(err)
		}

		err = ioutil.WriteFile(pi, str, 0644)
		if err != nil {
			panic(err)
		}
	}
}
