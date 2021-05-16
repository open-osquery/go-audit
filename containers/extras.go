package containers

import (
	"github.com/golang/glog"
	"github.com/open-osquery/auditrd"
	"github.com/spf13/viper"
)

var extraParserConstructors = []func(config *viper.Viper) (ExtraParser, error){}

type ExtraParser interface {
	Parse(am *auditrd.AuditMessage)
}

type ExtraParsers []ExtraParser

func RegisterExtraParser(constructor func(config *viper.Viper) (ExtraParser, error)) {
	extraParserConstructors = append(extraParserConstructors, constructor)
}

func createExtraParsers(config *viper.Viper) ExtraParsers {
	var extraParsers ExtraParsers
	for _, constructor := range extraParserConstructors {
		cp, err := constructor(config)
		if err != nil {
			glog.Fatalf("Failed to create ExtraParser: %v", err)
		}
		if cp != nil {
			extraParsers = append(extraParsers, cp)
		}
	}
	return extraParsers
}

func (ps ExtraParsers) Parse(am *auditrd.AuditMessage) {
	for _, p := range ps {
		p.Parse(am)
	}
}
