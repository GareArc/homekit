package templating

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"text/template"
)

// Renderer renders templates using text/template with optional functions.
type Renderer struct {
	Funcs template.FuncMap
}

// RenderFile renders a template file into the writer.
func (r Renderer) RenderFile(vfs fs.FS, name string, data any, dest io.Writer) error {
	tmpl := template.New(name).Funcs(r.Funcs)
	content, err := fs.ReadFile(vfs, name)
	if err != nil {
		return err
	}
	parsed, err := tmpl.Parse(string(content))
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err := parsed.Execute(&buf, data); err != nil {
		return err
	}
	if _, err := io.Copy(dest, &buf); err != nil {
		return err
	}
	return nil
}

// RenderString renders a template string and returns the result.
func (r Renderer) RenderString(name, tpl string, data any) (string, error) {
	tmpl := template.New(name).Funcs(r.Funcs)
	parsed, err := tmpl.Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}

	var buf bytes.Buffer
	if err := parsed.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}

// Render executes a template read from reader and writes to destination.
func (r Renderer) Render(reader io.Reader, data any, dest io.Writer) error {
	tplBytes, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	result, err := r.RenderString("inline", string(tplBytes), data)
	if err != nil {
		return err
	}
	_, err = io.WriteString(dest, result)
	return err
}
