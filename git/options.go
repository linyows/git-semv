package semv

// Option is func
type Option func(*Semv) interface{}

// PreReleaseName for option
func PreReleaseName(s string) Option {
	return func(v *Semv) interface{} {
		v.preReleaseName = s
		return v
	}
}

// BuildName for option
func BuildName(s string) Option {
	return func(v *Semv) interface{} {
		v.buildName = s
		return v
	}
}
