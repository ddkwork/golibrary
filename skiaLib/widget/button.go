package widget

import (
	"fmt"
	"github.com/ddkwork/golibrary/mylog"
	"github.com/richardwilkes/toolbox/log/jot"
	"github.com/richardwilkes/unison"
)

func CreateImageButton(img *unison.Image, actionText string, panel *unison.Panel) *unison.Button {
	btn := unison.NewButton()
	btn.Drawable = img
	btn.ClickCallback = func() { jot.Info(actionText) }
	btn.Tooltip = unison.NewTooltipWithText(fmt.Sprintf("Tooltip for: %s", actionText))
	btn.SetLayoutData(unison.MiddleAlignment)
	panel.AddChild(btn)
	return btn
}
func MustImage(b []byte) *unison.Image {
	image, err := unison.NewImageFromBytes(b, 1)
	if !mylog.Error(err) {
		return nil
	}
	return image
}
