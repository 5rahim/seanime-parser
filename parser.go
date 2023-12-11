package seanime_parser

import (
	"strings"
)

type parser struct {
	filename     string
	tokenManager *tokenManager
	metadata     *Metadata
}

func newParser(filename string) *parser {
	return &parser{
		filename:     filename,
		tokenManager: newTokenManager(filename),
		metadata:     &Metadata{},
	}
}

func Parse(filename string) *Metadata {

	p := newParser(filename)
	p.parse()
	p.cleanUp()

	return p.metadata

}

func (p *parser) parse() {

	p.parseKeywords("normal")

	p.parseSeason()

	p.parseEpisode()

	p.parseEpisodeTitle()

	p.parseTitle()

	p.parseReleaseGroup()

	p.parseKeywords("")

	p.collectMetadata()

}

func (p *parser) parseKeywords(priority string) {

	for _, tkn := range *p.tokenManager.tokens {

		if tkn.isKeyword() || !tkn.isUnknown() { // Don't bother if token is already a keyword
			continue // Skip to next token
		}

		// Identify keyword
		_ = p.identifyKeyword(tkn, priority)

	}

}

// identifyKeyword identifies STANDALONE and multi-PART keywords for the given token
func (p *parser) identifyKeyword(tkn *token, priority string) bool {

	if tkn.Kind == tokenKindCrc32 {
		tkn.setIdentifiedKeywordCategory(keywordCatFileChecksum)
		return true
	}

	if tkn.Kind == tokenKindPossibleVideoRes {
		tkn.setIdentifiedKeywordCategory(keywordCatVideoResolution)
		return true
	}

	// Check if token is a known pre-defined keyword prefix (e.g. "Blu" for "Blu-ray")
	keywordParts, found := p.tokenManager.keywordManager.findKeywordPartGroups(tkn.getValue())
	foundParts := false
	if found {
		foundParts = false
		for _, keywordGroup := range keywordParts {
			if retTkns, found := p.tokenManager.tokens.peekValuesAfter(p.tokenManager.tokens.getIndexOf(tkn), keywordGroup.seqParts); found {
				// Update token value
				seqPartsStr := ""
				for _, t := range retTkns {
					seqPartsStr += t.getValue()
				}
				tkn.setValue(mergeValues(tkn.getValue(), []string{seqPartsStr}))
				tkn.setIdentifiedKeywordCategory(keywordGroup.category)
				tkn.setKind(tokenKindWord)
				// Remove subsequent tokens
				for _, retTkn := range retTkns {
					p.tokenManager.tokens.removeAt(p.tokenManager.tokens.getIndexOf(retTkn))
				}
				foundParts = true
				break
			}
		}
	}

	if foundParts {
		return true
	}

	// Check if token is a known pre-defined STANDALONE keyword (e.g. "60FPS")
	if len(tkn.getValue()) > 1 {
		if keyword, found := p.tokenManager.keywordManager.findStandaloneKeywordByValue(tkn.getValue()); found {

			// When the priority is "normal", we only want to identify STANDALONE keywords that are not an anime type
			// That is because those are prone to false positives
			if priority == "normal" && keyword.isAnimeType() {
				return false
			}

			// when the priority is "normal", we only want to identify STANDALONE keywords that are not ambiguous
			if priority == "normal" && p.tokenManager.keywordManager.isKeywordAmbiguous(keyword) {
				return false
			}

			tkn.setIdentifiedKeywordCategory(keyword.category)
			return true
		}
	}

	return false

}

