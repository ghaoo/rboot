package richman

import (
	"bytes"
	"errors"
	"image/jpeg"
	"image/png"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// ImageOptions represent the options to generate the image.
type ImageOptions struct {
	// BinaryPath the path to your wkhtmltoimage binary. REQUIRED
	//
	// Must be absolute path e.g /usr/local/bin/wkhtmltoimage
	BinaryPath string
	// Input is the content to turn into an image. REQUIRED
	//
	// Can be a url (http://example.com), a local file (/tmp/example.html), or html as a string (send "-" and set the Html value)
	Input string
	// Format is the type of image to generate
	//
	// jpg, png, svg, bmp supported. Defaults to local wkhtmltoimage default
	Format string
	// Height is the height of the screen used to render in pixels.
	//
	// Default is calculated from page content. Default 0 (renders entire page top to bottom)
	Height int
	// Width is the width of the screen used to render in pixels.
	//
	// Note that this is used only as a guide line. Default 1024
	Width int
	// Quality determines the final image quality.
	//
	// Values supported between 1 and 100. Default is 94
	Quality int
	// Crop-x determines the final image crop from x
	CropX int
	// Crop-y determines the final image crop from y
	CropY int
	// Crop-w determines the final image crop width
	CropW int
	// Crop-h determines the final image crop height
	CropH int
	// Html is a string of html to render into and image.
	//
	// Only needed to be set if Input is set to "-"
	Html string
	// Output controls how to save or return the image.
	//
	// Leave nil to return a []byte of the image. Set to a path (/tmp/example.png) to save as a file.
	Output string
}

// GenerateImage creates an image from an input.
// It returns the image ([]byte) and any error encountered.
func GenerateImage(options *ImageOptions) ([]byte, error) {
	arr, err := buildParams(options)
	if err != nil {
		return []byte{}, err
	}

	if options.BinaryPath == "" {
		return []byte{}, errors.New("BinaryPath not set")
	}

	cmd := exec.Command(options.BinaryPath, arr...)

	if options.Html != "" {
		cmd.Stdin = strings.NewReader(options.Html)
	}

	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Println("can't generate file ", err)
		return nil, err
	}
	if options.Output == "" && len(output) > 0 {
		trimmed := cleanupOutput(output, options.Format)
		return trimmed, err
	} else {
		return output, err
	}
}

// buildParams takes the image options set by the user and turns them into command flags for wkhtmltoimage
// It returns an array of command flags.
func buildParams(options *ImageOptions) ([]string, error) {
	a := []string{}

	if options.Input == "" {
		return []string{}, errors.New("Must provide input")
	}

	// silence extra wkhtmltoimage output
	// might want to add --javascript-delay too?
	a = append(a, "-q")
	a = append(a, "--disable-plugins")

	a = append(a, "--format")
	if options.Format != "" {
		a = append(a, options.Format)
	} else {
		a = append(a, "png")
	}

	if options.Height != 0 {
		a = append(a, "--height")
		a = append(a, strconv.Itoa(options.Height))
	}

	if options.Width != 0 {
		a = append(a, "--width")
		a = append(a, strconv.Itoa(options.Width))
	}

	if options.Quality != 0 {
		a = append(a, "--quality")
		a = append(a, strconv.Itoa(options.Quality))
	}

	if options.CropX != 0 {
		a = append(a, "--crop-x")
		a = append(a, strconv.Itoa(options.CropX))
	}

	if options.CropY != 0 {
		a = append(a, "--crop-y")
		a = append(a, strconv.Itoa(options.CropX))
	}

	if options.CropW != 0 {
		a = append(a, "--crop-w")
		a = append(a, strconv.Itoa(options.CropX))
	}

	if options.CropH != 0 {
		a = append(a, "--crop-h")
		a = append(a, strconv.Itoa(options.CropX))
	}

	// url and output come last
	if options.Input != "-" {
		// make sure we dont pass stdin if we aren't expecting it
		options.Html = ""
	}

	a = append(a, options.Input)

	if options.Output == "" {
		a = append(a, "-")
	} else {
		a = append(a, options.Output)
	}

	return a, nil
}

func cleanupOutput(img []byte, format string) []byte {
	buf := new(bytes.Buffer)
	switch {
	case format == "png":
		decoded, err := png.Decode(bytes.NewReader(img))
		for err != nil {
			img = img[1:]
			if len(img) == 0 {
				break
			}
			decoded, err = png.Decode(bytes.NewReader(img))
		}
		png.Encode(buf, decoded)
		return buf.Bytes()
	case format == "jpg":
		decoded, err := jpeg.Decode(bytes.NewReader(img))
		for err != nil {
			img = img[1:]
			if len(img) == 0 {
				break
			}
			decoded, err = jpeg.Decode(bytes.NewReader(img))
		}
		jpeg.Encode(buf, decoded, nil)
		return buf.Bytes()
		// case format == "svg":
		// 	return img
	}
	return img
}
