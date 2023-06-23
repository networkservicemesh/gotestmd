package gen

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/networkservicemesh/gotestmd/internal/linker"
	"github.com/networkservicemesh/gotestmd/internal/parser"
	"github.com/spf13/cobra"
)

const baseurl = "https://raw.githubusercontent.com/networkservicemesh/deployments-k8s/main/"

// gotestmd Kernel2Kernel --bash

func New() *cobra.Command {
	var c = &cobra.Command{
		Use:     "gen",
		Aliases: []string{"generate"},
		Short:   "Generates a test based on url",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("wrong number of arguments")
			}

			fmt.Printf("Args[0]: %s\n", args[0])
			u, err := url.Parse(args[0])
			if err != nil {
				return err
			}

			exampleDir := filepath.Dir(u.Path[strings.Index(u.Path, "examples"):])

			exampleFile, err := downloadFile(filepath.Join(exampleDir, "README.md"))
			if err != nil {
				return err
			}

			p := parser.New()
			ex, err := p.Parse(exampleFile)
			ex.Dir = exampleDir
			if err != nil {
				return err
			}

			linkedExample := linker.NewLinkedExample("", ex)

			if len(ex.Includes) != 0 {
				return errors.New("example is not a leaf")
			}

			requires, _ := downloadRequiresRecursive(p, exampleDir, ex.Requires)
			linkedExample.Parents = requires

			fmt.Print("\n\nExamples:\n")
			for _, l := range []*linker.LinkedExample{linkedExample} {
				fmt.Printf("\tExample Name: %s\n", l.Name)
				fmt.Println("\tParents:")

				for _, p := range l.Parents {
					fmt.Printf("\t\t%s", p.Name)
				}

				fmt.Print("\n\n")
			}
			for _, r := range requires {
				fmt.Println(r.Run)
			}
			return nil
		},
	}

	return c
}

func downloadFile(path string) (io.Reader, error) {
	u, err := url.JoinPath(baseurl, path)

	if err != nil {
		return nil, err
	}

	resp, err := http.Get(u)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func downloadRequiresRecursive(p *parser.Parser, exampleDir string, requires []string) ([]*linker.LinkedExample, error) {
	result := make([]*linker.LinkedExample, 0)
	for _, require := range requires {
		requireDir := filepath.Join(exampleDir, require)
		requirePath := filepath.Join(exampleDir, require, "README.md")
		requireFile, _ := downloadFile(requirePath)

		ex, _ := p.Parse(requireFile)
		ex.Dir = requireDir
		ex.Includes = []string{exampleDir}

		linkedExample := linker.NewLinkedExample("", ex)

		result = append(result, linkedExample)
		if len(ex.Requires) != 0 {
			parentRequires, _ := downloadRequiresRecursive(p, filepath.Join(exampleDir, require), ex.Requires)
			linkedExample.Parents = parentRequires
		}
	}
	return result, nil
}
