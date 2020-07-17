package handler_test

import (
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/assert"
	"github.com/whale-team/whaleEcho/pkg/echoproto"
	"google.golang.org/protobuf/proto"
)

var _ = Describe("Room Handler", func() {
	conn, err := suite.Dial()
	assert.Nil(suite.T, err)

	Describe("#CreateRoom", func() {
		command := echoproto.Command{}
		command.Type = echoproto.CommandType_CreateRoom

		Context("when ...", func() {
			It("should response ok", func() {
				room := echoproto.Room{Id: 123}
				roomData, _ := proto.Marshal(&room)
				command.Type = echoproto.CommandType_CreateRoom
				command.Payload = roomData

				err := suite.SendCommand(conn, &command)
				assert.Nil(suite.T, err)

				resp, err := suite.ReadResp(conn)
				assert.Nil(suite.T, err)

				assert.Equal(suite.T, echoproto.RespStatus_OK, resp.Status)
			})
		})
	})
})
