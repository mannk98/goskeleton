package interfaces

import "goske/models"

type GoskeProject interface {
	InitializeProject(args []string, viper bool, userLicense, license_header, license_text, year, author string) (string, error)
	GetAbsolutePath() string
	GetPkgName() string
	GetCopyright() string
	GetLegal() models.License
	GetAppName() string
}
