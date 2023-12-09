package seanime_parser

type metadataCategory = uint8

const (
	metadataUnknown metadataCategory = iota
	metadataSeason
	metadataEpisodeNumber
	metadataPart
	metadataTitle
	metadataAnimeType
	metadataYear
	metadataAudioTerm
	metadataDeviceCompat
	metadataEpisodeNumberAlt
	metadataEpisodeTitle
	metadataChecksum
	metadataLanguage
	metadataSubtitles
	metadataReleaseVersion
	metadataSource
	metadataVideoResolution
	metadataVideoTerm
	metadataVolumeNumber
	metadataOtherEpisodeNumber
)

type Metadata struct {
	// Slice of strings representing the season of anime. "S1-S3" would be represented as []string{"1", "3"}.
	Season []string `json:"season,omitempty"`
	Part   []string `json:"anime_part,omitempty"`

	// Represents the strings prefixing the season in the file, e.g in "SEASON 2" "SEASON" is the SeasonPrefix.
	SeasonPrefix []string `json:"season_prefix,omitempty"`
	// Represents the strings prefixing the season in the file, e.g in "PART 2" "PART" is the SeasonPrefix.
	PartPrefix []string `json:"part_prefix,omitempty"`

	// Title of the Anime. e.g in "[HorribleSubs] Boku no Hero Academia - 01 [1080p].mkv",
	// "Boku No Hero Academia" is the AnimeTitle.
	Title string `json:"title,omitempty"`

	// Slice of strings representing the types specified in the anime file, e.g ED, OP, Movie, etc.
	AnimeType []string `json:"anime_type,omitempty"`

	// Year the anime was released.
	Year string `json:"year,omitempty"`

	// Slice of strings representing the audio terms included in the filename, e.g FLAC, AAC, etc.
	AudioTerm []string `json:"audio_term,omitempty"`

	// Slice of strings representing devices the video is compatible with that are mentioned in the filename.
	DeviceCompatibility []string `json:"device_compatibility,omitempty"`

	// Slice of strings representing the episode numbers. "01-10" would be respresented as []string{"1", "10"}.
	EpisodeNumber []string `json:"episode_number,omitempty"`

	// Slice of strings representing the alternative episode number.
	// This is for cases where you may have an episode number relative to the season,
	// but a larger episode number as if it were all one season.
	// e.g in [Hatsuyuki]_Kuroko_no_Basuke_S3_-_01_(51)_[720p][10bit][619C57A0].mkv
	// 01 would be the EpisodeNumber, and 51 would be the EpisodeNumberAlt.
	EpisodeNumberAlt []string `json:"episode_number_alt,omitempty"`

	// Slice of strings representing the words prefixing the episode number in the file, e.g in "EPISODE 2", "EPISODE" is the prefix.
	EpisodePrefix []string `json:"episode_prefix,omitempty"`

	// Title of the episode. e.g in "[BM&T] Toradora! - 07v2 - Pool Opening [720p Hi10 ] [BD] [8F59F2BA]",
	// "Pool Opening" is the EpisodeTitle.
	EpisodeTitle string `json:"episode_title,omitempty"`

	// Checksum of the file, in [BM&T] Toradora! - 07v2 - Pool Opening [720p Hi10 ] [BD] [8F59F2BA],
	// "8F59F2BA" would be the FileChecksum.
	FileChecksum string `json:"file_checksum,omitempty"`

	// File extension, in [HorribleSubs] Boku no Hero Academia - 01 [1080p].mkv,
	// "mkv" would be the FileExtension.
	FileExtension string `json:"file_extension,omitempty"`

	// Full filename that was parsed.
	FileName string `json:"file_name,omitempty"`

	// Languages specified in the file name, e.g RU, JP, EN etc.
	Language []string `json:"language,omitempty"`

	// Terms that could not be parsed into other buckets, but were deemed identifiers.
	// In [chibi-Doki] Seikon no Qwaser - 13v0 (Uncensored Director's Cut) [988DB090].mkv,
	// "Uncensored" is parsed into Other.
	Other []string `json:"other,omitempty"`

	// The fan sub group that uploaded the file. In [HorribleSubs] Boku no Hero Academia - 01 [1080p],
	// "HorribleSubs" is the ReleaseGroup.
	ReleaseGroup string `json:"release_group,omitempty"`

	// Information about the release that wasn't a version.
	// In "[SubDESU-H] Swing out Sisters Complete Version (720p x264 8bit AC3) [3ABD57E6].mp4
	// "Complete" is parsed into ReleaseInformation.
	ReleaseInformation []string `json:"release_information,omitempty"`

	// Slice of strings representing the version of the release.
	// In [FBI] Baby Princess 3D Paradise Love 01v0 [BD][720p-AAC][457CC066].mkv, 0 is parsed into ReleaseVersion.
	ReleaseVersion []string `json:"release_version,omitempty"`

	// Slice of strings representing where the video was ripped from. e.g BLU-RAY, DVD, etc.
	Source []string `json:"source,omitempty"`

	// Slice of strings representing the type of subtitles included, e.g HARDSUB, BIG5, etc.
	Subtitles []string `json:"subtitles,omitempty"`

	// Resolution of the video. Can be formatted like 1920x1080, 1080, 1080p, etc depending
	// on how it is represented in the filename.
	VideoResolution string `json:"video_resolution,omitempty"`

	// Slice of strings representing the video terms included in the filename, e.g h264, x264, etc.
	VideoTerm []string `json:"video_term,omitempty"`

	// Slice of strings represnting the volume numbers. "01-10" would be represented as []string{"1", "10"}.
	VolumeNumber []string `json:"volume_number,omitempty"`

	// Slice of strings representing the words prefixing the volume number in the file, e.g in "VOLUME 2", "VOLUME" is the prefix.
	VolumePrefix []string `json:"volume_prefix,omitempty"`

	// Entries that could not be parsed into any other categories.
	Unknown []string `json:"unknown,omitempty"`

	// Bool determining if "EpisodeNumberAlt" should be parsed or not.
	checkAltNumber bool
}
