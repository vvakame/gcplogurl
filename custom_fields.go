package gcplogurl

type CustomFields []string

const (
	customFieldsParam = "lfeCustomFields"
)

func (cf CustomFields) marshalURL(vs values) {
	vs.Del(customFieldsParam)
	for _, f := range cf {
		vs.Add(customFieldsParam, escape(f))
	}
}
