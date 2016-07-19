package main

import (
	"fmt"
	"git.oschina.net/mkdym/docker-registry-viewer/client"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/url"
	"os"
	"sort"
)

var (
	g_client *client.RegistryClient
)

func main() {
	listen_port := os.Getenv("LISTEN_PORT")
	if listen_port == "" {
		fmt.Println("default listen_port 49110, specify env LISTEN_PORT to change it")
		listen_port = "49110"
	}

	registry_host := os.Getenv("REGISTRY_HOST")
	if registry_host == "" {
		panic("empty registry_host, specify env REGISTRY_HOST")
	}

	registry_port := os.Getenv("REGISTRY_PORT")
	if registry_port == "" {
		fmt.Println("default registry_port 5000, specify env REGISTRY_PORT to change it")
		registry_port = "5000"
	}

	registry_protocol := "http"
	if ssl := os.Getenv("REGISTRY_SSL"); ssl == "on" {
		fmt.Println("registry ssl on")
		registry_protocol = "https"
	}

	registry_client, err := client.NewRegistryClient(registry_protocol, fmt.Sprintf("%s:%s", registry_host, registry_port))
	if err != nil {
		panic(err)
	}
	if err := registry_client.Ping(); err != nil {
		panic(err)
	}
	g_client = registry_client

	r := gin.Default()
	r.Static("/assets", "./resources/assets")
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	r.LoadHTMLGlob("./resources/templates/*")

	r.GET("/", handleGetRepos)
	r.GET("/tags/:repo", handleGetTags)
	r.GET("/detail/:repo/:tag", handleGetDetail)
	r.GET("/layers/:repo/:tag", handleGetLayers)
	r.GET("delete/:repo/:tag", handleDeleteImage)

	r.Run(":" + listen_port)
}

type RepoCountPair struct {
	Repo  string
	Count int
}

func handleGetRepos(c *gin.Context) {
	catalog, err := g_client.GetCatalog()
	if err != nil {
		c.String(http.StatusInternalServerError, "%s", err.Error())
		return
	}
	sort.Strings(catalog)

	repos := make([]RepoCountPair, 0, len(catalog))
	for _, name := range catalog {
		if tags, err := g_client.GetTags(name); err == nil {
			if len(tags) == 0 {
				fmt.Fprintln(os.Stderr, fmt.Sprintf("get tag of [%s] success, but Zero image", name))
			} else {
				repos = append(repos, RepoCountPair{Repo: name, Count: len(tags)})
			}
		} else {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("get tag of [%s] fail, error: %s", name, err.Error()))
		}
	}

	c.HTML(http.StatusOK, "repos", gin.H{"repos": repos})
}

type TimeSorterOfImageInfos []*client.ImageInfo

func (slice TimeSorterOfImageInfos) Len() int {
	return len(slice)
}

func (slice TimeSorterOfImageInfos) Less(i, j int) bool {
	return slice[i].CreatedTime < slice[j].CreatedTime
}

func (slice TimeSorterOfImageInfos) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

func handleGetTags(c *gin.Context) {
	repo, err := url.QueryUnescape(c.Param("repo"))
	if err != nil {
		c.String(http.StatusInternalServerError, "%s", err.Error())
		return
	}

	//fmt.Println("repo:", repo)

	tags, err := g_client.GetTags(repo)
	if err != nil {
		c.String(http.StatusInternalServerError, "%s", err.Error())
		return
	}

	tagsInfo := make([]*client.ImageInfo, 0, len(tags))
	for _, tag := range tags {
		if info, err := g_client.GetImageInfo(repo, tag); err == nil {
			tagsInfo = append(tagsInfo, info)
		} else {
			fmt.Fprintln(os.Stderr, fmt.Sprintf("get [%s:%s] image info fail, error: %s", repo, tag, err.Error()))
			tagsInfo = append(tagsInfo, &client.ImageInfo{})
		}
	}

	sort.Sort(sort.Reverse(TimeSorterOfImageInfos(tagsInfo)))
	c.HTML(http.StatusOK, "tags", gin.H{"tags": tagsInfo})
}

func handleGetDetail(c *gin.Context) {
	repo, err := url.QueryUnescape(c.Param("repo"))
	if err != nil {
		c.String(http.StatusInternalServerError, "%s", err.Error())
		return
	}

	tag, err := url.QueryUnescape(c.Param("tag"))
	if err != nil {
		c.String(http.StatusInternalServerError, "%s", err.Error())
		return
	}

	//fmt.Println("repo:", repo, ",tag:", tag)

	info, err := g_client.GetImageInfo(repo, tag)
	if err != nil {
		c.String(http.StatusInternalServerError, "%s", err.Error())
		return
	}

	c.HTML(http.StatusOK, "detail", gin.H{"info": info})
}

func handleGetLayers(c *gin.Context) {
	repo, err := url.QueryUnescape(c.Param("repo"))
	if err != nil {
		c.String(http.StatusInternalServerError, "%s", err.Error())
		return
	}

	tag, err := url.QueryUnescape(c.Param("tag"))
	if err != nil {
		c.String(http.StatusInternalServerError, "%s", err.Error())
		return
	}

	//fmt.Println("repo:", repo, ",tag:", tag)

	info, err := g_client.GetImageInfo(repo, tag)
	if err != nil {
		c.String(http.StatusInternalServerError, "%s", err.Error())
		return
	}

	c.HTML(http.StatusOK, "layers", gin.H{"layers": info.Layers})
}

func handleDeleteImage(c *gin.Context) {
	repo, err := url.QueryUnescape(c.Param("repo"))
	if err != nil {
		c.String(http.StatusInternalServerError, "%s", err.Error())
		return
	}

	tag, err := url.QueryUnescape(c.Param("tag"))
	if err != nil {
		c.String(http.StatusInternalServerError, "%s", err.Error())
		return
	}

	//fmt.Println("repo:", repo, ",tag:", tag)

	if err := g_client.DeleteTag(repo, tag); err != nil {
		c.String(http.StatusInternalServerError, "%s", err.Error())
		return
	}

	c.String(http.StatusOK, "delete %s:%s success", repo, tag)
}
