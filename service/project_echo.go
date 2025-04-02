package service

import (
	"fmt"
	"github.com/spf13/cobra"
	"goske/models"
	"goske/tpl"
	"os"
	"path"
	"strings"
	"text/template"
)

type EchoProject struct {
	// v2
	PkgName      string
	Copyright    string
	AbsolutePath string
	Legal        models.License
	Viper        bool
	AppName      string
}

func NewProjectEcho() *Project {
	wd, _ := os.Getwd()
	return &Project{
		AbsolutePath: fmt.Sprintf("%s", wd),
	}
}

func NewProjectEchoTest() *Project {
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

func (p *EchoProject) createLicenseFile(year, author string) error {
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

func (p *EchoProject) Create(year, author string) error {
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
	////////////////////////////////////////////////////////////////////////////
	// create cmd/root.go
	if _, err = os.Stat(fmt.Sprintf("%s/%s", p.AbsolutePath, DIR_CMD)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/%s", p.AbsolutePath, DIR_CMD), PERMISSION_DIR))
	}
	rootFile, err := os.Create(fmt.Sprintf("%s/%s/root.go", p.AbsolutePath, DIR_CMD))
	if err != nil {
		return err
	}
	defer rootFile.Close()

	rootTemplate := template.Must(template.New("root").Parse(string(tpl.RootTemplate())))
	err = rootTemplate.Execute(rootFile, p)
	if err != nil {
		return err
	}

	// create cmd/server.go
	/*	serverFile, err := os.Create(fmt.Sprintf("%s/%s/server.go", p.AbsolutePath, DIR_CMD))
		if err != nil {
			return err
		}
		defer rootFile.Close()

		serverTemplate := template.Must(template.New("server").Parse(string(tpl.CmdServerTemplate())))
		err = serverTemplate.Execute(serverFile, p)
		if err != nil {
			return err
		}*/

	////////////////////////////////////////////////////////////////////////////
	// create docs
	if _, err = os.Stat(fmt.Sprintf("%s/%s", p.AbsolutePath, DIR_DOCS)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/%s", p.AbsolutePath, DIR_DOCS), PERMISSION_DIR))
	}

	// create handler
	if _, err = os.Stat(fmt.Sprintf("%s/%s", p.AbsolutePath, DIR_HANDLER)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/%s", p.AbsolutePath, DIR_HANDLER), PERMISSION_DIR))
	}

	// create interfaces
	if _, err = os.Stat(fmt.Sprintf("%s/%s", p.AbsolutePath, DIR_INTERFACES)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/%s", p.AbsolutePath, DIR_INTERFACES), PERMISSION_DIR))
	}

	// create service
	if _, err = os.Stat(fmt.Sprintf("%s/%s", p.AbsolutePath, DIR_SERVICE)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/%s", p.AbsolutePath, DIR_SERVICE), PERMISSION_DIR))
	}

	// create service
	if _, err = os.Stat(fmt.Sprintf("%s/%s", p.AbsolutePath, DIR_MODELS)); os.IsNotExist(err) {
		cobra.CheckErr(os.Mkdir(fmt.Sprintf("%s/%s", p.AbsolutePath, DIR_MODELS), PERMISSION_DIR))
	}

	// create license
	return p.createLicenseFile(year, author)
}

func (p *EchoProject) InitializeProject(args []string, viper bool, userLicense, license_header, license_text, year, author string) (string, error) {
	projectPath := args[0]
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if len(args) > 0 {
		if strings.Contains(projectPath, ".") || string(projectPath[0]) != "/" {
			wd = fmt.Sprintf("%s/%s", wd, projectPath)
		} else {
			wd = fmt.Sprintf("%s", projectPath)
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

func (p *EchoProject) GetAbsolutePath() string {
	return p.AbsolutePath
}

func (p *EchoProject) GetPkgName() string {
	return p.PkgName
}

func (p *EchoProject) GetCopyright() string {
	return p.Copyright
}

func (p *EchoProject) GetLegal() models.License {
	return p.Legal
}

func (p *EchoProject) GetAppName() string {
	return p.AppName
}
