package gui

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type (
	SubmitImage          func(image.Image) error
	RandomImageGenerator func(width, height int) image.Image
)

type App struct {
	fyneApp             fyne.App
	mainWindow          fyne.Window
	selectedImageCanvas *canvas.Image
	receivedImageCanvas *canvas.Image

	submitImage          SubmitImage
	randomImageGenerator RandomImageGenerator
}

type Option func(*App)

func WithRandomImageGenerator(generator RandomImageGenerator) Option {
	return func(app *App) {
		app.randomImageGenerator = generator
	}
}

func WithSubmitImageHandler(handler SubmitImage) Option {
	return func(app *App) {
		app.submitImage = handler
	}
}

func New(options ...Option) *App {
	fyneApp := app.New()
	mainWindow := fyneApp.NewWindow("Image Sender")

	appInstance := &App{
		fyneApp:             fyneApp,
		mainWindow:          mainWindow,
		selectedImageCanvas: canvas.NewImageFromImage(nil),
		receivedImageCanvas: canvas.NewImageFromImage(nil),
	}

	for _, option := range options {
		option(appInstance)
	}

	appInstance.setupUI()
	return appInstance
}

func (a *App) setupUI() {
	a.selectedImageCanvas.FillMode = canvas.ImageFillContain
	a.selectedImageCanvas.SetMinSize(fyne.NewSize(400, 300))
	a.receivedImageCanvas.FillMode = canvas.ImageFillContain
	a.receivedImageCanvas.SetMinSize(fyne.NewSize(400, 300))

	selectImageButton := widget.NewButton("Select Image", a.selectImageHandler)
	randomImageButton := widget.NewButton("Generate Random Image", a.generateImageHandler)
	submitImageButton := widget.NewButton("Submit Image", a.submitImageHandler)

	separator := canvas.NewLine(color.RGBA{255, 255, 255, 255})
	separator.StrokeWidth = 2

	// Layout
	content := container.NewVBox(
		selectImageButton,
		randomImageButton,
		a.selectedImageCanvas,
		submitImageButton,
		separator,
		a.receivedImageCanvas,
	)

	a.mainWindow.SetContent(content)
	a.mainWindow.Resize(fyne.NewSize(1200, 300))
}

func (a *App) selectImageHandler() {
	fileDialog := dialog.NewFileOpen(
		func(uc fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, a.mainWindow)
				return
			}
			if uc == nil {
				return
			}
			defer uc.Close()

			// Чтение файла
			data, err := ioutil.ReadAll(uc)
			if err != nil {
				dialog.ShowError(err, a.mainWindow)
				return
			}
			img, _, err := image.Decode(bytes.NewReader(data))
			if err != nil {
				dialog.ShowError(err, a.mainWindow)
				return
			}
			a.selectedImageCanvas.Image = img
			a.selectedImageCanvas.Refresh()
		}, a.mainWindow)
	fileDialog.SetFilter(
		storage.NewExtensionFileFilter([]string{".jpg", ".jpeg"}))
	fileDialog.Show()
}

func (a *App) submitImageHandler() {
	if a.submitImage != nil {
		err := a.submitImage(a.selectedImageCanvas.Image)
		if err != nil {
			dialog.ShowError(err, a.mainWindow)
		}
	}
}

func (a *App) generateImageHandler() {
	if a.randomImageGenerator != nil {
		img := a.randomImageGenerator(400, 300)
		a.selectedImageCanvas.Image = img
		a.selectedImageCanvas.Refresh()
	}
}

func (a *App) SetReceivedImage(img image.Image) {
	a.receivedImageCanvas.Image = img
	a.receivedImageCanvas.Refresh()
}

func (a *App) Run() {
	a.mainWindow.ShowAndRun()
}
