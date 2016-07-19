package client

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

var (
	ERR_IMAGE_NOT_FOUND = errors.New("image not found")
)

type RegistryClient struct {
	host       string
	httpClient *http.Client
}

type registryResp struct {
	StatusCode    int
	StatusString  string
	Digest        string
	ContentLength uint64
	Body          string
}

type ImageInfo struct {
	Name          string
	Tag           string
	DockerVersion string
	CreatedTime   string
	DigestV1      string
	DigestV2      string
	ExposedPorts  []string
	Envs          []string
	Cmd           string
	Volumes       []string
	WorkingDir    string
	Entrypoint    string
	Size          uint64
	HumanSize     string
	Layers        []ImageLayer
}

type ImageLayer struct {
	BlobSum     string
	CreatedTime string
	Size        uint64
	HumanSize   string
	Cmd         string
}

func NewRegistryClient(protocol string, host string) (*RegistryClient, error) {
	httpClient := &http.Client{}

	if protocol == "https" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = &http.Client{Transport: tr}
	}

	return &RegistryClient{host: protocol + "://" + strings.Trim(host, "/\\"), httpClient: httpClient}, nil
}

func (c *RegistryClient) doRequest(method string, path string, headers map[string]string) (*registryResp, error) {
	req, err := http.NewRequest(method, c.host+"/v2/"+strings.Trim(path, "/\\"), nil)

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	httpResp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	bodyLenth, err := strconv.ParseUint(httpResp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		bodyLenth = 0
	}

	body, err := ioutil.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}

	return &registryResp{StatusCode: httpResp.StatusCode,
		StatusString:  httpResp.Status,
		Digest:        httpResp.Header.Get("Docker-Content-Digest"),
		ContentLength: bodyLenth,
		Body:          string(body)}, nil
}

func (c *RegistryClient) Ping() error {
	if _, err := c.doRequest(http.MethodGet, "", nil); err != nil {
		return err
	}
	return nil
}

func (c *RegistryClient) GetTags(name string) ([]string, error) {
	r, err := c.doRequest(http.MethodGet, name+"/tags/list", nil)
	if err != nil {
		return nil, err
	}

	if r.StatusCode == 404 {
		return nil, ERR_IMAGE_NOT_FOUND
	}

	if r.StatusCode != 200 {
		return nil, errors.New(r.StatusString)
	}

	var tags TagsResp
	if err := json.Unmarshal([]byte(r.Body), &tags); err != nil {
		return nil, errors.New("can not Unmarshal string\n\n" + r.Body + "\n\nerror: " + err.Error())
	}

	return tags.Tags, nil
}

func (c *RegistryClient) GetManifestV1(name string, reference string) (*ManifestV1Resp, error) {
	r, err := c.doRequest(http.MethodGet, name+"/manifests/"+reference, nil)
	if err != nil {
		return nil, err
	}

	if r.StatusCode == 404 {
		return nil, ERR_IMAGE_NOT_FOUND
	}

	if r.StatusCode != 200 {
		return nil, errors.New(r.StatusString)
	}

	var manifest ManifestV1Resp
	if err := json.Unmarshal([]byte(r.Body), &manifest); err != nil {
		return nil, errors.New("can not Unmarshal string\n\n" + r.Body + "\n\nerror: " + err.Error())
	}

	manifest.Digest = r.Digest

	return &manifest, nil
}

func (c *RegistryClient) GetManifestV2(name string, reference string) (*ManifestV2Resp, error) {
	headers := make(map[string]string)
	headers["Accept"] = "application/vnd.docker.distribution.manifest.v2+json"

	r, err := c.doRequest(http.MethodGet, name+"/manifests/"+reference, headers)
	if err != nil {
		return nil, err
	}

	if r.StatusCode == 404 {
		return nil, ERR_IMAGE_NOT_FOUND
	}

	if r.StatusCode != 200 {
		return nil, errors.New(r.StatusString)
	}

	var manifest ManifestV2Resp
	if err := json.Unmarshal([]byte(r.Body), &manifest); err != nil {
		return nil, errors.New("can not Unmarshal string\n\n" + r.Body + "\n\nerror: " + err.Error())
	}

	manifest.Digest = r.Digest

	return &manifest, nil
}

func (c *RegistryClient) GetCatalog() ([]string, error) {
	r, err := c.doRequest(http.MethodGet, "_catalog", nil)
	if err != nil {
		return nil, err
	}

	if r.StatusCode != 200 {
		return nil, errors.New(r.StatusString)
	}

	var catalog CatalogResp
	if err := json.Unmarshal([]byte(r.Body), &catalog); err != nil {
		return nil, errors.New("can not Unmarshal string\n\n" + r.Body + "\n\nerror: " + err.Error())
	}

	return catalog.Repositories, nil
}

