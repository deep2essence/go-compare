package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	gu "github.com/deep2chain/goutils"
)

type params struct {
	FilePath      string
	IgnoreVersion bool
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

func analyze(deps map[string][]string, IgnoreVersion bool) (*[]string, error) {
	counter := make(map[string]int)
	for _, repos := range deps {
		for _, dep := range repos {
			var item = dep
			if IgnoreVersion {
				item = strings.Split(item, " ")[0]
			}
			v, ok := counter[item]
			if ok {
				counter[item] = v + 1
				continue
			}
			counter[item] = 1
		}
	}
	max_cnt := len(deps)
	fmt.Println(max_cnt)
	var commons []string
	for k, v := range counter {
		if v == max_cnt {
			commons = append(commons, k)
			fmt.Println(k)
		}
	}
	gu.List2File("common.lst", gu.ToInterfaceSlice(commons))
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
		var splited []string
		if splited = strings.Split(repo, "/"); len(splited) < 3 {
			continue
		}
		key := splited[2]
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
		lstfile := fmt.Sprintf("%s/%s.lst", curdir, key)
		script := fmt.Sprintf("cd %s;go list -mod=readonly -m all >> %s; cd -", repopath, lstfile)
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
		fmt.Println("export CODEPATH='Your Github Source Path'")
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
	_, err = analyze(*deps, params.IgnoreVersion)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

}
