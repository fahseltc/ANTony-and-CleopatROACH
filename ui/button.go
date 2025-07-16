package ui

import (
	"image"
	"image/color"

	"gamejam/fonts"
	"gamejam/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type BtnOptFunc func(*Button)

type Button struct {
	rect image.Rectangle

	text   string
	fonts  *fonts.All
	Hidden bool

	currentImg *ebiten.Image
	defaultImg *ebiten.Image
	pressedImg *ebiten.Image

	OnClick func()
	key     ebiten.Key

	tooltip TooltipInterface
}

//
// NewButton creates a new Button with the given environment and options.
//

func NewButton(fonts *fonts.All, opts ...BtnOptFunc) *Button {
	btn := defaultBtnOpts(fonts)
	for _, opt := range opts {
		opt(&btn)
	}
	return &btn
}

func defaultBtnOpts(fonts *fonts.All) Button {
	defaultWidth := float32(250.0)
	defaultHeight := float32(100.0)
	defaultImg := util.LoadImage("ui/btn/menu-btn.png")
	defaultImg = util.ScaleImage(defaultImg, defaultWidth, defaultHeight)
	pressed := util.LoadImage("ui/btn/menu-btn-pressed.png") // todo pressed
	pressed = util.ScaleImage(pressed, defaultWidth, defaultHeight)
	return Button{
		rect: image.Rectangle{
			Min: image.Point{
				X: 0,
				Y: 0,
			},
			Max: image.Point{
				X: 250,
				Y: 100,
			},
		},
		fonts:      fonts,
		currentImg: defaultImg,
		defaultImg: defaultImg,
		pressedImg: pressed,
		key:        999,
	}
}
func WithText(txt string) BtnOptFunc {
	return func(btn *Button) {
		btn.text = txt
	}
}
func WithRect(rect image.Rectangle) BtnOptFunc {
	return func(btn *Button) {
		btn.rect = rect
		defaultImg := util.LoadImage("ui/btn/menu-btn.png")
		defaultImg = util.ScaleImage(defaultImg, float32(rect.Dx()), float32(rect.Dy()))
		pressed := util.LoadImage("ui/btn/menu-btn-pressed.png")
		pressed = util.ScaleImage(pressed, float32(rect.Bounds().Dx()), float32(rect.Bounds().Dy()))

		btn.currentImg = defaultImg
		btn.defaultImg = defaultImg
		btn.pressedImg = pressed
	}
}
func WithClickFunc(f func()) BtnOptFunc {
	return func(btn *Button) {
		btn.OnClick = f
	}
}
func WithImage(defaultImg *ebiten.Image, pressedImg *ebiten.Image) BtnOptFunc {
	return func(btn *Button) {
		defaultBtn := util.ScaleImage(defaultImg, float32(btn.rect.Bounds().Dx()), float32(btn.rect.Bounds().Dy()))
		btn.currentImg = defaultBtn
		btn.defaultImg = defaultBtn
		btn.pressedImg = util.ScaleImage(pressedImg, float32(btn.rect.Bounds().Dx()), float32(btn.rect.Bounds().Dy()))
	}
}
func WithKeyActivation(key ebiten.Key) BtnOptFunc {
	return func(btn *Button) {
		btn.key = key
	}
}
func WithToolTip(tt TooltipInterface) BtnOptFunc {
	return func(btn *Button) {
		btn.tooltip = tt
		btn.tooltip.GetAlignment().Align(btn.rect, tt.GetRect())
	}
}

// func WithCenteredPos() BtnOptFunc {
// 	return func(btn *Button) {
// 		centeredX := float64(btn.rect.Min.X) - 0.5*float64(btn.rect.Dx())
// 		centeredY := float64(btn.rect.Min.Y) - 0.5*float64(btn.rect.Dy())
// 		btn.rect.Min.X = int(centeredX)
// 		btn.rect.Min.Y = int(centeredY)
// 	}
// }

//
// Class Functions
//

func (btn *Button) Draw(screen *ebiten.Image) {
	if btn.Hidden {
		return
	}
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(btn.rect.Min.X), float64(btn.rect.Min.Y))
	screen.DrawImage(btn.currentImg, op)

	if btn.text != "" {
		// draw text centered
		centerX, centerY := btn.GetCenter()
		if btn.currentImg == btn.pressedImg {
			util.DrawCenteredText(screen, btn.fonts.Med, btn.text, centerX, centerY+4, color.RGBA{R: 0, G: 0, B: 0, A: 255})
		} else {
			util.DrawCenteredText(screen, btn.fonts.Med, btn.text, centerX, centerY, color.RGBA{R: 0, G: 0, B: 0, A: 255})
		}

	}

	if btn.tooltip != nil && btn.MouseCollides() {
		btn.tooltip.OnHover(screen)
	}

	if btn.key != 999 {
		util.DrawCenteredText(screen, btn.fonts.XSmall, btn.key.String(), btn.rect.Min.X+6, btn.rect.Min.Y-4, color.RGBA{R: 255, G: 255, B: 255, A: 255})
	}
}

func (btn *Button) Update() {
	if btn.Hidden {
		return
	}
	// clicks
	if btn.OnClick != nil && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && btn.MouseCollides() {
		btn.currentImg = btn.pressedImg
	}
	if btn.OnClick != nil && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && btn.MouseCollides() {
		btn.OnClick()
		btn.currentImg = btn.defaultImg
	}
	// key presses
	if btn.key != 999 && inpututil.IsKeyJustPressed(btn.key) {
		btn.currentImg = btn.pressedImg
	}
	if btn.key != 999 && inpututil.IsKeyJustReleased(btn.key) {
		btn.OnClick()
		btn.currentImg = btn.defaultImg
	}
}

func (btn *Button) MouseCollides() bool {
	mx, my := ebiten.CursorPosition()
	collides := mx > int(btn.rect.Min.X) &&
		mx < int(btn.rect.Max.X) &&
		my > int(btn.rect.Min.Y) &&
		my < int(btn.rect.Max.Y)
	return collides
}

func (btn *Button) GetCenter() (x, y int) {
	centerX := btn.rect.Min.X + btn.rect.Dx()/2
	centerY := btn.rect.Min.Y + btn.rect.Dy()/2
	return int(centerX), int(centerY)
}
