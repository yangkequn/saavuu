## saavuu, the most concise, redis based web-server framework
    the name saavuu borrow from "杀悟",means kill bad wisdom。 I hate bad tools.
### main features
* All HTTP requests are transferd as binary msgpack data. It's compact and fast.
* No API version related problem. Just use redis api.
* Use msgpack to support structure data by default. Easily to upgrade data sturecture.
* Use no database but KEYDB which is redis compatible. Flash storage supportion brings both memory speed and whole disk capacity
* You don't need to write any GET Logic. Just use redis to query.
* You don't need to write any DELETE Logic. Just use redis to remove.
* You can focus only and alway on POST & PUT logic. 
    saavuu will put request data to redis queue, and the service listening the queue will process the data.
* You can use any programming language you like. python or go or not.
* redis pipeline  brings heavy batch process performance.  
### other features
* specify Content-Type in web side
* allow specify response fields in web client to reduce web traffic
* support JWT for authorization
* fully access control
* support CORS
  
## demo usage
### server, go example:
```
package main

import (
	"github.com/yangkequn/saavuu"
	"github.com/yangkequn/saavuu/api"
	"github.com/yangkequn/saavuu/data"
)

type Input struct {
	Data   []uint16
	Id   string `msgpack:"alias:JWT_id"`
}
var apiDemo=api.New("demo")
func init() {
	apiDemo.Serve(func(parmIn map[string]interface{}) (ret map[string]interface{}, err error) {
		var req *Input = &Input{}
		if err = data.MapsToStructure(parmIn, req); err != nil {
			return nil, err
		} else if req.Id == "" || len(req.Data) == 0 {
			return nil, saavuu.ErrInvalidInput
		}
		// your logic here
		return map[string]interface{}{"data": req.Data}, nil
	})
}
```

### server, python example:
```
class service_textToMp3(Service):
    def __init__(self):
        Service.__init__(self,"redis.vm:6379/0")
    def check_input_item(self, i):
            if "BackTo" not in i:
                return False
            return True
    def process(self,items):
        #your logic here
        for i in items:
            self.send_back(i,{"Result":input.value})
service_textToMp3().start()
```

### web client, javascript /typescript example:
```
HGET("UserInfo", id).then((data) => {
    //your logic here
})
```


## about configuration 
    configuration is read from enviroment variables. Make sure enviroment variables are added to your IDE (launch.json for vs code) or docker. 
    these are the default example:
```
    "REDIS_ADDR_PARAM": "127.0.0.1:6379",
    "REDIS_PASSWORD_PARAM": "",
    "REDIS_DB_PARAM": "0",
    "REDIS_ADDR_DATA": "127.0.0.1:6379",
    "REDIS_PASSWORD_DATA": "",
    "REDIS_DB_DATA": "0",
    "JWT_SECRET": "WyBJujUQzWg4YiQqLe9N36DA/7QqZcOkg2o=",
    "JWT_IGNORE_FIELDS": "iat,exp,nbf,iss,aud,sub,typ,azp,nonce,auth_time,acr,amr,at_hash,c_hash,updated_at,nonce,auth_time,acr,amr,at_hash,c_hash,updated_at",
    "CORS": "*",
    "MAX_BUFFER_SIZE": "3145728",
    "DEVELOP_MODE": "true",
    "APP_MODE":"framework",
```