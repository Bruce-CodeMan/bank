package mail

import (
	"testing"

	"github.com/BruceCompiler/bank/utils"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := utils.LoadConfig()
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email"

	content := `<h1>Hello world</h1><p>This is a test message</p>`

	to := []string{"1033684650@qq.com"}
	attachFiles := []string{"../readme.md"}

	err = sender.SendEmail(
		subject,
		content,
		to,
		nil,
		nil,
		attachFiles,
	)
	require.NoError(t, err)
}
