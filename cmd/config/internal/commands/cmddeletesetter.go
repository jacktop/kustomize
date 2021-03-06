// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package commands

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/spf13/cobra"
	"sigs.k8s.io/kustomize/cmd/config/ext"
	"sigs.k8s.io/kustomize/cmd/config/internal/generateddocs/commands"
	"sigs.k8s.io/kustomize/kyaml/fieldmeta"
	"sigs.k8s.io/kustomize/kyaml/setters2/settersutil"
)

// NewDeleteRunner returns a command runner.
func NewDeleteSetterRunner(parent string) *DeleteSetterRunner {
	r := &DeleteSetterRunner{}
	c := &cobra.Command{
		Use:     "delete-setter DIR NAME",
		Args:    cobra.ExactArgs(2),
		Short:   commands.DeleteSetterShort,
		Long:    commands.DeleteSetterLong,
		Example: commands.DeleteSetterExamples,
		PreRunE: r.preRunE,
		RunE:    r.runE,
	}
	c.Flags().BoolVarP(&r.RecurseSubPackages, "recurse-subpackages", "R", false,
		"deletes setter recursively in all the nested subpackages")
	fixDocs(parent, c)
	r.Command = c

	return r
}

func DeleteSetterCommand(parent string) *cobra.Command {
	return NewDeleteSetterRunner(parent).Command
}

type DeleteSetterRunner struct {
	Command            *cobra.Command
	DeleteSetter       settersutil.DeleterCreator
	OpenAPIFile        string
	RecurseSubPackages bool
}

func (r *DeleteSetterRunner) preRunE(c *cobra.Command, args []string) error {
	var err error
	r.DeleteSetter.Name = args[1]
	r.DeleteSetter.DefinitionPrefix = fieldmeta.SetterDefinitionPrefix

	r.OpenAPIFile, err = ext.GetOpenAPIFile(args)
	if err != nil {
		return err
	}

	return nil
}

func (r *DeleteSetterRunner) runE(c *cobra.Command, args []string) error {
	e := executeCmdOnPkgs{
		needOpenAPI:        true,
		writer:             c.OutOrStdout(),
		rootPkgPath:        args[0],
		recurseSubPackages: r.RecurseSubPackages,
		cmdRunner:          r,
	}
	err := e.execute()
	if err != nil {
		return handleError(c, err)
	}
	return nil
}

func (r *DeleteSetterRunner) executeCmd(w io.Writer, pkgPath string) error {
	openAPIFileName, err := ext.OpenAPIFileName()
	if err != nil {
		return err
	}
	r.DeleteSetter = settersutil.DeleterCreator{
		Name:               r.DeleteSetter.Name,
		DefinitionPrefix:   fieldmeta.SetterDefinitionPrefix,
		RecurseSubPackages: r.RecurseSubPackages,
		OpenAPIFileName:    openAPIFileName,
		OpenAPIPath:        filepath.Join(pkgPath, openAPIFileName),
		ResourcesPath:      pkgPath,
	}

	err = r.DeleteSetter.Delete()
	if err != nil {
		// return err if RecurseSubPackages is false
		if !r.DeleteSetter.RecurseSubPackages {
			return err
		} else {
			// print error message and continue if RecurseSubPackages is true
			fmt.Fprintf(w, "%s in package %q\n\n", err.Error(), pkgPath)
		}
	} else {
		fmt.Fprintf(w, "deleted setter %q in package %q\n\n", r.DeleteSetter.Name, pkgPath)
	}
	return nil
}
