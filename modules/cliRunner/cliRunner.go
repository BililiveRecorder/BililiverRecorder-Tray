package cliRunner

import (
	"bufio"
	"errors"

	"github.com/UserExistsError/conpty"
)

type CliConfig struct {
	ExecutablePath    string
	Workplace         string
	Bind              string
	HttpBasicUser     string
	TttpBasicPass     string
	EnableFileBrowser bool
	LogLevel          string
	LogFileLevel      string
	CertPemPath       string
	CertKeyPath       string
	CertPfxPath       string
	CertPassword      string
}

type LineListener struct {
	listener func(string)
	id       int
}

type CliRunner struct {
	Config             *CliConfig
	Start              func() error
	Write              func(p []byte) (int, error)
	AddLineListener    func(func(string)) int
	RemoveLineListener func(int)
	Stop               func() error
	pty                *conpty.ConPty
	lineListener       []LineListener
	last4Lines         []string
	lineListenerId     int
}

func Create(config *CliConfig) *CliRunner {
	runner := &CliRunner{
		Config:             config,
		Start:              nil,
		Write:              nil,
		AddLineListener:    nil,
		RemoveLineListener: nil,
		Stop:               nil,
		pty:                nil,
		lineListener:       []LineListener{},
		last4Lines:         []string{},
		lineListenerId:     0,
	}
	runner.Start = func() error {
		commandLine, err := setupCommandLine(config)
		if err != nil {
			return err
		}
		cpty, err := conpty.Start(commandLine)
		if err != nil {
			return err
		}
		runner.pty = cpty
		runner.Write = cpty.Write
		scanner := bufio.NewScanner(cpty)
		go func() {
			for scanner.Scan() {
				line := scanner.Text()
				runner.last4Lines = append(runner.last4Lines, line)
				if len(runner.last4Lines) > 4 {
					runner.last4Lines = runner.last4Lines[1:]
				}
				for _, listener := range runner.lineListener {
					listener.listener(line)
				}
			}
		}()
		return nil
	}
	runner.Write = func(p []byte) (int, error) {
		return 0, errors.New("Cli Runner is not running")
	}
	runner.Stop = func() error {
		if runner.pty == nil {
			return errors.New("pty is nil")
		}
		return runner.pty.Close()
	}
	runner.AddLineListener = func(listener func(string)) int {
		id := runner.lineListenerId
		runner.lineListenerId++
		runner.lineListener = append(runner.lineListener, LineListener{listener, id})
		return id
	}
	runner.RemoveLineListener = func(listenerId int) {
		for i, listener := range runner.lineListener {
			if listener.id == listenerId {
				runner.lineListener = append(runner.lineListener[:i], runner.lineListener[i+1:]...)
				return
			}
		}
	}
	return runner
}

func setupCommandLine(config *CliConfig) (string, error) {
	if config.ExecutablePath == "" {
		return "", errors.New("ExecutablePath is empty")
	}
	commandLine := config.ExecutablePath
	commandLine += " run "
	if config.Workplace != "" {
		commandLine += config.Workplace
	} else {
		return "", errors.New("Workplace is empty")
	}
	if config.Bind != "" {
		commandLine += " -b " + config.Bind
	}
	if config.HttpBasicUser != "" {
		commandLine += " --http-basic-user " + config.HttpBasicUser
	}
	if config.TttpBasicPass != "" {
		commandLine += " --http-basic-pass " + config.TttpBasicPass
	}
	if config.EnableFileBrowser {
		commandLine += " --enable-file-browser=True"
	} else {
		commandLine += " --enable-file-browser=False"
	}
	if config.LogLevel != "" {
		commandLine += " --log-level " + config.LogLevel
	}
	if config.LogFileLevel != "" {
		commandLine += " --log-file-label " + config.LogFileLevel
	}
	if config.CertPemPath != "" {
		commandLine += " --cert-pem-path " + config.CertPemPath
	}
	if config.CertKeyPath != "" {
		commandLine += " --cert-key-path " + config.CertKeyPath
	}
	if config.CertPfxPath != "" {
		commandLine += " --cert-pfx-path " + config.CertPfxPath
	}
	if config.CertPassword != "" {
		commandLine += " --cert-password " + config.CertPassword
	}
	return commandLine, nil
}
