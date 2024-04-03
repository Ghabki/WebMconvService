package main

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"time"

	"github.com/kardianos/service"
	"github.com/xfrr/goffmpeg/transcoder"
)

// Define a struct to hold the service configuration
type program struct{}

// Define the interface methods required by the service package
func (p *program) Start(s service.Service) error {
	go p.run()
	return nil
}

func (p *program) Stop(s service.Service) error {
	return nil
}

//
//
//
// Function to process the newest file

func main() {

	envVariable := "ffmpeg.exe"

	// Get the value of the environment variable
	_, exists := os.LookupEnv(envVariable)

	// Check if the value is empty
	if exists {
		panic(fmt.Sprintf("Environment variable %s is not set. You need to download the ffmpeg dependecy for the services go package", envVariable))
	}

	svcConfig := &service.Config{
		Name:        "WebmService",
		DisplayName: "WebmService",
		Description: "Converst mp4 files to Webm in the downloads folder and stores it WebmVideos",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		fmt.Println("Error creating service:", err)
		return
	}

	// If invoked with 'install' flag, install the service
	if len(os.Args) > 1 {
		if os.Args[1] == "install" {
			err = s.Install()
			if err != nil {
				fmt.Println("Error installing service:", err)
				return
			}
			fmt.Println("Service installed successfully.")
			return
		} else if os.Args[1] == "uninstall" {
			err = s.Uninstall()
			if err != nil {
				fmt.Println("Error uninstalling service:", err)
				return
			}
			fmt.Println("Service uninstalled successfully.")
			return
		}
	}

	// If not installing, run the service
	err = s.Run()
	if err != nil {
		fmt.Println("Error running service:", err)
	}
}

// Main function to run the service
func (p *program) run() {
	currentUser, err := user.Current()

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	desktopPath := filepath.Join(currentUser.HomeDir, "Desktop")
	DPath := filepath.Join(currentUser.HomeDir, "Downloads")

	downloadFolderPath := DPath
	FolderPath := filepath.Join(desktopPath, "WebmVideos")
	ProcessFolder := filepath.Join(DPath, "ProcessWEBM")

	err = checkFolder(ProcessFolder)
	err1 := checkFolder(FolderPath)

	if err != nil || err1 != nil {
		return
	}

	for {
		var timeWait int64
		timeWait, err = processMP4Files(downloadFolderPath, FolderPath, ProcessFolder)
		if err != nil {
			return
		}

		time.Sleep(time.Duration(timeWait) * time.Second)
	}
}

func checkFolder(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			return fmt.Errorf("error creating folder: %v", err)
		}
	}
	return nil
}

func processMP4Files(folderPath string, outputFolderPath string, proccessFolder string) (int64, error) {

	d, err := os.Open(folderPath)
	if err != nil {
		fmt.Println("Error opening directory:", err)
		return 0, err
	}
	defer d.Close()

	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return 0, err
	}

	var latestMp4File string

	//code to move latest file and to find the file
	for _, file := range files {
		if file.Mode().IsRegular() && filepath.Ext(file.Name()) == ".mp4" {
			if time.Since(file.ModTime()) < time.Minute*5 {
				latestMp4File = file.Name()
				break
			}
		}
	}
	if latestMp4File == "" {
		return 60, nil
	}

	mp4file := filepath.Join(folderPath, latestMp4File)
	mp4NewFolder := filepath.Join(proccessFolder, latestMp4File)

	fmt.Println(mp4file, mp4NewFolder)
	err = os.Rename(mp4file, mp4NewFolder)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}

	err = convertToWebM(mp4NewFolder, outputFolderPath)
	if err != nil {
		fmt.Println("Error converting file:", err)
		return 0, err
	}
	return 5, nil
}

func convertToWebM(inputFilePath string, outputFolderPath string) error {
	trans := new(transcoder.Transcoder)
	// fmt.Println(inputFilePath)
	// fmt.Println(outputFolderPath)

	// Initialize the transcoder with input and output files
	fileName := "output" + fmt.Sprint(time.Now().Unix()) + ".webm"
	err := trans.Initialize(inputFilePath, filepath.Join(outputFolderPath, fileName))
	if err != nil {
		return err
	}

	// Start the transcoding process
	done := trans.Run(true)
	progress := trans.Output()

	// Print transcoding progress
	for msg := range progress {
		fmt.Println(msg)
	}

	// Wait for transcoding to finish
	err = <-done
	if err != nil {
		return err
	}

	fmt.Println("File converted successfully:", inputFilePath, " to ", fileName)
	os.Remove(inputFilePath)
	return nil
}
