# Build
```
make
```

# Usage
```
viam-data-capture-util some_capture_filet.capture | jq '.metadata'
{
  "part_id": "unknown",
  "component_type": "rdk:component:sensor",
  "component_name": "boom",
  "method_name": "Readings",
  "type": 2,
  "file_name": "/Users/nicksanford/rdk/capture.capture",
  "method_parameters": {
    "some": {
      "type_url": "type.googleapis.com/google.protobuf.StringValue",
      "value": "CgZwYXJhbXM="
    }
  },
  "file_extension": ".dat"
}
```
