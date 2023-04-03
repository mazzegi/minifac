package eeui

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Events
type (
	MouseCallback  func(image.Point)
	MouseCallbacks []MouseCallback
	KeyCallback    func(ebiten.Key)
	KeyCallbacks   []KeyCallback
)

func (cbs MouseCallbacks) Call(p image.Point) {
	for _, cb := range cbs {
		cb(p)
	}
}

func (cbs KeyCallbacks) Call(k ebiten.Key) {
	for _, cb := range cbs {
		cb(k)
	}
}

// NewHandler
func NewHandler() *EventHandler {
	h := &EventHandler{}
	h.mousePos = h.cursorPosition()
	return h
}

type EventHandler struct {
	mousePos           image.Point
	mouseLeftDown      *image.Point
	mouseRightDown     *image.Point
	keysPressed        []ebiten.Key
	cbsMouseMove       MouseCallbacks
	cbsMouseLeftClick  MouseCallbacks
	cbsMouseRightClick MouseCallbacks
	cbsKeyDown         KeyCallbacks
	cbsKeyUp           KeyCallbacks
}

func (h *EventHandler) cursorPosition() image.Point {
	x, y := ebiten.CursorPosition()
	return image.Point{x, y}
}

func (h *EventHandler) OnMouseMove(cb MouseCallback) {
	h.cbsMouseMove = append(h.cbsMouseMove, cb)
}

func (h *EventHandler) OnMouseLeftClicked(cb MouseCallback) {
	h.cbsMouseLeftClick = append(h.cbsMouseLeftClick, cb)
}

func (h *EventHandler) OnMouseRightClicked(cb MouseCallback) {
	h.cbsMouseRightClick = append(h.cbsMouseRightClick, cb)
}

func (h *EventHandler) OnKeyDown(cb KeyCallback) {
	h.cbsKeyDown = append(h.cbsKeyDown, cb)
}

func (h *EventHandler) OnKeyUp(cb KeyCallback) {
	h.cbsKeyUp = append(h.cbsKeyUp, cb)
}

func keysContain(rs []ebiten.Key, r ebiten.Key) bool {
	for _, er := range rs {
		if er == r {
			return true
		}
	}
	return false
}

func (h *EventHandler) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if h.mouseLeftDown == nil {
			p := image.Pt(ebiten.CursorPosition())
			h.mouseLeftDown = &p
		}
	} else {
		if h.mouseLeftDown != nil {
			h.cbsMouseLeftClick.Call(*h.mouseLeftDown)
			h.mouseLeftDown = nil
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if h.mouseRightDown == nil {
			p := image.Pt(ebiten.CursorPosition())
			h.mouseRightDown = &p
		}
	} else {
		if h.mouseRightDown != nil {
			h.cbsMouseRightClick.Call(*h.mouseRightDown)
			h.mouseRightDown = nil
		}
	}

	mpos := h.cursorPosition()
	if mpos != h.mousePos {
		h.cbsMouseMove.Call(mpos)
		h.mousePos = mpos
	}

	// Keys
	inKeys := inpututil.AppendPressedKeys([]ebiten.Key{})
	for _, ik := range inKeys {
		if !keysContain(h.keysPressed, ik) {
			h.cbsKeyDown.Call(ik)
		}
	}
	for _, ek := range h.keysPressed {
		if !keysContain(inKeys, ek) {
			h.cbsKeyUp.Call(ek)
		}
	}
	h.keysPressed = inKeys
}
