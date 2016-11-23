package commands

import (
	"github.com/godwhoa/ezbot/ezbot"
	"github.com/godwhoa/wordnet-go"
)

/*
 */
type Def struct {
	ezbot.Command
	dict *wordnet.WordNet
}

func NewDef() *Def {
	def := new(Def)
	def.Pattern = "^.def"
	def.dict = &wordnet.WordNet{}
	def.dict.Init("")
	return def
}

func (d *Def) Execute(user string, msg string, args []string) {
	mtype, word := "", ""
	switch len(args) {
	case 2:
		// .def <word>
		word = args[1]
		mtype = "noun"
	case 3:
		// .def <word> [noun|verb|adj|adv]
		word = args[1]
		mtype = args[2]
	default:
		d.SChan <- "Usage: .def <word> [noun|verb|adj|adv]"
		return
	}
	r, err := d.dict.ByType(word, mtype, 1)
	if err != nil || len(r) == 0 {
		d.SChan <- "No definition. Note: Defaults to type 'noun' and uses wordnet as source."
		return
	}
	for _, def := range r {
		d.SChan <- def.Definition
	}
}
