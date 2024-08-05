package domain

type DataType string

const (
	TextType        DataType = "text"
	BinaryType      DataType = "binary"
	CredentialsType DataType = "credentials"
	BankType        DataType = "bank"
)

type DataID string

type DataContext struct {
	ID     DataID
	UserID UserID
	Meta   string
	Title  string
	Type   DataType
}

type TextData struct {
	Ctx  DataContext
	Data string
}

type BinaryData struct {
	Ctx  DataContext
	Data []byte
}

type Credentials struct {
	Username string
	Password string
}

type CredentialsData struct {
	Ctx  DataContext
	Cred Credentials
}

type Card struct {
	CardNumber string
	CardHolder string
	Cvv        int
}

type BankData struct {
	Ctx  DataContext
	Card Card
}
