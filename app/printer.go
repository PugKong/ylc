package app

type Printer interface {
	Println(args ...any)
	Printf(format string, args ...any)
}
