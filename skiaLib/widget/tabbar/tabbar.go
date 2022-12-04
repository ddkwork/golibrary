package tabbar

import "github.com/richardwilkes/unison"

type (
	Interface interface{ unison.Dockable }
	object    struct {
		unison.Panel
		Text  string
		Tip   string
		Color unison.Ink
	}
)

func New(title, tip string, background unison.Ink) Interface {
	d := &object{
		Text:  title,
		Tip:   tip,
		Color: background,
	}
	d.Self = d
	d.DrawCallback = d.draw
	d.GainedFocusCallback = d.MarkForRedraw
	d.LostFocusCallback = d.MarkForRedraw
	d.MouseDownCallback = d.mouseDown
	d.SetFocusable(true)
	d.SetSizer(func(_ unison.Size) (min, pref, max unison.Size) {
		pref.Width = 200
		pref.Height = 100
		return min, pref, unison.MaxSize(max)
	})
	return d
}

func (d *object) draw(gc *unison.Canvas, rect unison.Rect) {
	gc.DrawRect(rect, d.Color.Paint(gc, rect, unison.Fill))
	if d.Focused() {
		txt := unison.NewText("Focused", &unison.TextDecoration{
			Font:       unison.EmphasizedSystemFont,
			Foreground: unison.Black,
		})
		r := d.ContentRect(false)
		size := txt.Extents()
		txt.Draw(gc, r.X+(r.Width-size.Width)/2, r.Y+(r.Height-size.Height)/2+txt.Baseline())
	}
}

func (d *object) mouseDown(where unison.Point, button, clickCount int, mod unison.Modifiers) bool {
	if !d.Focused() {
		d.RequestFocus()
		d.MarkForRedraw()
	}
	return true
}

func (d *object) TitleIcon(suggestedSize unison.Size) unison.Drawable { return nil }
func (d *object) Title() string                                       { return d.Text }
func (d *object) Tooltip() string                                     { return d.Tip }
func (d *object) Modified() bool                                      { return false }
