package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"render-app/config"
	"strings"
	"time"
)

const (
	emptyStatus = "Empty"
)

type ServerResponse struct {
	Status       string `json:"status"`
	Format       string `json:"format"`
	Resolution   string `json:"resolution"`
	DownloadLink string `json:"download_link"`
}

func main() {
	log.Println("setting up config")
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get cwd")
	}
	configPath := path.Join(wd, "config.yaml")
	cfg := config.MustLoad(configPath)

	log.Println("setting up python path")
	pythonPath := path.Join(wd, "python", "render.py")

	for {
		fmt.Println("requesting new order")
		resp, err := getOrder(cfg)
		if err != nil {
			fmt.Println(err)
			time.Sleep(cfg.SleepTime)
			continue
		}

		log.Println("got new order")

		linkArr := strings.Split(resp.DownloadLink, string(os.PathSeparator))
		uid := linkArr[len(linkArr)-4]
		filename := linkArr[len(linkArr)-1]
		linkFilename := strings.Split(filename, ".")[0]
		imageName := fmt.Sprintf("%s.%s", linkFilename, resp.Format)

		log.Println("downloading file")

		err = downloadFile(resp.DownloadLink, filename)
		if err != nil {
			log.Println(err)
			updateStatus(cfg, uid, linkFilename, cfg.UpdateStatus.Error)
			time.Sleep(cfg.SleepTime)
			continue
		}

		log.Println("updating status: IN PROGRESS")

		err = updateStatus(cfg, uid, linkFilename, cfg.UpdateStatus.InProgress)
		if err != nil {
			log.Println(err)
			log.Println("updating status: ERROR")
			updateStatus(cfg, uid, linkFilename, cfg.UpdateStatus.Error)
			time.Sleep(cfg.SleepTime)
			continue
		}

		log.Println("running blender")

		err = runBlender(cfg, filename, pythonPath)
		if err != nil {
			log.Println(err)
			log.Println("failed to render scene")
			// os.Remove(filename)
			// os.Remove(fmt.Sprintf("frame0000.%s", resp.Format))
			updateStatus(cfg, uid, linkFilename, cfg.UpdateStatus.Error)
			continue
		}

		log.Println("render finished")
		os.Rename(fmt.Sprintf("frame0000.%s", resp.Format), imageName)

		status, err := uploadFile(cfg, imageName, uid)
		if err != nil {
			log.Println(err)
			log.Println("failed to send file")
		}

		log.Println("Response code:", status)
		// os.Remove(imageName)
	}
}

func getOrder(cfg *config.Config) (*ServerResponse, error) {
	res, err := http.Get(cfg.BaseURL + "/request")
	if err != nil {
		log.Println("error request")
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status: %d", res.StatusCode)
	}

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("bad server response")
		return nil, err
	}

	var response ServerResponse

	err = json.Unmarshal(resBody, &response)
	if err != nil {
		fmt.Println("bad server response")
		return nil, err
	}

	if response.Status == emptyStatus {
		fmt.Println("no orders")
		return nil, errors.New("no orders")
	}

	return &response, err
}

func downloadFile(url, filename string) error {
	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	res, err := http.Get(url)
	if err != nil {
		fmt.Println("error request")
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %d", res.StatusCode)
	}

	_, err = io.Copy(out, res.Body)
	if err != nil {
		fmt.Println("failed to write file")
		return err
	}

	return nil
}

func updateStatus(cfg *config.Config, uid, linkFilename, status string) error {
	req, err := http.NewRequest(
		http.MethodPut,
		fmt.Sprintf("%s/%s/blend/update/%s/%s", cfg.BaseURL, uid, linkFilename, status),
		nil)
	if err != nil {
		fmt.Println("failed to create request")
		return err
	}

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		fmt.Println("failed to do request")
		return err
	}

	return nil
}

func runBlender(cfg *config.Config, filename, script string) error {
	cmd := exec.Command(fmt.Sprintf("%s %s --python %s --render-output //frame --render-frame 0", cfg.BlenderPath, filename, script))
	return cmd.Run()
}

func uploadFile(cfg *config.Config, filename, uid string) (int, error) {
	postURL := cfg.BaseURL + fmt.Sprintf("/%s/image/upload", uid)

	wd, err := os.Getwd()
	if err != nil {
		return 0, err
	}
	filePath := path.Join(wd, filename)

	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	buf := &bytes.Buffer{}
	w := multipart.NewWriter(buf)
	defer w.Close()

	part, err := w.CreateFormFile("uploadfile", filepath.Base(file.Name()))
	if err != nil {
		return 0, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return 0, err
	}

	r, err := http.NewRequest("POST", postURL, buf)
	if err != nil {
		return 0, err
	}
	r.Header.Add("Content-Type", w.FormDataContentType())

	client := &http.Client{}
	res, err := client.Do(r)

	return res.StatusCode, err
}
