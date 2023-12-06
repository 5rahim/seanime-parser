package seanime_parser

type keywordCategory uint8

const (
	keywordCategoryUnknown keywordCategory = iota
	keywordCategorySeasonPrefix
	keywordCategoryAnimeType
	keywordCategoryYear
	keywordCategoryAudioTerm
	keywordCategoryDeviceCompat
	keywordCategoryEpisodePrefix
	keywordCategoryPartPrefix
	keywordCategoryVolumePrefix
	keywordCategoryFileChecksum
	keywordCategoryFileExtension
	keywordCategoryLanguage
	keywordCategoryReleaseGroup
	keywordCategoryReleaseInformation
	keywordCategoryReleaseVersion
	keywordCategorySource
	keywordCategorySubtitles
	keywordCategoryVideoResolution
	keywordCategoryVideoTerm
)

type keywordKind uint8

const (
	keywordKindCombinedWithNumber keywordKind = iota
	keywordKindSeparatedWithNumber
	keywordKindOrdinalSuffix
	keywordKindStandalone
)

type (
	keyword struct {
		value    string
		category keywordCategory
		kind     keywordKind
	}

	keywordManager struct {
		keywords []*keyword
	}
)

func newKeywordManager() *keywordManager {
	km := keywordManager{
		keywords: make([]*keyword, 0),
	}

	// Season

	km.addGroup(
		keywordCategorySeasonPrefix,
		keywordKindCombinedWithNumber,
		[]string{"S"},
	)
	km.addGroup(
		keywordCategorySeasonPrefix,
		keywordKindSeparatedWithNumber,
		[]string{"SEASON", "SAISON", "SEASONS", "SAISONS"},
	)
	km.addGroup(
		keywordCategorySeasonPrefix,
		keywordKindOrdinalSuffix,
		[]string{"SEASON", "SAISON", "SEASONS", "SAISONS"},
	)

	// Episode

	km.addGroup(
		keywordCategoryEpisodePrefix,
		keywordKindCombinedWithNumber,
		[]string{"E", "\x7B2C", "EP", "EPS", "EPISODE", "EPISODES", "CAPITULO", "EPISODIO", "FOLDE"},
	)
	km.addGroup(
		keywordCategoryEpisodePrefix,
		keywordKindSeparatedWithNumber,
		[]string{"EP", "EPS", "EPISODE", "EPISODES", "CAPITULO", "EPISODIO", "FOLDE"},
	)

	// Volume

	km.addGroup(
		keywordCategoryVolumePrefix,
		keywordKindCombinedWithNumber,
		[]string{"VOL", "VOLUME", "VOLUMES"},
	)
	km.addGroup(
		keywordCategoryVolumePrefix,
		keywordKindSeparatedWithNumber,
		[]string{"VOL", "VOLUME", "VOLUMES"},
	)

	// Part

	km.addGroup(
		keywordCategoryPartPrefix,
		keywordKindCombinedWithNumber,
		[]string{"PART", "PARTS", "COUR"},
	)
	km.addGroup(
		keywordCategoryPartPrefix,
		keywordKindSeparatedWithNumber,
		[]string{"PART", "PARTS", "COUR"},
	)
	km.addGroup(
		keywordCategoryPartPrefix,
		keywordKindOrdinalSuffix,
		[]string{"PART", "COUR"},
	)

	// Anime Type

	km.addGroup(
		keywordCategoryAnimeType,
		keywordKindCombinedWithNumber,
		[]string{"SP"},
	)
	km.addGroup(
		keywordCategoryAnimeType,
		keywordKindCombinedWithNumber,
		[]string{"SP", "MOVIE", "OAD", "OAV", "ONA", "OVA", "SPECIAL", "SPECIALS", "ED", "ENDING", "NCED", "NCOP", "OPED", "OP", "OPENING",
			"番外編", "總集編", "映像特典", "特典", "特典アニメ"},
	)
	km.addGroup(
		keywordCategoryAnimeType,
		keywordKindSeparatedWithNumber,
		[]string{"SP", "MOVIE", "OAD", "OAV", "ONA", "OVA", "SPECIAL", "SPECIALS", "ED", "ENDING", "NCED", "NCOP", "OPED", "OP", "OPENING",
			"番外編", "總集編", "映像特典", "特典", "特典アニメ"},
	)
	km.addGroup(
		keywordCategoryAnimeType,
		keywordKindStandalone,
		[]string{
			"MOVIE", "GEKIJOUBAN", "ONA", "OVA", "OAV", "OAD", "SPECIALS", "TV",
			"ED", "ENDING", "NCED", "NCOP", "OPED", "OP", "OPENING", "PREVIEW", "PV", "EVENT", "TOKUTEN", "LOGO", "CM", "SPOT", "MENU"},
	)

	// Audio Term

	km.addGroup(
		keywordCategoryAudioTerm,
		keywordKindStandalone,
		[]string{
			// Audio channels
			"2.0CH", "2CH", "5.1", "5.1CH", "DTS", "DTS-ES", "DTS5.1", "TRUEHD5.1",
			// Audio codec
			"AAC", "AACX2", "AACX3", "AACX4", "AC3", "EAC3", "E-AC-3", "FLAC",
			"FLACX2", "FLACX3", "FLACX4", "LOSSLESS", "MP3", "OGG", "VORBIS",
			"DD2", "DD2.0",
			// Audio language
			"DUALAUDIO", "DUAL-AUDIO",
		},
	)

	// Video Term

	km.addGroup(
		keywordCategoryVideoTerm,
		keywordKindStandalone,
		[]string{
			// Frame rate
			"23.976FPS", "24FPS", "29.97FPS", "30FPS", "60FPS", "120FPS",
			// Video codec
			"8BIT", "8-BIT", "10BIT", "10BITS", "10-BIT", "10-BITS",
			"HI10", "HI10P", "HI444", "HI444P", "HI444PP",
			"H264", "H265", "H.264", "H.265", "X264", "X265", "X.264",
			"AVC", "HEVC", "HEVC2", "DIVX", "DIVX5", "DIVX6", "XVID",
			"AV1",
			"HDR", "DV", "DOLBY VISION",
			// Video format
			"AVI", "RMVB", "WMV", "WMV3", "WMV9",
			// Video quality
			"HQ", "LQ",
			// Video resolution
			"HD", "SD", "4K",
		},
	)

	// Device Compat

	km.addGroup(
		keywordCategoryDeviceCompat,
		keywordKindStandalone,
		[]string{"IPAD3", "IPHONE5", "IPOD", "PS3", "XBOX", "XBOX360", "ANDROID"},
	)

	// File Extension
	// should be last

	km.addGroup(
		keywordCategoryFileExtension,
		keywordKindStandalone,
		[]string{"3GP", "AVI", "DIVX", "FLV", "M2TS", "MKV", "MOV", "MP4", "MPG",
			"OGM", "RM", "RMVB", "TS", "WEBM", "WMV"},
	)

	// Language
	// should be enclosed

	km.addGroup(
		keywordCategoryLanguage,
		keywordKindStandalone,
		[]string{"ENG", "ENGLISH", "ESPANOL", "JAP", "PT-BR", "SPANISH", "VOSTFR", "ESP", "ITA"},
	)

	// Release info

	km.addGroup(
		keywordCategoryReleaseInformation,
		keywordKindStandalone,
		[]string{"REMASTER", "REMASTERED", "UNCENSORED", "UNCUT", "TS", "VFR",
			"WIDESCREEN", "WS", "BATCH", "COMPLETE", "PATCH", "REMUX"},
	)

	km.addGroup(
		keywordCategorySubtitles,
		keywordKindStandalone,
		[]string{"ASS", "BIG5", "DUB", "DUBBED", "HARDSUB", "HARDSUBS", "RAW",
			"SOFTSUB", "SOFTSUBS", "SUB", "SUBBED", "SUBTITLED", "MULTISUB"},
	)

	// Source

	km.addGroup(
		keywordCategorySubtitles,
		keywordKindStandalone,
		[]string{"BD", "BDRIP", "BLURAY", "BLU-RAY", "DVD", "DVD5", "DVD9",
			"DVD-R2J", "DVDRIP", "DVD-RIP", "R2DVD", "R2J", "R2JDVD",
			"R2JDVDRIP", "HDTV", "HDTVRIP", "TVRIP", "TV-RIP",
			"WEBCAST", "WEBRIP"},
	)

	return &km
}

func (km *keywordManager) addGroup(category keywordCategory, kind keywordKind, group []string) {
	for _, value := range group {
		km.keywords = append(km.keywords, &keyword{
			value:    value,
			category: category,
			kind:     kind,
		})
	}
}
