package commands

import (
	"bytes"
	"fmt"

	pb "github.com/jukeizu/contacts/api/contacts"
)

func formatContact(contact *pb.Contact) string {
	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("**%s**\n", contact.Name))
	buffer.WriteString(fmt.Sprintf(":house: %s", contact.Address))
	buffer.WriteString(fmt.Sprintf("\n\n:iphone: %s", contact.Phone))
	buffer.WriteString("\n\n\n")

	return buffer.String()
}
