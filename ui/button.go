package ui

import (
	"image"

	"gamejam/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type BtnOptFunc func(*Button)

type Button struct {
	rect image.Rectangle

	text string
	font text.Face

	currentImg *ebiten.Image
	defaultImg *ebiten.Image
	pressedImg *ebiten.Image

	OnClick func()
}

//
// NewButton creates a new Button with the given environment and options.
//

func NewButton(font text.Face, opts ...BtnOptFunc) *Button {
	btn := defaultBtnOpts(font)
	for _, opt := range opts {
		opt(&btn)
	}
	return &btn
}

func defaultBtnOpts(font text.Face) Button {
	defaultWidth := float32(250.0)
	defaultHeight := float32(100.0)
	defaultImg := util.LoadImage("ui/btn/yellow-btn.png")
	defaultImg = util.ScaleImage(defaultImg, defaultWidth, defaultHeight)
	pressed := util.LoadImage("ui/btn/yellow-btn.png") // todo pressed
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
		font:       font,
		currentImg: defaultImg,
		defaultImg: defaultImg,
		pressedImg: pressed,
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
		defaultImg := util.LoadImage("ui/btn/yellow-btn.png")
		defaultImg = util.ScaleImage(defaultImg, float32(rect.Dx()), float32(rect.Dy()))
		pressed := util.LoadImage("ui/btn/yellow-btn.png")
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

//	func WithToolTip(tt TooltipInterface) BtnOptFunc {
//		return func(btn *Button) {
//			btn.ToolTip = tt
//			btn.ToolTip.GetAlignment().Align(btn.rect, tt.GetRect())
//		}
//	}
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
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(btn.rect.Min.X), float64(btn.rect.Min.Y))
	screen.DrawImage(btn.currentImg, op)

	if btn.text != "" {
		// draw text centered
		centerX, centerY := btn.GetCenter()
		util.DrawCenteredText(screen, btn.font, btn.text, centerX, centerY, nil)
	}
	// ebitenutil.DrawRect(screen, float64(btn.rect.Min.X), float64(btn.rect.Min.Y), float64(btn.rect.Dx()), float64(btn.rect.Dy()), color.RGBA{0, 255, 0, 255})
}

func (btn *Button) Update() {
	if btn.OnClick != nil && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && btn.MouseCollides() {
		btn.currentImg = btn.pressedImg
	}
	if btn.OnClick != nil && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) && btn.MouseCollides() {
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
