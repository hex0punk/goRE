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