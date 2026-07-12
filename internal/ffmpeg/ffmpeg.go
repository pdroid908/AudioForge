package ffmpeg

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"strings"

	"price-comparator/internal/models"
)

func Exec(inputPath, outputPath string, req models.ExportRequest, ffmpegPath string) error {
	args := buildArgs(inputPath, outputPath, req)
	cmd := exec.Command(ffmpegPath, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildArgs(inputPath, outputPath string, req models.ExportRequest) []string {
	args := []string{"-y", "-i", inputPath}
	filters := []string{}

	if req.Speed != 1 {
		speed := req.Speed
		if speed < 0.5 {
			speed = 0.5
		}
		if speed > 2 {
			speed = 2
		}
		filters = append(filters, fmt.Sprintf("atempo=%.2f", speed))
	}

	if req.Pitch != 0 {
		ratio := math.Pow(2, float64(req.Pitch)/12)
		filters = append(filters, fmt.Sprintf("asetrate=44100*%.4f", ratio))
	}

	if req.Volume != 100 {
		factor := req.Volume / 100.0
		filters = append(filters, fmt.Sprintf("volume=%.2f", factor))
	}

	if req.Bass != 0 {
		filters = append(filters, fmt.Sprintf("bass=%.2f", req.Bass))
	}

	if req.Treble != 0 {
		filters = append(filters, fmt.Sprintf("treble=%.2f", req.Treble))
	}

	if req.Echo {
		filters = append(filters, "aecho=0.8:0.9:1000:0.3")
	}

	if req.Reverb {
		filters = append(filters, "aecho=0.8:0.9:3000:0.2")
	}

	if req.FadeIn {
		filters = append(filters, "afade=t=in:ss=0:d=1")
	}

	if req.FadeOut {
		filters = append(filters, "afade=t=out:st=-1:d=1")
	}

	if req.Normalize {
		filters = append(filters, "dynaudnorm")
	}

	if len(filters) > 0 {
		args = append(args, "-af", strings.Join(filters, ","))
	}
	args = append(args, "-vn", "-codec:a", "libmp3lame", outputPath)
	return args
}
