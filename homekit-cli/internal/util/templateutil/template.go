package templateutil

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/homekit/homekit-cli/internal/util/bufutil"
)

func RenderTemplateInBytes(content []byte, data any, name string, bufPool *bufutil.Pool) ([]byte, error) {
	tplt := template.New(name)
	if _, err := tplt.Parse(string(content)); err != nil {
		return nil, fmt.Errorf("parse template: %w", err)
	}
	var writer *bytes.Buffer
	if bufPool != nil {
		writer = bufPool.Get()
		defer bufPool.Put(writer)
	} else {
		writer = bytes.NewBuffer(nil)
	}
	if err := tplt.Execute(writer, data); err != nil {
		return nil, err
	}
	return writer.Bytes(), nil
}
