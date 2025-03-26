/*
 * Copyright contributors to the Galasa project
 *
 * SPDX-License-Identifier: EPL-2.0
 */
package cmd

import (
	"log"

	"github.com/galasa-dev/cli/pkg/api"
	"github.com/galasa-dev/cli/pkg/galasaapi"
	"github.com/galasa-dev/cli/pkg/monitors"
	"github.com/galasa-dev/cli/pkg/spi"
	"github.com/galasa-dev/cli/pkg/utils"
	"github.com/spf13/cobra"
)

type MonitorsEnableCmdValues struct {
}

type MonitorsEnableCommand struct {
	values *MonitorsEnableCmdValues
    cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewMonitorsEnableCommand(
    factory spi.Factory,
    monitorsEnableCommand spi.GalasaCommand,
    commsFlagSet GalasaFlagSet,
) (spi.GalasaCommand, error) {

    cmd := new(MonitorsEnableCommand)

    err := cmd.init(factory, monitorsEnableCommand, commsFlagSet)
    return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *MonitorsEnableCommand) Name() string {
    return COMMAND_NAME_MONITORS_ENABLE
}

func (cmd *MonitorsEnableCommand) CobraCommand() *cobra.Command {
    return cmd.cobraCommand
}

func (cmd *MonitorsEnableCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *MonitorsEnableCommand) init(factory spi.Factory, monitorsCommand spi.GalasaCommand, commsFlagSet GalasaFlagSet) error {
    var err error

	cmd.values = &MonitorsEnableCmdValues{}
    cmd.cobraCommand, err = cmd.createCobraCmd(factory, monitorsCommand, commsFlagSet.Values().(*CommsFlagSetValues))

    return err
}

func (cmd *MonitorsEnableCommand) createCobraCmd(
    factory spi.Factory,
    monitorsCommand spi.GalasaCommand,
    commsFlagSetValues *CommsFlagSetValues,
) (*cobra.Command, error) {

    var err error

    monitorsCommandValues := monitorsCommand.Values().(*MonitorsCmdValues)
    monitorsEnableCobraCmd := &cobra.Command{
        Use:     "enable",
        Short:   "Enable a monitor in the Galasa service",
        Long:    "Enables a given monitor in the Galasa service",
        Aliases: []string{COMMAND_NAME_MONITORS_ENABLE},
        RunE: func(cobraCommand *cobra.Command, args []string) error {
			return cmd.executeMonitorsEnable(factory, monitorsCommand.Values().(*MonitorsCmdValues), commsFlagSetValues)
        },
    }

    addMonitorNameFlag(monitorsEnableCobraCmd, true, monitorsCommandValues)

    monitorsCommand.CobraCommand().AddCommand(monitorsEnableCobraCmd)

    return monitorsEnableCobraCmd, err
}

func (cmd *MonitorsEnableCommand) executeMonitorsEnable(
    factory spi.Factory,
    monitorsCmdValues *MonitorsCmdValues,
    commsFlagSetValues *CommsFlagSetValues,
) error {

    var err error
    // Operations on the file system will all be relative to the current folder.
    fileSystem := factory.GetFileSystem()

	err = utils.CaptureLog(fileSystem, commsFlagSetValues.logFileName)
	if err == nil {
		commsFlagSetValues.isCapturingLogs = true

		log.Println("Galasa CLI - Enable monitors from the Galasa service")

		env := factory.GetEnvironment()

		var galasaHome spi.GalasaHome
		galasaHome, err = utils.NewGalasaHome(fileSystem, env, commsFlagSetValues.CmdParamGalasaHomePath)
		if err == nil {

			var commsClient api.APICommsClient
			commsClient, err = api.NewAPICommsClient(
				commsFlagSetValues.bootstrap,
				commsFlagSetValues.maxRetries,
				commsFlagSetValues.retryBackoffSeconds,
				factory,
				galasaHome,
			)

			if err == nil {

				byteReader := factory.GetByteReader()

				enableMonitorsFunc := func(apiClient *galasaapi.APIClient) error {
					return monitors.EnableMonitor(monitorsCmdValues.name, apiClient, byteReader)
				}
				err = commsClient.RunAuthenticatedCommandWithRateLimitRetries(enableMonitorsFunc)
			}
		}
	}

    return err
}
