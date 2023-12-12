package file_explorer

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FileExplorer struct {
	// Ready indicates that the user has Selected a file/directory. Retrieve the Selected file/directory via Get().
	Ready         bool
	Directory     string
	entries       *[]os.DirEntry
	Selected      int
	selectedCache map[string]int
	wait          int
}

func New() *FileExplorer {
	f := &FileExplorer{
		selectedCache: make(map[string]int),
	}
	return f
}

func (f *FileExplorer) Get() (string, error) {
	if !f.Ready {
		return "", fmt.Errorf("not ready")
	}
	s := (*f.entries)[f.Selected]
	result, err := filepath.Abs(filepath.Join(f.Directory, s.Name()))
	f.Ready = false
	return result, err
}

func (f *FileExplorer) Update() error {
	if f.Ready {
		return nil
	}
	f.updateEntries()
	return nil
}

func (f *FileExplorer) Draw() *ebiten.Image {
	width := 32 * 8 * 4
	height := 30 * 8 * 4
	ret := ebiten.NewImage(width, height)
	ret.Fill(color.Gray{Y: 20})
	const pad = 20

	lines := []string{
		fmt.Sprintf("%s\n", f.Directory),
		"<LEFT> to go to the parent directory\n",
		"<RIGHT> to go into the selected directory\n",
		"<UP>/<DOWN> to browse through the current directory\n",
		"<ENTER> to choose the currently selected file/directory\n",
		"<A-Z>/<0-9> to quickly selected a file/directory starting with the letter/number \n",
		"<ESC> close file explorer",
	}
	text := textutil.New(basicfont.Face7x13, ret.Bounds().Dx()-2*pad, len(lines)*basicfont.Face7x13.Height+2*pad, pad, pad, 1)
	for _, line := range lines {
		plz.Just(text.WriteString(line))
	}
	text.Draw(ret)

	img := ebiten.NewImage(ret.Bounds().Dx()-2*pad, ret.Bounds().Dy()-(len(lines)*basicfont.Face7x13.Height)-3*pad)
	textImage := ebiten.NewImage(ret.Bounds().Dx()-2*pad, ret.Bounds().Dy()-(len(lines)*basicfont.Face7x13.Height)-3*pad)
	img.Fill(color.Gray{Y: 70})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(pad, float64(len(lines)*basicfont.Face7x13.Height+2*pad))
	ret.DrawImage(ebiten.NewImageFromImage(img), op)

	const scale = 2
	entryHeight := scale * basicfont.Face7x13.Height

	maxEntries := img.Bounds().Dy() / entryHeight
	if maxEntries != 0 {
		min := f.Selected - maxEntries/2
		max := f.Selected + maxEntries/2 - 1

		if min < 0 {
			min = 0
		}
		if max > len(*f.entries)-1 {
			max = len(*f.entries) - 1
		}

		text := textutil.New(basicfont.Face7x13, img.Bounds().Dx(), (max-min+1)*entryHeight, 0, 0, scale)
		for i := min; i <= max; i++ {
			text.Color(colornames.White)
			if i == f.Selected {
				text.Color(colornames.Green)
			}
			entry := (*f.entries)[i]
			plz.Just(text.WriteString(fmt.Sprintf("%s %s\n", entry.Type().String(), entry.Name())))
		}

		text.Draw(textImage)

		y := float64(0)
		y += float64(ret.Bounds().Dy() / 2)
		y -= float64(entryHeight / 2)
		y -= float64(f.Selected-min) * float64(entryHeight)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(pad, y)
		ret.DrawImage(textImage, op)
	}

	return ret
}

func (f *FileExplorer) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (f *FileExplorer) Select(directory string) error {
	f.selectedCache[f.Directory] = f.Selected
	absolutePath, err := filepath.Abs(directory)
	if err != nil {
		return err
	}
	f.Directory = absolutePath
	previouslySelected, ok := f.selectedCache[absolutePath]
	if !ok || (ok && previouslySelected > len(*f.entries)-1) {
		previouslySelected = 0
	}
	f.Selected = previouslySelected
	return nil
}

func (f *FileExplorer) OpenFolder() {
	s := (*f.entries)[f.Selected]
	if s.IsDir() {
		if err := f.Select(filepath.Join(f.Directory, s.Name())); err != nil {
			log.Println(err)
		}
	}
}

func (f *FileExplorer) CloseFolder() {
	if err := f.Select(filepath.Dir(f.Directory)); err != nil {
		log.Println(err)
	}
}

func (f *FileExplorer) TextInput(char rune) {

	for index := 0; index < len(*f.entries); index++ {
		entryIndex := (index + f.Selected + 1) % len(*f.entries)
		entry := (*f.entries)[entryIndex]
		if strings.HasPrefix(strings.ToUpper(entry.Name()), string(char)) {
			f.Selected = entryIndex
			break
		}
	}

}

func (f *FileExplorer) updateEntries() {
	entries, err := os.ReadDir(f.Directory)
	if err != nil {
		entries = make([]os.DirEntry, 0)
	}
	e := make([]os.DirEntry, len(entries))
	f.entries = &e

	i := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		(*f.entries)[i] = entry
		i++
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		(*f.entries)[i] = entry
		i++
	}

	l := len(*f.entries) - 1
	if f.Selected > l {
		f.Selected = l
	}
	if f.Selected < 0 {
		f.Selected = 0
	}
}
