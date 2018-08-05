// Package config fills SiteInformation struct that contains
// all the necessary configuration for creating the blog.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type BlogInformation interface {
}

// SiteInformation contains the information inside ConfigFile
type SiteInformation struct {
	Author            string `json:"Author"`
	BlogURL           string `json:"BlogURL"`
	BlogLanguage      string `json:"BlogLanguage"`
	BlogDescription   string `json:"BlogDescription"`
	DateFormat        string `json:"DateFormat"`
	Theme             Theme
	ThemeFolder       string `json:"ThemeFolder"`
	BlogTitle         string `json:"BlogTitle"`
	NumPostsFrontPage int    `json:"NumPostsFrontPage"`
	DataSource        DataSource
	Upload            Upload
	TempFolder        string       `json:"TempFolder"`
	DestFolder        string       `json:"DestFolder"`
	StaticPages       []StaticPage `json:"StaticPages"`
}

type Theme struct {
	Repository string `json:"Repository"`
	Type       string `json:"Type"`
}

type StaticPage struct {
	File       string `json:"File"`
	To         string `json:"To"`
	IsTemplate bool   `json:"IsTemplate"`
}

type DataSource struct {
	Type       string `json:"Type"`
	Repository string `json:"Repository"`
}

type Upload struct {
	Type     string `json:"Type"`
	URL      string `json:"URL"`
	Username string `json:"Username"`
	Password string `json:"Password"`
}

func New(configFile string) (SiteInformation, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return SiteInformation{}, fmt.Errorf("error accessing directory %s: %v", configFile, err)
	}
	siteInfo := new(SiteInformation)
	siteInfo.parseJSON(data)
	siteInfo.fillDefaultValues()
	return *siteInfo, nil
}

func (si *SiteInformation) parseJSON(b []byte) (err error) {
	return json.Unmarshal(b, &si)
}

func (si *SiteInformation) fillDefaultValues() {
	if si.TempFolder == "" {
		si.TempFolder = "./tmp"
	}
	if si.DestFolder == "" {
		si.DestFolder = "./public"
	}
	if si.ThemeFolder == "" {
		si.DestFolder = "./static"
	}
	if si.NumPostsFrontPage == 0 {
		si.NumPostsFrontPage = 10
	}
}
