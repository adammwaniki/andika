package views

type Plan struct {
	Title    string
	Price    string
	Benefits []string
	OffsetClass   string
}

// Capital P because it needs to be exportable
var Plans = []Plan{
	{
		Title:    		"Chini",
		Price:    		"$TBA",
		Benefits: 		[]string{"Okane", "Kasegu", "Orera wa suta"},
		OffsetClass:   	"md:top-24",
	},
	{
		Title:    		"Katikati",
		Price:    		"$TBA",
		Benefits: 		[]string{"Chapo nne", "Na ndengu", "Mazee nikona njaa"},
		OffsetClass:   	"md:top-12",
	},
	{
		Title:    		"Juu",
		Price:    		"$TBA",
		Benefits: 		[]string{"Patco mbili", "Chips mwitu", "Smocha kama tano"},
		OffsetClass:   	"md:top-24",
	},
}

