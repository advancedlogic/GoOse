package types

import (
	"time"
	"github.com/advancedlogic/GoOse/pkg/goose"
)

const defaultUserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_2) AppleWebKit/534.52.7 (KHTML, like Gecko) Version/5.1.2 Safari/534.52.7"

// Configuration is a wrapper for various config options
type Configuration struct {
	localStoragePath        string //not used in this version
	imagesMinBytes          int    //not used in this version
	TargetLanguage          string
	imageMagickConvertPath  string //not used in this version
	imageMagickIdentifyPath string //not used in this version
	BrowserUserAgent        string
	Debug                   bool
	ExtractPublishDate      bool
	AdditionalDataExtractor bool
	EnableImageFetching     bool
	UseMetaLanguage         bool

	//path to the stopwords folder
	stopWordsPath string
	StopWords     goose.StopWords
	Parser        *goose.Parser

	Timeout time.Duration
}

// GetDefaultConfiguration returns safe default configuration options
func GetDefaultConfiguration(args ...string) Configuration {
	if len(args) == 0 {
		return Configuration{
			localStoragePath:        "",   //not used in this version
			imagesMinBytes:          4500, //not used in this version
			EnableImageFetching:     true,
			UseMetaLanguage:         true,
			TargetLanguage:          "en",
			imageMagickConvertPath:  "/usr/bin/convert",  //not used in this version
			imageMagickIdentifyPath: "/usr/bin/identify", //not used in this version
			BrowserUserAgent:        defaultUserAgent,
			Debug:                   false,
			ExtractPublishDate:      true,
			AdditionalDataExtractor: false,
			stopWordsPath:           "resources/stopwords",
			StopWords:               goose.NewStopwords(), //TODO with path
			Parser:                  goose.NewParser(),
			Timeout:                 time.Duration(5 * time.Second),
		}
	}
	return Configuration{
		localStoragePath:        "",   //not used in this version
		imagesMinBytes:          4500, //not used in this version
		EnableImageFetching:     true,
		UseMetaLanguage:         true,
		TargetLanguage:          "en",
		imageMagickConvertPath:  "/usr/bin/convert",  //not used in this version
		imageMagickIdentifyPath: "/usr/bin/identify", //not used in this version
		BrowserUserAgent:        defaultUserAgent,
		Debug:                   false,
		ExtractPublishDate:      true,
		AdditionalDataExtractor: false,
		stopWordsPath:           "resources/stopwords",
		StopWords:               goose.NewStopwords(), //TODO with path
		Parser:                  goose.NewParser(),
		Timeout:                 time.Duration(5 * time.Second),
	}
}
