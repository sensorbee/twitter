package twitter

import (
	"fmt"
	_ "github.com/ChimeraCoder/anaconda"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"pfi/sensorbee/sensorbee/bql"
	"pfi/sensorbee/sensorbee/core"
	"pfi/sensorbee/sensorbee/data"
)

type publicStream struct {
	keys *apiKey
}

type apiKey struct {
	ConsumerKey       string `json:"consumer_key" yaml:"consumer_key"`
	ConsumerSecret    string `json:"consumer_secret" yaml:"consumer_secret"`
	AccessToken       string `json:"access_token" yaml:"access_token"`
	AccessTokenSecret string `json:"access_token_secret" yaml:"access_token_secret"`
}

func (p *publicStream) GenerateStream(ctx *core.Context, w core.Writer) error {
	// TODO: implement
	return nil
}

func (p *publicStream) Stop(ctx *core.Context) error {
	return nil
}

// CreatePublicStreamSource creates a new source that receives the public stream
// from Twitter's sampling API.
func CreatePublicStreamSource(ctx *core.Context,
	ioParams *bql.IOParams, params data.Map) (core.Source, error) {
	keys, err := getKeyParameters(params)
	if err != nil {
		return nil, err
	}
	return core.ImplementSourceStop(&publicStream{
		keys: keys,
	}), nil
}

func getKeyParameters(params data.Map) (*apiKey, error) {
	// "key_file" parameter is preferred to other key parameters.
	if v, ok := params["key_file"]; ok {
		path, err := data.AsString(v)
		if err != nil {
			return nil, fmt.Errorf("key_file parameter must be a string: %v", v)
		}
		return loadKeyFile(path)
	}

	// look for other key parameters
	keys := &apiKey{}
	keyVars := []*string{&keys.ConsumerKey, &keys.ConsumerSecret,
		&keys.AccessToken, &keys.AccessTokenSecret}
	for i, p := range []string{"consumer_key", "consumer_secret",
		"access_token", "access_token_secret"} {
		v, ok := params[p]
		if !ok {
			return nil, fmt.Errorf("key_file or %v parameter is missing", p)
		}

		k, err := data.AsString(v)
		if err != nil {
			return nil, fmt.Errorf("%v parameter must be a string: %v", p, v)
		}
		*keyVars[i] = k
	}
	return keys, nil
}

func loadKeyFile(path string) (*apiKey, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return loadKey(f)
}

func loadKey(r io.Reader) (*apiKey, error) {
	in, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	keys := &apiKey{}
	if err := yaml.Unmarshal(in, keys); err != nil {
		return nil, err
	}

	keyVars := []string{keys.ConsumerKey, keys.ConsumerSecret,
		keys.AccessToken, keys.AccessTokenSecret}
	for i, p := range []string{"consumer_key", "consumer_secret",
		"access_token", "access_token_secret"} {
		if keyVars[i] == "" {
			return nil, fmt.Errorf("%v is missing in key_file", p)
		}
	}
	return keys, nil
}
