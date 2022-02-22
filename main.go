package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	gu "github.com/deep2chain/goutils"
)

type params struct {
	FilePath string
	ModOnly  bool
}

func configParams(args []string) (*params, error) {
	switch len(args) {
	case 1:
		return nil, errors.New("filepath should be indicated")
	case 2:
		return &params{args[1], false}, nil
	default:
		return &params{args[1], true}, nil
	}
}

func analyze(deps map[string][]string) (*[]string, error) {
	counter := make(map[string]int)
	for _, items := range deps {
		for _, item := range items {
			v, ok := counter[item]
			if ok {
				counter[item] = v + 1
			}
			counter[item] = 1
		}
	}
	max_cnt := len(deps)
	var commons []string
	for k, v := range counter {
		if v == max_cnt {
			commons = append(commons, k)
			fmt.Println(v)
		}
	}
	return &commons, nil
}
func generateMods(repos []string) (*map[string][]string, error) {
	rootpath := gu.GetEnv("CODEPATH")
	curdir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	var deps = make(map[string][]string)
	for _, repo := range repos {
		key := strings.Split(repo, "/")[2]
		//
		repopath := fmt.Sprintf("%s/%s", rootpath, repo)

		if _, err := os.Stat(repopath); err != nil {
			fmt.Printf("%s not exists. You need to download.\n", repopath)
			continue
		}
		modfile := fmt.Sprintf("%s/go.mod", repopath)
		if _, err := os.Stat(modfile); err != nil {
			fmt.Printf("%s not exists. You need to download.\n", modfile)
			continue
		}
		lstfile := fmt.Sprintf("%s/%s", curdir, key)
		script := fmt.Sprintf("cd %s;go list -mod=readonly -m all >> %s.lst; cd -", repopath, lstfile)
		gu.RunScripts(script)
		items, err := gu.File2List(lstfile)
		if err != nil {
			continue
		}
		deps[key] = items
	}
	return &deps, nil
}

func main() {
	// chekc env $CODEPATH
	rootpath := gu.GetEnv("CODEPATH")
	if rootpath == "" {
		fmt.Println("set CODEPATH='Your Github Source Path'")
		return
	}
	// config params
	params, err := configParams(os.Args)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	// repo.lst to list
	items, err := gu.File2List(params.FilePath)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	deps, err := generateMods(items)
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	_, err = analyze(*deps)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

}
