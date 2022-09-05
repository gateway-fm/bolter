logger{
  ## logger type:
  ## 0: creating file in project root directory
  ## 1: logging using Logrus logger
  ## 2: logging using Zap logger

  logger_type = 0
  file_name = "serega.txt"
}
vegeta {
  url = "" ## e.g. "http://localhost:8181/"
  method = "POST"
  is_public = true ## Does the auth header needed?
  rate = 10
  duration = 1 ## must be int!!! no "1s" "time.Second" etc
  header {
    auth = "" ## Auth type if it's needed. Now only "bearer" is availiable
    bear = "" ## bearer token
  }
}

## example of "eth_blockNumber" request
request {
  jsonrpc = "2.0"
  method = "eth_blockNumber"
  params = []
  id = "0x1234"
}
##


