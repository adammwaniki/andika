package views

type Plan struct {
	Title    string
	Price    string
	Benefits []string
	OffsetClass   string
}

type FAQ struct {
	Question string
	Answer   string
}

// Capital P because it needs to be exportable
var Plans = []Plan{
	{
		Title:    		"Students",
		Price:    		"$TBA",
		Benefits: 		[]string{"Okane", "Kasegu", "Orera wa suta"},
		OffsetClass:   	"md:top-24",
	},
	{
		Title:    		"Bei ya Mwananchi",
		Price:    		"$TBA",
		Benefits: 		[]string{"Chapo nne", "Na ndengu", "Mazee nikona njaa"},
		OffsetClass:   	"md:top-12",
	},
	{
		Title:    		"Companies",
		Price:    		"$TBA",
		Benefits: 		[]string{"Patco mbili", "Chips mwitu", "Smocha kama tano"},
		OffsetClass:   	"md:top-24",
	},
}

var FAQsData = []FAQ{
	{
		Question: "What is Andika?",
		Answer:   "Andika is a lightweight version control system for personal notes, inspired by Git. It lets you create, edit, snapshot, and restore notes from the terminal or API.",
	},
	{
		Question: "Why not just use Git for notes?",
		Answer:   "Git is powerful, but it comes with overhead. Andika is purpose-built for quick note-taking with simple commands i.e., no commits, staging, or complex branching required.",
	},
	{
		Question: "How are snapshots stored?",
		Answer:   "Each note has its own subdirectory within a vcs_storage directory that we automatically create locally, and every snapshot is stored with a unique hash. This keeps history organized and makes restores easy.",
	},
	{
		Question: "Can I restore a previous version of a note?",
		Answer:   "Yes! You can list snapshots chronologically and restore a note to any previous version.",
	},
	{
		Question: "Is there an API version of this?",
		Answer:   "Absolutely. The core logic is reusable across CLI and HTTP API, so you can integrate Andika into other apps or services.",
	},
	{
		Question: "What are the real-world use cases?",
		Answer:   "Journaling, research notes, writing drafts, tracking ideas, or any workflow where you need a reliable “undo button” for text.",
	},
}

