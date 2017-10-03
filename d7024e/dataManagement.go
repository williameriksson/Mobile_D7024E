package d7024e

import (
  "time"
)

const (
  REPUBLISHINTERVAL time.Duration = time.Duration(60)*time.Second
  PURGEINTERVAL time.Duration = time.Duration(10)*time.Second
)

type PurgeInformation struct {
  FileName        string
  Pinned          bool
  PurgeTimeStamp  time.Time
}

type DataInformation struct {
  MyFileNames  []string
  PurgeInfos   []PurgeInformation
}

func (kademlia *Kademlia) RepublishData() {
  for _, myFile := range kademlia.Datainfo.MyFileNames {
    kademlia.PublishData([]byte(myFile))
  }
  time.AfterFunc(REPUBLISHINTERVAL, kademlia.RepublishData)
}


// Should periodically call itself (needs testing), could be changed to trigger on next event
// if sorting mechanism is implemented
func (kademlia *Kademlia) PurgeData() {
  
  for _, purgeInfo := range kademlia.Datainfo.PurgeInfos {
    if !purgeInfo.Pinned && time.Now().After(purgeInfo.PurgeTimeStamp){
      // TODO: Add functionality to remove the actual file also
      delete(kademlia.files, HashData([]byte(purgeInfo.FileName)))
    }
  }
  time.AfterFunc(PURGEINTERVAL, kademlia.PurgeData)

}

func (kademlia *Kademlia) SetPurgeStamp(purgeInfo *PurgeInformation) {
  duration := time.Duration(60)*time.Second
  purgeInfo.PurgeTimeStamp = time.Now().Add(duration)
}
