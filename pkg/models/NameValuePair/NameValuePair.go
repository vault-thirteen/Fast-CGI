package nvpair

import (
	"bytes"
	"fmt"
	"io"

	"github.com/vault-thirteen/Fast-CGI/pkg/models/VariableLength"
	"github.com/vault-thirteen/auxie/reader"
)

// NameValuePair is a name-value pair in accordance to the FastCGI
// Specification Version 1.1.
type NameValuePair struct {
	NameLength  *vl.VariableLength
	ValueLength *vl.VariableLength

	Name  []byte
	Value []byte
}

// NewNameValuePairWithTextValue creates a name-value pair using the specified
// textual arguments.
func NewNameValuePairWithTextValue(name string, value string) (nvp *NameValuePair, err error) {
	nvp = &NameValuePair{
		Name:  []byte(name),
		Value: []byte(value),
	}

	nvp.NameLength, err = vl.NewFromInteger(uint(len(nvp.Name)))
	if err != nil {
		return nil, err
	}

	nvp.ValueLength, err = vl.NewFromInteger(uint(len(nvp.Value)))
	if err != nil {
		return nil, err
	}

	return nvp, nil
}

// NewNameValuePairWithTextValueU is an unsafe version of the
// 'NewNameValuePairWithTextValue' function. It is a helper-function for fast
// creation of a name-value pair using textual arguments. Please note that this
// function returns a NULL pointer on error.
func NewNameValuePairWithTextValueU(name string, value string) (nvp *NameValuePair) {
	var err error
	nvp, err = NewNameValuePairWithTextValue(name, value)
	if err != nil {
		return nil
	}

	return nvp
}

// NewNameValuePairFromStream reads a name-value pair from a byte stream.
func NewNameValuePairFromStream(stream io.Reader) (nvp *NameValuePair, err error) {
	rdr := reader.New(stream)

	nvp = &NameValuePair{}

	nvp.NameLength, err = vl.NewFromStream(stream)
	if err != nil {
		return nil, err
	}

	nvp.ValueLength, err = vl.NewFromStream(stream)
	if err != nil {
		return nil, err
	}

	nvp.Name, err = rdr.ReadBytes(int(nvp.NameLength.RealValue()))
	if err != nil {
		return nil, err
	}

	nvp.Value, err = rdr.ReadBytes(int(nvp.ValueLength.RealValue()))
	if err != nil {
		return nil, err
	}

	return nvp, nil
}

// Measure calculates memory size, or content length, required for storing or
// transmitting of a single name-value pair as a FastCGI data.
func (nvp *NameValuePair) Measure() (n int) {
	return nvp.NameLength.Size() + len(nvp.Name) + nvp.ValueLength.Size() + len(nvp.Value)
}

// MeasureNameValuePairs calculates memory size, or content length, required for
// storing or transmitting the specified name-value pairs.
func MeasureNameValuePairs(params []*NameValuePair) (n int) {
	n = 0
	for _, p := range params {
		n = n + p.Measure()
	}
	return n
}

// ToBytes returns the bytes storing the name-value pair in accordance to the
// FastCGI Specification Version 1.1.
func (nvp *NameValuePair) ToBytes() (ba []byte, err error) {
	var buf bytes.Buffer

	_, err = buf.Write(nvp.NameLength.ToBytes())
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(nvp.ValueLength.ToBytes())
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(nvp.Name)
	if err != nil {
		return nil, err
	}

	_, err = buf.Write(nvp.Value)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// PrintParameters is used for debugging parameters.
func PrintParameters(params []*NameValuePair) {
	for i, p := range params {
		fmt.Println(fmt.Sprintf("%v.\t[%v]=[%v]", i+1, string(p.Name), string(p.Value)))
	}
}
