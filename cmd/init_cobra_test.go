package cmd

import (
	"fmt"
	"goske/interfaces"
	"goske/service"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
)

func TestGoldenInitCmd(t *testing.T) {

	var goSke interfaces.GoskeProject
	goSke = service.NewProjectTest()

	dir, err := ioutil.TempDir("", "cobra-init")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	tests := []struct {
		name      string
		args      []string
		pkgName   string
		expectErr bool
	}{
		{
			name:      "successfully creates a project based on module",
			args:      []string{"testproject"},
			pkgName:   "github.com/spf13/testproject",
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			viper.Set("useViper", true)
			viper.Set("license", "apache")
			projectPath, err := goSke.InitializeProject(tt.args, false, "MIT", "", "", "1998", "mannk")
			defer func() {
				if projectPath != "" {
					os.RemoveAll(projectPath)
				}
			}()

			if !tt.expectErr && err != nil {
				t.Fatalf("did not expect an error, got %s", err)
			}
			if tt.expectErr {
				if err == nil {
					t.Fatal("expected an error but got none")
				} else {
					// got an expected error nothing more to do
					return
				}
			}

			expectedFiles := []string{"LICENSE", "main.go", "cmd/root.go"}
			for _, f := range expectedFiles {
				generatedFile := fmt.Sprintf("%s/%s", projectPath, f)
				goldenFile := fmt.Sprintf("testdata/%s.golden", filepath.Base(f))
				err := compareFiles(generatedFile, goldenFile)
				if err != nil {
					t.Fatal(err)
				}
			}
		})
	}
}
