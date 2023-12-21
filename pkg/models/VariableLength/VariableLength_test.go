package vl

import (
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/vault-thirteen/auxie/tester"
)

func Test_NewFromStream(t *testing.T) {
	aTest := tester.New(t)

	type TestData struct {
		Data            []byte
		IsErrorExpected bool
		ExpectedValue   *VariableLength
	}

	tests := []TestData{
		{
			Data:            []byte{0},
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    0,
				rawValueSize: 1,
				rawValue1B:   0,
				rawValue4B:   0,
			},
		},
		{
			Data:            []byte{1},
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    1,
				rawValueSize: 1,
				rawValue1B:   1,
				rawValue4B:   0,
			},
		},
		{
			Data:            []byte{255 >> 1},
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    127,
				rawValueSize: 1,
				rawValue1B:   127,
				rawValue4B:   0,
			},
		},
		{
			Data:            []byte{128, 0, 0, 128},
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    128,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   2_147_483_776,
			},
		},
		{
			Data:            []byte{128, 0, 0, 129},
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    129,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   2_147_483_777,
			},
		},
		{
			Data:            []byte{128, 0, 0, 255},
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    255,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   2_147_483_903,
			},
		},
		{
			Data:            []byte{128, 0, 1, 0},
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    256,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   2_147_483_904,
			},
		},
		{
			Data:            []byte{128, 1, 0, 0},
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    65536,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   2_147_549_184,
			},
		},
		{
			Data:            []byte{129, 0, 0, 0},
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    16_777_216,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   2_164_260_864,
			},
		},
		{
			Data:            []byte{255, 0, 0, 0},
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    2_130_706_432,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   4_278_190_080,
			},
		},
		{
			Data:            []byte{255, 255, 255, 255},
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    MaxVariableLength,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   4_294_967_295,
			},
		},
	}

	var rdr io.Reader
	var result *VariableLength
	var err error

	for i, test := range tests {
		fmt.Printf("[%v]", i+1)

		rdr = bytes.NewReader(test.Data)
		result, err = NewFromStream(rdr)

		if test.IsErrorExpected {
			aTest.MustBeAnError(err)
		} else {
			aTest.MustBeNoError(err)
		}

		aTest.MustBeEqual(result, test.ExpectedValue)
	}
	fmt.Println()
}

func Test_NewFromInteger(t *testing.T) {
	aTest := tester.New(t)

	type TestData struct {
		Data            uint
		IsErrorExpected bool
		ExpectedValue   *VariableLength
	}

	tests := []TestData{
		{
			Data:            uint(MaxVariableLength + 1),
			IsErrorExpected: true,
			ExpectedValue:   (*VariableLength)(nil),
		},
		{
			Data:            0,
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    0,
				rawValueSize: 1,
				rawValue1B:   0,
				rawValue4B:   0,
			},
		},
		{
			Data:            1,
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    1,
				rawValueSize: 1,
				rawValue1B:   1,
				rawValue4B:   0,
			},
		},
		{
			Data:            255 >> 1,
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    127,
				rawValueSize: 1,
				rawValue1B:   127,
				rawValue4B:   0,
			},
		},
		{
			Data:            128,
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    128,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   2_147_483_776,
			},
		},
		{
			Data:            129,
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    129,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   2_147_483_777,
			},
		},
		{
			Data:            255,
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    255,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   2_147_483_903,
			},
		},
		{
			Data:            256,
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    256,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   2_147_483_904,
			},
		},
		{
			Data:            65536,
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    65536,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   2_147_549_184,
			},
		},
		{
			Data:            16_777_216,
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    16_777_216,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   2_164_260_864,
			},
		},
		{
			Data:            0x7F * 256 * 256 * 256,
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    2_130_706_432,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   4_278_190_080,
			},
		},
		{
			Data:            uint(MaxVariableLength),
			IsErrorExpected: false,
			ExpectedValue: &VariableLength{
				realValue:    MaxVariableLength,
				rawValueSize: 4,
				rawValue1B:   0,
				rawValue4B:   4_294_967_295,
			},
		},
	}

	var result *VariableLength
	var err error

	for i, test := range tests {
		fmt.Printf("[%v]", i+1)

		result, err = NewFromInteger(test.Data)

		if test.IsErrorExpected {
			aTest.MustBeAnError(err)
		} else {
			aTest.MustBeNoError(err)
		}

		aTest.MustBeEqual(result, test.ExpectedValue)
	}
	fmt.Println()
}

