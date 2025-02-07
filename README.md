# Module fake-slam 

Provide a description of the purpose of the module and any relevant information.

## Model cjnj193:fake-slam:fake

This model provides a fake slam that cycles through some maps, for simple testing.

### Configuration

The config for fake slam is trivial, and currently has no attributes:

```json
{}
```

#### Attributes

The following attributes are available for this model:

| Name          | Type   | Inclusion | Description                |
|---------------|--------|-----------|----------------------------|
| `attribute_1` | float  | Required  | Description of attribute 1 |
| `attribute_2` | string | Optional  | Description of attribute 2 |

#### Example Configuration

```json
{}
```

### DoCommand

If your model implements DoCommand, provide an example payload of each command that is supported and the arguments that can be used. If your model does not implement DoCommand, remove this section.

#### Example DoCommand

```json
{
  "command_name": {
    "arg1": "foo",
    "arg2": 1
  }
}
```
