package twitter

import (
	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/sensorbee/sensorbee.v0/data"
	"io/ioutil"
	"os"
	"testing"
)

func TestKeyParameters(t *testing.T) {
	Convey("Given api key parameters", t, func() {
		params, err := data.NewMap(map[string]interface{}{
			"consumer_key":        "abc",
			"consumer_secret":     "def",
			"access_token":        "ghi",
			"access_token_secret": "jkl",
		})
		So(err, ShouldBeNil)

		Convey("when creating apiKey", func() {
			keys, err := getKeyParameters(params)
			So(err, ShouldBeNil)

			Convey("consumer_key should have the correct value", func() {
				So(keys.ConsumerKey, ShouldEqual, "abc")
			})

			Convey("consumer_secret should have the correct value", func() {
				So(keys.ConsumerSecret, ShouldEqual, "def")
			})

			Convey("access_token should have the correct value", func() {
				So(keys.AccessToken, ShouldEqual, "ghi")
			})

			Convey("access_token_secret should have the correct value", func() {
				So(keys.AccessTokenSecret, ShouldEqual, "jkl")
			})
		})

		Convey("when creating a source from it", func() {
			// kind of a white-box testing; ctx and ioParams aren't used in func
			_, err := CreatePublicStreamSource(nil, nil, params)

			Convey("it should succeed", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

func TestKeyFile(t *testing.T) {
	tempPath := func() string {
		f, err := ioutil.TempFile("", "twitter-plugin-test-key-file")
		if err != nil {
			t.Fatal("Cannot create a temp file:", err)
		}
		defer f.Close()

		_, err = f.WriteString(`consumer_key: abc
consumer_secret: def
access_token: ghi
access_token_secret: jkl`)
		if err != nil {
			t.Fatal("Cannot write a key information to the temp file:", err)
		}
		return f.Name()
	}()
	defer os.Remove(tempPath)

	Convey("Given an api key file", t, func() {
		params := data.Map{
			"key_file": data.String(tempPath),
		}

		Convey("when creating apiKey", func() {
			keys, err := getKeyParameters(params)
			So(err, ShouldBeNil)

			Convey("consumer_key should have the correct value", func() {
				So(keys.ConsumerKey, ShouldEqual, "abc")
			})

			Convey("consumer_secret should have the correct value", func() {
				So(keys.ConsumerSecret, ShouldEqual, "def")
			})

			Convey("access_token should have the correct value", func() {
				So(keys.AccessToken, ShouldEqual, "ghi")
			})

			Convey("access_token_secret should have the correct value", func() {
				So(keys.AccessTokenSecret, ShouldEqual, "jkl")
			})
		})

		Convey("when creating a source from it", func() {
			_, err := CreatePublicStreamSource(nil, nil, params)

			Convey("it should succeed", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("when creating apiKey with other key parameters", func() {
			params, err := data.NewMap(map[string]interface{}{
				"key_file":            tempPath,
				"consumer_key":        "_abc",
				"consumer_secret":     "_def",
				"access_token":        "_ghi",
				"access_token_secret": "_jkl",
			})
			So(err, ShouldBeNil)
			keys, err := getKeyParameters(params)
			So(err, ShouldBeNil)

			Convey("key_file parameter should be preferred", func() {
				So(keys.ConsumerKey, ShouldEqual, "abc")
				So(keys.ConsumerSecret, ShouldEqual, "def")
				So(keys.AccessToken, ShouldEqual, "ghi")
				So(keys.AccessTokenSecret, ShouldEqual, "jkl")
			})
		})
	})
}
