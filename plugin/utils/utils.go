package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	bmxErr "github.ibm.com/SoftLayer/softlayer-cli/plugin/errors"

	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/terminal"
	"github.com/IBM-Cloud/ibm-cloud-cli-sdk/bluemix/trace"
	"github.com/softlayer/softlayer-go/datatypes"
	"github.com/softlayer/softlayer-go/services"
	"github.com/softlayer/softlayer-go/session"
	"github.com/softlayer/softlayer-go/sl"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
	"github.ibm.com/SoftLayer/softlayer-cli/plugin/progress_bar"
)

const (
	EMPTY_VALUE  = "-"
	EMPTY_STRING = ""
)

// TODO support resolving guid to integer id
func ResolveVirtualGuestId(identifier string) (int, error) {
	id, err := strconv.Atoi(identifier)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func ResolveImageId(session *session.Session, id int) (string, error) {
	service := services.GetVirtualGuestBlockDeviceTemplateGroupService(session)
	image, err := service.Id(id).GetObject()
	if err != nil {
		return EMPTY_STRING, err
	}
	if image.GlobalIdentifier == nil {
		return EMPTY_STRING, errors.New(T("Image global identifier not found"))
	}
	return *image.GlobalIdentifier, nil
}

func ResolveVlanId(identifier string) (int, error) {
	id, err := strconv.Atoi(identifier)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func ResolveSubnetId(identifier string) (int, error) {
	id, err := strconv.Atoi(identifier)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func ResolveGloablIPId(identifier string) (int, error) {
	id, err := strconv.Atoi(identifier)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func StringSliceToString(slice []string) string {
	if len(slice) == 0 {
		return EMPTY_STRING
	}
	return strings.Trim(strings.Replace(fmt.Sprint(slice), " ", ",", -1), "[]")
}

func IntSliceToString(slice []int) string {
	if len(slice) == 0 {
		return EMPTY_STRING
	}
	return strings.Trim(strings.Replace(fmt.Sprint(slice), " ", ",", -1), "[]")
}

func TagRefsToString(tags []datatypes.Tag_Reference) string {
	var names []string
	for _, tag := range tags {
		if tag.Tag != nil && tag.Tag.Name != nil {
			names = append(names, *tag.Tag.Name)
		}
	}
	return StringSliceToString(names)
}

func StringSliceToIntSlice(slice []string) ([]int, error) {
	intSlice := make([]int, len(slice))
	for index, str := range slice {
		value, err := strconv.Atoi(str)
		if err != nil {
			return intSlice, err
		}
		intSlice[index] = value
	}
	return intSlice, nil
}

func StringInSlice(value string, slice []string) int {
	for i, s := range slice {
		if value == s {
			return i
		}
	}
	return -1
}

func IntInSlice(value int, slice []int) int {
	for i, s := range slice {
		if value == s {
			return i
		}
	}
	return -1
}

func IntSliceToStringSlice(slice []int) []string {
	var result []string
	for _, i := range slice {
		result = append(result, strconv.Itoa(i))
	}
	return result
}

// the first bool value indicates whether the source slice are all in target slice
// the second int value indicates the index of the source slice value is not in the target slice
func SliceInSlice(source, target []string) (bool, int) {
	exist := make([]bool, len(source))
	for i, s := range source {
		exist[i] = false
		for _, t := range target {
			if s == t {
				exist[i] = true
				break
			}
		}
	}

	for index, value := range exist {
		if value == false {
			return false, index
		}
	}
	return true, -1
}

// merge 2 slices to 1 without duplicate value
func MergeSlice(s1, s2 []string) []string {
	for _, s := range s2 {
		if StringInSlice(s, s1) == -1 {
			s1 = append(s1, s)
		}
	}
	return s1
}

func MergeAndSortSlice(s1, s2 []string) []string {
	s := MergeSlice(s1, s2)
	sort.Strings(s)
	return s
}

// Converts number of bytes to a string in gigabytes.
func B2GB(bytes int) string {
	return fmt.Sprintf("%.2f%s", float32(bytes)/1024/1024/1024, "G")
}

func FormatBoolPointer(value *bool) string {
	if value == nil {
		return EMPTY_VALUE
	}
	return strconv.FormatBool(sl.Get(value).(bool))
}

func FormatBoolPointerToYN(value *bool) string {
	if value == nil {
		return EMPTY_VALUE
	}
	if *value == true {
		return T("Yes")
	}
	return T("No")
}

func FormatStringPointer(value *string) string {
	if value == nil {
		return EMPTY_VALUE
	}
	return sl.Get(value).(string)
}

func FormatStringPointerName(value *string) string {
	if value == nil {
		return EMPTY_STRING
	}
	return sl.Get(value).(string)
}

func FormatIntPointer(value *int) string {
	if value == nil {
		return EMPTY_VALUE
	}
	return strconv.Itoa(sl.Get(value).(int))
}

func FormatIntPointerName(value *int) string {
	if value == nil {
		return EMPTY_STRING
	}
	return strconv.Itoa(sl.Get(value).(int))
}

func FormatUIntPointer(value *uint) string {
	if value == nil {
		return EMPTY_VALUE
	}
	return fmt.Sprintf("%d", sl.Get(value).(uint))
}

func FormatSLFloatPointerToInt(value *datatypes.Float64) string {
	if value == nil {
		return EMPTY_VALUE
	}
	return fmt.Sprintf("%d", int(sl.Get(value).(datatypes.Float64)))
}

func FormatSLFloatPointerToFloat(value *datatypes.Float64) string {
	if value == nil {
		return EMPTY_VALUE
	}
	return fmt.Sprintf("%f", float64(sl.Get(value).(datatypes.Float64)))
}

func FormatSLTimePointer(value *datatypes.Time) string {
	if value == nil {
		return EMPTY_VALUE
	}
	return value.UTC().Format(time.RFC3339)
}

func Bool2Int(value bool) int {
	if value {
		return 1
	}
	return 0
}

func ReplaceUIntPointerValue(value *uint, newValue string) string {

	if UIntPointertoUInt(value) > 0 {
		return newValue
	}
	return EMPTY_VALUE
}

// TODO: Once refactor is done, remove ValidateColumns and rename ValidateColumns2 to it.
func ValidateColumns2(sortby string, columns []string, defaultColumns []string, optionalColumns, sortColumns []string) ([]string, error) {
	if sortby != EMPTY_STRING && StringInSlice(sortby, sortColumns) == -1 {
		return nil, bmxErr.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
	}
	allColumns := append(defaultColumns, optionalColumns...)
	if exist, index := SliceInSlice(columns, allColumns); len(columns) > 0 && exist == false {
		return nil, bmxErr.NewInvalidUsageError(T("--column {{.Column}} is not supported.", map[string]interface{}{"Column": columns[index]}))
	}

	if len(columns) == 0 {
		return defaultColumns, nil
	}
	return columns, nil
}

func GetMask(maskMap map[string]string, columns []string, sortBy string) string {

	if sortBy != EMPTY_STRING && StringInSlice(sortBy, columns) == -1 {
		columns = append(columns, sortBy)
	}

	var mask []string
	for _, column := range columns {
		mask = append(mask, maskMap[column])
	}
	return strings.Join(mask, ",")
}

func GetColumnHeader(showColumns []string) []string {
	var columns []string
	for _, col := range showColumns {
		columns = append(columns, T(col))
	}
	return columns
}

// FailWithError ...
func FailWithError(message string, ui terminal.UI) error {
	trace.Logger.Println(T("Failed due to error: "), message)
	ui.Print(terminal.FailureColor(T("FAILED")))
	msg := fmt.Sprintf("%s\n", message)
	ui.Print(msg)
	return errors.New(EMPTY_STRING)
}

func StructToMap(struc Access) (map[string]string, error) {
	result := make(map[string]string)

	j, err := json.Marshal(struc)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(j, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func UIntPointertoUInt(value *uint) uint {
	if value == nil {
		return 0
	}
	return *value
}
func StringPointertoString(value *string) string {
	if value == nil {
		return EMPTY_STRING
	}
	return *value
}
func IntPointertoInt(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}
func BoolPointertoBool(value *bool) bool {
	if value == nil {
		return false
	}
	return *value
}
func IsEmptyString(value string) bool {
	return value == EMPTY_STRING
}
func WordInList(wordList []string, key string) bool {
	for _, word := range wordList {
		if word == key {
			return true
		}
	}
	return false
}

func PrintTableWithTitle(ui terminal.UI, table terminal.Table, bufEvent *bytes.Buffer, title string, outputFormat string) string {
	tableTitle := ui.Table([]string{T(title)})
	if outputFormat == "JSON" {
		table.PrintJson()
		tableTitle.Add(bufEvent.String())
		tableTitle.PrintJson()
		return ""
	}
	if outputFormat == "CSV" {
		err := table.PrintCsv()
		if err != nil {
			fmt.Println("Error:", err)
			return ""
		}
		tableTitle.Add(bufEvent.String())
		err = tableTitle.PrintCsv()
		if err != nil {
			fmt.Println("Error:", err)
			return ""
		}
		return ""
	}
	table.Print()
	tableTitle.Add(bufEvent.String())
	tableTitle.Print()
	return ""

}

func PrintTable(ui terminal.UI, table terminal.Table, outputFormat string) {
	if outputFormat == "JSON" {
		table.PrintJson()
		return
	}
	if outputFormat == "CSV" {
		err := table.PrintCsv()
		if err != nil {
			fmt.Println("Error:", err)

		}
		return
	}
	table.Print()
	return
}

func ShortenString(ugly_string string) string {
	limit := 80
	if len(ugly_string) > limit {
		return ugly_string[:limit] + "..."
	}
	return ugly_string
}

func ShortenStringWithLimit(ugly_string string, limit int) string {
	if len(ugly_string) > limit {
		return ugly_string[:limit] + "..."
	}
	return ugly_string
}

func ArrayStringToString(array []string) string {
	if len(array) == 0 {
		return EMPTY_STRING
	}

	var valueToReturn string
	for _, value := range array {
		valueToReturn += "'" + value + "' "
	}

	return valueToReturn
}

func ProgressBar(title string, numberElements int) *progress_bar.ProgressBar {
	bar := progress_bar.NewProgressBar(numberElements).OptionTitle(title + ":").PrintBlankProgressBar()
	return bar
}

/*
longName key: 'Amsterdam 3'
pod from: pods, err := cmd.NetworkManager.GetClosingPods()
*/
func GetPodWithClosedAnnouncement(key string, pods []datatypes.Network_Pod) string {
	for _, pod := range pods {
		if key == *pod.DatacenterLongName {
			return T("closing soon: ") + *pod.Name
		}
	}
	return "-"
}

/*
Converts a data storage value to an appropriate unit.
:param str value: The value to convert.
:param str unit: The unit of the value ('B', 'KB', 'MB', 'GB', 'TB').
:param bool round_result: rounded result
:return: The converted value and its unit.
*/
func ConvertSizes(valueString string, unit string, roundResult bool) string {
	value, _ := strconv.ParseFloat(valueString, 64)
	if value == 0 {
		return "0.00 MB"
	}

	units := []string{"B", "KB", "MB", "GB", "TB"}
	unit = strings.ToUpper(unit)

	unitIndex := -1
	for i, u := range units {
		if u == unit {
			unitIndex = i
			break
		}
	}

	if unitIndex == -1 {
		return "Invalid unit. Must be one of 'B', 'KB', 'MB', 'GB', 'TB'"
	}

	for value > 999 && unitIndex < len(units)-1 {
		value /= 1024
		unitIndex++
	}

	for value < 1 && unitIndex > 0 {
		value *= 1024
		unitIndex--
	}

	if roundResult {
		value = math.Round(value/5) * 5
	}

	return fmt.Sprintf("%.2f %s", value, units[unitIndex])
}

/*
Sums two data storage values.
:param str size1: The first value and its unit.
:param str size2: The second value and its unit.
:return: The sum of the values and its unit.
*/
func SumSizes(size1, size2 string) string {
	if size1 == "0.00 MB" {
		return size2
	}
	if size2 == "0.00 MB" {
		return size1
	}

	value1Str, unit1 := splitSize(size1)
	value2Str, unit2 := splitSize(size2)

	value1, _ := strconv.ParseFloat(value1Str, 64)
	value2, _ := strconv.ParseFloat(value2Str, 64)

	units := []string{"B", "KB", "MB", "GB", "TB"}
	if !isValidUnit(unit1, units) || !isValidUnit(unit2, units) {
		return "Invalid unit in one of the sizes. Unit must be one of 'B', 'KB', 'MB', 'GB', 'TB'"
	}

	value1 *= math.Pow(1024, float64(indexOf(unit1, units)))
	value2 *= math.Pow(1024, float64(indexOf(unit2, units)))

	totalValue := value1 + value2
	totalUnit := "B"
	for totalValue > 999 && totalUnit != "TB" {
		totalValue /= 1024
		totalUnit = units[indexOf(totalUnit, units)+1]
	}

	return fmt.Sprintf("%.2f %s", totalValue, totalUnit)
}

func splitSize(size string) (string, string) {
	parts := strings.Split(size, " ")
	return parts[0], parts[1]
}

func isValidUnit(unit string, units []string) bool {
	for _, u := range units {
		if u == unit {
			return true
		}
	}
	return false
}

func indexOf(element string, data []string) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func FormatStringToTime(timestamp *string) string {
	timeInt, err := strconv.ParseInt(*timestamp, 10, 64)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}
	t := time.Unix(timeInt, 0)
	return t.Format("2006-01-02 15:04:05")

}
