package receiver

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"dream/internal/types"

	"github.com/stretchr/testify/assert"
)

type mockProducer struct {
	lastMessage []byte
	shouldFail  bool
}

func (m *mockProducer) SendMessage(message []byte) error {
	if m.shouldFail {
		return assert.AnError
	}
	m.lastMessage = message
	return nil
}

func TestFileReceiver_HandleUpload(t *testing.T) {
	validator := &MessageRequestValidator{}
	producer := &mockProducer{}
	receiver := NewFileReceiver(producer, validator)

	validReq := types.MessageRequest{
		MachineID:   "123",
		MachineName: "TestMachine",
		OSVersion:   "Ubuntu",
		CommandType: "ps",
		UserName:    "student",
		UserID:      "u1",
		Faculty:     "Law",
	}
	body, _ := json.Marshal(validReq)
	r := httptest.NewRequest(http.MethodPost, "/upload", bytes.NewReader(body))
	w := httptest.NewRecorder()
	receiver.HandleUpload(w, r)
	resp := w.Result()
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(respBody), "success")
	assert.NotNil(t, producer.lastMessage)

	// Invalid request (missing MachineID)
	invalidReq := validReq
	invalidReq.MachineID = ""
	body, _ = json.Marshal(invalidReq)
	r = httptest.NewRequest(http.MethodPost, "/upload", bytes.NewReader(body))
	w = httptest.NewRecorder()
	receiver.HandleUpload(w, r)
	resp = w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Invalid method
	r = httptest.NewRequest(http.MethodGet, "/upload", nil)
	w = httptest.NewRecorder()
	receiver.HandleUpload(w, r)
	resp = w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

	// Producer failure
	producer.shouldFail = true
	body, _ = json.Marshal(validReq)
	r = httptest.NewRequest(http.MethodPost, "/upload", bytes.NewReader(body))
	w = httptest.NewRecorder()
	receiver.HandleUpload(w, r)
	resp = w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
