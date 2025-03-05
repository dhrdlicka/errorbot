package repo

type ErrorInfo struct {
	Code        uint32 `yaml:"code"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
}
