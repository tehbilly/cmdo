package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	"runtime"

	"github.com/kardianos/osext"
	"github.com/spf13/cobra"
)

var (
	force   = false
	clobber = false
	copy    = false
)

func init() {
	installCmd.Flags().BoolVarP(&force, "force", "f", false, "Allow installation to existing directory.")
	installCmd.Flags().BoolVarP(&clobber, "clobber", "c", false, "Overwrite existing files in target destination.")
	installCmd.Flags().BoolVar(&copy, "copy", false, "Copy this file to the target destination as well. [Currently unimplemented]")
	installCmd.Flags().StringSliceP("skip", "s", []string{}, "Commands to skip during installation.")

	RootCmd.AddCommand(installCmd)
}

// Commands to not create shims for
var skipCommands = []string{
	"help",
	"install",
	"version",
}

var installCmd = &cobra.Command{
	Use:   "install [flags] targetdirectory",
	Short: "Install cmdo and all subcommand shims.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
			os.Exit(0)
		}

		// Installation target
		target, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Println("Unable to determine absolute path for:", args[0])
			fmt.Println(err)
			os.Exit(-1)
		}

		// Should never fail, but better safe than sorry
		fi, err := os.Stat(target)
		if err != nil && !os.IsNotExist(err) {
			fmt.Println("Unable to access target path:", err)
			os.Exit(-1)
		}

		// Supplied a file as target? Shame
		if fi != nil && fi.IsDir() == false {
			fmt.Println("Target is a file. Try again.")
			os.Exit(-1)
		}

		if fi != nil && fi.IsDir() && force == false {
			fmt.Println("Target directory exists. Use --force if you want to install to existing directory.")
			os.Exit(-1)
		}

		if toSkip, err := cmd.Flags().GetStringSlice("skip"); err == nil {
			skipCommands = append(skipCommands, toSkip...)
		}

		// Create target directory if it does not exist
		if fi == nil {
			if err := os.MkdirAll(target, 0755); err != nil {
				fmt.Println("Unable to create target directory:", err)
				os.Exit(-1)
			}
		}

		// Loop through all non-hidden, non-skipped commands and create shims
		fmt.Println("Installing shims.")
		for _, cmd := range RootCmd.Commands() {
			name := cmd.Name()
			// Make sure we shouldn't skip this command
			if stringSliceContains(skipCommands, name) {
				continue
			}

			fmt.Printf("  Creating shim for '%s'...\n", cmd.Name())
			if err := createShim(target, cmd.Name()); err != nil {
				fmt.Printf("Error creating shim for %s: %s\n", cmd.Name(), err.Error())
			}
		}

		// All done!
		fmt.Println("All done! Shims should be in:", target)
	},
}

func stringSliceContains(slice []string, search string) bool {
	for _, v := range slice {
		if v == search {
			return true
		}
	}
	return false
}

// TODO Implement target exists && clobber == false
func targetExists(target string) (bool, error) {
	_, err := os.Stat(target)
	if err == nil {
		return false, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func createShim(dir, command string) error {
	switch runtime.GOOS {
	case "windows":
		return winShim(dir, command)
	default:
		return shellShim(dir, command)
	}
}

func winShim(dir, command string) error {
	target := filepath.Join(dir, command+".cmd")
	if fi, err := os.Stat(target); err != nil && !os.IsNotExist(err) {
		return err
	} else if fi != nil && !clobber {
		fmt.Printf("    Skipping shim for '%s', target exists and clobber == false.\n", command)
		return nil
	}

	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	cmdo, _ := filepath.Abs(os.Args[0])
	cmd := fmt.Sprintf(`@"%s" %s %%*`, cmdo, command)
	if _, err := f.WriteString(cmd); err != nil {
		return err
	}

	return nil
}

func shellShim(dir, command string) error {
	target := filepath.Join(dir, command)
	if _, err := os.Stat(target); err != nil && !os.IsNotExist(err) {
		return err
	}

	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()

	tmpl := template.New("")
	tmpl, err = tmpl.Parse(shellShimTemplate)
	if err != nil {
		return err
	}

	cmdo, _ := osext.Executable()
	if err := tmpl.Execute(f, map[string]string{"CmdoPath": cmdo, "Command": command}); err != nil {
		// TODO Remove the created file
		return err
	}

	if err := os.Chmod(target, 0755); err != nil {
		return err
	}

	return nil
}

var shellShimTemplate = `#!/bin/sh

cmdo="{{.CmdoPath}}"

# TODO: Test in msys/mingw/etc.
case $(uname) in
    *CYGWIN*) cmdo=$(cygpath -w "${cmdo}");;
esac

if [ -x "${cmdo}" ]; then
  "${cmdo}" "{{.Command}}" "$@"
  ret=$?
  exit $ret
else
  echo "Unable to find base executable: ${cmdo}"
  exit 1
fi
`
