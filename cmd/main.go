package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/mkdym/docker-registry-viewer/client"
	"sort"
	"strings"
)

type Config struct {
	fn   string
	host string
	name string
	tag  string
	sort bool
}

func (c Config) String() string {
	return fmt.Sprintf("%#v", c)
}

var (
	g_config Config
)

func main() {
	HandleFlag()
	if err := Exec(); err != nil {
		panic(err)
	}
}

func HandleFlag() {
	flag.StringVar(&g_config.fn, "fn", "", `specify a function, can be one of followings:
		get_digest/get_digest2: get an image's digest. need name and tag
		list_tags: list a repo's tags. need name
		list_repos: list all repos
		list_all: list all repo and its tags
		delete: delete image tag. need name and tag
		get_info: get image info, need name and tag`)
	flag.StringVar(&g_config.host, "host", "", "specify registry host, eg, https://ep.wps.kingsoft.net:5000, 127.0.0.1:5000. if ssl on, must add 'https://'")
	flag.StringVar(&g_config.name, "name", "", "specify image name")
	flag.StringVar(&g_config.tag, "tag", "", "sepcify image tag")
	flag.BoolVar(&g_config.sort, "sort", false, "sort output")

	flag.Parse()
	if g_config.host == "" {
		panic("empty host")
	}

	//fmt.Println(g_config)
}

func Exec() error {
	host := strings.TrimPrefix(g_config.host, "http://")
	protocol := "http"

	if strings.HasPrefix(host, "https://") {
		protocol = "https"
		host = strings.TrimPrefix(host, "https://")
	}

	c, err := client.NewRegistryClient(protocol, host)
	if err != nil {
		return err
	}
	if err := c.Ping(); err != nil {
		return err
	}

	switch g_config.fn {
	case "get_digest":
		if g_config.name == "" || g_config.tag == "" {
			return errors.New("empty image name or tag")
		}

		resp, err := c.GetManifestV1(g_config.name, g_config.tag)
		if err != nil && err != client.ERR_IMAGE_NOT_FOUND {
			return err
		}

		if resp != nil {
			fmt.Print(resp.Digest)
		}

	case "get_digest2":
		if g_config.name == "" || g_config.tag == "" {
			return errors.New("empty image name or tag")
		}

		resp, err := c.GetManifestV2(g_config.name, g_config.tag)
		if err != nil && err != client.ERR_IMAGE_NOT_FOUND {
			return err
		}

		if resp != nil {
			fmt.Print(resp.Digest)
		}

	case "list_tags":
		if g_config.name == "" {
			return errors.New("empty image name")
		}

		tags, err := c.GetTags(g_config.name)
		if err != nil {
			return err
		}

		for _, s := range tags {
			fmt.Println(s)
		}

	case "list_repos":
		repos, err := c.GetCatalog()
		if err != nil {
			return err
		}

		for _, s := range repos {
			fmt.Println(s)
		}

	case "list_all":
		repos, err := c.GetCatalog()
		if err != nil {
			return err
		}

		for _, name := range repos {
			fmt.Printf("%-20s\t", name)

			if tags, err := c.GetTags(name); err == nil {
				fmt.Printf("%-2d\t\t", len(tags))

				if g_config.sort {
					sort.Strings(tags)
				}

				var tags_str string
				for _, tag := range tags {
					tags_str += tag + ","
				}

				tags_str = strings.Trim(tags_str, ",")
				fmt.Print(tags_str)
			}

			fmt.Print("\n")
		}

	case "delete":
		if g_config.name == "" || g_config.tag == "" {
			return errors.New("empty image name or tag")
		}

		if err := c.DeleteTag(g_config.name, g_config.tag); err != nil {
			return err
		}

		fmt.Println("success")

	case "get_info":
		if g_config.name == "" || g_config.tag == "" {
			return errors.New("empty image name or tag")
		}

		info, err := c.GetImageInfo(g_config.name, g_config.tag)
		if err != nil {
			return err
		}

		fmt.Println("name:", info.Name)
		fmt.Println("tag:", info.Tag)
		fmt.Println("DockerVersion:", info.DockerVersion)
		fmt.Println("DigestV1:", info.DigestV1)
		fmt.Println("DigestV2:", info.DigestV2)
		fmt.Println("ExposedPorts:", info.ExposedPorts)
		fmt.Println("Envs:", info.Envs)
		fmt.Println("Cmd:", info.Cmd)
		fmt.Println("Volumes:", info.Volumes)
		fmt.Println("WorkingDir:", info.WorkingDir)
		fmt.Println("Entrypoint:", info.Entrypoint)
		fmt.Println("Size:", info.HumanSize)

		fmt.Println("Layers", len(info.Layers))
		for index, layer := range info.Layers {
			fmt.Println("\tLayer", len(info.Layers)-index)
			fmt.Println("\t\tCreatedTime:", layer.CreatedTime)
			fmt.Println("\t\tBlobSum:", layer.BlobSum)
			fmt.Println("\t\tSize:", layer.HumanSize)
			fmt.Println("\t\tCmd:", layer.Cmd)
			fmt.Println("")
		}

	default:
		return errors.New("unknown function: " + g_config.fn)
	}

	return nil
}
