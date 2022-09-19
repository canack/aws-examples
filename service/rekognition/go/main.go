package main

import (
	"bytes"
	"context"
	"errors"
	"github.com/aws/aws-sdk-go-v2/config"                    // AWS SDK
	"github.com/aws/aws-sdk-go-v2/service/rekognition"       // AWS Rekognition
	"github.com/aws/aws-sdk-go-v2/service/rekognition/types" // AWS Rekognition types
	"github.com/fogleman/gg"                                 // To draw rectangle, used another library currently
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
)

// Errors
var (
	ErrorNoFaceFound = errors.New("recognized 0 faces")
	ErrorEncode      = errors.New("encode unsuccessful")
	ErrorDecode      = errors.New("decode unsuccessful")
	ErrorServer      = errors.New("connection error")
)

var Client *rekognition.Client // to access aws rekognition service globally

// SetupRekognition To initialize aws rekognition service
func SetupRekognition() error {
	// Load default config
	cfg, err := config.LoadDefaultConfig(context.TODO())

	// If a server error occurs, return error.
	if err != nil {
		return err
	}

	// Create rekognition client
	Client = rekognition.NewFromConfig(cfg)
	return nil

}

func ProcessImage(input io.Reader) (io.ReadCloser, error) {

	// Decode image
	img, _, err := image.Decode(input)
	if err != nil {
		return io.ReadCloser(nil), ErrorDecode
	}

	// To save images bounds.
	x := img.Bounds().Dx()
	y := img.Bounds().Dy()

	// Image encoding before sent to aws rekognition service
	imgBuffer := new(bytes.Buffer)

	// I used jpeg format because it is smaller than png.
	if err := jpeg.Encode(imgBuffer, img, &jpeg.Options{Quality: 80}); err != nil {
		return io.ReadCloser(nil), ErrorEncode
	}

	// Rekognition service request
	out, err := Client.DetectFaces(context.TODO(), &rekognition.DetectFacesInput{
		Image:      &types.Image{Bytes: imgBuffer.Bytes()},    // Image bytes
		Attributes: []types.Attribute{types.AttributeDefault}, // Attributes to return
	})

	// If a server error occurs, return error.
	if err != nil {
		return io.ReadCloser(nil), ErrorServer
	}

	// If no face is found, return error.
	if len(out.FaceDetails) == 0 {
		return io.ReadCloser(nil), ErrorNoFaceFound
	}

	// To draw rectangle and write age as text, used another library currently
	drw := gg.NewContextForImage(img)

	// To processing all faces in the image
	for _, v := range out.FaceDetails {
		bb := v.BoundingBox                                       // Bounding box of face
		x0 := float64(*bb.Left * float32(x))                      // Left
		y0 := float64(*bb.Top * float32(y))                       // Top
		x1 := float64(*bb.Width*float32(x) + *bb.Left*float32(x)) // Right
		y1 := float64(*bb.Height*float32(y) + *bb.Top*float32(y)) // Bottom

		dragRectangle(drw, x0, y0, x1, y1) // Draw rectangle around face
	}

	// Encoding
	buf := new(bytes.Buffer)

	// I used jpeg format because it is smaller than png.
	if err := jpeg.Encode(buf, drw.Image(), &jpeg.Options{Quality: 90}); err != nil {
		return io.ReadCloser(nil), ErrorDecode
	}

	// Return image
	return io.NopCloser(buf), nil
}

// To draw rectangle
func dragRectangle(drw *gg.Context, x0, y0, x1, y1 float64) {
	drw.DrawRectangle(x0, y0, x1-x0, y1-y0)
	drw.SetLineWidth(2)
	drw.SetHexColor("#f54281")
	drw.StrokePreserve()
	drw.SetRGBA(0, 0, 0, 0.5)
	drw.Fill()
}

func main() {
	// Setup aws rekognition service
	if err := SetupRekognition(); err != nil {
		panic(err)
	}

	input := os.Args[1]  // Input image path
	output := os.Args[2] // Output image path

	// open image file to process
	file, err := os.Open(input)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Recognize faces
	result, err := ProcessImage(file)
	if err != nil {
		panic(err)
	}

	// save result to file
	out, err := os.Create(output)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	// Copy result to file
	_, err = io.Copy(out, result)
	if err != nil {
		panic(err)
	}
	
}
