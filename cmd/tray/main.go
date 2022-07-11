package tray

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/signal"

	"github.com/BililiveRecorder/BililiveRecorder-Tray/modules/httpServer"
	"github.com/BililiveRecorder/BililiveRecorder-Tray/modules/systemTray"
	log "github.com/sirupsen/logrus"
)

var configFile = flag.String("config", "config.json", "Configuration file")

type Config struct {
	ListenPort       int    `json:"listenPort"`
	ListenHost       string `json:"listenHost"`
	RecorderPort     int    `json:"recorderPort"`
	RecorderHost     string `json:"recorderHost"`
	RecorderUser     string `json:"recorderUser"`
	RecorderPassword string `json:"recorderPassword"`
}

var config Config

func Main(assets fs.FS) {
	flag.Parse()
	readConfig()
	initLogger()
	log.Info("Config: ", config)
	log.Info("Staring...")
	mainChan := make(chan os.Signal)
	httpChan := make(chan os.Signal)
	go httpServer.Main(assets, config.ListenPort, config.ListenHost, httpChan)
	signal.Notify(mainChan, os.Interrupt)
	go systemTray.Setup(mainChan)
	<-mainChan
	log.Println("Shutdown Server ...")
	httpChan <- os.Interrupt
	systemTray.Quit()
}

func readConfig() {
	file, err := os.ReadFile(*configFile)
	if err != nil {
		switch err.(type) {
		case *os.PathError:
			fmt.Println("Config file not found. Creating new one.")
			loadDefaultConfig()
			saveConfig()
			return
		default:
			panic(err)
		}
	}
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println("Config file is not valid JSON. Creating new one.")
		loadDefaultConfig()
		saveConfig()
		return
	}
	fmt.Println("Config file loaded.")
}

func loadDefaultConfig() {
	config.ListenPort = 8687
	config.ListenHost = "127.0.0.1"
	config.RecorderPort = 8686
	config.RecorderHost = "127.0.0.1"
	config.RecorderUser = "admin"
	config.RecorderPassword = "admin"
}

func saveConfig() {
	file, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(*configFile, file, 0644)
	if err != nil {
		panic(err)
	}
}

func initLogger() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}
