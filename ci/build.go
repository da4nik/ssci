package ci

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
)

func buildImage(imageTag, workdir string) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		log().Debugf("Unable to create docker client: %v", err)
		return err
	}
	defer cli.Close()

	buildOptions := types.ImageBuildOptions{
		ForceRemove: true,
		PullParent:  true,
		Dockerfile:  "Dockerfile",
		Tags:        []string{imageTag},
	}
	resp, err := cli.ImageBuild(context.Background(), makeTarReader(workdir), buildOptions)
	if err != nil {
		return err
	}

	if logrus.GetLevel() == logrus.DebugLevel {
		bts, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(bts))
	}
	resp.Body.Close()

	return nil
}

// makeTarReader - Making memory tar for image build context
func makeTarReader(workdir string) *bytes.Reader {
	buffer := new(bytes.Buffer)

	tw := tar.NewWriter(buffer)
	defer tw.Close()

	filepath.Walk(workdir, func(file string, fi os.FileInfo, err error) error {
		// return on any error
		if err != nil {
			return err
		}

		// create a new dir/file header
		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		// update the name to correctly reflect the desired destination when untaring
		header.Name = strings.TrimPrefix(strings.Replace(file, workdir, "", -1), string(filepath.Separator))

		// write the header
		if err = tw.WriteHeader(header); err != nil {
			return err
		}

		// return on directories since there will be no content to tar
		if fi.Mode().IsDir() {
			return nil
		}

		// open files for taring
		f, err := os.Open(file)
		defer f.Close()
		if err != nil {
			return err
		}

		// copy file data into tar writer
		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		return nil
	})

	return bytes.NewReader(buffer.Bytes())
}
