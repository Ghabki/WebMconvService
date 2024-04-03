# Transforming MP4 Downloads into WebM

Automate the conversion of newly downloaded MP4 files into WebM format with this efficient background service.

## Deployment
This application is crafted in Go and requires minimal setup. Follow these steps for seamless deployment:
- Download and Install FFmpeg: Visit https://ffmpeg.org/download.html, download the appropriate version for your system, and extract it. Set the Windows environment variable to the bin folder path for easy access to FFmpeg functionalities.
- Install Dependencies: Utilize Go's package management to acquire necessary modules by executing the following commands in your terminal:- install two go modules

```powershell
  go get github.com/kardianos/service
```

```Powershell
  go get github.com/xfrr/goffmpeg/transcoder
```
- Build the Executable: Compile the code to generate the executable file:
```Powershell
  go build .
```



## Usage/Examples
Execute the service for a one-time run without installation (optional):
```powershell
<filename.exe> run
```

Install the service. Remember to start it post-installation:

```powershell
<filename.exe> install
```
```powershell
<filename.exe> start
```

To uninstall the service
```powershell
<filename.exe> uninstall
```



## Additional notes
- Files are converted one by one, ensuring accuracy and speed.
- The original files are removed post-conversion, maintaining a clutter-free environment.
- Checks are conducted every 60 seconds. If a file is detected for conversion, subsequent runs are accelerated until completion, reverting to the standard 60-second interval thereafter.
- Disclaimer: Please note that I am not liable for any consequences resulting from the execution of this code.
