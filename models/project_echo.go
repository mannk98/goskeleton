package models

import (
	"fmt"
	"github.com/mannk98/goske/tpl"
	"github.com/spf13/cobra"
	"os"
	"text/template"
)

type EchoProject struct {
	// v2
	PkgName      string
	Copyright    string
	AbsolutePath string
	Legal        License
	Viper        bool
	AppName      string
}

func (p *EchoProject) Create() error {
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
	return p.createLicenseFile()
}

func (p *EchoProject) createLicenseFile() error {
	data := map[string]interface{}{
		"copyright": copyrightLine(),
	}
	licenseFile, err := os.Create(fmt.Sprintf("%s/LICENSE", p.AbsolutePath))
	if err != nil {
		return err
	}
	defer licenseFile.Close()

	licenseTemplate := template.Must(template.New("license").Parse(p.Legal.Text))
	return licenseTemplate.Execute(licenseFile, data)
}
