package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"            // required to load config
	"github.com/aws/aws-sdk-go-v2/service/translate" // required to use the translation service
	"os"
)

func getCommandLineArgs() (string, string, string) {
	// Get the command line arguments
	args := os.Args[1:]
	if len(args) != 3 {
		fmt.Println("You must supply 3 command line arguments: SOURCE_LANG TARGET_LANG TEXT")
		os.Exit(1)
	}
	sourceLang := args[0]
	targetLang := args[1]
	text := args[2]
	return sourceLang, targetLang, text
}

func main() {
	// load the SDK configuration from environment and shared config
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// Get the command line arguments
	sourceLanguage, targetLanguage, text := getCommandLineArgs()

	// Create the service client with the config
	client := translate.NewFromConfig(cfg)

	// Make the Translate API call
	request, err := client.TranslateText(context.TODO(), &translate.TranslateTextInput{
		SourceLanguageCode: &sourceLanguage, // Required
		TargetLanguageCode: &targetLanguage, // Required
		Text:               &text,           // Required
	})

	if err != nil {
		fmt.Println("Got an error calling TranslateText: " + err.Error())
		os.Exit(1)
	}

	// Display the results
	result := request.TranslatedText
	fmt.Println("Original text: " + text)
	fmt.Println("Translated text: " + *result)
}
