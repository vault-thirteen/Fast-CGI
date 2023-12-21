package nvpair

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/vault-thirteen/auxie/reader"
	"github.com/vault-thirteen/auxie/tester"
)

func Test_NewNameValuePairWithTextValue(t *testing.T) {
	aTest := tester.New(t)

	type InputData struct {
		Name  string
		Value string
	}

	type OutputData struct {
		Name  []byte
		Value []byte

		NameLength_RealValue    uint32
		NameLength_RawValueSize int // Number of bytes to store the raw data.

		ValueLength_RealValue    uint32
		ValueLength_RawValueSize int // Number of bytes to store the raw data.
	}

	type TestData struct {
		Data            InputData
		IsErrorExpected bool
		ExpectedValue   OutputData
	}

	tests := []TestData{
		{
			Data: InputData{
				Name:  "",
				Value: "",
			},
			IsErrorExpected: false,
			ExpectedValue: OutputData{
				Name:                     []byte{},
				Value:                    []byte{},
				NameLength_RealValue:     0,
				NameLength_RawValueSize:  1,
				ValueLength_RealValue:    0,
				ValueLength_RawValueSize: 1,
			},
		},
		{
			Data: InputData{
				Name:  "The Planet Earth",
				Value: "The Vault Thirteen",
			},
			IsErrorExpected: false,
			ExpectedValue: OutputData{
				Name:                     []byte("The Planet Earth"),
				Value:                    []byte("The Vault Thirteen"),
				NameLength_RealValue:     16,
				NameLength_RawValueSize:  1,
				ValueLength_RealValue:    18,
				ValueLength_RawValueSize: 1,
			},
		},
		{
			Data: InputData{
				Name:  "The name of the range hails from the Sanskrit Himālaya (हिमालय 'abode of the snow'), from himá (हिम 'snow') and ā-laya (आलय 'home, dwelling').",
				Value: "The Himalayas, or Himalaya is a mountain range in Asia, separating the plains of the Indian subcontinent from the Tibetan Plateau.",
			},
			IsErrorExpected: false,
			ExpectedValue: OutputData{
				Name:                     []byte("The name of the range hails from the Sanskrit Himālaya (हिमालय 'abode of the snow'), from himá (हिम 'snow') and ā-laya (आलय 'home, dwelling')."),
				Value:                    []byte("The Himalayas, or Himalaya is a mountain range in Asia, separating the plains of the Indian subcontinent from the Tibetan Plateau."),
				NameLength_RealValue:     uint32(len("The name of the range hails from the Sanskrit Himālaya (हिमालय 'abode of the snow'), from himá (हिम 'snow') and ā-laya (आलय 'home, dwelling').")),
				NameLength_RawValueSize:  4,
				ValueLength_RealValue:    130,
				ValueLength_RawValueSize: 4,
			},
		},
	}

	var result *NameValuePair
	var err error

	for i, test := range tests {
		fmt.Printf("[%v]", i+1)

		result, err = NewNameValuePairWithTextValue(test.Data.Name, test.Data.Value)

		if test.IsErrorExpected {
			aTest.MustBeAnError(err)
		} else {
			aTest.MustBeNoError(err)
		}

		aTest.MustBeEqual(result.Name, test.ExpectedValue.Name)
		aTest.MustBeEqual(result.Value, test.ExpectedValue.Value)
		aTest.MustBeEqual(result.NameLength.RealValue(), test.ExpectedValue.NameLength_RealValue)
		aTest.MustBeEqual(result.NameLength.Size(), test.ExpectedValue.NameLength_RawValueSize)
		aTest.MustBeEqual(result.ValueLength.RealValue(), test.ExpectedValue.ValueLength_RealValue)
		aTest.MustBeEqual(result.ValueLength.Size(), test.ExpectedValue.ValueLength_RawValueSize)
	}
	fmt.Println()
}

