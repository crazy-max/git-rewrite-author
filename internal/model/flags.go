package model

// Flags holds flags from command line
type Flags struct {
	Repo      string
	LogLevel  string
	LogCaller bool

	GitUsername string
	GitEmail    string

	Old     string
	Correct string

	AuthorsJSON string
}
