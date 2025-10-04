package agent

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
)

func ProjectStructure(root string) (string, error) {
	files, err := os.ReadDir(root)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, file := range files {
		fmt.Fprintf(&sb, "|_%s\n", file.Name())
	}
	return sb.String(), nil
}

type Shell struct {
	Path string
	Name string
}

func DefaultShell() (*Shell, error) {
	var shellPath string
	var err error

	shellPath, err = shellFromEnv()
	if err == nil && shellPath != "" {
		return newShell(shellPath)
	}

	shellPath, err = shellFromPasswd()
	if err == nil && shellPath != "" {
		return newShell(shellPath)
	}

	shellPath, err = shellFromCommon()
	if err == nil && shellPath != "" {
		return newShell(shellPath)
	}

	return nil, fmt.Errorf("unable to determine user shell using any available method")
}

func shellFromEnv() (string, error) {
	shell := os.Getenv("SHELL")
	if shell == "" {
		return "", fmt.Errorf("SHELL environment variable not set")
	}

	if _, err := os.Stat(shell); err != nil {
		return "", fmt.Errorf("shell from environment variable does not exist: %w", err)
	}

	return shell, nil
}

func shellFromPasswd() (string, error) {
	if runtime.GOOS == "windows" {
		return "", fmt.Errorf("passwd file not available on Windows")
	}

	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	passwdFile, err := os.Open("/etc/passwd")
	if err != nil {
		return "", fmt.Errorf("failed to open passwd file: %w", err)
	}
	defer passwdFile.Close()

	scanner := bufio.NewScanner(passwdFile)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, ":")
		if len(fields) >= 7 && fields[0] == currentUser.Username {
			shell := fields[6]
			if shell != "" && shell != "/sbin/nologin" && shell != "/bin/false" {
				if _, err := os.Stat(shell); err != nil {
					return "", fmt.Errorf("shell from passwd file does not exist: %w", err)
				}
				return shell, nil
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("error reading passwd file: %w", err)
	}

	return "", fmt.Errorf("user shell not found in passwd file")
}

func shellFromCommon() (string, error) {
	switch runtime.GOOS {
	case "windows":
		psPath := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "WindowsPowerShell", "v1.0", "powershell.exe")
		if _, err := os.Stat(psPath); err == nil {
			return psPath, nil
		}

		output, err := exec.Command("where", "pwsh.exe").Output()
		if err == nil && len(output) > 0 {
			path := strings.TrimSpace(string(output))
			return path, nil
		}

		cmdPath := filepath.Join(os.Getenv("SYSTEMROOT"), "System32", "cmd.exe")
		if _, err := os.Stat(cmdPath); err == nil {
			return cmdPath, nil
		}

		return "", fmt.Errorf("could not find a valid shell on Windows")

	case "darwin":
		shells := []string{"/bin/zsh", "/bin/bash", "/bin/sh"}
		for _, shell := range shells {
			if _, err := os.Stat(shell); err == nil {
				return shell, nil
			}
		}

	default:
		shells := []string{"/bin/bash", "/bin/sh", "/usr/bin/bash", "/usr/bin/sh"}
		for _, shell := range shells {
			if _, err := os.Stat(shell); err == nil {
				return shell, nil
			}
		}
	}

	return "", fmt.Errorf("no default shell found for %s", runtime.GOOS)
}

func newShell(path string) (*Shell, error) {
	if path == "" {
		return nil, fmt.Errorf("empty shell path")
	}
	name := filepath.Base(path)

	if runtime.GOOS == "windows" {
		name = strings.TrimSuffix(name, filepath.Ext(name))
	}

	return &Shell{
		Path: path,
		Name: name,
	}, nil
}

type DevTools struct {
	VersionControl         []string
	PackageManagers        []string
	LanguageRuntimes       []string
	BuildTools             []string
	Testing                []string
	Database               []string
	ContainerOrchestration []string
	CloudInfrastructure    []string
	TextProcessing         []string
	FileOperations         []string
	NetworkHTTP            []string
	SystemMonitoring       []string
}

