package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

// ConfigFile contains information about the site
var ConfigFile string

// SiteInfo contains all the variables for the site
var SiteInfo SiteInformation

// SiteInformation contains the information inside ConfigFile
type SiteInformation struct {
	Author            string `json:Author`
	BlogURL           string `json:BlogURL`
	BlogLanguage      string `json:BlogLanguage`
	BlogDescription   string `json:BlogDescription`
	DateFormat        string `json:DateFormat`
	ThemePath         string `json:TemplatePath`
	BlogTitle         string `json:BlogTitle`
	NumPostsFrontPage int    `json:NumPostsFrontPage`
	DataSource        struct {
		Type       string `json:Type`
		Repository string `json:Repository`
	}
	Upload struct {
		URL      string `json:URL`
		Username string `json:Username`
		Password string `json:Password`
	}
	TempFolder string `json:TempFolder`
	DestFolder string `json:DestFolder`
}

func NewSiteInformation() SiteInformation {
	data, err := ioutil.ReadFile(ConfigFile)
	if err != nil {
		log.Fatal("error accessing directory %s: %v", ConfigFile, err)
	}
	siteInfo := SiteInformation{}
	siteInfo.ParseJSON(data)
	return fillDefaultValues(siteInfo)
}

func (c *SiteInformation) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}

func fillDefaultValues(siteInfo SiteInformation) SiteInformation {
	if siteInfo.TempFolder == "" {
		siteInfo.TempFolder = "./tmp"
	}
	if siteInfo.DestFolder == "" {
		siteInfo.DestFolder = "./public"
	}
	if siteInfo.NumPostsFrontPage == 0 {
		siteInfo.NumPostsFrontPage = 10
	}
	return siteInfo
}
