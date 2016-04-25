package gforms

import (
	"strings"
)

type Field interface {
	New() FieldInterface
	// Get field name
	GetName() string
	GetDescription() string
	GetLabel() string
	GetWidget() Widget
	GetValidators() Validators
	SetWidget(widget Widget)
}

type Fields struct {
	list    []Field
	nameMap map[string]Field
}

// Get ordered list for field object.
func (fs *Fields) List() []Field {
	return fs.list
}

// Get field by name.
func (fs *Fields) Get(name string) (Field, bool) {
	v, ok := fs.nameMap[name]
	return v, ok
}

func (fs *Fields) AddField(field Field) bool {
	name := field.GetName()
	_, exists := fs.Get(name)
	if !exists {
		fs.list = append(fs.list, field)
		fs.nameMap[name] = field
		return true
	}
	return false
}

func NewFields(fields ...Field) *Fields {
	fs := Fields{}
	fs.nameMap = map[string]Field{}
	for _, field := range fields {
		fs.nameMap[field.GetName()] = field
	}
	fs.list = fields
	return &fs
}

type BaseField struct {
	Name        string
	Description string
	HideLabel   bool
	Label       string
	Validators  Validators
	Widget      Widget
}

func (f *BaseField) GetName() string {
	return f.Name
}

func (f *BaseField) GetDescription() string {
	return f.Description
}

func (f *BaseField) GetLabel() string {
	if f.HideLabel {
		return ""
	}
	if f.Label != "" {
		return f.Label
	}
	return strings.Title(f.Name)
}

func (f *BaseField) GetWidget() Widget {
	return f.Widget
}

func (f *BaseField) GetValidators() Validators {
	return f.Validators
}

func (f *BaseField) SetWidget(widget Widget) {
	f.Widget = widget
}

type FieldInterface interface {
	GetModel() Field
	GetName() string
	GetDescription() string
	GetLabel() string
	GetV() *V
	GetWidget() Widget
	SetInitial(string)
	Clean(Data) error
	Validate(*FormInstance) []string
	Html() string
	html() string
	Errors() []string
	SetErrors([]string)
	HasError() bool
}

type FieldInstance struct {
	Model  Field
	errors []string
	V      *V
	FieldInterface
}

func (f *FieldInstance) GetModel() Field {
	return f.Model
}

func (f *FieldInstance) GetName() string {
	return f.Model.GetName()
}

func (f *FieldInstance) GetLabel() string {
	return f.Model.GetLabel()
}

func (f *FieldInstance) GetDescription() string {
	return f.Model.GetDescription()
}

func (f *FieldInstance) GetWidget() Widget {
	return f.Model.GetWidget()
}

func (f *FieldInstance) GetV() *V {
	return f.V
}

func (f *FieldInstance) Errors() []string {
	return f.errors
}

func (f *FieldInstance) SetErrors(errs []string) {
	f.errors = errs
}

func (f *FieldInstance) HasError() bool {
	return len(f.errors) != 0
}

func (f *FieldInstance) SetInitial(v string) {
	f.V.RawStr = v
}

func (f *FieldInstance) Validate(fo *FormInstance) []string {
	vs := f.Model.GetValidators()
	if vs == nil {
		return nil
	}
	var errs []string
	for _, v := range vs {
		err := v.Validate(f, fo)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	return errs
}

type FieldInterfaces struct {
	list    []FieldInterface
	nameMap map[string]FieldInterface
}

func newFieldInterfaces(fs *Fields) *FieldInterfaces {
	fis := new(FieldInterfaces)
	fis.list = []FieldInterface{}
	fis.nameMap = map[string]FieldInterface{}
	for _, f := range fs.list {
		nf := f.New()
		fis.nameMap[f.GetName()] = nf
		fis.list = append(fis.list, nf)
	}
	return fis
}

func fieldToHtml(f FieldInterface) string {
	gd := f.GetModel().GetWidget()
	if gd == nil {
		return f.html()
	} else {
		return gd.html(f)
	}
}
