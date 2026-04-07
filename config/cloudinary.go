package config

import (
	"github.com/NhatHaoDev3324/zizone-be/factory"
	"github.com/cloudinary/cloudinary-go/v2"
)

var Cloud *cloudinary.Cloudinary

func InitCloudinary() {

	cld, err := cloudinary.New()
	if err != nil {
		factory.LogError("Failed to connect to Cloudinary: " + err.Error())
	}

	Cloud = cld
	factory.LogSuccess("Connected to Cloudinary successfully!")
}
