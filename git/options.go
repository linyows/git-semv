package semv

type option func(*Semv) interface{}

func PreReleaseName(s string) option {
	return func(v *Semv) interface{} {
		v.preReleaseName = s
		return v
	}
}

func BuildName(s string) option {
	return func(v *Semv) interface{} {
		v.buildName = s
		return v
	}
}
