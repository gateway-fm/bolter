logger{
  ## logger type:
  ## 0: creating file in project root directory
  ## 1: logging using Logrus logger
  ## 2: logging using Zap logger

  logger_type = 0
  file_name = "serega.txt"
}
vegeta {
  url = "http://127.0.0.1:10000" ## e.g. "http://localhost:8181/" MUST STARTS WITH "http://" !!!
  method = "POST"
  is_public = false ## Does the auth header needed?
  rate = 500
  duration = 2 ## must be int!!! no "1s" "time.Second" etc
  header {
    auth = "Bearer" ## Auth type if it's needed. Now only (uppercase!) "Bearer" is availiable
    bear = "" ## bearer token
  }
}

requests "eth_getBlockByNumber"{
  ## example of "eth_gasPrice" request
  request  {
    jsonrpc = "2.0"
    method = "eth_getBlockByNumber"
    params = ["latest", "true"]
    id = "0x1234"
    hard_coded = true ## must be true ONLY for eth_getBlockByNumber !
  }
}
