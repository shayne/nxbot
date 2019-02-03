package main

import (
	"fmt"
	"os"

	"github.com/jacknx/nxbot/pkg/nxapi"
	"golang.org/x/crypto/ssh/terminal"
)

const usage = `Usage:
%s <nx-ip-port> <nx-user> [nx-pass]

Arguments:
  nx-ip-port: the IP address and port of the Nx server, e.g. 1.2.3.4:7001
  nx-user: the user name for the Nx server
  nx-pass: optionally, the password for the Nx user, prompted otherwise
`

type cameraInfo struct {
	ID   string `json:"cameraId"`
	Name string `json:"cameraName"`
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf(usage, os.Args[0])
		os.Exit(2)
	}
	nxIPPort := os.Args[1]
	nxUser := os.Args[2]
	var nxPass string
	if len(os.Args) == 4 {
		nxPass = os.Args[3]
	} else {
		fmt.Printf("Password: ")
		pw, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		nxPass = string(pw)
	}
	api, err := nxapi.NewAPI(nxIPPort, nxUser, nxPass)
	if err != nil {
		panic(err)
	}
	var infos []cameraInfo
	err = api.GETRequest("getCameraUserAttributesList", &infos)
	if err != nil {
		panic(err)
	}
	for _, cam := range infos {
		fmt.Println("---------------------")
		fmt.Printf("Camera ID: %s\n", cam.ID)
		fmt.Printf("Camera Name: %s\n", cam.Name)
	}
}
