package main

import (
	"errors"
	"fmt"
)

type Package struct {
	Name         string
	Dependencies []string
}

func buildDependencyGraph(packages map[string]Package) (map[string][]string, error) {
	dependencyGraph := make(map[string][]string)

	for pkgName, pkg := range packages {
		for _, dep := range pkg.Dependencies {

			if _, exists := packages[dep]; !exists {
				return nil, errors.New(fmt.Sprintf("зависимость %s для пакета %s не найдена", dep, pkgName))
			}

			dependencyGraph[pkgName] = append(dependencyGraph[pkgName], dep)
		}
	}

	return dependencyGraph, nil
}

func hasCycle(dependencyGraph map[string][]string, current string, visited map[string]bool, stack map[string]bool) bool {
	if stack[current] {
		return true
	}

	if visited[current] {
		return false
	}

	visited[current] = true
	stack[current] = true

	for _, dep := range dependencyGraph[current] {
		if hasCycle(dependencyGraph, dep, visited, stack) {
			return true
		}
	}

	stack[current] = false
	return false
}

func detectCycles(dependencyGraph map[string][]string) bool {
	visited := make(map[string]bool)
	stack := make(map[string]bool)

	for pkg := range dependencyGraph {
		if hasCycle(dependencyGraph, pkg, visited, stack) {
			return true
		}
	}

	return false
}

func main() {

	packages := map[string]Package{
		"packageA": {Name: "packageA", Dependencies: []string{"packageB", "packageC"}},
		"packageB": {Name: "packageB", Dependencies: []string{"packageD"}},
		"packageC": {Name: "packageC", Dependencies: []string{"packageA"}},
		"packageD": {Name: "packageD", Dependencies: []string{"packageE"}},
		"packageE": {Name: "packageE", Dependencies: []string{}},
	}

	dependencyGraph, err := buildDependencyGraph(packages)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Println("Граф зависимостей:")
	for pkg, deps := range dependencyGraph {
		fmt.Printf("%s: %v\n", pkg, deps)
	}

	if detectCycles(dependencyGraph) {
		fmt.Println("Обнаружены циклические зависимости")
	} else {
		fmt.Println("Циклов не обнаружено")
	}
}
