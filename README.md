# go-metacontext

In a micro-service environment, a request to one service will often result in sub-requests being made to other services. In these sub-requests you may want to pass along some contextual metadata so that the next service does not have to look it up again.  
Metacontext allows you to pass along a JSON map of metadata without changing the API contract. It is a small wrapper that simplifies the manipulation of JSON bodies to include a metadata section that can be parsed out by the receiving service.  
See [the example](_example/main.go) for usage.

The idea behind this is that it would be used here and there in places where something quick and dirty is all that's needed. A more full fledged solution should be considered if you are finding that this functionality is needed in a lot of places.
