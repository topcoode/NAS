package nas

import (
	"errors"
	"fmt"

	"free5gc/context"
	"free5gc/fsm"
	"free5gc/gmm"
	"free5gc/logger"
	"free5gc/openapi/models"
	"free5gc/util"

	"github.com/free5gc/nas"
)

var (
	fsmN2Queue *util.FSMQueue
)

/*
the N2MessageHandler function is responsible for handling
NAS (Network Access Stratum)
messages received on the N2 interface.
*/

/*
The value parameter represents the NAS message received on the N2 interface.
It is expected to be of type fsm.NasMessage.

The index parameter is an integer that identifies the specific
FSM (Finite State Machine)
instance associated with the NAS message. This index is used to

	route the message to the appropriate FSM.
*/
func N2MessageHandler(value interface{}, index int) {
	// fmt.Println("N2MessageHandler ", value);
	// fmt.Println("")

	var nasMsg = value.(fsm.NasMessage)
	gmm.GmmFSM.SendEvent(nasMsg.Nstate, nasMsg.Nevent, nasMsg.Nargs)
}

func SendToSendEvent2(msg fsm.NasMessage) {

	amfUe := msg.Nargs[gmm.ArgAmfUe].(*context.AmfUe)
	accessType := msg.Nargs[gmm.ArgAccessType].(models.AccessType)

	// fmt.Println("SendToSendEvent2 ", amfUe, accessType, msg);
	// fmt.Println("")
	// fmt.Println("- ", amfUe.RanUe[accessType] , "-")
	// fmt.Println("")

	if amfUe.RanUe[accessType] != nil {
		//fmt.Println(" if  -- ", fsmN2Queue.FSMIndex2(amfUe.RanUe[accessType].AmfUeNgapId), " -- ", amfUe.RanUe[accessType].AmfUeNgapId)

		fsmN2Queue.EnqueueFP(msg, N2MessageHandler, fsmN2Queue.FSMIndex2(amfUe.RanUe[accessType].AmfUeNgapId))
		//fmt.Println(" if end")
	} else {
		fsmN2Queue.EnqueueFP(msg, N2MessageHandler, 1)
		//fmt.Println(" else")
	}

	//fmt.Println("")
}

func InitN2Queue() {
	logger.NgapLog.Errorf("Initialized N2Queue with 100 Message Handlers and 1000 Queue")
	gmm.GmmFSM.SetFSMObj(SendToSendEvent2)
	fsmN2Queue = util.NewFSMQueue(N2MessageHandler, 300)
	fsmN2Queue.SetMaxQueue(1000)

}

func Perlog() {
	if fsmN2Queue != nil {
		logger.Perlog.Infof("NAS  Queue     Queued  %-10d     Dequeued  %-10d   Pending  %-10d   Threads  %-3d  Rejected %-10d", fsmN2Queue.Total, fsmN2Queue.Executed, fsmN2Queue.Total-fsmN2Queue.Executed, fsmN2Queue.Count, fsmN2Queue.Rejected)
	}
}

func NASPendingQueueCount() int64 {
	if fsmN2Queue != nil {
		return (fsmN2Queue.Total - fsmN2Queue.Executed)
	}
	return 0
}

func NASThreadCount() int64 {
	if fsmN2Queue != nil {
		return int64(fsmN2Queue.Count)
	}
	return 0
}

func NASRejectedCount() {
	if fsmN2Queue != nil {
		fsmN2Queue.Rejected++
	}
}

func Dispatch(ue *context.AmfUe, accessType models.AccessType, procedureCode int64, msg *nas.Message, RanUeNgapId int64) error {

	if msg.GmmMessage == nil {
		return errors.New("Gmm Message is nil")
	}

	if msg.GsmMessage != nil {
		return errors.New("GSM Message should include in GMM Message")
	}

	if ue.State[accessType] == nil {
		return fmt.Errorf("UE State is empty (accessType=%q). Can't send GSM Message", accessType)
	}

	/*
		nasMsg := fsm.NasMessage{ ue.State[accessType], gmm.GmmMessageEvent, fsm.ArgsType{
			gmm.ArgAmfUe:         ue,
			gmm.ArgAccessType:    accessType,
			gmm.ArgNASMessage:    msg.GmmMessage,
			gmm.ArgProcedureCode: procedureCode,
		}}
	*/
	//fsmN2Queue.EnqueueFP( nasMsg, N2MessageHandler, fsmN2Queue.FSMIndex2( RanUeNgapId) )
	//return nil

	context.AMF_Self().IncrementNASMessage()

	if NASPendingQueueCount() > NASThreadCount() {
		NASRejectedCount()
		return fmt.Errorf("Rejecting Message, NASPending Queue Count is greater than thread count")
	}

	return gmm.GmmFSM.SendEvent3(ue.State[accessType], gmm.GmmMessageEvent, fsm.ArgsType{
		gmm.ArgAmfUe:         ue,
		gmm.ArgAccessType:    accessType,
		gmm.ArgNASMessage:    msg.GmmMessage,
		gmm.ArgProcedureCode: procedureCode,
	})

	/*
		return gmm.GmmFSM.SendEvent(ue.State[accessType], gmm.GmmMessageEvent, fsm.ArgsType{
			gmm.ArgAmfUe:         ue,
			gmm.ArgAccessType:    accessType,
			gmm.ArgNASMessage:    msg.GmmMessage,
			gmm.ArgProcedureCode: procedureCode,
		})
	*/
}
