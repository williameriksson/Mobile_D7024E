package d7024e

import (
  "Mobile_D7024E/common"
  "time"
)

const (
  REPUBLISHINTERVAL time.Duration = time.Duration(60)*time.Second
  PURGEINTERVAL time.Duration = time.Duration(10)*time.Second
)

type PurgeInformation struct {
  Key             string
  Pinned          bool
  PurgeTimeStamp  time.Time
}

type DataInformation struct {
  MyKeys       []string
  PurgeInfos   []PurgeInformation
}

func (kademlia *Kademlia) RepublishData() {
  for _, myKey := range kademlia.Datainfo.MyKeys {
    kademlia.PublishData(myKey, kademlia.files[myKey])
  }
  time.AfterFunc(REPUBLISHINTERVAL, kademlia.RepublishData)
}


// Should periodically call itself (needs testing), could be changed to trigger on next event
// if sorting mechanism is implemented
func (kademlia *Kademlia) PurgeData() {

  for _, purgeInfo := range kademlia.Datainfo.PurgeInfos {
    if !purgeInfo.Pinned && time.Now().After(purgeInfo.PurgeTimeStamp){
      select {
      case kademlia.serverChannel <- common.NewHandle(common.CMD_REMOVE_FILE, purgeInfo.Key, ""):
      default:
      }
      delete(kademlia.files, purgeInfo.Key)
    }
  }
  time.AfterFunc(PURGEINTERVAL, kademlia.PurgeData)

}

func (kademlia *Kademlia) SetPurgeStamp(purgeInfo *PurgeInformation) {
  duration := time.Duration(60)*time.Second
  purgeInfo.PurgeTimeStamp = time.Now().Add(duration)
}
