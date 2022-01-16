package file_explorer

import (
	"fmt"
	"github.com/exp625/gones/internal/plz"
	"github.com/exp625/gones/internal/textutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
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

func (f *FileExplorer) Draw(screen *ebiten.Image) {
	const scale = 2
	entryHeight := scale * basicfont.Face7x13.Height
	img := ebiten.NewImage(screen.Bounds().Dx(), screen.Bounds().Dy())

	maxEntries := img.Bounds().Dy() / entryHeight
	if maxEntries == 0 {
		return
	}

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

	text.Draw(img)

	y := float64(0)
	y += float64(screen.Bounds().Dy() / 2)
	y -= float64(entryHeight / 2)
	y -= float64(f.Selected-min) * float64(entryHeight)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, y)

	screen.DrawImage(img, op)
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

func repeatingKeyPressed(key ebiten.Key) bool {
	const (
		delay    = 15
		interval = 3
	)
	d := inpututil.KeyPressDuration(key)
	if d == 1 {
		return true
	}
	if d >= delay && (d-delay)%interval == 0 {
		return true
	}
	return false
}
