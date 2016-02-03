package plugin

import (
	"github.com/sensorbee/twitter"
	"pfi/sensorbee/sensorbee/bql"
)

func init() {
	bql.MustRegisterGlobalSourceCreator("twitter_public_stream",
		bql.SourceCreatorFunc(twitter.CreatePublicStreamSource))
}
