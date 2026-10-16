package fetch

import (
	"encoding/json"
	"net/http"
)

type JavaManifest []struct {
	Tag           string   `json:"tag"`
	Notes         string   `json:"notes,omitempty"`
	Java          string   `json:"java,omitempty"`
	Distribution  string   `json:"distribution,omitempty"`
	Jvm           string   `json:"jvm,omitempty"`
	Architectures []string `json:"architectures,omitempty"`
	Lts           bool     `json:"lts,omitempty"`
	Deprecated    bool     `json:"deprecated,omitempty"`
	Jdk           bool     `json:"jdk,omitempty"`
}

func FetchJavaManifest() (JavaManifest, error) {
	javaVersionListURL := "https://raw.githubusercontent.com/itzg/docker-minecraft-server/refs/heads/master/images.json"

	resp, err := http.Get(javaVersionListURL)
	if err != nil {
		return JavaManifest{}, err
	}
	defer resp.Body.Close()

	var versionList JavaManifest
	if err := json.NewDecoder(resp.Body).Decode(&versionList); err != nil {
		return JavaManifest{}, err
	}

	return versionList, nil
}

func (jm JavaManifest) GetStrList() []string {
	versions := []string{}
	for _, value := range jm {
		versions = append(versions, value.Tag)
	}
	return versions
}
