package service

import (
	"fmt"
	"goske/interfaces"
	"goske/tpl"
	"os"
	"text/template"
)

type Command struct {
	CmdName   string
	CmdParent string
	Project   interfaces.GoskeProject
}

func (c *Command) Create() error {
	cmdFile, err := os.Create(fmt.Sprintf("%s/cmd/%s.go", c.Project.GetAbsolutePath(), c.CmdName))
	if err != nil {
		return err
	}
	defer cmdFile.Close()

	commandTemplate := template.Must(template.New("sub").Parse(string(tpl.AddCommandTemplate())))
	err = commandTemplate.Execute(cmdFile, c)
	if err != nil {
		return err
	}
	return nil
}
