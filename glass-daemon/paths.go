package daemon

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/hashicorp/errwrap"
)

// returns the system path were all
// timeglass related data is stored for
// this machine
func SystemTimeglassPath() (string, error) {
	if runtime.GOOS == "windows" {
		//@see http://blogs.msdn.com/b/patricka/archive/2010/03/18/where-should-i-store-my-data-and-configuration-files-if-i-target-multiple-os-versions.aspx
		//win 7/vista
		if path := os.Getenv("PROGRAMDATA"); path != "" {
			return filepath.Join(path, "Timeglass"), nil
		} else if path = os.Getenv("ALLUSERSPROFILE"); path != "" {
			return filepath.Join(path, "Timeglass"), nil
		}

		return "", fmt.Errorf("Expected environmnet variable 'PROGRAMDATA' or 'ALLUSERPROFILE'")
	} else if runtime.GOOS == "darwin" {
		//osx we can actually create user specific services, and as such, store data for the user specifically
		home, err := os.UserHomeDir()
		if(err!=nil) {
			return "", err
		}
		return filepath.Join(home, "Library", "PandaGlassHelper"), nil
	} else if runtime.GOOS == "linux" {
		return filepath.Join("/var/lib", "timeglass"), nil
	}

	return "", fmt.Errorf("Operating system is not yet supported")
}

func SystemTimeglassPathCreateIfNotExist() (string, error) {
	path, err := SystemTimeglassPath()
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(path, 0755)
	if err != nil {
		return "", errwrap.Wrapf(fmt.Sprintf("Failed to create Timeglass system dir '%s': {{err}}", path), err)
	}

	return path, nil
}

func DarwinTimeglassLogPaths() (string, string, error) {
	if runtime.GOOS == "windows" {
		return "", "", fmt.Errorf("Expected Darwin Env")
	} else if runtime.GOOS == "darwin" {
		path, err := SystemTimeglassPath();
		if err != nil {
			return "", "", err
		}
		return filepath.Join(path, "launch-out.log"), filepath.Join(path, "launch-err.log"), nil
	} else if runtime.GOOS == "linux" {
		return "", "", fmt.Errorf("Expected Darwin Env")
	}
	return "", "", fmt.Errorf("Operating system is not yet supported")
}