package xsd

import (
	"fmt"
	"github.com/labstack/gommon/log"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Workspace struct {
	Cache         map[string]*Schema
	GoModulesPath string
}

func NewWorkspace(goModulesPath, xsdPath string) (*Workspace, error) {
	ws := Workspace{
		Cache:         map[string]*Schema{},
		GoModulesPath: goModulesPath,
	}
	var err error
	_, err = ws.loadXsd(xsdPath, true)
	return &ws, err
}

func (ws *Workspace) loadXsd(xsdPath string, cache bool) (*Schema, error) {
	fmt.Println("\n\n--> " + xsdPath)
	cached, found := ws.Cache[xsdPath]
	if found {
		return cached, nil
	}

	fmt.Println("\tParsing::", xsdPath)

	var f *os.File
	u := xsdPath
	if strings.HasPrefix(xsdPath, `http`) {
		xsdPath = strings.ReplaceAll(xsdPath , `http:/` , `http://`)
		xsdPath = strings.ReplaceAll(xsdPath , `///` , `//`)
		fmt.Println(`fetch ` + xsdPath)
		resp, err := http.Get(xsdPath)
		if err != nil {
			panic(err)
		}
		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		dd := `/tmp/ddex/xxxx`
		os.MkdirAll(dd, 0777)
		fname := dd + `/` + filepath.Base(xsdPath)
		ioutil.WriteFile(fname, b, 0777)
		u = fname

	}
	f, err := os.Open(u)
	if err != nil {
		fmt.Println(`----------`)
		return nil, err
	}
	defer f.Close()

	schema, err := parseSchema(f)
	if err != nil {
		return nil, err
	}
	schema.ModulesPath = ws.GoModulesPath
	schema.filePath = xsdPath
	// Won't cache included schemas - we need to append contents to the current
	// schema.
	if cache {
		ws.Cache[xsdPath] = schema
	}

	dir := filepath.Dir(xsdPath)
	if strings.HasPrefix(xsdPath , `http`) {
		dir = ""
	} else {
		uu  , _ := url.Parse(xsdPath)
		dir = "http://" + uu.Host
	}

	for idx, _ := range schema.Includes {
		log.Info(`require...`)
		si := schema.Includes[idx]
		if err := si.load(ws, dir); err != nil {
			return nil, err
		}

		isch := si.IncludedSchema
		schema.Imports = append(isch.Imports, schema.Imports...)
		schema.Elements = append(isch.Elements, schema.Elements...)
		schema.Attributes = append(isch.Attributes, schema.Attributes...)
		schema.AttributeGroups = append(isch.AttributeGroups, schema.AttributeGroups...)
		schema.ComplexTypes = append(isch.ComplexTypes, schema.ComplexTypes...)
		schema.SimpleTypes = append(isch.SimpleTypes, schema.SimpleTypes...)
		schema.inlinedElements = append(isch.inlinedElements, schema.inlinedElements...)
		for key, sch := range isch.importedModules {
			schema.importedModules[key] = sch
		}
	}

	for idx, _ := range schema.Imports {
		x := schema.Imports[idx]
		log.Info(x.SchemaLocation)
		if !strings.HasPrefix(x.SchemaLocation, `http`) {
			current , noUrlErr := url.Parse(xsdPath)
			if noUrlErr == nil {
				host := current.Host
				pathDir := filepath.Dir(current.Path)
				override := `http://` + host + pathDir + `/` + x.SchemaLocation
				schema.Imports[idx].SchemaLocation = override
			}

		}

		if err := schema.Imports[idx].load(ws, dir); err != nil {
			return nil, err
		}
	}
	schema.compile()
	return schema, nil
}
