package i18n

import (
	"ValorantRPC/pkg/errors"
	"embed"
	"encoding/json"
	"fmt"
	"path"
)

//go:embed locales/*.json
var localeFS embed.FS

type Translator interface {
	T(key string, args ...any) string
	Lang() string
}

type Service struct {
	lang         string
	translations map[string]string
}

func NewService(lang string) (*Service, error) {
	filePath := path.Join("locales", lang+".json")
	data, err := localeFS.ReadFile(filePath)
	if err != nil {
		if lang == "en" {
			return nil, errors.New("i18n.NewService", errors.ErrLocaleNotFound)
		}

		return NewService("en")
	}

	translations := make(map[string]string)
	err = json.Unmarshal(data, &translations)
	if err != nil {
		return nil, errors.New("i18n.NewService", errors.ErrLocaleMalformed)
	}

	return &Service{lang: lang, translations: translations}, nil
}

func (s *Service) T(key string, args ...any) string {
	val, ok := s.translations[key]
	if !ok {
		return key
	}
	if len(args) > 0 {
		return fmt.Sprintf(val, args...)
	}
	return val
}

func (s *Service) Lang() string {
	return s.lang
}
