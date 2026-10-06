package mapping

import "encoding/json"

// JsonbFromMap converts a map to a JSONB for Postgresql.
func JsonbFromMap(m map[string]any) []byte {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		return nil
	}
	return jsonBytes
}

// MapFromJsonb converts a JSONB to a map.
func MapFromJsonb(jsonBytes []byte) (map[string]any, error) {
	var m map[string]any
	if err := json.Unmarshal(jsonBytes, &m); err != nil {
		return nil, err
	}
	return m, nil
}
