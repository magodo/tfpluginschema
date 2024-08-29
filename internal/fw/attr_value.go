package fw

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

func AttrValueToGo(v attr.Value) (interface{}, error) {
	if v.IsNull() {
		return nil, nil
	}
	if v.IsUnknown() {
		return nil, fmt.Errorf("unexpected unknown value")
	}
	var err error
	switch v := v.(type) {
	case basetypes.BoolValue:
		return v.ValueBool(), nil
	case basetypes.Int32Value:
		return v.ValueInt32(), nil
	case basetypes.Int64Value:
		return v.ValueInt64(), nil
	case basetypes.Float32Value:
		return v.ValueFloat32(), nil
	case basetypes.Float64Value:
		return v.ValueFloat64(), nil
	case basetypes.NumberValue:
		return v.ValueBigFloat(), nil
	case basetypes.StringValue:
		return v.ValueString(), nil
	case basetypes.ListValue:
		l := []interface{}{}
		for _, v := range v.Elements() {
			vv, err := AttrValueToGo(v)
			if err != nil {
				return nil, err
			}
			l = append(l, vv)
		}
		return l, nil
	case basetypes.SetValue:
		l := []interface{}{}
		for _, v := range v.Elements() {
			vv, err := AttrValueToGo(v)
			if err != nil {
				return nil, err
			}
			l = append(l, vv)
		}
		return l, nil
	case basetypes.TupleValue:
		l := []interface{}{}
		for _, v := range v.Elements() {
			vv, err := AttrValueToGo(v)
			if err != nil {
				return nil, err
			}
			l = append(l, vv)
		}
		return l, nil
	case basetypes.ObjectValue:
		m := map[string]interface{}{}
		for k, v := range v.Attributes() {
			m[k], err = AttrValueToGo(v)
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	case basetypes.MapValue:
		m := map[string]interface{}{}
		for k, v := range v.Elements() {
			m[k], err = AttrValueToGo(v)
			if err != nil {
				return nil, err
			}
		}
		return m, nil
	case basetypes.DynamicValue:
		return AttrValueToGo(v.UnderlyingValue())
	default:
		return nil, fmt.Errorf("unhandled type: %T", v)
	}
}
