package loggers

type Config struct {
	LogLevel           string `json:"log_level" yaml:"log_level"`
	EnableConsolWriter bool   `json:"enable_consol_writer" yaml:"enable_consol_writer"`
	EnableConsolColor  bool   `json:"enable_consol_color" yaml:"enable_consol_color"`
	EnableFileWriter   bool   `json:"enable_file_writer" yaml:"enable_file_writer"`
	FileName           string `json:"file_name,omitempty" yaml:"file_name,omitempty"`
	MaxSizeMB          int    `json:"max_size_mb,omitempty" yaml:"max_size_mb,omitempty"`
	MaxAgeDays         int    `json:"max_age_days,omitempty" yaml:"max_age_days,omitempty"`
	MaxBackups         int    `json:"max_backups,omitempty" yaml:"max_backups,omitempty"`
	UseLocalTime       bool   `json:"use_local_time,omitempty" yaml:"use_local_time,omitempty"`
	Compress           bool   `json:"compress,omitempty" yaml:"compress,omitempty"`
}

func NewConfig() *Config {
	return &Config{
		LogLevel:           "info",
		EnableConsolWriter: true,
		EnableConsolColor:  false,
		EnableFileWriter:   true,
		FileName:           "log.json",
		MaxSizeMB:          100,
		MaxAgeDays:         7,
		MaxBackups:         7,
		UseLocalTime:       true,
		Compress:           false,
	}
}
