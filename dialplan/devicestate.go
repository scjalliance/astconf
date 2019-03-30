package dialplan

// DeviceState returns a device state function.
func DeviceState(device Device) DeviceStateFunc {
	return DeviceStateFunc{Device: device}
}

// DeviceStateFunc is a device state function.
type DeviceStateFunc struct {
	Device Device
}

// Func returns the assembled function call.
func (f DeviceStateFunc) Func() FuncCall {
	return FuncCall{
		Name: "DEVICE_STATE",
		Args: []string{f.Device.String()},
	}
}

// Expr returns the device state function as an expression.
func (f DeviceStateFunc) Expr() ExprDef {
	return ExprDef{Content: f.Func().String()}
}

// Ref returns a reference to the device state function.
func (f DeviceStateFunc) Ref() string {
	return f.Func().String()
}
