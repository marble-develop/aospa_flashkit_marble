package main

import (
	"bufio"
	"fmt"

	"bytes"
	"os"
	"os/exec"

	"strings"

	"archive/zip"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"gopkg.in/src-d/go-git.v4"
)

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
	}
}

// Function to get fastboot info
func getfastbootinfo(destDir string, outputTextArea *widget.Entry) {
	var outputBuffer2 bytes.Buffer

	cmd2 := exec.Command(destDir+"/platform-tools-windows/fastboot", "getvar", "all")
	cmd2.Dir = destDir

	cmd2.Stdout = &outputBuffer2
	cmd2.Stderr = &outputBuffer2
	err3 := cmd2.Start()
	if err3 != nil {
		fmt.Println("Error starting command:", err3)
		return
	}
	cmd2.Wait()

	scanner2 := bufio.NewScanner(bytes.NewReader(outputBuffer2.Bytes()))
	scanner2.Split(bufio.ScanLines)

	for scanner2.Scan() {
		m := scanner2.Text()
		fmt.Println(m)
		outputTextArea.SetText(outputTextArea.Text + "\n" + m)
		outputTextArea.CursorRow = outputTextArea.CursorRow + 1
	}
}

// Function to flash ROM
func flashrom(destDir string, outputTextArea *widget.Entry) {
	var outputBuffer2 bytes.Buffer
	cmd3 := exec.Command("bash", "flash_aospa.sh")
	cmd3.Dir = destDir

	cmd3.Stdout = &outputBuffer2
	cmd3.Stderr = &outputBuffer2
	err3 := cmd3.Start()
	if err3 != nil {
		fmt.Println("Error starting command:", err3)
		return
	}
	cmd3.Wait()

	scanner2 := bufio.NewScanner(bytes.NewReader(outputBuffer2.Bytes()))
	scanner2.Split(bufio.ScanLines)

	for scanner2.Scan() {
		m := scanner2.Text()
		fmt.Println(m)
		outputTextArea.SetText(outputTextArea.Text + "\n" + m)
		outputTextArea.CursorRow = outputTextArea.CursorRow + 1
	}
}

// Function to flash firmware
func flashfirmware(destDir string, outputTextArea *widget.Entry) {
	var outputBuffer2 bytes.Buffer
	// lightTheme := theme.LightTheme()
	// app.Settings().SetTheme(lightTheme)
	cmd3 := exec.Command("bash", "flash_firmware.sh")
	cmd3.Dir = destDir

	cmd3.Stdout = &outputBuffer2
	cmd3.Stderr = &outputBuffer2
	err3 := cmd3.Start()
	if err3 != nil {
		fmt.Println("Error starting command:", err3)
		return
	}
	cmd3.Wait()

	scanner2 := bufio.NewScanner(bytes.NewReader(outputBuffer2.Bytes()))
	scanner2.Split(bufio.ScanLines)

	for scanner2.Scan() {
		m := scanner2.Text()
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
	destDir := filepath.Dir(ex) + "/aospa-flashing-kit"

	a := app.New()
	w := a.NewWindow("AOSPA[Marble] Flash Tool - NooB Edition")

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
	outputTextArea.Enable()

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
			return
		}
		outputTextArea.SetText(outputTextArea.Text + "Repository cloned successfully" + "\n")

		if flashCheckbox.Checked {

			outputTextArea.SetText(outputTextArea.Text + "Copying Firmware Source ...")

			// Copy File
			cmd := exec.Command("cp", "-v", srcFile2, destDir+"/firmware.zip")
			stderr, _ := cmd.StdoutPipe()
			cmd.Start()

			scanner := bufio.NewScanner(stderr)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println(m)
				outputTextArea.SetText(outputTextArea.Text + "\n" + m)
				outputTextArea.CursorRow = outputTextArea.CursorRow + 1

			}
			cmd.Wait()

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Getting Fastboot Info")

			getfastbootinfo(destDir, outputTextArea)
			
			validatezip(destDir + "/firmware.zip")

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Flashing Firmware ..." + "\n" + "Please wait ...")

			flashfirmware(destDir, outputTextArea)

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Flashing Completed. Reboot your device now" + "\n")
		}
		if flashCheckbox0.Checked {

			outputTextArea.SetText(outputTextArea.Text + "Copying ROM Source ...")

			// Copy File
			cmd := exec.Command("cp", "-v", srcFile, destDir+"/aospa.zip")
			stderr, _ := cmd.StdoutPipe()
			cmd.Start()

			scanner := bufio.NewScanner(stderr)
			scanner.Split(bufio.ScanLines)

			for scanner.Scan() {
				m := scanner.Text()
				fmt.Println(m)
				outputTextArea.SetText(outputTextArea.Text + "\n" + m)
				outputTextArea.CursorRow = outputTextArea.CursorRow + 1

			}
			cmd.Wait()

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Getting Fastboot Info")

			getfastbootinfo(destDir, outputTextArea)

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Flashing ROM ..." + "\n" + "Please wait ...")

			flashrom(destDir, outputTextArea)

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Flashing Completed. Reboot your device now" + "\n")
		}
	})

	// Create the Reboot button and add functionality
	submitButton2 := widget.NewButton("Reboot Phone", func() {
		outputTextArea.SetText(outputTextArea.Text + "\n" + "Rebooting Phone ...")

		cmd2, _ := exec.Command("fastboot", "reboot").CombinedOutput()

		scanner := bufio.NewScanner(bytes.NewReader(cmd2))
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
			outputTextArea.SetText(outputTextArea.Text + "\n" + m)
			outputTextArea.CursorRow = outputTextArea.CursorRow + 1
		}

		outputTextArea.SetText(outputTextArea.Text + "\n" + "Phone rebooted")
	})

	// Create a vertical layout for the central widget
	outputTextArea.SetText("Waiting for logs ...")

	buttonContainer := container.NewHBox(submitButton, submitButton2)
	buttonContainer.Layout = layout.NewHBoxLayout()
	layout := container.NewVBox(flashCheckbox, input3, browseButton2, flashCheckbox0, input2, browseButton, buttonContainer, outputTextArea)

	// Set the layout as the content of the window
	w.Resize(fyne.NewSize(920, 600))
	w.SetContent(layout)
	w.CenterOnScreen()

	// Show the main window
	w.ShowAndRun()
}
