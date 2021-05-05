package utils

import (
	"encoding/json"
	"errors"
	"fmt"
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
	"github.com/urfave/cli"
	. "github.ibm.com/SoftLayer/softlayer-cli/plugin/i18n"
)

//TODO support resolving guid to integer id
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
		return "", err
	}
	if image.GlobalIdentifier == nil {
		return "", errors.New(T("Image global identifier not found"))
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
		return ""
	}
	return strings.Trim(strings.Replace(fmt.Sprint(slice), " ", ",", -1), "[]")
}

func IntSliceToString(slice []int) string {
	if len(slice) == 0 {
		return ""
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

//the first bool value indicates whether the source slice are all in target slice
//the second int value indicates the index of the source slice value is not in the target slice
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

//Converts number of bytes to a string in gigabytes.
func B2GB(bytes int) string {
	return fmt.Sprintf("%.2f%s", float32(bytes)/1024/1024/1024, "G")
}

func FormatBoolPointer(value *bool) string {
	if value == nil {
		return "-"
	}
	return strconv.FormatBool(sl.Get(value).(bool))
}

func FormatStringPointer(value *string) string {
	if value == nil {
		return "-"
	}
	return sl.Get(value).(string)
}

func FormatStringPointerName(value *string) string {
	if value == nil {
		return ""
	}
	return sl.Get(value).(string)
}

func FormatIntPointer(value *int) string {
	if value == nil {
		return "-"
	}
	return strconv.Itoa(sl.Get(value).(int))
}

func FormatIntPointerName(value *int) string {
	if value == nil {
		return ""
	}
	return strconv.Itoa(sl.Get(value).(int))
}

func FormatUIntPointer(value *uint) string {
	if value == nil {
		return "-"
	}
	return fmt.Sprintf("%d", sl.Get(value).(uint))
}

func FormatSLFloatPointerToInt(value *datatypes.Float64) string {
	if value == nil {
		return "-"
	}
	return fmt.Sprintf("%d", int(sl.Get(value).(datatypes.Float64)))
}

func FormatSLFloatPointerToFloat(value *datatypes.Float64) string {
	if value == nil {
		return "-"
	}
	return fmt.Sprintf("%f", float64(sl.Get(value).(datatypes.Float64)))
}

func FormatSLTimePointer(value *datatypes.Time) string {
	if value == nil {
		return "-"
	}
	return value.UTC().Format(time.RFC3339)
}

func Bool2Int(value bool) int {
	if value {
		return 1
	}
	return 0
}

func ValidateColumns(sortby string, columns []string, defaultColumns []string, optionalColumns, sortColumns []string, context *cli.Context) ([]string, error) {
	if sortby != "" && StringInSlice(sortby, sortColumns) == -1 {
		return nil, bmxErr.NewInvalidUsageError(T("--sortby {{.Column}} is not supported.", map[string]interface{}{"Column": sortby}))
	}
	allColumns := append(defaultColumns, optionalColumns...)
	if exist, index := SliceInSlice(columns, allColumns); len(columns) > 0 && exist == false {
		if context.IsSet("columns") {
			return nil, bmxErr.NewInvalidUsageError(T("--columns {{.Column}} is not supported.", map[string]interface{}{"Column": columns[index]}))
		}
		return nil, bmxErr.NewInvalidUsageError(T("--column {{.Column}} is not supported.", map[string]interface{}{"Column": columns[index]}))
	}

	if len(columns) == 0 {
		return defaultColumns, nil
	}
	return columns, nil
}

func GetMask(maskMap map[string]string, columns []string, sortBy string) string {

	if sortBy != "" && StringInSlice(sortBy, columns) == -1 {
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
	return cli.NewExitError("", 1)
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
		return ""
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
