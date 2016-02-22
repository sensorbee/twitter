package plugin

import (
	"github.com/sensorbee/twitter"
	"gopkg.in/sensorbee/sensorbee.v0/bql"
)

func init() {
	bql.MustRegisterGlobalSourceCreator("twitter_public_stream",
		bql.SourceCreatorFunc(twitter.CreatePublicStreamSource))
}
