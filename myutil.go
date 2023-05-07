package pglmmyutil

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/vaughan0/go-ini"
)

const version = "1.0.0.0"

// Constant of Application
const AppIni string = "app.ini"

var isInitialize bool = false

func GetVersion() string {
	return version
}
func GetExecDir() (exPath string) {
	path, err := os.Executable() //os.Getwd() //os.Executable()
	if err != nil {
		log.Println(err)
	}
	exPath = filepath.Dir(path)
	//fmt.Println("my path: " + exPath)
	return
}

func ReadINI(section string, key string) (value string) {
	initializeMyUtil()
	//var err error
	file, err := ini.LoadFile(AppIni)
	if err != nil {
		log.Println(err)
	}
	value, ok := file.Get(section, key)
	if !ok {
		panic("'" + key + "' variable missing from '" + section + "' section")
	}
	return
}

/*
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
*/

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

// //////////////////////////////////////////////////////////////////////////////////////////////////////
// Initialization modules
// /////////////////////////////////////////////////////////////////////////////////////////////////////
func initializeMyUtil() bool {
	//fmt.Println("Init: " + strconv.FormatBool(isInitialize))
	if isInitialize {
		return true
	}
	if !FileExist(AppIni) {
		isInitialize = false
		panic("File INI '" + AppIni + "' is not exist!")
	} else {
		isInitialize = true
		return true
	}

}
