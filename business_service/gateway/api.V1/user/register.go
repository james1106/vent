package user

import (
	"github.com/kataras/iris"
	"github.com/golang/protobuf/proto"
	pb "github.com/chenshaobo/vent/business_service/proto"
	"github.com/chenshaobo/vent/business_service/rpclient"
	"github.com/chenshaobo/vent/business_service/utils"
	"golang.org/x/net/context"
	"github.com/chenshaobo/vent/business_service/gateway/api.V1/apiUtils"
	"strconv"
	"github.com/jbrodriguez/mlog"
)


//REGISTER FOR IS API FOR  /api/v1/user/check_code/register
func Register(c *iris.Context) {
	body := c.PostBody()
	c2s := &pb.RegisterC2S{}
	s2c := &pb.RegisterS2C{}

	err := proto.Unmarshal(body, c2s)
	if err != nil {
		s2c.ErrCode = utils.ErrParams
		apiUtils.SetBody(c,s2c)
		return
	}
	conn := rpclient.Get(utils.UserSer)
	if conn == nil {
		s2c.ErrCode = utils.ErrServer
		apiUtils.SetBody(c,s2c)
		return
	}
	rc := pb.NewRegisterClient(conn)
	s2cTmp, err := rc.Register(context.Background(), c2s)
	utils.PrintErr(err)
	if err != nil {
		s2c.ErrCode = utils.ErrServer
		mlog.Error(err)
		apiUtils.SetBody(c,s2c)
		return
	}

	s2c = s2cTmp
	if s2c.ErrCode > 0 {
		apiUtils.SetBody(c,s2c)
		return
	}
	session,errCode := getSession(s2c.UserID)
	if errCode > 0 {
		s2c.ErrCode =  errCode
		apiUtils.SetBody(c,s2c)
		return
	}
	s2c.Session = session
	apiUtils.SetBody(c,s2c)
}

func GetRegisterCaptcha(c * iris.Context){
	phoneNumber := c.Param("phoneNumber")
	c2s := &pb.RegCaptchaC2S{}
	s2c := &pb.RegCaptchaS2C{}
	phoneNumberUint,err := strconv.ParseUint(phoneNumber,10,64)
	if err  != nil{
		s2c.ErrCode = utils.ErrParams
		apiUtils.SetBody(c,s2c)
		return
	}

	c2s.PhoneNumber = phoneNumberUint
	conn := rpclient.Get(utils.UserSer)
	if conn == nil {
		s2c.ErrCode = utils.ErrServer
		apiUtils.SetBody(c,s2c)
		return
	}
	rc := pb.NewRegisterClient(conn)
	s2cTmp ,err  := rc.GetRegCaptcha(context.Background(),c2s)
	if err != nil {
		s2c.ErrCode = utils.ErrServer
		mlog.Error(err)
		apiUtils.SetBody(c,s2c)
		return
	}
	apiUtils.SetBody(c,s2cTmp)
}