package main

import (
	"fmt"
	"os"

	"bitbucket.org/nordcloud/tagmanager/internal/azure"
	"bitbucket.org/nordcloud/tagmanager/internal/rules"

	log "github.com/sirupsen/logrus"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Mapping file not given")
		os.Exit(1)
	}

	t, err := rules.NewRulesFromFile(os.Args[1])

	if err != nil {
		log.WithError(err).Fatalf("can't open rules file: %s", os.Args[1])
	}

	tagger, err := azure.NewAzureTagger(t)

	if err != nil {
		log.WithError(err).Fatal("Can't create tagger")
	}

	scanner := azure.ResourceGroupScanner{Session: tagger.Session}
	res, err := scanner.GetResources()

	if err != nil {
		log.WithError(err).Fatalf("can't scan resources")
	}

	err = tagger.EvaluteRules(&res)
	if err != nil {
		log.WithError(err).Fatal("can't eval rules")
	}

	err = tagger.ExecuteActions()

	if err != nil {
		log.WithError(err).Fatal("can't exec actions")
	}
}
