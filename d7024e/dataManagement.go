package d7024e

import (
  "time"
  "fmt"
)

const (
  TTL time.Duration = time.Duration(40)*time.Second
  REPUBLISH_MY_FILES_INTERVAL time.Duration = time.Duration(20)*time.Second
  REPUBLISH_INTERVAL time.Duration = time.Duration(20)*time.Second
  PURGE_INTERVAL time.Duration = time.Duration(10)*time.Second
)

type PurgeInformation struct {
  Key             string
  Pinned          bool
  PurgeTimeStamp  time.Time
  TimeToLive      time.Duration
}

type DataInformation struct {
  MyKeys       []string
  PurgeInfos   map[string]PurgeInformation
}

func (kademlia *Kademlia) RepublishMyData() {
  for _, myKey := range kademlia.Datainfo.MyKeys {
    tmp := kademlia.Datainfo.PurgeInfos[myKey]
    tmp.TimeToLive = TTL
    kademlia.Datainfo.PurgeInfos[myKey] = tmp
    kademlia.PublishData(kademlia.Datainfo.PurgeInfos[myKey], kademlia.files[myKey])
  }
  time.AfterFunc(REPUBLISH_MY_FILES_INTERVAL, kademlia.RepublishMyData)
}

func (kademlia *Kademlia) RepublishData() {
  for key, purgeInfo := range kademlia.Datainfo.PurgeInfos {
    purgeInfo.TimeToLive = purgeInfo.PurgeTimeStamp.Sub(time.Now())
    kademlia.PublishData(purgeInfo, kademlia.files[key])
  }
  time.AfterFunc(REPUBLISH_INTERVAL, kademlia.RepublishData)
}


// Should periodically call itself (needs testing), could be changed to trigger on next event
// if sorting mechanism is implemented
func (kademlia *Kademlia) PurgeData() {

  for _, purgeInfo := range kademlia.Datainfo.PurgeInfos {
    if !purgeInfo.Pinned && time.Now().After(purgeInfo.PurgeTimeStamp){
      select {
      case kademlia.ServerChannel <- NewHandle(CMD_REMOVE_FILE, purgeInfo, kademlia.files[purgeInfo.Key]):
      case <-time.After(time.Millisecond * 50):
        fmt.Println("Could not purge the data, handler did not read the channel")
      }
      delete(kademlia.files, purgeInfo.Key)
      delete(kademlia.Datainfo.PurgeInfos, purgeInfo.Key)
    }
  }
  time.AfterFunc(PURGE_INTERVAL, kademlia.PurgeData)

}

func (kademlia *Kademlia) SetPurgeStamp(purgeInfo *PurgeInformation) {
  purgeInfo.PurgeTimeStamp = time.Now().Add(purgeInfo.TimeToLive)
}
