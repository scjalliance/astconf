package dialplan

// Channel returns a dialplan function that reads an item of information
// about the current channel, such as "channeltype".
func Channel(item string) ChannelFunc {
	return ChannelFunc{Item: item}
}

// ChannelFunc is a dialplan function that reads information about the
// current channel.
//
// https://docs.asterisk.org/Asterisk_22_Documentation/API_Documentation/Dialplan_Functions/CHANNEL/
type ChannelFunc struct {
	Item string
}

// Func returns the assembled function call.
func (f ChannelFunc) Func() FuncCall {
	return FuncCall{
		Name: "CHANNEL",
		Args: []string{f.Item},
	}
}

// Expr returns the function as an expression.
func (f ChannelFunc) Expr() ExprDef {
	return ExprDef{Content: f.Func().String()}
}

// QuotedContent returns true, which indicates that the function call should
// be wrapped in double quotes when compared.
func (f ChannelFunc) QuotedContent() bool {
	return true
}

// Validate returns an error if the item contains invalid characters.
func (f ChannelFunc) Validate() error {
	return errorIfAny("channel item", f.Item, invalidArgChars)
}
