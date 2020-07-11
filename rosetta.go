package rosetta

// Marshaler interface wraps the MarshalMap function, that allows an object to "marshal" itself into a generic map (map[string]interface{})
type Marshaler interface {

	// MarshalMap function generates a generic map (map[string]interface{}) using the data from that object.
	MarshalMap() map[string]interface{}
}

// Unmarshaler interface wraps the UnmarshalMap function, that allows an object to "unmarshal" (populate) itself from a generic map (map[string]interface{})
type Unmarshaler interface {

	// UnmarshalMap function populates an object from a generic map (map[string]interface{}) or returns an error
	UnmarshalMap(map[string]interface{}) error
}
