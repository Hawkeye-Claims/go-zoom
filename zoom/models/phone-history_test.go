package models_test

import (
	"encoding/json"
	"testing"

	"github.com/Hawkeye-Claims/go-zoom/zoom/models"
)

func TestCallHistoryIdentifiers(t *testing.T) {
	for _, test := range []struct{ body, uuid, path string }{
		{`{"call_path_id":"legacy"}`, "legacy", "legacy"},
		{`{"call_history_uuid":"modern","call_path_id":"legacy"}`, "modern", "legacy"},
		{`{"call_history_uuid":"modern"}`, "modern", ""},
	} {
		var history models.CallHistory
		if err := json.Unmarshal([]byte(test.body), &history); err != nil {
			t.Fatal(err)
		}
		if history.CallHistoryUUID != test.uuid || history.CallPathID != test.path {
			t.Fatalf("%s decoded to %+v", test.body, history)
		}
	}
}

func TestCallElementModernMetadata(t *testing.T) {
	var element models.CallElement
	if err := json.Unmarshal([]byte(`{"call_element_id":"el","call_history_uuid":"history","result":"future_result","result_reason":"future_reason","talk_time":5,"wait_time":7,"duration":12,"voicemail_id":"vm"}`), &element); err != nil {
		t.Fatal(err)
	}
	if element.CallHistoryUUID != "history" || element.TalkTime != 5 || element.WaitTime != 7 || element.Duration != 12 || element.Result != "future_result" || element.ResultReason != "future_reason" || element.VoicemailID != "vm" {
		t.Fatalf("element: %+v", element)
	}
}
