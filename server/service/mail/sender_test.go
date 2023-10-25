package mail

import (
	"testing"

	"github.com/punkzberryz/todo/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	//This will skip the test if we run test with -short command
	if testing.Short() {
		t.Skip()
	}
	config, err := util.LoadConfig("../../.")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"
	content := `
	<h1>Hello world</h1>
	<p>This is a test message from <a href="http://techschool.guru">Tech School</a></p>
	`
	to := []string{"kangtlee90@gmail.com"}
	attachFiles := []string{"../../../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)
}
