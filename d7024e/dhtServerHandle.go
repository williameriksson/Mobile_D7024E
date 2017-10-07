package d7024e

const (
  CMD_FOUND_FILE = "FOUND_FILE"
  CMD_RETRIEVE_FILE = "RETRIEVE_FILE"
  CMD_REMOVE_FILE = "REMOVE_FILE"
)

type Handle struct {
  Command string
  PurgeInfo  PurgeInformation
  Ip    string
}

func NewHandle(cmd string, purgeInfo PurgeInformation, ip string) Handle {
  return Handle{cmd, purgeInfo, ip}
}