func Test_ToBytes(t *testing.T) {
	aTest := tester.New(t)

	type TestData struct {
		Data            []byte
		IsErrorExpected bool
		ExpectedValue   []byte
	}

	tests := []TestData{
		{
			Data:            []byte{0},
			IsErrorExpected: false,
			ExpectedValue:   []byte{0},
		},
		{
			Data:            []byte{1},
			IsErrorExpected: false,
			ExpectedValue:   []byte{1},
		},
		{
			Data:            []byte{255 >> 1},
			IsErrorExpected: false,
			ExpectedValue:   []byte{127},
		},
		{
			Data:            []byte{128, 0, 0, 128},
			IsErrorExpected: false,
			ExpectedValue:   []byte{128, 0, 0, 128},
		},
		{
			Data:            []byte{128, 0, 0, 129},
			IsErrorExpected: false,
			ExpectedValue:   []byte{128, 0, 0, 129},
		},
		{
			Data:            []byte{255, 255, 255, 255},
			IsErrorExpected: false,
			ExpectedValue:   []byte{255, 255, 255, 255},
		},
	}

	var rdr io.Reader
	var vl *VariableLength
	var result []byte
	var err error

	for i, test := range tests {
		fmt.Printf("[%v]", i+1)

		rdr = bytes.NewReader(test.Data)
		vl, err = NewFromStream(rdr)
		aTest.MustBeNoError(err)
		result = vl.ToBytes()

		aTest.MustBeEqual(result, test.ExpectedValue)
	}
	fmt.Println()
}

func Test_RealValue(t *testing.T) {
	aTest := tester.New(t)

	type TestData struct {
		Data            []byte
		IsErrorExpected bool
		ExpectedValue   uint32
	}

	tests := []TestData{
		{
			Data:            []byte{0},
			IsErrorExpected: false,
			ExpectedValue:   0,
		},
		{
			Data:            []byte{1},
			IsErrorExpected: false,
			ExpectedValue:   1,
		},
		{
			Data:            []byte{255 >> 1},
			IsErrorExpected: false,
			ExpectedValue:   127,
		},
		{
			Data:            []byte{128, 0, 0, 128},
			IsErrorExpected: false,
			ExpectedValue:   128,
		},
		{
			Data:            []byte{128, 0, 0, 129},
			IsErrorExpected: false,
			ExpectedValue:   129,
		},
		{
			Data:            []byte{255, 255, 255, 255},
			IsErrorExpected: false,
			ExpectedValue:   MaxVariableLength,
		},
	}

	var rdr io.Reader
	var vl *VariableLength
	var result uint32
	var err error

	for i, test := range tests {
		fmt.Printf("[%v]", i+1)

		rdr = bytes.NewReader(test.Data)
		vl, err = NewFromStream(rdr)
		aTest.MustBeNoError(err)
		result = vl.RealValue()

		aTest.MustBeEqual(result, test.ExpectedValue)
	}
	fmt.Println()
}

func Test_Size(t *testing.T) {
	aTest := tester.New(t)

	type TestData struct {
		Data            []byte
		IsErrorExpected bool
		ExpectedValue   int
	}

	tests := []TestData{
		{
			Data:            []byte{0},
			IsErrorExpected: false,
			ExpectedValue:   1,
		},
		{
			Data:            []byte{1},
			IsErrorExpected: false,
			ExpectedValue:   1,
		},
		{
			Data:            []byte{255 >> 1},
			IsErrorExpected: false,
			ExpectedValue:   1,
		},
		{
			Data:            []byte{128, 0, 0, 128},
			IsErrorExpected: false,
			ExpectedValue:   4,
		},
		{
			Data:            []byte{128, 0, 0, 129},
			IsErrorExpected: false,
			ExpectedValue:   4,
		},
		{
			Data:            []byte{255, 255, 255, 255},
			IsErrorExpected: false,
			ExpectedValue:   4,
		},
	}

	var rdr io.Reader
	var vl *VariableLength
	var result int
	var err error

	for i, test := range tests {
		fmt.Printf("[%v]", i+1)

		rdr = bytes.NewReader(test.Data)
		vl, err = NewFromStream(rdr)
		aTest.MustBeNoError(err)
		result = vl.Size()

		aTest.MustBeEqual(result, test.ExpectedValue)
	}
	fmt.Println()
}
