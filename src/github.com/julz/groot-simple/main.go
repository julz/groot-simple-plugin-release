package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"code.cloudfoundry.org/lager"
	"github.com/cloudfoundry/groot"
	specs "github.com/opencontainers/runtime-spec/specs-go"
)

type driver struct{}

func main() {
	groot.Run(driver{}, os.Args)
}

func (d driver) Unpack(logger lager.Logger, id string, parentId string, tar io.Reader) error {
	f, err := os.OpenFile("/opt/rootfs/unpack-args", os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return err
	}

	defer f.Close()
	f.WriteString(fmt.Sprintf("id <%s>, parentId <%s>", id, parentId))

	dest, err := os.OpenFile("/opt/rootfs/"+id+"-layer.tar", os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	_, err = io.Copy(dest, tar)
	return err
}

func (d driver) Bundle(logger lager.Logger, id string, layerIDs []string) (specs.Spec, error) {
	f, err := os.OpenFile("/opt/rootfs/bundle-args", os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return specs.Spec{}, err
	}

	defer f.Close()
	f.WriteString(fmt.Sprintf("id <%s>, layerIds <%s>", id, strings.Join(layerIDs, ",")))

	return specs.Spec{
		Root: &specs.Root{
			Path: "/opt/rootfs",
		},
	}, nil
}
