package gpio

var relays = map[string]*pin{
	"1": newPin("18"),
	"2": newPin("24"),
	"3": newPin("23"),
	"4": newPin("25"),
	"5": newPin("17"),
	"6": newPin("27"),
}