func Test_NewNameValuePairWithTextValueU(t *testing.T) {
	aTest := tester.New(t)

	type InputData struct {
		Name  string
		Value string
	}

	type OutputData struct {
		Name  []byte
		Value []byte

		NameLength_RealValue    uint32
		NameLength_RawValueSize int // Number of bytes to store the raw data.

		ValueLength_RealValue    uint32
		ValueLength_RawValueSize int // Number of bytes to store the raw data.
	}

	type TestData struct {
		Data            InputData
		IsErrorExpected bool
		ExpectedValue   OutputData
	}

	tests := []TestData{
		{
			Data: InputData{
				Name:  "",
				Value: "",
			},
			IsErrorExpected: false,
			ExpectedValue: OutputData{
				Name:                     []byte{},
				Value:                    []byte{},
				NameLength_RealValue:     0,
				NameLength_RawValueSize:  1,
				ValueLength_RealValue:    0,
				ValueLength_RawValueSize: 1,
			},
		},
		{
			Data: InputData{
				Name:  "The Planet Earth",
				Value: "The Vault Thirteen",
			},
			IsErrorExpected: false,
			ExpectedValue: OutputData{
				Name:                     []byte("The Planet Earth"),
				Value:                    []byte("The Vault Thirteen"),
				NameLength_RealValue:     16,
				NameLength_RawValueSize:  1,
				ValueLength_RealValue:    18,
				ValueLength_RawValueSize: 1,
			},
		},
		{
			Data: InputData{
				Name:  "The name of the range hails from the Sanskrit Himālaya (हिमालय 'abode of the snow'), from himá (हिम 'snow') and ā-laya (आलय 'home, dwelling').",
				Value: "The Himalayas, or Himalaya is a mountain range in Asia, separating the plains of the Indian subcontinent from the Tibetan Plateau.",
			},
			IsErrorExpected: false,
			ExpectedValue: OutputData{
				Name:                     []byte("The name of the range hails from the Sanskrit Himālaya (हिमालय 'abode of the snow'), from himá (हिम 'snow') and ā-laya (आलय 'home, dwelling')."),
				Value:                    []byte("The Himalayas, or Himalaya is a mountain range in Asia, separating the plains of the Indian subcontinent from the Tibetan Plateau."),
				NameLength_RealValue:     uint32(len("The name of the range hails from the Sanskrit Himālaya (हिमालय 'abode of the snow'), from himá (हिम 'snow') and ā-laya (आलय 'home, dwelling').")),
				NameLength_RawValueSize:  4,
				ValueLength_RealValue:    130,
				ValueLength_RawValueSize: 4,
			},
		},
	}

	var result *NameValuePair

	for i, test := range tests {
		fmt.Printf("[%v]", i+1)

		result = NewNameValuePairWithTextValueU(test.Data.Name, test.Data.Value)

		aTest.MustBeEqual(result.Name, test.ExpectedValue.Name)
		aTest.MustBeEqual(result.Value, test.ExpectedValue.Value)
		aTest.MustBeEqual(result.NameLength.RealValue(), test.ExpectedValue.NameLength_RealValue)
		aTest.MustBeEqual(result.NameLength.Size(), test.ExpectedValue.NameLength_RawValueSize)
		aTest.MustBeEqual(result.ValueLength.RealValue(), test.ExpectedValue.ValueLength_RealValue)
		aTest.MustBeEqual(result.ValueLength.Size(), test.ExpectedValue.ValueLength_RawValueSize)
	}
	fmt.Println()
}

func Test_NewNameValuePairFromStream(t *testing.T) {
	aTest := tester.New(t)

	type InputData = []byte

	type OutputData struct {
		Name  []byte
		Value []byte

		NameLength_RealValue    uint32
		NameLength_RawValueSize int // Number of bytes to store the raw data.

		ValueLength_RealValue    uint32
		ValueLength_RawValueSize int // Number of bytes to store the raw data.
	}

	type TestData struct {
		Data            InputData
		IsErrorExpected bool
		ExpectedValue   OutputData
	}

	tests := []TestData{
		{
			Data:            []byte{0, 0},
			IsErrorExpected: false,
			ExpectedValue: OutputData{
				Name:                     []byte{},
				Value:                    []byte{},
				NameLength_RealValue:     0,
				NameLength_RawValueSize:  1,
				ValueLength_RealValue:    0,
				ValueLength_RawValueSize: 1,
			},
		},
		{
			Data: append(append(append(append(make([]byte, 0),
				[]byte{byte(16)}...),
				[]byte{byte(18)}...),
				[]byte("The Planet Earth")...),
				[]byte("The Vault Thirteen")...),
			IsErrorExpected: false,
			ExpectedValue: OutputData{
				Name:                     []byte("The Planet Earth"),
				Value:                    []byte("The Vault Thirteen"),
				NameLength_RealValue:     16,
				NameLength_RawValueSize:  1,
				ValueLength_RealValue:    18,
				ValueLength_RawValueSize: 1,
			},
		},
		{
			Data: append(append(append(append(make([]byte, 0),
				[]byte{128, 0, 0, 169}...),
				[]byte{128, 0, 0, 130}...),
				[]byte("The name of the range hails from the Sanskrit Himālaya (हिमालय 'abode of the snow'), from himá (हिम 'snow') and ā-laya (आलय 'home, dwelling').")...),
				[]byte("The Himalayas, or Himalaya is a mountain range in Asia, separating the plains of the Indian subcontinent from the Tibetan Plateau.")...),
			IsErrorExpected: false,
			ExpectedValue: OutputData{
				Name:                     []byte("The name of the range hails from the Sanskrit Himālaya (हिमालय 'abode of the snow'), from himá (हिम 'snow') and ā-laya (आलय 'home, dwelling')."),
				Value:                    []byte("The Himalayas, or Himalaya is a mountain range in Asia, separating the plains of the Indian subcontinent from the Tibetan Plateau."),
				NameLength_RealValue:     169,
				NameLength_RawValueSize:  4,
				ValueLength_RealValue:    130,
				ValueLength_RawValueSize: 4,
			},
		},
	}

	var rdr io.Reader
	var result *NameValuePair
	var err error

	for i, test := range tests {
		fmt.Printf("[%v]", i+1)

		rdr = reader.New(bytes.NewReader(test.Data))
		result, err = NewNameValuePairFromStream(rdr)

		if test.IsErrorExpected {
			aTest.MustBeAnError(err)
		} else {
			aTest.MustBeNoError(err)
		}

		aTest.MustBeEqual(result.Name, test.ExpectedValue.Name)
		aTest.MustBeEqual(result.Value, test.ExpectedValue.Value)
		aTest.MustBeEqual(result.NameLength.RealValue(), test.ExpectedValue.NameLength_RealValue)
		aTest.MustBeEqual(result.NameLength.Size(), test.ExpectedValue.NameLength_RawValueSize)
		aTest.MustBeEqual(result.ValueLength.RealValue(), test.ExpectedValue.ValueLength_RealValue)
		aTest.MustBeEqual(result.ValueLength.Size(), test.ExpectedValue.ValueLength_RawValueSize)
	}
	fmt.Println()
}

