package models

import (
	"encoding/json"
	"fmt"
	"github.com/mannk98/goske/tpl"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

// Project contains name, license and paths to projects.
type Project struct {
	// v2
	PkgName      string
	Copyright    string
	AbsolutePath string
	Legal        License
	Viper        bool
	AppName      string
}

func (p *Project) Create() error {
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
	return p.createLicenseFile()
}

func (p *Project) createLicenseFile() error {
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

func NewProject() *Project {
	return &Project{}
}

func NewProject_Test() *Project {
	wd, _ := os.Getwd()
	return &Project{
		AbsolutePath: fmt.Sprintf("%s/testproject", wd),
		Legal:        getLicense(),
		Copyright:    copyrightLine(),
		AppName:      "cmd",
		PkgName:      "github.com/spf13/cobra-cli/cmd/cmd",
		Viper:        true,
	}
}

func (p *Project) InitializeProject(args []string) (string, error) {
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
		Legal:        getLicense(),
		Copyright:    copyrightLine(),
		Viper:        viper.GetBool("useViper"),
		AppName:      path.Base(modName),
	}

	if err := project.Create(); err != nil {
		return "", err
	}

	return project.AbsolutePath, nil
}

type Mod struct {
	Path  string
	Dir   string
	GoMod string
}

type CurDir struct {
	Dir string
}

func parseModInfo() (Mod, CurDir) {
	var mod Mod
	var dir CurDir

	m := modInfoJSON("-m")
	cobra.CheckErr(json.Unmarshal(m, &mod))

	// Unsure why, but if no module is present Path is set to this string.
	if mod.Path == "command-line-arguments" {
		cobra.CheckErr("Please run `go mod init <MODNAME>` before `cobra-cli init`")
	}

	e := modInfoJSON("-e")
	cobra.CheckErr(json.Unmarshal(e, &dir))

	return mod, dir
}

func getModImportPath() string {
	mod, cd := parseModInfo()
	fmt.Printf(fileToURL(strings.TrimPrefix(cd.Dir, mod.Dir)))
	return path.Join(mod.Path, fileToURL(strings.TrimPrefix(cd.Dir, mod.Dir)))
}

func fileToURL(in string) string {
	i := strings.Split(in, string(filepath.Separator))
	return path.Join(i...)
}

func modInfoJSON(args ...string) []byte {
	cmdArgs := append([]string{"list", "-json"}, args...)
	out, err := exec.Command("go", cmdArgs...).Output()
	cobra.CheckErr(err)

	return out
}
