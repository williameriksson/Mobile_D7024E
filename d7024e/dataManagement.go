package d7024e

import (
  "time"
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

}


// Should periodically call itself (needs testing), could be changed to trigger on next event
// if sorting mechanism is implemented
func (kademlia *Kademlia) PurgeData() {

  for _, purgeInfo := range kademlia.Datainfo.PurgeInfos {
    if !purgeInfo.Pinned && time.Now().After(purgeInfo.PurgeTimeStamp){
      // TODO: Add functionality to remove the actual file also
      delete(kademlia.files, purgeInfo.Key)
    }
  }
  time.AfterFunc(time.Duration(10)*time.Second, kademlia.PurgeData)

}

func (kademlia *Kademlia) SetPurgeStamp(purgeInfo PurgeInformation) {
  duration := time.Duration(60)*time.Second
  purgeInfo.PurgeTimeStamp = time.Now().Add(duration)
}
