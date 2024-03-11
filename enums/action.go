package enums

type Action string

const (
	Create Action = "CREATE"
	Rename Action = "RENAME"
	Delete Action = "DELETE"
	Write  Action = "WRITE"
	Chmod  Action = "CHMOD"
)
