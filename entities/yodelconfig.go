package entities

import (
	"fmt"

	"github.com/spf13/viper"
)

type HeaderLable = string

// RuteCsvConfig: Rute untuk konfigurasi berkas CSV
type RuteCsvConfig struct {
	Alamat       *string
	Blacklist    []string
	PrimaryIndex HeaderLable
}

// YodelConfig: YOur DELivery CONFIGuration
type YodelConfig struct {
	NomorPort     uint16
	FolderIndukan string
	RuteCsv       map[string]RuteCsvConfig
}

// DefaultFolderIndukan :nodoc:
const DefaultFolderIndukan = "./csvs"

// DefaultNomorPort :nodoc:
const DefaultNomorPort uint16 = 8125

// NewYodelConfig: membuat instance baru dari struct YodelConfig
func NewYodelConfig() *YodelConfig {
	// setup lokasi dan type config
	viper.SetConfigType("hcl")
	viper.SetConfigFile(".yodelconf.tf")

	// baca koding dari file config
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}

	yodel := new(YodelConfig)

	yodel.FolderIndukan = viper.GetString("folder_indukan")
	if yodel.FolderIndukan == "" {
		yodel.FolderIndukan = DefaultFolderIndukan
	}

	yodel.NomorPort = uint16(viper.GetUint("nomor_port"))
	if yodel.NomorPort == 0 {
		yodel.NomorPort = DefaultNomorPort
	}

	return yodel
}
