package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cloudfoundry/cli/cf/configuration/config_helpers"
	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/plugin"
)

type AppSearchResoures struct {
	Metadata AppSearchMetaData `json:"metadata"`
}

type AppSearchMetaData struct {
	Guid string `json:"guid"`
	Url  string `json:"url"`
}

type AppSearchResults struct {
	Resources []AppSearchResoures `json:"resources"`
}

type AppEnv struct {
	System map[string]interface{} `json:"system_env_json"`
}

type Service struct {
	Name        string                 // name of the service
	Label       string                 // label of the service
	Tags        []string               // tags for the service
	Plan        string                 // plan of the service
	Credentials map[string]interface{} // credentials for the service
}

type Ports map[string]string

type Services map[string][]Service

func fatalIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stdout, "error:", err)
		os.Exit(1)
	}
}

func main() {
	plugin.Start(&LogsearchPlugin{})
}

type LogsearchPlugin struct{}

func (c *LogsearchPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	err := checkArgs(cliConnection, args)
	if err != nil {
		os.Exit(1)
	}

	if args[0] == "search-logs" {
		checkService(cliConnection, args[1])
	}
}

func (c *LogsearchPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "logsearch",
		Version: plugin.VersionType{
			Major: 0,
			Minor: 1,
			Build: 0,
		},
		Commands: []plugin.Command{
			plugin.Command{
				Name:     "search-logs",
				HelpText: "search and display application logs",

				UsageDetails: plugin.Usage{
					Usage: "search-logs <appname>",
				},
			},
		},
	}
}

func findAppGuid(cliConnection plugin.CliConnection, appName string) string {

	confRepo := core_config.NewRepositoryFromFilepath(config_helpers.DefaultFilePath(), fatalIf)
	spaceGuid := confRepo.SpaceFields().Guid

	appQuery := fmt.Sprintf("/v2/spaces/%v/apps?q=name:%v&inline-relations-depth=1", spaceGuid, appName)
	cmd := []string{"curl", appQuery}

	output, _ := cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	res := &AppSearchResults{}
	json.Unmarshal([]byte(strings.Join(output, "")), &res)

	return res.Resources[0].Metadata.Guid
}

func checkService(cliConnection plugin.CliConnection, appName string) {
	guid := findAppGuid(cliConnection, appName)
	appQuery := fmt.Sprintf("/v2/apps/%v/env", guid)
	cmd := []string{"curl", appQuery}
	output, _ := cliConnection.CliCommandWithoutTerminalOutput(cmd...)
	appEnvs := AppEnv{}
	json.Unmarshal([]byte(output[0]), &appEnvs)
	str, err := json.Marshal(appEnvs.System["VCAP_SERVICES"])
	if err != nil {
		return
	}
	var services Services
	json.Unmarshal([]byte(str), &services)
	fmt.Println(services)
	ports := services["logstash14"][0].Credentials["ports"].(map[string]interface{})
	fmt.Println(services["logstash14"][0].Credentials["hostname"])
	fmt.Println(ports["514/tcp"].(string))
	return
}

func checkArgs(cliConnection plugin.CliConnection, args []string) error {
	if len(args) < 2 {
		cliConnection.CliCommand(args[0], "-h")
		return errors.New("Appname is needed")
	}
	return nil
}