// collectMetadata collects the metadata elements from the parsed tokens.
// de-duplicates elements
func (p *parser) collectMetadata() {

	p.metadata.FileName = p.filename

	for _, tkn := range *p.tokenManager.tokens {

		switch tkn.IdentifiedKeywordCategory {
		case keywordCatYear:
			p.metadata.Year = tkn.getValue()
		case keywordCatReleaseVersion:
			p.metadata.ReleaseVersion = append(p.metadata.ReleaseVersion, tkn.getValue())
		case keywordCatFileChecksum:
			p.metadata.FileChecksum = tkn.getValue()
		case keywordCatVideoResolution:
			p.metadata.VideoResolution = tkn.getValue()
		case keywordCatReleaseGroup:
			p.metadata.ReleaseGroup = tkn.getValue()
		case keywordCatAudioTerm:
			p.metadata.AudioTerm = append(p.metadata.AudioTerm, tkn.getValue())
		case keywordCatAnimeType:
			p.metadata.AnimeType = append(p.metadata.AnimeType, tkn.getValue())
		case keywordCatVideoTerm:
			p.metadata.VideoTerm = append(p.metadata.VideoTerm, tkn.getValue())
		case keywordCatDeviceCompat:
			p.metadata.DeviceCompatibility = append(p.metadata.DeviceCompatibility, tkn.getValue())
		case keywordCatLanguage:
			p.metadata.Language = append(p.metadata.Language, tkn.getValue())
		case keywordCatSubtitles:
			p.metadata.Subtitles = append(p.metadata.Subtitles, tkn.getValue())
		case keywordCatSource:
			p.metadata.Source = append(p.metadata.Source, tkn.getValue())
		case keywordCatFileExtension:
			p.metadata.FileExtension = tkn.getValue()
		default:
		}

		switch tkn.MetadataCategory {
		case metadataTitle:
			p.metadata.Title = tkn.getValue()
		case metadataEpisodeTitle:
			p.metadata.EpisodeTitle = tkn.getValue()
		case metadataEpisodeNumber:
			p.metadata.EpisodeNumber = append(p.metadata.EpisodeNumber, tkn.getValue())
		case metadataOtherEpisodeNumber:
			p.metadata.OtherEpisodeNumber = append(p.metadata.OtherEpisodeNumber, tkn.getValue())
		case metadataEpisodeNumberAlt:
			p.metadata.EpisodeNumberAlt = append(p.metadata.EpisodeNumberAlt, tkn.getValue())
		case metadataSeason:
			p.metadata.Season = append(p.metadata.Season, tkn.getValue())
		case metadataPart:
			p.metadata.Part = append(p.metadata.Part, tkn.getValue())
		case metadataVolumeNumber:
			p.metadata.VolumeNumber = append(p.metadata.VolumeNumber, tkn.getValue())
		case metadataAnimeType:
			p.metadata.AnimeType = append(p.metadata.AnimeType, tkn.getValue())
		case metadataAudioTerm:
			p.metadata.AudioTerm = append(p.metadata.AudioTerm, tkn.getValue())
		case metadataDeviceCompat:
			p.metadata.DeviceCompatibility = append(p.metadata.DeviceCompatibility, tkn.getValue())
		case metadataLanguage:
			p.metadata.Language = append(p.metadata.Language, tkn.getValue())
		case metadataSubtitles:
			p.metadata.Subtitles = append(p.metadata.Subtitles, tkn.getValue())
		case metadataReleaseGroup:
			p.metadata.ReleaseGroup = tkn.getValue()
		case metadataReleaseVersion:
			p.metadata.ReleaseVersion = append(p.metadata.ReleaseVersion, tkn.getValue())
		case metadataSource:
			p.metadata.Source = append(p.metadata.Source, tkn.getValue())
		case metadataVideoResolution:
			p.metadata.VideoResolution = tkn.getValue()
		case metadataVideoTerm:
			p.metadata.VideoTerm = append(p.metadata.VideoTerm, tkn.getValue())
		default:
		}
	}

	if len(p.metadata.EpisodeNumber) == 0 && len(p.metadata.OtherEpisodeNumber) > 0 {
		p.metadata.EpisodeNumber = p.metadata.OtherEpisodeNumber
		p.metadata.OtherEpisodeNumber = nil
	}

}

func (p *parser) cleanUp() {

	if p.metadata.EpisodeNumber != nil {
		ret, vers := cleanNumbers(p.metadata.EpisodeNumber)
		p.metadata.EpisodeNumber = ret
		if len(vers) > 0 {
			p.metadata.ReleaseVersion = append(p.metadata.ReleaseVersion, vers)
		}
	}
	if p.metadata.Season != nil {
		ret, vers := cleanNumbers(p.metadata.Season)
		p.metadata.Season = ret
		if len(vers) > 0 {
			p.metadata.ReleaseVersion = append(p.metadata.ReleaseVersion, vers)
		}
	}
	if p.metadata.Part != nil {
		ret, vers := cleanNumbers(p.metadata.Part)
		p.metadata.Part = ret
		if len(vers) > 0 {
			p.metadata.ReleaseVersion = append(p.metadata.ReleaseVersion, vers)
		}
	}
	if p.metadata.VolumeNumber != nil {
		ret, vers := cleanNumbers(p.metadata.VolumeNumber)
		p.metadata.VolumeNumber = ret
		if len(vers) > 0 {
			p.metadata.ReleaseVersion = append(p.metadata.ReleaseVersion, vers)
		}
	}
}

func cleanNumbers(numbers []string) ([]string, string) {
	var cleaned []string
	var vers string
	for _, number := range numbers {
		num, v := cleanNumber(number)
		cleaned = append(cleaned, num)
		vers = v
	}
	return cleaned, vers
}

func cleanNumber(number string) (string, string) {
	if isDigitsOnly(number) {
		return number, ""
	}
	sepIdx := strings.IndexByte(number, '.')
	if sepIdx != -1 {
		number = number[:sepIdx]
		return number, "2"
	}
	sepIdx = strings.IndexByte(number, 'v')
	if sepIdx != -1 {
		s := number
		number = number[:sepIdx]
		pre, ok := strings.CutPrefix(s, number+"v")
		if !ok {
			pre = ""
		}
		return number, pre
	}
	sepIdx = strings.IndexByte(number, '\'')
	if sepIdx != -1 {
		number = number[:sepIdx]
		return number, "2"
	}
	return number, ""
}
