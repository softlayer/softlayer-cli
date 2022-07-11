package hardware

import (
	"strconv"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/urfave/cli"
	slErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/managers"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/metadata"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/utils"
)

type SensorCommand struct {
	UI              terminal.UI
	HardwareManager managers.HardwareServerManager
}

func NewSensorCommand(ui terminal.UI, hardwareManager managers.HardwareServerManager) (cmd *SensorCommand) {
	return &SensorCommand{
		UI:              ui,
		HardwareManager: hardwareManager,
	}
}

func (cmd *SensorCommand) Run(c *cli.Context) error {
	if c.NArg() != 1 {
		return errors.NewInvalidUsageError(T("This command requires one argument."))
	}

	outputFormat, err := metadata.CheckOutputFormat(c, cmd.UI)
	if err != nil {
		return err
	}

	hardwareId, err := strconv.Atoi(c.Args()[0])
	if err != nil {
		return slErr.NewInvalidSoftlayerIdInputError("Hardware ID")
	}

	displayDiscrateTable := false
	if c.IsSet("discrete") && c.Bool("discrete") {
		displayDiscrateTable = true
	}

	sensorsData, err := cmd.HardwareManager.GetSensorData(hardwareId, "")
	if err != nil {
		return cli.NewExitError(T("Failed to get hardware sensor data.\n")+err.Error(), 2)
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

func HardwareSensorMetaData() cli.Command {
	return cli.Command{
		Category:    "hardware",
		Name:        "sensor",
		Description: T("Retrieve a server’s hardware state via its internal sensors."),
		Usage: T(`${COMMAND_NAME} sl hardware sensor IDENTIFIER [OPTIONS]

EXAMPLE: 
   ${COMMAND_NAME} sl hardware sensor 123456`),
		Flags: []cli.Flag{
			cli.BoolFlag{
				Name:  "discrete",
				Usage: T("Show discrete units associated hardware sensor"),
			},
			metadata.OutputFlag(),
		},
	}
}
