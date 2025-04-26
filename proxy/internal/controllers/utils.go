package controllers

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

// shouldSelectReplica determines whether the next replica is selected based on a probability
func shouldSelectReplica(probability float64) bool {
	return probability > 0 && probability >= (float64(rand.Intn(100))/100.0)
}

// constructFullPath replaces placeholders in the Media template with actual values
func constructFullPath(template, representationID string, number int) string {
	// replace $RepresentationID$ with the actual representation ID
	path := strings.ReplaceAll(template, "$RepresentationID$", representationID)

	// replace $Number%05d$ with the formatted number
	numberPlaceholder := regexp.MustCompile(`\$Number%0(\d+)d\$`)
	path = numberPlaceholder.ReplaceAllStringFunc(path, func(match string) string {
		width, _ := strconv.Atoi(match[8 : len(match)-2]) // Extract width from %05d
		return fmt.Sprintf("%0*d", width, number)
	})

	// remove "/stream" from the path
	path = strings.ReplaceAll(path, "/stream", "")

	// construct the relative path by trimming the "/api/" prefix
	relativePath := strings.TrimPrefix(path, "/api/")

	return relativePath
}
