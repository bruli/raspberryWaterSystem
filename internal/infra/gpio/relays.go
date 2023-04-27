package gpio

var relays = map[string]*pin{
	"18": newPin("18"),
	"24": newPin("24"),
	"23": newPin("23"),
	"25": newPin("25"),
	"17": newPin("17"),
	"27": newPin("27"),
}
