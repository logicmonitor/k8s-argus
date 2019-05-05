package helm

import (
	"bytes"
	"fmt"
	"github.com/howardchn/argus-cli/pkg/conf"
	"log"
	"os/exec"
)

type Client struct {
	conf *conf.LMConf
}

func NewClient(conf *conf.LMConf) *Client {
	return &Client{conf}
}

func (client *Client) Clean() error {
	var err error

	log.Println("deleting argus")
	err = client.deleteArgus()
	if newErr := handleError(err, "argus"); newErr != nil {
		log.Println(newErr)
		return newErr
	} else {
		log.Println("deleted argus")
	}

	log.Println("deleting collectorset-controller")
	err = client.deleteCollectorSetController()
	if newErr := handleError(err, "collectorset-controller"); newErr != nil {
		log.Println(newErr)
		return newErr
	} else {
		log.Println("deleted collectorset-controller")
	}

	return nil
}

func (client *Client) deleteArgus() error {
	return deleteRelease("argus")
}

func (client *Client) deleteCollectorSetController() error {
	return deleteRelease("collectorset-controller")
}

func deleteRelease(name string) error {
	cmd := exec.Command("helm", "delete", name, "--purge")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}

	fmt.Println(out.String())
	fmt.Printf("%s deleted\n", name)

	return nil
}

func handleError(err error, name string) error {
	if err != nil && err.Error() == "exec: \"helm\": executable file not found in $PATH" {
		return fmt.Errorf("delete %s ignored, helm not installed", name)
	} else if err != nil {
		return err
	}

	return nil
}
