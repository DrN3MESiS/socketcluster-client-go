package tests

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/daominah/socketcluster-client-go/scclient/utils"
	"github.com/daominah/socketcluster-client-go/scclient/models"
)

func TestShouldCheckEqual(t *testing.T) {
	expected := [] byte("mystring")
	assert.True(t, utils.IsEqual("mystring", expected), "String and byte [] should be equal")

}

func _TestShouldSerializeData(t *testing.T) {

	emitEvent := models.EmitEvent{Cid: 2, Data: "My sample data", Event: "chat"}
	expectedData := "{\"event\":\"chat\",\"data\":\"My sample data\",\"cid\":2}"
	// changed to binary msgpack
	assert.Equal(t, expectedData, string(utils.SerializeData(emitEvent)))
}
