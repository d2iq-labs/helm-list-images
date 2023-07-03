package pkg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/cheynewallace/tabby"
	"github.com/ghodss/yaml"

	"github.com/d2iq-labs/helm-list-images/pkg/k8s"
)

func (image *Images) render(images []*k8s.Image) error {
	imagesFiltered := image.FilterImagesByRegistries(images)

	if image.JSON {
		return image.ToJSON(imagesFiltered)
	}

	if image.YAML {
		return image.ToYAML(imagesFiltered)
	}

	if image.Table {
		image.toTABLE(imagesFiltered)

		return nil
	}

	image.log.Debug("no format was specified for rendering images, defaulting to list")

	imags, err := GetImagesFromKind(imagesFiltered)
	if err != nil {
		return err
	}

	if image.UniqueImages {
		imags = GetUniqEntries(imags)
	}

	if image.SortImages {
		sort.Stable(sort.StringSlice(imags))
	}

	if _, err := fmt.Fprintf(image.writer, "%s\n", strings.Join(imags, "\n")); err != nil {
		image.log.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			image.log.Fatalln(err)
		}
	}(image.writer)

	return nil
}

func (image *Images) toTABLE(imagesFiltered []*k8s.Image) {
	image.log.Debug("rendering the images in table format since --table is enabled")

	table := tabby.New()
	table.AddHeader("Name", "Kind", "Image")

	for _, img := range imagesFiltered {
		table.AddLine(img.Name, img.Kind, strings.Join(img.Image, ", "))
	}

	table.Print()
}

func (image *Images) ToYAML(imagesFiltered []*k8s.Image) error {
	image.log.Debug("rendering the images in yaml format since --yaml is enabled")

	kindYAML, err := yaml.Marshal(imagesFiltered)
	if err != nil {
		return err
	}

	yamlString := "---" + "\n" + string(kindYAML)

	if _, err = image.writer.WriteString(yamlString); err != nil {
		image.log.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err = writer.Flush()
		if err != nil {
			image.log.Fatalln(err)
		}
	}(image.writer)

	return nil
}

func (image *Images) ToJSON(imagesFiltered []*k8s.Image) error {
	image.log.Debug("rendering the images in json format since --json is enabled")

	kindJSON, _ := json.MarshalIndent(imagesFiltered, " ", " ")

	if _, err := image.writer.Write(kindJSON); err != nil {
		image.log.Fatalln(err)
	}

	defer func(writer *bufio.Writer) {
		err := writer.Flush()
		if err != nil {
			image.log.Fatalln(err)
		}
	}(image.writer)

	return nil
}
