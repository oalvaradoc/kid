package factory

import (
	"git.multiverse.io/eventkit/kit/log"
	"github.com/beego/i18n"
	"path/filepath"
	"strings"
)

//InitLocaleLang initialize the i18n config files
func InitLocaleLang(pattern string) error {
	langFiles, err := filepath.Glob(pattern)

	if err != nil {
		return err
	}

	if len(langFiles) < 1 {
		log.Infos("Can't found locale language file, skip!")
		return nil
	}

	log.Infosf("start init %++v language files.", langFiles)

	for _, langFile := range langFiles {
		lang := strings.ReplaceAll(strings.Split(langFile, "e_")[1], ".ini", "")

		if err := i18n.SetMessage(lang, langFile); err != nil {
			return err
		}

		log.Debugsf("init %s successfully.", lang)
	}

	log.Infosf("init I18n config successfully.")
	return nil
}
