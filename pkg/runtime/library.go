package runtime

// only *ClassModel & *Function are exportable
type ExportableElement interface {
	GetProperty(name string) (Element, error)
	SetProperty(name string, value Element) error
	ExecMethod(name string, params []Element) (Element, error)
	String() string
	Exportable() bool
}

type Library struct {
	name         string
	exportValues map[string]ExportableElement
}

func NewLibrary(name string) *Library {
	return &Library{
		name:         name,
		exportValues: map[string]ExportableElement{},
	}
}

func (l *Library) GetName() string {
	return l.name
}

func (l *Library) GetAllExportValues() map[string]ExportableElement {
	return l.exportValues
}

func (l *Library) RegisterClass(name string, ref ExportableElement) *Library {
	l.addExportValue(name, ref)
	return l
}

func (l *Library) RegisterFunction(name string, fn ExportableElement) *Library {
	l.addExportValue(name, fn)
	return l
}

func (l *Library) addExportValue(name string, value ExportableElement) {
	l.exportValues[name] = value
}
