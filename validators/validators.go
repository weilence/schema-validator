package validators

func init() {
	RegisterDefault(defaultRegistry)
}

func RegisterDefault(r *Registry) {
	registerField(r)
	registerFormat(r)
	registerNetwork(r)
	registerOther(r)
}
