package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"text/template"

	"github.com/hsiaosiyuan0/mole/util"
)

type EntitiesData struct {
	MaxKeyLen int
	Entities  map[string]interface{}
	Keys      map[string]bool
}

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	distFile := path.Join(wd, "html_entity.go")
	_, err = os.Stat(distFile)
	if err != nil {
		return
	}

	// https://html.spec.whatwg.org/multipage/named-characters.html#named-character-references
	resp, err := http.Get("https://html.spec.whatwg.org/entities.json")
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	var rawEntities map[string]interface{}
	err = json.Unmarshal(body, &rawEntities)
	if err != nil {
		log.Fatal(err)
	}

	maxKeyLen := 0
	entities := make(map[string]interface{})
	quotedKeys := make(map[string]bool)
	for key := range rawEntities {
		kl := len(key)
		if kl > maxKeyLen {
			maxKeyLen = kl
		}
		val := rawEntities[key].(map[string]interface{})
		cs := val["characters"].(string)
		if !strings.HasSuffix(key, ";") {
			continue
		}
		if _, ok := entities[cs]; ok {
			continue
		}
		entities[key] = val

		qk := strconv.Quote(cs)
		quotedKeys[qk] = true
	}

	tpl, err := template.New("html_entities").Parse(`
// Code generated by scripts/html_entity. DO NOT EDIT.

//go:generate go run github.com/hsiaosiyuan0/mole/script/html_entity

package parser

type HTMLEntity struct {
	Name       string
	CodePoints []rune
}

var MaxHTMLEntityName int = {{ .MaxKeyLen }}

var HTMLEntities = map[string]HTMLEntity{
  {{- range $key, $value := .Entities }}
  "{{ $key }}": {"{{ $key }}", []rune{ {{ range $v := $value.codepoints }} {{ $v }}, {{ end }} }},
  {{- end }}
}

var HTMLEntityNames = map[string]bool{
  {{- range $key, $value := .Keys }}
  {{ $key }}: {{ $value }},
  {{- end }}
}
  `)

	if err != nil {
		log.Fatal(err)
	}

	var gen bytes.Buffer
	err = tpl.Execute(&gen, &EntitiesData{maxKeyLen, entities, quotedKeys})
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(distFile, gen.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}

	util.Shell("gofmt", "-w", distFile)
}
