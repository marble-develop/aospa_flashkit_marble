//go:build windows
// +build windows

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net/http"

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

func cloneRepo(repoURL, destDir string, outputTextArea *widget.Entry) {
	_, err := git.PlainClone(destDir, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		fmt.Println("Error cloning repository", err)
	}
	fmt.Printf("Cloning ghostrider-reborn Kit to: %s\n", destDir)
	outputTextArea.SetText("Cloning ghostrider-reborn Kit to: " + destDir + "\n")
	outputTextArea.SetText(outputTextArea.Text + "Repository cloned successfully" + "\n")
}

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

func extractFileFromZip(zipFile, fileName string) error {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.Name == fileName {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			defer rc.Close()

			path := filepath.Dir(ex) + "\\" + fileName
			w, err := os.Create(path)
			if err != nil {
				return err
			}
			defer w.Close()

			_, err = io.Copy(w, rc)
			if err != nil {
				return err
			}

			return nil
		}
	}

	return nil
}

func downloadFile(url, zipFile string) error {
	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the zip file
	out, err := os.Create(zipFile)
	if err != nil {
		return err
	}
	defer out.Close()

	// Copy the downloaded file to the zip file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// Extract all files from the zip file
	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer r.Close()

	return nil
}

func Unzip(source, destination string) error {
	archive, err := zip.OpenReader(source)
	if err != nil {
		return err
	}
	defer archive.Close()
	for _, file := range archive.Reader.File {
		reader, err := file.Open()
		if err != nil {
			return err
		}
		defer reader.Close()
		path := filepath.Join(destination, file.Name)
		// Remove file if it already exists; no problem if it doesn't; other cases can error out below
		_ = os.Remove(path)
		// Create a directory at path, including parents
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
		// If file is _supposed_ to be a directory, we're done
		if file.FileInfo().IsDir() {
			continue
		}
		// otherwise, remove that directory (_not_ including parents)
		err = os.Remove(path)
		if err != nil {
			return err
		}
		// and create the actual file.  This ensures that the parent directories exist!
		// An archive may have a single file with a nested path, rather than a file for each parent dir
		writer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer writer.Close()
		_, err = io.Copy(writer, reader)
		if err != nil {
			return err
		}
	}
	return nil
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
	input4 := widget.NewEntry()

	// Create the checkbox2
	flashCheckbox := widget.NewCheck("Flash Firmware", nil)
	flashCheckbox1 := widget.NewCheck("Custom Kernel", nil)

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
	kernelsrcFile := ""
	browseButton3 := widget.NewButton("Open Kernel Zip", func() {
		dialog := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
			if err == nil && file != nil {
				// Handle the selected file
				selectedFile := file.URI().String()
				kernelsrcFile = strings.TrimPrefix(selectedFile, "file://")
				fmt.Printf("Selected file: %s\n", selectedFile)
				input4.SetText(selectedFile)
			}
		}, w)
		dialog.Show()
	})

	// Create the text area to display printf outputs
	outputTextArea := widget.NewMultiLineEntry()
	outputTextArea.Wrapping = fyne.TextWrapBreak
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

	// Hide the text input fields and browse button
	input4.Hide()
	browseButton3.Hide()
	flashCheckbox1.OnChanged = func(checked bool) {
		if checked {
			input2.Show()
			browseButton.Show()
			input4.Show()
			browseButton3.Show()
		} else {
			input4.Hide()
			browseButton3.Hide()
			input2.Hide()
			browseButton.Hide()
			w.Resize(fyne.NewSize(920, 400))
			w.Content().Refresh()

		}
	}

	// Create the Flash button and add functionality
	submitButton := widget.NewButton("Start Flash", func() {
		cloneRepo(repoURL, destDir, outputTextArea)
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

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Flashing Completed." + "\n")
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

			// getfastbootinfo(destDir, outputTextArea)
			outputTextArea.CursorRow = outputTextArea.CursorRow + 1

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Flashing ROM ..." + "\n" + "Please wait ...")
			outputTextArea.CursorRow = outputTextArea.CursorRow + 1

			flashrom(destDir, outputTextArea)

			outputTextArea.SetText(outputTextArea.Text + "\n" + "Flashing Completed. Reboot your device now!" + "\n")
			outputTextArea.CursorRow = outputTextArea.CursorRow + 1
		}
	})

	// Create the Reboot button and add functionality
	submitButton2 := widget.NewButton("Reboot Phone", func() {
		cloneRepo(repoURL, destDir, outputTextArea)
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

	// Create the Reboot button and add functionality
	submitButton3 := widget.NewButton("Flash Kernel", func() {
		cloneRepo(repoURL, destDir, outputTextArea)
		outputTextArea.SetText(outputTextArea.Text + "\n" + "Extracting boot.img from AOSPA source...")
		outputTextArea.CursorRow = outputTextArea.CursorRow + 1
		extractFileFromZip(srcFile, "boot.img")
		outputTextArea.SetText(outputTextArea.Text + "\n" + "Downloading magiskboot...")
		outputTextArea.CursorRow = outputTextArea.CursorRow + 1
		downloadFile("https://github.com/svoboda18/magiskboot/releases/download/1.0-3/magiskboot.zip", filepath.Dir(ex)+"\\Magisk.zip")
		outputTextArea.SetText(outputTextArea.Text + "\n" + "Extracting magiskboot...")
		outputTextArea.CursorRow = outputTextArea.CursorRow + 1
		Unzip(filepath.Dir(ex)+"\\Magisk.zip", filepath.Dir(ex)+"\\magiskboot")

		outputTextArea.SetText(outputTextArea.Text + "\n" + "Patching boot.img ...")

		cmd := exec.Command(filepath.Dir(ex)+"\\magiskboot\\magiskboot.exe", "unpack", filepath.Dir(ex)+"\\boot.img")
		cmd.Dir = filepath.Dir(ex)

		output, _ := cmd.CombinedOutput()
		scanner := bufio.NewScanner(bytes.NewReader(output))
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
			outputTextArea.SetText(outputTextArea.Text + "\n" + m)
			outputTextArea.CursorRow = outputTextArea.CursorRow + 1
		}
		outputTextArea.SetText(outputTextArea.Text + "\n" + "Extracting kernel Image from kernel source ...")
		extractFileFromZip(kernelsrcFile, "Image")
		outputTextArea.SetText(outputTextArea.Text + "\n" + "Kernel extracted.")
		outputTextArea.CursorRow = outputTextArea.CursorRow + 1

		// Copy File
		bytesRead, err := os.ReadFile(filepath.Dir(ex) + "\\Image")
		if err != nil {
			log.Fatal(err)
		}

		err = os.WriteFile(filepath.Dir(ex)+"\\kernel", bytesRead, 0644)

		if err != nil {
			log.Fatal(err)
		}

		outputTextArea.SetText(outputTextArea.Text + "\n" + "Patching kernel ...")
		cmd = exec.Command(filepath.Dir(ex)+"\\magiskboot\\magiskboot.exe", "repack", filepath.Dir(ex)+"\\boot.img")
		cmd.Dir = filepath.Dir(ex)
		cmd.Run()
		outputTextArea.SetText(outputTextArea.Text + "\n" + "Patching completed.")
		outputTextArea.CursorRow = outputTextArea.CursorRow + 1

		outputTextArea.SetText(outputTextArea.Text + "\n" + "Installing new kernel...")
		outputTextArea.CursorRow = outputTextArea.CursorRow + 1
		cmd = exec.Command(destDir+"\\platform-tools-windows\\fastboot.exe", "reboot", "fastboot")
		cmd.Dir = filepath.Dir(ex)
		cmd.Run()
		cmd2, _ := exec.Command(destDir+"\\platform-tools-windows\\fastboot.exe", "flash", "boot", filepath.Dir(ex)+"\\new-boot.img").CombinedOutput()

		scanner = bufio.NewScanner(bytes.NewReader(cmd2))
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
			outputTextArea.SetText(outputTextArea.Text + "\n" + m)
			outputTextArea.CursorRow = outputTextArea.CursorRow + 1
		}
		outputTextArea.SetText(outputTextArea.Text + "\n" + "Completed. Please reboot your device.")
		outputTextArea.CursorRow = outputTextArea.CursorRow + 1
	})

	// Create a vertical layout for the central widget
	outputTextArea.SetText("Waiting for logs ...")

	buttonContainer := container.NewHBox(submitButton, submitButton2, submitButton3)
	buttonContainer.Layout = layout.NewHBoxLayout()
	// layout := container.NewVBox(flashCheckbox, input3, browseButton2, flashCheckbox0, input2, browseButton, buttonContainer, outputTextArea)
	layout := container.NewBorder(container.NewVBox(flashCheckbox, input3, browseButton2, flashCheckbox0, input2, browseButton, flashCheckbox1, input4, browseButton3), outputTextArea, nil, nil, container.NewVBox(buttonContainer))

	// Set the layout as the content of the window
	w.Resize(fyne.NewSize(920, 400))
	w.SetContent(layout)
	w.CenterOnScreen()

	// Show the main window
	w.ShowAndRun()
}
