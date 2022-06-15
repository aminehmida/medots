package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/aminehmida/medots/shell"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Source         string `yaml:"source" omitempty`
	Destination    string `yaml:"destination" omitempty`
	Run            string `yaml:"run" omitempty`
	RunInteractive string `yaml:"run_interactive" omitempty`
	IfOS           string `yaml:"if_os" omitempty`
}

// TODO: make this a config flag
const ShellToUse = "/bin/sh"

func ReadConfig(configPath string) (*map[string][]AppConfig, error) {
	config := make(map[string][]AppConfig)
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func symlinkWithBackup(src, dst string) error {
	if fileInfo, err := os.Lstat(dst); err == nil {
		// Check if the file is a symlink
		if fileInfo.Mode()&os.ModeSymlink == os.ModeSymlink {
			fmt.Printf("%s is already a symbolic link. Skipping.\n", dst)
			return nil
		}
		if fileInfo.Mode().IsRegular() {
			fmt.Println("Backing up " + dst)
			os.Rename(dst, dst+".bak")
		}
	}
	// Create the symlink
	srcAbs, err := filepath.Abs(src)
	if err != nil {
		return err
	}
	err = os.Symlink(srcAbs, dst)
	if err != nil {
		return err
	}
	fmt.Println("Symlink created for " + dst)
	return nil
}

func (obj AppConfig) Link() (*string, *string, error) {
	// fmt.Println("==>", obj.IfOS, "==>", obj.Template)
	if obj.IfOS == "" || obj.IfOS == runtime.GOOS {
		if obj.Destination != "" && obj.Source != "" {
			dest, err := homedir.Expand(filepath.Dir(obj.Destination))
			if err != nil {
				return nil, nil, err
			}
			// Check if destination is a directory
			destIsDir := obj.Destination[len(obj.Destination)-1] == '/'

			// List all files in case source is a glob
			files, err := filepath.Glob(obj.Source)
			if err != nil {
				return nil, nil, err
			}

			//Error if no files found
			if len(files) == 0 {
				return nil, nil, fmt.Errorf("No files found for %s", obj.Source)
			}
			// Error if multiple source files and destination is not a directory
			if len(files) > 1 && !destIsDir {
				return nil, nil, fmt.Errorf("multiple files found for %s, but destination is not a directory", obj.Source)
			}

			// Create the destination directory if it doesn't exist
			if _, err := os.Stat(dest); os.IsNotExist(err) {
				fmt.Println("Creating directory: " + dest)
				os.MkdirAll(dest, 0750)
			}

			if destIsDir {
				for _, file := range files {
					destFile := filepath.Join(dest, filepath.Base(file))
					fmt.Println("Linking " + file + " to " + destFile)
					symlinkWithBackup(file, destFile)
				}
			} else {
				dest, err := homedir.Expand(obj.Destination)
				if err != nil {
					return nil, nil, err
				}
				symlinkWithBackup(files[0], dest)
			}
		}
		if obj.Run != "" {
			fmt.Println("Running: " + obj.Run)
			stdout, stderr, err := shell.Run(obj.Run, ShellToUse, false)
			if err != nil {
				return nil, nil, err
			}
			return &stdout, &stderr, nil
		}
		if obj.RunInteractive != "" {
			fmt.Println("Running: " + obj.RunInteractive)
			stdout, stderr, err := shell.Run(obj.RunInteractive, ShellToUse, true)
			if err != nil {
				return nil, nil, err
			}
			return &stdout, &stderr, nil
		}
	} else {
		fmt.Printf("Skipping %s because it's for a different os: %s\n", obj.Source, obj.IfOS)
	}
	return nil, nil, nil
}
