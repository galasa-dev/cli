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

type MonitorsDisableCmdValues struct {
}

type MonitorsDisableCommand struct {
	values *MonitorsDisableCmdValues
    cobraCommand *cobra.Command
}

// ------------------------------------------------------------------------------------------------
// Constructors methods
// ------------------------------------------------------------------------------------------------
func NewMonitorsDisableCommand(
    factory spi.Factory,
    monitorsDisableCommand spi.GalasaCommand,
    commsFlagSet GalasaFlagSet,
) (spi.GalasaCommand, error) {

    cmd := new(MonitorsDisableCommand)

    err := cmd.init(factory, monitorsDisableCommand, commsFlagSet)
    return cmd, err
}

// ------------------------------------------------------------------------------------------------
// Public methods
// ------------------------------------------------------------------------------------------------
func (cmd *MonitorsDisableCommand) Name() string {
    return COMMAND_NAME_MONITORS_DISABLE
}

func (cmd *MonitorsDisableCommand) CobraCommand() *cobra.Command {
    return cmd.cobraCommand
}

func (cmd *MonitorsDisableCommand) Values() interface{} {
	return cmd.values
}

// ------------------------------------------------------------------------------------------------
// Private methods
// ------------------------------------------------------------------------------------------------
func (cmd *MonitorsDisableCommand) init(factory spi.Factory, monitorsCommand spi.GalasaCommand, commsFlagSet GalasaFlagSet) error {
    var err error

	cmd.values = &MonitorsDisableCmdValues{}
    cmd.cobraCommand, err = cmd.createCobraCmd(factory, monitorsCommand, commsFlagSet.Values().(*CommsFlagSetValues))

    return err
}

func (cmd *MonitorsDisableCommand) createCobraCmd(
    factory spi.Factory,
    monitorsCommand spi.GalasaCommand,
    commsFlagSetValues *CommsFlagSetValues,
) (*cobra.Command, error) {

    var err error

    monitorsCommandValues := monitorsCommand.Values().(*MonitorsCmdValues)
    monitorsDisableCobraCmd := &cobra.Command{
        Use:     "disable",
        Short:   "Disable a monitor in the Galasa service",
        Long:    "Disables a monitor with the given name in the Galasa service",
        Aliases: []string{COMMAND_NAME_MONITORS_DISABLE},
        RunE: func(cobraCommand *cobra.Command, args []string) error {
			return cmd.executeMonitorsDisable(factory, monitorsCommand.Values().(*MonitorsCmdValues), commsFlagSetValues)
        },
    }

    addMonitorNameFlag(monitorsDisableCobraCmd, true, monitorsCommandValues)

    monitorsCommand.CobraCommand().AddCommand(monitorsDisableCobraCmd)

    return monitorsDisableCobraCmd, err
}

func (cmd *MonitorsDisableCommand) executeMonitorsDisable(
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

		log.Println("Galasa CLI - Disable monitors from the Galasa service")

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

				disableMonitorsFunc := func(apiClient *galasaapi.APIClient) error {
					return monitors.DisableMonitor(monitorsCmdValues.name, apiClient, byteReader)
				}
				err = commsClient.RunAuthenticatedCommandWithRateLimitRetries(disableMonitorsFunc)
			}
		}
	}

    return err
}
