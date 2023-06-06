package config

type MySQL struct {
	DSN           string
	SlowThreshold int
	LogLevel      int
	Colorful      bool
}