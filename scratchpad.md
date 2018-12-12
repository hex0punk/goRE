# TODO
- Add "Registry" struct for every module to describe itself
- Add get registry 
- Add an object for module
- Make sure  to include path for file
- Module Option should include the Flag property as shown for merlin exploit framework

- Impletnt showOptions as shown here: https://github.com/Ne0nd0g/merlin/blob/master/pkg/modules/modules.go
- From the above, add show and set options as well, once reayd to add options to modules


```
	patterns[0] = &gcdapi.NetworkRequestPattern{
		UrlPattern: "*" + s.Options.Scope + "/*",
		ResourceType: "Document",
		InterceptionStage: "HeadersReceived",
	}
	patterns[1] = &gcdapi.NetworkRequestPattern{
		UrlPattern:        "*" + s.Options.Scope + "*.js",
		ResourceType:      "Script",
		InterceptionStage: "HeadersReceived",
	}

	s.Target.Subscribe("Network.requestWillBeSent", func(target *gcd.ChromeTarget, v []byte) {
		msg := &gcdapi.NetworkRequestInterceptedEvent{}
		err := json.Unmarshal(v, msg)
		iid := msg.Params.InterceptionId
		//rtype := msg.Params.ResourceType
		reason := msg.Params.ResponseErrorReason
		if err != nil {
			log.Fatalf("error unmarshalling event data: %v\n", err)
		}
		log.Println(msg)
		s.Target.Network.ContinueInterceptedRequest(iid, reason, "", "", "", "", nil, nil)
	})
```
