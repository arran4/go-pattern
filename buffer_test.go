package pattern

import (
	"image"
	"image/color"
	"testing"
	"time"
)

type MockSource struct {
	Null
	dirty bool
	color color.Color
}

func (m *MockSource) IsDirty() bool {
	return m.dirty
}

func (m *MockSource) At(x, y int) color.Color {
	return m.color
}

func TestBuffer_IsDirty(t *testing.T) {
	src := &MockSource{
		color: color.Black,
		dirty: false,
	}

	// 1. Basic Buffer
	b := NewBuffer(src)

	// Should be dirty initially because Cached is nil?
	// The implementation of IsDirty only checks Expiry and Source.IsDirty.
	// It doesn't check if Cached is nil.
	// However, At() checks cached == nil.

	if b.IsDirty() {
		// Default expiry is 0, source is not dirty.
		t.Error("New buffer with clean source should not be dirty")
	}

	// 2. Source becomes dirty
	src.dirty = true
	if !b.IsDirty() {
		t.Error("Buffer should be dirty when source is dirty")
	}

	// 3. Expiry
	b = NewBuffer(src, SetExpiry(100*time.Millisecond))
	src.dirty = false
	b.Refresh() // Sets LastRefresh

	if b.IsDirty() {
		t.Error("Buffer should not be dirty immediately after refresh")
	}

	time.Sleep(150 * time.Millisecond)
	if !b.IsDirty() {
		t.Error("Buffer should be dirty after expiry")
	}
}

func TestBuffer_At_Caching(t *testing.T) {
	// Source is red
	src := &MockSource{
		color: color.RGBA{255, 0, 0, 255},
		Null: Null{bounds: image.Rect(0, 0, 10, 10)},
	}

	b := NewBuffer(src)
	b.SetBounds(image.Rect(0, 0, 10, 10))

	// Before Refresh: At should pass through to Source (Red)
	c := b.At(0, 0)
	if r, _, _, _ := c.RGBA(); r != 0xFFFF {
		t.Error("Expected Red from source before refresh")
	}

	// Refresh: Caches Red
	b.Refresh()

	// Change Source to Blue
	src.color = color.RGBA{0, 0, 255, 255}

	// After Refresh, Source changed, but Buffer is not dirty (expiry 0, source not dirty yet)
	// So At should return Cached (Red)
	c = b.At(0, 0)
	if r, _, _, _ := c.RGBA(); r != 0xFFFF {
		t.Error("Expected Cached Red after source change but before dirty")
	}

	// Mark Source as Dirty
	src.dirty = true
	// IsDirty is true. At should pass through to Source (Blue)
	c = b.At(0, 0)
	if _, _, bl, _ := c.RGBA(); bl != 0xFFFF {
		t.Error("Expected Source Blue when dirty")
	}

	// Refresh: Caches Blue
	b.Refresh()
	src.dirty = false

	c = b.At(0, 0)
	if _, _, bl, _ := c.RGBA(); bl != 0xFFFF {
		t.Error("Expected Cached Blue after refresh")
	}
}
