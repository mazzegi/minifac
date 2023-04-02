package ui

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Pos struct {
	X, Y int
}

func (p Pos) String() string {
	return fmt.Sprintf("%d, %d", p.X, p.Y)
}

func NewPos(x, y int) Pos {
	return Pos{x, y}
}

// Events
type (
	MouseCallback  func(Pos)
	MouseCallbacks []MouseCallback
	KeyCallback    func(ebiten.Key)
	KeyCallbacks   []KeyCallback
)

func (cbs MouseCallbacks) Call(p Pos) {
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
func NewHandler() *Handler {
	h := &Handler{}
	h.mousePos = h.cursorPosition()
	return h
}

type Handler struct {
	mousePos           Pos
	mouseLeftDown      *Pos
	mouseRightDown     *Pos
	keysPressed        []ebiten.Key
	cbsMouseMove       MouseCallbacks
	cbsMouseLeftClick  MouseCallbacks
	cbsMouseRightClick MouseCallbacks
	cbsKeyDown         KeyCallbacks
	cbsKeyUp           KeyCallbacks
}

func (h *Handler) cursorPosition() Pos {
	x, y := ebiten.CursorPosition()
	return Pos{x, y}
}

func (h *Handler) OnMouseMove(cb MouseCallback) {
	h.cbsMouseMove = append(h.cbsMouseMove, cb)
}

func (h *Handler) OnMouseLeftClicked(cb MouseCallback) {
	h.cbsMouseLeftClick = append(h.cbsMouseLeftClick, cb)
}

func (h *Handler) OnMouseRightClicked(cb MouseCallback) {
	h.cbsMouseRightClick = append(h.cbsMouseRightClick, cb)
}

func (h *Handler) OnKeyDown(cb KeyCallback) {
	h.cbsKeyDown = append(h.cbsKeyDown, cb)
}

func (h *Handler) OnKeyUp(cb KeyCallback) {
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

func (h *Handler) Update() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		if h.mouseLeftDown == nil {
			x, y := ebiten.CursorPosition()
			h.mouseLeftDown = &Pos{x, y}
		}
	} else {
		if h.mouseLeftDown != nil {
			h.cbsMouseLeftClick.Call(*h.mouseLeftDown)
			h.mouseLeftDown = nil
		}
	}

	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
		if h.mouseRightDown == nil {
			x, y := ebiten.CursorPosition()
			h.mouseRightDown = &Pos{x, y}
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
