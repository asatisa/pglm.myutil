package pglmmyutil

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"

	"github.com/vaughan0/go-ini"
)

// Constant of Application
const APP_INI string = "app.ini"
const version = "1.0.0.01" //Update version of AutoCorrectDistance 4_1

var isInitialize bool = false

func GetVersion() string {
	return version
}

// //////////////////////////////////////////////////////////////////////////////////////////////////////
// Initialization modules
// /////////////////////////////////////////////////////////////////////////////////////////////////////
func initializeMyUtil() bool {
	//fmt.Println("Init: " + strconv.FormatBool(isInitialize))
	if isInitialize {
		return true
	}
	if !FileExist(APP_INI) {
		isInitialize = false
		panic("File INI '" + APP_INI + "' is not exist!")
	} else {
		isInitialize = true
		return true
	}

}

// /////////////////////////////////////////////////////////////////////////////////////////////////////
// Initialize function and Configuration declarations
// /////////////////////////////////////////////////////////////////////////////////////////////////////
func ReadINI(section string, key string) (value string) {
	initializeMyUtil()
	//var err error
	file, err := ini.LoadFile(APP_INI)
	if err != nil {
		log.Println(err)
	}
	value, ok := file.Get(section, key)
	if !ok {
		panic("'" + key + "' variable missing from '" + section + "' section")
	}
	return
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////

// /////////////////////////////////////////////////////////////////////////////////////////////////////
// File Management Functions
// /////////////////////////////////////////////////////////////////////////////////////////////////////

func GetExecDir() (exPath string) {
	path, err := os.Executable() //os.Getwd() //os.Executable()
	if err != nil {
		log.Println(err)
	}
	exPath = filepath.Dir(path)
	//fmt.Println("my path: " + exPath)
	return
}

// Check file is exist.
func FileExist(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		//fmt.Printf("File '" + filename + "' exists\n")
		return true
	} else {
		//fmt.Printf("File '" + filename + "' does not exist\n")
		return false
	}
}

func CopyFile(src string, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func CopyReplaceFile(src string, dest string, forceReplace bool) error {
	if forceReplace {
		RemoveFile(dest)
	}

	ret := CopyFile(src, dest)
	return ret
}

func CopyReplaceFileAndDeleteSource(src string, dest string, forceReplace bool, forceDeleteSrc bool) error {
	if forceReplace {
		RemoveFile(dest) //remove destination file
	}

	ret := CopyFile(src, dest)

	if forceDeleteSrc {
		RemoveFile(src) //after copy finished, delete source file
	}
	return ret
}

func RemoveFile(filePath string) error {
	var err error
	if FileExist(filePath) {
		err = os.Remove(filePath)
		fmt.Println("Remove file:'", filePath, "' is done.")
	} else {
		fmt.Println("Remove file:'", filePath, "' is done, but file does not exist.")
	}

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func MoveFileAndReplace(src string, dest string, forceReplace bool) error {
	if forceReplace {
		RemoveFile(dest)
	}

	err := os.Rename(src, dest)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func CreateFolderIfNotExist(folderPath string) error {
	_, err := os.Stat(folderPath)
	if os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0755)
		if err != nil {
			return err
		}
		fmt.Println("Folder created:", folderPath)
	} else if err != nil {
		return err
	} else {
		fmt.Println("Folder already exists:", folderPath)
	}
	return nil
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////

// /////////////////////////////////////////////////////////////////////////////////////////////////////
// Terminal manipulation functions
// /////////////////////////////////////////////////////////////////////////////////////////////////////
func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func ClearTerminal() {
	switch runtime.GOOS {
	case "darwin":
		runCmd("clear")
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cmd", "/c", "cls")
	default:
		runCmd("clear")
	}
}

func ClearScreen() {
	const (
		cursorPosition = 0
		clearLength    = 10000
	)

	var (
		kernel32                = syscall.NewLazyDLL("kernel32.dll")
		procGetStdHandle        = kernel32.NewProc("GetStdHandle")
		procFillConsole         = kernel32.NewProc("FillConsoleOutputCharacterW")
		stdOutHandle     int32  = -11 // STD_OUTPUT_HANDLE
		blank            uint16 = ' '
	)

	// Get the console handle.
	handle, _, _ := procGetStdHandle.Call(uintptr(stdOutHandle))

	// Clear the console.
	var written uint32
	procFillConsole.Call(
		uintptr(handle),
		uintptr(blank),
		uintptr(clearLength),
		uintptr(cursorPosition),
		uintptr(unsafe.Pointer(&written)),
	)

	// Move the cursor to the top left of the console.
	fmt.Fprintf(os.Stdout, "\033[0;0H")
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////

// Convert a path name that include back slash to standard format of golang
func ConvertPathBackslashToStandard(pathWithBackSlash string) string {
	standardPath := strings.ReplaceAll(pathWithBackSlash, "\\", "/")
	return standardPath
}

func ConvertPathStandardToBackslash(path string) string {
	standardPath := strings.ReplaceAll(path, "/", "\\")
	return standardPath
}

func TrimAll(str string) string {
	trimmedStr := strings.TrimSpace(str)
	return trimmedStr
}

// /////////////////////////////////////////////////////////////////////////////////////////////////////
