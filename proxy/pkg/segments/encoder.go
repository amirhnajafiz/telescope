package segments

import (
	"fmt"
	"os"
	"os/exec"
)

// Encode encodes a video file into segments using MP4Box.
func Encode(input, output string) error {
	cmd := exec.Command(
		"MP4Box",
		"-dash", "4000", // 4s segments
		"-frag", "4000",
		"-rap",
		"-profile", "dashavc264:onDemand",
		"-out", fmt.Sprintf("%s/manifest.mpd", output),
		input,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
