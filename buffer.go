package pattern

import (
	"image"
	"image/color"
	"image/draw"
	"sync"
	"time"
)

// Ensure Buffer implements the image.Image interface.
var _ image.Image = (*Buffer)(nil)

// Ensure Buffer implements the DirtyAware interface.
var _ DirtyAware = (*Buffer)(nil)

// DirtyAware is an interface for patterns that can report if they are dirty.
type DirtyAware interface {
	IsDirty() bool
}

// Buffer acts as a cache for a source image.
type Buffer struct {
	Null
	Expiry
	Source      image.Image
	Cached      image.Image
	LastRefresh time.Time
	mutex       sync.RWMutex
	dirty       bool
}

// IsDirty checks if the buffer needs refreshing.
func (b *Buffer) IsDirty() bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	if b.dirty {
		return true
	}

	if b.Expiry.Expiry > 0 && !b.LastRefresh.IsZero() && time.Since(b.LastRefresh) > b.Expiry.Expiry {
		return true
	}

	if source, ok := b.Source.(DirtyAware); ok {
		return source.IsDirty()
	}

	return false
}

// SetDirty explicitly marks the buffer as dirty.
func (b *Buffer) SetDirty() {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.dirty = true
}

// Refresh updates the cached image from the source.
func (b *Buffer) Refresh() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	bounds := b.Bounds()
	dst := image.NewRGBA(bounds)
	// If the source has different bounds, we might want to respect that,
	// but usually a pattern fills the target.
	// We draw source into dst.
	draw.Draw(dst, bounds, b.Source, bounds.Min, draw.Src)
	b.Cached = dst
	b.LastRefresh = time.Now()
	b.dirty = false
}

// At returns the color of the pixel at (x, y).
func (b *Buffer) At(x, y int) color.Color {
	// We do not hold a read lock here because At is called frequently
	// and we want to avoid contention. We rely on atomic pointer swap or similar if concurrency was high,
	// but here we just check references. Note that b.Cached might be replaced during Refresh.
	// However, usually Patterns are not thread-safe for mutation during read.
	// But Refresh uses a mutex.

	// We use a light read check.
	b.mutex.RLock()
	cached := b.Cached
	isDirty := b.dirty || (b.Expiry.Expiry > 0 && !b.LastRefresh.IsZero() && time.Since(b.LastRefresh) > b.Expiry.Expiry)
	b.mutex.RUnlock()

	// "It has IsDirty() but is a pass through" logic:
	// If dirty or no cache, pass through to source.
	if cached == nil || isDirty {
		// Also check source dirty?
		// The prompt says "It has IsDirty() but is a pass through".
		// This suggests we don't automatically Refresh() on At().
		// If we passed through, we return the source's current value.
		return b.Source.At(x, y)
	}

	// If source is dirty, we also treat as dirty?
	// IsDirty() checks source.IsDirty().
	// If source is dirty, IsDirty() is true.
	if source, ok := b.Source.(DirtyAware); ok && source.IsDirty() {
		return b.Source.At(x, y)
	}

	return cached.At(x, y)
}

// NewBuffer creates a new Buffer pattern.
func NewBuffer(source image.Image, ops ...func(any)) *Buffer {
	b := &Buffer{
		Null: Null{
			bounds: image.Rect(0, 0, 255, 255),
		},
		Source: source,
	}
	for _, op := range ops {
		op(b)
	}
	return b
}
