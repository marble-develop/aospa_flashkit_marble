//go:build windows
// +build windows

package main

import (
	"bufio"
	"fmt"
	"log"

	"bytes"
	"os"
	"os/exec"

	"strings"

	"archive/zip"
	"image/color"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"

	// "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/src-d/go-git.v4"
)

type myTheme struct{}

func (myTheme) Color(c fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	switch c {
	case theme.ColorNameDisabledButton:
		return color.NRGBA{R: 0x26, G: 0x26, B: 0x26, A: 0xff}
	case theme.ColorNameDisabled:
		return color.NRGBA{R: 0x42, G: 0xfd, B: 0xe1, A: 0xfd}
	case theme.ColorNameScrollBar:
		return color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0x99}
	case theme.ColorNameShadow:
		return color.NRGBA{R: 0x0, G: 0x0, B: 0x0, A: 0x66}
	default:
		return theme.DefaultTheme().Color(c, v)
	}
}

func (myTheme) Font(s fyne.TextStyle) fyne.Resource {
	if s.Monospace {
		return theme.DefaultTheme().Font(s)
	}
	if s.Bold {
		if s.Italic {
			return theme.DefaultTheme().Font(s)
		}
		return theme.DefaultTheme().Font(s)
	}
	if s.Italic {
		return theme.DefaultTheme().Font(s)
	}
	return theme.DefaultTheme().Font(s)
}

func (myTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(n)
}

func (myTheme) Size(s fyne.ThemeSizeName) float32 {
	switch s {
	case theme.SizeNameCaptionText:
		return 11
	case theme.SizeNameInlineIcon:
		return 20
	case theme.SizeNamePadding:
		return 4
	case theme.SizeNameScrollBar:
		return 16
	case theme.SizeNameScrollBarSmall:
		return 12
	case theme.SizeNameSeparatorThickness:
		return 1
	case theme.SizeNameText:
		return 14
	case theme.SizeNameInputBorder:
		return 2
	default:
		return theme.DefaultTheme().Size(s)
	}
}

// Function to validate the firmware zip file
func validatezip(zipFile string) {
	zipFilePath := zipFile
	targetFile := "firmware-update/modem.img"

	// Open the zip file
	fmFile, err := zip.OpenReader(zipFilePath)
	if err != nil {
		fmt.Println("Invalid Firmware File..aborting", err)
		os.Exit(1)
	}
	defer fmFile.Close()

	// Check if the specific file exists in the zip file
	found := false
	for _, file := range fmFile.File {
		if file.Name == targetFile {
			found = true
			break
		}
	}

	// Throw an error if the specific file is not found
	if !found {
		fmt.Println("Invalid Firmware ZIP detected")
		os.Exit(1)
	} else {
		fmt.Println("Firmware ZIP detected")
	}

}

// Function to get fastboot info
func getfastbootinfo(destDir string, outputTextArea *widget.Entry) {
	var outputBuffer2 bytes.Buffer

	cmd := exec.Command(destDir+"\\platform-tools-windows\\fastboot.exe", "getvar", "all")
	cmd.Dir = destDir

	cmd.Stdout = &outputBuffer2
	cmd.Stderr = &outputBuffer2
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return
	}
	cmd.Wait()

	scanner := bufio.NewScanner(bytes.NewReader(outputBuffer2.Bytes()))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		outputTextArea.SetText(outputTextArea.Text + "\n" + m)
		outputTextArea.CursorRow = outputTextArea.CursorRow + 1
	}
}

// Function to flash ROM
func flashrom(destDir string, outputTextArea *widget.Entry) {
	var outputBuffer2 bytes.Buffer
	cmd := exec.Command("cmd.exe", "/C", destDir+"\\flash_aospa_windows.cmd")
	cmd.Dir = destDir

	cmd.Stdout = &outputBuffer2
	cmd.Stderr = &outputBuffer2
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return
	}
	cmd.Wait()

	scanner := bufio.NewScanner(bytes.NewReader(outputBuffer2.Bytes()))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		outputTextArea.SetText(outputTextArea.Text + "\n" + m)
		outputTextArea.CursorRow = outputTextArea.CursorRow + 1
	}
}

