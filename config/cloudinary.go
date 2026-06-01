package config

import (
	"github.com/NhatHaoDev3324/zizone-be/pkg/log"
	"github.com/cloudinary/cloudinary-go/v2"
)

var Cloud *cloudinary.Cloudinary

func InitCloudinary() {

	cld, err := cloudinary.New()
	if err != nil {
		log.LogError("Failed to connect to Cloudinary: " + err.Error())
	}

	Cloud = cld
	log.LogSuccess("Connected to Cloudinary successfully!")
}
