package lockfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"ValorantRPC/pkg/errors"
)

type Lockfile struct {
	ProcessName string
	PID         int
	Port        int
	Password    string
	Protocol    string
}

func getLockfilePath() (string, error) {
	localAppdata := os.Getenv("LOCALAPPDATA")
	if localAppdata == "" {
		return "", errors.New("lockfile.getLockfilePath", errors.ErrEnvironmentVars)
	}

	lockfilePath := filepath.Join(localAppdata, "Riot Games", "Riot Client", "Config", "lockfile")
	return lockfilePath, nil
}

func Read() (*Lockfile, error) {

	lockfilePath, err := getLockfilePath()
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(lockfilePath)
	if err != nil {
		return nil, errors.New("lockfile.Read", errors.ErrLockfileNotFound)
	}

	data := strings.Split(string(content), ":")

	if len(data) < 5 {
		return nil, errors.New("lockfile.Read", errors.ErrLockfileMalformed)
	}

	processName := data[0]
	password := data[3]
	protocol := strings.TrimSpace(data[4])

	pid, err := strconv.Atoi(data[1])
	if err != nil {
		cause := fmt.Errorf("%w: %v", errors.ErrLockfileMalformed, err)
		return nil, errors.New("lockfile.Read", cause)
	}

	port, err := strconv.Atoi(data[2])
	if err != nil {
		cause := fmt.Errorf("%w: %v", errors.ErrLockfileMalformed, err)
		return nil, errors.New("lockfile.Read", cause)
	}

	return &Lockfile{
		ProcessName: processName,
		PID:         pid,
		Port:        port,
		Password:    password,
		Protocol:    protocol,
	}, nil
}