// Function to flash firmware
func flashfirmware(destDir string, outputTextArea *widget.Entry) {
	var outputBuffer2 bytes.Buffer
	cmd := exec.Command("cmd.exe", "/C", destDir+"\\flash_firmware_windows.cmd")
	cmd.Dir = destDir

	cmd.Stdout = &outputBuffer2
	cmd.Stderr = &outputBuffer2
	err := cmd.Start()
	if err != nil {
		fmt.Println("Error starting command:", err)
		return
	}
	cmd.Wait()

	scanner := bufio.NewScanner(bytes.NewReader(outputBuffer2.Bytes()))
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		outputTextArea.SetText(outputTextArea.Text + "\n" + m)
		outputTextArea.CursorRow = outputTextArea.CursorRow + 1
	}
}

// Main function
func main() {
	// Set the repository URL
	repoURL := "https://github.com/ghostrider-reborn/aospa-flashing-kit.git"

	// Set the destination directory
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	destDir := filepath.Dir(ex) + "\\aospa-flashing-kit"

	a := app.New()
	a.Settings().SetTheme(&myTheme{})
	w := a.NewWindow("AOSPA Fastboot Flashing Kit - Marble Edition")

	// Create the checkbox1
	flashCheckbox0 := widget.NewCheck("Flash ROM", nil)
	// flashCheckbox0.Checked = false

	// Create the text input fields
	input2 := widget.NewEntry()
	input3 := widget.NewEntry()

	// Create the checkbox2
	flashCheckbox := widget.NewCheck("Flash Firmware", nil)

	// Create the browse button2
	srcFile2 := ""
	browseButton2 := widget.NewButton("Open Firmware Zip", func() {
		dialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
			if err == nil && file != nil {
				// Handle the selected file
				selectedFile := file.URI().String()
				srcFile2 = strings.TrimPrefix(selectedFile, "file://")
				fmt.Printf("Selected file: %s\n", selectedFile)
				input3.SetText(selectedFile) // Store the selected file in input3 text field
			}
		}, w)
		dialog.Show()
	})

	// Create the text area to display printf outputs
	outputTextArea := widget.NewMultiLineEntry()
	// outputTextArea.Wrapping = fyne.TextWrapBreak
	outputTextArea.SetMinRowsVisible(24)

	outputTextArea.TextStyle.Monospace = true
	// outputTextArea.TextStyle.Italic = true
	outputTextArea.Disable()

	srcFile := ""
	browseButton := widget.NewButton("Open Fastboot ROM Zip", func() {
		dialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
			if err == nil && file != nil {
				selectedFile := file.URI().String()
				srcFile = strings.TrimPrefix(selectedFile, "file://")
				fmt.Printf("Selected file: %s\n", selectedFile)
				input2.SetText(selectedFile)
			}
		}, w)
		dialog.Show()
	})

	// Hide the text input fields and browse button
	input2.Hide()
	browseButton.Hide()
	flashCheckbox0.OnChanged = func(checked bool) {
		if checked {
			input2.Show()
			browseButton.Show()
		} else {
			input2.Hide()
			browseButton.Hide()
			w.Resize(fyne.NewSize(920, 400))
			w.Content().Refresh()

		}
	}

	// Hide the text input fields and browse button2
	input3.Hide()
	browseButton2.Hide()
	flashCheckbox.OnChanged = func(checked bool) {
		if checked {
			input3.Show()
			browseButton2.Show()
		} else {
			input3.Hide()
			browseButton2.Hide()
			w.Resize(fyne.NewSize(920, 400))
			w.Content().Refresh()
			// w.Content().Resize(fyne.NewSize(920, 400))
		}
	}

	// Create the Flash button and add functionality
	submitButton := widget.NewButton("Start Flash", func() {

		fmt.Printf("Cloning ghostrider-reborn Kit to: %s\n", destDir)
		outputTextArea.SetText("Cloning ghostrider-reborn Kit to: " + destDir + "\n")

		// Clone the repository

		_, err2 := git.PlainClone(destDir, false, &git.CloneOptions{
			URL:      repoURL,
			Progress: os.Stdout,
		})
		if err2 != nil {
			fmt.Println("Error cloning repository:", err)
		}
		outputTextArea.SetText(outputTextArea.Text + "Repository cloned successfully" + "\n")

		if flashCheckbox.Checked {

			outputTextArea.SetText(outputTextArea.Text + "Copying Firmware Source ..." + "\n")
			outputTextArea.SetText(outputTextArea.Text + "Copying " + srcFile2 + " " + destDir + "\\firmware.zip")

			// Copy File
			bytesRead, err := os.ReadFile(srcFile2)
			if err != nil {
				log.Fatal(err)
			}

			err = os.WriteFile(destDir+"\\firmware.zip", bytesRead, 0644)

			if err != nil {
				log.Fatal(err)
			}

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Getting Fastboot Info")

			getfastbootinfo(destDir, outputTextArea)

			validatezip(destDir + "\\firmware.zip")

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Flashing Firmware ..." + "\n" + "Please wait ...")

			flashfirmware(destDir, outputTextArea)

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Flashing Completed. Reboot your device now" + "\n")
			outputTextArea.CursorRow = outputTextArea.CursorRow + 1
		}
		if flashCheckbox0.Checked {

			outputTextArea.SetText(outputTextArea.Text + "Copying ROM Source ...")

			// Copy File
			bytesRead, err := os.ReadFile(srcFile)
			if err != nil {
				log.Fatal(err)
			}

			err = os.WriteFile(destDir+"\\aospa.zip", bytesRead, 0644)

			if err != nil {
				log.Fatal(err)
			}

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Getting Fastboot Info")

			getfastbootinfo(destDir, outputTextArea)
			outputTextArea.CursorRow = outputTextArea.CursorRow + 1

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Flashing ROM ..." + "\n" + "Please wait ...")
			outputTextArea.CursorRow = outputTextArea.CursorRow + 1

			flashrom(destDir, outputTextArea)

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Flashing Completed. Reboot your device now" + "\n")
			outputTextArea.CursorRow = outputTextArea.CursorRow + 1
		}
	})

	// Create the Reboot button and add functionality
	submitButton2 := widget.NewButton("Reboot Phone", func() {
		if _, err := os.Stat(destDir + "\\platform-tools-windows\\fastboot.exe"); err == nil || os.IsExist(err) {

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Rebooting Phone ...")

			cmd2, _ := exec.Command(destDir+"\\platform-tools-windows\\fastboot.exe", "reboot").CombinedOutput()

			scanner := bufio.NewScanner(bytes.NewReader(cmd2))
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println(m)
				outputTextArea.SetText(outputTextArea.Text + "\n" + m)
				outputTextArea.CursorRow = outputTextArea.CursorRow + 1
			}
		} else {
			outputTextArea.SetText(outputTextArea.Text + "\n" + "Fastboot not found. Please flash and invoke this command again.")
			return
		}

		// outputTextArea.SetText(outputTextArea.Text + "\n" + "Phone rebooted")
	})

	// Create a vertical layout for the central widget
	outputTextArea.SetText("Waiting for logs ...")

	buttonContainer := container.NewHBox(submitButton, submitButton2)
	buttonContainer.Layout = layout.NewHBoxLayout()
	// layout := container.NewVBox(flashCheckbox, input3, browseButton2, flashCheckbox0, input2, browseButton, buttonContainer, outputTextArea)
	layout := container.NewBorder(container.NewVBox(flashCheckbox, input3, browseButton2, flashCheckbox0, input2, browseButton), outputTextArea, nil, nil, container.NewVBox(buttonContainer))

	// Set the layout as the content of the window
	w.Resize(fyne.NewSize(920, 400))
	w.SetContent(layout)
	w.CenterOnScreen()

	// Show the main window
	w.ShowAndRun()
}