func AvailableDevTools() *DevTools {
	toolCategories := map[string][]string{
		"VersionControl": {
			"git", "gh",
		},
		"PackageManagers": {
			"npm", "yarn", "pnpm", "pip", "pip3", "pipenv", "poetry",
			"cargo", "go", "composer", "gem", "bundle", "maven", "gradle",
			"brew",
		},
		"LanguageRuntimes": {
			"node", "python", "python3", "java", "go", "php", "ruby", "swift", "kotlin",
			"scala", "dotnet",
		},
		"BuildTools": {
			"make", "cmake", "ninja", "webpack", "vite", "rollup", "parcel",
			"gulp", "bazel",
		},
		"Testing": {
			"jest", "mocha", "pytest", "phpunit", "rspec",
		},
		"Database": {
			"mysql", "psql", "sqlite3", "mongo", "redis-cli",
		},
		"ContainerOrchestration": {
			"docker", "docker-compose", "podman", "kubectl", "helm", "minikube", "nerdctl",
		},
		"CloudInfrastructure": {
			"aws", "gcloud", "az", "terraform", "ansible", "pulumi",
		},
		"TextProcessing": {
			"grep", "rg", "ag", "ack", "sed", "awk", "jq", "yq",
		},
		"FileOperations": {
			"find", "fd", "locate", "rsync", "scp", "tar", "zip", "unzip", "gzip",
		},
		"NetworkHTTP": {
			"curl", "wget", "nc", "ping", "dig", "nslookup",
		},
		"SystemMonitoring": {
			"ps", "top", "htop", "lsof", "netstat", "ss",
			"df", "du",
		},
	}

	result := &DevTools{}

	for _, tool := range toolCategories["VersionControl"] {
		if isToolAvailable(tool) {
			result.VersionControl = append(result.VersionControl, tool)
		}
	}

	for _, tool := range toolCategories["PackageManagers"] {
		if isToolAvailable(tool) {
			result.PackageManagers = append(result.PackageManagers, tool)
		}
	}

	for _, tool := range toolCategories["LanguageRuntimes"] {
		if isToolAvailable(tool) {
			result.LanguageRuntimes = append(result.LanguageRuntimes, tool)
		}
	}

	for _, tool := range toolCategories["BuildTools"] {
		if isToolAvailable(tool) {
			result.BuildTools = append(result.BuildTools, tool)
		}
	}

	for _, tool := range toolCategories["Testing"] {
		if isToolAvailable(tool) {
			result.Testing = append(result.Testing, tool)
		}
	}

	for _, tool := range toolCategories["Database"] {
		if isToolAvailable(tool) {
			result.Database = append(result.Database, tool)
		}
	}

	for _, tool := range toolCategories["ContainerOrchestration"] {
		if isToolAvailable(tool) {
			result.ContainerOrchestration = append(result.ContainerOrchestration, tool)
		}
	}

	for _, tool := range toolCategories["CloudInfrastructure"] {
		if isToolAvailable(tool) {
			result.CloudInfrastructure = append(result.CloudInfrastructure, tool)
		}
	}

	for _, tool := range toolCategories["TextProcessing"] {
		if isToolAvailable(tool) {
			result.TextProcessing = append(result.TextProcessing, tool)
		}
	}

	for _, tool := range toolCategories["FileOperations"] {
		if isToolAvailable(tool) {
			result.FileOperations = append(result.FileOperations, tool)
		}
	}

	for _, tool := range toolCategories["NetworkHTTP"] {
		if isToolAvailable(tool) {
			result.NetworkHTTP = append(result.NetworkHTTP, tool)
		}
	}

	for _, tool := range toolCategories["SystemMonitoring"] {
		if isToolAvailable(tool) {
			result.SystemMonitoring = append(result.SystemMonitoring, tool)
		}
	}

	return result
}

func isToolAvailable(tool string) bool {
	_, err := exec.LookPath(tool)
	return err == nil
}
