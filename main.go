package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/bwmarrin/discordgo"
	// "regexp"
	"strings"
)

var (
	sess, _         = discordgo.New("Bot MYTOKEN")
  RootWindow      = app.New().NewWindow("Discork!")
	DisplayBody     = container.NewVBox()
	MessageEntry    = widget.NewEntry()
	SendButton      = widget.NewButton("Send", func() { sendMessage("CHANNEL_ID", MessageEntry.Text); MessageEntry.SetText("") }) // Only works with 1 channel at the moment
	DisplayScroller = container.NewVScroll(DisplayBody)
)

// helpful to render images
type imageType struct {
	path string
	w    float32
	h    float32
}

func onMessage(s *discordgo.Session, message *discordgo.MessageCreate) {
	var imagesLinks []imageType
  // store the images URL and dimensions
	for i := 0; i < len(message.Attachments); i++ {
		imagesLinks = append(imagesLinks, imageType{message.Attachments[i].URL, float32(message.Attachments[i].Width), float32(message.Attachments[i].Height)})
	}
  // call the function to add sent/received message to the GUI
	DisplayBody.Add(newMessageWidget(message.Author.AvatarURL(""), message.Author.Username, message.Content, imagesLinks))
	DisplayBody.Add(widget.NewSeparator())
	DisplayScroller.ScrollToBottom() // to always scroll to the latest received message
}

func sendMessage(channel string, message string) {
	sess.ChannelMessageSend(channel, message)
}

func newMessageWidget(profilePicPath string, username string, message string, imageArray []imageType) fyne.CanvasObject {
	// create a profile picture widget
	a, _ := fyne.LoadResourceFromURLString(profilePicPath)
	profilePic := canvas.NewImageFromResource(a)
	profilePic.SetMinSize(fyne.NewSize(50, 50))
	profilePic.Resize(fyne.NewSize(50, 50))

	// create a username widget with padding
	usernameLabel := widget.NewLabelWithStyle(username, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})

	// create a message widget
	messageLabel := widget.NewLabel(wrap(message))

	imagesLabel := container.NewVBox()
	for i := 0; i < len(imageArray); i++ {
		a, _ = fyne.LoadResourceFromURLString(imageArray[i].path)
		img := canvas.NewImageFromResource(a)
		img.SetMinSize(fyne.NewSize(imageArray[i].w/2.5, imageArray[i].h/2.5)) // scaled down 2.5 times
		imagesLabel.Add(img)
	}

	// create a container for the profile picture, username, message and images widgets
	return container.NewHBox(
		container.NewVBox(
			profilePic,
			widget.NewLabel(""),
		),
		widget.NewLabel(""), // add a spacer widget here
		container.NewVBox(
			usernameLabel,
			messageLabel,
			imagesLabel,
		),
	)
}

// this wrap worked the best for me
func wrap(s string) string {
	l := len(s)
	var s1 []string
	for i := 0; i < l; i++ {
		s1 = append(s1, s[i:i+1])
		if (i > 0) && (i%200 == 0) {
			s1 = append(s1, "\n")
		}
	}
	return strings.Join(s1, "")
}

func main() {
	fmt.Println("Program starts here.")
	sess.Identify.Intents = discordgo.IntentsAll
	sess.AddHandler(onMessage)
	sess.Open()
	fmt.Println("The bot is online.")

	MessageEntry.SetPlaceHolder("Type message here...")
	MessageEntry.OnSubmitted = func(c string) { sendMessage("CHANNEL_ID", c); MessageEntry.SetText("") }

	RootWindow.SetContent(container.NewBorder(
		widget.NewLabel("Discork!"), // top
		container.NewBorder(
			nil,
			nil,
			nil,
			SendButton, // bottom right
			MessageEntry, // bottom left
		),
		nil,
		nil,
		DisplayScroller, // right side
	))

	RootWindow.ShowAndRun()
	defer sess.Close()
	fmt.Println("Program ends here.")
}
