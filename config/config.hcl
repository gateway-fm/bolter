logger{
  ## logger type:
  ## 0: creating file in project root directory
  ## 1: logging using Logrus logger
  ## 2: logging using Zap logger

  logger_type = 0
  file_name = "serega.txt"
}
vegeta {
  url = "http://127.0.0.1:10000/v4/ethereum/non-archival/mainnet" ## e.g. "http://localhost:8181/"
  method = "POST"
  is_public = false ## Does the auth header needed?
  rate = 10
  duration = 5 ## must be int!!! no "1s" "time.Second" etc
  header {
    auth = "Bearer" ## Auth type if it's needed. Now only (uppercase!) "Bearer" is availiable
    bear = "sSRspJoNahtpyS29LtVt5XuJJV1Mzza5.FatUKsu6MBQXHEzI" ## bearer token
  }
}

requests "eth_getBlockByHash" {
  request {
    jsonrpc = "2.0"
    method  = "eth_gasPrice"
    params  = []
    id      = "1"
  }
}
#requests "eth_getBlockByNumber"{
  ## example of "eth_blockNumber" request
 # request  {
  #  jsonrpc = "2.0"
   # method = "eth_gasPrice"
    #params = []
    #id = "0x1234"
  #}
#}
