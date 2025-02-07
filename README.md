# Module fake-slam 

Provide a description of the purpose of the module and any relevant information.

## Model cjnj193:fake-slam:fake

This model provides a fake slam that cycles through some maps, for simple testing.

### Configuration

No attributes are required. 'calls_till_next_map' controls how many PointCloudMap calls need to happen before the map updates and 'is_localizing' sets the mapping mode to localizing.

```json
{
  "calls_till_next_map": <int>,
  "is_localizing": <bool>,
}
```

#### Attributes

The following attributes are available for this model:

| Name          | Type   | Inclusion | Description                |
|---------------|--------|-----------|----------------------------|
| `calls_till_next_map` | int  | Optional  | how many calls to PointCloudMap before fake slam switches to the next map. Default 5 |
| `is_localizing` | bool  | Optional  | Tell fake slam to return in localizing mode or mapping mode, which can effect how the ui and motion service treat the service. Default false |

#### Example Configuration

```json
{
  "calls_till_next_map": 1,
  "is_localizing": false
}
```
