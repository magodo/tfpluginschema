package fw_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/magodo/tfpluginschema/internal/fw"
	"github.com/stretchr/testify/require"
)

func TestAttrValueToGo(t *testing.T) {
	cases := []struct {
		name   string
		input  attr.Value
		expect interface{}
		err    error
	}{
		{
			name:   "bool",
			input:  basetypes.NewBoolValue(true),
			expect: true,
		},
		{
			name:   "int32",
			input:  basetypes.NewInt32Value(123),
			expect: int32(123),
		},
		{
			name:   "int64",
			input:  basetypes.NewInt64Value(123),
			expect: int64(123),
		},
		{
			name:   "float32",
			input:  basetypes.NewFloat32Value(123),
			expect: float32(123),
		},
		{
			name:   "float64",
			input:  basetypes.NewFloat64Value(123),
			expect: float64(123),
		},
		{
			name:   "number",
			input:  basetypes.NewNumberValue(big.NewFloat(123)),
			expect: big.NewFloat(123),
		},
		{
			name:   "string",
			input:  basetypes.NewStringValue("abc"),
			expect: "abc",
		},
		{
			name:   "list",
			input:  basetypes.NewListValueMust(basetypes.BoolType{}, []attr.Value{basetypes.NewBoolValue(true)}),
			expect: []interface{}{true},
		},
		{
			name:   "set",
			input:  basetypes.NewSetValueMust(basetypes.BoolType{}, []attr.Value{basetypes.NewBoolValue(true)}),
			expect: []interface{}{true},
		},
		{
			name:   "tuple",
			input:  basetypes.NewTupleValueMust([]attr.Type{basetypes.BoolType{}}, []attr.Value{basetypes.NewBoolValue(true)}),
			expect: []interface{}{true},
		},
		{
			name:   "map",
			input:  basetypes.NewMapValueMust(basetypes.BoolType{}, map[string]attr.Value{"a": basetypes.NewBoolValue(true)}),
			expect: map[string]interface{}{"a": true},
		},
		{
			name: "object",
			input: basetypes.NewObjectValueMust(
				map[string]attr.Type{"a": basetypes.BoolType{}},
				map[string]attr.Value{"a": basetypes.NewBoolValue(true)},
			),
			expect: map[string]interface{}{"a": true},
		},
		{
			name:   "dynamic",
			input:  basetypes.NewDynamicValue(basetypes.NewBoolValue(true)),
			expect: true,
		},
		{
			name: "list of object",
			input: basetypes.NewListValueMust(
				basetypes.ObjectType{AttrTypes: map[string]attr.Type{"a": basetypes.BoolType{}}},
				[]attr.Value{basetypes.NewObjectValueMust(
					map[string]attr.Type{"a": basetypes.BoolType{}},
					map[string]attr.Value{"a": basetypes.NewBoolValue(true)},
				)},
			),
			expect: []interface{}{map[string]interface{}{"a": true}},
		},
		{
			name:   "null value",
			input:  basetypes.NewBoolNull(),
			expect: nil,
		},
		{
			name:  "unknown value",
			input: basetypes.NewBoolUnknown(),
			err:   errors.New("unexpected unknown value"),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			got, err := fw.AttrValueToGo(tt.input)
			if tt.err != nil {
				require.Error(t, tt.err, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.expect, got)
		})
	}
}
