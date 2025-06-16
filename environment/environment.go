package environment

type Env struct {
	Fonts
}

var ENV *Env

func NewEnv() *Env {
	env := &Env{
		Fonts: *NewFontsCollection(),
	}
	ENV = env
	return env
}
