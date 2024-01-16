# Installation Steps
```
   go get stash.bms.bz/merchandise/go-logger
```

### Usage Example

#### Initialize logger
```
log := logger.New("test")
```

#### Set Masking
```
log := logger.New("test")
pattern := map[string]string{
		`(^|\D)([6-9]\d{9})(\D|$)`:                         `"****"`,
		`([a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+)`: "@#$%",
	}

log.EnableMask(pattern)
log.Debug("TEST", map[string]interface{}{"key": "9774007164",
		"text": "My no. 9774007164 is on waybill 777919076389908 and 967400716499087655467111 and shourie.dbz2007@gmail.com"})
```

#### Logger types
```
log.Info("starting server", map[string]interface{}{"port": "80"})
log.Debug("order payload", map[string]interface{}{"id": "1"})
log.Warn("deprecated", map[string]interface{}{"warn": "will be removed in future releases"})
log.Error("CODE", "description", "priority 1", map[string]interface{}{"caller":"func","file":"func.go","line":1,"stackTrace": "trace","errorString": "shown outside"}) {

```

**Note:** If stash.bms.bz/merchandise/go-errors is used as the error handling lib you wont need to pass the 4th param(error source) in the log.Error