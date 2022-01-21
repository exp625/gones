package input

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"math"
	"os"
)

type FileBindings struct {
	Groups FileBindingGroups `json:"groups"`
}
type FileBindingGroup map[BindingName]FileBinding
type FileBindingGroups map[GroupName]FileBindingGroup
type FileBinding struct {
	BoundKey              string `json:"key"`
	BoundControllerButton string `json:"button"`
	BoundControllerAxis   string `json:"axis"`
}

func (b *Bindings) convertToFileFormat() *FileBindings {
	file := &FileBindings{}
	file.Groups = make(map[GroupName]FileBindingGroup, 0)
	for group, _ := range b.Groups {
		file.Groups[group] = make(map[BindingName]FileBinding, 0)
		for binding, _ := range b.Groups[group] {
			key := ""
			button := ""
			axis := ""

			if b.Groups[group][binding].HasBoundKey {
				key = b.Groups[group][binding].BoundKey.String()
			}

			if b.Groups[group][binding].HasBoundControllerButton {
				button = StandardGamepadButton(b.Groups[group][binding].BoundControllerButton).String()
			}

			if b.Groups[group][binding].HasBoundControllerAxis {
				key = StandardGamepadAxis(b.Groups[group][binding].BoundControllerAxis).String(math.Signbit(b.Groups[group][binding].BoundControllerAxisSign))
			}

			file.Groups[group][binding] = FileBinding{
				BoundKey:              key,
				BoundControllerButton: button,
				BoundControllerAxis:   axis,
			}
		}
	}
	return file
}

func (b *Bindings) loadFromFileFormat(file *FileBindings) {
	for group, _ := range file.Groups {
		for binding, _ := range file.Groups[group] {
			if key, ok := keyNameToKeyCode(file.Groups[group][binding].BoundKey); ok {
				b.Groups[group][binding].HasBoundKey = true
				b.Groups[group][binding].BoundKey = key
			}

			if btn, ok := buttonNameToButtonCode(file.Groups[group][binding].BoundControllerButton); ok {
				b.Groups[group][binding].HasBoundControllerButton = true
				b.Groups[group][binding].BoundControllerButton = btn
			}

			if axis, sign, ok := axisNameToAxisCode(file.Groups[group][binding].BoundControllerAxis); ok {
				b.Groups[group][binding].HasBoundControllerAxis = true
				b.Groups[group][binding].BoundControllerAxis = axis
				b.Groups[group][binding].BoundControllerAxisSign = sign
			}

		}
	}
}

func (b *FileBindings) Save() {
	file, err := json.MarshalIndent(b, "", " ")
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile("bindings.json", file, 0644)
	if err != nil {
		log.Print(err)
		return
	}
}

func (b *FileBindings) Load() {
	if _, err := os.Stat("bindings.json"); errors.Is(err, os.ErrNotExist) {
		return
	}
	bytes, err := ioutil.ReadFile("bindings.json")
	if err != nil {
		log.Print(err)
		return
	}
	err = json.Unmarshal(bytes, b)
	if err != nil {
		log.Print(err)
		return
	}
}
