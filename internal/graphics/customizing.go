package graphics

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

// Кастомная тема
type CustomTheme struct{}

func (CustomTheme) Font(s fyne.TextStyle) fyne.Resource {
	// Подключение пользовательского шрифта
	if s.Bold {
		return theme.DefaultTextBoldFont()
	}
	return theme.DefaultTextFont()
}

func (CustomTheme) Size(n fyne.ThemeSizeName) float32 {
	// Размеры стандартных элементов интерфейса
	switch n {
	case theme.SizeNamePadding:
		return 8 // Отступы
	default:
		return theme.DefaultTheme().Size(n)
	}
}

func (CustomTheme) Color(n fyne.ThemeColorName, v fyne.ThemeVariant) color.Color {
	// Цвета элементов интерфейса
	return theme.DefaultTheme().Color(n, v)
}

func (CustomTheme) Icon(n fyne.ThemeIconName) fyne.Resource {
	// Иконки
	return theme.DefaultTheme().Icon(n)
}

// Квадратная кнопка
type SquareButton struct {
	widget.BaseWidget
	Text     string
	OnTapped func()
}

func NewSquareButton(text string, tapped func()) *SquareButton {
	btn := &SquareButton{Text: text, OnTapped: tapped}
	btn.ExtendBaseWidget(btn)
	return btn
}

func (b *SquareButton) CreateRenderer() fyne.WidgetRenderer {
	label := canvas.NewText(b.Text, theme.ForegroundColor())
	label.Alignment = fyne.TextAlignCenter
	rect := canvas.NewRectangle(theme.ButtonColor())
	return &squareButtonRenderer{
		button: b,
		label:  label,
		rect:   rect,
	}
}

type squareButtonRenderer struct {
	button *SquareButton
	label  *canvas.Text
	rect   *canvas.Rectangle
}

func (r *squareButtonRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
	r.label.Move(fyne.NewPos(
		size.Width/2-r.label.Size().Width/2,
		size.Height/2-r.label.Size().Height/2,
	))
}

func (r *squareButtonRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 40) // Минимальный размер кнопки
}

func (r *squareButtonRenderer) Refresh() {
	r.label.Text = r.button.Text
	canvas.Refresh(r.label)
}

func (r *squareButtonRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.rect, r.label}
}

func (r *squareButtonRenderer) Destroy() {}

/*func main() {
	a := app.New()

	w := a.NewWindow("Customized Fyne App")
	w.Resize(fyne.NewSize(400, 400))

	// Фон экрана
	background := canvas.NewRectangle(color.RGBA{R: 220, G: 220, B: 255, A: 255})

	// Кнопка
	btn := NewSquareButton("Login", func() {
		w.SetContent(container.NewVBox(
			canvas.NewText("Welcome!", theme.ForegroundColor()),
		))
	})

	w.SetContent(container.NewStack(
		background,
		container.NewVBox(
			canvas.NewText("Login Screen", theme.ForegroundColor()),
			btn,
		),
	))
	w.ShowAndRun()
}
*/
