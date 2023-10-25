package mail

import "fmt"

//generate email content for Email Sender

func MakeEmailForPasswordReset(
	otp string,
	to string,
) (string, string, []string, []string, []string, []string) {
	subject := "Reset Password Otp"
	content := fmt.Sprintf(`
	<h1>Reset password</h1>
	<p>Otp code: %s</p>
	`, otp)
	return subject, content, []string{to}, nil, nil, nil
}
