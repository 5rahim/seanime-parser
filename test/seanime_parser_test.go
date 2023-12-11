package test

import (
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"seanime-parser"
	"testing"
)

func TestSeanimeParser(t *testing.T) {

	data := getData()
	assert.NotNil(t, data)

	for _, tt := range data {
		t.Run(tt.FileName, func(t *testing.T) {

			metadata := seanime_parser.Parse(tt.FileName)
			assert.NotNil(t, metadata)

			assertMetadataEquals(t, metadata.Season, tt.Season, "Season")
			assertMetadataEquals(t, metadata.EpisodeNumber, tt.EpisodeNumber, "Episode")
			assertMetadataEquals(t, metadata.Part, tt.Part, "Part")
			assertMetadataEquals(t, metadata.Title, tt.Title, "Title")
			assertMetadataEquals(t, metadata.AnimeType, tt.AnimeType, "AnimeType")
			assertMetadataEquals(t, metadata.Year, tt.Year, "Year")
			assertMetadataEquals(t, metadata.AudioTerm, tt.AudioTerm, "AudioTerm")
			assertMetadataEquals(t, metadata.DeviceCompatibility, tt.DeviceCompatibility, "DeviceCompatibility")
			assertMetadataEquals(t, metadata.EpisodeNumberAlt, tt.EpisodeNumberAlt, "EpisodeNumberAlt")
			assertMetadataEquals(t, metadata.EpisodeTitle, tt.EpisodeTitle, "EpisodeTitle")
			assertMetadataEquals(t, metadata.FileChecksum, tt.FileChecksum, "FileChecksum")
			assertMetadataEquals(t, metadata.FileExtension, tt.FileExtension, "FileExtension")
			assertMetadataEquals(t, metadata.FileName, tt.FileName, "FileName")
			assertMetadataEquals(t, metadata.Language, tt.Language, "Language")
			assertMetadataEquals(t, metadata.ReleaseGroup, tt.ReleaseGroup, "ReleaseGroup")
			assertMetadataEquals(t, metadata.ReleaseInformation, tt.ReleaseInformation, "ReleaseInformation")
			assertMetadataEquals(t, metadata.ReleaseVersion, tt.ReleaseVersion, "ReleaseVersion")
			assertMetadataEquals(t, metadata.Source, tt.Source, "Source")
			assertMetadataEquals(t, metadata.Subtitles, tt.Subtitles, "Subtitles")
			assertMetadataEquals(t, metadata.VideoResolution, tt.VideoResolution, "VideoResolution")
			assertMetadataEquals(t, metadata.VideoTerm, tt.VideoTerm, "VideoTerm")
			assertMetadataEquals(t, metadata.VolumeNumber, tt.VolumeNumber, "VolumeNumber")

		})
	}

}

func assertMetadataEquals(t *testing.T, received interface{}, expected interface{}, kind string) {
	if expected == nil {
		if received == nil {
			return
		} else {
			assert.Failf(t, "Expected %s to be nil but got %s", kind, received)
		}
	}

	assert.Equalf(t, expected, received, "Expected %s to be %s but got %s", kind, expected, received)
}

func getData() []*seanime_parser.Metadata {

	file, err := os.Open("data.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var metadata []*seanime_parser.Metadata
	if err := decoder.Decode(&metadata); err != nil {
		log.Fatalf("Error decoding JSON: %s", err)
	}

	return metadata
}
