package models

type SmtpServerConfig struct {
	SmtpServer   string `yaml:"smtpserver"`
	SmtpPort     string `yaml:"smtpport"`
	AuthAddress  string `yaml:"authaddress"`
	AuthPassword string `yaml:"authpassword"`
}
