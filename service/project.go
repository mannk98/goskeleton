package service

import (
	"fmt"
	"github.com/spf13/cobra"
	"goske/models"
	"goske/tpl"
	"os"
	"path"
	"text/template"
)

// Project contains name, license and paths to projects.
type Project struct {
	// v2
	PkgName      string
	Copyright    string
	AbsolutePath string
	Legal        models.License
	Viper        bool
	AppName      string
}

func NewProject() *Project {
	wd, _ := os.Getwd()
	return &Project{
		AbsolutePath: fmt.Sprintf("%s", wd),
	}
}

func NewProjectTest() *Project {
	wd, _ := os.Getwd()
	return &Project{
		AbsolutePath: fmt.Sprintf("%s/testproject", wd),
		Legal:        getLicense("", "", ""),
		Copyright:    copyrightLine("2004", "mannk"),
		AppName:      "cmd",
		PkgName:      "github.com/spf13/cobra-cli/cmd/cmd",
		Viper:        true,
	}
}

func (p *Project) createLicenseFile(year, author string) error {
	data := map[string]interface{}{
		"copyright": copyrightLine(year, author),
	}
	licenseFile, err := os.Create(fmt.Sprintf("%s/LICENSE", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer licenseFile.Close()

	licenseTemplate := template.Must(template.New("license").Parse(p.Legal.Text))
	return licenseTemplate.Execute(licenseFile, data)
}

func (p *Project) create(year, author string) error {
	// check if AbsolutePath exists
	if _, err := os.Stat(p.AbsolutePath); os.IsNotExist(err) {
		// create directory
		if err := os.Mkdir(p.AbsolutePath, PERMISSION_DIR); err != nil {
			return err
		}
	}

	// create main.go
	mainFile, err := os.Create(fmt.Sprintf("%s/main.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer mainFile.Close()

	mainTemplate := template.Must(template.New("main").Parse(string(tpl.MainTemplate())))
	err = mainTemplate.Execute(mainFile, p)
	if err != nil {
		return err
	}

	// create cmd/root.go
	if _, err = os.Stat(fmt.Sprintf("%s/cmd", p.AbsolutePath)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/cmd", p.AbsolutePath), PERMISSION_DIR))
	}
	rootFile, err := os.Create(fmt.Sprintf("%s/cmd/root.go", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer rootFile.Close()

	rootTemplate := template.Must(template.New("root").Parse(string(tpl.RootTemplate())))
	err = rootTemplate.Execute(rootFile, p)
	if err != nil {
		return err
	}

	// create license
	return p.createLicenseFile(year, author)
}

func (p *Project) InitializeProject(args []string, viper bool, userLicense, license_header, license_text, year, author string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// if init [path] is not empty
	if len(args) > 0 {
		projectPath := args[0]
		if string(projectPath[0]) != "/" {
			wd = fmt.Sprintf("%s", projectPath)
		} else if string(projectPath[0]) == "." {
			wd = fmt.Sprintf("%s/%s", wd, projectPath)
		}
	}

	modName := getModImportPath()

	project := &Project{
		AbsolutePath: wd,
		PkgName:      modName,
		Legal:        getLicense(userLicense, license_header, license_text),
		Copyright:    copyrightLine(year, author),
		Viper:        viper,
		AppName:      path.Base(modName),
	}

	if err := project.create(year, author); err != nil {
		return "", err
	}

	return project.AbsolutePath, nil
}

func (p *Project) GetAbsolutePath() string {
	return p.AbsolutePath
}

func (p *Project) GetPkgName() string {
	return p.PkgName
}

func (p *Project) GetCopyright() string {
	return p.Copyright
}

func (p *Project) GetLegal() models.License {
	return p.Legal
}
func (p *Project) GetAppName() string {
	return p.AppName
}
