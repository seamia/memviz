package memviz

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os/user"
	"path"
	"sync"
	"syscall"
)

const (
	optionsFileName = "memviz.options"
)

type Settings struct {
	MaxStringLength          int                          `json:"maxStringLength"`
	MaxSliceLength           int                          `json:"maxSliceLength"`
	MaxMapEntries            int                          `json:"maxMapEntries"`
	Discard                  map[string]int               `json:"discard"`
	Substitute               map[string]map[string]string `json:"substitute"`
	Colors                   interface{}                  `json:"colors"`
	SuppresHeader            bool                         `json:"suppresHeader"`
	CollapsePointerNodes     bool                         `json:"collapsePointerNodes"`
	CollapseSingleSliceNodes bool                         `json:"collapseSingleSliceNodes"`
	ColorBackground          string                       `json:"colorBackground"` // transparent
	ColorDefault             string                       `json:"colorDefault"`    // whitesmoke
	FontName                 string                       `json:"fontName"`
	FontSize                 string                       `json:"fontSize"`
	LinkPointer              string                       `json:"link.pointer"`
	LinkArray                string                       `json:"link.array"`

	LoadedFrom string `json:"-"`
}

var (
	settings = Settings{
		MaxStringLength:          64,
		MaxSliceLength:           100,
		MaxMapEntries:            32,
		CollapsePointerNodes:     true,
		CollapseSingleSliceNodes: true,
		ColorBackground:          "transparent",
		ColorDefault:             "whitesmoke",
		FontName:                 "Cascadia Code",
		FontSize:                 "10",
	}

	guard          sync.Mutex
	colors         = make(map[string]string)
	settingsLoaded bool
)

func Options() *Settings {
	if !settingsLoaded {
		guard.Lock()
		defer guard.Unlock()

		if !settingsLoaded {
			settingsLoaded = true

			loadedFrom := "./" + optionsFileName
			data, err := ioutil.ReadFile(loadedFrom)
			if err != nil {
				loadedFrom = homeDir(optionsFileName)
				data, err = ioutil.ReadFile(loadedFrom)
			}

			if err == nil {
				if err := json.Unmarshal(data, &settings); err != nil {
					warning("error while loading config file (%v)", err)
				} else {
					settings.LoadedFrom = loadedFrom
				}
			} else {
				if perr, found := err.(*fs.PathError); found && perr.Err == syscall.ERROR_FILE_NOT_FOUND {
					// it is okay to have config file missing --> do not report this fact
				} else {
					warning("error while reading config file (%v)", err)
				}
			}

			if settings.Colors != nil {
				loadColors()
			}
		}
	}
	return &settings
}

/*
	the "colors" entry of the config file can be of these three types:
	1. string - name of the file containing color definitions
	2. []string - names of the files containing color definitions to be combined
	3. map[string]string - actual color definitions
*/
func loadColors() {
	input := settings.Colors
	if input == nil {
		return
	}

	var list []string
	switch actual := input.(type) {
	case string:
		list = append(list, actual)

	case []interface{}:
		for _, entry := range actual {
			if txt, converts := entry.(string); converts {
				list = append(list, txt)
			}
		}
	case map[string]interface{}:
		for k, v := range actual {
			if txt, converts := v.(string); converts {
				colors[k] = txt
			}
		}
		return

	default:
		warning("unrecognized format of Colors section of the config file (%v)", actual)
		return
	}

	for _, entry := range list {
		data, err := ioutil.ReadFile(entry)
		if err == nil {
			var loaded map[string]string
			if err := json.Unmarshal(data, &loaded); err == nil {
				for k, v := range loaded {
					colors[k] = v
				}
			} else {
				warning("error (%v) while loading config file (%v)", err, entry)
			}
		}
	}
}

func GetColor(name string) (string, bool) {
	if len(name) == 0 || len(colors) == 0 {
		return "", false
	}
	result, found := colors[name]
	return result, found
}

func homeDir(name string) string {
	if current, err := user.Current(); err == nil {
		return path.Join(current.HomeDir, name)
	}
	return name
}
