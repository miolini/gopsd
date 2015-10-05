package gopsd

import (
	"encoding/json"
	"errors"
	"github.com/miolini/gopsd/util"
	"image"
	"io"
	"io/ioutil"
)

// TODO all INT -> INT64 (**PSB**)
// TODO make([]interface{}, 0) -> var name []interface{}
type Document struct {
	IsLarge bool `json:"is_large"`

	Channels  int16       `json:"channels"`
	Height    int32       `json:"height"`
	Width     int32       `json:"width"`
	Depth     int16       `json:"depth"`
	ColorMode string      `json:"color_mode"`
	Image     image.Image `json:"-"`

	Resources map[int16]interface{} `json:"-"`
	Layers    []*Layer              `json:"layers"`
}

var (
	reader *util.Reader
)

func (d *Document) GetLayersByName(name string) []*Layer {
	var layers []*Layer
	for _, layer := range d.Layers {
		if layer.Name == name {
			layers = append(layers, layer)
		}
	}
	return layers
}

func (d *Document) GetLayerByID(id int) *Layer {
	for _, layer := range d.Layers {
		if layer.ID == int32(id) {
			return layer
		}
	}
	return nil
}

func (d *Document) GetLayer(index int) *Layer {
	if index >= len(d.Layers) {
		return nil
	}
	return d.Layers[index]
}

func (d *Document) ToJSON() ([]byte, error) {
	return json.Marshal(d)
}

func ParseFromBuffer(buffer []byte) (doc *Document, err error) {
	defer func() {
		if r := recover(); r != nil {
			if r == io.EOF {
				err = nil
				return
			}
			switch value := r.(type) {
			case string:
				err = errors.New(value)
			case error:
				err = value
			}
		}
	}()

	reader = util.NewReader(buffer)
	doc = new(Document)
	readHeader(doc)
	readColorMode(doc)

	readResources(doc)
	readLayers(doc)
	readImageData(doc)

	return doc, nil
}

func ParseFromPath(path string) (*Document, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	doc, err := ParseFromBuffer(data)
	if err != nil {
		return nil, err
	}
	return doc, nil
}
