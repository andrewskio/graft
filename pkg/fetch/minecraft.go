package fetch

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type MinecraftVersionDetails struct {
	ID          string `json:"id"`
	JavaVersion struct {
		Component    string `json:"component"`
		MajorVersion int    `json:"majorVersion"`
	} `json:"javaVersion"`
}

type MinecraftManifest struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`
	Versions []struct {
		ID          string    `json:"id"`
		Type        string    `json:"type"`
		URL         string    `json:"url"`
		Time        time.Time `json:"time"`
		ReleaseTime time.Time `json:"releaseTime"`
	} `json:"versions"`
}

func FetchMinecraftManifest() (MinecraftManifest, error) {
	manifestUrl := "https://launchermeta.mojang.com/mc/game/version_manifest.json"

	// fetch manifest
	resp, err := http.Get(manifestUrl)
	if err != nil {
		return MinecraftManifest{}, err
	}
	defer resp.Body.Close()

	// decode manifest
	var manifest MinecraftManifest
	if err := json.NewDecoder(resp.Body).Decode(&manifest); err != nil {
		return MinecraftManifest{}, err
	}
	return manifest, nil
}

func (mcm MinecraftManifest) FetchVersionDetails(mcVersion string) (MinecraftVersionDetails, error) {

	if mcVersion == "LATEST" {
		mcVersion = mcm.Latest.Release
	}

	if mcVersion == "SNAPSHOT" {
		mcVersion = mcm.Latest.Snapshot
	}

	versionURL := ""
	for _, version := range mcm.Versions {
		if version.ID == mcVersion {
			versionURL = version.URL
		}
	}
	if versionURL == "" {
		return MinecraftVersionDetails{}, errors.New("Version URL not found in manifest")
	}

	resp, err := http.Get(versionURL)
	if err != nil {
		return MinecraftVersionDetails{}, err
	}
	defer resp.Body.Close()

	var versionDetails MinecraftVersionDetails
	if err := json.NewDecoder(resp.Body).Decode(&versionDetails); err != nil {
		return MinecraftVersionDetails{}, err
	}
	return versionDetails, nil
}

func (mcm MinecraftManifest) GetStrList() []string {
	versions := []string{}
	versions = append(versions, "LATEST", "SNAPSHOT")
	for _, value := range mcm.Versions {
		versions = append(versions, value.ID)
	}
	return versions
}

// Take minecraft manifest version and returns compatible java versions
func (mcv MinecraftVersionDetails) GetJavaVersions(javaVersionList []string) []string {
	matchedJavaVersions := []string{}

	targetJavaVersion := strconv.Itoa(mcv.JavaVersion.MajorVersion)

	for _, jv := range javaVersionList {
		if strings.Contains(jv, targetJavaVersion) {
			matchedJavaVersions = append(matchedJavaVersions, jv)
		}
	}

	return matchedJavaVersions
}
