package mail

import (
	"bytes"
	"fmt"
	"log/slog"
	"net/smtp"
	"os"

	"github.com/DillonEnge/service-vault-api/internal/bundles"
	"github.com/euank/gotmpl"
)

func SendMail(
	templatePath string,
	data map[string]string,
	recipients []string,
	subject string,
) error {
	f, err := os.Open(templatePath)
	if err != nil {
		return err
	}

	var b bytes.Buffer

	gotmpl.Template(f, &b, gotmpl.MapLookup(data))

	hostname := "mail.privateemail.com"
	port := "587"

	pw, ok := os.LookupEnv("MAIL_PASSWORD")
	if !ok {
		slog.Error("Failed to load MAIL_PASSWORD from env")
		return bundles.ErrEnvarNotFound
	}

	auth := smtp.PlainAuth("", "no-reply@engehost.net", pw, hostname)

	err = smtp.SendMail(hostname+":"+port, auth, "no-reply@engehost.net", recipients, []byte("From: Vault Engehost <no-reply@engehost.net>\r\n"+
		fmt.Sprintf("To: <%s>\r\n", recipients[0])+
		fmt.Sprintf("Subject: %s\r\n", subject)+
		"Content-Type: text/html\r\n"+
		"MIME-Version: 1.0\r\n"+
		"\r\n"+
		fmt.Sprintf("%s\r\n", string(b.Bytes())),
	))
	if err != nil {
		return err
	}

	return nil
}