func Test_Measure(t *testing.T) {
	aTest := tester.New(t)

	type InputData = []byte

	type OutputData = int

	type TestData struct {
		Data            InputData
		IsErrorExpected bool
		ExpectedValue   OutputData
	}

	tests := []TestData{
		{
			Data:            []byte{0, 0},
			IsErrorExpected: false,
			ExpectedValue:   1 + 1 + 0 + 0,
		},
		{
			Data: append(append(append(append(make([]byte, 0),
				[]byte{byte(16)}...),
				[]byte{byte(18)}...),
				[]byte("The Planet Earth")...),
				[]byte("The Vault Thirteen")...),
			IsErrorExpected: false,
			ExpectedValue:   1 + 1 + 16 + 18,
		},
		{
			Data: append(append(append(append(make([]byte, 0),
				[]byte{128, 0, 0, 169}...),
				[]byte{128, 0, 0, 130}...),
				[]byte("The name of the range hails from the Sanskrit Himālaya (हिमालय 'abode of the snow'), from himá (हिम 'snow') and ā-laya (आलय 'home, dwelling').")...),
				[]byte("The Himalayas, or Himalaya is a mountain range in Asia, separating the plains of the Indian subcontinent from the Tibetan Plateau.")...),
			IsErrorExpected: false,
			ExpectedValue:   4 + 4 + 169 + 130,
		},
	}

	var rdr io.Reader
	var nvp *NameValuePair
	var err error
	var result int

	for i, test := range tests {
		fmt.Printf("[%v]", i+1)

		rdr = reader.New(bytes.NewReader(test.Data))
		nvp, err = NewNameValuePairFromStream(rdr)
		aTest.MustBeNoError(err)
		result = nvp.Measure()

		aTest.MustBeEqual(result, test.ExpectedValue)
	}
	fmt.Println()
}

func Test_MeasureNameValuePairs(t *testing.T) {
	aTest := tester.New(t)
	nvp1 := NewNameValuePairWithTextValueU("T1", "V1")
	nvp2 := NewNameValuePairWithTextValueU("T22", "V22")

	aTest.MustBeEqual(MeasureNameValuePairs([]*NameValuePair{}), 0)
	aTest.MustBeEqual(MeasureNameValuePairs([]*NameValuePair{nvp1}), 6)
	aTest.MustBeEqual(MeasureNameValuePairs([]*NameValuePair{nvp1, nvp2}), 6+8)
	aTest.MustBeEqual(MeasureNameValuePairs([]*NameValuePair{nvp2, nvp1}), 8+6)
}

func Test_ToBytes(t *testing.T) {
	aTest := tester.New(t)
	nvp := NewNameValuePairWithTextValueU("A", "BCD")
	var err error
	var ba []byte

	ba, err = nvp.ToBytes()
	aTest.MustBeNoError(err)
	aTest.MustBeEqual(ba, []byte{1, 3, 'A', 'B', 'C', 'D'})
}
