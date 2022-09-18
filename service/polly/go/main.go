package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	"github.com/aws/aws-sdk-go-v2/service/polly/types"
	"io"
	"os"
)

func main() {
	// Creating a context to handle cancellation if needed
	ctx := context.Background()

	// Load the SDK's configuration from environment and shared config
	cfg, err := config.LoadDefaultConfig(ctx)

	if err != nil {
		panic("configuration error, " + err.Error())
	}

	// Create the client with the configuration loaded
	client := polly.NewFromConfig(cfg)

	// Make the TextToSpeech request
	out, err := client.SynthesizeSpeech(ctx, &polly.SynthesizeSpeechInput{
		Text:         argsToString(),         // Pass the text to synthesize from the command line arguments
		OutputFormat: types.OutputFormatMp3,  // Set the output format to MP3
		VoiceId:      types.VoiceIdJoanna,    // Set the voice to Joanna. You can also use other voices
		LanguageCode: types.LanguageCodeEnUs, // Set the language to US English
		Engine:       types.EngineNeural,     // Set the engine to neural (neural voices are only available in some regions)
		TextType:     types.TextTypeText,     // Set the text type to text (text or ssml)
	})

	if err != nil {
		panic("failed to synthesize speech, " + err.Error())
	}

	// Save the audio stream to a file
	err = saveToFile(out.AudioStream, "speech.mp3")
	if err != nil {
		panic("failed to save to file, " + err.Error())
	}

}

// Get the text to synthesize from the command line arguments
func argsToString() *string {
	var s string
	for _, arg := range os.Args[1:] {
		s += arg + " "
	}
	return &s
}

// Save the audio stream to a file
func saveToFile(r io.Reader, filename string) error {
	f, err := os.Create(filename) // Create the file
	if err != nil {
		return err
	}
	defer f.Close() // Close the file when we're done

	_, err = io.Copy(f, r) // Copy the audio stream to the file
	if err != nil {
		return err
	}

	return nil
}
