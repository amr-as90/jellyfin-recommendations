package recommender

import "fmt"

type CollectionNamer struct {
	Language string
}

// NewCollectionNamerFromState queries Jellyfin for the server UICulture
// field (to get the language) and initializes a CollectionNamer.
func (s *StateManager) NewCollectionNamerFromState() (*CollectionNamer, error) {
	config, err := getJellyfin[systemConfiguration](s, "/System/Configuration")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch server configuration: %w", err)
	}

	lang := config.UICulture
	if lang == "" {
		lang = config.PreferredMetadataLanguage
	}

	return &CollectionNamer{Language: lang}, nil
}

// FormatName accepts a username and formats the collection name based on the system's language and
// the user's name. For example, in English, the user 'Amr' would return 'Amr's recommendations'.
func (cn *CollectionNamer) FormatName(userName string) string {
	switch cn.Language {
	case "es":
		return fmt.Sprintf("Recomendaciones de %s", userName)
	case "de":
		return fmt.Sprintf("%ss Empfehlungen", userName)
	case "en":
		fallthrough
	default:
		return fmt.Sprintf("%s's Recommendations", userName)
	}
}
