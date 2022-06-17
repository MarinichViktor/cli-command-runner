package cli

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"strings"
)

type VirtualView struct {
	Offset int
	Data   []string
}

func (v *VirtualView) Reset() {
	v.Offset = 0
	v.Data = []string{}
}

func (v *VirtualView) AppendData(d []string) {
	if v.Offset > 0 {
		v.Offset += len(d)
	}

	v.Data = append(v.Data, d...)
}

func (v *VirtualView) VisibleData(xSize int, ySize int) []string {
	dataSize := len(v.Data)

	if dataSize <= ySize {
		return v.Data
	}

	if v.Offset < dataSize {
		startIdx := dataSize - v.Offset - ySize - 1

		if startIdx < 0 {
			v.Offset = dataSize - ySize
			return v.Data[:ySize]
		} else {
			return v.Data[startIdx : startIdx+ySize]
		}
	}

	// todo: to be reviewed
	log.Panicf("Offset: %d bigger than dsize: %d", v.Offset, dataSize)

	return []string{}
}

func (v *VirtualView) ScrollPageUp(view *gocui.View) error {
	_, ySize := view.Size()
	dataSize := len(v.Data)

	if ySize >= dataSize {
		return nil
	}

	availForScroll := dataSize - v.Offset - 1
	if availForScroll > ySize {
		v.Offset += ySize
	} else {
		v.Offset += availForScroll
	}

	if v.Offset > dataSize {
		log.Panicf("Offset: %d bigger than dsize: %d, avail %d", v.Offset, dataSize, availForScroll)
	}

	return v.Draw(view)
}

func (v *VirtualView) ScrollUp(view *gocui.View) error {
	_, ySize := view.Size()
	dataSize := len(v.Data)

	if ySize >= dataSize {
		return nil
	}

	availForScroll := dataSize - ySize - v.Offset - 1
	if availForScroll > 0 {
		v.Offset++
	}

	return v.Draw(view)
}

func (v *VirtualView) ScrollPageDown(view *gocui.View) error {
	_, ySize := view.Size()

	if v.Offset == 0 {
		return nil
	}

	if v.Offset > ySize {
		v.Offset -= ySize
	} else {
		v.Offset -= v.Offset
	}

	return v.Draw(view)
}

func (v *VirtualView) ScrollDown(view *gocui.View) error {
	if v.Offset > 0 {
		v.Offset--
	}

	return v.Draw(view)
}

func (v *VirtualView) Draw(view *gocui.View) error {
	xSize, ySize := view.Size()

	view.Clear()
	_, e := fmt.Fprint(view, strings.Join(v.VisibleData(xSize, ySize), "\n"))

	return e
}