func (c *RegistryClient) deleteByDigest(name string, digest string) error {
	headers := make(map[string]string)
	headers["Accept"] = "application/vnd.docker.distribution.manifest.v2+json"

	r, err := c.doRequest(http.MethodDelete, name+"/manifests/"+digest, headers)
	if err != nil {
		return err
	}

	if r.StatusCode != 202 {
		if r.Body != "" {
			return errors.New(r.Body)
		}
		return errors.New(r.StatusString)
	}

	return nil
}

func (c *RegistryClient) DeleteTag(name string, tag string) error {
	m, err := c.GetManifestV2(name, tag)
	if err != nil {
		return errors.New("can not get image[" + name + ":" + tag + "] digest for delete, error: " + err.Error())
	}

	return c.deleteByDigest(name, m.Digest)
}

func (c *RegistryClient) getBlobSize(name string, digest string) (uint64, error) {
	r, err := c.doRequest(http.MethodHead, name+"/blobs/"+digest, nil)
	if err != nil {
		return 0, err
	}

	if r.StatusCode == 404 {
		return 0, ERR_IMAGE_NOT_FOUND
	}

	if r.StatusCode != 200 {
		return 0, errors.New(r.StatusString)
	}

	return r.ContentLength, nil
}

func (c *RegistryClient) GetImageInfo(name string, tag string) (*ImageInfo, error) {
	m_v1, err := c.GetManifestV1(name, tag)
	if err != nil {
		return nil, errors.New("can not get image[" + name + ":" + tag + "] manifest(V1), error: " + err.Error())
	}

	if len(m_v1.FSLayers) == 0 || len(m_v1.Historys) == 0 || len(m_v1.FSLayers) != len(m_v1.Historys) {
		return nil, errors.New("invalid manifest(V1), empty layers or history or not equal numbers")
	}

	m_v2, err := c.GetManifestV2(name, tag)
	if err != nil {
		return nil, errors.New("can not get image[" + name + ":" + tag + "] manifest(V2), error: " + err.Error())
	}

	/*
		if len(m_v2.Layers) != len(m_v1.FSLayers) {
			return nil, errors.New("invalid manifest(V2), layers number of v1 and v2 not equal")
		}
	*/

	var info ImageInfo
	info.Name = name
	info.Tag = tag
	info.DockerVersion = m_v1.Historys[0].V1Compatibility.DockerVersion
	info.CreatedTime = m_v1.Historys[0].V1Compatibility.CreatedTime
	info.DigestV1 = m_v1.Digest
	info.DigestV2 = m_v2.Digest

	for k, _ := range m_v1.Historys[0].V1Compatibility.Config.ExposedPorts {
		info.ExposedPorts = append(info.ExposedPorts, k)
	}

	info.Envs = m_v1.Historys[0].V1Compatibility.Config.Envs
	info.Cmd = strings.Join(m_v1.Historys[0].V1Compatibility.Config.Cmds, ", ")

	for k, _ := range m_v1.Historys[0].V1Compatibility.Config.Volumes {
		info.Volumes = append(info.Volumes, k)
	}

	info.WorkingDir = m_v1.Historys[0].V1Compatibility.Config.WorkingDir
	info.Entrypoint = strings.Join(m_v1.Historys[0].V1Compatibility.Config.Entrypoint, ", ")

	//m_v1.FSLayers 是有顺序的，时间倒序
	for index, _ := range m_v1.FSLayers {
		var layer ImageLayer
		layer.BlobSum = m_v1.FSLayers[index].BlobSum
		layer.CreatedTime = m_v1.Historys[index].V1Compatibility.CreatedTime
		layer.Cmd = strings.Join(m_v1.Historys[index].V1Compatibility.ContainerConfig.Cmds, ", ")

		//v1中的blobsum在v2中不一定有，所以还是取v1中blob的length
		/*
			for _, v2layer := range m_v2.Layers {
				if v2layer.Digest == layer.BlobSum {
					layer.Size = v2layer.Size
					layer.HumanSize = humanSize(layer.Size)
					break
				}
			}
		*/
		layer.Size, _ = c.getBlobSize(info.Name, layer.BlobSum)
		layer.HumanSize = humanSize(layer.Size)

		info.Layers = append(info.Layers, layer)
		info.Size += layer.Size
	}
	info.HumanSize = humanSize(info.Size)

	return &info, nil
}
