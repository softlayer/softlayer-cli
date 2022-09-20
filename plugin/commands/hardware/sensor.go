package hardware

import (
	"strconv"

	"github.com/spf13/cobra"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SensorCommand struct {
	*metadata.SoftlayerCommand
	HardwareManager managers.HardwareServerManager
	Command         *cobra.Command
	Discrete        bool
}

func NewSensorCommand(sl *metadata.SoftlayerCommand) (cmd *SensorCommand) {
	thisCmd := &SensorCommand{
		SoftlayerCommand: sl,
		HardwareManager:  managers.NewHardwareServerManager(sl.Session),
	}

	cobraCmd := &cobra.Command{
		Use:   "sensor " + T("IDENTIFIER"),
		Short: T("Retrieve a server’s hardware state via its internal sensors."),
		Args:  metadata.OneArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return thisCmd.Run(args)
		},
	}

	cobraCmd.Flags().BoolVar(&thisCmd.Discrete, "discrete", false, T("Show discrete units associated hardware sensor"))

	thisCmd.Command = cobraCmd
	return thisCmd
}

func (cmd *SensorCommand) Run(args []string) error {
	outputFormat := cmd.GetOutputFlag()

	hardwareId, err := strconv.Atoi(args[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware ID")
	}

	displayDiscrateTable := false
	if cmd.Discrete {
		displayDiscrateTable = true
	}

	sensorsData, err := cmd.HardwareManager.GetSensorData(hardwareId, "")
	if err != nil {
		return errors.NewAPIError(T("Failed to get hardware sensor data.\n"), err.Error(), 2)
	}

	temperatureTable := cmd.UI.Table([]string{T("Temperature (°C) Sensor"), T("Status"), T("Reading"), T("Critical Min"), T("Min"), T("Max"), T("Critical Max")})
	voltsTable := cmd.UI.Table([]string{T("Volts Sensor"), T("Status"), T("Reading"), T("Critical Min"), T("Min"), T("Max"), T("Critical Max")})
	wattsTable := cmd.UI.Table([]string{T("Watts Sensor"), T("Status"), T("Reading"), T("Critical Min"), T("Min"), T("Max"), T("Critical Max")})
	rpmTable := cmd.UI.Table([]string{T("RPM Sensor"), T("Status"), T("Reading"), T("Critical Min"), T("Min"), T("Max"), T("Critical Max")})
	discreteTable := cmd.UI.Table([]string{T("Discrete Sensor"), T("Status"), T("Reading")})

	for _, sensor := range sensorsData {
		if *sensor.SensorUnits == "degrees C" {
			temperatureTable.Add(
				utils.FormatStringPointer(sensor.SensorId),
				utils.FormatStringPointer(sensor.Status),
				utils.FormatStringPointer(sensor.SensorReading),
				utils.FormatStringPointer(sensor.LowerCritical),
				utils.FormatStringPointer(sensor.LowerNonCritical),
				utils.FormatStringPointer(sensor.UpperNonCritical),
				utils.FormatStringPointer(sensor.UpperCritical),
			)
		}

		if *sensor.SensorUnits == "Volts" {
			voltsTable.Add(
				utils.FormatStringPointer(sensor.SensorId),
				utils.FormatStringPointer(sensor.Status),
				utils.FormatStringPointer(sensor.SensorReading),
				utils.FormatStringPointer(sensor.LowerCritical),
				utils.FormatStringPointer(sensor.LowerNonCritical),
				utils.FormatStringPointer(sensor.UpperNonCritical),
				utils.FormatStringPointer(sensor.UpperCritical),
			)
		}

		if *sensor.SensorUnits == "Watts" {
			wattsTable.Add(
				utils.FormatStringPointer(sensor.SensorId),
				utils.FormatStringPointer(sensor.Status),
				utils.FormatStringPointer(sensor.SensorReading),
				utils.FormatStringPointer(sensor.LowerCritical),
				utils.FormatStringPointer(sensor.LowerNonCritical),
				utils.FormatStringPointer(sensor.UpperNonCritical),
				utils.FormatStringPointer(sensor.UpperCritical),
			)
		}

		if *sensor.SensorUnits == "RPM" {
			rpmTable.Add(
				utils.FormatStringPointer(sensor.SensorId),
				utils.FormatStringPointer(sensor.Status),
				utils.FormatStringPointer(sensor.SensorReading),
				utils.FormatStringPointer(sensor.LowerCritical),
				utils.FormatStringPointer(sensor.LowerNonCritical),
				utils.FormatStringPointer(sensor.UpperNonCritical),
				utils.FormatStringPointer(sensor.UpperCritical),
			)
		}

		if displayDiscrateTable {
			if *sensor.SensorUnits == "discrete" {
				discreteTable.Add(
					utils.FormatStringPointer(sensor.SensorId),
					utils.FormatStringPointer(sensor.Status),
					utils.FormatStringPointer(sensor.SensorReading),
				)
			}
		}
	}

	utils.PrintTable(cmd.UI, temperatureTable, outputFormat)
	cmd.UI.Print("\n")
	utils.PrintTable(cmd.UI, voltsTable, outputFormat)
	cmd.UI.Print("\n")
	utils.PrintTable(cmd.UI, wattsTable, outputFormat)
	cmd.UI.Print("\n")
	utils.PrintTable(cmd.UI, rpmTable, outputFormat)
	if displayDiscrateTable {
		cmd.UI.Print("\n")
		utils.PrintTable(cmd.UI, discreteTable, outputFormat)
	}
	return nil
}
