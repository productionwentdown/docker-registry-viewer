package client

import (
	"encoding/json"
	"errors"
	"strconv"
)

type TagsResp struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

type ManifestV1Resp struct {
	Digest       string
	Name         string        `json:"name"`
	Tag          string        `json:"tag"`
	Architecture string        `json:"architecture"`
	FSLayers     []V1Layer     `json:"fsLayers"`
	Historys     []V1History   `json:"history"`
	Signatures   []V1Signature `json:"signatures"`
}

type V1Layer struct {
	BlobSum string `json:"blobSum"`
}

type V1History struct {
	V1Compatibility V1Compatibility `json:"v1Compatibility"`
}

type V1Compatibility struct {
	Architecture    string   `json:"architecture"`
	Author          string   `json:"author"`
	Config          V1Config `json:"config"`
	Container       string   `json:"container"`
	DockerVersion   string   `json:"docker_version"`
	Id              string   `json:"id"`
	Parent          string   `json:"parent"`
	CreatedTime     string   `json:"created"`
	ContainerConfig V1Config `json:"container_config"`
}

// http://attilaolah.eu/2013/11/29/json-decoding-in-go/
// Used to avoid recursion in UnmarshalJSON below.
type v1Compatibility V1Compatibility

func (v *V1Compatibility) UnmarshalJSON(b []byte) (err error) {
	unquotedstr, err := strconv.Unquote(string(b))
	if err != nil {
		return errors.New("unquote string fail, error: " + err.Error() + ", string:\n\n" + string(b) + "\n")
	}

	var v1Compatibility v1Compatibility
	if err := json.Unmarshal([]byte(unquotedstr), &v1Compatibility); err != nil {
		return err
	}
	*v = V1Compatibility(v1Compatibility)

	return nil
}

type V1Config struct {
	Hostname     string                 `json:"Hostname"`
	Domainname   string                 `json:"Domainname"`
	User         string                 `json:"User"`
	AttachStdin  bool                   `json:"AttachStdin"`
	AttachStdout bool                   `json:"AttachStdout"`
	AttachStderr bool                   `json:"AttachStderr"`
	ExposedPorts map[string]interface{} `json:"ExposedPorts"`
	Tty          bool                   `json:"Tty"`
	OpenStdin    bool                   `json:"OpenStdin"`
	StdinOnce    bool                   `json:"StdinOnce"`
	Envs         []string               `json:"Env"`
	Cmds         []string               `json:"Cmd"`
	Image        string                 `json:"Image"`
	Volumes      map[string]interface{} `json:"Volumes"`
	WorkingDir   string                 `json:"WorkingDir"`
	Entrypoint   []string               `json:"Entrypoint"`
	//OnBuild      []interface{}          `json:"OnBuild"`
	//Labels       interface{}            `json:"Labels"`
}

type V1Signature struct {
	//Header    map[string]interface{} `json:"header"`
	Signature string `json:"signature"`
	Protected string `json:"protected"`
}

type ManifestV2Resp struct {
	Digest    string
	MediaType string    `json:"mediaType"`
	Config    V2Config  `json:"config"`
	Layers    []V2Layer `json:"layers"`
}

type V2Config struct {
	MediaType string `json:"mediaType"`
	Size      uint64 `json:"size"`
	Digest    string `json:"digest"`
}

type V2Layer struct {
	MediaType string `json:"mediaType"`
	Size      uint64 `json:"size"`
	Digest    string `json:"digest"`
}

type CatalogResp struct {
	Repositories []string `json:"repositories"`
}
