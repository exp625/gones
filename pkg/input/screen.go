package input

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"math"
	"sort"
)

func (b *Bindings) Draw(screen *ebiten.Image) {
	gray := color.Gray{Y: 20}
	screen.Fill(gray)

	const pad = 20
	lines := []string{
		"<LEFT> to go to the parent\n",
		"<RIGHT> select group or binding\n",
		"<UP>/<DOWN> move selection up or down\n",
		"<ENTER> select new key for the selected input\n",
		"<ESC> close the keybindings\n",
		"<R> reset to defaults",
	}
	helperText := textutil.New(basicfont.Face7x13, screen.Bounds().Dx()-2*pad, len(lines)*basicfont.Face7x13.Height+2*pad, pad, pad, 1)
	for _, line := range lines {
		plz.Just(helperText.WriteString(line))
	}
	helperText.Draw(screen)

	img := ebiten.NewImage(screen.Bounds().Dx()-2*pad, screen.Bounds().Dy()-(len(lines)*basicfont.Face7x13.Height)-3*pad)
	img.Fill(color.Gray{Y: 70})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pad, float64(len(lines)*basicfont.Face7x13.Height+2*pad))
	width, height := ebiten.WindowSize()
	text := textutil.New(basicfont.Face7x13, width, height, 4, 24, 2)

	keys := sortGroupKeys(b.Groups)

	for _, key := range keys {
		if key == string(b.selectedGroup) {
			text.Color(colornames.Green)
			plz.Just(text.WriteString(fmt.Sprintf("%s: \n", key)))
			text.Color(colornames.White)
			keys := sortBindingKeys(b.Groups[b.selectedGroup])
			for _, key := range keys {
				binding := b.Groups[b.selectedGroup][BindingName(key)]
				if key == string(b.selectedBinding) {
					text.Color(colornames.Green)
				}
				if key == string(b.toChange) {
					text.Color(colornames.Yellow)
				}
				plz.Just(text.WriteString(fmt.Sprintf("    %s:\n", binding.Help)))
				plz.Just(text.WriteString(fmt.Sprintf("        <%s>", binding.Key().String())))
				if has, btn := binding.ControllerButton(); has {
					plz.Just(text.WriteString(fmt.Sprintf(", <%s>", StandardGamepadButton(btn).String())))
				} else {
					plz.Just(text.WriteString(fmt.Sprintf(", <%s>", "None")))
				}
				if has, axis, sign := binding.ControllerAxis(); has {
					plz.Just(text.WriteString(fmt.Sprintf(", <%s>\n", StandardGamepadAxis(axis).String(math.Signbit(sign)))))
				} else {
					plz.Just(text.WriteString(fmt.Sprintf(", <%s>\n", "None")))
				}

				text.Color(colornames.White)
			}
		} else {
			plz.Just(text.WriteString(fmt.Sprintf("%s: \n", key)))
		}
	}
	text.Draw(img)
	screen.DrawImage(img, op)
}

func (b *Bindings) MoveSelectionUp() {

	if b.selectedBinding == "" {
		keys := sortGroupKeys(b.Groups)
		for index, key := range keys {
			if key == string(b.selectedGroup) {
				if index == 0 {
					b.selectedGroup = GroupName(keys[len(keys)-1])
				} else {
					b.selectedGroup = GroupName(keys[index-1])
				}
				return
			}

		}
	} else {
		keys := sortBindingKeys(b.Groups[b.selectedGroup])
		for index, key := range keys {
			if key == string(b.selectedBinding) {
				if index == 0 {
					b.selectedBinding = BindingName(keys[len(keys)-1])
				} else {
					b.selectedBinding = BindingName(keys[index-1])
				}
				return
			}
		}
	}
}

func (b *Bindings) MoveSelectionDown() {

	if b.selectedBinding == "" {
		keys := sortGroupKeys(b.Groups)
		for index, key := range keys {
			if key == string(b.selectedGroup) {
				if index == len(keys)-1 {
					b.selectedGroup = GroupName(keys[0])
				} else {
					b.selectedGroup = GroupName(keys[index+1])
				}
				return
			}
		}
	} else {
		keys := sortBindingKeys(b.Groups[b.selectedGroup])
		for index, key := range keys {
			if key == string(b.selectedBinding) {
				if index == len(keys)-1 {
					b.selectedBinding = BindingName(keys[0])
				} else {
					b.selectedBinding = BindingName(keys[index+1])
				}
				return
			}
		}
	}
}

func (b *Bindings) Select() {
	if b.selectedBinding == "" {
		b.selectedBinding = BindingName(sortBindingKeys(b.Groups[b.selectedGroup])[0])
	}
}

func (b *Bindings) ChangeKeybinding() {
	if b.selectedBinding != "" {
		b.toChange = b.selectedBinding
	}
}

func (b *Bindings) ApplyKeybinding(key ebiten.Key, btn ebiten.StandardGamepadButton, axis ebiten.StandardGamepadAxis, sign float64) {
	if b.toChange == "" {
		return
	}
	if key != -1 {
		if b.Groups[b.selectedGroup][b.selectedBinding].DefaultKey != key {
			b.Groups[b.selectedGroup][b.selectedBinding].HasBoundKey = true
			b.Groups[b.selectedGroup][b.selectedBinding].BoundKey = key
		} else {
			b.Groups[b.selectedGroup][b.selectedBinding].HasBoundKey = false
			b.Groups[b.selectedGroup][b.selectedBinding].BoundKey = 0
		}
	}
	if btn != -1 {
		if b.Groups[b.selectedGroup][b.selectedBinding].DefaultControllerButton != btn {
			b.Groups[b.selectedGroup][b.selectedBinding].HasBoundControllerButton = true
			b.Groups[b.selectedGroup][b.selectedBinding].BoundControllerButton = btn
		} else {
			b.Groups[b.selectedGroup][b.selectedBinding].HasBoundControllerButton = false
			b.Groups[b.selectedGroup][b.selectedBinding].BoundControllerButton = 0
		}
	}
	if axis != -1 {
		if b.Groups[b.selectedGroup][b.selectedBinding].DefaultControllerAxis != axis || b.Groups[b.selectedGroup][b.selectedBinding].DefaultControllerAxisSign != sign {
			b.Groups[b.selectedGroup][b.selectedBinding].HasBoundControllerAxis = true
			b.Groups[b.selectedGroup][b.selectedBinding].BoundControllerAxis = axis
			b.Groups[b.selectedGroup][b.selectedBinding].BoundControllerAxisSign = sign
		} else {
			b.Groups[b.selectedGroup][b.selectedBinding].HasBoundControllerAxis = false
			b.Groups[b.selectedGroup][b.selectedBinding].BoundControllerAxis = 0
			b.Groups[b.selectedGroup][b.selectedBinding].BoundControllerAxisSign = 0
		}
	}
	b.toChange = ""
}

func (b *Bindings) Deselect() {
	b.selectedBinding = ""
}

func sortGroupKeys(g BindingGroups) []string {
	keys := make([]string, len(g))
	i := 0
	for k := range g {
		keys[i] = string(k)
		i++
	}
	sort.Strings(keys)
	return keys
}

func sortBindingKeys(g BindingGroup) []string {
	keys := make([]string, len(g))
	i := 0
	for k := range g {
		keys[i] = string(k)
		i++
	}
	sort.Strings(keys)
	return keys
}
