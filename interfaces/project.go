package interfaces

type GoskeProject interface {
	InitializeProject(args []string) (string, error)
}
