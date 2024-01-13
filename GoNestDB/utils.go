package GoNestDB

import (
	"os/exec"
	"os/user"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/mr-tron/base58"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func EncodeBase58(input []byte) string {
	return base58.Encode(input)
}

func DecodeBase58(input string) ([]byte, error) {
	return base58.Decode(input)
}

func Compare(a, b interface{}) bool {
	aType := reflect.TypeOf(a)
	bType := reflect.TypeOf(b)

	if aType == bType {
		return a == b
	} else {
		bValue := reflect.ValueOf(b)
		convertedB := bValue.Convert(aType).Interface()

		return a == convertedB
	}
}

func GetRoamingDataPath() string {
	// Get the current user's information
	currentUser, err := user.Current()
	if err != nil {
		panic(err)
	}

	// Get the current module name
	cmd := exec.Command("go", "list", "-m")
	output, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	moduleName := strings.TrimSpace(string(output))

	// Normalize and capitalize the module name
	capitalizer := cases.Title(language.English)
	moduleName = capitalizer.String(strings.ReplaceAll(moduleName, "/", "-"))

	// Construct the path to the Roaming folder with the module name appended
	roamingDir := filepath.Join(currentUser.HomeDir, "AppData", "Roaming", moduleName)
	return roamingDir
	// fmt.Println(roamingDir)

	// You can also get the path to the Local AppData folder with the module name appended like this:
	// localDir := filepath.Join(currentUser.HomeDir, "AppData", "Local", moduleName)
	// fmt.Println(localDir)
}
