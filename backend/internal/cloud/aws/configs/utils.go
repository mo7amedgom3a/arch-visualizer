package configs
import "encoding/json"
// contains global configurations for AWS resources
type Tag struct {
	Key string `json:"key"`
	Value string `json:"value"`
}


// marchal and unmarchal the tag to and from json
func (t Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]string{
		"key": t.Key,
		"value": t.Value,
	})
}

func (t *Tag) UnmarshalJSON(data []byte) error {
	var v map[string]string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	t.Key = v["key"]
	t.Value = v["value"]
	return nil
}