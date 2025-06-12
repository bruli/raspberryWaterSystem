package gpio

var relays = map[string]*pin{
	"18": newPin("18"),
	"17": newPin("17"),
	"23": newPin("23"),
	"24": newPin("24"),
}
